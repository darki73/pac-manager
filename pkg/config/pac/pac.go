package pac

// Pac represents the pac configuration.
type Pac struct {
	// Name represents the name of the pac.
	Name string `mapstructure:"name" yaml:"name" json:"name" toml:"name"`
	// Path represents the path of the pac.
	Path string `mapstructure:"path" yaml:"path" json:"path" toml:"path"`
}

// InitializeDefaults initializes the default values for the pac configuration.
func InitializeDefaults() *Pac {
	return &Pac{
		Name: "proxy.pac",
		Path: "/tmp",
	}
}

// GetName returns the name of the pac.
func (pac *Pac) GetName() string {
	return pac.Name
}

// GetPath returns the path of the pac.
func (pac *Pac) GetPath() string {
	return pac.Path
}
