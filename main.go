package main

import(
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"encoding/json"
)

func ServerHealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
    })
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    hits := cfg.fileserverHits.Load()
	htmlTemplate := `<html>
  		<body>
   	 		<h1>Welcome, Chirpy Admin</h1>
   		    <p>Chirpy has been visited %d times!</p>
  		</body>
    </html>`
    hits_string := fmt.Sprintf(htmlTemplate, hits)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(hits_string))
}

func (cfg *apiConfig) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
}

func ValidateMessage(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Something went wrong")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		return
    }
    
	if (len(params.Body) > 140){
		log.Printf("Chirp is too long")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Chirp is too long"})
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}

func main() {
    mux := http.NewServeMux()
    server := http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    apiCfg := apiConfig{}
    
    mux.HandleFunc("GET /api/healthz", ServerHealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetMetricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", ValidateMessage)


    fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))

    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}