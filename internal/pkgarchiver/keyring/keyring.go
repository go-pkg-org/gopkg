package keyring

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"os"
)

//go:generate mockgen -destination=../keyring_mock/keyring_mock.go -package=keyring_mock . Keyring

// Keyring represent the maintainer keyring
type Keyring interface {
	// CheckSignature check signature sig against file
	// and return corresponding Maintainer if exist, error otherwise
	CheckSignature(file, sig []byte) (Maintainer, error)
}

// Maintainer represent a maintainer
type Maintainer struct {
	Name string
}

type keyring struct {
	el openpgp.EntityList
}

func (k *keyring) CheckSignature(file, sig []byte) (Maintainer, error) {
	who, err := openpgp.CheckDetachedSignature(k.el, bytes.NewReader(file), bytes.NewReader(sig))
	if err != nil {
		return Maintainer{}, err
	}

	return Maintainer{Name: getMaintainerName(who)}, nil
}

// FromFile attempt to load keyring from given file
func FromFile(path string) (Keyring, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load keyring %s err: %s", path, err)
	}
	defer f.Close()

	k, err := openpgp.ReadKeyRing(f)
	if err != nil {
		return nil, fmt.Errorf("unable to load keyring %s err: %s", path, err)
	}

	return &keyring{el: k}, nil
}

func getMaintainerName(entity *openpgp.Entity) string {
	name := ""
	for id := range entity.Identities {
		name = id
	}

	return name
}
