package tests

import (
	"database/sql"
	"dklautomationgo/models"
	"fmt"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModelDatabaseAlignment verifies that Go models match database schema
func TestModelDatabaseAlignment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping model-database alignment tests in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	t.Run("ParticipantModel", func(t *testing.T) {
		testParticipantModelAlignment(t, db)
	})

	t.Run("EventRegistrationModel", func(t *testing.T) {
		testEventRegistrationModelAlignment(t, db)
	})

	t.Run("ParticipantRoleModel", func(t *testing.T) {
		testParticipantRoleModelAlignment(t, db)
	})

	t.Run("DistanceModel", func(t *testing.T) {
		testDistanceModelAlignment(t, db)
	})

	t.Run("StatusTypeModels", func(t *testing.T) {
		testStatusTypeModelsAlignment(t, db)
	})

	t.Run("EventModel", func(t *testing.T) {
		testEventModelAlignment(t, db)
	})
}

// testParticipantModelAlignment verifies Participant model matches database
func testParticipantModelAlignment(t *testing.T, db *sql.DB) {
	// Get table columns from database
	columns := getTableColumns(t, db, "participants")

	expectedColumns := map[string]bool{
		"id":           true,
		"naam":         true,
		"email":        true,
		"telefoon":     true,
		"terms":        true,
		"gebruiker_id": true,
		"test_mode":    true,
		"created_at":   true,
		"updated_at":   true,
	}

	for col := range expectedColumns {
		assert.Contains(t, columns, col, "participants table should have column %s", col)
	}

	// Verify columns that should NOT be on participants (moved to event_registrations)
	movedColumns := []string{"afstand", "rol", "ondersteuning", "bijzonderheden", "steps", "status"}
	for _, col := range movedColumns {
		assert.NotContains(t, columns, col,
			"participants table should NOT have column %s (moved to event_registrations)", col)
	}

	// Verify Participant struct has matching fields
	var p models.Participant
	structFields := getStructFields(p)

	// Check key fields exist in struct
	requiredFields := []string{"ID", "Naam", "Email", "Telefoon", "Terms", "GebruikerID", "TestMode"}
	for _, field := range requiredFields {
		assert.Contains(t, structFields, field, "Participant struct should have field %s", field)
	}
}

// testEventRegistrationModelAlignment verifies EventRegistration model
func testEventRegistrationModelAlignment(t *testing.T, db *sql.DB) {
	columns := getTableColumns(t, db, "event_registrations")

	// Key columns that should exist
	expectedColumns := []string{
		"id",
		"event_id",
		"participant_id",
		"registered_at",
		"steps",
		"participant_role_name",
		"distance_route",
		"status",
		"ondersteuning",
		"bijzonderheden",
		"test_mode",
	}

	for _, col := range expectedColumns {
		assert.Contains(t, columns, col, "event_registrations should have column %s", col)
	}

	// Verify EventRegistration struct
	var er models.EventRegistration
	structFields := getStructFields(er)

	requiredFields := []string{
		"ID", "EventID", "ParticipantID", "Steps",
		"ParticipantRoleName", "DistanceRoute", "Status",
	}
	for _, field := range requiredFields {
		assert.Contains(t, structFields, field,
			"EventRegistration struct should have field %s", field)
	}

	// Verify FK columns match model fields
	t.Run("ForeignKeyFields", func(t *testing.T) {
		// Check Status field points to correct column
		assert.Contains(t, structFields, "Status", "Should have Status field")
		assert.Contains(t, columns, "status", "Database should have status column")

		// Check DistanceRoute field
		assert.Contains(t, structFields, "DistanceRoute", "Should have DistanceRoute field")
		assert.Contains(t, columns, "distance_route", "Database should have distance_route column")

		// Check ParticipantRoleName field
		assert.Contains(t, structFields, "ParticipantRoleName", "Should have ParticipantRoleName field")
		assert.Contains(t, columns, "participant_role_name", "Database should have participant_role_name column")
	})
}

