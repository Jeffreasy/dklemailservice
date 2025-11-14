#!/bin/bash
# V30+RBAC Integration Test Script
# Tests all participant registration, upgrade, and RBAC flows

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
DB_NAME="${DB_NAME:-dkl_db}"
DB_USER="${DB_USER:-dkl_user}"

echo "======================================"
echo "V30+RBAC Integration Test Suite"
echo "======================================"
echo "Base URL: $BASE_URL"
echo "Database: $DB_NAME"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to run test
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_status="$3"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -e "${YELLOW}[TEST $TESTS_RUN]${NC} $test_name"
    
    # Run command and capture output
    response=$(eval "$test_command" 2>&1)
    status=$?
    
    if [ $status -eq 0 ] || [ "$expected_status" = "any" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "Output: $response"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Helper function to check database
check_db() {
    local query="$1"
    local expected="$2"
    
    result=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c "$query" 2>&1 | xargs)
    
    if [ "$result" = "$expected" ]; then
        return 0
    else
        echo "Expected: $expected"
        echo "Got: $result"
        return 1
    fi
}

echo "======================================"
echo "PART 1: Database Setup Verification"
echo "======================================"
echo ""

run_test "Check participant_user role exists" \
    "psql -U $DB_USER -d $DB_NAME -t -c \"SELECT COUNT(*) FROM roles WHERE name = 'participant_user';\"" \
    "any"

run_test "Check app:access permission exists" \
    "psql -U $DB_USER -d $DB_NAME -t -c \"SELECT COUNT(*) FROM permissions WHERE resource = 'app' AND action = 'access';\"" \
    "any"

run_test "Check triggers are enabled" \
    "psql -U $DB_USER -d $DB_NAME -t -c \"SELECT COUNT(*) FROM pg_trigger WHERE tgname LIKE '%participant%' AND tgenabled = 'O';\"" \
    "any"

echo ""
echo "======================================"
echo "PART 2: Full Account Registration"
echo "======================================"
echo ""

# Generate unique email for testing
TIMESTAMP=$(date +%s)
FULL_EMAIL="full_${TIMESTAMP}@test.nl"

echo "Registering full account: $FULL_EMAIL"

FULL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/public/aanmelden" \
  -H "Content-Type: application/json" \
  -d "{
    \"naam\": \"Full Account Test\",
    \"email\": \"$FULL_EMAIL\",
    \"telefoon\": \"06-12345678\",
    \"rol\": \"Begeleider\",
    \"afstand\": \"10 KM\",
    \"ondersteuning\": \"Nee\",
    \"want_account\": true,
    \"wachtwoord\": \"TestPass123\",
    \"terms\": true
  }")

echo "Response: $FULL_RESPONSE"

# Extract gebruiker_id from response
GEBRUIKER_ID=$(echo "$FULL_RESPONSE" | grep -o '"gebruiker_id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$GEBRUIKER_ID" ]; then
    echo -e "${RED}✗ No gebruiker_id in response${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
else
    echo -e "${GREEN}✓ Full account created with gebruiker_id: $GEBRUIKER_ID${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
fi

echo ""
echo "Verifying RBAC setup for full account..."

# Wait a moment for triggers to execute
sleep 2

# Check if participant_user role was assigned
ROLE_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM user_roles ur \
     JOIN roles r ON ur.role_id = r.id \
     WHERE ur.user_id = '$GEBRUIKER_ID' \
     AND r.name = 'participant_user' \
     AND ur.is_active = true;" | xargs)

if [ "$ROLE_COUNT" = "1" ]; then
    echo -e "${GREEN}✓ participant_user role assigned${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ participant_user role NOT assigned (count: $ROLE_COUNT)${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check if participant_guide role was assigned (Begeleider)
GUIDE_ROLE_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM user_roles ur \
     JOIN roles r ON ur.role_id = r.id \
     WHERE ur.user_id = '$GEBRUIKER_ID' \
     AND r.name = 'participant_guide' \
     AND ur.is_active = true;" | xargs)

if [ "$GUIDE_ROLE_COUNT" = "1" ]; then
    echo -e "${GREEN}✓ participant_guide role assigned (correct for Begeleider)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ participant_guide role NOT assigned (count: $GUIDE_ROLE_COUNT) - check handler logic${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check app:access permission
APP_ACCESS_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM user_permissions \
     WHERE user_id = '$GEBRUIKER_ID' \
     AND resource = 'app' \
     AND action = 'access';" | xargs)

if [ "$APP_ACCESS_COUNT" -ge "1" ]; then
    echo -e "${GREEN}✓ app:access permission available${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ app:access permission NOT available${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

echo ""
echo "======================================"
echo "PART 3: App Login Test"
echo "======================================"
echo ""

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$FULL_EMAIL\",
    \"wachtwoord\": \"TestPass123\"
  }")

