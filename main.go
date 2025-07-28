package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Employee struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	BirthDate string `json:"birth_date"`
	EntryDate string `json:"entry_date"`
}

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/employees?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static") // serves /static/index.html
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html") // serves it at /
	})
		

	// API endpoint: /employee?name=Alice
	router.GET("/employee", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'name' parameter"})
			return
		}

		rows, err := db.Query(`SELECT name, surname, birth_date, entry_date FROM employees_data WHERE name = $1`, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var results []Employee
		for rows.Next() {
			var e Employee
			if err := rows.Scan(&e.Name, &e.Surname, &e.BirthDate, &e.EntryDate); err == nil {
				results = append(results, e)
			}
		}

		c.JSON(http.StatusOK, results)
	})

	log.Println("Gin server running at http://localhost:8080")
	router.Run(":8080")
}
