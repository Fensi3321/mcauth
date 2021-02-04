package main

import "database/sql"

func userInLPDBase(name string, db *sql.DB) bool {
	var count int
	_, err := db.QueryRow(`select count(username) from LUCKPERMS_PLAYERS where username = ?`, name).Scan(&count)
	if err != nil {
		panic(err)
	}

	return count > 0
}
