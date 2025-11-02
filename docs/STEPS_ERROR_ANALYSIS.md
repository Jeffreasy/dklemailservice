# Steps Service Error Analysis & Fix - "deelnemer niet gevonden"

## 🔄 UPDATE: Admin Access Fixed (2025-11-02)

**Second Issue Discovered:** Admins konden geen stappen bijwerken voor participanten omdat de handler eerst checkte of de ingelogde gebruiker een participant was.

**Root Cause:** De logica prioriteerde `userID` (participant flow) boven `id` parameter (admin flow), waardoor admins altijd de "NOT_A_PARTICIPANT" error kregen.

**Solution Implemented:** Prioriteit omgedraaid in [`handlers/steps_handler.go`](../handlers/steps_handler.go):
- ✅ ID parameter (admin flow) heeft nu voorrang
- ✅ POST `/api/steps/:id` - admin werkt stappen bij voor specifieke participant
- ✅ POST `/api/steps` - participant werkt eigen stappen bij (geen ID nodig)
- ✅ Zelfde logica voor dashboard endpoints

---

# Steps Service Error Analysis - "deelnemer niet gevonden"

## 📋 Probleemanalyse

### Error Details
```json
{
  "niveau": "ERROR",
  "tijd": "2025-11-02T17:59:38.579Z",
  "caller": "logger/logger.go:212",
  "bericht": "Fout bij bijwerken stappen",
  "error": "deelnemer niet gevonden: record not found",
  "user_id": "0197cfc3-7ca2-403b-ae4d-32627cd47222"
}
```

## 🔍 Root Cause Analysis

### Situatie
De fout treedt op wanneer een **ingelogde gebruiker** probeert stappen bij te werken via de [`POST /api/steps`](../handlers/steps_handler.go:39) endpoint, maar deze gebruiker **geen bijbehorende aanmelding** heeft in de `aanmeldingen` tabel.

### Technische Details

#### Database Schema
```
gebruikers
├── id (UUID)           ← De user_id in de logging
├── naam
├── email
└── ...

aanmeldingen
├── id (UUID)
├── gebruiker_id (UUID, OPTIONAL)  ← Link naar gebruikers.id
├── naam
├── email
├── steps
└── ...
```

#### Flow van de Error
1. Gebruiker logt in → krijgt JWT token met `user_id`
2. Gebruiker roept [`POST /api/steps`](../handlers/steps_handler.go:39) aan
3. Handler haalt `userID` uit context: [`handlers/steps_handler.go:79`](../handlers/steps_handler.go:79)
4. Service zoekt aanmelding met `gebruiker_id = userID`: [`services/steps_service.go:58`](../services/steps_service.go:58)
5. **GORM retourneert** `record not found` → geen match gevonden
6. Error wordt gelogd met originele foutmelding

### Waarom Gebeurt Dit?

#### Scenario 1: Staff/Admin Accounts
- Staff en admin gebruikers hebben accounts in `gebruikers` tabel
- Deze accounts zijn **niet** bedoeld als deelnemers
- Ze hebben **geen** record in `aanmeldingen` tabel
- Ze mogen **geen** stappen bijwerken

#### Scenario 2: Gebruiker zonder Aanmelding
- Gebruiker heeft zich geregistreerd voor het systeem
- Gebruiker heeft zich **niet** (nog) aangemeld als deelnemer
- Geen link tussen `gebruikers.id` en `aanmeldingen.gebruiker_id`

#### Scenario 3: Data Inconsistentie
- Aanmelding is verwijderd maar gebruikersaccount bestaat nog
- Database migratie waarbij de link niet correct is ingesteld

## ✅ Geïmplementeerde Oplossing

### 1. Verbeterde Error Handling in Service Layer

**Bestand:** [`services/steps_service.go`](../services/steps_service.go)

