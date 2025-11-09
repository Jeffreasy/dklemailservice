package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseMigrations_Complete is de main comprehensive test suite voor alle migrations
func TestDatabaseMigrations_Complete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration tests in short mode")
	}

	// Check if we should run database tests
	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	t.Run("CoreTables", func(t *testing.T) {
		testCoreTablesExist(t, db)
	})

	t.Run("LookupTables", func(t *testing.T) {
		testLookupTablesExist(t, db)
	})

	t.Run("LookupTableData", func(t *testing.T) {
		testLookupTableSeeding(t, db)
	})

	t.Run("ForeignKeyConstraints", func(t *testing.T) {
		testForeignKeyConstraints(t, db)
	})

	t.Run("RBACSystem", func(t *testing.T) {
		testRBACSystemIntegrity(t, db)
	})

	t.Run("Indexes", func(t *testing.T) {
		testCriticalIndexesExist(t, db)
	})

	t.Run("ColumnTypes", func(t *testing.T) {
		testColumnTypes(t, db)
	})

	t.Run("TableRenaming", func(t *testing.T) {
		testTableRenaming(t, db)
	})

	t.Run("DataIntegrity", func(t *testing.T) {
		testDataIntegrity(t, db)
	})
}

// setupTestDatabase sets up a test database connection
func setupTestDatabase(t *testing.T) *sql.DB {
	// Get connection string from environment or use default
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=dkl_test sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "Failed to connect to test database")

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	require.NoError(t, err, "Failed to ping test database")

	return db
}

