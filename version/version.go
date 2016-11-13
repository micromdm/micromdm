// Package version provides utilities for managing and printing the current
// version of the MicroMDM binary.
package version

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// version flags are set at build time with the go build tool.
// example:
// go build -ldflags "-X github.com/micromdm/micromdm/version.version=v1.0.0-dev"
var (
	version   = "unknown"
	gitBranch = "unknown"
	gitRev    = "unknown"
	goVersion = "unknown"
	buildTime = "unknown"
)

// Info holds version and build info for the current MicroMDM binary.
type Info struct {
	Version   string `json:"version"`
	Branch    string `json:"branch"`
	Revision  string `json:"revision"`
	GoVersion string `json:"go_version"`
	BuildDate string `json:"build_date"`
}

// Version returns a struct with the current version information.
func Version() Info {
	return Info{
		Version:   version,
		Branch:    gitBranch,
		Revision:  gitRev,
		GoVersion: goVersion,
		BuildDate: buildTime,
	}
}

// Print outputs the app name and version string.
func Print() {
	v := Version()
	fmt.Printf("MicroMDM - version %s\n", v.Version)
}

// PrintFull outputs the app name and detailed version information.
func PrintFull() {
	v := Version()
	fmt.Printf("MicroMDM - version %s\n", v.Version)
	fmt.Printf("  branch: \t%s\n", v.Branch)
	fmt.Printf("  revision: \t%s\n", v.Revision)
	fmt.Printf("  build date: \t%s\n", v.BuildDate)
	fmt.Printf("  go version: \t%s\n", v.GoVersion)
}

// Handler provides an HTTP Handler that returns JSON formatted version info.
func Handler() http.Handler {
	v := Version()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(v)

	})
}
