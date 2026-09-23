package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) fileServerHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) resetFileServerHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FileServerHits have been set back to 0"))
}

func healtzHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	apiCfg := apiConfig{atomic.Int32{}}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", healtzHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.fileServerHitsHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.resetFileServerHits)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
