package remove

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/cache"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/rs/zerolog/log"
	"os"
)

// Remove given package
func Remove(pkgName string) error {
	cachePath, err := config.GetCachePath()
	if err != nil {
		return err
	}

	c, err := cache.Read(cachePath)
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
		log.Debug().Str("file", file).Msg("Removing file")
	}

	c.RemovePackage(pkgName)

	if err := cache.Write(cachePath, c); err != nil {
		return err
	}

	log.Info().Str("package", pkgName).Msg("Successfully removed package")
	return nil
}
