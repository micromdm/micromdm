// package appmanifest provides utilities for managing app manifest files
// used by MDM InstallApplication commands.
package appmanifest

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// DefaultMD5Size is the default size of each file chunk that needs to be hashed
const DefaultMD5Size = 10 << 20 // 10MB

// http://help.apple.com/deployment/osx/#/ior5df10f73a
type Manifest struct {
	ManifestItems []Item `plist:"items"`
}

type Item struct {
	Assets []Asset `plist:"assets"`
	// Apple claims the metadata struct is required,
	// but testing shows otherwise.
	Metadata *Metadata `plist:"metadata,omitempty"`
}

type Asset struct {
	Kind    string   `plist:"kind"`
	MD5Size int64    `plist:"md5-size"`
	MD5s    []string `plist:"md5s"`
	URL     string   `plist:"url"`
}

type Metadata struct {
	BundleInfo
	Items    []BundleInfo `plist:"items,omitempty"`
	Kind     string       `plist:"kind"`
	Subtitle string       `plist:"subtitle"`
	Title    string       `plist:"title"`
}

type BundleInfo struct {
	BundleIdentifier string `plist:"bundle-identifier"`
	BundleVersion    string `plist:"bundle-version"`
}

// File is an io.Reader which knows its size.
type File interface {
	io.Reader
	Size() int64
}

// Create an AppManifest and write it to an io.Writer.
func Create(file File, url string, md5Size int64) (*Manifest, error) {
	fSize := file.Size()
	if md5Size > fSize {
		md5Size = fSize
	}

	// create a list of md5s
	md5s, err := calculateMD5s(file, md5Size)
	if err != nil {
		return nil, errors.Wrap(err, "calculate appmanifest md5s")
	}

	// create an asset
	ast := Asset{
		Kind:    "software-package",
		MD5Size: md5Size,
		MD5s:    md5s,
		URL:     url,
	}

	// make a manifest
	m := Manifest{
		ManifestItems: []Item{
			Item{
				Assets: []Asset{ast},
			},
		},
	}

	return &m, nil
}

// reads a file and returns a slice of hashes, one for each
// 10mb chunk
func calculateMD5s(f io.Reader, s int64) ([]string, error) {
	h := md5.New()
	var md5s []string
	for {
		n, err := io.CopyN(h, f, s)
		if n > 0 {
			md5s = append(md5s, fmt.Sprintf("%x", h.Sum(nil)))
			h.Reset()
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return md5s, err
		}
	}
}
