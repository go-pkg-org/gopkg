package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Entry is a tiny struct to contain data for a specific
// entry that will be archived into a pkg file.
type Entry struct {
	FilePath    string
	ArchivePath string
}

// CreateFileMap creates a slice with all files in a specific directory that should be added to the archive.
// The resulting value is a Entry, which maps a filepath to an archive path.
func CreateFileMap(path string, pathPrefix string, fileTypes []string) ([]Entry, error) {
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Create file list.
	var fileList []Entry
	for _, file := range dirContent {
		if strings.Index(file.Name(), ".") == 0 {
			// No dot-files or directories will be added.
			continue
		}

		if file.IsDir() {
			tmp, err := CreateFileMap(filepath.Join(path, file.Name()), pathPrefix, fileTypes)
			if err != nil {
				return nil, err
			}
			for _, p := range tmp {
				fileList = append(fileList, Entry{
					FilePath:    p.FilePath,
					ArchivePath: filepath.Join(pathPrefix, file.Name(), p.ArchivePath),
				})
			}
		} else {
			if len(fileTypes) != 0 && !util.Contains(fileTypes, filepath.Ext(file.Name())) {
				continue
			}
			fileList = append(fileList, Entry{
				FilePath:    filepath.Join(path, file.Name()),
				ArchivePath: filepath.Join(pathPrefix, file.Name()),
			})
		}
	}

	return fileList, nil
}

// Create creates a tar file from a set of ArchiveEntries.
func Create(path string, files []Entry, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("failed to create new tar source (file already exist)")
		}
	}

	var buffer bytes.Buffer
	tw := tar.NewWriter(&buffer)

	for _, file := range files {
		fileBody, err := ioutil.ReadFile(file.FilePath)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: file.ArchivePath,
			Mode: 0644,
			Size: int64(len(fileBody)),
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tw.Write(fileBody); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0644)
}
