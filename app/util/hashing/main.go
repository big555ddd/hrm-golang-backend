package hashing

import (
	"crypto/rand"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {

	r := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), r)
	log.Println(err)
	return bytes, err
}

func CheckPasswordHash(hash []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	log.Println(err)
	return err == nil
}

func GenerateNumber(min, max int64) int {
	bg := big.NewInt(max - min)
	nBig, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}
	n := nBig.Int64() + min

	return int(n)
}
