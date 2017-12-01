package backup

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"time"

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
	fileName := fmt.Sprintf("/var/db/micromdm/micromdm-backup-%v", time.Now().Unix())
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		stdlog.Println(err)
	}
	defer file.Close()

	svc.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(file)
		return err
	})

	if err := file.Close(); err != nil {
		stdlog.Println(err)
	}
	return nil
}
