package database

import (
	"database/sql"
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/storage"
)

// GetConnection returns the connection to the database.
func (database *Database) GetConnection() *sql.DB {
	return database.connection
}

// GetStorage returns the storage.
func (database *Database) GetStorage() *storage.Storage {
	return database.storage
}

// GetName returns the name of the database.
func (database *Database) GetName() string {
	return database.name
}

// GetDomainsTableName returns the name of the table.
func (database *Database) GetDomainsTableName() string {
	return database.domainsTableName
}

func (database *Database) GetProxiesTableName() string {
	return database.proxiesTableName
}

// CreateDatabase creates the database file.
func (database *Database) CreateDatabase() error {
	if !database.GetStorage().FileExists(database.GetName()) {
		return database.GetStorage().CreateFile(database.GetName())
	}
	return nil
}

// CreateTables creates the tables.
func (database *Database) CreateTables() error {
	if err := database.CreateDomainsTable(); err != nil {
		return err
	}

	if err := database.CreateProxiesTable(); err != nil {
		return err
	}

	return nil
}

// CreateDomainsTable creates the domains table.
func (database *Database) CreateDomainsTable() error {
	if !database.TableExists(database.GetDomainsTableName()) {
		queryString := fmt.Sprintf(
			"CREATE TABLE %s (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, domain TEXT NOT NULL, proxify INTEGER NOT NULL DEFAULT 0, UNIQUE(domain))",
			database.GetDomainsTableName(),
		)

		query, err := database.GetConnection().Prepare(queryString)
		if err != nil {
			logger.Errorf("database", "error occurred while preparing create query: %s", err.Error())
			return err
		}

		if _, err := query.Exec(); err != nil {
			logger.Errorf("database", "error occurred while creating table: %s", err.Error())
			return err
		}
	}
	return nil
}

// CreateProxiesTable creates the proxies table.
func (database *Database) CreateProxiesTable() error {
	if !database.TableExists(database.GetProxiesTableName()) {
		queryString := fmt.Sprintf(
			"CREATE TABLE %s (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, type TEXT NOT NULL, host TEXT NOT NULL, port INTEGER NOT NULL, username TEXT, password TEXT, UNIQUE(type, host, port))",
			database.GetProxiesTableName(),
		)

		query, err := database.GetConnection().Prepare(queryString)
		if err != nil {
			logger.Errorf("database", "error occurred while preparing create query: %s", err.Error())
			return err
		}

		if _, err := query.Exec(); err != nil {
			logger.Errorf("database", "error occurred while creating table: %s", err.Error())
			return err
		}
	}
	return nil
}

// TableExists checks if the table exists.
func (database *Database) TableExists(table string) bool {
	queryString := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", table)

	result := database.GetConnection().QueryRow(queryString)

	if result.Err() != nil {
		logger.Errorf("database", "error occurred while checking if table exists: %s", result.Err().Error())
		return false
	}

	switch err := result.Scan(); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return true
	}
}
