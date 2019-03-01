package migrate

var migrationsQuery = `
CREATE TABLE IF NOT EXISTS migrations(
id text NOT NULL UNIQUE, 
created_at timestamp without time zone default (now() at time zone 'utc') NOT NULL, 
applied boolean NOT NULL
);`

// createMigrationsTable creates migrations table only if it doesn't exist already
func (m Migrate) createMigrationsTable() error {
	_, err := m.db.Exec(migrationsQuery)
	return err
}
