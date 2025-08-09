package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime"       // ← ADDED
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// ---- ADDED: build-time metadata ----
var (
	Version   = "dev" // overridden via -ldflags
	BuildTime = "unknown"
)
// ------------------------------------

type Employee struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	BirthDate string `json:"birth_date"`
	EntryDate string `json:"entry_date"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// local-only fallback for docker-compose on Linux
		dsn = "postgres://postgres:postgres@172.17.0.1:5432/employees?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://emp.athenalabo.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Origin"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			c.String(http.StatusServiceUnavailable, "db not ready")
			return
		}
		c.String(http.StatusOK, "ready")
	})

	// ---- ADDED: version endpoint ----
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":    "api-emp",
			"version":    Version,
			"build_time": BuildTime,
			"go_version": runtime.Version(),
		})
	})
	// ---------------------------------

	// API
	r.GET("/api/employee", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'name'"})
			return
		}
		rows, err := db.Query(`
			SELECT name, surname, birth_date, entry_date
			FROM employees_data
			WHERE name ILIKE '%' || $1 || '%'`, name)
		if err != nil {
			log.Println("DB query error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var out []Employee
		for rows.Next() {
			var e Employee
			if err := rows.Scan(&e.Name, &e.Surname, &e.BirthDate, &e.EntryDate); err == nil {
				out = append(out, e)
			}
		}
		c.JSON(http.StatusOK, out)
	})

	log.Println("Gin server :8080")
	r.Run(":8080")
}
