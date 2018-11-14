package migration

import (
	"database/sql"
)

type Step struct {
	StepsUp   []string
	StepsDown []string
}

var initialMigration = Step{
	StepsUp: []string{
		"CREATE TABLE _meta_versions (version int PRIMARY KEY NOT NULL);",
	},
	StepsDown: []string{
		"DROP TABLE _meta_versions;",
	},
}

type Set map[int]Step

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
	for _, r := range m.StepsUp {
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
	for _, r := range m.StepsDown {
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
	err := db.QueryRow("SELECT 1 FROM _meta_versions WHERE 1=2;").Scan(&e)
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