#### Voor:
```go
func (s *StepsService) UpdateStepsByUserID(userID string, deltaSteps int) (*models.Aanmelding, error) {
    var participant models.Aanmelding
    err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
    if err != nil {
        return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
    }
    // ...
}
```

#### Na:
```go
func (s *StepsService) UpdateStepsByUserID(userID string, deltaSteps int) (*models.Aanmelding, error) {
    var participant models.Aanmelding
    err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("geen deelnemersregistratie gevonden voor gebruiker %s - gebruiker is mogelijk geen deelnemer", userID)
        }
        return nil, fmt.Errorf("fout bij ophalen deelnemer: %w", err)
    }
    // ...
}
```

**Changes:**
- ✅ Onderscheid tussen "niet gevonden" en andere database errors
- ✅ Duidelijke foutmelding die de situatie uitlegt
- ✅ Bevat user_id voor debugging
- ✅ Toegepast op beide functies: `UpdateStepsByUserID` en `GetParticipantDashboardByUserID`

### 2. Verbeterde HTTP Response in Handler Layer

**Bestand:** [`handlers/steps_handler.go`](../handlers/steps_handler.go)

#### Voor:
```go
if ok && userID != "" {
    participant, err := h.stepsService.UpdateStepsByUserID(userID, req.Steps)
    if err != nil {
        logger.Error("Fout bij bijwerken stappen", "error", err, "user_id", userID)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Kon stappen niet bijwerken",
        })
    }
    return c.JSON(participant)
}
```

#### Na:
```go
if ok && userID != "" {
    participant, err := h.stepsService.UpdateStepsByUserID(userID, req.Steps)
    if err != nil {
        logger.Error("Fout bij bijwerken stappen", "error", err, "user_id", userID)
        // Check if this is a "not found" error (user is not a participant)
        if strings.Contains(err.Error(), "geen deelnemersregistratie gevonden") {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "U bent niet geregistreerd als deelnemer. Alleen deelnemers kunnen stappen bijwerken.",
                "code":  "NOT_A_PARTICIPANT",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Kon stappen niet bijwerken",
            "code":  "INTERNAL_ERROR",
        })
    }
    return c.JSON(participant)
}
```

**Changes:**
- ✅ HTTP 404 (Not Found) i.p.v. HTTP 500 voor "niet gevonden" scenario
- ✅ Gebruiksvriendelijke foutmelding voor frontend
- ✅ Error code `NOT_A_PARTICIPANT` voor programmatische afhandeling
- ✅ HTTP 500 alleen voor echte internal server errors
- ✅ Toegepast op zowel `UpdateSteps` als `GetParticipantDashboard`

## 📊 Impact & Benefits

### Voor Gebruikers
- ✅ Duidelijke feedback: "U bent niet geregistreerd als deelnemer"
- ✅ Begrijpen waarom actie niet mogelijk is
- ✅ Frontend kan specifieke UI tonen (bijv. registratie prompt)

### Voor Developers
- ✅ Betere logging met contextuele informatie
- ✅ Onderscheid tussen verschillende error scenarios
- ✅ Makkelijker debuggen met user_id in error message

### Voor Monitoring
- ✅ HTTP 404 in plaats van HTTP 500 → geen false alarms
- ✅ Duidelijk verschil tussen gebruikersfouten en systeemfouten
- ✅ Betere metrics voor monitoring dashboards

## 🧪 Testing Scenarios

### Test 1: Admin Werkt Steps Bij Voor Participant (✅ WERKT NU)
```bash
# Login als admin
POST /api/auth/login
{
  "email": "admin@example.com",
  "wachtwoord": "password"
}

# Update steps voor specifieke participant via ID
POST /api/steps/<participant-aanmelding-id>
Authorization: Bearer <admin-token>
{
  "steps": 1000
}

# Verwacht resultaat:
HTTP 200 OK
{
  "id": "<participant-aanmelding-id>",
  "steps": 1000,
  "naam": "Participant Naam",
  ...
}
```

