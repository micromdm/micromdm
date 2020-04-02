package mysql

import (
	"context"
	"testing"
	"fmt"
	"net/url"
	"encoding/json"
	
	"io/ioutil"
	
	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micromdm/micromdm/dep"
	
	"github.com/fullsailor/pkcs7"
	//"github.com/micromdm/micromdm/platform/pubsub"
	"github.com/micromdm/micromdm/platform/pubsub/inmem"
	"github.com/micromdm/micromdm/platform/config"
)

func Test_SavePushCert(t *testing.T) {
	db,err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	
	cert := []byte("Cert")
	key := []byte("key")
	err = db.SavePushCertificate(ctx, cert, key)
	if err != nil {
		t.Fatal(err)
	}
	
	//serverConfig := 
}

func Test_ApplyWrongDEPToken(t *testing.T) {
	db,err := setup(t)
	ctx := context.Background()

	p7mBytes, err := ioutil.ReadFile("dep_token.p7m")
	if err != nil {
		t.Fatal(err)
	}
	
	unwrapped, err := config.UnwrapSMIME(p7mBytes)
	if err != nil {
		t.Fatal(err)
	}
	key, cert, err := db.DEPKeypair(ctx)
	if err != nil {
		t.Fatal(err)
	}

	p7, err := pkcs7.Parse(unwrapped)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := p7.Decrypt(cert, key)
	// expecting an error, because the token should not match to the generated one!
	if err == nil {
		t.Fatal()
	}
	if decrypted != nil {
		t.Fatal()
	}
}

func Test_GetNewDEPToken(t *testing.T) {
	db,err := setup(t)
	ctx := context.Background()

	_, cert, err := db.DEPKeypair(ctx)
	if err != nil {
		t.Fatal(err)
	}
	
	var certBytes []byte
	if cert != nil {
		certBytes = cert.Raw
	}

	tokens, err := db.DEPTokens(ctx)
	if err != nil {
		t.Fatal(err)
		//return nil, certBytes, err
	}

	if certBytes == nil {
		t.Fatal()
	}
	
	// Expecting to currently have no tokens
	if tokens != nil {
		t.Fatal()
	}
}

func Test_DEPKeyPairGeneratedOnce(t *testing.T) {
	db,err := setup(t)
	ctx := context.Background()

	key, cert, err := db.DEPKeypair(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if (key == nil || cert == nil) {
		t.Fatal()
	}
	
	key2, cert2, err := db.DEPKeypair(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if (key2 == nil || cert2 == nil) {
		t.Fatal()
	}
	
	keyDescr := fmt.Sprintf("%d", key)
	key2Descr := fmt.Sprintf("%d", key2)
	certDescr := fmt.Sprintf("%d", cert.Raw)
	cert2Descr := fmt.Sprintf("%d", cert2.Raw)
	if (keyDescr != key2Descr) || (certDescr != cert2Descr) {
		t.Fatal()
	}
}

func Test_ApplyDEPToken(t *testing.T) {
	return
	db,err := setup(t)
	ctx := context.Background()

	p7mBytes, err := ioutil.ReadFile("dep_token.p7m")
	if err != nil {
		t.Fatal(err)
	}
	
	unwrapped, err := config.UnwrapSMIME(p7mBytes)
	if err != nil {
		t.Fatal(err)
	}
	key, cert, err := db.DEPKeypair(ctx)
	if err != nil {
		t.Fatal(err)
	}

	p7, err := pkcs7.Parse(unwrapped)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := p7.Decrypt(cert, key)
	if err != nil {
		t.Fatal(err)
	}
	tokenJSON, err := config.UnwrapTokenJSON(decrypted)
	if err != nil {
		t.Fatal(err)
	}
	var depToken config.DEPToken
	err = json.Unmarshal(tokenJSON, &depToken)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddToken(ctx, depToken.ConsumerKey, tokenJSON)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("stored DEP token with ck", depToken.ConsumerKey)
}

func Test_DEPToken(t *testing.T) {
	db,err := setup(t)
	ctx := context.Background()
	
	var (
		conf           dep.OAuthParameters
		hasTokenConfig bool
		opts           []dep.Option
	)

	// try getting the oauth config from bolt
	tokens, err := db.DEPTokens(ctx)
	if err != nil {
		t.Fatal(err)
	}
	
	if len(tokens) >= 1 {
		hasTokenConfig = true
		conf.ConsumerSecret = tokens[0].ConsumerSecret
		conf.ConsumerKey = tokens[0].ConsumerKey
		conf.AccessSecret = tokens[0].AccessSecret
		conf.AccessToken = tokens[0].AccessToken
		// TODO: handle expiration
	}

	
	hasTokenConfig = true
	conf = dep.OAuthParameters{
		ConsumerKey:    "CK_48dd68d198350f51258e885ce9a5c37ab7f98543c4a697323d75682a6c10a32501cb247e3db08105db868f73f2c972bdb6ae77112aea803b9219eb52689d42e6",
		ConsumerSecret: "CS_34c7b2b531a600d99a0e4edcf4a78ded79b86ef318118c2f5bcfee1b011108c32d5302df801adbe29d446eb78f02b13144e323eb9aad51c79f01e50cb45c3a68",
		AccessToken:    "AT_927696831c59ba510cfe4ec1a69e5267c19881257d4bca2906a99d0785b785a6f6fdeb09774954fdd5e2d0ad952e3af52c6d8d2f21c924ba0caf4a031c158b89",
		AccessSecret:   "AS_c31afd7a09691d83548489336e8ff1cb11b82b6bca13f793344496a556b1f4972eaff4dde6deb5ac9cf076fdfa97ec97699c34d515947b9cf9ed31c99dded6ba",
	}
	depsimurl, err := url.Parse("https://test.mdm.com/")
	if err != nil {
		t.Fatal(err)
	}
	opts = append(opts, dep.WithServerURL(depsimurl))
	

	if !hasTokenConfig {
		t.Fatal()
	}
}

func setup(t *testing.T) (*Mysql, error) {
	// https://stackoverflow.com/a/23550874/464016
	db, err := dbutil.OpenDBX(
		"mysql",
		"micromdm:micromdm@tcp(127.0.0.1:3306)/micromdm_test?parseTime=true",
		//"host=127.0.0.1 port=3306 user=micromdm dbname=micromdm_test password=micromdm sslmode=disable",
		dbutil.WithLogger(log.NewNopLogger()),
		dbutil.WithMaxAttempts(1),
	)
	
	if err != nil {
		t.Fatal(err)
	}
	return NewDB(db, inmem.NewPubSub())
}
