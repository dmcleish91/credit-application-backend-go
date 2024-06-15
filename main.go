package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	PhoneNumber           string  `json:"phonenumber"`
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

func addCustomer(ctx echo.Context) error {
	customer := Customer{}

	defer ctx.Request().Body.Close()
	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return ctx.String(http.StatusInternalServerError, "Error parsing data")
	}

	err = json.Unmarshal(b, &customer)
	if err != nil {
		log.Printf("Failed unmarshaling: %s", err)
		return ctx.String(http.StatusInternalServerError, "Some error occured")
	}
	result, err := insertCustomer(customer)
	if err != nil {
		log.Printf("Failed insert into database: %s", err)
		return ctx.String(http.StatusInternalServerError, "Some error occured")
	}

	if result == 1 {
		return ctx.String(http.StatusOK, "Customer data accepted")
	} else {
		return ctx.String(http.StatusInternalServerError, "Some error occured")
	}
}

func insertCustomer(customer Customer) (int64, error) {

	result, err := conn.Exec(context.Background(), "INSERT INTO public.customers (uuid, first_name, last_name, address, phone_number, email, employer, "+
		"annual_income, request_credit_amount, additional_information) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);", customer.UUID, customer.FirstName,
		customer.LastName, customer.Address, customer.PhoneNumber, customer.Email, customer.Employer, customer.AnnualIncome, customer.RequestedCreditAmount, customer.AdditionalInformation)
	if err != nil {
		conn.Close()
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

// func fetchAllCustomers(conn *pgxpool.Pool) {
// 	var customers []Customer
// 	rows, err := conn.Query(context.Background(), "SELECT * FROM customers")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 		os.Exit(1)
// 	}

// 	for rows.Next() {
// 		var customer Customer
// 		err := rows.Scan(&customer.ID, &customer.UUID, &customer.FirstName, &customer.LastName, &customer.Address,
// 			&customer.PhoneNumber, &customer.Email, &customer.Employer, &customer.AnnualIncome, &customer.RequestedCreditAmount, &customer.AdditionalInformation)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		customers = append(customers, customer)
// 	}

// 	fmt.Println(customers)
// }

func uploadPDF(c echo.Context) error {
	// Source
	file, err := c.FormFile("pdf")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Check if the uploaded file is a PDF
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.String(http.StatusBadRequest, "The uploaded file is not a PDF")
	}

	// Destination
	dst, err := os.Create(fmt.Sprintf("upload/%s", file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}
