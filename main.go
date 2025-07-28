package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Serve static HTML from ./static/
	http.Handle("/", http.FileServer(http.Dir("./static")))
	
	http.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("name")
		if query == "" {
			http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`SELECT name, surname, birth_date, entry_date FROM employees_data WHERE name = $1`, query)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []Employee
		for rows.Next() {
			var e Employee
			err := rows.Scan(&e.Name, &e.Surname, &e.BirthDate, &e.EntryDate)
			if err != nil {
				log.Println("Row scan error:", err)
				continue
			}
			results = append(results, e)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
