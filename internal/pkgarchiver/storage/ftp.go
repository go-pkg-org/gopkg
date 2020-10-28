package storage

import (
	"bytes"
	"encoding/json"
	"github.com/go-pkg-org/gopkg/internal/archive"
	"github.com/jlaffaye/ftp"
	"path/filepath"
	"strings"
	"time"
)

type ftpStorage struct {
	conn *ftp.ServerConn
}

func (f *ftpStorage) GetIndex() (archive.Index, error) {
	resp, err := f.conn.Retr("index.json")
	if err != nil {
		// No index exist at the time, create new one
		if err.Error() == "550 Can't open index.json: No such file or directory" {
			return archive.Index{Packages: map[string]archive.Package{}}, nil
		}

		return archive.Index{}, err
	}

	var index archive.Index
	if err := json.NewDecoder(resp).Decode(&index); err != nil {
		return archive.Index{}, err
	}
	resp.Close() // Need to be close or we are failing the FTP client

	return index, nil
}

func (f *ftpStorage) UpdateIndex(index archive.Index) error {
	b, err := json.Marshal(index)
	if err != nil {
		return err
	}

	return f.Upload(b, "index.json")
}

func (f *ftpStorage) Upload(file []byte, path string) error {
	// first of all create any missing directories
	if err := f.makeMissingDirectories(filepath.Dir(path)); err != nil {
		return err
	}

	if err := f.conn.Stor(path, bytes.NewReader(file)); err != nil {
		return err
	}

	return nil
}

func (f *ftpStorage) makeMissingDirectories(target string) error {
	parts := strings.Split(target, "/")
	path := ""
	for _, part := range parts {
		path = filepath.Join(path, part)
		if err := f.conn.MakeDir(path); err != nil && err.Error() != "550 Can't create directory: File exists" {
			return err
		}
	}

	return nil
}

// NewFTPStorage create a brand new storage using FTP has backend
func NewFTPStorage(host, user, pass, baseDir string) (Storage, error) {
	ftpClient, err := ftp.Dial(host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	// Authenticate
	if err := ftpClient.Login(user, pass); err != nil {
		return nil, err
	}
	if baseDir != "" {
		if err := ftpClient.ChangeDir(baseDir); err != nil {
			return nil, err
		}
	}

	return &ftpStorage{conn: ftpClient}, nil
}
