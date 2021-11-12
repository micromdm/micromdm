//go:build !darwin
// +build !darwin

package main

import (
	"fmt"
	"os"

	goerrors "errors"

	macospkg "github.com/korylprince/go-macos-pkg"
	"github.com/pkg/errors"
)

func signPackage(path, outpath, developerID string) error {
	fmt.Println("[WARNING] package signing with -sign only implemented on macOS. Use -sign-identity instead")
	return nil
}

func checkSignature(pkgpath string) (bool, error) {
	buf, err := os.ReadFile(pkgpath)
	if err != nil {
		return false, errors.Wrap(err, "reading package")
	}
	if err = macospkg.VerifyPkg(buf); err != nil {
		if goerrors.Is(err, macospkg.ErrNotSigned) {
			return false, nil
		}
		return false, errors.Wrap(err, "verifying package")
	}
	return true, nil
}
