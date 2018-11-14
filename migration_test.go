package migration

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestSet_Upgrade(t *testing.T) {

	os.Remove("./test.db")

	s := Set{
		1: Step{
			Up: []string{
				"CREATE TABLE test (id int PRIMARY KEY);",
			},
			Down: []string{
				"DROP TABLE test;",
			},
		},
	}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = s.Upgrade(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test (id) VALUES ($1)", 1)
	if err != nil {
		t.Fatal(err)
	}

	var id int
	err = db.QueryRow("SELECT id FROM test;").Scan(&id)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Fatal("Wrong id written to database")
	}

	db.Close()

	os.Remove("./test.db")

}

func TestSet_UpgradeToVersion(t *testing.T) {

	os.Remove("./test.db")

	s := Set{
		1: Step{
			Up: []string{
				"CREATE TABLE test (id int PRIMARY KEY);",
			},
			Down: []string{
				"DROP TABLE test;",
			},
		},
		2: Step{
			Up: []string{
				"CREATE TABLE test2 (id int PRIMARY KEY);",
			},
			Down: []string{
				"DROP TABLE test2;",
			},
		},
	}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}

	from, to, err := s.UpgradeToVersion(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if from != 0 && to != 1 {
		t.Fatalf("Upgraded targets didn't met the design from: %d (needed: %d), to: %d (needed: %d)", from, 0, to, 1)
	}

	_, err = db.Exec("INSERT INTO test (id) VALUES ($1)", 1)
	if err != nil {
		t.Fatal(err)
	}

	var id1 int
	err = db.QueryRow("SELECT id FROM test;").Scan(&id1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test2 (id) VALUES ($1)", 1)
	if err == nil {
		t.Fatal("Table from second migration already exists")
	}

	from, to, err = s.UpgradeToVersion(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	if from != 1 && to != 2 {
		t.Fatalf("Upgraded targets didn't met the design from: %d (needed: %d), to: %d (needed: %d)", from, 1, to, 2)
	}

	err = db.QueryRow("SELECT id FROM test;").Scan(&id1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test2 (id) VALUES ($1)", 1)
	if err != nil {
		t.Fatal(err)
	}

	from, to, err = s.UpgradeToVersion(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	if from != 2 && to != 1 {
		t.Fatalf("Upgraded targets didn't met the design from: %d (needed: %d), to: %d (needed: %d)", from, 2, to, 1)
	}

	err = db.QueryRow("SELECT id FROM test;").Scan(&id1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test2 (id) VALUES ($1)", 1)
	if err == nil {
		t.Fatal("Table from second migration still exists after downgrade")
	}

	db.Close()

	os.Remove("./test.db")

}
