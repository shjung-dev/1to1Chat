package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)


func GenerateRandomKey() string {
	bytes:=make([]byte ,32) //Allocate space for 256 bits / 32 bytes

	_, err := rand.Read(bytes) //Fills in the space with random sequence of 0 and 1

	if err != nil {
		log.Fatal("Failed to generate key" , err)
	}

	return base64.URLEncoding.EncodeToString(bytes)  //Turn into Base64 string -> 44 characters -> 44 x 6 = 256 bits
}