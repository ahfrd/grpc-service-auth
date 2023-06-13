package db

import (
	"database/sql"
	"fmt"

	"github.com/ahfrd/grpc/auth-client/config"
	_ "github.com/go-sql-driver/mysql"
)

// Database is a
type Database struct{}

// ConnectDB is a
func (o Database) ConnectDB() (*sql.DB, error) {
	c, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed connection to DB : %v", err)

	}
	DB := c.DB
	URI := c.DBUrl
	fmt.Println("...")
	fmt.Println(DB)
	fmt.Println(URI)
	// db, err := apmsql.Open(DB, URI)
	db, err := sql.Open(DB, URI)

	if err != nil {
		return nil, fmt.Errorf("failed connection to DB : %v", err)
	}
	fmt.Println(db)

	return db, nil
}
