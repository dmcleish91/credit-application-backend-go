package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

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

func uploadPDF(c echo.Context) error {
	// Read form fields
	folder := c.FormValue("uuid")

	// Read multipart form file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("%s", err)
	}

	// Check if the uploaded file is a PDF
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.String(http.StatusBadRequest, "The uploaded file is not a PDF")
	}

	// opens multipart form file
	src, err := file.Open()
	if err != nil {
		log.Printf("%s", err)
	}
	defer src.Close()

	// create destination folder in uploads
	createDirectory(folder)

	// creates the file at the target location
	fileLocation := fmt.Sprintf("uploads/%s/%s", folder, file.Filename)
	dst, err := os.Create(fileLocation)
	if err != nil {
		log.Printf("%s", err)
	}
	defer dst.Close()

	// copy the file from memory to the destination
	if _, err = io.Copy(dst, src); err != nil {
		log.Printf("%s", err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}
