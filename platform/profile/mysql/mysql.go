package mysql

import (
	"context"
	"strings"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/micromdm/micromdm/platform/profile"
)

type Mysql struct{ db *sqlx.DB }

const tableName = "profiles"

func columns() []string {
	return []string{
		"identifier",
		"mobileconfig",
	}
}


func NewDB(db *sqlx.DB) (*Mysql, error) {
	// Required for TIMESTAMP DEFAULT 0
	_,err := db.Exec(`SET sql_mode = '';`)
	
	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS profiles (
			profile_id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
		    identifier TEXT DEFAULT NULL,
		    mobileconfig BLOB DEFAULT NULL
		);`)
	if err != nil {
	   return nil, errors.Wrap(err, "creating profile sql table failed")
	}
	
	return &Mysql{db: db}, nil
}

func (d *Mysql) List(ctx context.Context) ([]profile.Profile, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(columns()...).
		From(tableName).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}
	var list []profile.Profile
	err = d.db.SelectContext(ctx, &list, query, args...)
	return list, errors.Wrap(err, "list profiles")
}

func (d *Mysql) Save(ctx context.Context, p *profile.Profile) error {

	_, err := d.ProfileById(ctx,p.Identifier)
	// Empty object => insert
	if (err != nil) {
		
		query, args, err := sq.StatementBuilder.
			PlaceholderFormat(sq.Question).
			Insert(tableName).
			Columns(columns()...).
			Values(
				p.Identifier,
				p.Mobileconfig,
			).
			//Suffix(updateQuery).
			ToSql()
		
		var all_args = append(args, args...)
		if err != nil {
			return errors.Wrap(err, "building profile save query")
		}
		
		_, err = d.db.ExecContext(ctx, query, all_args...)
		
	} else {
		// Update existing entry
		updateQuery, args_update, err := sq.StatementBuilder.
			PlaceholderFormat(sq.Question).
			Update(tableName).
			//Prefix("ON DUPLICATE KEY").
			//Set("identifier", p.Identifier).
			Set("mobileconfig", p.Mobileconfig).
			Where("identifier LIKE ?", fmt.Sprint("", p.Identifier, "")).
			ToSql()
		if err != nil {
			return errors.Wrap(err, "building update query for device save")
		}
		
		// MySql Convention
		// Replace "ON DUPLICATE KEY UPDATE TABLE_NAME SET" to "ON DUPLICATE KEY UPDATE"
		updateQuery = strings.Replace(updateQuery, tableName+" SET ", "", -1)
		
		_, err = d.db.ExecContext(ctx, updateQuery, args_update...)
	}
	
	return errors.Wrap(err, "exec profile save in mysql")
}

func (d *Mysql) ProfileById(ctx context.Context, id string) (*profile.Profile, error) {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select(columns()...).
		From(tableName).
		Where("identifier LIKE ?", fmt.Sprint("", id, "")).
		//Where(sq.Eq{"identifier": id}).
		ToSql()
	
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	var p profile.Profile
	
	err = d.db.QueryRowxContext(ctx, query, args...).StructScan(&p)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, profileNotFoundErr{}
	}
	return &p, errors.Wrap(err, "finding profile by identifier")
}

func (d *Mysql) Delete(ctx context.Context, id string) error {
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Question).
		Delete(tableName).
		Where("identifier LIKE ?", fmt.Sprint("", id, "")).
		//Where(sq.Eq{"identifier": id}).		
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}
	_, err = d.db.ExecContext(ctx, query, args...)
	return errors.Wrap(err, "delete profile by identifier")
}

type profileNotFoundErr struct{}

func (e profileNotFoundErr) Error() string {
	return "profile not found"
}

func (e profileNotFoundErr) NotFound() bool {
	return true
}