package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/MrBhop/Chirpy/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable is not set")
	}
	
	dbConnection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error, opening db: %v", err)
	}
	
	dbQueries := database.New(dbConnection)

	apiCfg := handlers.ApiConfig{
		FileServerHits: atomic.Int32{},
		Db: dbQueries,
		Platform: platform,
		Secret: secret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	mux.HandleFunc("POST /api/users", apiCfg.HandlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUsersChange)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerChirpsCreate)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{" + handlers.ChirpIDParameter +"}", apiCfg.HandlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{" + handlers.ChirpIDParameter +"}", apiCfg.HandlerChirpsDelete)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
