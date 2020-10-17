package util

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func CreateFileMap (path string) (map[string]string, error) {
	dirContent, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Create file list.
	var fileList = map[string]string {}
	for _, file := range dirContent {
		if strings.Index(file.Name(), ".") == 0 {
			continue
		}

		if file.IsDir()  {
			tmp, err := CreateFileMap(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
			for t, p := range tmp {
				fileList[filepath.Join(file.Name(), t)] = p
			}
		} else {
			fileList[file.Name()] = filepath.Join(path, file.Name())
		}
	}

	return fileList, nil
}


// Create a tar file from a set of files.
// The expected input `files` should be 'filename  => filePath' where the filename
// is the path that will be created inside the tar file, and the filePath is the path
// to the file to add.
// The path is the path to the where the tar file will be created.
func CreateTar(path string, files map[string]string, overwrite bool) error {
	if !overwrite {
		if _, err := ioutil.ReadFile(path); err != nil {
			return fmt.Errorf("failed to create new tar source (file already exist)")
		}
	}

	var buffer bytes.Buffer
	tw := tar.NewWriter(&buffer)

	for fileName, file := range files {
		fileBody, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Name: fileName,
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

	return ioutil.WriteFile(path, buffer.Bytes(), 0755)
}