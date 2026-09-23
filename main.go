package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zeckitoz/chirpy/internal/database"
)

type apiConfig struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error logging into database: %s", err)
		return
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{atomic.Int32{}, dbQueries, platform}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", healtzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.fileServerHitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetServer)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
