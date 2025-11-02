# Aanmelding Permissions Hotfix

## 🚨 Probleem
Staff users krijgen 403 Forbidden errors bij het ophalen van aanmeldingen omdat de `aanmelding` permissions nog niet bestaan in de productie database.

## ✅ Oplossing
Voer het hotfix script uit op de productie database om de permissions onmiddellijk toe te voegen.

## 📋 Instructies

### Optie 1: Via Render.com Dashboard (Aanbevolen)

1. **Open Render.com Dashboard:**
   - Ga naar https://dashboard.render.com
   - Selecteer je PostgreSQL database

2. **Open Shell:**
   - Klik op "Shell" tab
   - Of gebruik "Connect" → "External Connection" voor lokale psql

3. **Voer het script uit:**
   ```bash
   # Copy de inhoud van add_aanmelding_permissions_hotfix.sql
   # Paste in de shell en voer uit
   ```

### Optie 2: Via lokale psql

1. **Connecteer met productie database:**
   ```bash
   psql "postgresql://[username]:[password]@[host]/[database]?sslmode=require"
   ```

2. **Voer het script uit:**
   ```bash
   \i database/scripts/add_aanmelding_permissions_hotfix.sql
   ```

### Optie 3: Via pgAdmin

1. **Connect met productie database**
2. **Open Query Tool**
3. **Load en execute:** `database/scripts/add_aanmelding_permissions_hotfix.sql`

## 🔍 Verificatie

Na het uitvoeren van het script, controleer of de permissions zijn toegevoegd:

```sql
-- Check permissions
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    p.description
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.resource = 'aanmelding'
ORDER BY r.name, p.action;
```

**Verwachte output:**
```
 role_name | resource   | action  | description
-----------+------------+---------+----------------------------
 admin     | aanmelding | delete  | Aanmeldingen verwijderen
 admin     | aanmelding | read    | Aanmeldingen bekijken
 admin     | aanmelding | write   | Aanmeldingen bewerken
 staff     | aanmelding | read    | Aanmeldingen bekijken
 staff     | aanmelding | write   | Aanmeldingen bewerken
```

## 🧪 Test na Fix

1. **Login als staff user** in de frontend
2. **Navigate naar Dashboard**
3. **Controleer dat aanmeldingen zichtbaar zijn** (geen 403 error)

## 📝 Wat doet dit script?

1. ✅ Voegt 3 permissions toe: `aanmelding:read`, `aanmelding:write`, `aanmelding:delete`
2. ✅ Geeft admin role alle 3 permissions
3. ✅ Geeft staff role `read` en `write` permissions (geen delete)
4. ✅ Toont verificatie query om te bevestigen dat het werkt

## ⚠️ Belangrijk

- Dit is een **hotfix** voor het onmiddellijke probleem
- De V1.50 migration zal later automatisch worden uitgevoerd
- Het script is **idempotent** - veilig om meerdere keren uit te voeren
- Gebruikt `ON CONFLICT DO NOTHING` om duplicaten te voorkomen

## 🔄 Toekomstige Deploys

Bij de volgende deployment zal de V1.50 migration automatisch draaien via het normale migration systeem. Deze hotfix zorgt ervoor dat staff users **nu** toegang hebben zonder te wachten op een volledige redeploy.

## 📞 Bij Problemen

Als je na het uitvoeren van het script nog steeds 403 errors ziet:

1. **Check user roles:**
   ```sql
   SELECT u.email, u.naam, r.name as role_name
   FROM user_roles ur
   JOIN gebruikers u ON ur.user_id = u.id
   JOIN roles r ON ur.role_id = r.id
   WHERE u.email = '[staff-email@example.com]'
     AND ur.is_active = true;
   ```

2. **Verify permissions:**
   ```sql
   SELECT p.resource, p.action
   FROM user_roles ur
   JOIN role_permissions rp ON ur.role_id = rp.role_id
   JOIN permissions p ON rp.permission_id = p.id
   WHERE ur.user_id = '[user-id]'
     AND ur.is_active = true
     AND p.resource = 'aanmelding';
   ```

3. **Force logout/login:** Staff user moet mogelijk opnieuw inloggen om nieuwe permissions te krijgen