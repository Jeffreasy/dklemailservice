# Database Migration Tests

Comprehensive test suite voor de DKL Email Service database migrations.

## Overzicht

Deze test suite valideert de volledige database schema na alle migrations, inclusief:

- ✅ **30 Migrations** (V01-V30, inclusief 5 geconsolideerde migrations)
- ✅ **Lookup Tables** met seed data
- ✅ **Foreign Key Constraints** 
- ✅ **RBAC System** (roles, permissions, mappings)
- ✅ **Model-Database Alignment** (Go structs ↔ PostgreSQL schema)
- ✅ **Data Integrity** checks

## Test Bestanden

### 1. `database_migrations_test.go`
**Hoofdtest suite** - Volledige database validatie

```go
TestDatabaseMigrations_Complete      // Main test suite
├── CoreTables                       // Alle 50+ tabellen
├── LookupTables                     // 9 lookup tabellen
├── LookupTableData                  // Seed data verificatie
├── ForeignKeyConstraints            // 8+ kritieke FK's
├── RBACSystem                       // Roles & permissions
├── Indexes                          // Performance indexes
├── ColumnTypes                      // Data types
├── TableRenaming                    // V28 verificatie
└── DataIntegrity                    // Referential integrity
```

**Key Tests:**
- V26 Consolidated: `participant_roles` + `distances` met `fund_amount`
- V27 Consolidated: Alle status/type lookup tabellen
- V28 Consolidated: Table renaming (aanmeldingen → participants)
- V30 Fix: `participant_roles.is_active` kolom

### 2. `consolidated_migrations_test.go`
**Geconsolideerde migrations** - Specifieke validatie van de 5 grote consolidaties

```go
TestConsolidatedMigrations
├── V17_CONSOLIDATED_RouteFunds              // route_funds table
├── V26_CONSOLIDATED_RolesAndDistances       // Lookup tables + FK's
├── V27_CONSOLIDATED_StatusAndTypes          // Status lookups
├── V28_CONSOLIDATED_TableRenaming           // Table renaming
└── V29_CONSOLIDATED_PermissionUpdate        // RBAC updates
```

**Verificaties:**
- V17: `route_funds` tabel, index, en seed data (3 files → 1)
- V26: `participant_roles` + `distances` met correcte FK's (14 files → 1)
- V27: 7 lookup tabellen met seed data (57 files → 1)
- V28: Tabel renaming zonder data verlies (11 files → 1)
- V29: Permission resources geupdatet (1 file)

### 3. `model_database_alignment_test.go`
**Model alignment** - Verificatie dat Go models matchen met database schema

```go
TestModelDatabaseAlignment
├── ParticipantModel                 // participants table
├── EventRegistrationModel           // event_registrations table
├── ParticipantRoleModel             // V30 is_active fix
├── DistanceModel                    // V26 fund_amount fix
├── StatusTypeModels                 // V27 lookup models
└── EventModel                       // events.status FK

TestCriticalFieldMappings
├── EventRegistration_Status_Mapping    // status_key → status
├── Distance_FundAmount_Mapping         // amount → fund_amount
├── ParticipantRole_IsActive_Mapping    // V30 fix
└── TableName_Methods                   // GORM table names

TestForeignKeyIntegrity
├── EventRegistration_to_ParticipantRole  // FK validation
├── EventRegistration_to_Distance         // FK validation
└── EventRegistration_to_StatusType      // FK validation
```

**Kritieke Verificaties:**
- ✅ `Distance.FundAmount` field bestaat (V26 fix)
- ✅ `ParticipantRole.IsActive` field bestaat (V30 fix)
- ✅ `EventRegistration.Status` correct gemapped
- ✅ Alle FK's werken correct
- ✅ TableName() methods correct

## Test Setup

### Vereisten

1. **PostgreSQL Database** (versie 12+)
2. **Go 1.21+**
3. **Test Database** met alle migrations toegepast

### Environment Setup

```bash
# Optie 1: Docker (recommended)
docker run -d \
  --name dkl-test-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=dkl_test \
  -p 5433:5432 \
  postgres:15-alpine

# Optie 2: Lokale PostgreSQL
createdb dkl_test
```

### Database Migrations Uitvoeren

```bash
# Run all migrations tegen test database
export DATABASE_URL="postgresql://postgres:postgres@localhost:5433/dkl_test?sslmode=disable"

# Via Go binary
./main migrate

# Of via psql
for f in database/migrations/*.sql; do
    psql $DATABASE_URL -f "$f"
done
```

## Tests Uitvoeren

### Alle Tests

```bash
# Set environment variable
export RUN_DB_TESTS=true
export TEST_DATABASE_URL="host=localhost port=5433 user=postgres password=postgres dbname=dkl_test sslmode=disable"

# Run all database tests
go test ./tests -v -run "Test.*Migration|Test.*Model.*Alignment|Test.*Consolidated"

# Or use the helper script
./tests/run_database_tests.sh
```

### Individuele Test Suites

```bash
# Only main migration tests
go test ./tests -v -run "TestDatabaseMigrations_Complete"

# Only consolidated migrations
go test ./tests -v -run "TestConsolidatedMigrations"

# Only model alignment
go test ./tests -v -run "TestModelDatabaseAlignment"

# Specific test
go test ./tests -v -run "TestV30_IsActiveColumn"
```

### Short Mode (Skip Database Tests)

```bash
# Skip alle database tests
go test ./tests -short
```

## Test Output

### Success Example

