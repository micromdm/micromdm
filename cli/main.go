package cli

import (
	"fmt"
	"os"

	"github.com/micromdm/micromdm/version"
)

func Main() {
	switch os.Args[1] {
	case "version":
		version.Print()
		return
	case "gencert":
		fmt.Printf("coming soon\n")
		return
	default:
		fmt.Printf("no such command")
		return
	}

}
