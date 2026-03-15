package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
_ 	"github.com/lib/pq" 
	"github.com/joho/godotenv"
	"fmt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

		if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}
	log.Println("Connected to database")

	repo := NewRepository(db)
	handler := NewHandler(repo)

	mux := http.NewServeMux()
	// GET /users
	//   ?page=1&page_size=10&order_by=name&order_dir=asc
	//   &id=5&name=alice&email=gmail&gender=female&birth_date=1995-06-15
	mux.HandleFunc("/users", handler.GetUsersHandler)
 
	// GET /users/common-friends?user_id_1=1&user_id_2=2
	mux.HandleFunc("/users/common-friends", handler.GetCommonFriendsHandler)

	addr := ":8080"
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

	
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}