package util

import "database/sql"

func EmailExists(db *sql.DB, email string) error {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)",
		email,
	).Scan(&exists)

	if err != nil {
		return err
	}

	return nil
}
