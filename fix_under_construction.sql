-- Fix under construction records
-- Delete empty records (IDs 2-6)
DELETE FROM under_construction WHERE id IN (2, 3, 4, 5, 6);

-- Update record ID 1 to be active
UPDATE under_construction
SET is_active = true
WHERE id = 1;