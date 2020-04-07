package pg

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/micromdm/micromdm/pkg/crypto"

	"github.com/go-kit/kit/log"
	"github.com/kolide/kit/dbutil"
	_ "github.com/lib/pq"
)

// createDepot creates a Bolt database in a temporary location.
func createDB(t *testing.T) *Depot {
	db, err := dbutil.OpenDBX(
		"postgres",
		"host=localhost port=5432 user=micromdm dbname=micromdm_test password=micromdm sslmode=disable",
		dbutil.WithLogger(log.NewNopLogger()),
		dbutil.WithMaxAttempts(1),
	)
	
	if err != nil {
		t.Fatal(err)
	}
	
	// Clean up first
	_, _ = db.Exec(`DROP TABLE IF EXISTS scep_certificates;`)
	_, _ = db.Exec(`DROP SEQUENCE IF EXISTS scep_certificates_scep_id_seq CASCADE;`)

	d, err := NewDB(db)
	if err != nil {
		t.Fatal(err)
	}
	
	//return nil
	return d
}

func TestDepot_Serial(t *testing.T) {
	db := createDB(t)
	tests := []struct {
		name    string
		want    *big.Int
		wantErr bool
	}{
		{
			name: "two is the default value.",
			want: big.NewInt(2),
		},
	}
	for _, tt := range tests {
		got, err := db.Serial()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Depot.Serial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Depot.Serial() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDepot_Put(t *testing.T) {
	db := createDB(t)
	
	_, cert, err := crypto.SimpleSelfSignedRSAKeypair("micromdm-dep-token", 365)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    *big.Int
		wantErr bool
	}{
		{
			name: "two is the default value.",
			want: big.NewInt(2),
		},
		{
			name: "After Put, expecting increment",
			want: big.NewInt(3),
		},
	}
	for _, tt := range tests {
		got, err := db.Serial()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Depot.Serial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err := db.Put("cn", cert); (err != nil) != tt.wantErr {
			t.Errorf("%q. Depot.Serial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		
		got, err = db.Serial()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Depot.Serial() = %v, want %v", tt.name, got, tt.want)
		}
	}
}


/*
func TestDepot_writeSerial(t *testing.T) {
	db := createDB(t)
	type args struct {
		s *big.Int
	}
	tests := []struct {
		name    string
		args    *big.Int
		wantErr bool
	}{
		{
			args: big.NewInt(5),
		},
		{
			args: big.NewInt(3),
		},
	}
	for _, tt := range tests {
		if err := db.writeSerial(tt.args); (err != nil) != tt.wantErr {
			t.Errorf("%q. Depot.writeSerial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
*/

/*
func TestDepot_incrementSerial(t *testing.T) {
	db := createDB(t)
	type args struct {
		s *big.Int
	}
	tests := []struct {
		name    string
		args    *big.Int
		want    *big.Int
		wantErr bool
	}{
		{
			args: big.NewInt(2),
			want: big.NewInt(3),
		},
		{
			args: big.NewInt(3),
			want: big.NewInt(4),
		},
	}
	for _, tt := range tests {
		if err := db.incrementSerial(tt.args); (err != nil) != tt.wantErr {
			t.Errorf("%q. Depot.incrementSerial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		got, _ := db.Serial()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Depot.Serial() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
*/

func TestDepot_CreateOrLoadKey(t *testing.T) {
	db := createDB(t)
	
	tests := []struct {
		bits    int
		wantErr bool
	}{
		{
			bits: 1024,
		},
		{
			bits: 2048,
		},
	}
	for i, tt := range tests {
		if _, err := db.CreateOrLoadKey(tt.bits); (err != nil) != tt.wantErr {
			t.Errorf("%d. Depot.CreateOrLoadKey() error = %v, wantErr %v", i, err, tt.wantErr)
		}
	}
}

func TestDepot_CreateOrLoadCA(t *testing.T) {
	db := createDB(t)
	tests := []struct {
		wantErr bool
	}{
		{},
		{},
	}
	for i, tt := range tests {
		key, err := db.CreateOrLoadKey(1024)
		if err != nil {
			t.Fatalf("%d. Depot.CreateOrLoadKey() error = %v", i, err)
		}
		if key == nil {
			panic("Expected public key not to be nil when creating new Key")
		}

		if ca, err := db.CreateOrLoadCA(key, 10, "MicroMDM", "US"); (err != nil) != tt.wantErr {
			if ca == nil {
				panic("Expected private key not to be nil when creating new CA for Key")
			}
			t.Errorf("%d. Depot.CreateOrLoadCA() error = %v, wantErr %v", i, err, tt.wantErr)
		}
	}
}
