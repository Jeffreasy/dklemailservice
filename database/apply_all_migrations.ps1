# Apply All Migrations Script
$ErrorActionPreference = "Stop"

Write-Host "=========================================================" -ForegroundColor Blue
Write-Host "      Applying All 30 Migrations to Test Database       " -ForegroundColor Blue
Write-Host "=========================================================" -ForegroundColor Blue
Write-Host ""

$DATABASE = "dklemailservice_test"
$CONTAINER = "dkl-postgres"

$migrations = @(
    "V01__initial_schema.sql",
    "V02__seed_data.sql",
    "V03__add_test_data.sql",
    "V04__create_chat_tables.sql",
    "V05__add_newsletter_tables.sql",
    "V06__create_rbac_tables.sql",
    "V07__seed_rbac_tables.sql",
    "V08__migrate_and_assign_user_roles.sql",
    "V09__create_refresh_tokens_table.sql",
    "V10__create_uploaded_images_table.sql",
    "V11__migrate_cms_data.sql",
    "V12__add_gebruiker_id_to_aanmeldingen.sql",
    "V13__add_steps_permissions.sql",
    "V14__add_remaining_cms_permissions.sql",
    "V15__replace_title_sections.sql",
    "V16__add_steps_to_aanmeldingen.sql",
    "V17_CONSOLIDATED__create_route_funds.sql",
    "V18__performance_optimizations.sql",
    "V19__advanced_optimizations.sql",
    "V20__update_staff_aanmelding_permissions.sql",
    "V21__add_gamification_tables.sql",
    "V22__improve_user_participant_linking.sql",
    "V23__add_events_table.sql",
    "V24__create_notulen_module.sql",
    "V25__create_leaderboard_materialized_view.sql",
    "V26_CONSOLIDATED__normalize_roles_and_distances.sql",
    "V27_CONSOLIDATED__normalize_status_and_type_fields.sql",
    "V28_CONSOLIDATED__rename_tables_participant_refactor.sql",
    "V29_CONSOLIDATED__update_permissions_for_participants.sql",
    "V30__add_is_active_to_participant_roles.sql"
)

$migrationsDir = "migrations"
$success = 0
$failed = 0

foreach ($migration in $migrations) {
    $filePath = Join-Path $migrationsDir $migration
    
    if (Test-Path $filePath) {
        Write-Host "-> Applying: $migration" -ForegroundColor Cyan
        
        docker cp $filePath "${CONTAINER}:/tmp/migration.sql" 2>&1 | Out-Null
        $output = docker exec $CONTAINER psql -U postgres -d $DATABASE -f /tmp/migration.sql 2>&1
        $exitCode = $LASTEXITCODE
        
        if ($exitCode -eq 0) {
            Write-Host "   [OK]" -ForegroundColor Green
            $success++
        } else {
            Write-Host "   [FAIL]" -ForegroundColor Red
            if ($output) {
                Write-Host "   Error: $output" -ForegroundColor Red
            }
            $failed++
        }
    } else {
        Write-Host "   [NOT FOUND] $filePath" -ForegroundColor Yellow
        $failed++
    }
}

Write-Host ""
Write-Host "=========================================================" -ForegroundColor Blue
Write-Host "                  Migration Summary                      " -ForegroundColor Blue
Write-Host "=========================================================" -ForegroundColor Blue
Write-Host ""
Write-Host "  Total: $($migrations.Count)" -ForegroundColor White
Write-Host "  Success: $success" -ForegroundColor Green
Write-Host "  Failed: $failed" -ForegroundColor $(if ($failed -gt 0) { "Red" } else { "Green" })
Write-Host ""

if ($failed -eq 0) {
    Write-Host "[OK] All migrations applied!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "[FAIL] Some migrations failed" -ForegroundColor Red
    exit 1
}