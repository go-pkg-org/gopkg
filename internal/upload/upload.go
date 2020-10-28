package upload

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// Upload upload given package to given archive
func Upload(pkgPath, archive string) error {
	log.Info().Str("package", pkgPath).Str("archive", archive).Msg("Uploading package")

	pkgAscPath := fmt.Sprintf("%s.asc", pkgPath)

	// Validate files exist
	if _, err := os.Stat(pkgPath); err != nil {
		return err
	}
	if _, err := os.Stat(pkgAscPath); err != nil {
		return err
	}

	// Upload the package
	if err := uploadPackage(pkgPath, pkgAscPath, fmt.Sprintf("%s/packages", archive)); err != nil {
		log.Err(err).Msg("error while uploading package")
		return err
	}

	log.Info().Str("package", pkgPath).Str("archive", archive).Msg("Package successfully uploaded")

	return nil
}

func uploadPackage(pkgPath, pkgAscPath, where string) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if err := addFileFormParam(writer, "package", pkgPath); err != nil {
		return err
	}
	if err := addFileFormParam(writer, "packageAsc", pkgAscPath); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", where, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error while uploading file: %s", resp.Status)
	}

	return nil
}

func addFileFormParam(w *multipart.Writer, param, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()

	part, err := w.CreateFormFile(param, fi.Name())
	if err != nil {
		return err
	}
	if _, err := part.Write(fileContents); err != nil {
		return err
	}

	return nil
}
