package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hash-walker/chirpy/internal/auth"
	"github.com/hash-walker/chirpy/internal/database"
)

type parameter struct{
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}


func checkForProfanes(chirp string) string{

		wordsExclude := make([]string, 3)
		wordsExclude = append(wordsExclude, "kerfuffle")
		wordsExclude = append(wordsExclude, "sharbert")
		wordsExclude = append(wordsExclude, "fornax")

		words := strings.Split(chirp, " ")

		for i, word := range words{
			if slices.Contains(wordsExclude, strings.ToLower(word)){
				words[i] = "****"
			}
		}

		chirp = strings.Join(words, " ")

		return chirp
}


func (apiCfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request){


	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("Cannot get the authorization token %v", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Cannot get the authorization token"})
		return 
	}

	userID, err := auth.ValidateJWT(token, apiCfg.SecretToken)

	if err != nil {
		log.Printf("Cannot validate the token %v", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Cannot validate the token"})
		return 
	}

	var params parameter

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	
	bodyLength := len(params.Body)

	if bodyLength > 140 {
		writeJSON(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}else{

		cleaned_body := checkForProfanes(params.Body)

		chirp, err := apiCfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Body: cleaned_body,
			UserID: userID,
		})

		if err != nil {
			log.Printf("Error creating chirp: %v", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error creating chirp"})
			return
		}

		writeJSON(w, http.StatusOK, databaseChirpToChirp(chirp))

	}	

}

func (apiCfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request){

	s := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")
	

	if s != ""{
		author_id, err := uuid.Parse(s)

		if err != nil {
			writeJSON(w, 400, fmt.Sprintf("Couldn't parse the chirp id: %v", err))
			return 
		}

		chirps, err := apiCfg.DB.GetChirpByAuthor(r.Context(), author_id)

		if err != nil {
			log.Printf("Error getting chirps: %v", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error getting chirps"})
			return
		}

		if order != "" && order == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		} else if order == "asc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		}

		writeJSON(w, http.StatusOK, chirps)


	}else{

		chirps, err := apiCfg.DB.GetAllChirps(r.Context())

		if err != nil {
			log.Printf("Error getting chirps: %v", err)
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error getting chirps"})
			return
		}

		if order != "" && order == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		} else if order == "asc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		}

		writeJSON(w, http.StatusOK, chirps)
	}
}

func (apiCfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	chirpID, err := uuid.Parse(r.PathValue("chirpyID"))

	if err != nil {
		writeJSON(w, 400, fmt.Sprintf("Couldn't parse the chirp id: %v", err))
		return 
	}


	chirps, err := apiCfg.DB.GetChirp(r.Context(), chirpID)

	if err != nil {
		writeJSON(w, 400, fmt.Sprintf("Couldn't get the chirp with id: %v", chirpID))
		return
	}

	writeJSON(w, http.StatusOK, chirps)
	
}


func (apiCfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request){
	chirpID, err := uuid.Parse(r.PathValue("chirpyID"))

	if err != nil {
		writeJSON(w, 400, fmt.Sprintf("Couldn't parse the chirp id: %v", err))
		return 
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return 
	}

	userID, err := auth.ValidateJWT(token, apiCfg.SecretToken)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return 
	}

	chirp, err := apiCfg.DB.GetChirp(r.Context(), chirpID)

	if err != nil {
		writeJSON(w, 404, fmt.Sprintf("Can't find the chirp with id: %v", chirpID))
		return
	}

	if chirp.UserID != userID{
		writeJSON(w, http.StatusUnauthorized, "Unauthorized user")
		return
	}

	err = apiCfg.DB.DeleteChirp(r.Context(), chirpID)

	if chirp.UserID != userID{
		writeJSON(w, http.StatusUnauthorized, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}