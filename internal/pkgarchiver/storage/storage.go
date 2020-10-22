package storage

import (
	"github.com/go-pkg-org/gopkg/internal/archive"
)

//go:generate mockgen -destination=../storage_mock/storage_mock.go -package=storage_mock . Storage

// Storage represent a storage support for the archive
type Storage interface {
	// GetIndex retrieve the index from the storage
	GetIndex() (archive.Index, error)
	// UpdateIndex update the index with given one
	UpdateIndex(index archive.Index) error
	// Upload upload given file to the storage
	Upload(file []byte, path string) error
}
