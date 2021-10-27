//go:build !darwin
// +build !darwin

package main

import "fmt"

func signPackage(path, outpath, developerID string) error {
	fmt.Println("[WARNING] package signing only implemented on macOS")
	return nil
}

func checkSignature(pkgpath string) (bool, error) {
	fmt.Println("[WARNING] package signing only implemented on macOS. An unsigned macOS package will not install with MDM.")
	return true, nil
}
