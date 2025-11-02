#!/bin/bash
# =====================================================
# RBAC Database Verification Script
# Version: 1.49
# Datum: 2025-11-02
# =====================================================
# 
# Dit script voert automatisch de RBAC database verificatie uit
# in je Docker of lokale PostgreSQL omgeving
#
# Gebruik:
#   chmod +x database/scripts/check_rbac_db.sh
#   ./database/scripts/check_rbac_db.sh
# =====================================================

set -e  # Exit op eerste error

# Kleuren voor output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================================="
echo -e "🔍 RBAC Database Verification"
echo -e "===================================================${NC}"
echo ""

# =====================================================
# STAP 1: Detecteer Database Environment
# =====================================================
echo -e "${BLUE}[1/5] Detecteren database environment...${NC}"

# Check voor docker compose
if docker-compose ps db &> /dev/null; then
    echo -e "${GREEN}✓ Docker Compose gevonden${NC}"
    DB_CONTAINER=$(docker-compose ps -q db)
    DB_USER=${DB_USER:-dkluser}
    DB_NAME=${DB_NAME:-dklemailservice}
    CONNECTION_TYPE="docker-compose"
    
elif docker ps --filter "name=postgres" --format "{{.Names}}" | head -1 &> /dev/null; then
    echo -e "${GREEN}✓ Docker container gevonden${NC}"
    DB_CONTAINER=$(docker ps --filter "name=postgres" --format "{{.Names}}" | head -1)
    DB_USER=${DB_USER:-dkluser}
    DB_NAME=${DB_NAME:-dklemailservice}
    CONNECTION_TYPE="docker"
    
elif command -v psql &> /dev/null; then
    echo -e "${GREEN}✓ Lokale psql gevonden${NC}"
    DB_HOST=${DB_HOST:-localhost}
    DB_PORT=${DB_PORT:-5432}
    DB_USER=${DB_USER:-dkluser}
    DB_NAME=${DB_NAME:-dklemailservice}
    CONNECTION_TYPE="local"
    
else
    echo -e "${RED}✗ Geen database verbinding gevonden!${NC}"
    echo "  Installeer Docker of PostgreSQL"
    exit 1
fi

echo -e "Connection type: ${GREEN}${CONNECTION_TYPE}${NC}"
echo ""

# =====================================================
# STAP 2: Test Database Connectivity
# =====================================================
echo -e "${BLUE}[2/5] Testen database connectiviteit...${NC}"

if [[ "$CONNECTION_TYPE" == "docker-compose" ]] || [[ "$CONNECTION_TYPE" == "docker" ]]; then
    if docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Database verbinding succesvol${NC}"
    else
        echo -e "${RED}✗ Database verbinding mislukt${NC}"
        echo "  Container: $DB_CONTAINER"
        echo "  User: $DB_USER"
        echo "  Database: $DB_NAME"
        exit 1
    fi
else
    if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Database verbinding succesvol${NC}"
    else
        echo -e "${RED}✗ Database verbinding mislukt${NC}"
        exit 1
    fi
fi

echo ""

# =====================================================
# STAP 3: Quick Health Check
# =====================================================
echo -e "${BLUE}[3/5] Quick health check...${NC}"

QUERY="
SELECT 
    (SELECT COUNT(*) FROM roles WHERE is_system_role = true) as system_roles,
    (SELECT COUNT(*) FROM permissions WHERE is_system_permission = true) as system_permissions,
    (SELECT COUNT(DISTINCT user_id) FROM user_roles WHERE is_active = true) as users_with_roles,
    (SELECT COUNT(*) FROM gebruikers WHERE is_actief = true) as active_users;
"

if [[ "$CONNECTION_TYPE" == "docker-compose" ]] || [[ "$CONNECTION_TYPE" == "docker" ]]; then
    RESULT=$(docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -t -c "$QUERY")
else
    RESULT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "$QUERY")
fi

