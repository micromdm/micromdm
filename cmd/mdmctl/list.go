package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/go-kit/kit/log"
	"github.com/micromdm/micromdm/core/list"
)

type tableOutput struct{ w *tabwriter.Writer }

func (out *tableOutput) BasicHeader() {
	fmt.Fprintf(out.w, "UDID\tSerialNumber\tEnrollmentStatus\tLastSeen\n")
}

func (out *tableOutput) BasicFooter() {
	out.w.Flush()
}

func listDevices(args []string) error {
	if len(args) == 0 || args[0] != "devices" {
		return errors.New("must specify resource name as an argument (ex: devices)")
	}
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &tableOutput{w}
	out.BasicHeader()
	defer out.BasicFooter()
	logger := log.NewLogfmtLogger(os.Stderr)
	// TODO needs some config for authentication and server url client side.
	instance := os.Getenv("MICROMDM_SERVER_URL")
	if instance == "" {
		return errors.New("MICROMDM_SERVER_URL not set")
	}

	token := os.Getenv("MICROMDM_API_TOKEN")
	if token == "" {
		return errors.New("MICROMDM_API_TOKEN not set")
	}
	svc, err := list.NewClient(instance, logger, token)
	if err != nil {
		return err
	}

	ctx := context.Background()
	devices, err := svc.ListDevices(ctx, list.ListDevicesOption{})
	if err != nil {
		return err
	}
	for _, d := range devices {
		fmt.Fprintf(out.w, "%s\t%s\t%v\t%s\n", d.UDID, d.SerialNumber, d.EnrollmentStatus, d.LastSeen)
	}
	return nil
}