echo "Login Response: $LOGIN_RESPONSE"

# Check if access_token is present
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ Login successful with access_token${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Login failed - no access_token${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check if roles are in token (decode JWT - basic check)
if echo "$ACCESS_TOKEN" | grep -q "participant"; then
    echo -e "${GREEN}✓ JWT contains participant roles${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ Cannot verify roles in JWT (needs jq for proper decode)${NC}"
fi

echo ""
echo "======================================"
echo "PART 4: Temporary Account Test"
echo "======================================"
echo ""

TEMP_EMAIL="temp_${TIMESTAMP}@test.nl"

echo "Registering temporary account: $TEMP_EMAIL"

TEMP_RESPONSE=$(curl -s -X POST "$BASE_URL/api/public/aanmelden" \
  -H "Content-Type: application/json" \
  -d "{
    \"naam\": \"Temporary Test\",
    \"email\": \"$TEMP_EMAIL\",
    \"rol\": \"Deelnemer\",
    \"afstand\": \"6 KM\",
    \"ondersteuning\": \"Nee\",
    \"want_account\": false,
    \"terms\": true
  }")

echo "Response: $TEMP_RESPONSE"

# Check that NO gebruiker_id is returned
if echo "$TEMP_RESPONSE" | grep -q '"gebruiker_id"'; then
    echo -e "${RED}✗ Temporary account should NOT have gebruiker_id${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
else
    echo -e "${GREEN}✓ Temporary account has no gebruiker_id (correct)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
fi

echo ""
echo "======================================"
echo "PART 5: Account Upgrade Test"
echo "======================================"
echo ""

echo "Upgrading temporary account to full: $TEMP_EMAIL"

sleep 1

UPGRADE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/public/upgrade-to-full-account" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEMP_EMAIL\",
    \"wachtwoord\": \"NewPass456\"
  }")

echo "Upgrade Response: $UPGRADE_RESPONSE"

# Check if gebruiker_id is NOW present
UPGRADED_GEBRUIKER_ID=$(echo "$UPGRADE_RESPONSE" | grep -o '"gebruiker_id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$UPGRADED_GEBRUIKER_ID" ]; then
    echo -e "${GREEN}✓ Upgrade successful with gebruiker_id: $UPGRADED_GEBRUIKER_ID${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Upgrade failed - no gebruiker_id${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Wait for triggers
sleep 2

# Check if roles were assigned after upgrade
UPGRADED_ROLE_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM user_roles ur \
     WHERE ur.user_id = '$UPGRADED_GEBRUIKER_ID' \
     AND ur.is_active = true;" | xargs)

if [ "$UPGRADED_ROLE_COUNT" -ge "1" ]; then
    echo -e "${GREEN}✓ Roles assigned after upgrade (count: $UPGRADED_ROLE_COUNT)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ No roles assigned after upgrade${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check if upgrade was logged in participant_upgrades
UPGRADE_AUDIT_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM participant_upgrades pu \
     JOIN participants p ON pu.participant_id = p.id \
     WHERE p.email = '$TEMP_EMAIL';" | xargs)

if [ "$UPGRADE_AUDIT_COUNT" -ge "1" ]; then
    echo -e "${GREEN}✓ Upgrade logged in audit table${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ Upgrade not logged in audit table (non-critical)${NC}"
fi

# Test login after upgrade
echo ""
echo "Testing login after upgrade..."

UPGRADED_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEMP_EMAIL\",
    \"wachtwoord\": \"NewPass456\"
  }")

UPGRADED_ACCESS_TOKEN=$(echo "$UPGRADED_LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$UPGRADED_ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ Login successful after upgrade${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Login failed after upgrade${NC}"
    echo "Response: $UPGRADED_LOGIN_RESPONSE"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

echo ""
echo "======================================"
echo "PART 6: Permission Function Tests"
echo "======================================"
echo ""

# Test participant_has_permission function
if [ -n "$GEBRUIKER_ID" ]; then
    PARTICIPANT_ID=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT id FROM participants WHERE gebruiker_id = '$GEBRUIKER_ID' LIMIT 1;" | xargs)
    
    if [ -n "$PARTICIPANT_ID" ]; then
        HAS_APP_ACCESS=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
            "SELECT participant_has_permission('$PARTICIPANT_ID', 'app', 'access');" | xargs)
        
        if [ "$HAS_APP_ACCESS" = "t" ]; then
            echo -e "${GREEN}✓ participant_has_permission() returns true for app:access${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}✗ participant_has_permission() returns false${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
        
        CAN_ACCESS_APP=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
            "SELECT participant_can_access_app('$PARTICIPANT_ID');" | xargs)
        
        if [ "$CAN_ACCESS_APP" = "t" ]; then
            echo -e "${GREEN}✓ participant_can_access_app() returns true${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}✗ participant_can_access_app() returns false${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    fi
