package pkg

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/control"
	"github.com/go-pkg-org/gopkg/internal/util"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const pkgFileExt = "pkg"

// Entry is a tiny struct to contain data for a specific
// entry that will be archived into a pkg file.
type Entry struct {
	FilePath    string
	ArchivePath string
}

// CreateEntries creates a slice with all files in a specific directory that should be added to the archive.
// The resulting value is a Entry, which maps a filepath to an archive path.
func CreateEntries(path string, pathPrefix string, fileTypes []string) ([]Entry, error) {
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Write file list.
	var fileList []Entry
	for _, file := range dirContent {
		if file.IsDir() {
			if file.Name() == ".git" || file.Name() == control.GoPkgDir {
				// The above directories should _not_ be included.
				// TODO: https://github.com/go-pkg-org/gopkg/issues/23
				continue
			}

			tmp, err := CreateEntries(filepath.Join(path, file.Name()), pathPrefix, fileTypes)
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

// Read reads a package and returns content.
func Read(path string) (map[string][]byte, error) {
	result := map[string][]byte{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(file)

	tr := tar.NewReader(buffer)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		out := bytes.NewBuffer(make([]byte, header.Size))
		if _, err := io.Copy(out, tr); err != nil {
			return nil, err
		}

		result[header.Name] = out.Bytes()
	}
	return result, nil
}

// Write creates a tar file from a set of ArchiveEntries.
func Write(path string, files []Entry, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("failed to create new tar source (file already exist)")
		}
	}

	var buffer bytes.Buffer
	tw := tar.NewWriter(&buffer)

	for _, file := range files {
		log.Trace().Str("file-path", file.FilePath).Str("archive-path", file.ArchivePath).Msg("Writing file")
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

// GetName return package name corresponding to given information
func GetName(pkgName, version, os, arch string, isSrc bool) (string, error) {
	// validate common information
	if pkgName == "" || version == "" {
		return "", fmt.Errorf("missing information to build package name")
	}

	if isSrc {
		return fmt.Sprintf("%s_%s-dev.%s", pkgName, version, pkgFileExt), nil
	}

	// validate binary specific information
	if os == "" || arch == "" {
		return "", fmt.Errorf("missing information to build package name")
	}

	return fmt.Sprintf("%s_%s_%s_%s.%s", pkgName, version, os, arch, pkgFileExt), nil
}

// ParseName parse existing package name and return found information
func ParseName(pkg string) (string, string, string, string, bool, error) {
	cleanPkg := strings.TrimSuffix(pkg, ".pkg")
	if strings.HasSuffix(cleanPkg, "-dev") {
		cleanPkg = strings.TrimSuffix(cleanPkg, "-dev")
		parts := strings.Split(cleanPkg, "_")
		if len(parts) != 2 {
			return "", "", "", "", false, fmt.Errorf("wrong source package name")
		}

		return parts[0], parts[1], "", "", true, nil
	}

	parts := strings.Split(cleanPkg, "_")
	if len(parts) != 4 {
		return "", "", "", "", false, fmt.Errorf("wrong binary package name")
	}

	return parts[0], parts[1], parts[2], parts[3], false, nil
}
