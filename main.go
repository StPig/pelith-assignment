package main

import (
	"fmt"
	"os"
	"pelith-assignment/api"
	"pelith-assignment/database"
)

func main() {
	// 初始化 DB
	fmt.Println("Init database")
	db, err := database.InitDB()
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	// 初始化 server
	fmt.Println("Init api server")
	api.InitAPIService()
}
