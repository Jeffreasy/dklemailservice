# Test Albums API - Production
# Dit script test het ophalen van albums in production

$ErrorActionPreference = "Continue"
$ProductionUrl = "https://dklemailservice.onrender.com"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Testing Albums in PRODUCTION" -ForegroundColor Cyan
Write-Host "Base URL: $ProductionUrl" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Test 1: Health Check
Write-Host "1. Health Check..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$ProductionUrl/api/health" -Method Get -ErrorAction Stop
    Write-Host "SUCCESS: Health check successful" -ForegroundColor Green
    Write-Host "  Status: $($health.status)" -ForegroundColor Gray
    Write-Host "  Database: $($health.database)" -ForegroundColor Gray
} catch {
    Write-Host "FAILED: Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Could not connect to production server" -ForegroundColor Red
    exit 1
}

# Test 2: Get All Visible Albums
Write-Host "`n2. Get All Visible Albums (Public)..." -ForegroundColor Yellow
$albums = $null
try {
    $albums = Invoke-RestMethod -Uri "$ProductionUrl/api/albums" -Method Get -ErrorAction Stop
    Write-Host "SUCCESS: Fetched albums" -ForegroundColor Green
    Write-Host "  Total albums: $($albums.Count)" -ForegroundColor Gray
    
    if ($albums.Count -gt 0) {
        Write-Host "`n  Albums in Production:" -ForegroundColor Cyan
        foreach ($album in $albums) {
            Write-Host "    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
            Write-Host "    Album: $($album.title)" -ForegroundColor White
            Write-Host "       ID: $($album.id)" -ForegroundColor Gray
            Write-Host "       Description: $($album.description)" -ForegroundColor Gray
            Write-Host "       Visible: $($album.visible)" -ForegroundColor Gray
            Write-Host "       Order: $($album.order_number)" -ForegroundColor Gray
        }
        Write-Host "    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor DarkGray
    } else {
        Write-Host "  No albums found" -ForegroundColor Yellow
    }
} catch {
    Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails) {
        Write-Host "  Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

# Test 3: Get Albums with Cover Photos
Write-Host "3. Get Albums with Cover Photos..." -ForegroundColor Yellow
try {
    $albumsWithCovers = Invoke-RestMethod -Uri "$ProductionUrl/api/albums?include_covers=true" -Method Get -ErrorAction Stop
    Write-Host "SUCCESS: Fetched albums with covers" -ForegroundColor Green
    Write-Host "  Total albums: $($albumsWithCovers.Count)" -ForegroundColor Gray
    
    if ($albumsWithCovers.Count -gt 0) {
        Write-Host "`n  Cover Photos:" -ForegroundColor Cyan
        foreach ($album in $albumsWithCovers) {
            if ($album.cover_photo) {
                Write-Host "    Album '$($album.title)':" -ForegroundColor White
                Write-Host "       Photo ID: $($album.cover_photo.id)" -ForegroundColor Gray
                Write-Host "       URL: $($album.cover_photo.url)" -ForegroundColor Blue
                Write-Host ""
            }
        }
    }
} catch {
    Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Get Photos for Each Album
if ($albums -and $albums.Count -gt 0) {
    Write-Host "`n4. Get Photos for Albums..." -ForegroundColor Yellow
    
    $totalPhotos = 0
    foreach ($album in $albums) {
        Write-Host "  Album: $($album.title)" -ForegroundColor Cyan
        try {
            $photos = Invoke-RestMethod -Uri "$ProductionUrl/api/albums/$($album.id)/photos" -Method Get -ErrorAction Stop
            Write-Host "    Found $($photos.Count) photos" -ForegroundColor Green
            $totalPhotos += $photos.Count
            
            if ($photos.Count -gt 0 -and $photos.Count -le 3) {
                foreach ($photo in $photos) {
                    Write-Host "      - $($photo.title)" -ForegroundColor Gray
                }
            } elseif ($photos.Count -gt 3) {
                foreach ($photo in ($photos | Select-Object -First 3)) {
                    Write-Host "      - $($photo.title)" -ForegroundColor Gray
                }
                Write-Host "      ... and $($photos.Count - 3) more" -ForegroundColor DarkGray
            }
        } catch {
            Write-Host "    FAILED: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host "`n  Total photos across all albums: $totalPhotos" -ForegroundColor Cyan
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "PRODUCTION TEST SUMMARY" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Albums found: $($albums.Count)" -ForegroundColor Green
Write-Host "Production server is accessible" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Test completed!" -ForegroundColor Green