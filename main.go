package main

import(
	"net/http"
	"log"
)

func ServerHealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	//Instantiation 
	mux := http.NewServeMux()
    server := http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

	//Check if the server is ready
	mux.HandleFunc("/healthz", ServerHealthHandler)
	
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	//Run the server
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}