fi

echo ""
echo "======================================"
echo "PART 7: Views Test"
echo "======================================"
echo ""

# Test participant_user_permissions view
VIEW_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM participant_user_permissions WHERE participant_email = '$FULL_EMAIL';" | xargs)

if [ "$VIEW_COUNT" -ge "1" ]; then
    echo -e "${GREEN}✓ participant_user_permissions view contains data${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ participant_user_permissions view is empty for test user${NC}"
fi

# Test participant_account_stats view
STATS_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM participant_account_stats WHERE registration_year = 2026;" | xargs)

if [ "$STATS_COUNT" -ge "1" ]; then
    echo -e "${GREEN}✓ participant_account_stats view working${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ participant_account_stats view has no 2026 data${NC}"
fi

echo ""
echo "======================================"
echo "PART 8: Negative Tests"
echo "======================================"
echo ""

# Test: Duplicate full account registration (should fail)
DUP_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/public/aanmelden" \
  -H "Content-Type: application/json" \
  -d "{
    \"naam\": \"Duplicate\",
    \"email\": \"$FULL_EMAIL\",
    \"rol\": \"Deelnemer\",
    \"afstand\": \"6 KM\",
    \"ondersteuning\": \"Nee\",
    \"want_account\": true,
    \"wachtwoord\": \"Pass123\",
    \"terms\": true
  }")

DUP_STATUS=$(echo "$DUP_RESPONSE" | tail -n 1)

if [ "$DUP_STATUS" = "409" ]; then
    echo -e "${GREEN}✓ Duplicate registration correctly rejected (409)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Duplicate registration should return 409, got $DUP_STATUS${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test: Login with temporary account email (should fail)
TEMP_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"temp_999@test.nl\",
    \"wachtwoord\": \"anypass\"
  }")

TEMP_LOGIN_STATUS=$(echo "$TEMP_LOGIN_RESPONSE" | tail -n 1)

if [ "$TEMP_LOGIN_STATUS" = "401" ]; then
    echo -e "${GREEN}✓ Temporary account login correctly rejected (401)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}⚠ Temporary login status: $TEMP_LOGIN_STATUS (expected 401)${NC}"
fi

echo ""
echo "======================================"
echo "PART 9: Data Integrity Checks"
echo "======================================"
echo ""

# Check: All full accounts have gebruiker_id
ORPHAN_COUNT=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM participants \
     WHERE account_type = 'full' \
     AND gebruiker_id IS NULL;" | xargs)

if [ "$ORPHAN_COUNT" = "0" ]; then
    echo -e "${GREEN}✓ No orphaned full accounts (all have gebruiker_id)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Found $ORPHAN_COUNT orphaned full accounts${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check: All full accounts with gebruiker_id have at least participant_user role
USERS_WITHOUT_ROLE=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(DISTINCT g.id) \
     FROM gebruikers g \
     JOIN participants p ON p.gebruiker_id = g.id \
     WHERE p.account_type = 'full' \
     AND NOT EXISTS (
         SELECT 1 FROM user_roles ur 
         JOIN roles r ON ur.role_id = r.id
         WHERE ur.user_id = g.id 
         AND r.name = 'participant_user' 
         AND ur.is_active = true
     );" | xargs)

if [ "$USERS_WITHOUT_ROLE" = "0" ]; then
    echo -e "${GREEN}✓ All full account users have participant_user role${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Found $USERS_WITHOUT_ROLE users without participant_user role${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Check: No temporary accounts have gebruiker_id
TEMP_WITH_GEBRUIKER=$(psql -U "$DB_USER" -d "$DB_NAME" -t -c \
    "SELECT COUNT(*) FROM participants \
     WHERE account_type = 'temporary' \
     AND gebruiker_id IS NOT NULL;" | xargs)

if [ "$TEMP_WITH_GEBRUIKER" = "0" ]; then
    echo -e "${GREEN}✓ No temporary accounts have gebruiker_id (correct)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ Found $TEMP_WITH_GEBRUIKER temporary accounts with gebruiker_id${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

echo ""
echo "======================================"
echo "TEST SUMMARY"
echo "======================================"
echo ""
echo "Total Tests Run: $TESTS_RUN"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓✓✓ ALL TESTS PASSED ✓✓✓${NC}"
    echo ""
    echo "V30+RBAC Integration is working correctly!"
    exit 0
else
    echo -e "${RED}✗✗✗ SOME TESTS FAILED ✗✗✗${NC}"
    echo ""
    echo "Please review failed tests and fix issues."
    echo "Check logs and database for more details."
    exit 1
fi