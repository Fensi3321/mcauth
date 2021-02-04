package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type RegisterState string

// registration states
const (
	OK RegisterState = "OK"
	BadPassword RegisterState = "BadPassword"
	UserExists RegisterState = "UserExists"
	UserDoesntExists RegisterState = "UserDoesntExists"
	UserNotInLPDB RegisterState = "UserNotInLPDB"
)

func register(w http.ResponseWriter, r *http.Request) {
	db, err := conntectToDB("./bfs")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	nickname := r.FormValue("mcnm")
	serverpasswd := md5Hash(r.FormValue("srvpw"))

	switch handleUser(nickname, serverpasswd, db) {
	case OK :
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Changed Perms")

		return

	case BadPassword: 
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Bad Password")

		return
	case UserExists:
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, "User Exists in registered db")

		return
	case UserNotInLPDB:
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "User not in LP DB")
		return

	case UserDoesntExists:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "User not in REG DB")
		return

	default:
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "chuj w sumie wie")
	}
}

func handleUser(name string, hash string, db *sql.DB) RegisterState {
	name = strings.ToLower(name)

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

			defer lpDB.Close()

			isInLPDB := userInLPDBase(name, lpDB)
			isInRegDB := userAlreadyRegistered(name, db)

			if isInLPDB && !isInRegDB {
				
				var uuid string
				lpDB.QueryRow(`select uuid from luckperms_players where username = ?`, name).Scan(&uuid)

				changePermissions(uuid, "paleciak", lpDB)

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
