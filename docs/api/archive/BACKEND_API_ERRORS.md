# Backend API Errors - Analyse & Oplossingen

**Document Versie:** 1.1
**Datum:** 2025-01-08
**Status:** ✅ Geanalyseerd & Gefixed

---

## 📋 Executive Summary

Van de 2 gerapporteerde backend API errors:
- **1 echte bug gevonden & GEFIXED** ✅ (Participant Emails - route ordering probleem)
- **1 verwachte behavior** (Under Construction - geen actieve record)

> **UPDATE:** De participant emails bug is opgelost. Zie [BACKEND_FIX_PARTICIPANT_EMAILS.md](../BACKEND_FIX_PARTICIPANT_EMAILS.md) voor details.

---

## 1️⃣ Participant Emails Endpoint - 500 Error

### ✅ Status: **GEFIXED** - Route ordering gecorrigeerd

> **FIX TOEGEPAST:** 2025-01-08
> **Zie:** [BACKEND_FIX_PARTICIPANT_EMAILS.md](../BACKEND_FIX_PARTICIPANT_EMAILS.md)

### Probleem Beschrijving

**Error:**
```
GET http://localhost:8080/api/participant/emails 500 (Internal Server Error)
```

**Locatie:**
- Handler: [`handlers/participant_handler.go:408-451`](../handlers/participant_handler.go:408)
- Route Registratie: [`handlers/participant_handler.go:71-73`](../handlers/participant_handler.go:71)

### Root Cause Analyse

Het probleem zit in de **route registratie volgorde** in `participant_handler.go`:

```go
// ❌ PROBLEMATISCHE VOLGORDE
func (h *ParticipantHandler) RegisterRoutes(app *fiber.App) {
    participantApi := app.Group("/api/participant", AuthMiddleware(h.authService))
    
    // Regel 54-56: Generieke lijst route
    participantApi.Get("/",
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.ListParticipants)
    
    // Regel 58-60: ⚠️ Deze route matched ALLE paden inclusief "/emails"!
    participantApi.Get("/:id",
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.GetParticipant)
    
    // Regel 71-73: Deze route wordt NOOIT bereikt!
    participantApi.Get("/emails",
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.GetParticipantEmails)
}
```

**Wat gebeurt er:**
1. Frontend roept `/api/participant/emails` aan
2. Fiber matcht dit eerst met route `/:id` (waar `id = "emails"`)
3. De `GetParticipant` handler probeert een participant te vinden met ID "emails"
4. Dit faalt → 500 Internal Server Error of 404 Not Found

### ✅ Oplossing

**Route volgorde MOET specifiek → generiek zijn:**

```go
// ✅ CORRECTE VOLGORDE
func (h *ParticipantHandler) RegisterRoutes(app *fiber.App) {
    participantApi := app.Group("/api/participant", AuthMiddleware(h.authService))
    
    // 1. Eerst generieke lijst route
    participantApi.Get("/",
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.ListParticipants)
    
    // 2. Dan specifieke routes (VOOR :id route!)
    participantApi.Get("/emails",  // ← MOET VOOR /:id komen!
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.GetParticipantEmails)
    
    // 3. Tot slot generieke :id routes
    participantApi.Get("/:id",
        PermissionMiddleware(h.permissionService, "participant", "read"),
        h.GetParticipant)
    
    participantApi.Post("/:id/antwoord",
        PermissionMiddleware(h.permissionService, "participant", "write"),
        h.AddParticipantAntwoord)
    
    participantApi.Delete("/:id",
        PermissionMiddleware(h.permissionService, "participant", "delete"),
        h.DeleteParticipant)
}
```

### Implementatie Details

De handler `GetParticipantEmails` zelf is correct geïmplementeerd:

**Handler Code:** [`handlers/participant_handler.go:408-451`](../handlers/participant_handler.go:408)

