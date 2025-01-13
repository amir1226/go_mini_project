package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amir1226/go_mini_project/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// Import the postgres driver
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Default to .env file in the current directory
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		log.Fatal("DB URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(apiCfg.DB, 10, time.Minute)

	log.Printf("Listening on port %s\n", port)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFolLow))
	v1Router.Delete("/feed_follows/{feed_follow_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
