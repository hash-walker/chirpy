package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/hash-walker/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
	Platform string
	SecretToken string
}



func main(){

	

	err := godotenv.Load()

	if err != nil {
		return
	}

	dbURL := os.Getenv("DB_URL")
	secretToken := os.Getenv("SECRETTOKEN")

	if dbURL == "" {
		log.Fatal("DB_URL env variable not found")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB: db,
		Platform: os.Getenv("PLATFORM"),
		SecretToken: secretToken,
	}


	const port = "8080"
	const filepathroot = "."

	

	serveMux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathroot))))
	serveMux.Handle("/app/", fsHandler)


	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metric", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	serveMux.HandleFunc("POST /api/users", apiCfg.users)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpyID}", apiCfg.getChirp)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)


	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}