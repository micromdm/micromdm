package mysql

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
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
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/platform/config"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

//type Mysql struct{ db *sqlx.DB }
type Depot struct {
	db *sqlx.DB
}

type SCEPCertificate struct {
	SCEPID int `db:"scep_id"`
	CertName string `db:"cert_name"`
	SCEPCert []byte `db:"scep_cert"`
}

func NewBoltDepot(db *sqlx.DB) (*Depot, error) {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)
	
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS server_config (
			config_id INT PRIMARY KEY,
		    push_certificate BLOB DEFAULT NULL,
		    private_key BLOB DEFAULT NULL
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating server_config sql table failed")
	}
	
	// TODO create table for serials/identities
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS scep_certificates (
			scep_id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			cert_name TEXT NULL,
			scep_cert BLOB DEFAULT NULL
		);`)
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
		PlaceholderFormat(sq.Question).
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
	
	
/*
	chain := []*x509.Certificate{}
	var key *rsa.PrivateKey
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		// get ca_certificate
		caCert := bucket.Get([]byte("ca_certificate"))
		if caCert == nil {
			return fmt.Errorf("no ca_certificate in bucket")
		}
		// we need to make a copy of the byte slice because the asn.Unmarshal
		// method called by ParseCertificate will retain a reference to the original.
		// The slice should no longer be referenced once the BoltDB transaction is closed.
		caCertBytes := append([]byte(nil), caCert...)
		cert, err := x509.ParseCertificate(caCertBytes)
		if err != nil {
			return err
		}
		chain = append(chain, cert)

		// get ca_key
		caKey := bucket.Get([]byte("ca_key"))
		if caKey == nil {
			return fmt.Errorf("no ca_key in bucket")
		}
		key, err = x509.ParsePKCS1PrivateKey(caKey)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return chain, key, nil
*/
}

func (d *Depot) Put(cn string, crt *x509.Certificate) error {
	// TODO
	if crt == nil || crt.Raw == nil {
		return fmt.Errorf("%q does not specify a valid certificate for storage", cn)
	}
	
	// No need to get serial, we have an auto_increment in place
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("scep_certificates").
		Columns("cert_name", "scep_cert").
		Values(
			cn,
			crt.Raw,
		).
		ToSql()
	
	if err != nil {
		return errors.Wrap(err, "building scep_certificates save query")
	}
	ctx := context.Background()
	_, err = d.db.ExecContext(ctx, query, args...)
	
	return errors.Wrap(err, "exec scep_certificates save in mysql")
}

func (db *Depot) Serial() (*big.Int, error) {
	return big.NewInt(2), nil
/*
	fmt.Println("--- scep/builtin/db.go Serial()")
		
	s := big.NewInt(2)
	if !db.hasKey([]byte("serial")) {
		if err := db.writeSerial(s); err != nil {
			return nil, err
		}
		fmt.Println(s)
		return s, nil
	}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(certBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found!", certBucket)
		}
		k := bucket.Get([]byte("serial"))
		if k == nil {
			return fmt.Errorf("key %q not found", "serial")
		}
		s = s.SetBytes(k)
		fmt.Println(s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(s)
	return s, nil
*/
}

func (d *Depot) HasCN(cn string, allowTime int, cert *x509.Certificate, revokeOldCertificate bool) (bool, error) {
	fmt.Println("--- scep/builtin/db.go HasCN")
	fmt.Println(cn)
	fmt.Println(cert.Subject.CommonName)
	
	// TODO: implement allowTime
	// TODO: implement revocation
	if cert == nil {
		return false, errors.New("nil certificate provided")
	}
	
	prefix := []byte(cert.Subject.CommonName)
		
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select("scep_id", "cert_name", "scep_cert").
		From("scep_certificates").
		Where("cert_name LIKE ?", fmt.Sprint("", prefix, "%")).
		ToSql()
		
	if err != nil {
		return false, errors.Wrap(err, "building sql")
	}

	var scep_cert SCEPCertificate
	ctx := context.Background()
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&scep_cert)
	if errors.Cause(err) == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (d *Depot) CreateOrLoadKey(bits int) (*rsa.PrivateKey, error) {
	var (
		key *rsa.PrivateKey
		err error
	)
	
	ctx := context.Background()
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
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
	fmt.Println("--- scep/builtin/mysql.go CreateOrLoadKey")
	fmt.Println("ca_key")
	fmt.Println(x509.MarshalPKCS1PrivateKey(key))
	return key, err
}

func generateAndStoreKey(ctx context.Context, d *Depot, bits int) (key *rsa.PrivateKey, err error) {
	key, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update("server_config").
		Prefix("ON DUPLICATE KEY").
		Set("config_id", 4).
		Set("push_certificate", x509.MarshalPKCS1PrivateKey(key)).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building update query for server_config save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, "server_config SET ", "", -1)
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("server_config").
		Columns("config_id", "push_certificate").
		Values(
			4,
			x509.MarshalPKCS1PrivateKey(key),
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return nil, errors.Wrap(err, "building server_config save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return key, errors.Wrap(err, "exec server_config save in mysql")
}

func (d *Depot) CreateOrLoadCA(key *rsa.PrivateKey, years int, org, country string) (*x509.Certificate, error) {
	var (
		cert *x509.Certificate
		err  error
	)
	ctx := context.Background()
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
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
	fmt.Println("--- 	scep/builtin/db.go CreateOrLoadKey")
	fmt.Println("ca_certificate")
	fmt.Println(certBytes)

	return cert, err
}

func generateAndStoreCA(ctx context.Context, d *Depot, key *rsa.PrivateKey, years int, org string, country string) (cert *x509.Certificate, err error) {
	subject := pkix.Name{
		Country:            []string{country},
		Organization:       []string{org},
		OrganizationalUnit: []string{"Abacus Research AG SCEP CA"},
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

	
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update("server_config").
		Prefix("ON DUPLICATE KEY").
		Set("config_id", 4).
		Set("private_key", crtBytes).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building update query for server_config save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, "server_config SET ", "", -1)
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("server_config").
		Columns("config_id", "private_key").
		Values(
			4,
			crtBytes,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	if err != nil {
		return nil, errors.Wrap(err, "building server_config save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
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
	fmt.Println("--- scep/builtin/db.go generateSubjectKeyID")
	fmt.Println(pub)
	
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