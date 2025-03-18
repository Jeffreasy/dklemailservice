package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Verbindingsgegevens voor de database
	dbHost := "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com"
	dbPort := "5432"
	dbUser := "dekoninklijkeloopdatabase_user"
	dbPassword := "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB"
	dbName := "dekoninklijkeloopdatabase"
	
	// Maak verbinding met de database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	fmt.Println("Verbinden met database...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Kon niet verbinden met database: %v", err)
	}
	defer db.Close()
	
	// Test de verbinding
	fmt.Println("Testen van de verbinding...")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Kon database niet pingen: %v", err)
	}
	fmt.Println("Database verbinding succesvol!")
	
	// Controleer eerst de huidige wachtwoord hashes
	fmt.Println("\nHuidige gebruikers in de database:")
	rows, err := db.Query(`
		SELECT id, naam, email, wachtwoord_hash, rol, is_actief
		FROM gebruikers
		WHERE email IN ('admin@dekoninklijkeloop.nl', 'jeffrey@dekoninklijkeloop.nl')
	`)
	if err != nil {
		log.Fatalf("Kon gebruikers niet ophalen: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, naam, email, wachtwoordHash, rol string
		var isActief bool
		err := rows.Scan(&id, &naam, &email, &wachtwoordHash, &rol, &isActief)
		if err != nil {
			log.Fatalf("Kon rij niet scannen: %v", err)
		}
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Naam: %s\n", naam)
		fmt.Printf("Email: %s\n", email)
		fmt.Printf("WachtwoordHash: %s\n", wachtwoordHash)
		fmt.Printf("Rol: %s\n", rol)
		fmt.Printf("IsActief: %t\n\n", isActief)
	}
	
	// Update de wachtwoorden
	fmt.Println("Updaten van wachtwoorden...")
	
	// Wachtwoord: admin123
	// Hash: $2a$10$3o5nRdG9E8SitO.Zz81M8.z7D5GXNCXl9ZZmQzuR5S5EZ1DUzXfbG
	query := `
	UPDATE gebruikers 
	SET wachtwoord_hash = '$2a$10$3o5nRdG9E8SitO.Zz81M8.z7D5GXNCXl9ZZmQzuR5S5EZ1DUzXfbG', 
		is_actief = TRUE 
	WHERE email IN ('admin@dekoninklijkeloop.nl', 'jeffrey@dekoninklijkeloop.nl')
	`
	
	result, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Kon wachtwoorden niet updaten: %v", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Wachtwoorden geupdate voor %d gebruikers!\n", rowsAffected)
	fmt.Println("Nieuwe wachtwoord is: admin123")
	
	// Controleer of de update succesvol was
	fmt.Println("\nGebruikers na update:")
	rows, err = db.Query(`
		SELECT id, naam, email, wachtwoord_hash, rol, is_actief
		FROM gebruikers
		WHERE email IN ('admin@dekoninklijkeloop.nl', 'jeffrey@dekoninklijkeloop.nl')
	`)
	if err != nil {
		log.Fatalf("Kon gebruikers niet ophalen: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, naam, email, wachtwoordHash, rol string
		var isActief bool
		err := rows.Scan(&id, &naam, &email, &wachtwoordHash, &rol, &isActief)
		if err != nil {
			log.Fatalf("Kon rij niet scannen: %v", err)
		}
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Naam: %s\n", naam)
		fmt.Printf("Email: %s\n", email)
		fmt.Printf("WachtwoordHash: %s\n", wachtwoordHash)
		fmt.Printf("Rol: %s\n", rol)
		fmt.Printf("IsActief: %t\n\n", isActief)
	}
	
	fmt.Println("Programma voltooid. Probeer nu in te loggen met het nieuwe wachtwoord: admin123")
}
