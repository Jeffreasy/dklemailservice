package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
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
