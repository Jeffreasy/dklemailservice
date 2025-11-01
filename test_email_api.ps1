# Test Email API - Complete Test Suite
# Dit script test alle email-gerelateerde endpoints

param(
    [string]$Environment = "production",
    [string]$TestType = "public"
)

$ErrorActionPreference = "Continue"

function Get-BaseUrl {
    param([string]$Env)
    
    if ($Env -eq "docker") {
        return "http://localhost:8082"
    } else {
        return "https://dklemailservice.onrender.com"
    }
}

function Test-PublicEmailEndpoints {
    param(
        [string]$BaseUrl,
        [string]$EnvName
    )
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "PUBLIC EMAIL API TESTS - $EnvName" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    
    # Test 1: Contact Form Submission (Test Mode)
    Write-Host "1. Contact Form Submission (Test Mode)..." -ForegroundColor Yellow
    try {
        $contactData = @{
            naam = "Test Gebruiker"
            email = "test@example.com"
            telefoon = "06-12345678"
            bericht = "Dit is een test bericht vanuit het contact formulier"
            privacy_akkoord = $true
            test_mode = $true
        } | ConvertTo-Json
        
        $response = Invoke-RestMethod -Uri "$BaseUrl/api/contact-email" -Method Post -Body $contactData -ContentType "application/json" -ErrorAction Stop
        
        Write-Host "SUCCESS: Contact form submitted" -ForegroundColor Green
        Write-Host "  Message: $($response.message)" -ForegroundColor Gray
        if ($response.test_mode) {
            Write-Host "  Test Mode: Active" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # Test 2: Aanmelding Form Submission (Test Mode)
    Write-Host "`n2. Aanmelding Form Submission (Test Mode)..." -ForegroundColor Yellow
    try {
        $aanmeldingData = @{
            naam = "Test Deelnemer"
            email = "deelnemer@example.com"
            telefoon = "06-98765432"
            rol = "deelnemer"
            afstand = "10km"
            ondersteuning = "Geen"
            bijzonderheden = "Eerste keer deelnemer"
            terms = $true
            test_mode = $true
        } | ConvertTo-Json
        
        $response = Invoke-RestMethod -Uri "$BaseUrl/api/aanmelding-email" -Method Post -Body $aanmeldingData -ContentType "application/json" -ErrorAction Stop
        
        Write-Host "SUCCESS: Aanmelding submitted" -ForegroundColor Green
        Write-Host "  Message: $($response.message)" -ForegroundColor Gray
        if ($response.test_mode) {
            Write-Host "  Test Mode: Active" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # Test 3: Contact Form Validation Test
    Write-Host "`n3. Contact Form Validation (Missing Fields)..." -ForegroundColor Yellow
    try {
        $invalidData = @{
            naam = ""
            email = "test@example.com"
            bericht = ""
            privacy_akkoord = $false
        } | ConvertTo-Json
        
        $null = Invoke-RestMethod -Uri "$BaseUrl/api/contact-email" -Method Post -Body $invalidData -ContentType "application/json" -ErrorAction Stop
        
        Write-Host "UNEXPECTED: Should have failed validation" -ForegroundColor Yellow
    } catch {
        if ($_.Exception.Response.StatusCode -eq 'BadRequest') {
            Write-Host "SUCCESS: Validation working correctly" -ForegroundColor Green
            try {
                $errorResponse = $_.ErrorDetails.Message | ConvertFrom-Json
                Write-Host "  Error: $($errorResponse.error)" -ForegroundColor Gray
            } catch {
                Write-Host "  Validation error returned" -ForegroundColor Gray
            }
        } else {
            Write-Host "FAILED: Unexpected error: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    # Test 4: Aanmelding Validation Test
    Write-Host "`n4. Aanmelding Validation (Missing Terms)..." -ForegroundColor Yellow
    try {
        $invalidAanmelding = @{
            naam = "Test"
            email = "test@example.com"
            rol = "deelnemer"
            afstand = "10km"
            terms = $false
        } | ConvertTo-Json
        
        $null = Invoke-RestMethod -Uri "$BaseUrl/api/aanmelding-email" -Method Post -Body $invalidAanmelding -ContentType "application/json" -ErrorAction Stop
        
        Write-Host "UNEXPECTED: Should have failed validation" -ForegroundColor Yellow
    } catch {
        if ($_.Exception.Response.StatusCode -eq 'BadRequest') {
            Write-Host "SUCCESS: Terms validation working" -ForegroundColor Green
            try {
                $errorResponse = $_.ErrorDetails.Message | ConvertFrom-Json
                Write-Host "  Error: $($errorResponse.error)" -ForegroundColor Gray
            } catch {
                Write-Host "  Validation error returned" -ForegroundColor Gray
            }
        } else {
            Write-Host "FAILED: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
}

# Main execution
$baseUrl = Get-BaseUrl -Env $Environment

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "   DKL Email Service API Test Suite    " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Health Check
Write-Host "`nHealth Check..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/api/health" -Method Get -ErrorAction Stop
    Write-Host "SUCCESS: Service is $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "FAILED: Cannot connect to service" -ForegroundColor Red
    exit 1
}

# Test public endpoints
if ($TestType -eq "all" -or $TestType -eq "public") {
    Test-PublicEmailEndpoints -BaseUrl $baseUrl -EnvName $Environment
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "         Test Suite Completed           " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

Write-Host "`nTested Environment: $Environment" -ForegroundColor White
Write-Host "Base URL: $baseUrl" -ForegroundColor Blue
Write-Host "`nAll tests completed!`n" -ForegroundColor Green