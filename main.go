package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// connect reads the password and returns a DB connection
func connect() (*sql.DB, error) {
	bin, err := ioutil.ReadFile("/app/secrets/db-password")
	if err != nil {
		return nil, err
	}
	password := strings.TrimSpace(string(bin)) // remove newline/whitespace
	return sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/example?sslmode=disable", password))
}

// blogHandler reuses the global DB connection
func blogHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT title FROM blog")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		titles = append(titles, title)
	}
	json.NewEncoder(w).Encode(titles)
}

func main() {
	var err error
	log.Print("Connecting to DB...")
	db, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping DB until ready
	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Printf("Waiting for DB... attempt %d", i+1)
		time.Sleep(time.Second)
	}

	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening on 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", blogHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func prepare() error {
	if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id SERIAL, title VARCHAR)"); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if _, err := db.Exec("INSERT INTO blog (title) VALUES ($1);", fmt.Sprintf("Blog post #%d", i)); err != nil {
			return err
		}
	}
	return nil
}
// test
