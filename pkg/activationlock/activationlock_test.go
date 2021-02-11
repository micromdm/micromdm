package activationlock

import (
	"encoding/hex"
	"testing"
)

func TestBypassCode(t *testing.T) {
	tests := []struct {
		key           string
		humanReadable string
		hash          string
	}{
		{
			key:           "1ea841db5edfafe6075b5ae0d845d254",
			humanReadable: "3UM43-PUYVY-QYD1-UVCC-HEHJ-FKA4",
			hash:          "6ab40d5eabe7218ec04182f461005600c7e3426bddd82cdb405bde9a1e0014b5",
		},
		{
			key:           "44ebe63375969fec2da67e87e7317946",
			humanReadable: "8LNYD-DVNKU-GYRC-E6GU-3YFD-CT86",
			hash:          "c1968cb4c013ea893f1922bb5c39f81e35012c0bd9ce3c01cc2a05873a2499e6",
		},
		{
			key:           "cb84798c3ca85a674194550a2e96aed8",
			humanReadable: "TF27L-31WN1-E6FH-DMAM-52X5-NFV0",
			hash:          "23cf8b7873425fd8efe31dc5b6ab9c357eb98a2a59c82ea1084ca8af58cc480a",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()

			key, _ := hex.DecodeString(tt.key)
			code, err := Create(key)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := code.Hash(), tt.hash; got != want {
				t.Errorf("Hash(): got %q, want %q", got, want)
			}

			if got, want := code.String(), tt.humanReadable; got != want {
				t.Errorf("Human Readable String: got %q, want %q", got, want)
			}
		})
	}
}
