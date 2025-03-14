@echo off
REM Script om tests uit te voeren met CGO enabled voor SQLite ondersteuning

REM Zorg ervoor dat CGO enabled is
set CGO_ENABLED=1

REM Ga naar de project root directory
cd %~dp0\..

REM Voer alle tests uit
echo Running tests with CGO_ENABLED=1...
go test ./tests/... -v

REM Voer tests uit met coverage
if "%1"=="--coverage" (
    echo Running tests with coverage...
    go test ./tests/... -coverprofile=coverage -v
    go tool cover -html=coverage -o coverage.html
    echo Coverage report generated at coverage.html
) 