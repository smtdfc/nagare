package helpers

import (
	"log"

	"github.com/smtdfc/nagare/shared/helpers"
)

func GetKeyPair() *helpers.RSAKeys {
	service := "my-app"
	user := "anon"
	password := "secret"

	// set password
	err := keyring.Set(service, user, password)
	if err != nil {
		log.Fatal(err)
	}

	// get password
	secret, err := keyring.Get(service, user)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(secret)
}
