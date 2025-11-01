# 🔐 Lokale Login Credentials

**⚠️ ALLEEN VOOR LOKALE DEVELOPMENT - NIET VOOR PRODUCTIE!**

## Admin Account (Lokaal)

```
Email:     admin@dekoninklijkeloop.nl
Wachtwoord: admin
Rol:       admin (volledige toegang)
```

### Test Login

**Via cURL:**
```bash
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "admin"
  }'
```

**Verwachte Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "...",
  "user": {
    "id": "1cdff2e3-bcce-4d58-b956-fb7c2281cd00",
    "email": "admin@dekoninklijkeloop.nl",
    "naam": "Admin",
    "rol": "admin"
  }
}
```

**Via Frontend (JavaScript/TypeScript):**
```typescript
const response = await fetch('http://localhost:8082/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@dekoninklijkeloop.nl',
    wachtwoord: 'admin'
  })
});

const data = await response.json();
console.log('Token:', data.token);
localStorage.setItem('auth_token', data.token);
```

## 🔑 Wachtwoord Details

**Bcrypt Hash:**
```
$2a$10$5Yse5i2BJV.bwTzbmywa9e/3G.XxzQPayGPlTsut/nBrZr05pKMCK
```

**Plain Text:** `admin`

**⚠️ BELANGRIJK:** API verwacht `wachtwoord` (Nederlands) niet `password`!

**Bron:** [`database/migrations/002_seed_data.sql:9`](../database/migrations/002_seed_data.sql:9)

## 📝 Extra Test Accounts Aanmaken

### Methode 1: Via SQL

```bash
# Connect to database
docker exec dkl-postgres psql -U postgres -d dklemailservice
```

```sql
-- Maak een test staff gebruiker (wachtwoord: staff123)
INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief)
VALUES (
  'Staff Test',
  'staff@localhost.dev',
  '$2a$10$4QK0pPQhJ9jBm5pZ8XW8qOE8xRm5VZD3gZKnLs4L2N4xZKdGBfF8G',
  'user',
  true
);

-- Assign staff rol
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM gebruikers u, roles r
WHERE u.email = 'staff@localhost.dev'
AND r.name = 'staff';
```

### Methode 2: Via Go Script

Gebruik het password generation script:

```bash
cd scripts/password_gen
go run generate_password_hash.go "mijn_wachtwoord"
```

Dit geeft een bcrypt hash die je kunt gebruiken in een INSERT query.

## 🎯 Permissies Admin Account

De admin account heeft **alle permissies** via de admin rol:

- ✅ Contact formulieren (read, write, delete)
- ✅ Aanmeldingen (read, write, delete) 
- ✅ Gebruikersbeheer (read, write, delete, manage_roles)
- ✅ Nieuwsbrieven (read, write, send, delete)
- ✅ Albums & Foto's (read, write, delete)
- ✅ Videos & Sponsors (read, write, delete)
- ✅ Steps tracking (read, write)
- ✅ Chat (read, write, manage_channel, moderate)
- ✅ Email management (read, write, delete, fetch)
- ✅ Admin emails (send)

**Totaal:** Toegang tot alle 68 permissies

## 🧪 Test Scenario's

### 1. Basis Login Test
```bash
# Login en sla token op
TOKEN=$(curl -s -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","password":"admin"}' \
  | jq -r '.token')

echo $TOKEN
```

### 2. Authenticated Request Test
```bash
# Gebruik token voor protected endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8082/api/contact
```

### 3. Profile Check
```bash
# Check gebruikersprofiel
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8082/api/auth/profile
```

## 🔄 Wachtwoord Resetten

### Via Database
```sql
-- Reset naar "admin" wachtwoord
UPDATE gebruikers 
SET wachtwoord_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
WHERE email = 'admin@dekoninklijkeloop.nl';
```

### Via API (als je ingelogd bent)
```bash
curl -X POST http://localhost:8082/api/auth/reset-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "huidig_wachtwoord": "admin",
    "nieuw_wachtwoord": "nieuw_wachtwoord"
  }'
```

## ⚠️ Veiligheid

**Lokaal Development:**
- Simpel wachtwoord ("admin") is OK voor lokale testing
- Database reset erased alles dus geen risico

**Productie:**
- Gebruik STERKE wachtwoorden
- Verander ALTIJD standaard credentials
- Enable 2FA indien mogelijk
- Monitor login attempts

## 📊 Verificatie Queries

```sql
-- Check alle gebruikers
SELECT naam, email, rol, is_actief FROM gebruikers;

-- Check user roles
SELECT u.naam, r.name as rol
FROM gebruikers u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id;

-- Check user permissions
SELECT u.naam, COUNT(DISTINCT p.id) as permissie_count
FROM gebruikers u
JOIN user_roles ur ON ur.user_id = u.id
JOIN role_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.id = rp.permission_id
GROUP BY u.id, u.naam;
```

## 🎓 Common Password Hashes (voor testing)

| Wachtwoord | Bcrypt Hash |
|-----------|-------------|
| `admin` | `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy` |

Voor andere wachtwoorden, gebruik `scripts/password_gen/generate_password_hash.go`

---

**⚠️ BELANGRIJK:**  
Deze credentials zijn ALLEEN voor lokale development.  
Productie gebruikt andere, veilige credentials.  
**Commit deze file NOOIT naar productie!**