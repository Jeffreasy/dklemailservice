# Database Migrations - Final Overview

**Laatste Update:** 2025-01-08  
**Status:** ✅ VOLLEDIG GECONSOLIDEERD EN GEVALIDEERD

---

## 📊 Migrations Overzicht

### Huidige Structuur (30 bestanden)

| Versie | Bestand | Type | Beschrijving |
|--------|---------|------|--------------|
| V01 | `V01__initial_schema.sql` | Core | Basis tabellen (users, contact, aanmelding, emails, notifications) |
| V02 | `V02__seed_data.sql` | Core | Seed admin user + email templates |
| V03 | `V03__add_test_data.sql` | Test | Test aanmeldingen en contactformulieren |
| V04 | `V04__create_chat_tables.sql` | Feature | Chat tabellen (channels, messages, reactions) |
| V05 | `V05__add_newsletter_tables.sql` | Feature | Newsletter functionaliteit |
| V06 | `V06__create_rbac_tables.sql` | Core | RBAC tabellen (roles, permissions, mappings) |
| V07 | `V07__seed_rbac_tables.sql` | Core | Seed system roles + permissions |
| V08 | `V08__migrate_and_assign_user_roles.sql` | Core | Migreer legacy roles naar RBAC |
| V09 | `V09__create_refresh_tokens_table.sql` | Auth | JWT refresh tokens |
| V10 | `V10__create_uploaded_images_table.sql` | Feature | Cloudinary image tracking |
| V11 | `V11__migrate_cms_data.sql` | Feature | CMS tabellen (photos, albums, videos, etc.) |
| V12 | `V12__add_gebruiker_id_to_aanmeldingen.sql` | Link | Link aanmeldingen aan users |
| V13 | `V13__add_steps_permissions.sql` | Feature | Steps tracking permissions |
| V14 | `V14__add_remaining_cms_permissions.sql` | Feature | Extra CMS permissions |
| V15 | `V15__replace_title_sections.sql` | Feature | Title section content |
| V16 | `V16__add_steps_to_aanmeldingen.sql` | Feature | Steps kolom toevoegen |
| **V17** | **`V17_CONSOLIDATED__create_route_funds.sql`** | **Consolidated** | **Route funds (3 files → 1)** |
| V18 | `V18__performance_optimizations.sql` | Optimize | Database performance |
| V19 | `V19__advanced_optimizations.sql` | Optimize | Geavanceerde optimalisaties |
| V20 | `V20__update_staff_aanmelding_permissions.sql` | Auth | Staff aanmelding permissions |
| V21 | `V21__add_gamification_tables.sql` | Feature | Gamification (badges, achievements) |
| V22 | `V22__improve_user_participant_linking.sql` | Link | Verbeterde user-participant link |
| V23 | `V23__add_events_table.sql` | Feature | Events tracking |
| V24 | `V24__create_notulen_module.sql` | Feature | Notulen systeem |
| V25 | `V25__create_leaderboard_materialized_view.sql` | Feature | Leaderboard view |
| **V26** | **`V26_CONSOLIDATED__normalize_roles_and_distances.sql`** | **Consolidated** | **Roles & distances (14 files → 1)** |
| **V27** | **`V27_CONSOLIDATED__normalize_status_and_type_fields.sql`** | **Consolidated** | **Status/types normalization (57 files → 1)** |
| **V28** | **`V28_CONSOLIDATED__rename_tables_participant_refactor.sql`** | **Consolidated** | **Table renaming (11 files → 1)** |
| **V29** | **`V29_CONSOLIDATED__update_permissions_for_participants.sql`** | **Consolidated** | **Permission update (1 file)** |
| **V30** | **`V30__add_is_active_to_participant_roles.sql`** | **Fix** | **is_active kolom fix** |

---

## 🎯 Consolidatie Resultaten

