// Package user exposes methods for representing users in a database.
package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/scrypt"

	"micromdm.io/v2/pkg/id"
)

// User is a MicroMDM user.
type User struct {
	ID               string
	Username         string
	Email            string
	Password         []byte
	Salt             []byte
	ConfirmationHash *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func columns() []string {
	return []string{
		"id",
		"username",
		"email",
		"password",
		"salt",
		"confirmation_hash",
		"created_at",
		"updated_at",
	}
}

// IsConfirmed returns true after the user confirms registration to the site.
// The user creation process always creates a confirmation token,
// which is deleted when the user accepts, by following the URL in email.
// Confirmation methods other than email should be possible to implement.
func (u *User) IsConfirmed() bool { return u.ConfirmationHash == nil }

// ValidatePassword compares the plaintext password with a stored password hash.
func (u *User) ValidatePassword(plaintext string) error {
	derived, err := scrypt.Key([]byte(plaintext), u.Salt, scryptCost, 8, 1, 32)
	if err != nil {
		return err
	}

	if 1 != subtle.ConstantTimeCompare(u.Password, derived) {
		return errors.New("password does not match")
	}

	return nil
}

func (u *User) setPassword(plaintext string) error {
	if plaintext == "" {
		return errors.New("password cannot be empty")
	}

	salt, err := random(32, base64.StdEncoding)
	if err != nil {
		return err
	}

	derived, err := scrypt.Key([]byte(plaintext), salt, scryptCost, 8, 1, 32)
	if err != nil {
		return err
	}

	u.Salt = salt
	u.Password = derived
	return nil
}

func random(keySize int, enc *base64.Encoding) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, enc.EncodedLen(len(key)))
	enc.Encode(buf, key)
	return buf, nil
}

// DB

func create(username, email, password string) (*User, error) {
	u := &User{
		ID:       id.New(),
		Username: username,
		Email:    email,
	}

	confirmationHash, err := random(32, base64.RawURLEncoding)
	if err != nil {
		return nil, fmt.Errorf("create confirmation_hash (email %s): %w", email, err)
	}
	cs := string(confirmationHash)
	u.ConfirmationHash = &cs

	if err := u.setPassword(password); err != nil {
		return nil, fmt.Errorf("set password for user (email %s): %w", email, err)
	}

	return u, nil
}

var constraints = map[string]string{
	"chk_email_not_empty": "You must provide an email address.",
}
