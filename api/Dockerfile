FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

ENV DB_USER=postgres
ENV DB_PASSWORD=manish
ENV DB_NAME=postgres
ENV DB_HOST=host.docker.internal
ENV DB_PORT=5432
ENV CLIENT_URL=http://localhost:3000
ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./main"]
