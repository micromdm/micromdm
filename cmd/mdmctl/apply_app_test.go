package main

import (
	"testing"
)

func TestPkgURL(t *testing.T) {
	baseURL := "https://mdm.acme.co/repo"
	pkgPath := "/Users/user/Desktop/dep_pkg-foo_1.1.pkg"
	u := pkgURL(baseURL, pkgPath)
	if have, want := u, "https://mdm.acme.co/repo/dep_pkg-foo_1.1.pkg"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

}
