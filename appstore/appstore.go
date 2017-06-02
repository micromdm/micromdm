// package appstore provides an abstraction for uploading files and manifests
// to a repository.
package appstore

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/groob/plist"
	"github.com/micromdm/micromdm/appmanifest"
	"github.com/pkg/errors"
)

type AppStore interface {
	SaveFile(name string, f io.Reader) error
	Manifest(name string) (*appmanifest.Manifest, error)
	Apps() ([]string, error)
}

type Repo struct {
	Path string
}

func (r *Repo) SaveFile(name string, f io.Reader) error {
	fname := filepath.Join(r.Path, name)
	file, err := os.Create(fname)
	if err != nil {
		return errors.Wrapf(err, "saving file %s", name)
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	return err
}

func (r *Repo) Manifest(name string) (*appmanifest.Manifest, error) {
	manifestName := name
	if !strings.HasSuffix(name, ".plist") {
		manifestName = name + ".plist"

	}
	fname := filepath.Join(r.Path, manifestName)
	file, err := os.Open(fname)
	if err != nil {
		return nil, errors.Wrapf(err, "reading manifest %s", name)
	}
	defer file.Close()

	var m appmanifest.Manifest
	if err := plist.NewDecoder(file).Decode(&m); err != nil {
		return nil, errors.Wrap(err, "decoding manifest file")
	}

	return &m, nil
}

func (r *Repo) Apps() ([]string, error) {
	files, err := ioutil.ReadDir(r.Path)
	if err != nil {
		return nil, err
	}

	var manifests []string
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".plist" {
			continue
		}
		manifests = append(manifests, filepath.Base(file.Name()))
	}

	return manifests, nil
}
