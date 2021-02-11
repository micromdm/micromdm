// Package activationlock supports Apple device activation locking and unlocking.
// https://developer.apple.com/documentation/devicemanagement/activation_lock_a_device
package activationlock

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// The available character set for a bypass code string.
const charset = "0123456789ACDEFGHJKLMNPQRTUVWXYZ"

var (
	// Positions to insert "-" in the bypass code string.
	dashPositions = []int{5, 10, 14, 18, 22}

	// The salt for the hash is always the same.
	salt = []uint8{0, 0, 0, 0}
)

// BypassCode is the ActivationLock code.
type BypassCode struct {
	Key [16]byte

	format string
}

// Create generates a usable bypass code.
func Create(key []byte) (BypassCode, error) {
	var code BypassCode

	if key == nil {
		key = make([]byte, 16)
		if _, err := rand.Read(key); err != nil {
			return code, err
		}
	}

	copy(code.Key[:], key)

	// Format human readable version of key.
	values, err := convertBits(key, 8, 5)
	if err != nil {
		return code, err
	}

	dashIdx := 0
	var str strings.Builder
	for i, p := range values {
		if dashIdx < len(dashPositions) && i == dashPositions[dashIdx] {
			str.WriteString("-")
			dashIdx++
		}
		str.WriteByte(charset[p])
	}
	code.format = str.String()

	return code, nil
}

// Hash returns a PBKKDF2 derived hash of the bypass code.
// The hex encoded string is sent to Apple to lock a device.
func (c BypassCode) Hash() string {
	return hex.EncodeToString(pbkdf2.Key(c.Key[:], salt, 50000, sha256.Size, sha256.New))
}

// String returns the bypass code in human readable format.
func (c BypassCode) String() string {
	return c.format
}

// convert binary data from one bits-per-byte arrangement to another.
// Ex: re-arrange 8 bit bytes to groups of 5 when converting to base32.
// This is a modified helper from a Go implementation of the bech32 format.
// https://github.com/FiloSottile/age/blob/c9a35c072716b5ac6cd815366999c9e189b0c317/internal/bech32/bech32.go#L79-L105
func convertBits(data []byte, frombits, tobits byte) ([]byte, error) {
	var ret []byte
	acc := uint32(0)
	bits := byte(0)
	maxv := byte(1<<tobits - 1)
	for idx, value := range data {
		if value>>frombits != 0 {
			return nil, fmt.Errorf("invalid data range: data[%d]=%d (frombits=%d)", idx, value, frombits)
		}
		acc = acc<<frombits | uint32(value)
		bits += frombits
		for bits >= tobits {
			bits -= tobits
			ret = append(ret, byte(acc>>bits)&maxv)
		}
	}

	// zero out most significant bits of the last value, until we get to remaining bits
	if bits > 0 {
		for bit := frombits; bit >= bits; bit-- {
			acc = acc &^ (1 << bit)
		}
		ret = append(ret, byte(acc)&maxv)
	}

	return ret, nil
}
