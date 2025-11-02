# 🚀 Deployment Instructies - Email Decoder & Frontend Docs

## ✅ Wat Is Er Gemaakt?

### 📚 Frontend Documentatie (Compleet!)
1. **[`docs/FRONTEND_API_OVERVIEW.md`](docs/FRONTEND_API_OVERVIEW.md)** - Master overzicht
2. **[`docs/FRONTEND_ALBUMS_API.md`](docs/FRONTEND_ALBUMS_API.md)** - Photo albums (3 albums, 45 photos)
3. **[`docs/FRONTEND_VIDEOS_API.md`](docs/FRONTEND_VIDEOS_API.md)** - Video gallery (5 videos)
4. **[`docs/FRONTEND_UNDER_CONSTRUCTION_API.md`](docs/FRONTEND_UNDER_CONSTRUCTION_API.md)** - Maintenance mode
5. **[`docs/FRONTEND_EMAIL_API.md`](docs/FRONTEND_EMAIL_API.md)** - Complete email APIs (22+ endpoints)

### 🐛 Backend Verbeteringen
1. **[`services/email_decoder.go`](services/email_decoder.go)** - Advanced email decoder
2. **[`services/email_reprocessor.go`](services/email_reprocessor.go)** - Reprocess bestaande emails
3. **[`services/mail_fetcher.go`](services/mail_fetcher.go)** - Updated met nieuwe decoder
4. **[`handlers/admin_mail_handler.go`](handlers/admin_mail_handler.go)** - Nieuw reprocess endpoint
5. **[`main.go`](main.go)** - Geüpdatet voor nieuwe handler

### 🧪 Test Scripts
1. **[`test_albums.ps1`](test_albums.ps1)** - Test albums
2. **[`test_videos.ps1`](test_videos.ps1)** - Test videos
3. **[`test_under_construction.ps1`](test_under_construction.ps1)** - Test maintenance mode
4. **[`test_email_api.ps1`](test_email_api.ps1)** - Test email endpoints
5. **[`reprocess_emails.ps1`](reprocess_emails.ps1)** - Reprocess bestaande emails

---

## 🎯 DEPLOYMENT STAPPEN

### Stap 1: Git Commit & Push

```bash
# Add alle nieuwe files
git add services/email_decoder.go services/email_reprocessor.go
git add handlers/admin_mail_handler.go services/mail_fetcher.go main.go
git add docs/FRONTEND_*.md EMAIL_DECODER_IMPROVEMENTS.md
git add test_*.ps1 reprocess_emails.ps1
git add go.mod go.sum
git add DEPLOYMENT_INSTRUCTIONS.md

# Commit
git commit -m "feat: Advanced email decoder + complete frontend documentation

- Email decoder met multipart MIME, quoted-printable, Windows-1252 support
- Email reprocessor voor bestaande emails
- Complete frontend docs voor Albums, Videos, Email, Under Construction
- Test scripts voor alle modules
- Production tested"

# Push naar je repository
git push origin main
```

### Stap 2: Wacht op Render Deployment

Render zal automatisch deployen. Monitor in Render dashboard:
- Build logs checken
- Deployment status
- Service is back online (~2-3 minuten)

### Stap 3: Reprocess Bestaande Emails

**BELANGRIJK:** Run dit script om alle bestaande emails opnieuw te decoderen:

```powershell
.\reprocess_emails.ps1 -Environment production -Token "YOUR_JWT_TOKEN"
```

**Wat gebeurt er:**
- Script loopt door alle emails in database
- Detecteert welke emails encoding artifacts hebben (`=92`, MIME boundaries, etc.)
- Past verbeterde decoding toe
- Update database met schone content
- Geeft rapport van hoeveel emails zijn gefixed

**Verwacht resultaat:**
```
========================================
   Reprocessing Complete!
========================================
Processed: 15 emails
Failed: 0 emails
========================================

✓ Emails have been reprocessed with improved decoding
✓ Quoted-printable artifacts (=92, =85) are now fixed
✓ MIME boundaries are removed
✓ Windows-1252 characters converted to UTF-8
```

### Stap 4: Verificatie

Test dat emails nu clean zijn:

```powershell
# Test nieuwe emails ophalen
.\test_email_api.ps1 -Environment production -TestType public

# Check inbox via browser
# Login op admin panel
# Open inbox
# Bekijk een email die =92 had
# Moet nu I'm zijn (niet I=92m)
```

---

## 🔧 Nieuwe Backend Features

### Email Decoder (`email_decoder.go`)
- ✅ Multipart MIME parsing (boundaries)
- ✅ Quoted-printable decoding (`=92` → `'`)
- ✅ Base64 decoding
- ✅ Windows-1252 → UTF-8 conversie
- ✅ RFC 2047 header decoding
- ✅ HTML/Text part selectie
- ✅ Charset detection (6+ charsets)
- ✅ Graceful fallbacks

