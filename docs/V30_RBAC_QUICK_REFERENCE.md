# V30+RBAC Quick Reference Guide

**TL;DR:** Participant full accounts krijgen automatisch RBAC rollen. Temporary accounts hebben geen RBAC rollen.

---

## 🚀 Quick Commands

### Check Participant RBAC Status

```sql
-- Check of participant correct is ingesteld
SELECT 
    p.email,
    p.account_type,
    p.has_app_access,
    p.gebruiker_id IS NOT NULL as has_gebruiker,
    p.wachtwoord_hash IS NOT NULL as has_password,
    COUNT(ur.id) as roles_count
FROM participants p
LEFT JOIN user_roles ur ON p.gebruiker_id = ur.user_id AND ur.is_active = true
WHERE p.email = 'user@example.com'
GROUP BY p.id;
```

### Check User Roles

```sql
-- Welke rollen heeft een user?
SELECT 
    g.email,
    r.name as role_name,
    ur.assigned_at,
    ur.is_active
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE g.email = 'user@example.com'
ORDER BY ur.assigned_at;
```

### Check User Permissions

```sql
-- Welke permissions heeft een user?
SELECT DISTINCT
    p.resource,
    p.action,
    p.description,
    r.name as via_role
FROM user_permissions up
JOIN permissions p ON up.resource = p.resource AND up.action = p.action
JOIN roles r ON up.role_name = r.name
WHERE up.email = 'user@example.com'
ORDER BY p.resource, p.action;
```

### Manual Role Assignment (Als Trigger Faalt)

```sql
-- Wijs participant_user rol toe aan gebruiker
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT 
    g.id,
    r.id,
    CURRENT_TIMESTAMP,
    true
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'user@example.com'
  AND r.name = 'participant_user'
ON CONFLICT (user_id, role_id) DO NOTHING;
```

---

## 📊 RBAC Rollen Overzicht

### participant_user (Basis Rol)

**Auto-assigned:** ✅ Ja (bij full account creatie)  
**Permissions:** 12 permissions

**Kan:**
- ✅ App inloggen
- ✅ Stappen tracken
- ✅ Eigen data bekijken/wijzigen
- ✅ Achievements/badges verdienen
- ✅ Leaderboard bekijken
- ✅ Community deelnemen

**Kan NIET:**
- ❌ Anderen bekijken/bewerken
- ❌ Event support functies
- ❌ Admin functies

### participant_guide (Begeleiders)

**Auto-assigned:** ✅ Ja (als event rol = "Begeleider")  
**Inherits:** Alle participant_user permissions  
**Extra Permissions:** 1

**Kan (extra):**
- ✅ Community modereren
- ✅ Groepen beheren

### participant_volunteer (Vrijwilligers)

**Auto-assigned:** ✅ Ja (als event rol = "Vrijwilliger")  
**Inherits:** Alle participant_user permissions  
**Extra Permissions:** 1

**Kan (extra):**
- ✅ Event support functies
- ✅ Deelnemers assisteren (beperkt)

---

## 🔑 Kritieke Permissions

### Voor App Login (VERPLICHT)

```
app:access  ← Moet ALTIJD aanwezig zijn
app:login   ← Moet ALTIJD aanwezig zijn
```

**Zonder deze permissions:** App login faalt

### Voor Stappen Tracking (Core Feature)

```
steps:track     ← Synchroniseren van stappen
steps:view_own  ← Bekijken van eigen historie
```

### Voor Gamification

```
achievements:earn  ← Verdienen van achievements
badges:earn        ← Verdienen van badges
leaderboard:participate ← Deelnemen aan rankings
```

---

## 🐛 Common Issues & Fixes

### "Geen app toegang" error bij login

**Oorzaak:** Een van deze is niet OK:
1. account_type != 'full'
2. has_app_access = false
3. gebruiker_id IS NULL
4. Geen app:access permission

