-- SQL Script om nieuwe admin gebruikers toe te voegen
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, created_at, updated_at)
VALUES 
(gen_random_uuid(), 'Jeffrey', 'jeffrey@dekoninklijkeloop.nl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', true, NOW(), NOW()),
(gen_random_uuid(), 'Salih', 'salih@dekoninklijkeloop.nl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', true, NOW(), NOW()),
(gen_random_uuid(), 'Marieke', 'marieke@dekoninklijkeloop.nl', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', true, NOW(), NOW()); 