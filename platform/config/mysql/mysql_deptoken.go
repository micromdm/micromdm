package mysql

import (
	"context"
//	"bytes"
	"crypto/rsa"
	"crypto/x509"
//	"encoding/json"

	"fmt"

//	"github.com/go-kit/kit/log"
//	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/pkg/crypto"
	"github.com/micromdm/micromdm/platform/config"
)

const dep_tableName = "dep_tokens"

func dep_columns() []string {
	return []string{
		"key",
		"private_key",
	}
}

func (d *Mysql) AddToken(ctx context.Context, consumerKey string, json []byte) error {
	
	fmt.Println(consumerKey)
	fmt.Println(json)
	
	return nil
	
/*
	updateQuery, args_update, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Update(tableName).
		Prefix("ON DUPLICATE KEY").
		Set("key", consumerKey).
		Set("value", json).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building update query for server_config save")
	}
	// MySql Convention
	// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
	updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)
	
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert(tableName).
		Columns(columns()...).
		Values(
			cert,
			key,
		).
		Suffix(updateQuery).
		ToSql()
	
	var all_args = append(args, args_update...)
	
	if err != nil {
		return errors.Wrap(err, "building server_config save query")
	}
	
	_, err = d.db.ExecContext(ctx, query, all_args...)
	
	return errors.Wrap(err, "exec server_config save in mysql")
	
	
	
	
	
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(depTokenBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(consumerKey), json)
	})
	if err != nil {
		return err
	}
	err = db.Publisher.Publish(context.TODO(), config.DEPTokenTopic, json)
	return err
*/
}

func (d *Mysql) DEPTokens(ctx context.Context) ([]config.DEPToken, error) {

	var result []config.DEPToken
	return result, nil
/*
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		prefix := []byte("CK_")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var depToken config.DEPToken
			err := json.Unmarshal(v, &depToken)
			if err != nil {
				// TODO: log problematic DEP token, or remove altogether?
				continue
			}
			result = append(result, depToken)
		}
		return nil
	})
	return result, err
*/
}

func (d *Mysql) DEPKeypair(ctx context.Context) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
/*
	var keyBytes, certBytes []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(depTokenBucket))
		if b == nil {
			return nil
		}
		keyBytes = b.Get([]byte("key"))
		certBytes = b.Get([]byte("certificate"))
		return nil
	})
	if err != nil {
		return
	}
	if keyBytes == nil || certBytes == nil {
		// if there is no certificate or private key then generate
		key, cert, err = generateAndStoreDEPKeypair(db)
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
*/
	return
}

func generateAndStoreDEPKeypair(d *Mysql) (ctx context.Context, key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, cert, err = crypto.SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		return
	}

	pkBytes := x509.MarshalPKCS1PrivateKey(key)
	certBytes := cert.Raw
	fmt.Println(pkBytes)
	fmt.Println(certBytes)
/*
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(depTokenBucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte("key"), pkBytes)
		if err != nil {
			return err
		}
		err = b.Put([]byte("certificate"), certBytes)
		if err != nil {
			return err
		}
		return nil
	})
*/

	return
}
