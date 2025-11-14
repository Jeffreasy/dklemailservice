# V30 Duaal Registratiesysteem - Quick Start Guide

**TL;DR:** Backend implementatie voor dual registration (full vs temporary accounts) is compleet. Frontend moet nog geïmplementeerd worden.

---

## ⚡ Quick Reference

### Nieuwe API Endpoints

```bash
# 1. Publieke Registratie (GEEN auth vereist)
POST /api/public/aanmelden
Body: { naam, email, rol, afstand, ondersteuning, want_account, wachtwoord?, terms }

# 2. Account Upgrade (GEEN auth vereist)
POST /api/public/upgrade-to-full-account
Body: { email, wachtwoord }

# 3. Actief Event Info (GEEN auth vereist)
GET /api/public/events/active
```

---

## 🚀 Backend Testen

### 1. Test Full Account Registratie
```bash
curl -X POST http://localhost:8080/api/public/aanmelden \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "Full Test",
    "email": "full@test.nl",
    "rol": "Deelnemer",
    "afstand": "10 KM",
    "ondersteuning": "Nee",
    "want_account": true,
    "wachtwoord": "TestPass123",
    "terms": true
  }'
```

**Verwacht:** 201 Created + `account_type: "full"` + `gebruiker_id`

### 2. Test Temporary Registratie
```bash
curl -X POST http://localhost:8080/api/public/aanmelden \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "Temp Test",
    "email": "temp@test.nl",
    "rol": "Begeleider",
    "telefoon": "06-12345678",
    "afstand": "6 KM",
    "ondersteuning": "Nee",
    "want_account": false,
    "terms": true
  }'
```

**Verwacht:** 201 Created + `account_type: "temporary"` + geen `gebruiker_id`

### 3. Test Account Upgrade
```bash
curl -X POST http://localhost:8080/api/public/upgrade-to-full-account \
  -H "Content-Type: application/json" \
  -d '{
    "email": "temp@test.nl",
    "wachtwoord": "NewPass123"
  }'
```

**Verwacht:** 200 OK + `gebruiker_id` + `has_app_access: true`

---

## 📁 Nieuwe Bestanden Overzicht

### Backend (✅ Compleet)
```
database/migrations/
  └─ V30__dual_registration_system.sql          # Database migratie

handlers/
  └─ public_registration_handler.go              # Publieke registratie logica

models/
  └─ participant.go                              # Updated met V30 velden

templates/
  ├─ registration_full_account_email.html        # Full account email
  ├─ registration_temporary_account_email.html   # Temporary email
  └─ account_upgrade_email.html                  # Upgrade email

docs/
  ├─ V30_DUAL_REGISTRATION_SYSTEM.md            # Volledige backend docs
  ├─ V30_IMPLEMENTATION_SUMMARY.md              # Deze samenvatting
  ├─ V30_QUICK_START.md                         # Deze guide
  └─ frontend/
      └─ V30_FRONTEND_IMPLEMENTATION_GUIDE.md   # Frontend stap-voor-stap
```

### Frontend (⏳ Te Implementeren)
```
DKL25/src/pages/Aanmelden/
  ├─ components/
  │   ├─ AccountTypeSelector.tsx                # NIEUW - Maak aan
  │   ├─ PasswordField.tsx                      # NIEUW - Maak aan
  │   ├─ FormContainer.tsx                      # UPDATE - Voeg account choice toe
  │   └─ SuccessMessage.tsx                     # UPDATE - Conditionele content
  ├─ types/
  │   └─ schema.ts                              # UPDATE - Nieuwe velden
  └─ Upgrade/
      └─ UpgradeAccount.tsx                     # NIEUW - Upgrade pagina
```

---

## 🎯 Frontend Implementatie (5 Stappen)

### Stap 1: Update Schema (5 min)
Voeg toe aan `schema.ts`:
```typescript
want_account: z.boolean(),
wachtwoord: z.string().min(8).optional(),
```

### Stap 2: Maak Components (20 min)
- `AccountTypeSelector.tsx` - Keuze tussen full/temporary
- `PasswordField.tsx` - Wachtwoord input met strength

### Stap 3: Update FormContainer (15 min)
- Import nieuwe components
- Voeg account choice sectie toe (na contactgegevens)
- Voeg conditioneel password veld toe
- Update API call naar `/api/public/aanmelden`

### Stap 4: Update SuccessMessage (10 min)
- Voeg app download sectie toe (voor full)
- Voeg upgrade CTA toe (voor temporary)

### Stap 5: Maak Upgrade Pagina (30 min)
- Nieuwe route `/upgrade`
- Form voor email + password
- API call naar `/api/public/upgrade-to-full-account`

**Totaal: ~80 minuten werk**

---

## 📊 Database Verificatie

Na backend start, check of migratie succesvol was:

```sql
-- Check nieuwe kolommen
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'participants'
AND column_name IN ('account_type', 'registration_year', 'wachtwoord_hash', 'has_app_access');

-- Moet 4 rijen retourneren
```

---

## 🐛 Troubleshooting

### Backend start faalt
```bash
# Check logs
docker-compose logs app

# Zoek naar migratie errors
docker-compose logs app | grep "V30"
```

### Endpoint geeft 404
```bash
# Check of handler geregistreerd is
docker-compose logs app | grep "Public registration routes"
```

### Email wordt niet verzonden
```bash
# Check email service logs
docker-compose logs app | grep "confirmation email"
```

---

## 📞 Hulp Nodig?

**Backend Problemen:**
- Zie: [`docs/V30_DUAL_REGISTRATION_SYSTEM.md`](V30_DUAL_REGISTRATION_SYSTEM.md)
- Check: Database migratie logs
- Verify: Alle environment variables set

**Frontend Implementatie:**
- Zie: [`docs/frontend/V30_FRONTEND_IMPLEMENTATION_GUIDE.md`](frontend/V30_FRONTEND_IMPLEMENTATION_GUIDE.md)
- Voorbeelden van alle components included
- Copy-paste ready code snippets

**API Testen:**
- Gebruik: Postman/Insomnia
- Collection: Import curl commands uit docs
- Check: Response structures match DTOs

---

## ✅ Deployment Checklist

### Development
- [x] Database migratie gemaakt
- [x] Backend code geïmplementeerd
- [x] Email templates gemaakt
- [ ] Migratie getest lokaal
- [ ] Endpoints getest met curl
- [ ] Frontend implementatie gestart

### Staging
- [ ] Database migratie uitgevoerd
- [ ] Backend deployed
- [ ] Endpoints functioneel
- [ ] Emails worden verzonden
- [ ] Frontend deployed
- [ ] End-to-end test

### Production
- [ ] Final QA approval
- [ ] Database backup gemaakt
- [ ] Migratie uitgevoerd
- [ ] Deployment verified
- [ ] Monitoring enabled
- [ ] Support team  geïnformeerd

---

**Versie:** 1.0  
**Laatste Update:** 2025-11-10  
**Status:** Backend Ready ✅ | Frontend Pending ⏳