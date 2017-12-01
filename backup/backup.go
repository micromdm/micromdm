package backup

import (
	"context"
	"fmt"

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

func (svc *BackupService) Backup(ctx context.Context) error {
	fmt.Println("backing up")
	return nil
}