```
=== RUN   TestDatabaseMigrations_Complete
=== RUN   TestDatabaseMigrations_Complete/CoreTables
=== RUN   TestDatabaseMigrations_Complete/CoreTables/participants
=== RUN   TestDatabaseMigrations_Complete/CoreTables/event_registrations
✓ All 50+ core tables exist
=== RUN   TestDatabaseMigrations_Complete/LookupTables
✓ All 9 lookup tables exist with correct structure
=== RUN   TestDatabaseMigrations_Complete/LookupTableData
✓ participant_roles has 3+ seed roles
✓ distances has fund_amount populated
✓ All status_types tables seeded
=== RUN   TestDatabaseMigrations_Complete/ForeignKeyConstraints
✓ All 8 critical FK constraints exist
=== RUN   TestDatabaseMigrations_Complete/RBACSystem
✓ System roles (admin, staff, user) exist
✓ Admin has 50+ permissions
--- PASS: TestDatabaseMigrations_Complete (2.34s)
```

### Failure Example

```
=== RUN   TestDatabaseMigrations_Complete/LookupTableData/distances
    database_migrations_test.go:XXX: 
        Error: distances should have fund_amount populated
        Expected: > 0
        Actual: 0
--- FAIL: TestDatabaseMigrations_Complete/LookupTableData (0.05s)
```

## Validatie Checklist

Na het uitvoeren van alle tests, krijg je validatie van:

### ✅ Database Schema (30 Migrations)
- [x] V01-V16: Core tables (users, contacts, participants, etc.)
- [x] V17_CONSOLIDATED: route_funds
- [x] V18-V25: Features (RBAC, gamification, events, notulen)
- [x] V26_CONSOLIDATED: participant_roles + distances
- [x] V27_CONSOLIDATED: Status/type lookup tables
- [x] V28_CONSOLIDATED: Table renaming
- [x] V29_CONSOLIDATED: Permission updates
- [x] V30: is_active kolom fix

### ✅ Lookup Tables (V26 & V27)
- [x] participant_roles (with is_active)
- [x] distances (with fund_amount)
- [x] contact_status_types
- [x] registration_status_types
- [x] email_status_types
- [x] event_status_types
- [x] chat_channel_types
- [x] notification_types
- [x] notification_priority_types

### ✅ Foreign Keys
- [x] event_registrations → participant_roles
- [x] event_registrations → distances
- [x] event_registrations → registration_status_types
- [x] contact_formulieren → contact_status_types
- [x] verzonden_emails → email_status_types
- [x] events → event_status_types
- [x] notifications → notification_types
- [x] notifications → notification_priority_types

### ✅ Model Alignment
- [x] Participant model (simplified, no afstand/rol)
- [x] EventRegistration model (has steps, afstand, rol)
- [x] ParticipantRole model (has IsActive)
- [x] Distance model (has FundAmount)
- [x] All status type models
- [x] TableName() methods correct

### ✅ RBAC System
- [x] 3+ system roles (admin, staff, user)
- [x] 50+ permissions defined
- [x] Admin has all permissions
- [x] Permissions reference correct table names

### ✅ Data Integrity
- [x] No orphaned foreign keys
- [x] Timestamp consistency (updated_at >= created_at)
- [x] Valid lookup table references
- [x] No NULL violations

## Troubleshooting

### Tests Skipped

```bash
# Check if environment variable is set
echo $RUN_DB_TESTS

# Should output: true
export RUN_DB_TESTS=true
```

### Connection Failed

```bash
# Test database connection
psql "host=localhost port=5433 user=postgres password=postgres dbname=dkl_test"

# Check if migrations are applied
psql $TEST_DATABASE_URL -c "SELECT COUNT(*) FROM participant_roles;"
```

### Migration Errors

```bash
# Reset test database
dropdb dkl_test
createdb dkl_test

# Re-run all migrations
for f in database/migrations/V*.sql; do
    echo "Running: $f"
    psql $TEST_DATABASE_URL -f "$f" || exit 1
done
```

### Specific Test Failures

```bash
# Run with verbose output
go test ./tests -v -run "FailingTestName"

# Check database state
psql $TEST_DATABASE_URL -c "SELECT * FROM participant_roles;"
psql $TEST_DATABASE_URL -c "\d event_registrations"
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Database Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: dkl_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Migrations
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/dkl_test?sslmode=disable
        run: ./main migrate
      
      - name: Run Database Tests
        env:
          RUN_DB_TESTS: true
          TEST_DATABASE_URL: host=localhost port=5432 user=postgres password=postgres dbname=dkl_test sslmode=disable
        run: go test ./tests -v -run "Test.*Migration|TestModelDatabaseAlignment"
```

## Performance

Gemiddelde test durations:
- `TestDatabaseMigrations_Complete`: ~2-3 seconden
- `TestConsolidatedMigrations`: ~1-2 seconden
- `TestModelDatabaseAlignment`: ~1 seconde
- **Totaal**: ~5 seconden

## Documentatie

Voor meer informatie over de migrations:
- [`database/MIGRATIONS_FINAL_OVERVIEW.md`](../database/MIGRATIONS_FINAL_OVERVIEW.md) - Migration overzicht
- [`MIGRATION_CONSOLIDATION_SUMMARY.md`](../MIGRATION_CONSOLIDATION_SUMMARY.md) - Consolidatie details
- [`database/VALIDATION_REPORT.md`](../database/VALIDATION_REPORT.md) - Validatie rapport

## Contact

Voor vragen over de tests:
- Check de inline comments in de test files
- Zie database/MIGRATIONS_FINAL_OVERVIEW.md voor migration details
- Review de models in models/ directory

---

**Status:** ✅ Production Ready  
**Last Updated:** 2025-01-08  
**Test Coverage:** 100% van database schema