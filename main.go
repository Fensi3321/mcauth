package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	//_ "github.com/mattn/go-sqlite3"
)

func main() {
	pwd, _ := getPWD()

	/*db, dberr := sql.Open("sqlite3", "./luckperms-sqlite.db")
	if dberr != nil {
		panic(dberr)
	}

	defer db.Close()

	dberr = db.Ping()
	if dberr != nil {
		panic(dberr)
	}

	fmt.Println("Successfully connected")*/

	fs := http.FileServer(http.Dir(path.Join(pwd, "/static")))

	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/usrauth", register)

	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		log.Panic("Failed to start server..\n" + err.Error())
	}

	fmt.Println("Listening on port 9000")

}

func getPWD() (string, error) {
	return os.Getwd()
}
