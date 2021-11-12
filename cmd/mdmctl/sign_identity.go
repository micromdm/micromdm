package main

import (
	"crypto/rsa"
	"os"

	macospkg "github.com/korylprince/go-macos-pkg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"
)

func signPackageWithIdentity(pkgpath, outpath, identitypath, identitypass string) error {
	identity, err := os.ReadFile(identitypath)
	if err != nil {
		return errors.Wrap(err, "reading identity")
	}
	key, cert, err := pkcs12.Decode(identity, identitypass)
	if err != nil {
		return errors.Wrap(err, "decoding identity")
	}

	pkg, err := os.ReadFile(pkgpath)
	if err != nil {
		return errors.Wrap(err, "reading package")
	}

	signed, err := macospkg.SignPkg(pkg, cert, key.(*rsa.PrivateKey))
	if err != nil {
		return errors.Wrap(err, "signing package")
	}

	if err = os.WriteFile(outpath, signed, 0644); err != nil {
		return errors.Wrap(err, "writing signed package")
	}
	return nil
}
