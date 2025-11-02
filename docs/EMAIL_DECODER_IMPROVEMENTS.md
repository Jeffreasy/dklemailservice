# Email Decoder Verbeteringen

## 🎯 Probleem
Inkomende emails werden niet correct gedecodeerd:
- `=92` characters (quoted-printable encoding)
- MIME boundaries niet verwerkt
- Multipart emails niet correct geparsed
- Windows-1252 charset niet geconverteerd naar UTF-8
- Headers niet gedecodeerd (RFC 2047)

## ✅ Oplossing

### Nieuwe Files
1. **[`services/email_decoder.go`](services/email_decoder.go)** - Complete email decoder service

### Geüpdatete Files
1. **[`services/mail_fetcher.go`](services/mail_fetcher.go)** - Gebruikt nu de nieuwe decoder

### Dependencies Toegevoegd
```
golang.org/x/text/encoding/charmap
golang.org/x/text/transform
```

---

## 🔧 Nieuwe Features

### 1. Multipart MIME Parsing
Parseert multipart/mixed en multipart/alternative emails correct:
```
--_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_
Content-Type: text/plain
Content-Type: text/html
```
→ Extraheert de juiste text of HTML part

### 2. Charset Conversie
Ondersteunt conversie van verschillende charsets naar UTF-8:
- **Windows-1252** (meest voorkomend voor Nederlandse emails)
- **ISO-8859-1** (Latin-1)
- **ISO-8859-15** (Latin-9)
- **Windows-1250** (Oost-Europees)
- **Windows-1251** (Cyrillisch)

Voorbeeld:
```
Content-Type: text/plain; charset="Windows-1252"
I=92m writing...
```
→ Wordt correct geconverteerd naar: `I'm writing...`

### 3. Transfer Encoding Decoding
Correct decoderen van:
- **Quoted-Printable** (`=92`, `=85`, etc.)
- **Base64**
- **7bit, 8bit, binary**

### 4. RFC 2047 Header Decoding
Decodeert encoded headers:
```
Subject: =?utf-8?Q?Hallo_=C3=A9=C3=A9n?=
From: =?windows-1252?Q?Jos=E9_Butler?= <email@example.com>
```
→ Correct gedecodeerd naar UTF-8

### 5. Plain Text naar HTML Conversie
Converteert plain text emails naar HTML met:
- HTML escaping (voorkomen XSS)
- Line breaks → `<br>` tags
- Monospace font voor formatting behoud

---

## 🧪 Wat Er Verbeterd Is

### Voor:
```
I=92m writing to follow up. Since I haven't received a reply...
Thank You=85!!
```

### Na:
```
I'm writing to follow up. Since I haven't received a reply...
Thank You…!!
```

### Voor (multipart):
```
--_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_
Content-Type: text/plain; charset="Windows-1252"
...complex MIME structure...
```

### Na:
```html
<div>Clean HTML content zonder MIME boundaries</div>
```

---

## 📋 Wat Moet Je Doen?

### 1. Service Herstarten

#### Docker (Development):
```bash
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up --build -d
```

#### Production (Render):
```bash
# Push naar repository
git add .
git commit -m "feat: Advanced email decoder with multipart, charset, and encoding support"
git push origin main

# Render zal automatisch deployen
```

### 2. Bestaande Emails Reprocessen (Belangrijk!)

**Na deployment**, gebruik het nieuwe admin endpoint om ALL bestaande emails opnieuw te decoderen:

#### Via PowerShell Script (Aanbevolen):
```powershell
.\reprocess_emails.ps1 -Environment production -Token "YOUR_JWT_TOKEN"
```

#### Of via curl:
```bash
curl -X POST https://dklemailservice.onrender.com/api/admin/mail/reprocess \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "success": true,
  "message": "Email reprocessing voltooid",
  "processed": 15,
  "failed": 0
}
```

Dit zal:
- ✅ Alle bestaande emails in de database vinden
- ✅ Detecteren welke emails encoding artifacts hebben
- ✅ Opnieuw decoderen met de verbeterde decoder
- ✅ Database updaten met schone content

### 3. Nieuwe Emails Testen

Test dat nieuwe emails ook correct worden gedecodeerd:

