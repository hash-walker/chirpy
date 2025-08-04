package main

import (
	"fmt"
	"net/http"
)


func (apiCfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
  	<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %d times!</p>
  	</body>
	</html>
	`, apiCfg.fileserverHits.Load())))
}

func (apiCfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		apiCfg.fileserverHits.Add(1)
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}