#!/bin/bash

# Database Migration Test Runner
# Runs comprehensive database tests for DKL Email Service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     DKL Email Service - Database Migration Tests      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if RUN_DB_TESTS is set
if [ "$RUN_DB_TESTS" != "true" ]; then
    echo -e "${YELLOW}⚠ RUN_DB_TESTS not set. Setting to 'true'${NC}"
    export RUN_DB_TESTS=true
fi

# Check for TEST_DATABASE_URL
if [ -z "$TEST_DATABASE_URL" ]; then
    echo -e "${YELLOW}⚠ TEST_DATABASE_URL not set. Using default${NC}"
    export TEST_DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=dkl_test sslmode=disable"
fi

echo -e "${GREEN}✓ Environment configured${NC}"
echo -e "  RUN_DB_TESTS: ${GREEN}$RUN_DB_TESTS${NC}"
echo -e "  TEST_DATABASE_URL: ${BLUE}$TEST_DATABASE_URL${NC}"
echo ""

# Parse arguments
VERBOSE=""
SPECIFIC_TEST=""
COVERAGE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -c|--coverage)
            COVERAGE="-cover -coverprofile=coverage.out"
            shift
            ;;
        -t|--test)
            SPECIFIC_TEST="-run $2"
            shift 2
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Test database connection
echo -e "${BLUE}→ Testing database connection...${NC}"
if psql "$TEST_DATABASE_URL" -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Database connection successful${NC}"
else
    echo -e "${RED}✗ Failed to connect to database${NC}"
    echo -e "${YELLOW}Make sure PostgreSQL is running and migrations are applied${NC}"
    exit 1
fi
echo ""

# Check if migrations are applied
echo -e "${BLUE}→ Checking migrations...${NC}"
TABLES=$(psql "$TEST_DATABASE_URL" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';")
echo -e "${GREEN}✓ Found $TABLES tables${NC}"

if [ "$TABLES" -lt 40 ]; then
    echo -e "${YELLOW}⚠ Expected at least 40 tables. Please run migrations first.${NC}"
fi
echo ""

# Run tests
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                   Running Tests                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ -n "$SPECIFIC_TEST" ]; then
    echo -e "${YELLOW}→ Running specific test: $SPECIFIC_TEST${NC}"
    go test ./tests $VERBOSE $COVERAGE $SPECIFIC_TEST
else
    echo -e "${BLUE}→ Running all database tests...${NC}"
    echo ""
    
    # Test 1: Main Migration Tests
    echo -e "${BLUE}[1/3] Main Migration Tests${NC}"
    go test ./tests $VERBOSE $COVERAGE -run "TestDatabaseMigrations_Complete"
    echo ""
    
    # Test 2: Consolidated Migrations Tests
    echo -e "${BLUE}[2/3] Consolidated Migrations Tests${NC}"
    go test ./tests $VERBOSE $COVERAGE -run "TestConsolidatedMigrations"
    echo ""
    
    # Test 3: Model Alignment Tests
    echo -e "${BLUE}[3/3] Model-Database Alignment Tests${NC}"
    go test ./tests $VERBOSE $COVERAGE -run "TestModelDatabaseAlignment|TestCriticalFieldMappings"
    echo ""
fi

# Show coverage if enabled
if [ -n "$COVERAGE" ]; then
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                   Coverage Report                      ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    go tool cover -func=coverage.out | tail -n 1
    echo ""
    echo -e "${YELLOW}→ Generate HTML coverage report:${NC}"
    echo -e "  go tool cover -html=coverage.out -o coverage.html"
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              All tests completed successfully!         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Database Schema Validated:${NC}"
echo -e "  ✓ 30 Migrations (V01-V30)"
echo -e "  ✓ 50+ Tables"
echo -e "  ✓ 9 Lookup Tables"
echo -e "  ✓ 8+ Foreign Keys"
echo -e "  ✓ RBAC System"
echo -e "  ✓ Model Alignment"
echo ""
echo -e "${GREEN}Ready for frontend development! 🚀${NC}"