package server

import (
	"strconv"
)

// Server contains the data for a server.
type Server struct {
	// Host is the server's IP address.
	Host string `mapstructure:"host" json:"host" yaml:"host""`
	// Port is the server's port.
	Port uint `mapstructure:"port" json:"port" yaml:"port"`
}

// InitializeDefault creates a new server with default values.
func InitializeDefault() *Server {
	return &Server{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

// GetHost returns the server's IP address.
func (server *Server) GetHost() string {
	return server.Host
}

// GetPort returns the server's port.
func (server *Server) GetPort() uint {
	return server.Port
}

// GetAddress returns the server's address.
func (server *Server) GetAddress() string {
	return server.Host + ":" + strconv.Itoa(int(server.Port))
}
