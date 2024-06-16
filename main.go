package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Customer struct {
	ID                    int     `json:"id"`
	UUID                  string  `json:"uuid"`
	FirstName             string  `json:"firstname"`
	LastName              string  `json:"lastname"`
	Address               string  `json:"address"`
	PhoneNumber           string  `json:"phoneNumber"`
	Email                 string  `json:"email"`
	Employer              string  `json:"employer"`
	AnnualIncome          float32 `json:"annualIncome"`
	RequestedCreditAmount float32 `json:"requestedCreditAmount"`
	AdditionalInformation string  `json:"additionalInformation"`
}

var DATABASE_URL string
var conn *pgxpool.Pool

func main() {
	godotenv.Load()
	DATABASE_URL = os.Getenv("DATABASE_URL")
	user := os.Getenv("username")
	pass := os.Getenv("password")

	conn = createDatabaseConnection(DATABASE_URL)
	defer conn.Close()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339} ${status} ${method} ${host}${path} ${latency_human}]` + "\n",
	}))
	e.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.POST("/newcredit", addCustomer)
	e.POST("/uploadPDF", uploadPDF)

	initTempDir()

	e.Logger.Fatal(e.Start(":8000"))
}

func createDatabaseConnection(url string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
