package database

import (
	"database/sql"
	"fmt"
	"github.com/darki73/pac-manager/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database configuration.
type Database struct {
	// storage represents the storage.
	storage *storage.Storage
	// name represents the name of the database.
	name string
	// domainsTableName represents the name of the domains table.
	domainsTableName string
	// proxiesTableName represents the name of the proxies table.
	proxiesTableName string
	// connection represents the connection to the database.
	connection *sql.DB
}

// NewDatabase creates a new database.
func NewDatabase(storage *storage.Storage, name string) *Database {
	return &Database{
		storage:          storage,
		name:             fmt.Sprintf("%s.db", name),
		domainsTableName: "domains",
		proxiesTableName: "proxies",
		connection:       nil,
	}
}

// Start creates the database and opens the connection.
func (database *Database) Start() error {
	if err := database.CreateDatabase(); err != nil {
		return err
	}

	connection, err := sql.Open("sqlite3", database.GetStorage().GetPath(database.GetName()))
	if err != nil {
		return err
	}

	database.connection = connection

	if err := database.CreateTables(); err != nil {
		return err
	}

	return nil
}

// Stop closes the connection to the database.
func (database *Database) Stop() error {
	return database.GetConnection().Close()
}
