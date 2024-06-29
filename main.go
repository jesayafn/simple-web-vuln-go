package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {

	router := gin.Default()
	port := 9090
	// log.Println("Starting Go SQL Injection Demo application...")

	//HealthCheck path
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	//Vulnerable path
	router.GET("/vuln-path", func(c *gin.Context) {

		db := dbConn()
		defer db.Close()
		username := c.Query("username")
		query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username string
			var password string
			err := rows.Scan(&id, &username, &password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ID":       id,
				"Username": username,
				"Password": password,
			})
		}
	})

	//Secured path
	router.GET("/secured-path", func(c *gin.Context) {
		db := dbConn()
		defer db.Close()

		username := c.Query("username")

		// Use parameterized, and prepared query to avoid SQL injection
		stmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		rows, err := stmt.Query(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username string
			var password string
			if err := rows.Scan(&id, &username, &password); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ID":       id,
				"Username": username,
				"Password": password,
			})
		}
	})

	// log.Printf("Listening on port %d...", port)
	router.Run(fmt.Sprintf(":%d", port))
}
