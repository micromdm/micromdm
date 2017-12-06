package backup

import (
	"bytes"
	"context"

	"github.com/boltdb/bolt"
)

// type Service interface {
// 	Backup(ctx context.Context) (data io.Reader, err error)
// }

func NewDB(db *bolt.DB, path string) *BackupService {
	return &BackupService{
		db:   db,
		path: path,
	}
}

type BackupService struct {
	db   *bolt.DB
	path string
}

func (svc *BackupService) Backup(ctx context.Context) ([]byte, error) {
	var buf bytes.Buffer

	err := svc.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(&buf)
		return err
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
