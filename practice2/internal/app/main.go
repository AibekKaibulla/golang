package app

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "practice2/internal/handler"
    "practice2/internal/middleware"
    "practice2/internal/repository"
    _postgres "practice2/internal/repository/_postgres"
    "practice2/internal/usecase"
    "practice2/pkg/modules"
)


func Run() {
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    db := _postgres.NewPGXDialect(ctx, initPostgreConfig())
    repos := repository.NewRepositories(db)
    uc := usecase.NewUserUsecase(repos.UserRepository)
    h := handler.NewUserHandler(uc)

    r := mux.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Auth)

    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    }).Methods("GET")

    r.HandleFunc("/users", h.GetUsers).Methods("GET")
    r.HandleFunc("/users/{id}", h.GetUserByID).Methods("GET")
    r.HandleFunc("/users", h.CreateUser).Methods("POST")
    r.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func initPostgreConfig() *modules.PostgresConfig {
    return &modules.PostgresConfig{
        Host:        "localhost",
        Port:        "5432",
        Username:    "postgres",
        Password:    "postgres",
        DBName:      "practice3",
        SSLMode:     "disable",
        ExecTimeout: 5 * time.Second,
    }
}