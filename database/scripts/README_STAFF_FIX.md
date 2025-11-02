# Staff Role en Aanmelding Permissions Fix

## 🚨 Probleem

Staff users krijgen 403 Forbidden errors bij het ophalen van aanmeldingen. Dit kan verschillende oorzaken hebben:

1. ❌ Staff role bestaat niet in de database
2. ❌ Aanmelding permissions bestaan niet
3. ❌ Staff role heeft geen aanmelding permissions
4. ❌ Gebruikers hebben geen staff role toegewezen (legacy vs RBAC)
5. ❌ @dekoninklijkeloop.com users hebben geen staff role

## ✅ Oplossing: 3-Stappen Aanpak

### Stap 1: Diagnostiek (Optioneel maar Aanbevolen)

Voer eerst de diagnostic script uit om te zien wat het probleem is:

```bash
psql "your-database-url" -f database/scripts/diagnose_staff_permissions.sql
```

Dit script laat zien:
- ✓ Of staff role bestaat
- ✓ Welke permissions staff heeft
- ✓ Of aanmelding permissions zijn toegewezen
- ✓ Welke users staff role hebben
- ✓ Legacy vs RBAC mismatches

### Stap 2: Comprehensive Fix (Aanbevolen)

Voer het comprehensive fix script uit dat ALLES fixt:

```bash
psql "your-database-url" -f database/scripts/comprehensive_staff_fix.sql
```

**Dit script doet:**
1. ✅ Creëert staff role (als niet bestaat)
2. ✅ Creëert aanmelding permissions (als niet bestaan)
3. ✅ Kent aanmelding:read + aanmelding:write toe aan staff
4. ✅ Migreert legacy staff users (gebruikers.rol = 'staff') naar RBAC
5. ✅ Kent staff role toe aan ALLE @dekoninklijkeloop.com users
6. ✅ Verifieert en rapporteert resultaten

**Veiligheid:**
- ✅ Idempotent (veilig om meerdere keren uit te voeren)
- ✅ Gebruikt `ON CONFLICT DO NOTHING` 
- ✅ Geen data wordt verwijderd
- ✅ Alleen permissions worden toegevoegd

### Stap 3: Verificatie

Na het uitvoeren van de fix, controleer:

```sql
-- Check staff role permissions
SELECT 
    r.name as role_name,
    p.resource,
    p.action
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff' AND p.resource = 'aanmelding'
ORDER BY p.action;
```

**Verwachte output:**
```
 role_name | resource   | action
-----------+------------+--------
 staff     | aanmelding | read
 staff     | aanmelding | write
```

```sql
-- Check which users have staff role
SELECT 
    u.email,
    u.naam,
    r.name as role_name
FROM user_roles ur
JOIN gebruikers u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE r.name = 'staff' 
  AND ur.is_active = true
ORDER BY u.email;
```

## 🔧 Render.com Instructies

### Via Render.com Dashboard:

1. **Open je database:**
   - Ga naar https://dashboard.render.com
   - Selecteer je PostgreSQL database

2. **Open Shell:**
   - Klik op "Shell" tab in de database

3. **Voer het comprehensive fix script uit:**
   ```sql
   -- Copy-paste de volledige inhoud van comprehensive_staff_fix.sql
   -- Let op: dit is een groot script, maar kan in één keer uitgevoerd worden
   ```

### Via Lokale Connection:

1. **Get connection string** van Render.com dashboard

2. **Connect en execute:**
   ```bash
   psql "postgresql://user:pass@host/db?sslmode=require" \
     -f database/scripts/comprehensive_staff_fix.sql
   ```

## 🧪 Test na Fix

1. **Logout/Login:** Staff users moeten opnieuw inloggen om nieuwe permissions te krijgen

2. **Test API call:**
   ```bash
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     https://dklemailservice.onrender.com/api/aanmeldingen
   ```

3. **Verwachte response:** `200 OK` met lijst van aanmeldingen

4. **Als nog steeds 403:**
   - Check JWT token bevat staff role
   - Check user is_actief = true
   - Check user_roles.is_active = true
   - Run diagnostic script opnieuw

## 📊 Wat doet elke Script?

### diagnose_staff_permissions.sql
- **Doel:** Analyseren wat het probleem is
- **Output:** Gedetailleerd rapport van staff setup
- **Veilig:** Alleen SELECT queries, geen wijzigingen
- **Gebruik:** Voor troubleshooting

### comprehensive_staff_fix.sql  
- **Doel:** Alles fixen in één keer
- **Output:** Stap-voor-stap voortgang + final report
- **Veilig:** Idempotent, alleen INSERT/UPDATE met conflict handling
- **Gebruik:** Primaire fix script

### add_aanmelding_permissions_hotfix.sql
- **Doel:** Alleen aanmelding permissions toevoegen (minimale fix)
- **Output:** Verificatie query resultaten
- **Veilig:** Idempotent
- **Gebruik:** Als je zeker weet dat alleen permissions ontbreken

## ⚠️ Belangrijke Notes

### Legacy vs RBAC Roles

Het systeem ondersteunt BEIDE:
- **Legacy:** `gebruikers.rol` kolom (oud systeem)
- **RBAC:** `user_roles` tabel (nieuw systeem)

**De comprehensive fix migreert legacy naar RBAC maar behoudt legacy for backward compatibility.**

### @dekoninklijkeloop.com Users

Het comprehensive fix script kent **automatisch** staff role toe aan ALLE active users met @dekoninklijkeloop.com email. Dit omdat:
- Deze users zijn typisch internal staff
- Ze moeten toegang hebben tot aanmeldingen
- Safer to give access than to manually assign each time

Als je dit NIET wil, comment deze sectie uit in het script (STEP 5).

### Permission Strategy

**Staff krijgt:**
- ✅ `aanmelding:read` - Kan aanmeldingen bekijken
- ✅ `aanmelding:write` - Kan aanmeldingen bewerken/beantwoorden  
- ❌ `aanmelding:delete` - Kan GEEN aanmeldingen verwijderen (only admin)

**Admin krijgt:**
- ✅ Alle permissions (via V1_21 migration)

## 🔍 Troubleshooting

### "Role 'staff' does not exist"
```sql
INSERT INTO roles (name, description, is_system_role) VALUES
('staff', 'Ondersteunend personeel met beperkte beheerrechten', true);
```

### "Permission 'aanmelding:read' does not exist"  
```sql
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken', true);
```

### "User has staff but still 403"
1. Check JWT contains staff role:
   ```javascript
   // In browser console on frontend
   const token = localStorage.getItem('token');
   console.log(JSON.parse(atob(token.split('.')[1])));
   // Should show: { roles: ['staff'], ... }
   ```

2. Force new token by logging out/in

3. Check user_roles is_active:
   ```sql
   SELECT ur.is_active, ur.expires_at 
   FROM user_roles ur
   JOIN gebruikers u ON ur.user_id = u.id
   WHERE u.email = 'your-email@dekoninklijkeloop.com';
   ```

### "Script fails halfway"
- **Safe:** All operations use `ON CONFLICT DO NOTHING`
- **Solution:** Simply run the script again
- **It will skip already completed steps**

## 📞 Support

Als na het uitvoeren van alle fixes het probleem blijft:

1. Run diagnostic script en save output
2. Check backend logs for permission errors
3. Verify JWT token structure
4. Check if migrations V1_20-V1_26 are applied

## 🔗 Related Files

- Migration: `V1_50__add_aanmelding_permissions.sql`
- Handler: `handlers/aanmelding_handler.go`  
- Middleware: `handlers/permission_middleware.go`
- Auth Service: `services/auth_service.go`