// testParticipantRoleModelAlignment verifies ParticipantRole model
func testParticipantRoleModelAlignment(t *testing.T, db *sql.DB) {
	columns := getTableColumns(t, db, "participant_roles")

	expectedColumns := []string{"name", "description", "is_active", "created_at", "updated_at"}
	for _, col := range expectedColumns {
		assert.Contains(t, columns, col, "participant_roles should have column %s", col)
	}

	// Verify ParticipantRole struct has IsActive field (V30 fix)
	var pr models.ParticipantRole
	structFields := getStructFields(pr)
	assert.Contains(t, structFields, "IsActive",
		"ParticipantRole struct should have IsActive field (V30 fix)")

	// Verify primary key is on 'name'
	var pkColumn string
	err := db.QueryRow(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = 'participant_roles'::regclass AND i.indisprimary
	`).Scan(&pkColumn)
	require.NoError(t, err)
	assert.Equal(t, "name", pkColumn, "participant_roles PK should be on 'name'")
}

// testDistanceModelAlignment verifies Distance model
func testDistanceModelAlignment(t *testing.T, db *sql.DB) {
	columns := getTableColumns(t, db, "distances")

	// Critical: fund_amount column must exist (V26 consolidation)
	expectedColumns := []string{"route", "distance_km", "fund_amount", "description"}
	for _, col := range expectedColumns {
		assert.Contains(t, columns, col, "distances should have column %s", col)
	}

	// Verify Distance struct has FundAmount field
	var d models.Distance
	structFields := getStructFields(d)
	assert.Contains(t, structFields, "FundAmount",
		"Distance struct should have FundAmount field (V26 fix)")

	// Verify primary key is on 'route'
	var pkColumn string
	err := db.QueryRow(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = 'distances'::regclass AND i.indisprimary
	`).Scan(&pkColumn)
	require.NoError(t, err)
	assert.Equal(t, "route", pkColumn, "distances PK should be on 'route'")

	// Verify data type of fund_amount
	var dataType string
	err = db.QueryRow(`
		SELECT data_type 
		FROM information_schema.columns 
		WHERE table_name = 'distances' AND column_name = 'fund_amount'
	`).Scan(&dataType)
	require.NoError(t, err)
	assert.Equal(t, "numeric", dataType, "fund_amount should be numeric type")
}

// testStatusTypeModelsAlignment verifies all status type lookup models
func testStatusTypeModelsAlignment(t *testing.T, db *sql.DB) {
	statusTypes := []struct {
		table    string
		model    interface{}
		pkColumn string
	}{
		{"contact_status_types", models.ContactStatusType{}, "status"},
		{"registration_status_types", models.RegistrationStatusType{}, "status"},
		{"email_status_types", models.EmailStatusType{}, "status"},
		{"event_status_types", models.EventStatusType{}, "status"},
		{"notification_types", models.NotificationType{}, "type"},
		{"notification_priority_types", models.NotificationPriorityType{}, "priority"},
	}

	for _, st := range statusTypes {
		t.Run(st.table, func(t *testing.T) {
			columns := getTableColumns(t, db, st.table)

			// All lookup tables should have description
			assert.Contains(t, columns, "description",
				"%s should have description column", st.table)

			// Verify PK column
			assert.Contains(t, columns, st.pkColumn,
				"%s should have %s column", st.table, st.pkColumn)

			// Verify struct has matching fields
			structFields := getStructFields(st.model)
			assert.Contains(t, structFields, "Description",
				"%s model should have Description field", st.table)
		})
	}
}

// testEventModelAlignment verifies Event model status reference
func testEventModelAlignment(t *testing.T, db *sql.DB) {
	columns := getTableColumns(t, db, "events")

	// Events should have status column
	assert.Contains(t, columns, "status", "events should have status column")

	// Verify FK to event_status_types
	var fkExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
				ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.constraint_column_usage ccu
				ON ccu.constraint_name = tc.constraint_name
			WHERE tc.constraint_type = 'FOREIGN KEY'
				AND tc.table_name = 'events'
				AND kcu.column_name = 'status'
				AND ccu.table_name = 'event_status_types'
		)
	`).Scan(&fkExists)
	require.NoError(t, err)
	assert.True(t, fkExists, "events should have FK to event_status_types")

	// Verify Event struct
	var e models.Event
	structFields := getStructFields(e)
	assert.Contains(t, structFields, "Status", "Event struct should have Status field")
}

// Helper function to get table columns from database
func getTableColumns(t *testing.T, db *sql.DB, tableName string) map[string]bool {
	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = $1
	`
	rows, err := db.Query(query, tableName)
	require.NoError(t, err)
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var colName string
		err := rows.Scan(&colName)
		require.NoError(t, err)
		columns[colName] = true
	}

	return columns
}

// Helper function to get struct field names using reflection
func getStructFields(model interface{}) map[string]bool {
	fields := make(map[string]bool)

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields[field.Name] = true
	}

	return fields
}

