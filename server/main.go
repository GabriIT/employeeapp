package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	_ "github.com/lib/pq"
)

type Employee struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	BirthDate string `json:"birth_date"`
	EntryDate string `json:"entry_date"`
}

func main() {
	connStr := "postgres://postgres:postgres@172.17.0.1:5432/employees?sslmode=disable"


	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	log.Println("Connected to the database successfully")



	defer db.Close()

	router := gin.Default()
	router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost"},
    AllowMethods:     []string{"GET"},
    AllowHeaders:     []string{"Origin"},
	}))



	router.GET("/api/employee", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'name'"})
			return
		}
		rows, err := db.Query(`SELECT name, surname, birth_date, entry_date FROM employees_data WHERE name ILIKE '%' || $1 || '%'`, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
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

	log.Println("Gin server at :8080")
	router.Run(":8080")
}
