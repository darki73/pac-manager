package server

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/pac"
	"github.com/darki73/pac-manager/pkg/storage"
	"net/http"
	"strings"
)

type Server struct {
	// cfg is the global configuration.
	cfg *config.Config
	// db is the database
	db *database.Database
}

// NewServer creates a new server.
func NewServer(cfg *config.Config, db *database.Database) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

// GetAddress returns the server's IP address.
func (server *Server) GetAddress() string {
	return server.GetConfiguration().GetServer().GetHost()
}

// GetPort returns the server's port.
func (server *Server) GetPort() uint {
	return server.GetConfiguration().GetServer().GetPort()
}

// GetConfiguration returns the configuration.
func (server *Server) GetConfiguration() *config.Config {
	return server.cfg
}

// GetDatabase returns the database.
func (server *Server) GetDatabase() *database.Database {
	return server.db
}

// generatePacFile generates the PAC file.
func (server *Server) generatePacFile() {
	pac.NewPac(
		storage.NewStorage(server.GetConfiguration().GetPac().GetPath(), 0755),
		server.GetConfiguration().GetPac().GetName(),
	).Generate(server.GetDatabase().ProxyList(), server.GetDatabase().DomainList())
}

// indexHandler handles the index page.
func (server *Server) indexHandler(writer http.ResponseWriter, request *http.Request) {
	body := []string{
		"<center>",
		"<h2>PAC Manager</h2>",
		"</center>",
	}

	renderTemplate("", strings.Join(body, "\n"), writer)
}

// Run starts the server.
func (server *Server) Run() {
	handler := http.Server{
		Addr: fmt.Sprintf("%s:%d", server.GetAddress(), server.GetPort()),
	}

	http.HandleFunc("/", server.indexHandler)

	http.HandleFunc("/domains", server.domainListHandler)
	http.HandleFunc("/domains/edit/", server.domainEditHandler)
	http.HandleFunc("/domains/delete/", server.domainDeleteHandler)
	http.HandleFunc("/domains/add", server.domainAddHandler)
	http.HandleFunc("/domains/save", server.domainSaveHandler)

	http.HandleFunc("/proxies", server.proxyListHandler)
	http.HandleFunc("/proxies/delete/", server.proxyDeleteHandler)
	http.HandleFunc("/proxies/add", server.proxyAddHandler)
	http.HandleFunc("/proxies/save", server.proxySaveHandler)

	fmt.Println(fmt.Sprintf("🚀 Starting server on %s:%d", server.GetAddress(), server.GetPort()))

	if err := handler.ListenAndServe(); err != nil {
		logger.Fatalf(
			"server:run",
			"failed to start server: %s",
			err.Error(),
		)
	}
}
