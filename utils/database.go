package utils

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

// DBConfig holds database configuration
type DBConfig struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
	Trusted  bool
}

// ConnectDB establishes a connection to the MS SQL Server database
func ConnectDB(config DBConfig) (*sql.DB, error) {
	server := strings.ReplaceAll(config.Server, "\\", "\\\\")
	connString := fmt.Sprintf("server=%s;port=%d; user id=%s;password=%s;database=%s;encrypt=disable",
		server, config.Port, config.User, config.Password, config.Database)
	fmt.Println(connString)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}
