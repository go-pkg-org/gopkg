package pkgarchiver

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestHandleAcceptedPackage_NoMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//signer := signing_mock.NewMockSigner(ctrl)
	//storer := storage_mock.NewMockStorage(ctrl)

	/*if err := handleAcceptedPackage(signer, storer, archive.Index{}, f); err != ErrMissingPkgDefinition {
		t.Error("handleAcceptedPackage() should have returned ErrMissingPkgDefinition")
	}*/
}

func TestHandleAcceptedPackage(t *testing.T) {

}