### Voor Consolidatie:
- **V01-V16:** 16 bestanden (al redelijk geconsolideerd)
- **V17:** 3 bestanden (route_funds)
- **V18-V25:** 8 bestanden (standalone features)
- **V26-V29:** 83 bestanden (zeer gefragmenteerd)
- **Totaal:** 110 migration bestanden

### Na Consolidatie:
- **V01-V16:** 16 bestanden (behouden)
- **V17:** 1 bestand (✅ geconsolideerd)
- **V18-V25:** 8 bestanden (behouden)
- **V26-V29:** 4 bestanden (✅ geconsolideerd)
- **V30:** 1 bestand (nieuwe fix)
- **Totaal:** 30 migration bestanden

**Reductie:** 110 → 30 bestanden (73% minder files!)

---

## ✅ Code Alignment Status

### Alle Kritieke Fixes Toegepast:

1. ✅ **Distance Model** - `FundAmount` field toegevoegd
   - File: [`models/distance.go`](../models/distance.go:6)
   - Matching: `distances.fund_amount` kolom

2. ✅ **ParticipantRole Model** - `IsActive` field toegevoegd
   - File: [`models/participant_role.go`](../models/participant_role.go:6)
   - Migration: [`V30__add_is_active_to_participant_roles.sql`](V30__add_is_active_to_participant_roles.sql)
   - Database: Kolom toegevoegd en gevuld

3. ✅ **Status Type Models** - Alle field names gecorrigeerd
   - File: [`models/status_types.go`](../models/status_types.go:1)
   - Changes: `Type` → `Status` of `Name` (afhankelijk van tabel)

4. ✅ **FK References** - Alle 7 references gecorrigeerd
   - [`models/event_registration.go`](../models/event_registration.go:39)
   - [`models/contact.go`](../models/contact.go:20)
   - [`models/event.go`](../models/event.go:52)
   - [`models/verzonden_email.go`](../models/verzonden_email.go:17)
   - [`models/notification.go`](../models/notification.go:48)
   - [`models/event.go`](../models/event.go:138) - Status check
   - [`models/event.go`](../models/event.go:221) - Status check

### Build Status:
```bash
go build -o test.exe
# Exit code: 0 ✅ - Geen compiler errors!
```

---

## 📋 Migration Execution Volgorde

Voor een **nieuwe database** (van scratch):

```
1.  V01__initial_schema.sql
2.  V02__seed_data.sql
3.  V03__add_test_data.sql
4.  V04__create_chat_tables.sql
5.  V05__add_newsletter_tables.sql
6.  V06__create_rbac_tables.sql
7.  V07__seed_rbac_tables.sql
8.  V08__migrate_and_assign_user_roles.sql
9.  V09__create_refresh_tokens_table.sql
10. V10__create_uploaded_images_table.sql
11. V11__migrate_cms_data.sql
12. V12__add_gebruiker_id_to_aanmeldingen.sql
13. V13__add_steps_permissions.sql
14. V14__add_remaining_cms_permissions.sql
15. V15__replace_title_sections.sql
16. V16__add_steps_to_aanmeldingen.sql
17. V17_CONSOLIDATED__create_route_funds.sql       ⭐ Consolidated
18. V18__performance_optimizations.sql
19. V19__advanced_optimizations.sql
20. V20__update_staff_aanmelding_permissions.sql
21. V21__add_gamification_tables.sql
22. V22__improve_user_participant_linking.sql
23. V23__add_events_table.sql
24. V24__create_notulen_module.sql
25. V25__create_leaderboard_materialized_view.sql
26. V26_CONSOLIDATED__normalize_roles_and_distances.sql       ⭐ Consolidated
27. V27_CONSOLIDATED__normalize_status_and_type_fields.sql    ⭐ Consolidated
28. V28_CONSOLIDATED__rename_tables_participant_refactor.sql  ⭐ Consolidated  
29. V29_CONSOLIDATED__update_permissions_for_participants.sql ⭐ Consolidated
30. V30__add_is_active_to_participant_roles.sql              🆕 New Fix
```

