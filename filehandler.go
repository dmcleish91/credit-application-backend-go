package main

import (
	"log"
	"os"
	"path/filepath"
)

func initTempDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	path := "uploads"
	currentPath := filepath.Join(cwd, path)
	err = os.Mkdir(currentPath, os.ModePerm)
	if err != nil {
		log.Println("Directory probably already exists", err)
	} else {
		log.Println("initialized", path)
	}
}

// function to create a new path in the same directory as the running go executable
func createDirectory(folder string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	currentPath := filepath.Join(cwd, "uploads", folder)
	log.Println(currentPath)
	err = os.Mkdir(currentPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Directory created successfully at", currentPath)
	}
}
