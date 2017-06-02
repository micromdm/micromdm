package main

import (
	"bytes"
	"compress/zlib"
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
	"github.com/micromdm/micromdm/appmanifest"
	"github.com/pkg/errors"
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
	serverURL, err := url.Parse(cmd.config.ServerURL)
	if err != nil {
		return "", err
	}
	serverURL.Path = "/repo"
	return serverURL.String(), nil
}

func pkgURL(repoURL, pkgPath string) string {
	return path.Join(repoURL, filepath.Base(pkgPath))
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
