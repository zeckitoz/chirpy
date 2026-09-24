package auth

import (
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
		return "", err
	}
	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		log.Panicf("Couldn't check password and hash: %v", err)
		return "", err
	}
	if !ok {
		log.Panic("Password and hash doesn't match.")
		return "", errors.New("Password and hash doesn't match.")
	}

	return hash, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatalf("Error checking hash: %v", err)
		return false, err
	}
	return match, nil
}