**Fix:**
```sql
-- Check status
SELECT 
    account_type, 
    has_app_access, 
    gebruiker_id,
    (SELECT COUNT(*) FROM user_permissions up 
     WHERE up.user_id = (SELECT id FROM gebruikers WHERE email = p.email)
     AND up.resource = 'app' AND up.action = 'access') as has_app_permission
FROM participants p
WHERE email = 'user@example.com';

-- Fix account type
UPDATE participants 
SET account_type = 'full', has_app_access = true
WHERE email = 'user@example.com';

-- Fix rol (als ontbreekt)
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT g.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers g, roles r
WHERE g.email = 'user@example.com' AND r.name = 'participant_user'
ON CONFLICT DO NOTHING;
```

### Permission werkt niet (cache issue)

**Fix:**
```bash
# Clear Redis cache
redis-cli FLUSHDB

# Of selectief
redis-cli --scan --pattern "perm:*" | xargs redis-cli DEL
```

### Rol automatisch toewijzing werkt niet

**Check triggers:**
```sql
-- Zijn triggers actief?
SELECT 
    tgname as trigger_name,
    tgenabled as is_enabled,
    proname as function_name
FROM pg_trigger t
JOIN pg_proc p ON t.tgfoid = p.oid
WHERE tgname LIKE '%participant%';

-- Expected:
-- trigger_assign_participant_user_role | O (enabled) | assign_participant_user_role
-- trigger_sync_participant_role_to_rbac | O (enabled) | sync_participant_role_to_rbac
```

**Re-enable trigger:**
```sql
ALTER TABLE gebruikers ENABLE TRIGGER trigger_assign_participant_user_role;
ALTER TABLE participants ENABLE TRIGGER trigger_sync_participant_role_to_rbac;
```

---

## 🧪 Quick Tests

### Test 1: Full Account Has Correct Roles

```bash
# Register
curl -X POST http://localhost:8080/api/public/aanmelden \
  -H "Content-Type: application/json" \
  -d '{"naam":"Test","email":"test@x.nl","rol":"Begeleider","afstand":"10 KM","ondersteuning":"Nee","want_account":true,"wachtwoord":"Pass1234","terms":true}'

# Check roles (should have participant_user AND participant_guide)
psql -d dkl_db -c "SELECT r.name FROM user_roles ur JOIN roles r ON ur.role_id = r.id JOIN gebruikers g ON ur.user_id = g.id WHERE g.email = 'test@x.nl';"
```

**Expected Output:**
```
      name       
-----------------
 participant_user
 participant_guide
```

### Test 2: Can Login With Full Account

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@x.nl","wachtwoord":"Pass1234"}'
```

**Expected:** 200 OK with access_token

### Test 3: Cannot Login With Temporary Account

```bash
# Register temporary
curl -X POST http://localhost:8080/api/public/aanmelden \
  -H "Content-Type: application/json" \
  -d '{"naam":"Temp","email":"temp@x.nl","rol":"Deelnemer","afstand":"6 KM","ondersteuning":"Nee","want_account":false,"terms":true}'

# Try login (should fail)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"temp@x.nl","wachtwoord":"anypass"}'
```

**Expected:** 401 Unauthorized (no gebruiker record)

---

## 📞 Emergency Procedures

### Alle Users Hebben Geen Permissions (Redis Down)

**Symptoom:** Permission checks falen voor iedereen

**Check:**
```bash
redis-cli ping
```

**Fix:**
```bash
# Herstart Redis
docker-compose restart redis

# Of in code: permission checks vallen terug op database
# Geen code wijziging nodig - dit is al ingebouwd
```

### Trigger Werkt Niet (Rollen Niet Auto-Assigned)

**Quick Fix:**
```sql
-- Run handmatig voor alle users zonder rol
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT g.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'participant_user'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  );
```

---

## 📚 Meer Informatie

- **Volledige Docs:** [`V30_RBAC_INTEGRATION.md`](V30_RBAC_INTEGRATION.md)
- **V30 Core:** [`V30_DUAL_REGISTRATION_SYSTEM.md`](V30_DUAL_REGISTRATION_SYSTEM.md)
- **RBAC API:** [`api/PERMISSIONS.md`](api/PERMISSIONS.md)

---

**Last Updated:** 2025-11-10  
**Version:** 1.0