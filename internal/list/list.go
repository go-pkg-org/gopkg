package list

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/cache"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/rs/zerolog/log"
)

// List packages from cache
func List(onlyInstalled bool) error {
	if !onlyInstalled {
		return fmt.Errorf("not implemented at the moment")
	}

	cachePath, err := config.GetCachePath()
	if err != nil {
		return err
	}

	c, err := cache.Read(cachePath)
	if err != nil {
		return err
	}

	if len(c.Packages) == 0 {
		log.Info().Msg("No packages installed")
		return nil
	}

	for pkg := range c.Packages {
		log.Info().Str("package", pkg).Msg("")
	}

	return nil
}
