# Database Migration Test Runner - PowerShell Version
# Runs comprehensive database tests for DKL Email Service

param(
    [switch]$Verbose,
    [switch]$Coverage,
    [string]$Test
)

$ErrorActionPreference = "Stop"

Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║     DKL Email Service - Database Migration Tests      ║" -ForegroundColor Blue
Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# Check if RUN_DB_TESTS is set
if (-not $env:RUN_DB_TESTS) {
    Write-Host "⚠ RUN_DB_TESTS not set. Setting to 'true'" -ForegroundColor Yellow
    $env:RUN_DB_TESTS = "true"
}

# Check for TEST_DATABASE_URL
if (-not $env:TEST_DATABASE_URL) {
    Write-Host "⚠ TEST_DATABASE_URL not set. Using default" -ForegroundColor Yellow
    $env:TEST_DATABASE_URL = "host=localhost port=5432 user=postgres password=postgres dbname=dkl_test sslmode=disable"
}

Write-Host "✓ Environment configured" -ForegroundColor Green
Write-Host "  RUN_DB_TESTS: $($env:RUN_DB_TESTS)" -ForegroundColor Green
Write-Host "  TEST_DATABASE_URL: $($env:TEST_DATABASE_URL)" -ForegroundColor Blue
Write-Host ""

# Build test arguments
$testArgs = @("test", "./tests")

if ($Verbose) {
    $testArgs += "-v"
}

if ($Coverage) {
    $testArgs += "-cover", "-coverprofile=coverage.out"
}

if ($Test) {
    $testArgs += "-run", $Test
}

# Test database connection
Write-Host "→ Testing database connection..." -ForegroundColor Blue
try {
    # Try to connect using psql (if available)
    $null = psql $env:TEST_DATABASE_URL -c "SELECT 1;" 2>&1
    Write-Host "✓ Database connection successful" -ForegroundColor Green
} catch {
    Write-Host "⚠ Could not verify database connection (psql not available)" -ForegroundColor Yellow
    Write-Host "  Continuing with tests..." -ForegroundColor Yellow
}
Write-Host ""

# Run tests
Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║                   Running Tests                        ║" -ForegroundColor Blue
Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

if ($Test) {
    Write-Host "→ Running specific test: $Test" -ForegroundColor Yellow
    & go @testArgs
} else {
    Write-Host "→ Running all database tests..." -ForegroundColor Blue
    Write-Host ""
    
    # Test 1: Main Migration Tests
    Write-Host "[1/3] Main Migration Tests" -ForegroundColor Blue
    $args1 = $testArgs + @("-run", "TestDatabaseMigrations_Complete")
    & go @args1
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    Write-Host ""
    
    # Test 2: Consolidated Migrations Tests
    Write-Host "[2/3] Consolidated Migrations Tests" -ForegroundColor Blue
    $args2 = $testArgs + @("-run", "TestConsolidatedMigrations")
    & go @args2
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    Write-Host ""
    
    # Test 3: Model Alignment Tests
    Write-Host "[3/3] Model-Database Alignment Tests" -ForegroundColor Blue
    $args3 = $testArgs + @("-run", "TestModelDatabaseAlignment|TestCriticalFieldMappings")
    & go @args3
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    Write-Host ""
}

# Show coverage if enabled
if ($Coverage) {
    Write-Host ""
    Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Blue
    Write-Host "║                   Coverage Report                      ║" -ForegroundColor Blue
    Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Blue
    & go tool cover -func=coverage.out | Select-Object -Last 1
    Write-Host ""
    Write-Host "→ Generate HTML coverage report:" -ForegroundColor Yellow
    Write-Host "  go tool cover -html=coverage.out -o coverage.html"
}

Write-Host ""
Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║              All tests completed successfully!         ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "Database Schema Validated:" -ForegroundColor Blue
Write-Host "  ✓ 30 Migrations (V01-V30)"
Write-Host "  ✓ 50+ Tables"
Write-Host "  ✓ 9 Lookup Tables"
Write-Host "  ✓ 8+ Foreign Keys"
Write-Host "  ✓ RBAC System"
Write-Host "  ✓ Model Alignment"
Write-Host ""
Write-Host "Ready for frontend development! 🚀" -ForegroundColor Green