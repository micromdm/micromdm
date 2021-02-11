package pg

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	//"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	//"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"context"

	"strings"
	"database/sql"

	//"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/micromdm/micromdm/platform/config"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type Depot struct {
	db *sqlx.DB
}

type SCEPCertificate struct {
	SCEPID int `db:"scep_id"`
	CertName string `db:"cert_name"`
	SCEPCert []byte `db:"scep_cert"`
}

func NewDB(db *sqlx.DB) (*Depot, error) {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS server_config (
			config_id INT PRIMARY KEY,
		    push_certificate bytea DEFAULT NULL,
		    private_key bytea DEFAULT NULL
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating server_config sql table failed")
	}

	// This Order is important, else the start value will be wrong...
	// https://asktom.oracle.com/pls/apex/f?p=100:11:0::::P11_QUESTION_ID:9529339800346302436
	// FIRST create Sequence with Start, then create referencing table, with integer instead of Serial
	_,err = db.Exec(`CREATE SEQUENCE IF NOT EXISTS scep_certificates_scep_id_seq
			INCREMENT 1
			MINVALUE 1
			NO MAXVALUE
			START 2;
			`)

	if err != nil {
	   return nil, errors.Wrap(err, "creating sequence with start = 2")
	}

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS scep_certificates (
			scep_id integer PRIMARY KEY DEFAULT nextval('scep_certificates_scep_id_seq'),
			cert_name TEXT NULL,
			scep_cert bytea DEFAULT NULL
		)`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating scep_certificates sql table failed")
	}

	store := &Depot{db: db}
	return store, err
}


func (d *Depot) CA(pass []byte) ([]*x509.Certificate, *rsa.PrivateKey, error) {
	ctx := context.Background()
	chain := []*x509.Certificate{}

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 4}).
		ToSql()

	var config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if err != nil {
		return nil, nil, err
	}
	var keyBytes, certBytes []byte
	keyBytes = config.PushCertificate
	certBytes = config.PrivateKey

	key, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		return chain, key, err
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return chain, key, err
	}
	chain = append(chain, cert)
	return chain, key, nil
}

func (d *Depot) Put(cn string, crt *x509.Certificate) error {
	// TODO
	if crt == nil || crt.Raw == nil {
		return fmt.Errorf("%q does not specify a valid certificate for storage", cn)
	}

	serial, err := d.Serial()
	if err != nil {
		return err
	}

	name := cn + "." + serial.String()

	// No need to get serial, we have an auto_increment in place
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("scep_certificates").
		Columns("cert_name", "scep_cert").
		Values(
			name,
			crt.Raw,
		).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building scep_certificates save query")
	}
	ctx := context.Background()
	_, err = d.db.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "exec scep_certificates save in pg")
}

type AutoIncrement struct {
	Index int64 `db:"currval"`
}

type AutoIncrement2 struct {
	Index int64 `db:"last_value"`
}

func (d *Depot) Serial() (*big.Int, error) {

	// currval only available, if we called nextval in the current session, but this would increment our sequence...
// 	query, args, err := sq.StatementBuilder.
// 		Select("currval('scep_certificates_scep_id_seq')").
// 		ToSql()

	query, args, err := sq.StatementBuilder.
		Select("last_value").
		From("scep_certificates_scep_id_seq").
		ToSql()

	if err != nil {
	   return nil, errors.Wrap(err, "retrieving serial sequence failed")
	}
	var auto_increment AutoIncrement2
	ctx := context.Background()
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&auto_increment)
	if errors.Cause(err) == sql.ErrNoRows {
		return big.NewInt(2), nil
	}
	return big.NewInt(auto_increment.Index), nil
}

func (d *Depot) HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error) {
	// TODO: implement allowTime
	// TODO: implement revocation
	if cert == nil {
		return false, errors.New("nil certificate provided")
	}

	scep_certs, err := d.listCertificates(context.Background(),cert.Subject.CommonName)
	if err != nil {
		return false, err
	}
	for _, v := range scep_certs {
		if bytes.Compare(v.SCEPCert, cert.Raw) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (d *Depot) listCertificates(ctx context.Context, prefix string) ([]SCEPCertificate, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("scep_id", "cert_name", "scep_cert").
		From("scep_certificates").
		Where("cert_name LIKE ?", fmt.Sprint("", prefix, "%")).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}
	var list []SCEPCertificate
	err = d.db.SelectContext(ctx, &list, query, args...)
	return list, errors.Wrap(err, "list scep certs")
}

func (d *Depot) CreateOrLoadKey(bits int) (*rsa.PrivateKey, error) {
	var (
		key *rsa.PrivateKey
		err error
	)

	ctx := context.Background()
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 4}).
		ToSql()

	var config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}
	var keyBytes []byte
	keyBytes = config.PushCertificate
	if keyBytes == nil {
		// if there is no certificate or private key then generate
		key, err = generateAndStoreKey(ctx, d, bits)
	} else {
		key, err = x509.ParsePKCS1PrivateKey(keyBytes)
	}
	return key, err
}

func generateAndStoreKey(ctx context.Context, d *Depot, bits int) (key *rsa.PrivateKey, err error) {
	key, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("server_config").
		Prefix("ON CONFLICT (config_id) DO").
		Set("config_id", 4).
		Set("push_certificate", x509.MarshalPKCS1PrivateKey(key)).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building update query for server_config save")
	}
	updateQuery = strings.Replace(updateQuery, "server_config", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("server_config").
		Columns("config_id", "push_certificate").
		Values(
			4,
			x509.MarshalPKCS1PrivateKey(key),
		).
		Suffix(updateQuery).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building server_config save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)

	return key, errors.Wrap(err, "exec server_config save in pg")
}

func (d *Depot) CreateOrLoadCA(key *rsa.PrivateKey, years int, org, country string) (*x509.Certificate, error) {
	var (
		cert *x509.Certificate
		err  error
	)
	ctx := context.Background()
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 4}).
		ToSql()

	var config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		fmt.Println("CreateOrLoadCA ERR")
		return nil, err
	}
	var certBytes []byte
	certBytes = config.PrivateKey
	if cert == nil {
		cert, err = generateAndStoreCA(ctx,d,key,years,org,country)
	} else {
		cert, err = x509.ParseCertificate(certBytes)
	}
	//fmt.Println(cert)

	return cert, err
}

func generateAndStoreCA(ctx context.Context, d *Depot, key *rsa.PrivateKey, years int, org string, country string) (cert *x509.Certificate, err error) {
	subject := pkix.Name{
		Country:            []string{country},
		Organization:       []string{org},
		OrganizationalUnit: []string{"MicroMDM SCEP CA"},
		Locality:           nil,
		Province:           nil,
		StreetAddress:      nil,
		PostalCode:         nil,
		SerialNumber:       "",
		CommonName:         org,
	}

	subjectKeyID, err := generateSubjectKeyID(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	authTemplate := x509.Certificate{
		SerialNumber:       big.NewInt(1),
		Subject:            subject,
		NotBefore:          time.Now().Add(-600).UTC(),
		NotAfter:           time.Now().AddDate(years, 0, 0).UTC(),
		KeyUsage:           x509.KeyUsageCertSign,
		ExtKeyUsage:        nil,
		UnknownExtKeyUsage: nil,

		BasicConstraintsValid: true,
		IsCA:                        true,
		MaxPathLen:                  0,
		SubjectKeyId:                subjectKeyID,
		DNSNames:                    nil,
		PermittedDNSDomainsCritical: false,
		PermittedDNSDomains:         nil,
	}

	crtBytes, err := x509.CreateCertificate(rand.Reader, &authTemplate, &authTemplate, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}


	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("server_config").
		Prefix("ON CONFLICT (config_id) DO").
		Set("config_id", 4).
		Set("private_key", crtBytes).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building update query for server_config save")
	}
	updateQuery = strings.Replace(updateQuery, "server_config", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("server_config").
		Columns("config_id", "private_key").
		Values(
			4,
			crtBytes,
		).
		Suffix(updateQuery).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building server_config save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)

	return x509.ParseCertificate(crtBytes)
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyID generates SubjectKeyId used in Certificate
// ID is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
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
