package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}



func main(){

	const port = "8080"
	const filepathroot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathroot))))
	serveMux.Handle("/app/", fsHandler)


	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metric", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)


	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}