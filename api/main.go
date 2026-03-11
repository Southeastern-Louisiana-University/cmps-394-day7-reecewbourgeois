package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Response struct {
	Message string `json:"message"`
	Source  string `json:"source"`
}

func getDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "hello")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", dsn)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := getDB()
	if err != nil {
		log.Printf("db open error: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var message string
	err = db.QueryRow("SELECT content FROM messages LIMIT 1").Scan(&message)
	if err != nil {
		log.Printf("query error: %v", err)
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{
		Message: message,
		Source:  "postgresql",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/health", healthHandler)

	port := getEnv("PORT", "8080")
	log.Printf("API server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