```go
func (h *ParticipantHandler) GetParticipantEmails(c *fiber.Ctx) error {
    ctx := c.Context()
    
    // Haal ALLE participants op (max 10000)
    participants, err := h.participantRepo.List(ctx, 10000, 0)
    if err != nil {
        logger.Error("Fout bij ophalen participants voor emails", "error", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Kon participant emails niet ophalen",
            "code":  "INTERNAL_ERROR",
        })
    }
    
    // Extract participant emails
    participantEmails := make([]string, len(participants))
    for i, participant := range participants {
        participantEmails[i] = participant.Email
    }
    
    // Haal systeem emails op uit environment variables
    systemEmails := []string{}
    if infoEmail := os.Getenv("INFO_EMAIL"); infoEmail != "" {
        systemEmails = append(systemEmails, infoEmail)
    }
    if inschrijvingEmail := os.Getenv("INSCHRIJVING_EMAIL"); inschrijvingEmail != "" {
        systemEmails = append(systemEmails, inschrijvingEmail)
    }
    
    // Combineer alle emails
    allEmails := append(participantEmails, systemEmails...)
    
    // Return uitgebreide response
    return c.JSON(fiber.Map{
        "participant_emails": participantEmails,
        "system_emails":      systemEmails,
        "all_emails":         allEmails,
        "counts": fiber.Map{
            "participants": len(participantEmails),
            "system":       len(systemEmails),
            "total":        len(allEmails),
        },
    })
}
```

**Response Format:**
```json
{
  "participant_emails": ["email1@example.com", "email2@example.com"],
  "system_emails": ["info@dekoninklijkeloop.nl", "inschrijving@dekoninklijkeloop.nl"],
  "all_emails": ["email1@example.com", "email2@example.com", "info@...", "inschrijving@..."],
  "counts": {
    "participants": 2,
    "system": 2,
    "total": 4
  }
}
```

### Impact

**Getroffen Functionaliteit:**
1. **Email suggesties in admin panel** - Email autocomplete werkt niet
2. **Aanmelding email fallback** - Fallback mechanisme in `adminEmailService.ts:625` faalt
3. **Admin email compose** - Kan participant emails niet ophalen voor TO veld

**Severity:** 🔴 **HIGH** - Admin functionaliteit is gebroken

---

## 2️⃣ Under Construction Endpoint - 404 Error

### 🟢 Status: **GEEN BUG** - Dit is correct gedrag

### Probleem Beschrijving

**Error:**
```
GET http://localhost:8080/api/under-construction/active 404 (Not Found)
```

**Locatie:**
- Handler: [`handlers/under_construction_handler.go:67-93`](../handlers/under_construction_handler.go:67)
- Route Registratie: [`handlers/under_construction_handler.go:37`](../handlers/under_construction_handler.go:37)
- Main Registration: [`main.go:662-667`](../main.go:662)

### Root Cause Analyse

Dit is **GEEN BUG** maar **correct gedrag by design**!

**Handler Implementatie:**

```go
func (h *UnderConstructionHandler) GetActiveUnderConstruction(c *fiber.Ctx) error {
    ctx := c.Context()
    uc, err := h.underConstructionRepo.GetActive(ctx)
    
    if err != nil {
        // Check if this is a "record not found" error
        if err.Error() == "record not found" {
            // ✅ Return 404 WITHOUT logging error - this is EXPECTED
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "No active under construction found",
            })
        }
        // For other errors (database issues), log and return 500
        logger.Error("Failed to fetch active under construction", "error", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch under construction",
        })
    }
    
    if uc == nil {
        // ✅ Return 404 WITHOUT logging error - this is EXPECTED
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "No active under construction found",
        })
    }
    
    return c.JSON(uc)
}
```

**Wat gebeurt er:**
1. De website is **NIET** in maintenance mode
2. Er is **GEEN** actieve "under construction" record in database
3. Backend returnt **404** met message "No active under construction found"
4. Dit is **EXPECTED BEHAVIOR** - geen error log!

### ✅ Oplossing: Frontend Aanpassing

De frontend moet deze 404 **gracefully handlen**:

**Huidige Frontend Code** (waarschijnlijk):
```typescript
// ❌ SLECHT - Logt elke 404 als error
const response = await fetch('/api/under-construction/active');
if (!response.ok) {
  console.error('Failed to fetch under construction'); // ← Onnodig!
}
```

**Correcte Frontend Code:**
```typescript
// ✅ GOED - Handle 404 als normale state
const response = await fetch('/api/under-construction/active');

if (response.status === 404) {
  // Normal: website is not in maintenance mode
  return { isActive: false, mode: null };
}

if (!response.ok) {
  // Only log non-404 errors
  console.error('Server error fetching under construction:', response.status);
  return { isActive: false, mode: null };
}

const data = await response.json();
return { isActive: true, mode: data };
```

### Route Registratie

**Code:** [`handlers/under_construction_handler.go:33-56`](../handlers/under_construction_handler.go:33)

