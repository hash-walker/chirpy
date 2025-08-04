package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type parameter struct{
		Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Valid bool `json:"valid"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}


func handlerValidateChirp(w http.ResponseWriter, r *http.Request){

	

	var params parameter

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	
	bodyLength := len(params.Body)

	if bodyLength > 400 {
		
		writeJSON(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	writeJSON(w, http.StatusOK, successResponse{Valid: true})

}