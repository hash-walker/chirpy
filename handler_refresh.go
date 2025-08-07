package main

import (
	"net/http"
	"time"

	"github.com/hash-walker/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return 
	}

	refreshToken, err := apiCfg.DB.CheckToken(r.Context(), token)

	if err != nil {
		writeJSON(w, http.StatusUnauthorized, err)
		return 
	}


	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		writeJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type responseParam struct {
		Token string `json:"token"`
	}

	response := responseParam{
		Token: refreshToken.Token,
	}


	writeJSON(w, http.StatusOK, response)

}