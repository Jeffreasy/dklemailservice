package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Gebruik de verbindingsstring met de volledige hostname
	connStr := "postgresql://dekoninklijkeloopdatabase_user:I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB@dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com/dekoninklijkeloopdatabase?sslmode=require"

	// Verbinding maken met de database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		fmt.Printf("Fout bij verbinden met database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Controleren of de verbinding werkt
	err = db.Ping()
	if err != nil {
		fmt.Printf("Fout bij pingen van database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verbinding met database succesvol!")

	// Query om alle tabellen op te halen
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)
	if err != nil {
		fmt.Printf("Fout bij ophalen tabellen: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	// Resultaten verwerken
	fmt.Println("\nTabellen in database:")
	fmt.Println("====================")

	var tableName string
	var tableCount int
	var tableNames []string

	for rows.Next() {
		tableCount++
		err := rows.Scan(&tableName)
		if err != nil {
			fmt.Printf("Fout bij lezen tabel naam: %v\n", err)
			continue
		}
		tableNames = append(tableNames, tableName)
		fmt.Printf("%d. %s\n", tableCount, tableName)
	}

	if tableCount == 0 {
		fmt.Println("Geen tabellen gevonden in de database.")
	}

	// Optioneel: toon aantal rijen per tabel
	fmt.Println("\nAantal rijen per tabel:")
	fmt.Println("=====================")

	rows, err = db.Query(`
		SELECT 
			table_name,
			(SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name AND table_schema = 'public') AS column_count
		FROM 
			information_schema.tables t
		WHERE 
			table_schema = 'public'
		ORDER BY 
			table_name
	`)
	if err != nil {
		fmt.Printf("Fout bij ophalen tabel details: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var columnCount int
	var tableCounts = make(map[string]int)

	for rows.Next() {
		err := rows.Scan(&tableName, &columnCount)
		if err != nil {
			fmt.Printf("Fout bij lezen tabel details: %v\n", err)
			continue
		}

		// Aantal rijen in de tabel ophalen
		var rowCount int
		countErr := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&rowCount)
		if countErr != nil {
			fmt.Printf("%-30s: %d kolommen, fout bij tellen rijen: %v\n", tableName, columnCount, countErr)
			continue
		}

		tableCounts[tableName] = rowCount
		fmt.Printf("%-30s: %d kolommen, %d rijen\n", tableName, columnCount, rowCount)
	}

	// Toon inhoud van tabellen met data
	fmt.Println("\nInhoud van tabellen met data:")
	fmt.Println("============================")

	for _, tableName := range tableNames {
		rowCount := tableCounts[tableName]
		if rowCount > 0 {
			fmt.Printf("\nTabel: %s (%d rijen)\n", tableName, rowCount)
			fmt.Println(strings.Repeat("-", 80))

			// Haal kolomnamen op
			colRows, err := db.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '%s' ORDER BY ordinal_position", tableName))
			if err != nil {
				fmt.Printf("Fout bij ophalen kolommen voor %s: %v\n", tableName, err)
				continue
			}

			var columns []string
			for colRows.Next() {
				var colName string
				if err := colRows.Scan(&colName); err != nil {
					fmt.Printf("Fout bij lezen kolomnaam: %v\n", err)
					continue
				}
				columns = append(columns, colName)
			}
			colRows.Close()

			// Toon kolomnamen
			fmt.Println("Kolommen: " + strings.Join(columns, ", "))
			fmt.Println(strings.Repeat("-", 80))

			// Haal data op (max 10 rijen)
			dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 10", tableName))
			if err != nil {
				fmt.Printf("Fout bij ophalen data voor %s: %v\n", tableName, err)
				continue
			}

			// Bereid scanvariabelen voor
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			// Toon data
			rowNum := 1
			for dataRows.Next() {
				if err := dataRows.Scan(valuePtrs...); err != nil {
					fmt.Printf("Fout bij scannen rij: %v\n", err)
					continue
				}

				fmt.Printf("Rij %d:\n", rowNum)
				for i, col := range columns {
					// Converteer interface{} naar string voor weergave
					var v interface{} = values[i]
					var strVal string
					b, ok := v.([]byte)
					if ok {
						strVal = string(b)
					} else {
						strVal = fmt.Sprintf("%v", v)
					}
					fmt.Printf("  %-20s: %s\n", col, strVal)
				}
				fmt.Println()
				rowNum++
			}
			dataRows.Close()
		}
	}
}
