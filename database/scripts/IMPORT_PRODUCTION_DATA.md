# 📥 Productie Data Import & RBAC Sync Guide

> **Versie:** 1.49  
> **Datum:** 2025-11-02

---

## 🎯 Doel

Importeer productie gebruikers data naar Docker en wijs automatisch de correcte RBAC roles toe.

---

## 📋 Stap-voor-Stap Instructies

### Stap 1: Export Productie Gebruikers Data

**Via pgAdmin of productie database:**

```sql
-- Export gebruikers table naar CSV
\copy gebruikers TO 'productie_gebruikers.csv' CSV HEADER;

-- Of via psql:
psql -h productie-db-host -U user -d dbname -c "\copy gebruikers TO 'productie_gebruikers.csv' CSV HEADER"
```

**Of als SQL INSERT statements:**

```sql
-- Export als SQL statements
pg_dump -h productie-db-host -U user -d dbname \
    --table=gebruikers \
    --data-only \
    --column-inserts \
    > productie_gebruikers.sql
```

### Stap 2: Importeer naar Docker Database

**Optie A: Via SQL File**

```bash
# Kopieer SQL file naar container
docker cp productie_gebruikers.sql dkl-postgres:/tmp/

# Importeer
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -f /tmp/productie_gebruikers.sql
```

**Optie B: Via CSV File**

```bash
# Kopieer CSV file
docker cp productie_gebruikers.csv dkl-postgres:/tmp/

# Importeer in database  
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -c "\copy gebruikers FROM '/tmp/productie_gebruikers.csv' CSV HEADER"
```

**Optie C: Direct INSERT (Snelst voor weinig users)**

Gebruik de data die je hebt getoond en maak INSERT statements:

```sql
-- Voer dit uit in Docker database
docker exec -i dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase << 'EOF'

-- Jeffrey
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed)
VALUES (
    '748320dd-5b5e-4434-ad7e-8a405fd6266f',
    'Jeffrey',
    'jeffrey@dekoninklijkeloop.nl',
    '$2a$10$SxGtzqId5ZwGhvhmU4ys0O4tzEhsi3HYljq3ObRaEjppjohtta.2a',
    'staff',
    true,
    '2025-11-01 06:47:34.468919',
    '2025-03-14 19:40:14.755163',
    '2025-11-01 06:47:34.472282',
    false
) ON CONFLICT (id) DO UPDATE SET
    naam = EXCLUDED.naam,
    wachtwoord_hash = EXCLUDED.wachtwoord_hash,
    rol = EXCLUDED.rol,
    is_actief = EXCLUDED.is_actief,
    laatste_login = EXCLUDED.laatste_login,
    updated_at = NOW();

-- Lida
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed)
VALUES (
    '0197cfc3-7ca2-403b-ae4d-32627cd47222',
    'Lida',
    'Lida@dekoninklijkeloop.nl',
    '$2a$10$mrJh4LXjKO6/bcC/ZPJdQ.TQFGLUhiQU4XPJvbjebmnJ/eu5KgjZm',
    'staff',
    true,
    '2025-11-02 14:52:04.124096',
    '2025-11-02 12:54:43.931013',
    '2025-11-02 14:52:04.123658',
    false
) ON CONFLICT (id) DO UPDATE SET
    naam = EXCLUDED.naam,
    wachtwoord_hash = EXCLUDED.wachtwoord_hash,
    rol = EXCLUDED.rol,
    is_actief = EXCLUDED.is_actief,
    laatste_login = EXCLUDED.laatste_login,
    updated_at = NOW();

-- Marieke
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed)
VALUES (
    '11a3ec93-b159-473b-86d2-d3979b9c9e3a',
    'Marieke',
    'Marieke@dekoninklijkeloop.nl',
    '$2a$10$AJUC49EoFrusB9GOGP1aHucLItd5OZfKMWwgF5dsNwKRvF/lPcxse',
    'staff',
    true,
    '2025-10-05 17:42:17.59877',
    '2024-12-27 15:40:16.942007',
    '2025-10-05 17:42:17.599415',
    false
) ON CONFLICT (id) DO UPDATE SET
    naam = EXCLUDED.naam,
    wachtwoord_hash = EXCLUDED.wachtwoord_hash,
    rol = EXCLUDED.rol,
    is_actief = EXCLUDED.is_actief,
    laatste_login = EXCLUDED.laatste_login,
    updated_at = NOW();

-- Salih
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed)
VALUES (
    'b1bec8cb-1709-420f-88cb-fc74e0c6eec2',
    'Salih',
    'Salih@dekoninklijkeloop.nl',
    '$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku',
    'staff',
    true,
    NULL,
    '2024-12-28 02:43:22',
    '2025-10-05 17:15:03.315274',
    false
) ON CONFLICT (id) DO UPDATE SET
    naam = EXCLUDED.naam,
    wachtwoord_hash = EXCLUDED.wachtwoord_hash,
    rol = EXCLUDED.rol,
    is_actief = EXCLUDED.is_actief,
    laatste_login = EXCLUDED.laatste_login,
    updated_at = NOW();

-- Event deelnemers (sample - voeg meer toe zoals nodig)
-- Joyce Thielen (begeleider)
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, created_at, updated_at, newsletter_subscribed)
VALUES (
    '00a26c26-4932-4091-91ba-4385a251e285',
    'Joyce Thielen',
    'Joyce.thielen@sheerenloo.nl',
    '$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku',
    'begeleider',
    true,
    '2025-10-25 11:58:19.59437',
    '2025-10-25 12:09:50.868231',
    false
) ON CONFLICT (id) DO NOTHING;

-- JeffTest (socialmedia)
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed)
VALUES (
    '8f333073-8d5a-4202-9627-02863000b822',
    'JeffTest',
    'laventejeffrey@gmail.com',
    '$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku',
    'socialmedia',
    true,
    '2025-11-02 12:58:24.606474',
    '2025-10-25 11:58:19.59437',
    '2025-11-02 12:58:24.609733',
    false
) ON CONFLICT (id) DO NOTHING;

-- Voeg hier meer INSERT statements toe voor andere gebruikers...

COMMIT;
EOF
```