### Test 2: Participant Werkt Eigen Steps Bij
```bash
# Login als participant (heeft aanmelding)
POST /api/auth/login
{
  "email": "participant@example.com",
  "wachtwoord": "password"
}

# Update steps
POST /api/steps
Authorization: Bearer <token>
{
  "steps": 1000
}

# Verwacht resultaat:
HTTP 200 OK
{
  "id": "...",
  "steps": 1000,
  ...
}
```

### Test 3: Staff Probeert Eigen Steps Bij Te Werken (Zonder ID)
```bash
# Login als staff (geen participant)
POST /api/auth/login
{
  "email": "staff@example.com",
  "wachtwoord": "password"
}

# Probeer eigen steps bij te werken (GEEN ID parameter)
POST /api/steps
Authorization: Bearer <staff-token>
{
  "steps": 1000
}

# Verwacht resultaat:
HTTP 404 Not Found
{
  "error": "U bent niet geregistreerd als deelnemer. Alleen deelnemers kunnen hun eigen stappen bijwerken.",
  "code": "NOT_A_PARTICIPANT"
}
```

### Test 4: Admin Haalt Dashboard Op Voor Participant
```bash
GET /api/participant/<participant-aanmelding-id>/dashboard
Authorization: Bearer <admin-token>

# Verwacht resultaat:
HTTP 200 OK
{
  "steps": 1000,
  "route": "10 KM",
  "allocatedFunds": 75,
  "naam": "Participant Naam",
  "email": "participant@example.com"
}
```

## 🔧 Preventieve Maatregelen

### 1. Frontend Validatie
De frontend kan de error code checken:
```typescript
try {
  await api.post('/api/steps', { steps: 1000 });
} catch (error) {
  if (error.response?.data?.code === 'NOT_A_PARTICIPANT') {
    // Toon registratie prompt
    showRegistrationPrompt();
  } else {
    // Toon generieke error
    showError(error);
  }
}
```

### 2. Role-Based UI
Verberg steps-functionaliteit voor niet-deelnemers:
```typescript
// Check role in JWT token
if (user.roles.includes('participant')) {
  <StepsTracker />
} else {
  <RegisterAsParticipant />
}
```

### 3. Database Constraints
Overweeg een check constraint:
```sql
-- Optioneel: Zorg dat participant role een aanmelding heeft
ALTER TABLE gebruikers ADD CONSTRAINT check_participant_has_aanmelding
  CHECK (
    rol != 'participant' OR 
    EXISTS (SELECT 1 FROM aanmeldingen WHERE gebruiker_id = id)
  );
```

## 📝 Aanbevelingen

### Short Term
1. ✅ **GEDAAN:** Betere error messages
2. ✅ **GEDAAN:** HTTP status codes correct gebruiken
3. **TODO:** Frontend error handling implementeren
4. **TODO:** UI aanpassen voor niet-deelnemers

### Long Term
1. **TODO:** Role-based route filtering op API niveau
2. **TODO:** Periodic data integrity checks
3. **TODO:** Automated tests voor edge cases
4. **TODO:** Monitoring dashboard voor deze specifieke error

## 🔗 Related Files

- [`models/aanmelding.go`](../models/aanmelding.go) - Aanmelding model
- [`models/gebruiker.go`](../models/gebruiker.go) - Gebruiker model
- [`services/steps_service.go`](../services/steps_service.go) - Steps business logic
- [`handlers/steps_handler.go`](../handlers/steps_handler.go) - Steps HTTP handlers
- [`repository/aanmelding_repository.go`](../repository/aanmelding_repository.go) - Database queries

## 📞 Contact

Bij verdere vragen of problemen, neem contact op met het development team.

---
**Last Updated:** 2025-11-02 18:18 CET
**Author:** Development Team
**Status:** ✅ Fully Resolved
**Commits:**
- Initial fix: `20f510c` - Better error handling for non-participants
- Admin flow fix: `e0b1e0c` - Admin can now update participant steps via ID parameter