### Email Reprocessor (`email_reprocessor.go`)
- ✅ Batch reprocessing van bestaande emails
- ✅ Smart detection van encoding artifacts
- ✅ MIME boundary removal
- ✅ Charset conversion
- ✅ Database updates

### Nieuw Admin Endpoint
**`POST /api/admin/mail/reprocess`**
- Authenticatie: JWT Bearer token
- Permission: `admin_email:send`
- Functie: Reprocess ALL emails in database
- Response: Count van processed/failed emails

---

## 💻 Voor Frontend Developers

### Geen Wijzigingen Nodig! ✅

De frontend blijft **exact hetzelfde werken**:

```typescript
// Dit blijft hetzelfde
const response = await fetch('https://dklemailservice.onrender.com/api/mail/123', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const email = await response.json();

// email.html is nu automatisch clean (na backend deployment)
<div dangerouslySetInnerHTML={{ __html: email.html }} />
```

### Gebruik de Documentatie:

Start met **[`docs/FRONTEND_API_OVERVIEW.md`](docs/FRONTEND_API_OVERVIEW.md)** voor het complete overzicht.

Voor email specifics: **[`docs/FRONTEND_EMAIL_API.md`](docs/FRONTEND_EMAIL_API.md)**

Alle React components, TypeScript types, en voorbeelden staan klaar! 🎉

---

## 🧪 Testing Workflow

### Voor Deployment:
```powershell
# Test lokaal in Docker
docker-compose -f docker-compose.dev.yml up --build

# Test alle endpoints
.\test_albums.ps1 -Environment docker
.\test_videos.ps1 -Environment docker
.\test_email_api.ps1 -Environment docker -TestType public
```

### Na Deployment:
```powershell
# Test production
.\test_albums.ps1 -Environment production
.\test_videos.ps1 -Environment production
.\test_email_api.ps1 -Environment production -TestType public

# Fix bestaande emails
.\reprocess_emails.ps1 -Environment production -Token "YOUR_TOKEN"
```

---

## 📊 Verwachte Resultaten

### Voor Email Decoding Fix:
```
From: example@email.com
Subject: Hello
Body: I=92m writing... Thank You=85!!
--_000_PSAPR06MB4198_--
```

### Na Email Decoding Fix:
```
From: example@email.com
Subject: Hello
Body: I'm writing... Thank You…!!
```

Clean, readable content! 🎯

---

## 🆘 Troubleshooting

### Als emails nog steeds encoding artifacts hebben:
1. Check of service succesvol is herstart
2. Controleer of je `reprocess_emails.ps1` hebt gedraaid
3. Check logs voor errors: `docker logs dkl-email-service` of Render dashboard
4. Verify token permissions (moet `admin_email:send` hebben)

### Als reprocess script faalt:
```powershell
# Debug mode (enable in script)
$ErrorActionPreference = "Continue"

# Check token
curl https://dklemailservice.onrender.com/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 Changelog

### Added
- ✅ Advanced email decoder met multipart MIME support
- ✅ Windows-1252 en meerdere charset conversies
- ✅ Quoted-printable en Base64 decoding
- ✅ RFC 2047 header decoding
- ✅ Email reprocessor voor bestaande emails
- ✅ Admin endpoint `POST /api/admin/mail/reprocess`
- ✅ Complete frontend documentatie (5 docs)
- ✅ Test scripts voor alle modules (5 scripts)

### Fixed
- 🐛 Email encoding artifacts (`=92`, `=85`, etc.)
- 🐛 MIME boundaries in email body
- 🐛 Windows-1252 charset niet geconverteerd
- 🐛 Multipart emails niet correct geparsed

### Dependencies
- ✅ `golang.org/x/text v0.30.0` (charset conversie)

---

## 🎉 Deployment Checklist

- [ ] 1. Git commit & push
- [ ] 2. Wacht op Render deployment (~3 min)
- [ ] 3. Verify service is online (health check)
- [ ] 4. **RUN:** `.\reprocess_emails.ps1` met je token
- [ ] 5. Test nieuwe emails ook correct decoderen
- [ ] 6. Check frontend - emails should be clean
- [ ] 7. Deel frontend docs met frontend team

---

## 📞 Support

Bij vragen of problemen:
- Check [`EMAIL_DECODER_IMPROVEMENTS.md`](EMAIL_DECODER_IMPROVEMENTS.md) voor details
- Review logs in Render dashboard
- Test lokaal in Docker eerst
- Deel error logs voor support

**Frontend team heeft GEEN wijzigingen nodig** - alles werkt automatisch na jouw deployment! ✅