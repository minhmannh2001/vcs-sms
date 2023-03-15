package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const (
	// CREATE - operation to create database
	CREATE = "CREATE"
	// DROP - operation to drop database
	DROP = "DROP"
)

// MockedDB is used in unit tests to mock db
func MockedDB(operation string) error {
	// Load .env file
	err := godotenv.Load("./../.env")
	if err != nil {
		godotenv.Load()
		panic("Failed to load env file")
	}

	dbName := os.Getenv("MOCK_DBNAME")
	pgUser := os.Getenv("DB_USERNAME")
	pgPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// createdb => https://www.postgresql.org/docs/7.0/app-createdb.htm
	// dropdb => https://www.postgresql.org/docs/7.0/app-dropdb.htm
	var command string

	if operation == CREATE {
		command = "createdb"
	} else {
		command = "dropdb"
	}

	// createdb & dropdb commands have same configuration syntax.
	cmd := exec.Command(command, "-h", dbHost, "-p", dbPort, "-U", pgUser, "-e", dbName)
	cmd.Env = os.Environ()

	/*
	   if we normally execute createdb/dropdb, we will be propmted to provide password.
	   To inject password automatically, we have to set PGPASSWORD as prefix.
	*/
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", pgPassword))

	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v on %v.\n%v", command, dbName, err)
		return err
	}

	return nil
	/*
	   Alternatively instead of createdb/dropdb, you can use
	   psql -c "CREATE/DROP DATABASE DBNAME" "DATABASE_URL"
	*/
}
