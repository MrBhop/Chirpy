package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/MrBhop/Chirpy/internal/database"
	"github.com/google/uuid"
)

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("Couldn't generate random bytes: %w", err)
	}

	return hex.EncodeToString(randomBytes), nil
}

func MakeRegisteredRefreshToken(databaseConnection *database.Queries, r *http.Request, userId uuid.UUID) (string, error) {
	token, err := MakeRefreshToken()
	if err != nil {
		return "", err
	}

	tokenRecord, err := databaseConnection.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: token,
		UserID: userId,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		return "", err
	}

	fmt.Printf("\nReturning token: %s\n\n", tokenRecord.Token)
	return tokenRecord.Token, nil
}
