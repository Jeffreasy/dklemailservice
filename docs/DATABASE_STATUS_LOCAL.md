# Lokale Database Status

Status van de lokale Docker PostgreSQL database voor development.

## 📊 Database Overzicht

**Database:** `dklemailservice`  
**Host:** `localhost:5433`  
**Gebruiker:** `postgres`  
**Container:** `dkl-postgres`

## 📈 Data Samenvatting

| Tabel | Aantal Records | Beschrijving |
|-------|----------------|--------------|
| **gebruikers** | 1 | Admin account voor testing |
| **contact_formulieren** | 2 | Test contactformulieren |
| **aanmeldingen** | 32 | Test aanmeldingen (deelnemers, begeleiders) |
| **albums** | 1 | Test album |
| **photos** | 45 | Test foto's |
| **videos** | 5 | Test video's |
| **sponsors** | 5 | Test sponsors |
| **roles** | 9 | Systeemrollen |
| **permissions** | 68 | Granulaire permissies |

## 👤 Gebruikers Account

### Admin Account (voor frontend login)

```
Email:    admin@dekoninklijkeloop.nl
Rol:      admin
Status:   actief
ID:       1cdff2e3-bcce-4d58-b956-fb7c2281cd00
```

**Wachtwoord:** Check de database migraties voor het gehashte wachtwoord, of maak een nieuw account.

## 📋 Test Data

### Contactformulieren (2)

| Naam | Email | Status |
|------|-------|--------|
| Bas heijenk | basheijenk96@gmail.com | nieuw |
| je geheime liefde | de.konining@willem.alexander.nl | nieuw |

### Aanmeldingen (32 totaal, eerste 5):

| Naam | Email | Rol | Afstand | Steps |
|------|-------|-----|---------|-------|
| TGTest | laventejeffrey@gmail.com | Begeleider | 15 KM | 0 |
| Manuela van Zwam | rik.van-harxen@sheerenloo.nl | Deelnemer | 2.5 KM | 0 |
| Bas heijenk | basheijenk96@gmail.com | Deelnemer | 2.5 KM | 0 |
| Salih | topraks@gmail.com | Deelnemer | 2.5 KM | 0 |
| Manuela van zwam | benjaminlaan.64a@sheerenloo.nl | Deelnemer | 2.5 KM | 0 |

## 🎭 Rollen & Permissies

### Beschikbare Rollen (9)

| Rol | Beschrijving |
|-----|--------------|
| **admin** | Volledige beheerder met toegang tot alle functies |
| **staff** | Ondersteunend personeel met beperkte beheerrechten |
| **user** | Standaard gebruiker |
| **owner** | Chat kanaal eigenaar |
| **chat_admin** | Chat kanaal beheerder |
| **member** | Chat kanaal lid |
| **deelnemer** | Evenement deelnemer |
| **begeleider** | Evenement begeleider |
| **vrijwilliger** | Evenement vrijwilliger |

### Permissie Systeem

**Totaal:** 68 granulaire permissies

**Categorieën:**
- Contact management (read, write, delete)
- Aanmelding management (read, write, delete)
- User management (read, write, delete, manage_roles)
- Newsletter (read, write, send, delete)
- Email management (read, write, delete, fetch)
- Chat (read, write, manage_channel, moderate)
- Album/Photo/Video management
- Steps tracking (read, write)
- Partner/Sponsor management
- Program schedule management

**Admin heeft toegang tot alle 68 permissies.**

## 📸 Content Data

### Albums
- 1 album beschikbaar
- 45 foto's in totaal
- Album-photo relaties in `album_photos` tabel

### Media
- 5 video's (Streamable embeds)
- 5 sponsors met logo's
- Social embeds en links beschikbaar

### Program
- Program schedule items beschikbaar
- Radio recordings beschikbaar
- Under construction pages beschikbaar

## 🆕 Nieuwe Features

### Steps Tracking (NIEUW!)

```sql
-- Steps veld toegevoegd aan aanmeldingen
SELECT naam, email, steps FROM aanmeldingen WHERE steps > 0;

-- Route funds tabel aangemaakt
SELECT * FROM route_funds;
```

Alle aanmeldingen hebben momenteel 0 steps - perfect voor testing!

## 🔧 Database Toegang

### Via Command Line

