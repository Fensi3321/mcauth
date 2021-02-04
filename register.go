package main

import (
	"database/sql"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

type RegisterState string

// registration states
const (
	OK RegisterState = "OK"
	BadPassword RegisterState = "BadPassword"
	UserExists RegisterState = "UserExists"
	UserNotInLPDB RegisterState = "UserNotInLPDB"
)

func register(w http.ResponseWriter, r *http.Request) {
	db, err := conntectToDB("./bfs")
	if err != nil {
		panic(err)
	}

	nickname := r.FormValue("mcnm")
	serverpasswd := md5Hash(r.FormValue("srvpw"))

	switch handleUser(nickname, serverpasswd, db) {
	case OK :
		fmt.Fprintf(w, "OK")
	case BadPassword: 
		fmt.Fprintf(w, "Bad Password")
	case UserExists:
		fmt.Fprintf(w, "User Exists")
	case UserNotInLPDB:
		fmt.Fprintf(w, "User not in LP DB")
	default:
		fmt.Fprintf(w, "chuj w sumie wie")
	}
}

func handleUser(name string, hash string, db *sql.DB) RegisterState {
	rows, err := db.Query(`select count(*) from player where nick = ?`, name)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var passwd string
	db.QueryRow(`select password from serverpassword`).Scan(&passwd)

	userCount := checkRowCount(rows)

	if userCount == 0 {
		if hash == passwd {
			
			lpDB, lpdbErr := conntectToDB("./luckperms-sqlite.db")
			if lpdbErr != nil {
				panic(lpdbErr)
			}

			if userInLPDBase(name, lpDB) && !userAlreadyRegistered(name, db) {
				var lastid int
				db.QueryRow(`select max(id) from player`).Scan(&lastid)

				_, insertErr := db.Exec(`insert into player values(?, ?)`,lastid+1, name)
				if insertErr != nil {
					panic(insertErr)
				}
				
				return OK;
			} 
			
			
			return UserNotInLPDB;
		}

		return BadPassword;

	} else if userCount > 0 && hash != passwd {
		return BadPassword
	}
		return UserExists;
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

func userAlreadyRegistered(name string, db *sql.DB) bool {
	var count int
	db.QueryRow(`select count(*) from player where nick = ?`, name).Scan(&count)

	return count > 0
}

func md5Hash(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
