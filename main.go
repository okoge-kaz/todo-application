package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/okoge-kaz/todo-application/db"
	"github.com/okoge-kaz/todo-application/router"
)

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize engine
	engine := router.Init()

	// start server
	const port = 8000
	engine.Run(fmt.Sprintf(":%d", port))
}