# Parse result
read -r SYSTEM_ROLES SYSTEM_PERMS USERS_WITH_ROLES ACTIVE_USERS <<< "$RESULT"

echo "System Roles:        $SYSTEM_ROLES (expected: 9)"
echo "System Permissions:  $SYSTEM_PERMS (expected: 58+)"
echo "Users with Roles:    $USERS_WITH_ROLES"
echo "Active Users:        $ACTIVE_USERS"

# Validate
HEALTH_OK=true

if [ "$SYSTEM_ROLES" -ne 9 ]; then
    echo -e "${RED}✗ System roles count mismatch!${NC}"
    HEALTH_OK=false
else
    echo -e "${GREEN}✓ System roles OK${NC}"
fi

if [ "$SYSTEM_PERMS" -lt 58 ]; then
    echo -e "${YELLOW}⚠ Fewer system permissions than expected${NC}"
    HEALTH_OK=false
else
    echo -e "${GREEN}✓ System permissions OK${NC}"
fi

if [ "$USERS_WITH_ROLES" -lt "$ACTIVE_USERS" ]; then
    echo -e "${YELLOW}⚠ Some users don't have RBAC roles${NC}"
    HEALTH_OK=false
else
    echo -e "${GREEN}✓ All users have roles${NC}"
fi

echo ""

# =====================================================
# STAP 4: Run Full Verification Script
# =====================================================
echo -e "${BLUE}[4/5] Running full verification script...${NC}"
echo -e "${YELLOW}(This may take a moment)${NC}"
echo ""

SCRIPT_PATH="database/scripts/verify_rbac_tables.sql"

if [ ! -f "$SCRIPT_PATH" ]; then
    echo -e "${RED}✗ Verification script not found: $SCRIPT_PATH${NC}"
    exit 1
fi

# Run verification
if [[ "$CONNECTION_TYPE" == "docker-compose" ]] || [[ "$CONNECTION_TYPE" == "docker" ]]; then
    docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < $SCRIPT_PATH | tee verification_output.log
else
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $SCRIPT_PATH | tee verification_output.log
fi

echo ""

# =====================================================
# STAP 5: Summary
# =====================================================
echo -e "${BLUE}[5/5] Generating summary...${NC}"
echo ""

# Check voor failures in output
FAIL_COUNT=$(grep -c "✗ FAIL" verification_output.log || true)
WARN_COUNT=$(grep -c "⚠ WARNING" verification_output.log || true)

echo -e "${BLUE}==================================================="
echo -e "📊 VERIFICATION SUMMARY"
echo -e "===================================================${NC}"
echo ""
echo "Environment:     $CONNECTION_TYPE"
echo "Database:        $DB_NAME"
echo "User:            $DB_USER"
if [[ "$CONNECTION_TYPE" == "docker-compose" ]] || [[ "$CONNECTION_TYPE" == "docker" ]]; then
    echo "Container:       $DB_CONTAINER"
else
    echo "Host:            $DB_HOST:$DB_PORT"
fi
echo ""
echo "System Roles:    $SYSTEM_ROLES / 9"
echo "Permissions:     $SYSTEM_PERMS / 58+"
echo "Users w/ Roles:  $USERS_WITH_ROLES / $ACTIVE_USERS"
echo ""

if [ "$FAIL_COUNT" -gt 0 ]; then
    echo -e "${RED}✗ FAILED CHECKS: $FAIL_COUNT${NC}"
    echo -e "${RED}URGENT: Review verification_output.log for details${NC}"
    exit 1
elif [ "$WARN_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}⚠ WARNINGS: $WARN_COUNT${NC}"
    echo -e "${YELLOW}Review verification_output.log for details${NC}"
    exit 0
else
    echo -e "${GREEN}✓ ALL CHECKS PASSED${NC}"
    echo -e "${GREEN}RBAC database is healthy!${NC}"
    exit 0
fi

echo ""
echo "Full output saved to: verification_output.log"
echo ""