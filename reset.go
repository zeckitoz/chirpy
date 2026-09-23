package main

import "net/http"

func (cfg *apiConfig) resetFileServerHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FileServerHits have been set back to 0"))
}
