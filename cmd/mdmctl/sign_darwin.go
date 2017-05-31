package main

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func signPackage(pkgpath, outpath string, developerID string) error {
	cmd := exec.Command("/usr/bin/productsign", "--sign", developerID, pkgpath, outpath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "signing package")
	}
	return nil
}

func checkSignature(pkgpath string) (bool, error) {
	cmd := exec.Command("pkgutil", "--check-signature", pkgpath)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil && !isNoSignature(out) {
		return false, errors.Wrap(err, "checking signature")
	}

	if bytes.Contains(out, []byte(`Status: signed`)) {
		return true, nil
	}

	return false, nil
}

func isNoSignature(out []byte) bool {
	return bytes.Contains(out, []byte(`Status: no signature`))
}
