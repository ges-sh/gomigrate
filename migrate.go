// Package migrate provided easy way for applying PostgreSQL database migrations.
// Expected filename format for migrations is for example: 1548938237393.up.sql.
package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// DB represents database which allows Executing and Querying queries
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Migrate is executing migrations on database
type Migrate struct {
	db DB
}

// ErrNilDB is returned when nil database is passed to migrate
var ErrNilDB = errors.New("migrate: database cannot be nil")

// New returns new instance of Migrate
func New(db DB) (Migrate, error) {
	if db == nil {
		return Migrate{}, ErrNilDB
	}

	m := Migrate{db: db}

	m.createMigrationsTable()

	return m, nil
}

// MustNew returns new instance of Migrate and panics if there's any error
func MustNew(db DB) Migrate {
	m, err := New(db)
	if err != nil {
		panic(err)
	}
	return m
}

// Metadata contains information about single executed migration
type Metadata struct {
	ID        uint64
	CreatedAt time.Time
	Applied   bool
}

// Up executes migrations inside of dir
func (m Migrate) Up(dir string) (map[string]Metadata, error) {
	meta := make(map[string]Metadata)

	files, err := filepath.Glob(dir + "/*.up.sql")
	if err != nil {
		return meta, err
	}

	sort.Strings(files)

	lastID, err := m.lastMigrationID()
	if err != nil {
		return meta, err
	}

	var newMigrationsIndex int // specifies at what index new migrations starts

	for i, fileName := range files {
		migrationID, err := id(fileName)
		if err != nil {
			return meta, err
		}

		if migrationID == lastID { //migrations from this point are new
			newMigrationsIndex = i + 1
			break
		}
	}

	return meta, m.runMigrations(files[newMigrationsIndex:], meta)
}

// ErrorDirtyDatabase is returned when database contains dirty migration. In that case, migration needs to be fixed, removed from migrations table and applied again.
type ErrorDirtyDatabase struct {
	ID uint64
}

func (e ErrorDirtyDatabase) Error() string {
	return fmt.Sprintf("migrate: dirty migration with ID %v. Fix it, and try again", e.ID)
}

// lastMigrationID returns ID of the last applied migration. If the migration is dirty, ErrorDirtyDatabase error is returned.
func (m Migrate) lastMigrationID() (id uint64, err error) {
	var applied bool
	err = m.db.QueryRow(`SELECT id, applied FROM migrations ORDER BY created_at DESC LIMIT 1`).
		Scan(&id, &applied)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if !applied {
		return 0, ErrorDirtyDatabase{ID: id}
	}
	return
}

// runMigrations runs migrations provided within files slice.
func (m Migrate) runMigrations(files []string, meta map[string]Metadata) (err error) {
	for _, fileName := range files {
		meta[fileName], err = m.run(fileName)
		if err != nil {
			err = fmt.Errorf("error while applying %s migration: %s", fileName, err)
			return
		}
	}
	return
}

// run runs single migration.
func (m Migrate) run(path string) (meta Metadata, err error) {
	meta.ID, err = id(path)
	if err != nil {
		return
	}

	query, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	_, mErr := m.db.Exec(string(query))

	err = m.db.QueryRow(`INSERT INTO migrations(id, applied) VALUES($1, $2)
		RETURNING created_at, applied`,
		meta.ID, mErr == nil).Scan(&meta.CreatedAt, &meta.Applied)
	if err != nil {
		err = fmt.Errorf("migrate: %v: %v", err, mErr)
		return
	}

	return meta, mErr
}

// ErrInvalidFilename is returned when file name is in invalid format.
var ErrInvalidFilename = errors.New("migrate: unexpected file name format")

// id extracts migration id from filename
func id(fileName string) (uint64, error) {
	fileName = filepath.Base(fileName)

	if len(fileName) < 7 {
		return 0, ErrInvalidFilename
	}

	return strconv.ParseUint(fileName[:10], 10, 64)
}
