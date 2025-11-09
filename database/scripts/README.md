# Database Scripts Consolidatie

**Status:** Volledig geconsolideerd - 2025-11-08

## 📋 Overzicht

Na grondige analyse zijn alle 16 database scripts beoordeeld en geconsolideerd. Dit document beschrijft welke scripts behouden blijven en waarom.

---

## ✅ Scripts die BEHOUDEN blijven (5 scripts)

### 1. `check_rbac_db.sh` - RBAC Database Verificatie Automatisering
**Status:** ✅ ESSENTIEEL - BEHOUDEN
- **Functie:** Shell wrapper die automatisch RBAC verificatie uitvoert
- **Wanneer gebruiken:** Wekelijks of na RBAC wijzigingen
- **Features:**
  - Auto-detecteert Docker/local PostgreSQL
  - Roept `verify_rbac_tables.sql` aan
  - Genereert overzichtelijk rapport
- **Gebruik:** `./database/scripts/check_rbac_db.sh`

### 2. `verify_rbac_tables.sql` - Comprehensive RBAC Verificatie
**Status:** ✅ ESSENTIEEL - BEHOUDEN
- **Functie:** 17-punts verificatie van complete RBAC setup
- **Wanneer gebruiken:** 
  - Na elke RBAC migration
  - Bij permission/role problemen
  - Maandelijkse health check
- **Controleert:**
  - Table existence en structure
  - System roles (9 verwacht)
  - Permissions coverage (58+ verwacht)
  - User role assignments
  - Legacy vs RBAC sync
  - Foreign key constraints
  - Orphaned records
- **Gebruik:** `psql $DATABASE_URL -f database/scripts/verify_rbac_tables.sql`

### 3. `data_cleanup.sql` - Data Archivering & Cleanup
**Status:** ✅ REGULIER ONDERHOUD - BEHOUDEN
- **Functie:** Maandelijks onderhoud voor data archivering
- **Wanneer gebruiken:** Maandelijks
- **Acties:**
  - Archiveert oude verzonden emails (>1 jaar)
  - Archiveert processed incoming emails (>6 maanden)
  - Verwijdert expired refresh tokens (>30 dagen)
  - Permanent verwijderen soft-deleted images (>3 maanden)
  - VACUUM ANALYZE voor ruimte terugwinning
- **⚠️ BELANGRIJK:** ALTIJD eerst backup maken!
- **Gebruik:** `psql $DATABASE_URL -f database/scripts/data_cleanup.sql`

### 4. `vacuum_analyze.sql` - Database Performance Onderhoud
**Status:** ✅ REGULIER ONDERHOUD - BEHOUDEN
- **Functie:** Wekelijks performance onderhoud
- **Wanneer gebruiken:** Wekelijks of na grote data wijzigingen
- **Acties:**
  - VACUUM ANALYZE op alle tables
  - Rapporteert table sizes
  - Toont dead tuples (bloat)
  - Index usage statistics
  - Identificeert ongebruikte indexes
- **Gebruik:** `psql $DATABASE_URL -f database/scripts/vacuum_analyze.sql`

### 5. `maintenance_rbac_users.sql` - RBAC User Onderhoud (NIEUW)
**Status:** ✅ GECONSOLIDEERD - NIEUW SCRIPT
- **Functie:** Consolidated best practices voor RBAC user management
- **Wanneer gebruiken:**
  - Na toevoegen nieuwe staff members
  - Bij domain-based role toewijzingen
  - Voor user-role synchronisatie
- **Features:**
  - Assign staff role to @dekoninklijkeloop.nl users
  - Sync legacy roles naar RBAC
  - Diagnostic checks
  - Verificatie rapportage
- **Gebruik:** `psql $DATABASE_URL -f database/scripts/maintenance_rbac_users.sql`

---

## 🗑️ Scripts die VERWIJDERD zijn (11 scripts)

