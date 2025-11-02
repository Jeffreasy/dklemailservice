# 🚀 HOE TE DEPLOYEN NAAR RENDER.COM POSTGRESQL

## ⚡ Snelle Instructies (2 minuten)

### Stap 1: Open Render.com Dashboard
1. Ga naar https://dashboard.render.com
2. Login met je account
3. Selecteer je **PostgreSQL database** uit de lijst

### Stap 2: Open Database Shell
1. Klik op de database naam
2. Klik op **"Shell"** tab bovenaan
3. Of klik op **"Connect"** → **"psql"**

### Stap 3: Copy-Paste & Execute
1. **Open dit bestand:** [`RENDER_DEPLOY_NOW.sql`](RENDER_DEPLOY_NOW.sql)
2. **Select ALL** (Ctrl+A / Cmd+A)
3. **Copy** (Ctrl+C / Cmd+C)
4. **Paste** in Render.com Shell
5. **Press Enter** om uit te voeren

## 📋 Wat Gebeurt Er?

Het script voert 6 stappen uit:

1. ✅ Creëert staff role (als niet bestaat)
2. ✅ Creëert aanmelding permissions (read, write, delete)
3. ✅ Kent read + write toe aan staff role
4. ✅ Kent alle permissions toe aan admin role
5. ✅ Migreert legacy staff users naar RBAC
6. ✅ Kent staff role toe aan @dekoninklijkeloop.com users

## ✅ Verificatie

Na uitvoering zie je:

```
✓ Staff Permissions:
  staff | aanmelding | read   | GRANTED
  staff | aanmelding | write  | GRANTED

✓ Users with Staff Role:
  lida@dekoninklijkeloop.com | Lida | true | 2025-11-02...
  [...andere staff users...]
```

## 🧪 Test Direct Na Uitvoering

1. **Staff user MOET logout/login** (om nieuwe JWT te krijgen)
2. **Test:** Navigate naar Dashboard
3. **Verwacht:** Aanmeldingen zijn zichtbaar (geen 403)

## ⚠️ Als Je Geen Shell Toegang Hebt

### Optie A: Via pgAdmin

1. Get connection string van Render.com:
   - Dashboard → Database → "Info"
   - Copy "External Database URL"

2. Open pgAdmin:
   - Create New Server
   - Paste connection details
   - Open Query Tool

3. Load SQL file:
   - File → Open → `RENDER_DEPLOY_NOW.sql`
   - Execute (F5)

### Optie B: Via DBeaver

1. New Connection → PostgreSQL
2. Paste connection details from Render
3. SQL Editor → Load `RENDER_DEPLOY_NOW.sql`
4. Execute

### Optie C: Via Lokale psql

```bash
# Get connection string from Render.com
# Should look like: postgresql://user:pass@host/dbname

psql "postgresql://your-connection-string" \
  -f database/scripts/RENDER_DEPLOY_NOW.sql
```

## 🔒 Veiligheid

- ✅ **Idempotent:** Veilig om meerdere keren uit te voeren
- ✅ **No Data Loss:** Gebruikt `ON CONFLICT DO NOTHING`
- ✅ **Read-Only Staff:** Staff krijgt alleen read + write (geen delete)
- ✅ **Verified:** Script is getest en veilig

## 📞 Bij Problemen

Als script niet werkt:

1. **Check permissions:** Je moet SUPERUSER of OWNER zijn van de database
2. **Check migration:** V1_20-V1_21 migrations moeten uitgevoerd zijn
3. **Manual table check:**
   ```sql
   -- Check if tables exist
   SELECT tablename FROM pg_tables 
   WHERE schemaname = 'public' 
   AND tablename IN ('roles', 'permissions', 'role_permissions', 'user_roles');
   ```

Als tables niet bestaan:
- V1.20 migration is niet uitgevoerd
- Render.com heeft migrations nog niet gedraaid
- Wacht op automatische deployment of run migrations handmatig

## ✅ Success Criteria

Na succesvol uitvoeren:

```sql
-- Dit moet 2 rows tonen
SELECT count(*) FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff' AND p.resource = 'aanmelding';
-- Expected: 2 (read + write)

-- Dit moet je staff users tonen
SELECT count(*) FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE r.name = 'staff' AND ur.is_active = true;
-- Expected: Aantal staff users (minimaal 1)
```

## 🎯 Resultaat

**NA dit script:**
- ✅ Staff users kunnen aanmeldingen zien
- ✅ Staff users kunnen aanmeldingen bewerken
- ✅ Staff users kunnen antwoorden sturen
- ❌ Staff users kunnen NIET verwijderen (only admin)

**Tijd:** ~30 seconden  
**Risico:** Geen (alleen INSERT met conflict handling)  
**Impact:** Onmiddellijk (na logout/login)