```bash
# Connect to database
docker exec -it dkl-postgres psql -U postgres -d dklemailservice

# Zonder -it (Windows)
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "YOUR SQL HERE"
```

### Via GUI Tools

**pgAdmin/DBeaver Configuratie:**
```
Host:     localhost
Port:     5433
Database: dklemailservice
Username: postgres
Password: postgres
SSL:      disable
```

## 📝 Handige Queries

### Check Gebruikers
```sql
SELECT id, naam, email, rol, is_actief FROM gebruikers;
```

### Recent Contactformulieren
```sql
SELECT naam, email, status, created_at 
FROM contact_formulieren 
ORDER BY created_at DESC 
LIMIT 10;
```

### Aanmeldingen per Rol
```sql
SELECT rol, COUNT(*) as aantal 
FROM aanmeldingen 
GROUP BY rol 
ORDER BY aantal DESC;
```

### User Permissions
```sql
SELECT u.naam, r.name as rol, p.name as permissie
FROM gebruikers u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE u.email = 'admin@dekoninklijkeloop.nl'
LIMIT 10;
```

### Album met Foto's
```sql
SELECT a.title, COUNT(ap.photo_id) as foto_count
FROM albums a
LEFT JOIN album_photos ap ON ap.album_id = a.id
GROUP BY a.id, a.title;
```

## 🔄 Database Reset

### Volledige Reset (Clean Slate)
```bash
# Verwijder alle data en start opnieuw
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d

# Wacht ~30 seconden voor migraties
docker-compose -f docker-compose.dev.yml logs -f app
```

### Selectieve Data Cleanup
```sql
-- Via psql
docker exec dkl-postgres psql -U postgres -d dklemailservice

-- Verwijder test contactformulieren
DELETE FROM contact_formulieren WHERE test_mode = true;

-- Verwijder test aanmeldingen  
DELETE FROM aanmeldingen WHERE test_mode = true;
```

## 📊 Database Statistieken

### Table Sizes
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 10;
```

### Active Connections
```sql
SELECT datname, count(*) as connections
FROM pg_stat_activity
WHERE datname = 'dklemailservice'
GROUP BY datname;
```

## 🎯 Frontend Testing Scenario's

### Scenario 1: Contact Form Testing
```sql
-- Voeg test contactformulier toe
INSERT INTO contact_formulieren (naam, email, bericht, status, beantwoord)
VALUES ('Test User', 'test@example.com', 'Test bericht', 'nieuw', false);
```

### Scenario 2: User Role Testing
```sql
-- Maak test staff gebruiker
INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief)
VALUES ('Staff User', 'staff@localhost.dev', '$2a$10$...', 'user', true);

-- Assign staff rol (gebruik juiste IDs)
INSERT INTO user_roles (user_id, role_id)
SELECT id, '5a1389ff-ed94-47cb-8c7f-8506ff517004'
FROM gebruikers 
WHERE email = 'staff@localhost.dev';
```

### Scenario 3: Steps Testing
```sql
-- Update steps voor deelnemer
UPDATE aanmeldingen 
SET steps = 15000 
WHERE email = 'laventejeffrey@gmail.com';

-- Check totalen
SELECT SUM(steps) as total_steps FROM aanmeldingen;
```

## 🛠️ Maintenance Commands

### Vacuum & Analyze
```sql
VACUUM ANALYZE;
```

### Check Indexes
```sql
SELECT 
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;
```

### Check Triggers
```sql
SELECT 
    trigger_name,
    event_object_table,
    action_statement
FROM information_schema.triggers
WHERE trigger_schema = 'public'
ORDER BY event_object_table;
```

## 📞 Troubleshooting

### Database Niet Bereikbaar?
```bash
# Check container status
docker-compose -f docker-compose.dev.yml ps postgres

# Check logs
docker-compose -f docker-compose.dev.yml logs postgres

# Restart
docker-compose -f docker-compose.dev.yml restart postgres
```

### Migratie Issues?
```bash
# Check migratie status
docker exec dkl-postgres psql -U postgres -d dklemailservice -c "SELECT versie, naam, toegepast FROM migraties ORDER BY toegepast DESC LIMIT 10;"
```

### Data Corrupt?
```bash
# Complete reset
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

---

**Database Type:** PostgreSQL 15  
**Schema Versie:** 1.48.0  
**Status:** ✅ Healthy  
**Last Checked:** 1 November 2025