### Eenmalige Hotfixes (al uitgevoerd in migrations)
- ❌ `add_aanmelding_permissions_hotfix.sql` - Nu in V1_22 migration
- ❌ `comprehensive_staff_fix.sql` - Geconsolideerd in maintenance_rbac_users.sql
- ❌ `RENDER_DEPLOY_NOW.sql` - Oude production hotfix, nu in migrations
- ❌ `fix_data_quality_for_v1_48.sql` - Pre-V1.48 fix, eenmalig
- ❌ `fix_v1_58_migration.sql` - Specifieke V1.58 fix, eenmalig

### Test & Diagnostic Scripts (obsolete)
- ❌ `diagnose_staff_permissions.sql` - Geconsolideerd in verify_rbac_tables.sql
- ❌ `test_v1_47_indexes.sql` - Oude migration test
- ❌ `check_event_date.sql` - Specifieke datum check

### Data Scripts (eenmalig of test data)
- ❌ `create_test_event_dronten.sql` - Test event data
- ❌ `update_event_date_to_16_mei.sql` - Specifieke datum update
- ❌ `sync_production_users_with_rbac.sql` - Geconsolideerd in maintenance_rbac_users.sql

---

## 📁 Optionele Scripts (aparte folder aanbevolen)

### `setup_partitioning.sql`
**Status:** ⚠️ OPTIONEEL - Toekomstige schaalbaarheid
- **Functie:** Table partitioning setup voor grote tables
- **Wanneer gebruiken:** Alleen wanneer tables >10 miljoen records
- **⚠️ BELANGRIJK:** 
  - Vereist downtime
  - Production backup verplicht
  - Niet nodig voor huidige schaal
- **Aanbeveling:** Verplaats naar `database/scripts/advanced/` folder

---

## 🔄 Gebruik Workflow

### Dagelijks
- Geen scripts nodig

### Wekelijks
```bash
# Performance onderhoud
./database/scripts/vacuum_analyze.sql
```

### Maandelijks
```bash
# RBAC verificatie
./database/scripts/check_rbac_db.sh

# Data cleanup (met backup!)
./database/scripts/data_cleanup.sql
```

### Bij Problemen
```bash
# RBAC issues
./database/scripts/verify_rbac_tables.sql

# User role sync
./database/scripts/maintenance_rbac_users.sql
```

### Na Migrations
```bash
# Verifieer RBAC integriteit
./database/scripts/check_rbac_db.sh
```

---

## 📊 Consolidatie Details

### Van 16 naar 5 scripts
- **Behouden:** 4 essentiële onderhoud scripts + 1 nieuw geconsolideerd script
- **Verwijderd:** 11 obsolete/eenmalige scripts
- **Reductie:** 69% minder scripts, 100% behoud van functionaliteit

### Geconsolideerde Functionaliteit
De nieuwe `maintenance_rbac_users.sql` consolideert:
- Beste delen van `comprehensive_staff_fix.sql`
- User sync logica van `sync_production_users_with_rbac.sql`
- Diagnostic queries van `diagnose_staff_permissions.sql`

---

## ⚡ Quick Reference

| Script | Frequentie | Doel |
|--------|-----------|------|
| `check_rbac_db.sh` | Wekelijks/Na wijzigingen | RBAC verificatie |
| `verify_rbac_tables.sql` | Maandelijks | Diepgaande RBAC check |
| `data_cleanup.sql` | Maandelijks | Data archivering |
| `vacuum_analyze.sql` | Wekelijks | Performance |
| `maintenance_rbac_users.sql` | Bij behoefte | User role sync |

---

## 🔐 Security Notes

- **data_cleanup.sql:** Maakt permanent backups, controleer eerst!
- **maintenance_rbac_users.sql:** Wijzigt user roles, test eerst in staging!
- **Alle scripts:** Geen credentials in scripts, gebruik environment variables

---

## 📚 Gerelateerde Documentatie

- [DATABASE_DOC.md](../../01DATABASE_DOC.md) - Database architectuur
- [AUTHENTICATION_DOC.md](../../02AUTHENTICATION_DOC.md) - RBAC system documentatie
- [MIGRATIONS_FINAL_OVERVIEW.md](../MIGRATIONS_FINAL_OVERVIEW.md) - Migration geschiedenis

---

**Laatst bijgewerkt:** 2025-11-08
**Consolidatie door:** Kilo Code AI
**Status:** ✅ Productie-klaar