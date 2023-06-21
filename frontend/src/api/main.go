package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// the structure of the leave form data
type Leave struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LeaveType string `json:"leaveType"`
	FromDate  string `json:"fromDate"`
	ToDate    string `json:"toDate"`
	Team      string `json:"team"`
	FilePath  string `json:"filePath"`
	Reporter  string `json:"reporter"`
}

var db *sql.DB

func main() {

	// creating database connection
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=manish dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/getleave", getLeaveData)
	router.POST("/postleave", saveLeaveData)
	router.GET("/leaveTypes", handleLeaveTypesAPI)
	router.GET("/file/:filepath", serveFile)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}

func getLeaveData(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM leaves")
	if err != nil {
		fmt.Println("Error executing SQL query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var leaves []Leave

	for rows.Next() {
		var leave Leave
		var toDate, fromDate time.Time
		err := rows.Scan(&leave.ID, &leave.Name, &leave.LeaveType, &fromDate, &toDate, &leave.Team, &leave.FilePath, &leave.Reporter)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows found")
			} else {
				if strings.Contains(err.Error(), "converting NULL to string") {
					leave.FilePath = ""
				} else {
					fmt.Println("Error scanning row:", err)
					c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Server Error"})
					return
				}
			}
		}

		leave.FromDate = fromDate.Format("2006-01-02")
		leave.ToDate = toDate.Format("2006-01-02")

		leaves = append(leaves, leave)
	}

	fmt.Println("Leave data retrieved successfully:", leaves)
	c.JSON(http.StatusOK, leaves)
}

func saveLeaveData(c *gin.Context) {
	var leave Leave

	err := c.Request.ParseMultipartForm(0)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave.Name = c.Request.FormValue("name")
	leave.LeaveType = c.Request.FormValue("leaveType")
	leave.FromDate = c.Request.FormValue("fromDate")
	leave.ToDate = c.Request.FormValue("toDate")
	leave.Team = c.Request.FormValue("team")
	leave.Reporter = c.Request.FormValue("reporter")

	// checking condition for sick leave
	if leave.LeaveType == "Sick Leave" {
		filePath, err := saveUploadedFile(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		leave.FilePath = filePath
	}

	_, err = db.Exec("INSERT INTO leaves (name, leave_type, from_date, to_date, team, filepath, reporter) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		leave.Name, leave.LeaveType, leave.FromDate, leave.ToDate, leave.Team, leave.FilePath, leave.Reporter)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mesage":"Success"})
}

func saveUploadedFile(c *gin.Context) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	err = c.SaveUploadedFile(file, "./public/"+file.Filename)
	if err != nil {
		return "", err
	}

	return file.Filename, nil
}

func handleLeaveTypesAPI(c *gin.Context) {
	leaveTypes := []string{"Casual Leave", "Earned Leave", "Sick Leave"}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"leaveTypes": leaveTypes})
}

func serveFile(c *gin.Context) {
	filePath := c.Param("filepath")
	filePath = "./public/" + filePath
	c.File(filePath)
}
