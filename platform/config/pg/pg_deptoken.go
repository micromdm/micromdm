package pg

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"

	"strings"
	"database/sql"

	"fmt"

	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
	
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/config"
)

const dep_tableName = "push_certificate"

func dep_columns() []string {
	return []string{
		"push_certificate",
		"private_key",
//		"consumer_key",
//		"consumer_secret",
//		"access_token",
//		"access_secret",
//		"access_token_expiry",
	}
}

func (d *Postgres) AddToken(ctx context.Context, consumerKey string, json []byte) error {
	
	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("server_config").
		Prefix("ON CONFLICT (config_id) DO").
		Set("config_id", 3).
		Set("push_certificate", []byte(consumerKey)).
		Set("private_key", json).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for dep_tokens save")
	}
	
	updateQuery = strings.Replace(updateQuery, "", "server_config", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("server_config").
		Columns(columns()...).
		Values(
			3,
			[]byte(consumerKey),
			json,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building dep_tokens save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "exec dep_tokens save in pg")
}

func (d *Postgres) DEPTokens(ctx context.Context) ([]config.DEPToken, error) {
	var result []config.DEPToken
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 3}).
		// If multiple keys need to be supported, use separate table or use something like
		// Where("push_certificate LIKE ?", fmt.Sprint("%CK_")).
		ToSql()

	var server_config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&server_config)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, nil
		}
		return result, err
	}

	
	var depToken config.DEPToken
	err = json.Unmarshal(server_config.PrivateKey, &depToken)
	if err != nil {
		// TODO: log problematic DEP token, or remove altogether?
		fmt.Println("Cannot Unmarshal Private Key of DEP Token")
	}
	result = append(result, depToken)
	return result, nil
}

func (d *Postgres) DEPKeypair(ctx context.Context) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 2}).
		ToSql()

	var config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return
	}
	var keyBytes, certBytes []byte
	keyBytes = config.PushCertificate
	certBytes = config.PrivateKey
	
	if keyBytes == nil || certBytes == nil {
		// if there is no certificate or private key then generate
		key, cert, err = generateAndStoreDEPKeypair(ctx, d)
	} else {
		key, err = x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			return
		}
		cert, err = x509.ParseCertificate(certBytes)
		if err != nil {
			return
		}
	}
	return 
}

func generateAndStoreDEPKeypair(ctx context.Context, d *Postgres) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, cert, err = crypto.SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		return
	}

	pkBytes := x509.MarshalPKCS1PrivateKey(key)
	certBytes := cert.Raw
	updateQuery, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("server_config").
		Prefix("ON CONFLICT (config_id) DO").
		Set("config_id", 2).
		Set("push_certificate", pkBytes).
		Set("private_key", certBytes).
		ToSql()
	if err != nil {
		return key, cert, errors.Wrap(err, "building update query for dep_tokens save")
	}
	
	updateQuery = strings.Replace(updateQuery, "server_config", "", -1)

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("server_config").
		Columns("config_id", "push_certificate", "private_key").
		Values(
			2,
			pkBytes,
			certBytes,
		).
		Suffix(updateQuery).
		ToSql()
	if err != nil {
		return key, cert, errors.Wrap(err, "building dep_tokens save query")
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return key, cert, errors.Wrap(err, "building dep_tokens save query")
}
