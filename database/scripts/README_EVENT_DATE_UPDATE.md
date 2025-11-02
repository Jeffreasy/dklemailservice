# Event Datum Update: 16 mei 2026

## 📅 Context

De Koninklijke Loop (DKL) 2026 vindt plaats op **16 mei 2026** (niet 17 mei).

Volgens de [`production_backup.sql`](../../production_backup.sql:2273) is de database **AL CORRECT** met "16 mei 2026".

## ✅ Huidige Status

De `title_section_content` tabel bevat:
- **detail_1_title:** "16 mei 2026" ✅ CORRECT
- **detail_1_description:** Vermeldt correcte datum
- **image_alt:** "Wij gaan 16 mei lopen voor hen" ✅ CORRECT

## 🔍 Verificatie Scripts

### 1. Check Event Datum
Gebruik [`check_event_date.sql`](check_event_date.sql) om de huidige datum te verifiëren:

```sql
-- Run in pgAdmin Query Tool
SELECT 
    id,
    event_title,
    detail_1_title,
    detail_1_description,
    created_at,
    updated_at
FROM title_section_content
WHERE detail_1_title LIKE '%mei 2026%';
```

**Verwacht resultaat:**
- `detail_1_title` = "16 mei 2026" ✅

**Als je "17 mei 2026" ziet**, ga dan naar stap 2.

### 2. Update Datum (alleen als nodig)
Als de datum incorrect is (17 mei), gebruik dan [`update_event_date_to_16_mei.sql`](update_event_date_to_16_mei.sql):

#### pgAdmin Instructies:
1. ✅ Open pgAdmin 4
2. ✅ Connect met database server
3. ✅ Selecteer database: `dklemailservice`
4. ✅ Klik: **Tools → Query Tool** (of druk F5)
5. ✅ Kopieer de **volledige inhoud** van `update_event_date_to_16_mei.sql`
6. ✅ Plak in Query Tool venster
7. ✅ Klik **Execute** (▶ knop) of druk F5
8. ✅ Bekijk de output in de Messages tab

#### Wat doet het script?
- ✅ Toont huidige staat VOOR update
- ✅ Update `detail_1_title` van "17 mei" → "16 mei"
- ✅ Update `detail_1_description` tekst
- ✅ Update `image_alt` tekst
- ✅ Set `updated_at` timestamp
- ✅ Toont verificatie resultaten NA update
- ✅ Toont alle event records voor overzicht
- ✅ Auto-commit

#### Verwachte Output:
```
NOTICE: ═══════════════════════════════════════════════════
NOTICE: HUIDIGE STAAT VOOR UPDATE:
NOTICE: ═══════════════════════════════════════════════════
NOTICE: Huidige detail_1_title: 17 mei 2026
NOTICE: Aantal records met "17 mei 2026": 1

UPDATE 1

NOTICE: ═══════════════════════════════════════════════════
NOTICE: UPDATE RESULTAAT:
NOTICE: ═══════════════════════════════════════════════════
NOTICE: Nieuwe detail_1_title: 16 mei 2026
NOTICE: Records met "17 mei 2026": 0 (moet 0 zijn)
NOTICE: Records met "16 mei 2026": 1
NOTICE: ✅ SUCCESS: Alle datums zijn correct bijgewerkt naar 16 mei 2026
```

## 🔧 Docker Database Access

Als je via Docker moet verbinden:

### Via Docker Compose:
```bash
# Development database
docker exec -it dkl-postgres psql -U postgres -d dklemailservice -f /path/to/script.sql

# Production/Render
# Gebruik pgAdmin of psql met juiste credentials
```

### Connection Info:
- **Host:** localhost (dev) of Render host (prod)
- **Port:** 5433 (dev) of 5432 (prod)
- **Database:** dklemailservice
- **User:** postgres
- **Password:** Zie `.env` of Render dashboard

## 📝 Betreffende Tabellen

### title_section_content
```sql
CREATE TABLE title_section_content (
    id UUID PRIMARY KEY,
    event_title VARCHAR(255),
    event_subtitle TEXT,
    image_url VARCHAR(500),
    image_alt VARCHAR(255),
    detail_1_title VARCHAR(100),      -- ← Event datum staat hier
    detail_1_description TEXT,        -- ← En hier in beschrijving
    detail_2_title VARCHAR(100),
    detail_2_description TEXT,
    detail_3_title VARCHAR(100),
    detail_3_description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    participant_count INTEGER
);
```

## 🚀 Frontend Impact

Na het updaten van de database:
- ✅ Frontend haalt nieuw data op via API
- ✅ Cloudinary images blijven werken (image URL wijzigt niet)
- ✅ Geen frontend code changes nodig
- ✅ Cache kan snel verversen

## ⚠️ Waarschuwing

**BELANGRIJK:** Dit script update de **PRODUCTION DATA**!

- ✅ Maak eerst een backup met pgAdmin
- ✅ Test eerst in development database
- ✅ Run tijdens onderhoudsvenster (indien mogelijk)
- ✅ Verificeer resultaat na uitvoeren

## 📊 Backup Maken (optioneel maar aanbevolen)

Voor de update:
```sql
-- Backup van huidige staat
CREATE TABLE title_section_content_backup AS 
SELECT * FROM title_section_content;

-- Of via pgAdmin:
-- Right-click op tabel → Backup...
```

Na succesvolle update:
```sql
-- Verwijder backup tabel
DROP TABLE IF EXISTS title_section_content_backup;
```

## 🔗 Gerelateerde Bestanden

- [`check_event_date.sql`](check_event_date.sql) - Verificatie script
- [`update_event_date_to_16_mei.sql`](update_event_date_to_16_mei.sql) - Update script
- [`production_backup.sql`](../../production_backup.sql) - Production backup met correcte datum

## 📞 Support

Bij problemen:
1. Check database connectie
2. Verify table bestaat: `SELECT * FROM title_section_content LIMIT 1;`
3. Check permissions: User moet UPDATE rechten hebben
4. Zie logs in pgAdmin Messages tab

---

**Last Updated:** 2025-11-02  
**Author:** Development Team  
**Status:** ✅ Database is correct (16 mei 2026)