package database

// Database represents the database configuration.
type Database struct {
	// Name represents the name of the database.
	Name string `mapstructure:"name" yaml:"name" json:"name" toml:"name"`
	// Path represents the path of the database.
	Path string `mapstructure:"path" yaml:"path" json:"path" toml:"path"`
}

// InitializeDefaults initializes the default values for the database configuration.
func InitializeDefaults() *Database {
	return &Database{
		Name: "pacm",
		Path: "/tmp",
	}
}

// GetName returns the name of the database.
func (database *Database) GetName() string {
	return database.Name
}

// GetPath returns the path of the database.
func (database *Database) GetPath() string {
	return database.Path
}
