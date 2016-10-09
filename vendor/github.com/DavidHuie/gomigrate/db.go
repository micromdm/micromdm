package gomigrate

import "strings"

type Migratable interface {
	SelectMigrationTableSql() string
	CreateMigrationTableSql() string
	GetMigrationSql() string
	MigrationLogInsertSql() string
	MigrationLogDeleteSql() string
	GetMigrationCommands(string) []string
}

// POSTGRES

type Postgres struct{}

func (p Postgres) SelectMigrationTableSql() string {
	return "SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1"
}

func (p Postgres) CreateMigrationTableSql() string {
	return `CREATE TABLE gomigrate (
                  id           SERIAL       PRIMARY KEY,
                  migration_id BIGINT       UNIQUE NOT NULL
                )`
}

func (p Postgres) GetMigrationSql() string {
	return `SELECT migration_id FROM gomigrate WHERE migration_id = $1`
}

func (p Postgres) MigrationLogInsertSql() string {
	return "INSERT INTO gomigrate (migration_id) values ($1)"
}

func (p Postgres) MigrationLogDeleteSql() string {
	return "DELETE FROM gomigrate WHERE migration_id = $1"
}

func (p Postgres) GetMigrationCommands(sql string) []string {
	return []string{sql}
}

// MYSQL

type Mysql struct{}

func (m Mysql) SelectMigrationTableSql() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_name = ? AND table_schema = (SELECT DATABASE())"
}

func (m Mysql) CreateMigrationTableSql() string {
	return `CREATE TABLE gomigrate (
                  id           INT          NOT NULL AUTO_INCREMENT,
                  migration_id BIGINT       NOT NULL UNIQUE,
                  PRIMARY KEY (id)
                )`
}

func (m Mysql) GetMigrationSql() string {
	return `SELECT migration_id FROM gomigrate WHERE migration_id = ?`
}

func (m Mysql) MigrationLogInsertSql() string {
	return "INSERT INTO gomigrate (migration_id) values (?)"
}

func (m Mysql) MigrationLogDeleteSql() string {
	return "DELETE FROM gomigrate WHERE migration_id = ?"
}

func (m Mysql) GetMigrationCommands(sql string) []string {
	count := strings.Count(sql, ";")
	commands := strings.SplitN(string(sql), ";", count)
	return commands
}

// MARIADB

type Mariadb struct {
	Mysql
}

// SQLITE3

type Sqlite3 struct{}

func (s Sqlite3) SelectMigrationTableSql() string {
	return "SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?"
}

func (s Sqlite3) CreateMigrationTableSql() string {
	return `CREATE TABLE gomigrate (
  id INTEGER PRIMARY KEY,
  migration_id INTEGER NOT NULL UNIQUE
)`
}

func (s Sqlite3) GetMigrationSql() string {
	return "SELECT migration_id FROM gomigrate WHERE migration_id = ?"
}

func (s Sqlite3) MigrationLogInsertSql() string {
	return "INSERT INTO gomigrate (migration_id) values (?)"
}

func (s Sqlite3) MigrationLogDeleteSql() string {
	return "DELETE FROM gomigrate WHERE migration_id = ?"
}

func (s Sqlite3) GetMigrationCommands(sql string) []string {
	return []string{sql}
}
