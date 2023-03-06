package storage

import (
	"os"
	"path"
)

// Storage represents a storage.
type Storage struct {
	// base represents the base path of the storage.
	base string
	// defaultPermissions represents the default permissions for new files.
	defaultPermissions os.FileMode
}

// NewStorage creates a new storage.
func NewStorage(base string, defaultPermissions os.FileMode) *Storage {
	return &Storage{
		base:               base,
		defaultPermissions: defaultPermissions,
	}
}

// GetBase returns the base path of the storage.
func (storage *Storage) GetBase() string {
	return storage.base
}

// GetDefaultPermissions returns the default permissions for new files.
func (storage *Storage) GetDefaultPermissions() os.FileMode {
	return storage.defaultPermissions
}

// GetPath returns the path of the target path.
func (storage *Storage) GetPath(targetPath string) string {
	return path.Join(storage.base, targetPath)
}

// FileExists returns true if the file exists.
func (storage *Storage) FileExists(targetPath string) bool {
	_, err := os.Stat(storage.GetPath(targetPath))
	return err == nil
}

// DirectoryExists returns true if the directory exists.
func (storage *Storage) DirectoryExists(targetPath string) bool {
	fileInfo, err := os.Stat(storage.GetPath(targetPath))
	return err == nil && fileInfo.IsDir()
}

// CreateFile creates the file.
func (storage *Storage) CreateFile(targetPath string) error {
	handle, err := os.OpenFile(storage.GetPath(targetPath), os.O_CREATE, storage.defaultPermissions)
	defer handle.Close()
	return err
}

// CreateDirectory creates the directory.
func (storage *Storage) CreateDirectory(targetPath string) error {
	return os.MkdirAll(storage.GetPath(targetPath), storage.defaultPermissions)
}

// OpenFile opens the file.
func (storage *Storage) OpenFile(targetPath string) (*os.File, error) {
	return os.OpenFile(storage.GetPath(targetPath), os.O_RDWR, storage.defaultPermissions)
}
