package migration

import (
	"database/sql"
)

// Step declares each atomic step of your migrations here
type Step struct {
	// Up declares all the queries that should be executed in order to increase the database version
	Up []string
	// Down declares all the queries that should be executed in order to downgrade the database version
	Down []string
}

var initialMigration = Step{
	Up: []string{
		"CREATE TABLE _meta_versions (version int PRIMARY KEY NOT NULL);",
	},
	Down: []string{
		"DROP TABLE _meta_versions;",
	},
}

// Set is used to declare all your database versions and the steps needed to upgrade or downgrade it
type Set map[int]Step

// Upgrade upgrades your database with the highest available version number
func (m Set) Upgrade(db *sql.DB) (int, int, error) {
	return m.UpgradeToVersion(db, m.maxVersion())
}

func (m Set) maxVersion() int {
	max := 0
	for n := range m {
		if n > max {
			max = n
		}
	}
	return max
}

// UpgradeToVersion upgrades or downgrades your database to a specific version
func (m Set) UpgradeToVersion(db *sql.DB, v int) (int, int, error) {
	err := ensureVersionStorageIsPresent(db)
	if err != nil {
		return 0, 0, err
	}
	cv, err := currentVersion(db)
	if err != nil {
		return 0, 0, err
	}
	if cv > v {
		for i := cv; i > v; i-- {
			tx, err := db.Begin()
			if err != nil {
				return 0, 0, err
			}
			err = m[i].unapply(i, tx)
			if err != nil {
				tx.Rollback()
				return 0, 0, err
			}
			tx.Commit()
		}
		return cv, v, nil
	} else if cv < v {
		for i := cv + 1; i <= v; i++ {
			tx, err := db.Begin()
			if err != nil {
				return 0, 0, err
			}
			err = m[i].apply(i, tx)
			if err != nil {
				tx.Rollback()
				return 0, 0, err
			}
			tx.Commit()
		}
		return cv, v, nil
	} else {
		return cv, cv, nil
	}
}

func (m Step) apply(v int, tx *sql.Tx) error {
	for _, r := range m.Up {
		_, err := tx.Exec(r)
		if err != nil {
			return err
		}
	}
	_, err := tx.Exec("INSERT INTO _meta_versions (version) VALUES ($1)", v)
	if err != nil {
		return err
	}
	return nil
}

func (m Step) unapply(v int, tx *sql.Tx) error {
	for _, r := range m.Down {
		_, err := tx.Exec(r)
		if err != nil {
			return err
		}
	}
	_, err := tx.Exec("DELETE FROM _meta_versions WHERE version = $1", v)
	if err != nil {
		return err
	}
	return nil
}

func currentVersion(db *sql.DB) (int, error) {
	var v int
	err := db.QueryRow("SELECT max(version) FROM _meta_versions;").Scan(&v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func ensureVersionStorageIsPresent(db *sql.DB) error {
	var e string
	err := db.QueryRow("SELECT 1 FROM _meta_versions WHERE 1=1;").Scan(&e)
	if err != nil {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		err = initialMigration.apply(0, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}
