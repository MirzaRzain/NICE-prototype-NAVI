package utils

import (
	"os"
	"strconv"

	"be-dashboard-nba/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Generate hash password with bcrypt.
func GenerateHashPassword(password string) (hashed string, err error) {
	costStr := os.Getenv("AUTH_BCRYPT_COST")

	cost, err := strconv.Atoi(costStr)
	if err != nil {
		err = errors.Wrapf(err, "error parse int on bcrypt cost env : %s", costStr)
		return
	}

	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}

	crypt, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		err = errors.Wrap(err, "error generate bcrypt hash password")
		return
	}

	hashed = string(crypt)

	return
}

// Compare bcrypt hashed password.
func CompareHashPassword(passwordInput, passwordDB string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(passwordInput))
	return
}

// Generate UUID.
func GenerateUUID() string {
	guid, err := uuid.NewV7()
	if err != nil {
		logger.Errorf("error generate UUID %v", err)
		return uuid.NewString()
	}

	return guid.String()
}