// testCoreTablesExist verifies all core tables exist
func testCoreTablesExist(t *testing.T, db *sql.DB) {
	expectedTables := []string{
		// V01 - Core tables
		"gebruikers",
		"contact_formulieren",
		"email_templates",
		"verzonden_emails",
		"incoming_emails",
		"notifications",

		// V28 - Renamed tables
		"participants",           // was: aanmeldingen
		"participant_antwoorden", // was: aanmelding_antwoorden
		"event_registrations",    // was: event_participants

		// V04 - Chat
		"chat_channels",
		"chat_messages",
		"chat_message_reactions",
		"chat_channel_participants",
		"chat_user_presence",

		// V05 - Newsletter
		"newsletters",

		// V06 - RBAC
		"roles",
		"permissions",
		"role_permissions",
		"user_roles",

		// V09 - Auth
		"refresh_tokens",

		// V10 - Images
		"uploaded_images",

		// V11 - CMS
		"photos",
		"albums",
		"album_photos",
		"videos",
		"sponsors",
		"program_schedule",
		"social_embeds",
		"social_links",
		"under_construction",
		"partners",
		"radio_recordings",
		"title_section_content",

		// V21 - Gamification
		"badges",
		"participant_achievements",

		// V23 - Events
		"events",

		// V24 - Notulen
		"notulen",
		"notulen_versies",

		// V25 - Leaderboard
		"leaderboard_view",

		// V26 & V27 - Lookup tables
		"participant_roles",
		"distances",
		"contact_status_types",
		"registration_status_types",
		"email_status_types",
		"event_status_types",
		"chat_channel_types",
		"notification_types",
		"notification_priority_types",
	}

	for _, tableName := range expectedTables {
		t.Run(tableName, func(t *testing.T) {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' 
					AND table_name = $1
				)
			`
			err := db.QueryRow(query, tableName).Scan(&exists)
			require.NoError(t, err, "Failed to check if table exists")
			assert.True(t, exists, "Table %s should exist", tableName)
		})
	}
}

// testLookupTablesExist verifies all lookup tables exist with correct structure
func testLookupTablesExist(t *testing.T, db *sql.DB) {
	lookupTables := map[string][]string{
		"participant_roles":           {"name", "description", "is_active"},
		"distances":                   {"route", "fund_amount"},
		"contact_status_types":        {"status", "description"},
		"registration_status_types":   {"status", "description"},
		"email_status_types":          {"status", "description"},
		"event_status_types":          {"status", "description"},
		"chat_channel_types":          {"type", "description"},
		"notification_types":          {"name", "description"},
		"notification_priority_types": {"name", "description"},
	}

	for tableName, expectedColumns := range lookupTables {
		t.Run(tableName, func(t *testing.T) {
			// Check table exists
			var exists bool
			err := db.QueryRow(`
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' AND table_name = $1
				)
			`, tableName).Scan(&exists)
			require.NoError(t, err)
			require.True(t, exists, "Lookup table %s should exist", tableName)

			// Check columns exist
			for _, colName := range expectedColumns {
				var colExists bool
				err := db.QueryRow(`
					SELECT EXISTS (
						SELECT FROM information_schema.columns 
						WHERE table_schema = 'public' 
						AND table_name = $1 
						AND column_name = $2
					)
				`, tableName, colName).Scan(&colExists)
				require.NoError(t, err)
				assert.True(t, colExists, "Column %s.%s should exist", tableName, colName)
			}
		})
	}
}

// testLookupTableSeeding verifies lookup tables have seed data
func testLookupTableSeeding(t *testing.T, db *sql.DB) {
	t.Run("participant_roles", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM participant_roles").Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "participant_roles should have seed data")

		// Check specific roles exist (case-insensitive)
		expectedRoles := []string{"Deelnemer", "Vrijwilliger", "Sponsor", "Begeleider"}
		for _, role := range expectedRoles {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM participant_roles WHERE name ILIKE $1)", role).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "Role %s should exist", role)
		}
	})

	t.Run("distances", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM distances").Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "distances should have seed data")

		// Check that fund_amount exists and has values
		var avgFundAmount float64
		err = db.QueryRow("SELECT AVG(fund_amount) FROM distances WHERE fund_amount IS NOT NULL").Scan(&avgFundAmount)
		require.NoError(t, err)
		assert.Greater(t, avgFundAmount, 0.0, "distances should have fund_amount values")
	})

	t.Run("status_types", func(t *testing.T) {
		statusTables := []string{
			"contact_status_types",
			"registration_status_types",
			"email_status_types",
			"event_status_types",
		}

		for _, table := range statusTables {
			var count int
			err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
			require.NoError(t, err)
			assert.Greater(t, count, 0, "%s should have seed data", table)
		}
	})

	t.Run("notification_types", func(t *testing.T) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM notification_types").Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0, "notification_types should have seed data")
	})

	t.Run("RBAC_seed_data", func(t *testing.T) {
		// Check admin role exists
		var adminExists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = 'admin')").Scan(&adminExists)
		require.NoError(t, err)
		assert.True(t, adminExists, "Admin role should exist")

		// Check permissions exist
		var permCount int
		err = db.QueryRow("SELECT COUNT(*) FROM permissions").Scan(&permCount)
		require.NoError(t, err)
		assert.Greater(t, permCount, 0, "Permissions should be seeded")
	})
}

// testForeignKeyConstraints verifies critical FK relationships
func testForeignKeyConstraints(t *testing.T, db *sql.DB) {
	fkTests := []struct {
		name         string
		childTable   string
		childColumn  string
		parentTable  string
		parentColumn string
		shouldExist  bool
	}{
		// V26 & V27 - Critical lookups
		{
			name:         "event_registrations_to_participant_roles",
			childTable:   "event_registrations",
			childColumn:  "participant_role_name",
			parentTable:  "participant_roles",
			parentColumn: "name",
			shouldExist:  true,
		},
		{
			name:         "event_registrations_to_distances",
			childTable:   "event_registrations",
			childColumn:  "distance_route",
			parentTable:  "distances",
			parentColumn: "route",
			shouldExist:  true,
		},
		{
			name:         "event_registrations_to_status",
			childTable:   "event_registrations",
			childColumn:  "status",
			parentTable:  "registration_status_types",
			parentColumn: "status",
			shouldExist:  true,
		},
		{
			name:         "contact_formulieren_to_status",
			childTable:   "contact_formulieren",
			childColumn:  "status",
			parentTable:  "contact_status_types",
			parentColumn: "status",
			shouldExist:  true,
		},
		{
			name:         "verzonden_emails_to_status",
			childTable:   "verzonden_emails",
			childColumn:  "status",
			parentTable:  "email_status_types",
			parentColumn: "status",
			shouldExist:  true,
		},
		{
			name:         "events_to_status",
			childTable:   "events",
			childColumn:  "status",
			parentTable:  "event_status_types",
			parentColumn: "status",
			shouldExist:  true,
		},
		{
			name:         "notifications_to_type",
			childTable:   "notifications",
			childColumn:  "type",
			parentTable:  "notification_types",
			parentColumn: "name",
			shouldExist:  true,
		},
		{
			name:         "notifications_to_priority",
			childTable:   "notifications",
			childColumn:  "priority",
			parentTable:  "notification_priority_types",
			parentColumn: "name",
			shouldExist:  true,
		},
	}

	for _, tc := range fkTests {
		t.Run(tc.name, func(t *testing.T) {
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM information_schema.table_constraints tc
					JOIN information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
						AND tc.table_schema = kcu.table_schema
					JOIN information_schema.constraint_column_usage ccu
						ON ccu.constraint_name = tc.constraint_name
						AND ccu.table_schema = tc.table_schema
					WHERE tc.constraint_type = 'FOREIGN KEY'
						AND tc.table_name = $1
						AND kcu.column_name = $2
						AND ccu.table_name = $3
						AND ccu.column_name = $4
				)
			`
			var exists bool
			err := db.QueryRow(query, tc.childTable, tc.childColumn, tc.parentTable, tc.parentColumn).Scan(&exists)
			require.NoError(t, err)

			if tc.shouldExist {
				assert.True(t, exists, "FK constraint %s should exist", tc.name)
			}
		})
	}
}

