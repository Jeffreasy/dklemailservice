-- Fix Data Quality Issues Before V1_48 Deployment
-- This script fixes invalid data that would violate V1_48 constraints

\echo 'Starting data quality fixes for V1_48...'
\echo ''

-- ============================================
-- SECTION 1: FIX INVALID EMAILS
-- ============================================

\echo 'Fixing invalid emails in gebruikers...'

-- Show invalid emails first
SELECT 
    id, 
    naam, 
    email as old_email,
    'fixed_' || REPLACE(email, '@', '_at_') || '@placeholder.invalid' as new_email
FROM gebruikers 
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
LIMIT 5;

-- Fix invalid emails in gebruikers
UPDATE gebruikers 
SET email = 'fixed_' || REPLACE(email, '@', '_at_') || '@placeholder.invalid'
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';

\echo 'Fixed emails in gebruikers.'
\echo ''

\echo 'Fixing invalid emails in contact_formulieren...'

-- Show invalid emails
SELECT 
    id, 
    naam, 
    email as old_email,
    'fixed_' || REPLACE(email, '@', '_at_') || '@placeholder.invalid' as new_email
FROM contact_formulieren 
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';

-- Fix invalid emails in contact_formulieren
UPDATE contact_formulieren 
SET email = 'fixed_' || REPLACE(email, '@', '_at_') || '@placeholder.invalid'
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';

\echo 'Fixed emails in contact_formulieren.'
\echo ''

\echo 'Fixing invalid emails in aanmeldingen (if any)...'

-- Fix invalid emails in aanmeldingen
UPDATE aanmeldingen 
SET email = 'fixed_' || REPLACE(email, '@', '_at_') || '@placeholder.invalid'
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';

-- ============================================
-- SECTION 2: FIX EMPTY NAMES
-- ============================================

\echo 'Fixing empty names...'

-- Fix empty names in gebruikers
UPDATE gebruikers 
SET naam = 'Onbekend_' || id::text
WHERE LENGTH(TRIM(naam)) = 0;

-- Fix empty names in contact_formulieren
UPDATE contact_formulieren 
SET naam = 'Onbekend_' || id::text
WHERE LENGTH(TRIM(naam)) = 0;

-- Fix empty names in aanmeldingen
UPDATE aanmeldingen 
SET naam = 'Onbekend_' || id::text
WHERE LENGTH(TRIM(naam)) = 0;

-- ============================================
-- SECTION 3: FIX NEGATIVE STEPS
-- ============================================

\echo 'Fixing negative steps (if any)...'

UPDATE aanmeldingen 
SET steps = 0
WHERE steps < 0;

-- ============================================
-- SECTION 4: FIX INVALID STATUS VALUES
-- ============================================

\echo 'Fixing invalid status values...'

-- Check and fix contact_formulieren status
UPDATE contact_formulieren 
SET status = 'nieuw'
WHERE status NOT IN ('nieuw', 'in_behandeling', 'beantwoord', 'gesloten');

-- Check and fix aanmeldingen status
UPDATE aanmeldingen 
SET status = 'nieuw'
WHERE status NOT IN ('nieuw', 'bevestigd', 'geannuleerd', 'voltooid');

-- ============================================
-- SECTION 5: FIX EMAIL CONSISTENCY
-- ============================================

\echo 'Fixing email consistency...'

-- Contact formulieren: If email_verzonden = TRUE but no timestamp
UPDATE contact_formulieren
SET email_verzonden_op = created_at
WHERE email_verzonden = TRUE AND email_verzonden_op IS NULL;

-- Aanmeldingen: If email_verzonden = TRUE but no timestamp
UPDATE aanmeldingen
SET email_verzonden_op = created_at
WHERE email_verzonden = TRUE AND email_verzonden_op IS NULL;

-- ============================================
-- SECTION 6: VERIFICATION
-- ============================================

\echo ''
\echo '=== VERIFICATION AFTER FIXES ==='

-- Verify no more invalid emails
SELECT 
    'Invalid emails in gebruikers' as check,
    COUNT(*) as count
FROM gebruikers 
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
UNION ALL
SELECT 
    'Invalid emails in contact_formulieren',
    COUNT(*)
FROM contact_formulieren 
WHERE email !~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
UNION ALL
SELECT 
    'Empty names',
    COUNT(*)
FROM gebruikers 
WHERE LENGTH(TRIM(naam)) = 0
UNION ALL
SELECT 
    'Negative steps',
    COUNT(*)
FROM aanmeldingen 
WHERE steps < 0;

\echo ''
\echo 'Data quality fixes completed!'
\echo 'All counts above should be 0.'
\echo ''
\echo 'Ready for V1_48 deployment!'