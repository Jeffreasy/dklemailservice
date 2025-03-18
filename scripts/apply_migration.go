package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

// loadEnv laadt de .env file en zet deze als omgevingsvariabelen
func loadEnv(filename string) error {
	// Lees het bestand
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse de inhoud en zet de variabelen
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Sla commentaren en lege regels over
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Verwijder aanhalingstekens als die er zijn
		value = strings.Trim(value, `"'`)

		// Stel de omgevingsvariabele in
		os.Setenv(key, value)
		if key != "DB_PASSWORD" {
			fmt.Printf("Geladen env var: %s\n", key)
		} else {
			fmt.Printf("Geladen env var: %s=********\n", key)
		}
	}

	return nil
}

func main() {
	// Laad de omgevingsvariabelen uit .env
	if err := loadEnv(".env"); err != nil {
		fmt.Printf("Kon .env bestand niet laden: %v\n", err)
		fmt.Println("Probeer handmatig de omgevingsvariabelen in te stellen.")
		// Ga door, misschien zijn de variabelen al ingesteld
	}

	// Migratie bestandspad
	migrationFile := "database/migrations/004_create_incoming_emails_table.sql"

	// Database verbindingsgegevens uit omgevingsvariabelen
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		fmt.Println("Database configuratie ontbreekt. Zorg dat de volgende omgevingsvariabelen zijn ingesteld:")
		fmt.Println("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE")
		os.Exit(1)
	}

	// Database connectie string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	fmt.Println("Verbinden met database:", dbHost, dbPort, dbUser, "********", dbName, dbSSLMode)

	// Verbind met database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Fout bij verbinden met database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test de verbinding
	err = db.Ping()
	if err != nil {
		fmt.Printf("Database connectie test mislukt: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verbonden met database, migratie wordt uitgevoerd...")

	// Lees de migratie file
	content, err := os.ReadFile(filepath.Join(".", migrationFile))
	if err != nil {
		fmt.Printf("Fout bij lezen migratie bestand: %v\n", err)
		os.Exit(1)
	}

	// Voer de migratie uit
	_, err = db.Exec(string(content))
	if err != nil {
		fmt.Printf("Fout bij uitvoeren migratie: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migratie succesvol uitgevoerd.")
}
