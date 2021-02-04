package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

func userInLPDBase(name string, db *sql.DB) bool {
	var count int
	db.QueryRow(`select count(username) from LUCKPERMS_PLAYERS where username = ?`, name).Scan(&count)

	fmt.Println("LPDB COUNT " + strconv.FormatInt(int64(count), 10) + " " + name)

	return count > 0
}
