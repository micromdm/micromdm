// Package app sets up and manages the MicroMDM server application which
// is composed of multiple independent components such as datastores and
// services.
package app

import "github.com/go-kit/kit/log"

// exit status
const (
	success   int = iota
	badReturn     // general errors
	badInput      // issue with CLI input
)

// Main is MicroMDM's main() function. Main parses the CLI options given
// by the user and sets up an HTTP service for all the services built into
// MicroMDM.
func Main(logger log.Logger) (status int, err error) {
	config, err := loadConfig()
	if err != nil {
		return badInput, err
	}
	sm, err := setupServices(config, logger)
	if err != nil {
		return badReturn, err
	}
	err = serveHTTP(
		logger,
		makeHTTPHandler(logger, sm),
		config.TLS.Enabled,
		config.Server.ListenURL,
		config.TLS.PrivateKeyPath,
		config.TLS.CertificatePath,
	)
	if err != nil {
		return badReturn, err
	}
	return success, nil
}
