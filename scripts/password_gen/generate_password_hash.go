package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin"

	// Genereer bcrypt hash (cost 10, zoals in de applicatie)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Fout bij genereren hash:", err)
	}

	fmt.Println("===========================================")
	fmt.Println("Wachtwoord Hash Generator")
	fmt.Println("===========================================")
	fmt.Printf("Wachtwoord: %s\n", password)
	fmt.Printf("Bcrypt Hash: %s\n", string(hash))
	fmt.Println("===========================================")
	fmt.Println("\nKopieer deze hash naar het SQL script:")
	fmt.Printf("'%s'\n", string(hash))
	fmt.Println("===========================================")

	// Verifieer dat de hash werkt
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("✓ Hash verificatie succesvol!")
	} else {
		fmt.Println("✗ Hash verificatie gefaald!")
	}
}
