# Test Videos API - Docker & Production
# Dit script test het ophalen van videos in zowel Docker als production

param(
    [string]$Environment = "docker"
)

$ErrorActionPreference = "Continue"

function Test-VideosEndpoint {
    param(
        [string]$BaseUrl,
        [string]$EnvName
    )
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Testing Videos in $EnvName" -ForegroundColor Cyan
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
    
    # Test 2: Get All Visible Videos
    Write-Host "`n2. Get All Visible Videos (Public)..." -ForegroundColor Yellow
    try {
        $videos = Invoke-RestMethod -Uri "$BaseUrl/api/videos" -Method Get -ErrorAction Stop
        Write-Host "SUCCESS: Fetched videos" -ForegroundColor Green
        Write-Host "  Total videos: $($videos.Count)" -ForegroundColor Gray
        
        if ($videos.Count -gt 0) {
            Write-Host "`n  Videos:" -ForegroundColor Cyan
            foreach ($video in $videos) {
                Write-Host "    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor DarkGray
                Write-Host "    Video: $($video.title)" -ForegroundColor White
                Write-Host "       ID: $($video.id)" -ForegroundColor Gray
                Write-Host "       Video ID: $($video.video_id)" -ForegroundColor Gray
                Write-Host "       URL: $($video.url)" -ForegroundColor Blue
                Write-Host "       Description: $($video.description)" -ForegroundColor Gray
                if ($video.thumbnail_url) {
                    Write-Host "       Thumbnail: $($video.thumbnail_url)" -ForegroundColor DarkBlue
                }
                Write-Host "       Visible: $($video.visible)" -ForegroundColor Gray
                Write-Host "       Order: $($video.order_number)" -ForegroundColor Gray
                Write-Host ""
            }
        } else {
            Write-Host "  No videos found" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.ErrorDetails) {
            Write-Host "  Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Completed testing $EnvName" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

# Main execution
switch ($Environment.ToLower()) {
    "docker" {
        Test-VideosEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
    }
    "production" {
        Test-VideosEndpoint -BaseUrl "https://dklemailservice.onrender.com" -EnvName "Production"
    }
    "both" {
        Test-VideosEndpoint -BaseUrl "http://localhost:8082" -EnvName "Docker (Dev)"
        Test-VideosEndpoint -BaseUrl "https://dklemailservice.onrender.com" -EnvName "Production"
    }
    default {
        Write-Host "Invalid environment. Use: docker, production, or both" -ForegroundColor Red
        Write-Host "Example: .\test_videos.ps1 -Environment production" -ForegroundColor Yellow
    }
}

Write-Host "`nTest completed!" -ForegroundColor Green