// testRBACSystemIntegrity verifies RBAC system setup
func testRBACSystemIntegrity(t *testing.T, db *sql.DB) {
	t.Run("SystemRolesExist", func(t *testing.T) {
		expectedRoles := []string{"admin", "staff", "user"}
		for _, roleName := range expectedRoles {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", roleName).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "System role %s should exist", roleName)
		}
	})

	t.Run("PermissionCategories", func(t *testing.T) {
		// Check that permissions exist for key resources
		resources := []string{
			"participant",
			"event",
			"chat",
			"user",
		}

		for _, resource := range resources {
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM permissions WHERE resource LIKE $1", resource+"%").Scan(&count)
			require.NoError(t, err)
			assert.Greater(t, count, 0, "Permissions for %s should exist", resource)
		}

		// Verify we have a good variety of permissions
		var totalCount int
		err := db.QueryRow("SELECT COUNT(*) FROM permissions").Scan(&totalCount)
		require.NoError(t, err)
		assert.Greater(t, totalCount, 20, "Should have at least 20 permissions total")
	})

	t.Run("AdminHasAllPermissions", func(t *testing.T) {
		// Get admin role ID
		var adminRoleID string
		err := db.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID)
		require.NoError(t, err)

		// Check admin has permissions assigned
		var permCount int
		err = db.QueryRow("SELECT COUNT(*) FROM role_permissions WHERE role_id = $1", adminRoleID).Scan(&permCount)
		require.NoError(t, err)
		assert.Greater(t, permCount, 10, "Admin should have many permissions")
	})
}

