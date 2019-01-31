Migrations for PostgreSQL created using Go language.

[![GoDoc](https://godoc.org/encoding/json?status.svg)](https://godoc.org/encoding/json)

Example usage:
```go
package migrate_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	migrate "github.com/ges-sh/gomigrate"
)

func Example(t *testing.T) {
	db, err := sql.Open("postgres", "dburl")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer db.Close()

	m := migrate.MustNew(db)

	meta, err := m.Up("migrationsDir")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(meta) // example output: map[migrationsDir/1548931630834.up.sql:{1548931630834 2019-01-31 12:37:40.307724 +0000 +0000 true} migrationsDir/1548938237393.up.sql:{1548938237393 2019-01-31 12:37:40.31107 +0000 +0000 true}]
}
```
