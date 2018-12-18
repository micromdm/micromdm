package mysql

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"

	"strings"
	"database/sql"

	"fmt"

	"github.com/pkg/errors"
	_ "github.com/go-sql-driver/mysql"
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

func (d *Mysql) AddToken(ctx context.Context, consumerKey string, json []byte) error {
	
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update("server_config").
		Prefix("ON DUPLICATE KEY").
		Set("config_id", 3).
		Set("push_certificate", []byte(consumerKey)).
		Set("private_key", json).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for dep_tokens save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(tableName).
		Columns("config_id", "push_certificate", "private_key").
		Values(
			3,
			[]byte(consumerKey),
			json,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return errors.Wrap(err, "building dep_tokens save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return errors.Wrap(err, "exec dep_tokens save in mysql")
}

func (d *Mysql) DEPTokens(ctx context.Context) ([]config.DEPToken, error) {
	var result []config.DEPToken
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
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

func (d *Mysql) DEPKeypair(ctx context.Context) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select("push_certificate", "private_key").
		From("server_config").
		Where(sq.Eq{"config_id": 2}).
		ToSql()
	if err != nil {
		return
	}

	var config config.ServerConfig
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&config)
	if errors.Cause(err) == sql.ErrNoRows {
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

func generateAndStoreDEPKeypair(ctx context.Context, d *Mysql) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, cert, err = crypto.SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		return
	}

	pkBytes := x509.MarshalPKCS1PrivateKey(key)
	certBytes := cert.Raw
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update("server_config").
		Prefix("ON DUPLICATE KEY").
		Set("config_id", 2).
		Set("push_certificate", pkBytes).
		Set("private_key", certBytes).
		ToSql()
	if err != nil {
		return key, cert, errors.Wrap(err, "building update query for dep_tokens save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, "server_config"+" SET ", "", -1)
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("server_config").
		Columns("config_id", "push_certificate", "private_key").
		Values(
			2,
			pkBytes,
			certBytes,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	if err != nil {
		return key, cert, errors.Wrap(err, "building dep_tokens save query")
	}
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return key, cert, errors.Wrap(err, "building dep_tokens save query")
}
