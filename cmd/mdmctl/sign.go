// +build !darwin

package main

import "fmt"

func signPackage(path, outpath, developerID string) error {
	return fmt.Errorf("package signing only implemented on macOS")
}

func checkSignature(pkgpath string) (bool, error) {
	return false, fmt.Errorf("package signing only implemented on macOS")
}