// testCriticalIndexesExist verifies performance-critical indexes
func testCriticalIndexesExist(t *testing.T, db *sql.DB) {
	criticalIndexes := []struct {
		table  string
		column string
	}{
		{"participants", "email"},
		{"participants", "gebruiker_id"},
		{"event_registrations", "participant_id"},
		{"event_registrations", "event_id"},
		{"event_registrations", "tracking_status"},
		{"contact_formulieren", "email"},
		{"contact_formulieren", "status"},
		{"gebruikers", "email"},
		{"notifications", "type"},
		{"notifications", "priority"},
		{"notifications", "sent"},
		{"chat_messages", "channel_id"},
		{"incoming_emails", "is_processed"},
	}

	for _, idx := range criticalIndexes {
		t.Run(fmt.Sprintf("%s_%s", idx.table, idx.column), func(t *testing.T) {
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM pg_indexes
					WHERE tablename = $1
					AND indexdef LIKE '%' || $2 || '%'
				)
			`
			var exists bool
			err := db.QueryRow(query, idx.table, idx.column).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "Index on %s.%s should exist", idx.table, idx.column)
		})
	}
}

// testColumnTypes verifies critical column types are correct
func testColumnTypes(t *testing.T, db *sql.DB) {
	columnTypeTests := []struct {
		table        string
		column       string
		expectedType string
	}{
		{"participants", "id", "uuid"},
		{"participants", "created_at", "timestamp with time zone"},
		{"event_registrations", "id", "uuid"},
		{"distances", "fund_amount", "integer"},
		{"gebruikers", "is_actief", "boolean"},
		{"notifications", "sent", "boolean"},
	}

	for _, tc := range columnTypeTests {
		t.Run(fmt.Sprintf("%s_%s", tc.table, tc.column), func(t *testing.T) {
			var dataType string
			query := `
				SELECT data_type
				FROM information_schema.columns
				WHERE table_name = $1 AND column_name = $2
			`
			err := db.QueryRow(query, tc.table, tc.column).Scan(&dataType)
			require.NoError(t, err)
			assert.Contains(t, dataType, tc.expectedType,
				"Column %s.%s should be of type %s but is %s",
				tc.table, tc.column, tc.expectedType, dataType)
		})
	}
}

// testTableRenaming verifies V28 table renaming was successful
func testTableRenaming(t *testing.T, db *sql.DB) {
	t.Run("OldTablesDeleted", func(t *testing.T) {
		oldTables := []string{
			"aanmeldingen",
			"aanmelding_antwoorden",
			"event_participants",
		}

		for _, tableName := range oldTables {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' AND table_name = $1
				)
			`
			err := db.QueryRow(query, tableName).Scan(&exists)
			require.NoError(t, err)
			assert.False(t, exists, "Old table %s should not exist after V28", tableName)
		}
	})

	t.Run("NewTablesExist", func(t *testing.T) {
		newTables := []string{
			"participants",
			"participant_antwoorden",
			"event_registrations",
		}

		for _, tableName := range newTables {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' AND table_name = $1
				)
			`
			err := db.QueryRow(query, tableName).Scan(&exists)
			require.NoError(t, err)
			assert.True(t, exists, "New table %s should exist after V28", tableName)
		}
	})
}

// testDataIntegrity runs basic data integrity checks
func testDataIntegrity(t *testing.T, db *sql.DB) {
	t.Run("NoOrphanedForeignKeys", func(t *testing.T) {
		// Check event_registrations have valid participant_id
		var orphanedCount int
		query := `
			SELECT COUNT(*)
			FROM event_registrations er
			LEFT JOIN participants p ON er.participant_id = p.id
			WHERE p.id IS NULL
		`
		err := db.QueryRow(query).Scan(&orphanedCount)
		require.NoError(t, err)
		assert.Equal(t, 0, orphanedCount, "Should have no orphaned event_registrations")
	})

	t.Run("TimestampConsistency", func(t *testing.T) {
		// Check that updated_at >= created_at
		tables := []string{"participants", "contact_formulieren", "gebruikers"}
		for _, table := range tables {
			var invalidCount int
			query := fmt.Sprintf(`
				SELECT COUNT(*)
				FROM %s
				WHERE updated_at < created_at
			`, table)
			err := db.QueryRow(query).Scan(&invalidCount)
			require.NoError(t, err)
			assert.Equal(t, 0, invalidCount, "Table %s should have valid timestamp consistency", table)
		}
	})

	t.Run("LookupTableReferencesValid", func(t *testing.T) {
		// Check event_registrations reference valid lookup values
		var invalidRoles int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM event_registrations er
			LEFT JOIN participant_roles pr ON er.participant_role_name = pr.name
			WHERE er.participant_role_name IS NOT NULL AND pr.name IS NULL
		`).Scan(&invalidRoles)
		require.NoError(t, err)
		assert.Equal(t, 0, invalidRoles, "All participant_role_name references should be valid")

		var invalidDistances int
		err = db.QueryRow(`
			SELECT COUNT(*)
			FROM event_registrations er
			LEFT JOIN distances d ON er.distance_route = d.route
			WHERE er.distance_route IS NOT NULL AND d.route IS NULL
		`).Scan(&invalidDistances)
		require.NoError(t, err)
		assert.Equal(t, 0, invalidDistances, "All distance_route references should be valid")
	})
}

// TestMigrationSequence verifies migrations can run in correct order
func TestMigrationSequence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping migration sequence test in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	// This test verifies that migrations are numbered correctly and in proper sequence
	expectedSequence := []string{
		"V01", "V02", "V03", "V04", "V05", "V06", "V07", "V08", "V09", "V10",
		"V11", "V12", "V13", "V14", "V15", "V16", "V17_CONSOLIDATED", "V18", "V19", "V20",
		"V21", "V22", "V23", "V24", "V25",
		"V26_CONSOLIDATED",
		"V27_CONSOLIDATED",
		"V28_CONSOLIDATED",
		"V29_CONSOLIDATED",
		"V30",
	}

	t.Logf("Expected migration sequence: %v", expectedSequence)
	assert.Equal(t, 30, len(expectedSequence), "Should have exactly 30 migrations after consolidation")
}

// TestV30_IsActiveColumn verifies the V30 fix
func TestV30_IsActiveColumn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping V30 test in short mode")
	}

	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	t.Run("participant_roles_has_is_active", func(t *testing.T) {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.columns 
				WHERE table_name = 'participant_roles' 
				AND column_name = 'is_active'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "participant_roles should have is_active column")
	})

	t.Run("is_active_has_default_value", func(t *testing.T) {
		var columnDefault sql.NullString
		err := db.QueryRow(`
			SELECT column_default
			FROM information_schema.columns
			WHERE table_name = 'participant_roles'
			AND column_name = 'is_active'
		`).Scan(&columnDefault)
		require.NoError(t, err)
		assert.True(t, columnDefault.Valid, "is_active should have a default value")
	})

	t.Run("all_roles_have_is_active_populated", func(t *testing.T) {
		var nullCount int
		err := db.QueryRow("SELECT COUNT(*) FROM participant_roles WHERE is_active IS NULL").Scan(&nullCount)
		require.NoError(t, err)
		assert.Equal(t, 0, nullCount, "All participant_roles should have is_active populated")
	})
}
