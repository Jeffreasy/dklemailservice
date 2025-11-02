-- Check huidige datum in title_section_content tabel
-- Dit script controleert of de event datum correct is

SELECT 
    id,
    event_title,
    detail_1_title,
    detail_1_description,
    created_at,
    updated_at
FROM title_section_content
WHERE detail_1_title LIKE '%mei 2026%'
   OR detail_1_description LIKE '%mei 2026%'
   OR event_title LIKE '%2026%';

-- Expected: "16 mei 2026" (CORRECT)
-- If shows: "17 mei 2026" (INCORRECT) - dan moet je update_event_date.sql runnen