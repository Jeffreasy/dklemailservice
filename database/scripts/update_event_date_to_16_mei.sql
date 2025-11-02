-- ============================================
-- UPDATE EVENT DATUM: 17 mei 2026 → 16 mei 2026
-- ============================================
-- 
-- Dit script update de event datum in de title_section_content tabel
-- van 17 mei 2026 naar de correcte datum: 16 mei 2026
--
-- INSTRUCTIES VOOR PGADMIN:
-- 1. Open pgAdmin
-- 2. Connect met de database
-- 3. Klik op Tools → Query Tool
-- 4. Kopieer en plak deze volledige query
-- 5. Klik op Execute (F5)
--
-- ============================================

-- STAP 1: Check huidige staat (voor referentie)
DO $$
DECLARE
    current_date_title TEXT;
    record_count INTEGER;
BEGIN
    SELECT detail_1_title INTO current_date_title
    FROM title_section_content
    WHERE id = '550e8400-e29b-41d4-a716-446655440001';
    
    SELECT COUNT(*) INTO record_count
    FROM title_section_content
    WHERE detail_1_title LIKE '%17 mei 2026%';
    
    RAISE NOTICE '═══════════════════════════════════════════════════';
    RAISE NOTICE 'HUIDIGE STAAT VOOR UPDATE:';
    RAISE NOTICE '═══════════════════════════════════════════════════';
    RAISE NOTICE 'Huidige detail_1_title: %', current_date_title;
    RAISE NOTICE 'Aantal records met "17 mei 2026": %', record_count;
    RAISE NOTICE '';
END $$;

-- STAP 2: Update de datum van 17 mei naar 16 mei
UPDATE title_section_content
SET 
    detail_1_title = '16 mei 2026',
    detail_1_description = REPLACE(detail_1_description, '17 mei 2026', '16 mei 2026'),
    updated_at = CURRENT_TIMESTAMP
WHERE 
    detail_1_title LIKE '%17 mei 2026%'
    OR detail_1_description LIKE '%17 mei 2026%';

-- STAP 3: Update ook eventuele afbeelding alt-text
UPDATE title_section_content
SET 
    image_alt = REPLACE(image_alt, '17 mei', '16 mei'),
    updated_at = CURRENT_TIMESTAMP
WHERE 
    image_alt LIKE '%17 mei%';

-- STAP 4: Verificatie na update
DO $$
DECLARE
    new_date_title TEXT;
    old_date_count INTEGER;
    new_date_count INTEGER;
    rec RECORD;
BEGIN
    SELECT detail_1_title INTO new_date_title
    FROM title_section_content
    WHERE id = '550e8400-e29b-41d4-a716-446655440001';
    
    SELECT COUNT(*) INTO old_date_count
    FROM title_section_content
    WHERE detail_1_title LIKE '%17 mei 2026%'
       OR detail_1_description LIKE '%17 mei 2026%';
    
    SELECT COUNT(*) INTO new_date_count
    FROM title_section_content
    WHERE detail_1_title LIKE '%16 mei 2026%'
       OR detail_1_description LIKE '%16 mei 2026%';
    
    RAISE NOTICE '';
    RAISE NOTICE '═══════════════════════════════════════════════════';
    RAISE NOTICE 'UPDATE RESULTAAT:';
    RAISE NOTICE '═══════════════════════════════════════════════════';
    RAISE NOTICE 'Nieuwe detail_1_title: %', new_date_title;
    RAISE NOTICE 'Records met "17 mei 2026": % (moet 0 zijn)', old_date_count;
    RAISE NOTICE 'Records met "16 mei 2026": %', new_date_count;
    RAISE NOTICE '';
    
    IF old_date_count > 0 THEN
        RAISE NOTICE '⚠️  WAARSCHUWING: Er zijn nog % records met "17 mei 2026"!', old_date_count;
    ELSE
        RAISE NOTICE '✅ SUCCESS: Alle datums zijn correct bijgewerkt naar 16 mei 2026';
    END IF;
    RAISE NOTICE '';
    
    -- Toon alle relevante records
    RAISE NOTICE '═══════════════════════════════════════════════════';
    RAISE NOTICE 'OVERZICHT VAN ALLE EVENT RECORDS:';
    RAISE NOTICE '═══════════════════════════════════════════════════';
    
    FOR rec IN 
        SELECT 
            id,
            event_title,
            detail_1_title,
            LEFT(detail_1_description, 50) || '...' as description_preview,
            updated_at
        FROM title_section_content
        ORDER BY updated_at DESC
    LOOP
        RAISE NOTICE 'ID: %', rec.id;
        RAISE NOTICE 'Title: %', rec.event_title;
        RAISE NOTICE 'Date: %', rec.detail_1_title;
        RAISE NOTICE 'Description: %', rec.description_preview;
        RAISE NOTICE 'Updated: %', rec.updated_at;
        RAISE NOTICE '---';
    END LOOP;
    
    RAISE NOTICE '═══════════════════════════════════════════════════';
END $$;

-- STAP 5: Commit de transactie
-- (In pgAdmin wordt dit automatisch gedaan, maar je kunt het handmatig forceren)
COMMIT;