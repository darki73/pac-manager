package database

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
)

// DomainExists checks if the domain exists.
func (database *Database) DomainExists(domain string) bool {
	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE domain = ?", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing exists query: %s", err.Error())
		return false
	}

	var count int
	if err := query.QueryRow(domain).Scan(&count); err != nil {
		logger.Errorf("database", "error occurred while checking if domain exists: %s", err.Error())
		return false
	}

	return count > 0
}

// DomainUpdate updates the domain.
func (database *Database) DomainUpdate(domain string, proxify bool) bool {
	if !database.DomainExists(domain) {
		return false
	}

	queryString := fmt.Sprintf("UPDATE %s SET proxify = ? WHERE domain = ?", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing update query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(proxify, domain); err != nil {
		logger.Errorf("database", "error occurred while updating domain: %s", err.Error())
		return false
	}

	return true
}

// DomainCreate creates the domain.
func (database *Database) DomainCreate(domain string, proxify bool) bool {
	if database.DomainExists(domain) {
		return false
	}

	queryString := fmt.Sprintf("INSERT INTO %s (domain, proxify) VALUES (?, ?)", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing create query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(domain, proxify); err != nil {
		logger.Errorf("database", "error occurred while creating domain: %s", err.Error())
		return false
	}

	return true
}

// DomainDelete deletes the domain.
func (database *Database) DomainDelete(domain string) bool {
	if !database.DomainExists(domain) {
		return false
	}

	queryString := fmt.Sprintf("DELETE FROM %s WHERE domain = ?", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing delete query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(domain); err != nil {
		logger.Errorf("database", "error occurred while deleting domain: %s", err.Error())
		return false
	}

	return true
}

// DomainList lists the domains.
func (database *Database) DomainList() map[string]bool {
	queryString := fmt.Sprintf("SELECT domain, proxify FROM %s", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing list query: %s", err.Error())
		return nil
	}

	rows, err := query.Query()
	if err != nil {
		logger.Errorf("database", "error occurred while listing domains: %s", err.Error())
		return nil
	}

	domains := make(map[string]bool)

	for rows.Next() {
		var domain string
		var proxify bool
		if err := rows.Scan(&domain, &proxify); err != nil {
			logger.Errorf("database", "error occurred while scanning domain: %s", err.Error())
			continue
		}

		domains[domain] = proxify
	}

	return domains
}

// DomainGet gets the domain.
func (database *Database) DomainGet(domain string) (int, string, bool) {
	queryString := fmt.Sprintf("SELECT id, domain, proxify FROM %s WHERE domain = ?", database.GetDomainsTableName())
	query, err := database.GetConnection().Prepare(queryString)

	if err != nil {
		logger.Errorf("database", "error occurred while preparing get query: %s", err.Error())
		return 0, "", false
	}

	var id int
	var proxify bool
	if err := query.QueryRow(domain).Scan(&id, &domain, &proxify); err != nil {
		logger.Errorf("database", "error occurred while getting domain: %s", err.Error())
		return 0, "", false
	}

	return id, domain, proxify
}
