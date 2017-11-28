package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/device"
)

//TODO: select command

func main() {
	var (
		flPath = flag.String("path", "", "path to boltdb")
	)
	flag.Parse()

	db, err := bolt.Open(*flPath, 0644, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Fatal(err)
	}

	repl := &REPL{
		db:      db,
		devices: &device.DB{db},
	}

	if err := repl.Run(); err != nil {
		log.Fatal(err)
	}
}

type REPL struct {
	db      *bolt.DB
	devices *device.DB
}

func (r *REPL) DeviceByUDID(udid string) (*device.Device, error) {
	return r.devices.DeviceByUDID(udid)
}
func (r *REPL) ListDevices() ([]device.Device, error) {
	return r.devices.List()
}

func (r *REPL) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(r.Prompt())
		rawLine, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if rawLine == "" {
			continue
		}

		line := strings.TrimSpace((rawLine))
		split := strings.Split(line, " ")

		switch cmd := split[0]; cmd {
		case "devices":
			devices, err := r.ListDevices()
			if err != nil {
				log.Println(err)
				continue
			}
			if len(split) == 1 {
				printDevices(devices)
				continue
			}
			switch split[1] {
			case "-udid":
				if len(split) < 3 {
					fmt.Println("udid not entered")
					continue
				}
				udid := strings.TrimSpace((split[2]))
				if err != nil {
					log.Println(err)
					continue
				}
				dev, err := r.DeviceByUDID(udid)
				if err != nil {
					log.Println(err)
					continue
				}
				printDevice(*dev)

			//case "-serial":

			default:
				fmt.Println("wrong parameter: ", split[1])
			}
		default:
			fmt.Println("command not exists: ", split[0])

		}
	}
}

func (r *REPL) Prompt() string {
	out := bytes.NewBufferString("$ ")
	return out.String()
}

func printDevices(devices []device.Device) {
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &struct{ w *tabwriter.Writer }{w}
	fmt.Fprintf(out.w, "%40s\t%s\t%v\t%s\n", "UDID", "SerialNumber", "EnrollmentStatus", "LastSeen")
	for _, d := range devices {
		fmt.Fprintf(out.w, "%40s\t%s\t%v\t%s\n", d.UDID, d.SerialNumber,
			d.Enrolled, d.LastCheckin)
	}
	out.w.Flush()
}
func printDevice(dev device.Device) {
	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	out := &struct{ w *tabwriter.Writer }{w}
	fmt.Fprintf(out.w, "%40s\t%s\t%v\t%s\n", "UDID", "SerialNumber", "EnrollmentStatus", "LastSeen")
	fmt.Fprintf(out.w, "%40s\t%s\t%v\t%s\n", dev.UDID, dev.SerialNumber,
		dev.Enrolled, dev.LastCheckin)

	out.w.Flush()
}
