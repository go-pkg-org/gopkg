package pkgarchiver

import (
	"bytes"
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/archive"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/go-pkg-org/gopkg/internal/pkgarchiver/keyring"
	"github.com/go-pkg-org/gopkg/internal/pkgarchiver/signing"
	"github.com/go-pkg-org/gopkg/internal/pkgarchiver/storage"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// ErrMissingPkgDefinition is returned when missing package definition from archive
var ErrMissingPkgDefinition = fmt.Errorf("missing package definition (package.yaml)")

// Execute is the main entrypoint of pkgarchiver
func Execute(c *cli.Context) error {
	// Load archive signing key
	signer, err := signing.FromKeyFile(c.String("signing-key"))
	if err != nil {
		return err
	}

	// Load maintainers keyring
	maintainerKeyring, err := keyring.FromFile(c.String("maintainer-keyring"))
	if err != nil {
		return fmt.Errorf("error while loading keyring: %s", err)
	}

	// Open the storage session
	storer, err := storage.NewFTPStorage(c.String("ftp-host"), c.String("ftp-user"),
		c.String("ftp-pass"), c.String("ftp-dir"))
	if err != nil {
		return err
	}

	// Load existing index
	index, err := storer.GetIndex()
	if err != nil {
		return err
	}
	log.Debug().Int("count", len(index.Packages)).Msg("Loaded packages index")

	// Create HTTP server
	http.HandleFunc("/packages", handleUpload(maintainerKeyring, signer, index, storer))
	log.Info().Str("address", ":8888").Msg("Listening for packages")

	// Listen for packages
	return http.ListenAndServe(":8888", nil)
}

func handleUpload(maintainerKeyring keyring.Keyring, signer signing.Signer,
	index archive.Index, storer storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pkgFile, header, err := readFormFile(r, "package")
		if err != nil {
			log.Err(err).Msg("error while reading package")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pkgFileAsc, _, err := readFormFile(r, "packageAsc")
		if err != nil {
			log.Err(err).Msg("error while reading package")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info().Str("package", header.Filename).Msg("Handling package")

		// Validate signature
		maintainer, err := maintainerKeyring.CheckSignature(pkgFile, pkgFileAsc)
		if err != nil {
			log.Err(err).Msg("error while checking package signature")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info().
			Str("package", header.Filename).
			Str("maintainer", maintainer.Name).
			Msg("Accepted package")

		if err := handleAcceptedPackage(signer, storer, index, pkgFile); err != nil {
			log.Err(err).Msg("error while uploading package")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func handleAcceptedPackage(
	signer signing.Signer,
	storer storage.Storage,
	index archive.Index,
	pkgBytes []byte) error {
	// Read package
	pkgContent, err := pkg.Read(bytes.NewReader(pkgBytes))
	if err != nil {
		return err
	}

	meta, err := pkgContent.Metadata()
	if err != nil {
		return ErrMissingPkgDefinition
	}

	log.Debug().
		Str("alias", meta.Alias).
		Str("version", meta.ReleaseVersion).
		Str("os", meta.TargetOS).
		Str("arch", meta.TargetArch).
		Msg("Uploading package")

	// Control package cannot be allowed at the moment since doesn't contains package.yaml file
	// and thus GetMetadata() will fail prematurely
	var pkgType pkg.Type
	if meta.Source() {
		pkgType = pkg.Source
	} else {
		pkgType = pkg.Binary
	}

	// Compute file name
	fileName, err := pkg.GetFileName(meta.Alias, meta.ReleaseVersion, meta.TargetOS, meta.TargetArch, pkgType)
	if err != nil {
		return err
	}

	// Create the package signature
	sig, err := signer.Sign(pkgBytes)
	if err != nil {
		return fmt.Errorf("error while signing package: %s", err)
	}

	// Upload the package
	if err := storer.Upload(pkgBytes, fmt.Sprintf("%s/%s", meta.Alias, fileName)); err != nil {
		return fmt.Errorf("error while uploading package: %s", err)
	}

	// Upload the signature
	if err := storer.Upload(sig, fmt.Sprintf("%s/%s.asc", meta.Alias, fileName)); err != nil {
		return fmt.Errorf("error while uploading package signature: %s", err)
	}

	// Reflect changes to index
	var p archive.Package
	if _, ok := index.Packages[meta.Alias]; ok {
		p = index.Packages[meta.Alias]
	} else {
		p = archive.Package{
			Releases:      map[string][]archive.Release{},
			LatestRelease: "",
		}
	}

	// Update the package status
	p.Releases[meta.ReleaseVersion] = append(index.Packages[meta.Alias].Releases[meta.ReleaseVersion], archive.Release{
		OS:   meta.TargetOS,
		Arch: meta.TargetArch,
		Path: fmt.Sprintf("%s/%s", meta.Alias, fileName),
	})
	p.LatestRelease = meta.ReleaseVersion

	// Update index
	index.Packages[meta.Alias] = p
	if err := storer.UpdateIndex(index); err != nil {
		return err
	}

	log.Info().Str("alias", meta.Alias).
		Str("version", meta.ReleaseVersion).
		Msg("Successfully uploaded package & signature")

	return nil
}

func readFormFile(r *http.Request, paramName string) ([]byte, *multipart.FileHeader, error) {
	f, header, err := r.FormFile(paramName)
	if err != nil {
		return nil, nil, err
	}

	f.Seek(0, io.SeekStart)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	f.Close()

	return b, header, nil
}
