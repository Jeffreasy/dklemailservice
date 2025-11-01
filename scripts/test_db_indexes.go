package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment or use Render URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	fmt.Println("✅ Database connection successful!")
	fmt.Println()

	// Check if V1_47 migration exists
	var versie, naam string
	var toegepast string
	err = db.QueryRow(`
		SELECT versie, naam, toegepast::text 
		FROM migraties 
		WHERE versie = '1.47.0'
	`).Scan(&versie, &naam, &toegepast)

	if err == sql.ErrNoRows {
		fmt.Println("❌ V1_47 migration NOT found!")
		fmt.Println("   Migration has not been applied yet.")
		fmt.Println("   Check Render deployment logs.")
		os.Exit(1)
	} else if err != nil {
		log.Fatalf("Query error: %v", err)
	}

	fmt.Println("✅ V1_47 Migration found!")
	fmt.Printf("   Versie: %s\n", versie)
	fmt.Printf("   Naam: %s\n", naam)
	fmt.Printf("   Toegepast: %s\n", toegepast)
	fmt.Println()

	// Check new indexes
	fmt.Println("🔍 Checking new indexes from V1_47...")
	fmt.Println()

	expectedIndexes := []string{
		"idx_gebruikers_role_id",
		"idx_aanmeldingen_gebruiker_id",
		"idx_verzonden_emails_contact_id",
		"idx_verzonden_emails_aanmelding_id",
		"idx_contact_antwoorden_contact_id",
		"idx_contact_formulieren_status_created",
		"idx_aanmeldingen_status_created",
		"idx_contact_formulieren_fts",
		"idx_aanmeldingen_fts",
		"idx_chat_messages_fts",
	}

	found := 0
	missing := 0

	for _, indexName := range expectedIndexes {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM pg_indexes 
				WHERE schemaname = 'public' 
				AND indexname = $1
			)
		`, indexName).Scan(&exists)

		if err != nil {
			log.Printf("Error checking index %s: %v", indexName, err)
			continue
		}

		if exists {
			fmt.Printf("✅ %s\n", indexName)
			found++
		} else {
			fmt.Printf("❌ %s - NOT FOUND\n", indexName)
			missing++
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d/%d indexes found\n", found, len(expectedIndexes))
	fmt.Println()

	// Get index count per table
	fmt.Println("📊 Index count per table:")
	rows, err := db.Query(`
		SELECT 
			tablename, 
			COUNT(*) as index_count 
		FROM pg_indexes 
		WHERE schemaname = 'public' 
		GROUP BY tablename 
		HAVING COUNT(*) > 2
		ORDER BY index_count DESC
		LIMIT 10
	`)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		var indexCount int
		if err := rows.Scan(&tableName, &indexCount); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		fmt.Printf("   %s: %d indexes\n", tableName, indexCount)
	}

	fmt.Println()

	// Table sizes
	fmt.Println("💾 Top 5 largest tables:")
	rows, err = db.Query(`
		SELECT 
			tablename,
			pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY pg_total_relation_size('public.'||tablename) DESC
		LIMIT 5
	`)
	if err != nil {
		log.Fatalf("Failed to query table sizes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, size string
		if err := rows.Scan(&tableName, &size); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		fmt.Printf("   %s: %s\n", tableName, size)
	}

	fmt.Println()

	if missing == 0 {
		fmt.Println("🎉 All V1_47 indexes successfully created!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("1. Run: ANALYZE; (to update query planner statistics)")
		fmt.Println("2. Monitor query performance via Render Dashboard")
		fmt.Println("3. Review slow queries after 24 hours")
	} else {
		fmt.Printf("⚠️  %d indexes missing - check migration logs\n", missing)
	}
}