// TestCriticalFieldMappings tests specific critical field mappings
func TestCriticalFieldMappings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping critical field mappings test in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	t.Run("EventRegistration_Status_Mapping", func(t *testing.T) {
		// V27 fix: column renamed from status_key to status
		var hasOldColumn bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'event_registrations' AND column_name = 'status_key'
			)
		`).Scan(&hasOldColumn)
		require.NoError(t, err)
		assert.False(t, hasOldColumn, "Old column 'status_key' should not exist")

		var hasNewColumn bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'event_registrations' AND column_name = 'status'
			)
		`).Scan(&hasNewColumn)
		require.NoError(t, err)
		assert.True(t, hasNewColumn, "New column 'status' should exist")

		// Verify model field name
		var er models.EventRegistration
		structFields := getStructFields(er)
		assert.Contains(t, structFields, "Status", "EventRegistration should have Status field")
	})

	t.Run("Distance_FundAmount_Mapping", func(t *testing.T) {
		// V26 fix: column renamed from amount to fund_amount
		var hasOldColumn bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'distances' AND column_name = 'amount'
			)
		`).Scan(&hasOldColumn)
		require.NoError(t, err)
		assert.False(t, hasOldColumn, "Old column 'amount' should not exist")

		var hasNewColumn bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'distances' AND column_name = 'fund_amount'
			)
		`).Scan(&hasNewColumn)
		require.NoError(t, err)
		assert.True(t, hasNewColumn, "New column 'fund_amount' should exist")

		// Verify model field
		var d models.Distance
		structFields := getStructFields(d)
		assert.Contains(t, structFields, "FundAmount", "Distance should have FundAmount field")
	})

	t.Run("ParticipantRole_IsActive_Mapping", func(t *testing.T) {
		// V30 fix: is_active column added
		var hasColumn bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'participant_roles' AND column_name = 'is_active'
			)
		`).Scan(&hasColumn)
		require.NoError(t, err)
		assert.True(t, hasColumn, "participant_roles should have is_active column")

		// Verify model field
		var pr models.ParticipantRole
		structFields := getStructFields(pr)
		assert.Contains(t, structFields, "IsActive", "ParticipantRole should have IsActive field")
	})

	t.Run("TableName_Methods", func(t *testing.T) {
		// Verify models have correct TableName() methods
		tests := []struct {
			model    interface{ TableName() string }
			expected string
		}{
			{models.Participant{}, "participants"},
			{models.EventRegistration{}, "event_registrations"},
			{models.ParticipantAntwoord{}, "participant_antwoorden"},
		}

		for _, tt := range tests {
			actual := tt.model.TableName()
			assert.Equal(t, tt.expected, actual,
				"TableName() should return %s", tt.expected)
		}
	})
}

// TestForeignKeyIntegrity verifies FK relationships work correctly
func TestForeignKeyIntegrity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping FK integrity test in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	t.Run("EventRegistration_to_ParticipantRole", func(t *testing.T) {
		// Try to insert invalid participant_role_name (should fail)
		_, err := db.Exec(`
			INSERT INTO event_registrations (
				id, event_id, participant_id, participant_role_name, status
			) VALUES (
				gen_random_uuid(), 
				(SELECT id FROM events LIMIT 1),
				(SELECT id FROM participants LIMIT 1),
				'invalid_role_that_does_not_exist',
				'registered'
			)
		`)
		assert.Error(t, err, "Should fail to insert invalid participant_role_name")
	})

	t.Run("EventRegistration_to_Distance", func(t *testing.T) {
		// Try to insert invalid distance_route (should fail)
		_, err := db.Exec(`
			INSERT INTO event_registrations (
				id, event_id, participant_id, distance_route, status
			) VALUES (
				gen_random_uuid(),
				(SELECT id FROM events LIMIT 1),
				(SELECT id FROM participants LIMIT 1),
				'invalid_distance_9999km',
				'registered'
			)
		`)
		assert.Error(t, err, "Should fail to insert invalid distance_route")
	})

	t.Run("EventRegistration_to_StatusType", func(t *testing.T) {
		// Try to insert invalid status (should fail)
		_, err := db.Exec(`
			INSERT INTO event_registrations (
				id, event_id, participant_id, status
			) VALUES (
				gen_random_uuid(),
				(SELECT id FROM events LIMIT 1),
				(SELECT id FROM participants LIMIT 1),
				'invalid_status_xyz'
			)
		`)
		assert.Error(t, err, "Should fail to insert invalid status")
	})
}

// TestDataTypeConsistency verifies data types match between model and database
func TestDataTypeConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping data type consistency test in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	tests := []struct {
		table    string
		column   string
		dataType string
	}{
		{"participants", "id", "uuid"},
		{"participants", "terms", "boolean"},
		{"participants", "test_mode", "boolean"},
		{"event_registrations", "steps", "integer"},
		{"event_registrations", "test_mode", "boolean"},
		{"distances", "fund_amount", "numeric"},
		{"distances", "distance_km", "numeric"},
		{"participant_roles", "is_active", "boolean"},
		{"gebruikers", "is_actief", "boolean"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.table, tt.column), func(t *testing.T) {
			var dataType string
			err := db.QueryRow(`
				SELECT data_type 
				FROM information_schema.columns 
				WHERE table_name = $1 AND column_name = $2
			`, tt.table, tt.column).Scan(&dataType)
			require.NoError(t, err)
			assert.Contains(t, dataType, tt.dataType,
				"Column %s.%s should be of type %s", tt.table, tt.column, tt.dataType)
		})
	}
}
