package db

import (
	"io/ioutil"
	"testing"
)

func setup(t *testing.T) string {
	schemas["test"] = &Schema{
		name: "test",
		latest: `
PRAGMA user_version = 2;
CREATE TABLE Foo(
	id INT NOT NULL,
	name TEXT
)
`,
		migrations: []string{
			`
PRAGMA user_version = 1;
CREATE TABLE Foo(
	id INT NOT NULL
)
`,
			`
PRAGMA user_version = 2;
ALTER TABLE Foo ADD COLUMN name TEXT;
`,
		},
	}
	ftmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("ioutil.TempFile failed: %v", err)
	}
	path := ftmp.Name()
	ftmp.Close()
	return path
}

func TestNew(t *testing.T) {
	path := setup(t)

	db := NewDatabase(path, "test")
	if db.version != 2 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()
}

func TestMigration(t *testing.T) {
	path := setup(t)

	db := NewDatabase(path, "test")
	if db.version != 2 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()

	latest := `
PRAGMA user_version = 3;
CREATE TABLE Foo(
	id INT NOT NULL,
	name TEXT,
	author TEXT
)
	`

	migration := `
PRAGMA user_version = 3;
ALTER TABLE Foo ADD COLUMN author TEXT;
	`

	schema := schemas["test"]
	schema.latest = latest
	schema.migrations = append(schema.migrations, migration)
	db = NewDatabase(path, "test")
	if db.version != 3 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()
}

func TestBroken(t *testing.T) {
	path := setup(t)

	schemas["broken"] = &Schema{
		name: "broken",
		latest: `
CREATE TABLE Foo(
	id INT NOT NULL
)
`}

	// start with broken install (like from initial install party)
	db := NewDatabase(path, "broken")
	if db.version != 0 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()

	db = NewDatabase(path, "test")
	if db.version != 2 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()

	latest := `
PRAGMA user_version = 3;
CREATE TABLE Foo(
	id INT NOT NULL,
	name TEXT,
	author TEXT
)
	`

	migration := `
PRAGMA user_version = 3;
ALTER TABLE Foo ADD COLUMN author TEXT;
	`

	schema := schemas["test"]
	schema.latest = latest
	schema.migrations = append(schema.migrations, migration)
	db = NewDatabase(path, "test")
	if db.version != 3 {
		t.Fatalf("Invalid version: %d", db.version)
	}
	db.Close()
}