---

## 🎯 Key Consolidated Migrations

### V17_CONSOLIDATED (Route Funds)
**Consolidates:** V17_01, V17_02, V17_03  
**Purpose:** Create route_funds table, index, and seed data  
**File:** [`V17_CONSOLIDATED__create_route_funds.sql`](migrations/V17_CONSOLIDATED__create_route_funds.sql)

### V26_CONSOLIDATED (Roles & Distances)  
**Consolidates:** V26_01 through V26_14 (14 files)  
**Purpose:** Normalize participant roles and distances to lookup tables  
**File:** [`V26_CONSOLIDATED__normalize_roles_and_distances.sql`](migrations/V26_CONSOLIDATED__normalize_roles_and_distances.sql)

### V27_CONSOLIDATED (Status & Types)
**Consolidates:** V27_01 through V27_57 (57 files)  
**Purpose:** Normalize all status/type fields with lookup tables  
**File:** [`V27_CONSOLIDATED__normalize_status_and_type_fields.sql`](migrations/V27_CONSOLIDATED__normalize_status_and_type_fields.sql)

### V28_CONSOLIDATED (Table Renaming)
**Consolidates:** V28_01 through V28_11 (11 files)  
**Purpose:** Rename Dutch tables to English + refactor participant/event  
**File:** [`V28_CONSOLIDATED__rename_tables_participant_refactor.sql`](migrations/V28_CONSOLIDATED__rename_tables_participant_refactor.sql)

### V29_CONSOLIDATED (Permission Update)
**Consolidates:** V29 (1 file, already consolidated)  
**Purpose:** Update permission resources to match renamed tables  
**File:** [`V29_CONSOLIDATED__update_permissions_for_participants.sql`](migrations/V29_CONSOLIDATED__update_permissions_for_participants.sql)

### V30 (New Fix)
**Purpose:** Add is_active column to participant_roles  
**File:** [`V30__add_is_active_to_participant_roles.sql`](migrations/V30__add_is_active_to_participant_roles.sql)

---

## ✅ Database Validatie

Gevalideerd op `dkl-postgres` Docker container:

```sql
-- All lookup tables present
✓ participant_roles (with is_active)
✓ distances (with fund_amount)
✓ contact_status_types
✓ registration_status_types  
✓ email_status_types
✓ event_status_types
✓ chat_channel_types
✓ notification_types
✓ notification_priority_types

-- All tables renamed correctly
✓ participants (was: aanmeldingen)
✓ event_registrations (was: event_participants)
✓ participant_antwoorden (was: aanmelding_antwoorden)

-- All foreign keys working
✓ event_registrations.participant_role_name → participant_roles.name
✓ event_registrations.distance_route → distances.route
✓ event_registrations.status → registration_status_types.status
```

---

##  Testing

### Local Docker Test
```powershell
cd database
.\test_consolidated_migrations.ps1
```

### MCP Server Validation
Zie [`validate_with_mcp.json`](validate_with_mcp.json) voor alle validatie commando's.

---

## 📚 Documentation

- **[`MIGRATIONS_FINAL_OVERVIEW.md`](MIGRATIONS_FINAL_OVERVIEW.md)** - Deze overview
- **[`VALIDATION_REPORT.md`](VALIDATION_REPORT.md)** - Database validatie details
- **[`MIGRATION_CONSOLIDATION_SUMMARY.md`](../MIGRATION_CONSOLIDATION_SUMMARY.md)** - Consolidatie samenvatting

---

## 🚀 Production Ready

✅ Alle migrations geconsolideerd  
✅ Alle code fixes toegepast  
✅ Database gevalideerd  
✅ Build succesvol  
✅ 73% minder files  
✅ 100% database-code alignment

**Het systeem is klaar voor gebruik!**