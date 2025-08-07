package main

import (
	"net/http"
	"log"
)


func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request){

	if (apiCfg.Platform != "dev"){
		w.WriteHeader(http.StatusForbidden)	
		return
	}

	err := apiCfg.DB.DeleteUsers(r.Context())

	if err != nil {
		log.Printf("Error deleting users: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error deleting users"})
		return
	}
		
}