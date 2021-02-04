package main

import (
	"database/sql"
)

func userInLPDBase(name string, db *sql.DB) bool {
	var count int
	db.QueryRow(`select count(username) from LUCKPERMS_PLAYERS where username = ?`, name).Scan(&count)

	return count > 0
}

func changePermissions(uuid string, permName string, db *sql.DB) {
	db.Exec(`update luckperms_players set primary_group = ? where uuid = ?`, permName, uuid)
}
