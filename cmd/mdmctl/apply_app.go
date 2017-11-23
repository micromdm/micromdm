package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/groob/plist"
	"github.com/pkg/errors"

	"github.com/micromdm/micromdm/mdm/appmanifest"
)

func (cmd *applyCommand) applyApp(args []string) error {
	flagset := flag.NewFlagSet("app", flag.ExitOnError)
	var (
		flPkgPath     = flagset.String("pkg", "", "path to a distribution pkg.")
		flPkgURL      = flagset.String("pkg-url", "", "use custom pkg url")
		flAppManifest = flagset.String("manifest", "-", `path to an app manifest. optional,
		will be created if file does not exist.`)

		flHashSize = flagset.Int64("md5size", appmanifest.DefaultMD5Size, "md5 hash size in bytes (optional)")
		flSign     = flagset.String("sign", "", "sign package before importing, requires specifying a product ID (optional)")
		flUpload   = flagset.Bool("upload", false, "upload package and/or manifest to micromdm repository.")
	)
	flagset.Usage = usageFor(flagset, "mdmctl apply app [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	pkgurl := *flPkgURL
	if pkgurl == "" {
		su, err := cmd.serverRepoURL()
		if err != nil {
			return err
		}
		pkgurl = pkgURL(su, *flPkgPath)
	}

	pkg := *flPkgPath
	signed, err := checkSignature(pkg)
	if err != nil {
		return err
	}

	distribution, err := checkDistribution(*flPkgPath)
	if err != nil {
		return err
	}
	if !distribution {
		fmt.Println(`
[WARNING] The package you're importing is not a macOS distribution package. MDM requires distribution packages.
You can turn a flat package to a distribution one with the productbuild command:

productbuild --package someFlatPkg.pkg myNewPkg.pkg

Please rebuild the package and re-run the command.

		`)
	}

	if !signed {
		if *flSign == "" {
			flagset.Usage()
			return errors.New(`MDM packages must be signed. Provide signed package or Developer ID with -sign flag`)
		}
		outpath := filepath.Join(os.TempDir(), filepath.Base(*flPkgPath))
		if err := signPackage(*flPkgPath, outpath, *flSign); err != nil {
			return err
		}
		pkg = outpath // use signed package to create the manifest
	}

	// open pkg file
	f, err := os.Open(pkg)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := []appmanifest.Option{appmanifest.WithMD5Size(*flHashSize)}
	manifest, err := appmanifest.Create(&appFile{f}, pkgurl, opts...)
	if err != nil {
		return errors.Wrap(err, "creating manifest")
	}

	var buf bytes.Buffer
	enc := plist.NewEncoder(&buf)
	enc.Indent("  ")
	if err := enc.Encode(manifest); err != nil {
		return err
	}

	if *flUpload {
		// we read the file to generate the appmanifest, so we need to seek to the beginning of the file again.
		if _, err := f.Seek(0, 0); err != nil {
			return errors.Wrap(err, "reset pkg file reader")
		}
		err := cmd.applysvc.UploadApp(context.TODO(), nameMannifest(f.Name()), &buf, filepath.Base(f.Name()), f)
		if err != nil {
			return err
		}
	}

	switch *flAppManifest {
	case "":
	case "-":
		_, err := os.Stdout.Write(buf.Bytes())
		return err
	default:
		return ioutil.WriteFile(*flAppManifest, buf.Bytes(), 0644)
	}

	return nil
}

func nameMannifest(pkgName string) string {
	trimmed := strings.TrimSuffix(filepath.Base(pkgName), filepath.Ext(pkgName))
	return trimmed + ".plist"
}

func checkDistribution(pkgPath string) (bool, error) {
	const (
		xarHeaderMagic = 0x78617221
		xarHeaderSize  = 28
	)

	f, err := os.Open(pkgPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	hdr := make([]byte, xarHeaderSize)
	_, err = f.ReadAt(hdr, 0)
	if err != nil {
		return false, err
	}
	tocLenZlib := binary.BigEndian.Uint64(hdr[8:16])
	ztoc := make([]byte, tocLenZlib)
	_, err = f.ReadAt(ztoc, xarHeaderSize)
	if err != nil {
		return false, err
	}

	br := bytes.NewBuffer(ztoc)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return false, err
	}
	toc, err := ioutil.ReadAll(zr)
	if err != nil {
		return false, err
	}
	return bytes.Contains(toc, []byte(`<name>Distribution</name>`)), nil
}

func (cmd *applyCommand) serverRepoURL() (string, error) {
	return repoURL(cmd.config.ServerURL)
}

func pkgURL(repoURL, pkgPath string) string {
	u, _ := url.Parse(repoURL)
	newPath := path.Join(u.Path, filepath.Base(pkgPath))
	u.Path = newPath
	return u.String()
}

func repoURL(server string) (string, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return "", err
	}
	serverURL.Path = "/repo"
	return serverURL.String(), nil
}

// replaces .pkg with .plist
func manifestURL(repoURL, pkgPath string) string {
	pu := pkgURL(repoURL, pkgPath)
	trimmed := strings.TrimSuffix(pu, path.Ext(pu))
	return trimmed + ".plist"
}

type appFile struct {
	*os.File
}

func (af *appFile) Size() int64 {
	info, err := af.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return info.Size()
}
