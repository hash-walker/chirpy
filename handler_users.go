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

	writeJSON(w, http.StatusOK, databaseUserToUser(usr))
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


	type AuthResponse struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
    Token     string    `json:"token"`
	RefreshToken string `json:"refreshToken"`
	}

	response := AuthResponse{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email: usr.Email,
		Token: token,
		RefreshToken: refreshToken,
	}

	writeJSON(w, http.StatusOK, response)

}