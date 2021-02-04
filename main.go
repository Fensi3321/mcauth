package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	pwd, _ := getPWD()

	fmt.Println(pwd)

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
