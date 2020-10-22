package sign

import (
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/rs/zerolog/log"
	"os/exec"
)

// Sign sign given package
func Sign(pkgPath string) error {
	key := config.GetSigningKey()
	if key == "TODO" {
		return fmt.Errorf("please configure a signing key trough GOPKG_SIGNING_KEY env variable")
	}

	log.Debug().Str("signingKey", key).Str("package", pkgPath).Msg("Signing package")

	if err := signPackage(pkgPath, pkgPath+".asc", key, ""); err != nil {
		return err
	}

	log.Info().Str("signingKey", key).Str("package", pkgPath).Msg("Successfully signed package")

	return nil
}

// signPackage sign given package using given parameters
func signPackage(pkgPath, ascPath, key, keyring string) error {
	var args []string
	args = append(args, "--pinentry-mode", "loopback", "--default-key", key)

	if keyring != "" {
		args = append(args, "--no-default-keyring", "--keyring", keyring)
	}

	args = append(args, "--out", ascPath, "--detach-sign", pkgPath)

	cmd := exec.Command("gpg", args...)

	return cmd.Run()
}
