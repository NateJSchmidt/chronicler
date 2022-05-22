package main

import (
	"database/sql"
	"log"

	"github.com/NateJSchmidt/chronicler/internal/pkg/db"
	"github.com/NateJSchmidt/chronicler/internal/pkg/db/postgres"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func getLoggerConfiguration() zap.Config {
	loggerConfig := zap.NewDevelopmentConfig() // zap.NewProductionConfig()

	// set the output to stdout and stderr
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}
	// set the timestamp key to "timestamp"
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	// set the timestamp encoder to RFC3339 (ex: 2022-05-22T09:37:48-06:00)
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// set the level to debug for development
	loggerConfig.Level.SetLevel(zap.DebugLevel)

	return loggerConfig
}

// main entrypoint for the program
func main() {
	// initialize the logger
	loggerConfig := getLoggerConfiguration()
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Debug("Zap logger initialized")

	// create necessary resources
	r := gin.Default()
	// dbName := "testdb"

	// perform configuration for REST server
	configureRoutes(r)

	// perform configuration for database
	// createDatabase(dbName)
	// createTables(dbName)

	var pg db.DBInterface = postgres.Init(sugar, "chronicler_user", "password1234", "localhost", "5432", "chronicler", postgres.Disable)
	pg.WriteEvent("noun:verb-1234", "{\"hurray\": \"payload\"}", uuid.New().String(), "1.0.0")

	// start the server
	r.Run()
}
