package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request){
	chirps, err := apiCfg.DB.GetAllChirps(r.Context())

	if err != nil {
		log.Printf("Error getting chirps: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error getting chirps"})
		return
	}

	writeJSON(w, http.StatusOK, chirps)
}

func (apiCfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	chirpID, err := uuid.Parse(r.PathValue("chirpyID"))

	if err != nil {
		writeJSON(w, 400, fmt.Sprintf("Couldn't parse the chirp id: %v", err))
		return 
	}

	chirp, err := apiCfg.DB.GetChirp(r.Context(), chirpID)

	if err != nil {
		writeJSON(w, 400, fmt.Sprintf("Couldn't get the chirp with id: %v", chirpID))
		return
	}

	writeJSON(w, http.StatusOK, chirp)
}