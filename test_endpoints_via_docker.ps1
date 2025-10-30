# Test alle API endpoints via Docker container
Write-Host "Testing DKL Email Service API Endpoints via Docker..." -ForegroundColor Green
Write-Host ""

$baseUrl = "http://localhost:8080"
$container = "dkl-email-service"

# Array van endpoints om te testen
$endpoints = @(
    "/",
    "/api/partners",
    "/api/radio-recordings",
    "/api/photos",
    "/api/albums",
    "/api/videos",
    "/api/sponsors",
    "/api/program-schedule",
    "/api/social-embeds",
    "/api/social-links",
    "/api/under-construction/active",
    "/api/title_section_content"
)

$successCount = 0
$failCount = 0

foreach ($endpoint in $endpoints) {
    Write-Host "Testing: $endpoint" -ForegroundColor Cyan
    
    # Test via wget inside container
    $cmd = "wget -q -O - '$baseUrl$endpoint'"
    docker exec $container /bin/sh -c $cmd > $null 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Success" -ForegroundColor Green
        $successCount++
    } else {
        Write-Host "  ✗ Failed (Exit: $LASTEXITCODE)" -ForegroundColor Red
        $failCount++
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Test Results Summary" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Total Tests: $($endpoints.Count)" -ForegroundColor White
Write-Host "Passed: $successCount" -ForegroundColor Green
Write-Host "Failed: $failCount" -ForegroundColor Red