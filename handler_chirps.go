package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hash-walker/chirpy/internal/auth"
	"github.com/hash-walker/chirpy/internal/database"
)

type parameter struct{
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
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

	if userID != params.UserID{
		log.Printf("Cannot authenticate the user")
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Cannot authenticate the user"})
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