package main

import "net/http"

func main() {
	handler := http.NewServeMux()

	server := &http.Server{
		Handler: handler,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
