package database

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
)

// ProxyExists checks if a proxy exists.
func (database *Database) ProxyExists(proxyType, proxyHost string, proxyPort int) bool {
	queryString := fmt.Sprintf(
		"SELECT id FROM %s WHERE type = ? AND host = ? AND port = ?",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return false
	}

	rows, err := query.Query(proxyType, proxyHost, proxyPort)
	if err != nil {
		logger.Errorf("database", "error occurred while querying: %s", err.Error())
		return false
	}

	defer rows.Close()

	return rows.Next()
}

// ProxyCreate creates a new proxy.
func (database *Database) ProxyCreate(proxyType, proxyHost string, proxyPort int) bool {
	queryString := fmt.Sprintf(
		"INSERT INTO %s (type, host, port) VALUES (?, ?, ?)",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(proxyType, proxyHost, proxyPort); err != nil {
		logger.Errorf("database", "error occurred while inserting: %s", err.Error())
		return false
	}

	return true
}

// ProxyDelete deletes a proxy.
func (database *Database) ProxyDelete(proxyType, proxyHost string, proxyPort int) bool {
	queryString := fmt.Sprintf(
		"DELETE FROM %s WHERE type = ? AND host = ? AND port = ?",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(proxyType, proxyHost, proxyPort); err != nil {
		logger.Errorf("database", "error occurred while deleting: %s", err.Error())
		return false
	}

	return true
}

// ProxyEnableAuthentication enables authentication for a proxy.
func (database *Database) ProxyEnableAuthentication(proxyType, proxyHost string, proxyPort int, username, password string) bool {
	queryString := fmt.Sprintf(
		"UPDATE %s SET username = ?, password = ? WHERE type = ? AND host = ? AND port = ?",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(username, password, proxyType, proxyHost, proxyPort); err != nil {
		logger.Errorf("database", "error occurred while updating: %s", err.Error())
		return false
	}

	return true
}

// ProxyDisableAuthentication disables authentication for a proxy.
func (database *Database) ProxyDisableAuthentication(proxyType, proxyHost string, proxyPort int) bool {
	queryString := fmt.Sprintf(
		"UPDATE %s SET username = NULL, password = NULL WHERE type = ? AND host = ? AND port = ?",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return false
	}

	if _, err := query.Exec(proxyType, proxyHost, proxyPort); err != nil {
		logger.Errorf("database", "error occurred while updating: %s", err.Error())
		return false
	}

	return true
}

// ProxyServer represents a proxy server configuration.
type ProxyServer struct {
	Type     string
	Host     string
	Port     int
	Username *string
	Password *string
}

// ProxyFindByHost finds all proxies by host.
func (database *Database) ProxyFindByHost(proxyHost string) []*ProxyServer {
	queryString := fmt.Sprintf(
		"SELECT type, host, port, username, password FROM %s WHERE host = ?",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return nil
	}

	rows, err := query.Query(proxyHost)
	if err != nil {
		logger.Errorf("database", "error occurred while querying: %s", err.Error())
		return nil
	}

	defer rows.Close()

	var servers []*ProxyServer

	for rows.Next() {
		var server ProxyServer

		if err := rows.Scan(&server.Type, &server.Host, &server.Port, &server.Username, &server.Password); err != nil {
			logger.Errorf("database", "error occurred while scanning: %s", err.Error())
			return nil
		}

		servers = append(servers, &server)
	}

	return servers
}

// ProxyList lists all proxies.
func (database *Database) ProxyList() []string {
	queryString := fmt.Sprintf(
		"SELECT type, host, port FROM %s",
		database.GetProxiesTableName(),
	)

	query, err := database.GetConnection().Prepare(queryString)
	if err != nil {
		logger.Errorf("database", "error occurred while preparing query: %s", err.Error())
		return nil
	}

	rows, err := query.Query()
	if err != nil {
		logger.Errorf("database", "error occurred while querying: %s", err.Error())
		return nil
	}

	defer rows.Close()

	var servers []string

	for rows.Next() {
		var server ProxyServer

		if err := rows.Scan(&server.Type, &server.Host, &server.Port); err != nil {
			logger.Errorf("database", "error occurred while scanning: %s", err.Error())
			return nil
		}

		servers = append(servers, fmt.Sprintf("%s %s:%d", server.Type, server.Host, server.Port))

	}

	return servers
}
