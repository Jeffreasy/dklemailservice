package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gebruiker struct {
	ID                   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Naam                 string
	Email                string `gorm:"uniqueIndex"`
	WachtwoordHash       string
	Rol                  string
	IsActief             bool
	NewsletterSubscribed bool
}

type Aanmelding struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	Naam        string
	Email       string
	Rol         string  // Deelnemer, Begeleider, of Vrijwilliger
	GebruikerID *string `gorm:"type:uuid"`
}

func main() {
	// Laad .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Geen .env file gevonden, gebruik omgevingsvariabelen")
	}

	// Connecteer met database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Kon niet verbinden met database:", err)
	}

	fmt.Println("Verbonden met database")

	// Standaard wachtwoord voor nieuwe accounts (gebruikers moeten dit wijzigen)
	defaultPassword := "DKL2025!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Kon wachtwoord niet hashen:", err)
	}

	// Tel hoeveel aanmeldingen er zijn zonder gebruikersaccount
	var aanmeldingenZonderAccount []Aanmelding
	result := db.Where("gebruiker_id IS NULL AND email IS NOT NULL AND email != ''").
		Find(&aanmeldingenZonderAccount)

	if result.Error != nil {
		log.Fatal("Kon aanmeldingen niet ophalen:", result.Error)
	}

	fmt.Printf("Gevonden %d aanmeldingen zonder gebruikersaccount\n", len(aanmeldingenZonderAccount))

	createdCount := 0
	linkedCount := 0
	skippedCount := 0

	for _, aanmelding := range aanmeldingenZonderAccount {
		// Check of er al een gebruiker bestaat met dit email
		var existingUser Gebruiker
		err := db.Where("email = ?", aanmelding.Email).First(&existingUser).Error

		if err == gorm.ErrRecordNotFound {
			// Bepaal rol op basis van aanmelding.Rol (Deelnemer/Begeleider)
			rol := "deelnemer" // Default
			if strings.ToLower(aanmelding.Rol) == "begeleider" {
				rol = "begeleider"
			} else if strings.ToLower(aanmelding.Rol) == "vrijwilliger" {
				rol = "vrijwilliger"
			}

			// Maak nieuwe gebruiker aan
			newUser := Gebruiker{
				Naam:                 aanmelding.Naam,
				Email:                aanmelding.Email,
				WachtwoordHash:       string(hashedPassword),
				Rol:                  rol,
				IsActief:             true,
				NewsletterSubscribed: false,
			}

			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("Fout bij aanmaken gebruiker voor %s: %v\n", aanmelding.Email, err)
				skippedCount++
				continue
			}

			fmt.Printf("✓ Gebruiker aangemaakt voor: %s (%s)\n", aanmelding.Naam, aanmelding.Email)
			createdCount++

			// Link aanmelding aan nieuwe gebruiker
			if err := db.Model(&Aanmelding{}).Where("id = ?", aanmelding.ID).
				Update("gebruiker_id", newUser.ID).Error; err != nil {
				log.Printf("Fout bij linken aanmelding voor %s: %v\n", aanmelding.Email, err)
			} else {
				linkedCount++
			}
		} else if err == nil {
			// Gebruiker bestaat al, link alleen de aanmelding
			if err := db.Model(&Aanmelding{}).Where("id = ?", aanmelding.ID).
				Update("gebruiker_id", existingUser.ID).Error; err != nil {
				log.Printf("Fout bij linken bestaande gebruiker voor %s: %v\n", aanmelding.Email, err)
				skippedCount++
			} else {
				fmt.Printf("→ Aanmelding gelinkt aan bestaande gebruiker: %s (%s)\n", aanmelding.Naam, aanmelding.Email)
				linkedCount++
			}
		} else {
			log.Printf("Database fout voor %s: %v\n", aanmelding.Email, err)
			skippedCount++
		}
	}

	fmt.Println("\n=== Samenvatting ===")
	fmt.Printf("Nieuwe gebruikers aangemaakt: %d\n", createdCount)
	fmt.Printf("Aanmeldingen gelinkt: %d\n", linkedCount)
	fmt.Printf("Overgeslagen (fouten): %d\n", skippedCount)
	fmt.Printf("\nStandaard wachtwoord voor nieuwe accounts: %s\n", defaultPassword)
	fmt.Println("Gebruikers moeten hun wachtwoord wijzigen via de app!")
}