```powershell
# Fetch nieuwe emails
curl -X POST https://dklemailservice.onrender.com/api/mail/fetch \
  -H "Authorization: Bearer YOUR_TOKEN"

# Check of emails correct zijn gedecodeerd
curl https://dklemailservice.onrender.com/api/mail/unprocessed \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Verificatie

Check of emails nu correct worden weergegeven:
- Geen `=92` characters meer
- Geen `=85` characters meer
- Geen MIME boundaries in de body
- Correcte Nederlandse characters (é, ë, etc.)
- Clean HTML zonder raw encoding

---

## 🔍 Technische Details

### Email Decoder Architecture

```
Email Input
    ↓
Parse MIME Headers
    ↓
Is Multipart? ──Yes──→ Parse Each Part
    ↓                      ↓
    No              Decode Transfer Encoding
    ↓                      ↓
Decode Transfer       Convert Charset
Encoding                   ↓
    ↓              Choose Best Part (HTML > Text)
Convert Charset            ↓
    ↓                   ←──┘
Sanitize HTML (if needed)
    ↓
Clean Output
```

### Supported MIME Types
- `multipart/mixed` ✅
- `multipart/alternative` ✅
- `text/html` ✅
- `text/plain` ✅ (converted to HTML)

### Supported Encodings
- `quoted-printable` ✅
- `base64` ✅
- `7bit`, `8bit`, `binary` ✅

### Supported Charsets
- `UTF-8` ✅
- `Windows-1252` ✅ (meest voorkomend)
- `ISO-8859-1` ✅
- `ISO-8859-15` ✅
- `Windows-1250` ✅
- `Windows-1251` ✅

---

## 🐛 Known Edge Cases

De decoder handelt de volgende edge cases:
- ✅ Emails zonder explicit charset
- ✅ Mixed transfer encodings in multipart
- ✅ Invalid MIME boundaries (fallback to raw)
- ✅ Missing Content-Type headers
- ✅ Nested multipart messages
- ✅ Empty parts in multipart
- ✅ Malformed quoted-printable

---

## 📊 Performance Impact

- **Minimal overhead**: Alleen extra processing voor gecodeerde content
- **Memory efficient**: Streaming decoders gebruikt
- **Concurrent safe**: Decoder is stateless
- **Backward compatible**: Fallback naar oude methode bij errors

---

## 🔜 Mogelijke Toekomstige Verbeteringen

1. **Attachment extraction** - Parse en sla attachments op
2. **Inline images** - Extract en serve inline afbeeldingen
3. **More charsets** - UTF-16, Shift-JIS, etc.
4. **Email threading** - Group replies samen
5. **Spam detection** - Integreer spam filter
6. **Auto-categorization** - ML-based email categorisatie

---

## 🎓 Testing Tips

### Test Email Formats

1. **Quoted-Printable Test:**
```
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset="Windows-1252"

I=92m testing special characters: =E9 =EB =EF
```

2. **Base64 Test:**
```
Content-Transfer-Encoding: base64
Content-Type: text/plain

SGVsbG8gV29ybGQh
```

3. **Multipart Test:**
```
Content-Type: multipart/alternative; boundary="boundary123"

--boundary123
Content-Type: text/plain

Plain text version
--boundary123
Content-Type: text/html

<p>HTML version</p>
--boundary123--
```

### Verify in Database
```sql
SELECT 
  id,
  subject,
  "from",
  LEFT(body, 200) as body_preview,
  content_type
FROM incoming_emails
ORDER BY received_at DESC
LIMIT 5;
```

De body moet nu clean UTF-8 zijn zonder encoding artifacts!

---

## ✅ Checklist

- [x] EmailDecoder service gemaakt
- [x] mail_fetcher.go geüpdatet
- [x] Dependencies toegevoegd
- [x] Code compileert
- [ ] **Service herstarten** (Docker of Production)
- [ ] **Test met echte emails**
- [ ] **Verificeer decoding** in database
- [ ] **Check frontend display**

---

## 🆘 Troubleshooting

### Als emails nog steeds raar zijn:
1. Check logs voor decode errors
2. Verificeer Content-Type header van email
3. Check of charset bekend is
4. Test met verschillende email clients
5. Enable DEBUG logging voor meer details

### Logs Checken:
```bash
# Docker
docker-compose -f docker-compose.dev.yml logs -f app | grep -i "decode"

# Production (Render)
# Check Render dashboard logs
```

---

## 📧 Contact

Bij problemen met specifieke email formats, deel:
1. Content-Type header
2. Content-Transfer-Encoding header
3. Charset (uit Content-Type)
4. Sample van de raw email (gefilterd)

Dit helpt om de decoder verder te verbeteren!