package migrate

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var dbURL string

func init() {
	flag.StringVar(&dbURL, "dburl", "postgres://testuser:testpass@127.0.0.1/testdb?sslmode=disable", "URL to testing database")
	flag.Parse()
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestMigrate(t *testing.T) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer db.Close()

	m, err := New(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	meta, err := m.Up("migs")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(meta)
}

type database struct {
	// queries keeps all queries executed against database
	queries []string

	// err is an error we'd like to return
	err error
}

func (d *database) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.queries = append(d.queries, query)
	return nil, d.err
}

func (d database) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}