```go
func (h *UnderConstructionHandler) RegisterRoutes(app *fiber.App) {
    // Public routes (NO authentication required)
    public := app.Group("/api/under-construction")
    public.Get("/active", h.GetActiveUnderConstruction)
    public.Get("/", h.GetActiveUnderConstruction) // Alias for backwards compatibility
    
    // Admin routes (require authentication and permissions)
    admin := app.Group("/api/under-construction", AuthMiddleware(h.authService))
    
    // Read routes
    readGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "read"))
    readGroup.Get("/admin", h.ListUnderConstruction)
    readGroup.Get("/:id", h.GetUnderConstruction)
    
    // Write routes
    writeGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "write"))
    writeGroup.Post("/", h.CreateUnderConstruction)
    writeGroup.Put("/:id", h.UpdateUnderConstruction)
    
    // Delete routes
    deleteGroup := admin.Group("", PermissionMiddleware(h.permissionService, "under_construction", "delete"))
    deleteGroup.Delete("/:id", h.DeleteUnderConstruction)
}
```

**Registratie in main.go:** [`main.go:662-667`](../main.go:662)

```go
underConstructionHandler := handlers.NewUnderConstructionHandler(
    repoFactory.UnderConstruction,
    serviceFactory.AuthService,
    serviceFactory.PermissionService,
)
underConstructionHandler.RegisterRoutes(app)
```

✅ Routes zijn correct geregistreerd! Endpoint bestaat en werkt perfect.

### Impact

**Getroffen Functionaliteit:**
- Settings pagina toont "error" message (misleidend)
- Feature is optioneel dus geen kritische impact

**Severity:** 🟡 **LOW** - Frontend moet betere error handling toevoegen

---

## 🔧 Actie Items

### Voor Backend Developer

#### Priority 1: Fix Participant Emails Route
- [ ] Open [`handlers/participant_handler.go`](../handlers/participant_handler.go)
- [ ] Verplaats regel 71-73 (`Get("/emails")`) naar BOVEN regel 58 (`Get("/:id")`)
- [ ] Test endpoint: `curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/participant/emails`
- [ ] Verify response format matches expected structure
- [ ] Commit met message: `fix: correct route ordering for /api/participant/emails endpoint`

#### Priority 2: Verbeter Logging (optioneel)
- [ ] Voeg debug logging toe in `participant_handler.go:GetParticipant` om route matching te debuggen

### Voor Frontend Developer

#### Priority 1: Fix Under Construction Error Handling
- [ ] Update `underConstructionClient.ts` om 404 als normale state te handlen
- [ ] Verwijder error logging voor 404 responses
- [ ] Test dat settings pagina geen error toont wanneer maintenance mode uit staat

#### Priority 2: Update Email Service na Backend Fix
- [ ] Test `adminEmailService.getParticipantEmails()` na backend fix
- [ ] Verify email suggesties werken in admin panel
- [ ] Test fallback mechanisme in `fetchAanmeldingenEmails()`

---

## 📊 Testing Checklist

### Backend Tests

```bash
# Test 1: Participant Emails (na fix)
curl -X GET http://localhost:8080/api/participant/emails \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"

# Expected: 200 OK met participant_emails array

# Test 2: Under Construction Active (geen maintenance)
curl -X GET http://localhost:8080/api/under-construction/active \
  -H "Content-Type: application/json"

# Expected: 404 Not Found met "No active under construction found"

# Test 3: Under Construction Active (met maintenance)
# Eerst in database een actieve record aanmaken, dan:
curl -X GET http://localhost:8080/api/under-construction/active \
  -H "Content-Type: application/json"

# Expected: 200 OK met under construction object
```

### Frontend Tests

```typescript
// Test 1: Verify email suggesties werken
// Navigeer naar admin panel → Compose Email
// Typ in TO field → Verwacht autocomplete met participant emails

// Test 2: Verify under construction check niet logt
// Open browser console
// Refresh settings page
// Verwacht: GEEN error logs voor under-construction endpoint

// Test 3: Verify aanmelding fallback werkt
// Open admin panel → Email Inbox
// Klik "Load Aanmeldingen"
// Verwacht: Emails laden zonder errors
```

---

## 📚 Related Documentation

- [Participant Handler Implementation](../handlers/participant_handler.go)
- [Under Construction Handler Implementation](../handlers/under_construction_handler.go)
- [Main Routes Registration](../main.go)
- [Frontend Email Service](../../frontend/src/features/email/adminEmailService.ts) (indien beschikbaar)
- [API Quick Reference](./QUICK_REFERENCE.md)

---

## 📝 Change Log

| Datum | Versie | Auteur | Wijziging |
|-------|--------|--------|-----------|
| 2025-01-08 | 1.0 | Kilo Code | Initiële analyse en documentatie |
