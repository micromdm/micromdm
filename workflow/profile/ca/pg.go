package ca

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type Postgres struct {
	db            *sqlx.DB
	caKeyPassword []byte
}

func New(db *sqlx.DB) (*Postgres, error) {
	d := &Postgres{db: db}
	err := d.init()
	return d, err
}

const (
	tableName = "wf_profile_ca"
	caRowID   = "01CTC2SHTD0PYZ6S7Y0CZFJZKE"
	bits      = 4096
)

func (d *Postgres) init() error {
	_, _, err := d.CA(nil)
	if err == nil {
		return nil
	}
	if isNotFound(err) {
	} else if err != nil {
		return err
	}

	pkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return errors.Wrap(err, "init private key for CA")
	}

	var (
		authPkixName = pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"MicroMDM"},
			CommonName:   "MicroMDM Profile Workflow CA",
		}
		authTemplate = x509.Certificate{
			SerialNumber:          big.NewInt(1),
			Subject:               authPkixName,
			NotBefore:             time.Now().UTC(),
			NotAfter:              time.Now().AddDate(1, 0, 0).UTC(),
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			BasicConstraintsValid: true,
			IsCA: true,
		}
	)

	subjectKeyID, err := generateSubjectKeyID(&pkey.PublicKey)
	if err != nil {
		return err
	}
	authTemplate.SubjectKeyId = subjectKeyID

	crtBytes, err := x509.CreateCertificate(rand.Reader, &authTemplate, &authTemplate, &pkey.PublicKey, pkey)
	if err != nil {
		return err
	}

	pkeyBytes := x509.MarshalPKCS1PrivateKey(pkey)
	err = d.create(context.TODO(), crtBytes, pkeyBytes)
	return err
}

type rsaPublicKey struct {
	N *big.Int
	E int
}

func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)

	return hash[:], nil
}

type dbCA struct {
	ID               string    `db:"id"`
	CertificateBytes []byte    `db:"ca_certificate"`
	PrivateKeyBytes  []byte    `db:"ca_private_key"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func columns() []string {
	return []string{
		"id",
		"ca_certificate",
		"ca_private_key",
		"created_at",
		"updated_at",
	}
}

func (d *Postgres) create(ctx context.Context, cert []byte, key []byte) error {
	p := dbCA{
		ID:               caRowID,
		CertificateBytes: cert,
		PrivateKeyBytes:  key,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(columns()...).
		Values(
			p.ID,
			p.CertificateBytes,
			p.PrivateKeyBytes,
			p.CreatedAt,
			p.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building insert profile query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	if isAlreadyExistsPG(err) {
		return nil
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec device save in pg")
}

func (d *Postgres) CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error) {
	ctx := context.TODO()
	var ca dbCA
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(columns()...).
		From(tableName).
		Where(sq.Eq{"id": caRowID}).
		ToSql()
	if err != nil {
		return nil, nil, errors.Wrap(err, "building sql")
	}
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&ca)
	if err == sql.ErrNoRows {
		return nil, nil, notFoundErr{}
	} else if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(ca.CertificateBytes)
	if err != nil {
		return nil, nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(ca.PrivateKeyBytes)
	if err != nil {
		return nil, nil, err
	}
	chain := []*x509.Certificate{cert}
	return chain, key, nil
}

func (d *Postgres) Put(name string, crt *x509.Certificate) error {
	fmt.Println("PUT ", name)
	return nil
}

func (d *Postgres) Serial() (*big.Int, error) {
	return big.NewInt(time.Now().UTC().UnixNano()), nil
}

func (d *Postgres) HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error) {
	fmt.Println("HAS CN ", cn)
	return true, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string  { return "ca certificate not found" }
func (e notFoundErr) NotFound() bool { return true }

func isNotFound(err error) bool {
	err = errors.Cause(err)
	e, ok := errors.Cause(err).(interface {
		error
		NotFound() bool
	})
	return ok && e.NotFound()
}

func isAlreadyExistsPG(err error) bool {
	pqErr, ok := errors.Cause(err).(*pq.Error)
	return ok && pqErr.Code.Name() == "unique_violation"
}
