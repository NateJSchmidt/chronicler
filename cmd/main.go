package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// configureRoutes sets up the routes for the REST server
func configureRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
}

// createDatabase creates a database named dbName
func createDatabase(dbName string) {
	db, err := sql.Open("postgres", "postgres://postgres:password1234@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the database
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Fatal(err)
	}
}

// createTables creates the tables in the database specified by dbName
func createTables(dbName string) {
	db, err := sql.Open("postgres", "postgres://postgres:password1234@localhost:5432/"+dbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table
	sqlStatement := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
}

// main entrypoint for the program
func main() {
	// create necessary resources
	r := gin.Default()
	dbName := "testdb"

	// perform configuration for REST server
	configureRoutes(r)

	// perform configuration for database
	createDatabase(dbName)
	createTables(dbName)

	// start the server
	r.Run()
}
