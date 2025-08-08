package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/hash-walker/chirpy/internal/database"
)

type User struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsRed bool `json:"is_chirpy_red"`
}

type AuthPayload struct{
		Token string 
		RefreshToken string
}

func databaseUserToUser(dbUser database.User, payload interface{}) User{

	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
		IsRed: dbUser.IsChirpyRed,
	}

	if auth, ok := payload.(AuthPayload); ok {
		user.Token = auth.Token
		user.RefreshToken = auth.RefreshToken
	}

	return user
}

type Chirp struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID     uuid.UUID    `json:"user_id"`
}

func databaseChirpToChirp(dbChirp database.Chirp) Chirp{
	return Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}
}