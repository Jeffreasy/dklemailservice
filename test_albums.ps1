# Test Albums API - Docker & Production
# Dit script test het ophalen van albums in zowel Docker als production

param(
    [string]$Environment = "docker"
)

$ErrorActionPreference = "Continue"

function Test-AlbumsEndpoint {
    param(
        [string]$BaseUrl,
        [string]$EnvName
    )
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Testing Albums in $EnvName" -ForegroundColor Cyan
    Write-Host "Base URL: $BaseUrl" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    
    # Test 1: Health Check
    Write-Host "1. Health Check..." -ForegroundColor Yellow
    try {
        $health = Invoke-RestMethod -Uri "$BaseUrl/api/health" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Health check successful" -ForegroundColor Green
        Write-Host "  Status: $($health.status)" -ForegroundColor Gray
    } catch {
        Write-Host "FAILED: Health check failed: $($_.Exception.Message)" -ForegroundColor Red
        return
    }
    
    # Test 2: Get All Visible Albums
    Write-Host "`n2. Get All Visible Albums (Public)..." -ForegroundColor Yellow
    try {
        $albums = Invoke-RestMethod -Uri "$BaseUrl/api/albums" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Fetched albums" -ForegroundColor Green
        Write-Host "  Total albums: $($albums.Count)" -ForegroundColor Gray
        
        if ($albums.Count -gt 0) {
            Write-Host "`n  Albums:" -ForegroundColor Cyan
            foreach ($album in $albums) {
                Write-Host "    - ID: $($album.id)" -ForegroundColor White
                Write-Host "      Title: $($album.title)" -ForegroundColor White
                Write-Host "      Description: $($album.description)" -ForegroundColor Gray
                Write-Host "      Visible: $($album.visible)" -ForegroundColor Gray
                Write-Host "      Order: $($album.order_number)" -ForegroundColor Gray
                Write-Host ""
            }
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
    Write-Host "`n3. Get Albums with Cover Photos..." -ForegroundColor Yellow
    try {
        $albumsWithCovers = Invoke-RestMethod -Uri "$BaseUrl/api/albums?include_covers=true" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Fetched albums with covers" -ForegroundColor Green
        Write-Host "  Total albums: $($albumsWithCovers.Count)" -ForegroundColor Gray
        
        if ($albumsWithCovers.Count -gt 0) {
            foreach ($album in $albumsWithCovers) {
                if ($album.cover_photo) {
                    Write-Host "  Album '$($album.title)' has cover photo:" -ForegroundColor Cyan
                    Write-Host "    Photo ID: $($album.cover_photo.id)" -ForegroundColor Gray
                    Write-Host "    URL: $($album.cover_photo.url)" -ForegroundColor Gray
                }
            }
        }
    } catch {
        Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # Test 4: Get Photos for First Album
    if ($albums -and $albums.Count -gt 0) {
        $firstAlbumId = $albums[0].id
        Write-Host "`n4. Get Photos for Album '$($albums[0].title)'..." -ForegroundColor Yellow
        try {
            $photos = Invoke-RestMethod -Uri "$BaseUrl/api/albums/$firstAlbumId/photos" -Method Get -ErrorAction Stop
            Write-Host "SUCCESS: Fetched album photos" -ForegroundColor Green
            Write-Host "  Total photos: $($photos.Count)" -ForegroundColor Gray
            
            if ($photos.Count -gt 0) {
                Write-Host "`n  Sample photos:" -ForegroundColor Cyan
                foreach ($photo in ($photos | Select-Object -First 5)) {
                    Write-Host "    - Photo ID: $($photo.id)" -ForegroundColor White
                    Write-Host "      Title: $($photo.title)" -ForegroundColor White
                    Write-Host "      URL: $($photo.url)" -ForegroundColor Gray
                    if ($photo.thumbnail_url) {
                        Write-Host "      Thumbnail: $($photo.thumbnail_url)" -ForegroundColor Gray
                    }
                    Write-Host ""
                }
                if ($photos.Count -gt 5) {
                    Write-Host "    ... and $($photos.Count - 5) more photos" -ForegroundColor Gray
                }
            }
        } catch {
            Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Completed testing $EnvName" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

# Main execution
switch ($Environment.ToLower()) {
    "docker" {
        Test-AlbumsEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
    }
    "production" {
        Write-Host "Enter production URL (e.g., https://api.dekoninklijkeloop.nl):" -ForegroundColor Yellow
        $prodUrl = Read-Host
        if ($prodUrl) {
            Test-AlbumsEndpoint -BaseUrl $prodUrl -EnvName "Production"
        } else {
            Write-Host "No production URL provided" -ForegroundColor Red
        }
    }
    "both" {
        Test-AlbumsEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
        
        Write-Host "`nEnter production URL (e.g., https://api.dekoninklijkeloop.nl):" -ForegroundColor Yellow
        $prodUrl = Read-Host
        if ($prodUrl) {
            Test-AlbumsEndpoint -BaseUrl $prodUrl -EnvName "Production"
        }
    }
    default {
        Write-Host "Invalid environment. Use: docker, production, or both" -ForegroundColor Red
        Write-Host "Example: .\test_albums.ps1 -Environment docker" -ForegroundColor Yellow
    }
}

Write-Host "`nTest completed!" -ForegroundColor Green