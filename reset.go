package main

import (
	"fmt"
	"net/http"
)


func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request){
		apiCfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)	
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.fileserverHits.Load())))
}