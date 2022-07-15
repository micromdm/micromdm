// Package password provides utilities for creating and verifying macOS
// passwords for MDM account configuration.
package password

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"errors"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

// macKeyLen is the length of a PBKDF2 salt.
const macKeyLen = 128

// ErrNoMatch is returned by Verify if the password does not match.
var ErrNoMatch = errors.New("password does not match")

// SaltedSha512PBKDF2Dictionary is a SHA512 PBKDF2 dictionary.
type SaltedSHA512PBKDF2Dictionary struct {
	Iterations int    `plist:"iterations"`
	Salt       []byte `plist:"salt"`
	Entropy    []byte `plist:"entropy"`
}

// SaltedSHA512PBKDF2 creates a SALTED-SHA512-PBKDF2 dictionary
// from a plaintext password. The hash function will use a 128 bit
// salt.
func SaltedSHA512PBKDF2(plaintext string) (SaltedSHA512PBKDF2Dictionary, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return SaltedSHA512PBKDF2Dictionary{}, err
	}
	iterations, err := secureRandInt(20000, 40000)
	if err != nil {
		return SaltedSHA512PBKDF2Dictionary{}, err
	}
	return SaltedSHA512PBKDF2Dictionary{
		Iterations: iterations,
		Salt:       salt,
		Entropy: pbkdf2.Key([]byte(plaintext),
			salt, iterations, macKeyLen, sha512.New),
	}, nil
}

// Verify verifies a plaintext password against a existing SALTED-SHA512-PBKDF2
// password dictionary.
func Verify(plaintext string, h SaltedSHA512PBKDF2Dictionary) error {
	hashed := pbkdf2.Key([]byte(plaintext), h.Salt, h.Iterations, macKeyLen, sha512.New)
	if 1 != subtle.ConstantTimeCompare(h.Entropy, hashed) {
		return ErrNoMatch
	}
	return nil
}

// CCCalibratePBKDF uses a pseudorandom value returned within 100 milliseconds.
// Use a random int from crypto/rand between 20,000 and 40,000 instead.
func secureRandInt(min, max int64) (int, error) {
	var random int
	for {
		iter, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return 0, err
		}
		if iter.Int64() >= min {
			random = int(iter.Int64())
			break
		}
	}
	return random, nil
}
