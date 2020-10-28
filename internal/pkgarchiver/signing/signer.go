package signing

import (
	"bytes"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"os"
)

//go:generate mockgen -destination=../signing_mock/signing_mock.go -package=signing_mock . Signer

// Signer is something that can cryptographically sign data
type Signer interface {
	// Sign given message using configured key
	Sign(message []byte) ([]byte, error)
}

type signer struct {
	e *openpgp.Entity
}

func (s *signer) Sign(message []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := openpgp.DetachSign(&buf, s.e, bytes.NewReader(message), nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FromKeyFile instantiate signer from given key file
func FromKeyFile(path string) (Signer, error) {
	keyFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()

	signingKey, err := openpgp.ReadEntity(packet.NewReader(keyFile))
	if err != nil {
		return nil, err
	}

	return &signer{e: signingKey}, nil
}
