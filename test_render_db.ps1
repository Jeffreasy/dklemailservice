# Test Render Database - V1_47 Migration Verification
$env:PGPASSWORD = "I4QP3JwyCcEbn8tGl6k3ErEvjUZ9V5rB"
$psql = "C:\Program Files\PostgreSQL\12\bin\psql.exe"

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "RENDER DATABASE V1_47 VERIFICATION" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Check V1_47 migration
Write-Host "1. Checking V1_47 migration status..." -ForegroundColor Yellow
& $psql -h "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com" -U "dekoninklijkeloopdatabase_user" -d "dekoninklijkeloopdatabase" -c "SELECT versie, naam, toegepast FROM migraties WHERE versie = '1.47.0';"

Write-Host ""

# Test 2: Count indexes
Write-Host "2. Counting indexes per table..." -ForegroundColor Yellow
& $psql -h "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com" -U "dekoninklijkeloopdatabase_user" -d "dekoninklijkeloopdatabase" -c "SELECT tablename, COUNT(*) as index_count FROM pg_indexes WHERE schemaname = 'public' GROUP BY tablename HAVING COUNT(*) > 2 ORDER BY index_count DESC LIMIT 10;"

Write-Host ""

# Test 3: Check specific new indexes
Write-Host "3. Checking new V1_47 indexes..." -ForegroundColor Yellow
& $psql -h "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com" -U "dekoninklijkeloopdatabase_user" -d "dekoninklijkeloopdatabase" -c "SELECT COUNT(*) as v1_47_indexes FROM pg_indexes WHERE schemaname = 'public' AND indexname LIKE 'idx_%' AND indexname IN ('idx_gebruikers_role_id', 'idx_verzonden_emails_contact_id', 'idx_contact_formulieren_fts');"

Write-Host ""

# Test 4: Table sizes
Write-Host "4. Top 5 largest tables..." -ForegroundColor Yellow
& $psql -h "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com" -U "dekoninklijkeloopdatabase_user" -d "dekoninklijkeloopdatabase" -c "SELECT tablename, pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size('public.'||tablename) DESC LIMIT 5;"

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "VERIFICATION COMPLETE" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "If V1_47 appears above, run ANALYZE:" -ForegroundColor Yellow
Write-Host '& $psql -h "dpg-cva4c01c1ekc738q6q0g-a.oregon-postgres.render.com" -U "dekoninklijkeloopdatabase_user" -d "dekoninklijkeloopdatabase" -c "ANALYZE;"' -ForegroundColor Cyan
Write-Host ""