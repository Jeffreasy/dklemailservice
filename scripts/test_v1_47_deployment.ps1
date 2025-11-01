# Test V1_47 Deployment on Render
# This script tests if the V1_47 migration was successfully deployed

param(
    [string]$ApiUrl = "https://your-render-app-url.onrender.com"  # Update with your Render URL
)

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "V1_47 DEPLOYMENT VERIFICATION TEST" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health Check
Write-Host "1. Testing application health..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$ApiUrl/api/health" -Method Get -TimeoutSec 10
    Write-Host "✅ Application is running!" -ForegroundColor Green
    Write-Host "   Status: $($healthResponse.status)" -ForegroundColor Gray
} catch {
    Write-Host "❌ Application health check failed!" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Possible reasons:" -ForegroundColor Yellow
    Write-Host "- Render deployment still in progress" -ForegroundColor Gray
    Write-Host "- Application failed to start" -ForegroundColor Gray
    Write-Host "- Migration error occurred" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "1. Check Render Dashboard: https://dashboard.render.com/" -ForegroundColor Gray
    Write-Host "2. View deployment logs for errors" -ForegroundColor Gray
    Write-Host "3. Look for V1_47 migration in logs" -ForegroundColor Gray
    exit 1
}

Write-Host ""

# Test 2: Database Connection (via health endpoint)
Write-Host "2. Testing database connectivity..." -ForegroundColor Yellow
if ($healthResponse.database -eq "connected" -or $healthResponse.status -eq "healthy") {
    Write-Host "✅ Database is connected!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Database status unknown from health endpoint" -ForegroundColor Yellow
}

Write-Host ""

# Test 3: Check if application responds (indicates successful start)
Write-Host "3. Verifying application started successfully..." -ForegroundColor Yellow
try {
    $metricsResponse = Invoke-RestMethod -Uri "$ApiUrl/api/metrics" -Method Get -TimeoutSec 10 -ErrorAction SilentlyContinue
    Write-Host "✅ Metrics endpoint responding!" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Metrics endpoint not accessible (might be protected)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "DEPLOYMENT CHECK RESULTS" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "✅ Application is running and responding" -ForegroundColor Green
Write-Host ""
Write-Host "To verify V1_47 migration:" -ForegroundColor Yellow
Write-Host "1. Go to Render Dashboard: https://dashboard.render.com/" -ForegroundColor Gray
Write-Host "2. Click on your DKL Email Service" -ForegroundColor Gray
Write-Host "3. Go to 'Logs' tab" -ForegroundColor Gray
Write-Host "4. Search for:" -ForegroundColor Gray
Write-Host '   "Migratie succesvol uitgevoerd","file":"V1_47__performance_optimizations.sql"' -ForegroundColor Cyan
Write-Host ""
Write-Host "Then run ANALYZE to update query planner:" -ForegroundColor Yellow
Write-Host "   psql `$DATABASE_URL -c 'ANALYZE;'" -ForegroundColor Cyan
Write-Host ""
Write-Host "Or use the SQL test script:" -ForegroundColor Yellow
Write-Host "   psql `$DATABASE_URL -f database/scripts/test_v1_47_indexes.sql" -ForegroundColor Cyan
Write-Host ""