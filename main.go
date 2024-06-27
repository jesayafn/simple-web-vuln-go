package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbName := os.Getenv("MYSQL_DB")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(db:3306)/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	port := 9090
	// log.Println("Starting Go SQL Injection Demo application...")

	//HealthCheck path
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//Vulnerable path
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		db := dbConn()
		defer db.Close()

		username := r.URL.Query().Get("username")

		query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)

		rows, err := db.Query(query)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			var username string
			var password string
			err = rows.Scan(&id, &username, &password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "ID: %d, Username: %s, Password: %s\n", id, username, password)
		}
	},
	)
	log.Printf("Listening on port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
