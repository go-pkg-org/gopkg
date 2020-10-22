package remove

import (
	"fmt"
	"os"

	"github.com/go-pkg-org/gopkg/internal/cache"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/rs/zerolog/log"
)

// Remove given package
func Remove(pkgName string) error {
	config, err := config.Default()
	if err != nil {
		return err
	}

	c, err := cache.Read(config.CachePath)
	if err != nil {
		return err
	}

	files := c.GetFiles(pkgName)
	if files == nil {
		return fmt.Errorf("package %s is not installed", pkgName)
	}

	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Warn().Str("err", err.Error()).Str("file", file).Msg("Error while removing file")
		}
		log.Trace().Str("file", file).Msg("Removing file")
	}

	c.RemovePackage(pkgName)

	if err := cache.Write(config.CachePath, c); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully removed package")
	return nil
}
