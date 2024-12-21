package main

import(
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
    server := http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
	
	mux.Handle("/", http.FileServer(http.Dir(".")))

    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}