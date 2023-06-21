package main

import (
	"bytes"
	"database/sql"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetLeaveData(t *testing.T) {

	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=manish dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.GET("/getleave", getLeaveData)

	req, err := http.NewRequest("GET", "/getleave", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

}

func TestSaveLeaveData(t *testing.T) {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=manish dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)

    writer.WriteField("name", "Bipul Dubey")
    writer.WriteField("leaveType", "Casual Leave")
    writer.WriteField("fromDate", "2023-12-05")
    writer.WriteField("toDate", "2023-12-20")
    writer.WriteField("team", "AnalyticOps")
    writer.WriteField("reporter", "Surya Kant")
    writer.Close()

    req, err := http.NewRequest("POST", "/postleave", body)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    saveLeaveData(c)

    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSaveLeaveWithFileData(t *testing.T) {

	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=manish dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    fileWriter, err := writer.CreateFormFile("file", "sick_leave.pdf")
    if err != nil {
        t.Fatal(err)
    }
    fileContents := []byte("test file contents")
    _, err = fileWriter.Write(fileContents)
    if err != nil {
        t.Fatal(err)
    }
    writer.WriteField("name", "Bipul")
    writer.WriteField("leaveType", "Sick Leave")
    writer.WriteField("fromDate", "2022-01-01")
    writer.WriteField("toDate", "2022-12-31")
    writer.WriteField("team", "DataOps")
    writer.WriteField("reporter", "Surya Kant")

    writer.Close()

    req, err := http.NewRequest("POST", "/postleave", body)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req

    // Call the function being tested
    saveLeaveData(c)
    assert.Equal(t, http.StatusCreated, w.Code)
}