### Stap 3: Assign RBAC Roles

**Na import, voer sync script uit:**

```bash
docker exec -i dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase < database/scripts/sync_production_users_with_rbac.sql
```

Dit script wijst automatisch toe:
- ✅ `staff` role aan alle @dekoninklijkeloop.nl users
- ✅ `begeleider` role aan users met legacy rol 'begeleider'
- ✅ `deelnemer` role aan users met legacy rol 'deelnemer'
- ✅ `user` role aan overige users

### Stap 4: Verificeer Resultaten

```bash
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -c "
SELECT 
    g.email,
    g.rol as legacy_rol,
    STRING_AGG(DISTINCT r.name, ', ') as rbac_roles,
    COUNT(DISTINCT p.id) as permissions
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
GROUP BY g.email, g.rol
ORDER BY 
    CASE 
        WHEN g.email LIKE '%@dekoninklijkeloop.nl' THEN 1
        ELSE 2
    END,
    g.email;
"
```

**Verwachte Output:**

```
Email                          | Legacy Rol  | RBAC Roles  | Permissions
-------------------------------+-------------+-------------+-------------
admin@dekoninklijkeloop.nl     | admin       | admin       | 68
jeffrey@dekoninklijkeloop.nl   | staff       | staff       | 25
Lida@dekoninklijkeloop.nl      | staff       | staff       | 25
Marieke@dekoninklijkeloop.nl   | staff       | staff       | 25
Salih@dekoninklijkeloop.nl     | staff       | staff       | 25
Joyce.thielen@sheerenloo.nl    | begeleider  | begeleider  | 3
laventejeffrey@gmail.com       | socialmedia | user        | 2
... (event deelnemers)         | deelnemer   | deelnemer   | 3
```

---

## 🚀 Quick Import (All-in-One)

Als je alle gebruiker data hebt, maak een `import_users.sql` file en voer uit:

```bash
# 1. Maak import_users.sql met alle INSERT statements
# 2. Voer import uit
docker exec -i dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase < import_users.sql

# 3. Assign RBAC roles
docker exec -i dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase < database/scripts/sync_production_users_with_rbac.sql

# 4. Verify
docker exec dkl-postgres psql -U dekoninklijkeloopdatabase_user -d dekoninklijkeloopdatabase -c "SELECT COUNT(*) as total_users, COUNT(DISTINCT ur.user_id) as with_rbac FROM gebruikers g LEFT JOIN user_roles ur ON g.id = ur.user_id WHERE g.is_actief = true;"
```

---

## ✅ Expected Results Na Import

### @dekoninklijkeloop.nl Users
- admin@dekoninklijkeloop.nl → `admin` role (68 permissions)
- jeffrey@dekoninklijkeloop.nl → `staff` role (25 permissions)
- Lida@dekoninklijkeloop.nl → `staff` role (25 permissions)
- Marieke@dekoninklijkeloop.nl → `staff` role (25 permissions)
- Salih@dekoninklijkeloop.nl → `staff` role (25 permissions)

### Event Participants
- Users met rol 'begeleider' → `begeleider` role (3 steps permissions)
- Users met rol 'deelnemer' → `deelnemer` role (3 steps permissions)
- Users met rol 'vrijwilliger' → `vrijwilliger` role (3 steps permissions)

### Other Users
- Users met rol 'socialmedia' of geen rol → `user` role (2 chat permissions)

---

## 🔒 Security Note

**BELANGRIJK:** Wachtwoord hashes uit productie zijn veilig om te importeren, maar:
- ⚠️ Gebruik alleen voor development/testing
- ⚠️ Voor nieuwe productie deployment: reset alle wachtwoorden
- ⚠️ Informeer users over wachtwoord reset procedure

---

**Last Updated:** 2025-11-02  
**Compatibility:** PostgreSQL 12+  
**Docker Container:** dkl-postgres