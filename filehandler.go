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

	path := "upload"
	currentPath := filepath.Join(cwd, path)
	err = os.Mkdir(currentPath, os.ModePerm)
	if err != nil {
		log.Println("Directory probably already exists", err)
	} else {
		log.Println("initialized", path)
	}
}

// function to create a new path in the same directory as the running go executable
// func createDirectory(uuid string) {
// 	currentPath := filepath.Join("temp", uuid)
// 	err := os.Mkdir(currentPath, os.ModePerm)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println("Directory created successfully at", currentPath)
// 	}
// }
