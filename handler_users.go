package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hash-walker/chirpy/internal/auth"
	"github.com/hash-walker/chirpy/internal/database"
)

func (apiCfg *apiConfig) users(w http.ResponseWriter, r *http.Request){

	type parameter struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var params parameter

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error generating hashed password"})
		return
	}

	usr, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		HashedPassword: hashed_password,
		Email: params.Email,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error creating user"})
		return
	}

	writeJSON(w, http.StatusOK, databaseUserToUser(usr, nil))
}

func (apiCfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request){

	type parameter struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var params parameter
	expireTime := time.Duration(3600) * time.Second

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	usr, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Incorrect email or password"})
		return
	}

	err = auth.CheckPasswordHash(usr.HashedPassword, params.Password)

	if err != nil {
		log.Printf("Error comparing hash: %v", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Incorrect email or password"})
		return
	}


	token, err := auth.MakeJWT(usr.ID, apiCfg.SecretToken, expireTime)


	if err != nil {
		log.Printf("Cannot generate token: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Cannot generate token"})
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	expirationTime := time.Now().Add(60 * 24 * time.Hour)

	apiCfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: usr.ID,
		ExpiresAt: expirationTime,
		RevokedAt: sql.NullTime{Valid: false},
	})

	payload := AuthPayload{
		Token: token,
		RefreshToken: refreshToken,
	}
	

	writeJSON(w, http.StatusOK, databaseUserToUser(usr, payload))

}

func (apiCfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request){

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
	
	type parameter struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var params parameter

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error generating hashed password"})
		return
	}

	err = apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
		ID: userID,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Error updating user details"})
		return
	}

	usr, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Incorrect email or password"})
		return
	}

	writeJSON(w, http.StatusOK, databaseUserToUser(usr, nil))

}

func (apiCfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request){

	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
	}

	if apiKey != apiCfg.PolkaKey{
		w.WriteHeader(http.StatusUnauthorized)
		return 
	}


	type parameter struct{
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var params parameter

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Something went wrong"})
		return
	}

	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return 
	}

	err = apiCfg.DB.UpgradeUsers(r.Context(), params.Data.UserID)

	if err != nil{
		writeJSON(w, http.StatusNotFound, err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)

}