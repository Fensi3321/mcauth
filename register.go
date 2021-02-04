package main

import (
	"database/sql"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

func register(w http.ResponseWriter, r *http.Request) {
	db, err := conntectToDB("./bfs")
	if err != nil {
		panic(err)
	}

	nickname := r.FormValue("mcnm")
	serverpasswd := md5Hash(r.FormValue("srvpw"))

	handleUser(nickname, serverpasswd, db)
}

func handleUser(name string, hash string, db *sql.DB) bool {
	rows, err := db.Query(`select count(*) from player where nick = ?`, name)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if checkRowCount(rows) == 0 {
		var passwd string
		db.QueryRow(`select password from serverpassword`).Scan(&passwd)

		if hash == passwd {
			var lastid int
			db.QueryRow(`select max(id) from player`).Scan(&lastid)

			fmt.Println(lastid)

			_, insertErr := db.Exec(`insert into player values(?, ?)`,lastid+1, name)
			if insertErr != nil {
				panic(insertErr)
			}

			return true;
		}

		return false;
	} 
		return false;
}

func conntectToDB(connstring string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", connstring)

	return db, err
}

func checkRowCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}

	return count
}

func md5Hash(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
