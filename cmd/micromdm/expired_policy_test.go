package main

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"testing"
	"time"
)

type fakeDepot struct {
	ca  *x509.Certificate
	err error
}

func (f fakeDepot) CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error) {
	if f.err != nil {
		return nil, nil, f.err
	}
	return []*x509.Certificate{f.ca}, nil, nil
}

func TestExpiredDeviceCertPolicyFlag(t *testing.T) {
	now := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	ca := &x509.Certificate{}
	ok := fakeDepot{ca: ca}

	p, err := expiredDeviceCertPolicy("2027-01-31", ok, now)
	if err != nil || !p.Until.Equal(time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)) || p.Issuer != ca {
		t.Fatalf("date form: %v %v", p, err)
	}
	if p, err = expiredDeviceCertPolicy("2026-12-01T00:00:00Z", ok, now); err != nil || p.Until.Month() != 12 {
		t.Fatalf("rfc3339 form: %v %v", p, err)
	}
	for _, bad := range []string{"tomorrow", "2026-11-01", "2027-12-31"} {
		if _, err := expiredDeviceCertPolicy(bad, ok, now); err == nil {
			t.Errorf("%q should be rejected (invalid, past, or > 180 days)", bad)
		}
	}
	if _, err := expiredDeviceCertPolicy("2027-01-31", fakeDepot{err: errors.New("no ca")}, now); err == nil {
		t.Error("missing CA should fail")
	}
}
