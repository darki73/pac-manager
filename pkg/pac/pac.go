package pac

import "github.com/darki73/pac-manager/pkg/storage"

// Pac represents a pac service configuration.
type Pac struct {
	// storage represents the storage.
	storage *storage.Storage
	// name represents the name of the pac file.
	name string
}

// NewPac creates a new pac service.
func NewPac(storage *storage.Storage, name string) *Pac {
	pac := &Pac{
		storage: storage,
		name:    name,
	}

	pac.bootstrap()

	return pac
}

// GetStorage returns the storage.
func (pac *Pac) GetStorage() *storage.Storage {
	return pac.storage
}

// GetName returns the name of the pac file.
func (pac *Pac) GetName() string {
	return pac.name
}

// bootstrap bootstraps the pac service.
func (pac *Pac) bootstrap() {
	if !pac.GetStorage().FileExists(pac.GetName()) {
		pac.GetStorage().CreateFile(pac.GetName())
	}
}
