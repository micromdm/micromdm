package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/subtle"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/smallstep/pkcs7"
)

// ExpiredDeviceCertPolicy is a temporary, time-boxed exception that lets devices whose MDM
// identity certificate has already expired reach the server again, so they can be sent a
// renewal enrollment profile (InstallProfile → new SCEP certificate → re-Authenticate).
//
// iOS never renews the MDM identity certificate on its own. Without this exception a device
// with an expired certificate is rejected before any command can be delivered and can only
// be recovered by re-enrolling by hand.
//
// The exception is narrow: the CMS signature over the request must still be valid, the
// signer certificate must be issued by our SCEP CA, it must actually be expired (not
// "not yet valid"), and the policy switches itself off at Until. The UDID ↔ certificate
// binding (platform/device/udidauth.go) still applies after this check.
type ExpiredDeviceCertPolicy struct {
	// Until is when the exception ends. Zero means disabled.
	Until time.Time
	// Issuer is the SCEP CA certificate that must have signed the device certificate.
	Issuer *x509.Certificate
}

// Active reports whether the exception applies at now.
func (p *ExpiredDeviceCertPolicy) Active(now time.Time) bool {
	return p != nil && p.Issuer != nil && !p.Until.IsZero() && now.Before(p.Until)
}

// IsSignerValidityError reports whether err is the pkcs7 error for a signer certificate
// whose validity window does not contain the signing time.
func IsSignerValidityError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "is outside of certificate validity")
}

// VerifyExpiredSigner re-verifies a parsed CMS signature (p7.Content must hold the detached
// content) whose only failure was the signer certificate's validity, applying the policy.
// It returns the signer certificate only when every other check passes.
func (p *ExpiredDeviceCertPolicy) VerifyExpiredSigner(p7 *pkcs7.PKCS7, now time.Time) (*x509.Certificate, error) {
	if !p.Active(now) {
		return nil, errors.New("expired device certificate exception is not active")
	}
	if p7 == nil {
		return nil, errors.New("missing CMS signature")
	}
	signer := p7.GetOnlySigner()
	if signer == nil {
		return nil, errors.New("missing or ambiguous signer certificate")
	}
	if !now.After(signer.NotAfter) {
		return nil, errors.New("signer certificate is not expired")
	}
	if now.Before(signer.NotBefore) {
		return nil, errors.New("signer certificate is not yet valid")
	}
	if !bytes.Equal(signer.RawIssuer, p.Issuer.RawSubject) {
		return nil, errors.New("signer certificate was not issued by the SCEP CA")
	}
	if err := signer.CheckSignatureFrom(p.Issuer); err != nil {
		return nil, fmt.Errorf("signer certificate not signed by the SCEP CA: %w", err)
	}
	if err := VerifyIgnoringSignerValidity(p7, signer); err != nil {
		return nil, err
	}
	return signer, nil
}

// signedAttribute mirrors the pkcs7 package's (unexported) attribute type field for field, so
// the signed attributes re-encode to exactly the bytes the package verifies for valid signers.
type signedAttribute struct {
	Type  asn1.ObjectIdentifier
	Value asn1.RawValue `asn1:"set"`
}

var oidAttrMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}

func hashForOID(oid asn1.ObjectIdentifier) (crypto.Hash, error) {
	switch {
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}):
		return crypto.SHA1, nil
	case oid.Equal(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}):
		return crypto.SHA256, nil
	case oid.Equal(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}):
		return crypto.SHA384, nil
	case oid.Equal(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}):
		return crypto.SHA512, nil
	}
	return 0, fmt.Errorf("unsupported digest algorithm %v", oid)
}

// VerifyIgnoringSignerValidity repeats the pkcs7 package's signature verification for a single
// signer except for the signer certificate's validity window: the messageDigest attribute must
// match p7.Content and the signature over the signed attributes must verify with the signer's
// public key. The pkcs7 package has already converted iOS's BER encoding to DER when parsing.
func VerifyIgnoringSignerValidity(p7 *pkcs7.PKCS7, signer *x509.Certificate) error {
	if len(p7.Signers) != 1 {
		return fmt.Errorf("expected exactly one signer, got %d", len(p7.Signers))
	}
	si := p7.Signers[0]
	if len(si.AuthenticatedAttributes) == 0 {
		return errors.New("signer has no signed attributes")
	}
	hash, err := hashForOID(si.DigestAlgorithm.Algorithm)
	if err != nil {
		return err
	}
	attrs := make([]signedAttribute, len(si.AuthenticatedAttributes))
	var digest []byte
	for i, a := range si.AuthenticatedAttributes {
		attrs[i] = signedAttribute{Type: a.Type, Value: a.Value}
		if a.Type.Equal(oidAttrMessageDigest) {
			if _, err := asn1.Unmarshal(a.Value.Bytes, &digest); err != nil {
				return fmt.Errorf("parse messageDigest: %w", err)
			}
		}
	}
	if digest == nil {
		return errors.New("signed attributes have no messageDigest")
	}
	h := hash.New()
	h.Write(p7.Content)
	if subtle.ConstantTimeCompare(digest, h.Sum(nil)) != 1 {
		return errors.New("message digest mismatch")
	}

	// same encoding as pkcs7.marshalAttributes: the attributes as a DER SET
	enc, err := asn1.Marshal(struct {
		A []signedAttribute `asn1:"set"`
	}{A: attrs})
	if err != nil {
		return fmt.Errorf("encode signed attributes: %w", err)
	}
	var outer asn1.RawValue
	if _, err := asn1.Unmarshal(enc, &outer); err != nil {
		return fmt.Errorf("encode signed attributes: %w", err)
	}
	h = hash.New()
	h.Write(outer.Bytes)
	hashed := h.Sum(nil)

	switch pub := signer.PublicKey.(type) {
	case *rsa.PublicKey:
		if err := rsa.VerifyPKCS1v15(pub, hash, hashed, si.EncryptedDigest); err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(pub, hashed, si.EncryptedDigest) {
			return errors.New("signature verification failed")
		}
	default:
		return fmt.Errorf("unsupported signer key type %T", signer.PublicKey)
	}
	return nil
}
