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

	// Haal alle unieke emails op uit aanmeldingen zonder gebruikersaccount
	var uniqueEmails []struct {
		Email string
		Naam  string
		Rol   string
	}

	// Gebruik DISTINCT ON om per email de meest recente aanmelding te krijgen
	db.Raw(`
		SELECT DISTINCT ON (LOWER(email))
			email, naam, rol
		FROM aanmeldingen
		WHERE gebruiker_id IS NULL
		AND email IS NOT NULL
		AND email != ''
		ORDER BY LOWER(email), created_at DESC
	`).Scan(&uniqueEmails)

	fmt.Printf("Gevonden %d unieke emails zonder gebruikersaccount\n", len(uniqueEmails))

	createdCount := 0
	linkedCount := 0
	skippedCount := 0

	for _, emailData := range uniqueEmails {
		// Check of er al een gebruiker bestaat met dit email (case-insensitive)
		var existingUser Gebruiker
		err := db.Where("LOWER(email) = LOWER(?)", emailData.Email).First(&existingUser).Error

		if err == gorm.ErrRecordNotFound {
			// Bepaal rol op basis van aanmelding.Rol (Deelnemer/Begeleider)
			rol := "deelnemer" // Default
			if strings.ToLower(emailData.Rol) == "begeleider" {
				rol = "begeleider"
			} else if strings.ToLower(emailData.Rol) == "vrijwilliger" {
				rol = "vrijwilliger"
			}

			// Maak nieuwe gebruiker aan
			newUser := Gebruiker{
				Naam:                 emailData.Naam,
				Email:                emailData.Email,
				WachtwoordHash:       string(hashedPassword),
				Rol:                  rol,
				IsActief:             true,
				NewsletterSubscribed: false,
			}

			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("Fout bij aanmaken gebruiker voor %s: %v\n", emailData.Email, err)
				skippedCount++
				continue
			}

			fmt.Printf("✓ Gebruiker aangemaakt voor: %s (%s) [%s]\n", emailData.Naam, emailData.Email, rol)
			createdCount++

			// Link ALLE aanmeldingen met dit email aan nieuwe gebruiker
			result := db.Model(&Aanmelding{}).
				Where("LOWER(email) = LOWER(?) AND gebruiker_id IS NULL", emailData.Email).
				Update("gebruiker_id", newUser.ID)

			if result.Error != nil {
				log.Printf("Fout bij linken aanmeldingen voor %s: %v\n", emailData.Email, result.Error)
			} else {
				linkedCount += int(result.RowsAffected)
				if result.RowsAffected > 1 {
					fmt.Printf("  → %d aanmeldingen gelinkt voor dit email\n", result.RowsAffected)
				}
			}
		} else if err == nil {
			// Gebruiker bestaat al, link alle aanmeldingen met dit email
			result := db.Model(&Aanmelding{}).
				Where("LOWER(email) = LOWER(?) AND gebruiker_id IS NULL", emailData.Email).
				Update("gebruiker_id", existingUser.ID)

			if result.Error != nil {
				log.Printf("Fout bij linken bestaande gebruiker voor %s: %v\n", emailData.Email, result.Error)
				skippedCount++
			} else {
				linkedCount += int(result.RowsAffected)
				fmt.Printf("→ %d aanmelding(en) gelinkt aan bestaande gebruiker: %s (%s)\n",
					result.RowsAffected, emailData.Naam, emailData.Email)
			}
		} else {
			log.Printf("Database fout voor %s: %v\n", emailData.Email, err)
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
