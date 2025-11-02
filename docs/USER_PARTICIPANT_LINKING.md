# Gebruikers vs Participants - Koppeling Uitleg

Dit document legt uit hoe de koppeling tussen `gebruikers` (users) en `aanmeldingen` (participants) werkt, vooral belangrijk voor gamification features.

## Overzicht

Er zijn **twee gescheiden maar gekoppelde tabellen**:

1. **`gebruikers`** - Voor authenticatie en RBAC (admin/staff/deelnemer rollen)
2. **`aanmeldingen`** - Voor event participatie (deelnemers die meedoen aan het lopen)

## Relatie Schema

```
┌─────────────────┐              ┌──────────────────┐
│   gebruikers    │              │  aanmeldingen    │
├─────────────────┤              ├──────────────────┤
│ id (PK)        │◄─────┐       │ id (PK)          │
│ naam            │      │       │ naam             │
│ email           │      └───────┤ gebruiker_id (FK)│
│ wachtwoord_hash │              │ steps            │
│ is_actief       │              │ afstand          │
│ roles (M2M)     │              │ ...              │
└─────────────────┘              └──────────────────┘
                                          │
                                          │ participant_id (FK)
                                          ▼
                            ┌──────────────────────────┐
                            │ participant_achievements │
                            ├──────────────────────────┤
                            │ id                       │
                            │ participant_id           │
                            │ badge_id                 │
                            │ earned_at                │
                            └──────────────────────────┘
```

## Vier Scenario's

### 1. Participant ZONDER User Account
**Bijvoorbeeld:** Iemand die zich aanmeldt maar geen login account heeft.

```sql
-- In aanmeldingen:
INSERT INTO aanmeldingen (naam, email, afstand, gebruiker_id)
VALUES ('Jan Jansen', 'jan@example.com', '10 KM', NULL);
```

- ✅ Kan deelnemen aan event
- ✅ Kan stappen loggen (via admin/staff)
- ✅ Kan badges verdienen
- ✅ Zichtbaar op leaderboard
- ❌ Kan NIET inloggen
- ❌ Kan NIET zelf stappen updaten

### 2. Participant MET User Account (Deelnemer Role)
**Bijvoorbeeld:** Iemand die zich aanmeldt EN een login account krijgt.

```sql
-- Eerst gebruiker aanmaken:
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash)
VALUES (gen_random_uuid(), 'Piet Pietersen', 'piet@example.com', '...');

-- Deelnemer role toekennen:
INSERT INTO user_roles (user_id, role_id)
SELECT g.id, r.id FROM gebruikers g, roles r 
WHERE g.email = 'piet@example.com' AND r.name = 'deelnemer';

-- Dan aanmelding koppelen:
INSERT INTO aanmeldingen (naam, email, afstand, gebruiker_id)
VALUES ('Piet Pietersen', 'piet@example.com', '15 KM', 
        (SELECT id FROM gebruikers WHERE email = 'piet@example.com'));
```

- ✅ Kan deelnemen aan event
- ✅ Kan inloggen
- ✅ Kan ZELF stappen updaten
- ✅ Kan eigen dashboard bekijken
- ✅ Kan badges verdienen
- ✅ Zichtbaar op leaderboard

### 3. Admin/Staff die ALLEEN Beheren
**Bijvoorbeeld:** Een organisator die niet mee loopt.

```sql
-- Alleen gebruiker, geen aanmelding:
INSERT INTO gebruikers (naam, email, wachtwoord_hash)
VALUES ('Admin User', 'admin@example.com', '...');

-- Admin role toekennen:
INSERT INTO user_roles (user_id, role_id)
SELECT g.id, r.id FROM gebruikers g, roles r 
WHERE g.email = 'admin@example.com' AND r.name = 'admin';
```

- ✅ Kan inloggen
- ✅ Kan andere participants beheren
- ✅ Kan badges aanmaken/beheren
- ✅ Kan badges toekennen aan anderen
- ❌ Heeft GEEN eigen stappen
- ❌ Verschijnt NIET op leaderboard
- ❌ Kan GEEN badges verdienen voor zichzelf

### 4. Admin/Staff die OOK Deelneemt
**Bijvoorbeeld:** Een organisator die ook mee loopt.

```sql
-- Gebruiker bestaat al als admin:
-- Voeg nu een aanmelding toe:
INSERT INTO aanmeldingen (naam, email, afstand, gebruiker_id)
VALUES ('Admin User', 'admin@example.com', '20 KM',
        (SELECT id FROM gebruikers WHERE email = 'admin@example.com'));
```

- ✅ Kan inloggen (via gebruiker)
- ✅ Kan andere participants beheren (via admin role)
- ✅ Kan badges aanmaken/beheren (via admin role)
- ✅ Heeft OOK een participant profiel
- ✅ Kan eigen stappen loggen
- ✅ Verschijnt op leaderboard
- ✅ Kan badges verdienen

## Code Voorbeelden

### Handler: Stappen Updaten voor Ingelogde User

```go
func (h *StepsHandler) UpdateSteps(c *fiber.Ctx) error {
    // Haal userID uit JWT token
    userID := c.Locals("userID").(string)
    
    // Zoek participant_id voor deze user
    var participant models.Aanmelding
    err := db.Where("gebruiker_id = ?", userID).First(&participant).Error
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "U bent niet geregistreerd als deelnemer",
        })
    }
    
    // Update stappen voor participant
    participant.Steps += deltaSteps
    db.Save(&participant)
    
    return c.JSON(participant)
}
```

### Handler: Achievements voor Ingelogde User

```go
func (h *GamificationHandler) GetMyAchievements(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    
    // Zoek participant_id
    var participant models.Aanmelding
    err := db.Where("gebruiker_id = ?", userID).First(&participant).Error
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "U bent niet geregistreerd als deelnemer",
        })
    }
    
    // Haal achievements op via participant_id
    achievements := gamificationService.GetParticipantAchievements(
        ctx, 
        participant.ID,
    )
    
    return c.JSON(achievements)
}
```

## Helper Functies (Database)

De migratie `V1_52` voegt handige functies toe:

### get_participant_id_for_user(user_id)
```sql
-- Voorbeeld gebruik:
SELECT get_participant_id_for_user('user-uuid-here');

-- Retourneert participant_id als user een participant is
-- Retourneert NULL als user geen participant is
```

### get_user_id_for_participant(participant_id)
```sql
-- Voorbeeld gebruik:
SELECT get_user_id_for_participant('participant-uuid-here');

-- Retourneert user_id als participant een account heeft
-- Retourneert NULL als participant geen account heeft
```

## Helper Views (Database)

### user_participant_mapping
```sql
-- Toont alle users met hun participant info:
SELECT * FROM user_participant_mapping;

-- Resultaat:
-- gebruiker_id | gebruiker_naam | aanmelding_id | is_participant
-- uuid-1       | Admin User     | NULL          | false         (alleen admin)
-- uuid-2       | Piet Pietersen | uuid-a        | true          (ook participant)
```

### unified_leaderboard
```sql
-- Leaderboard met user info:
SELECT * FROM unified_leaderboard LIMIT 10;

-- Toont rank, steps, achievements, EN user account status
```

### users_without_participation
```sql
-- Vind users die kunnen deelnemen maar nog niet doen:
SELECT * FROM users_without_participation;

-- Handig voor admin UI: "Deze users kunnen zich nog aanmelden"
```

## Best Practices

### ✅ DO's

1. **Check altijd of user een participant is** voordat je participant-gerelateerde acties doet
2. **Gebruik de helper functies** om tussen user_id en participant_id te converteren
3. **Gebruik de views** voor UI die beide tabellen moet tonen
4. **Houd email addresses gesynchroniseerd** tussen gebruikers en aanmeldingen

### ❌ DON'Ts

1. **Neem niet aan dat elke user een participant is** - check eerst
2. **Gebruik niet de oude `gebruikers.rol` veld** - gebruik `user_roles` tabel
3. **Probeer niet badges toe te kennen aan user_id** - gebruik altijd participant_id

## Frontend Overwegingen

### Login Flow Check
```typescript
// Na login:
const response = await fetch('/api/participant/dashboard', {
    headers: { 'Authorization': `Bearer ${token}` }
});

if (response.status === 404) {
    // User is GEEN participant
    showMessage('Je bent nog niet aangemeld als deelnemer');
    showRegistrationButton();
} else {
    // User IS participant
    const dashboard = await response.json();
    showParticipantDashboard(dashboard);
}
```

### Admin die ook Participant is
```typescript
// Check roles:
const hasAdminRole = user.roles.includes('admin');
const isParticipant = await checkIfParticipant(user.id);

if (hasAdminRole && isParticipant) {
    // Toon beide interfaces:
    showAdminPanel();      // Badge management, etc.
    showParticipantDashboard();  // Eigen stappen, achievements
}
```

## Database Constraints

### Foreign Keys
```sql
-- aanmeldingen.gebruiker_id -> gebruikers.id
ALTER TABLE aanmeldingen 
ADD CONSTRAINT fk_aanmelding_gebruiker 
FOREIGN KEY (gebruiker_id) REFERENCES gebruikers(id) ON DELETE SET NULL;

-- participant_achievements.participant_id -> aanmeldingen.id
ALTER TABLE participant_achievements
ADD CONSTRAINT fk_participant
FOREIGN KEY (participant_id) REFERENCES aanmeldingen(id) ON DELETE CASCADE;
```

### Wat betekent dit?
- Als een **gebruiker** wordt verwijderd → `gebruiker_id` in aanmeldingen wordt NULL (aanmelding blijft bestaan)
- Als een **aanmelding** wordt verwijderd → alle achievements worden ook verwijderd (CASCADE)

## FAQs

### Q: Kan een participant meerdere user accounts hebben?
**A:** Nee, door de UNIQUE constraint op email in beide tabellen en de foreign key is dit niet mogelijk.

### Q: Kan een user meerdere participants zijn?
**A:** Technisch ja (geen UNIQUE constraint op gebruiker_id in aanmeldingen), maar dit is niet de bedoeling. Normaal heeft een user maximaal 1 aanmelding.

### Q: Hoe maak ik van een bestaande participant een user?
**A:**
```sql
-- 1. Maak gebruiker aan met zelfde email
INSERT INTO gebruikers (naam, email, wachtwoord_hash)
VALUES ('Naam', 'email@example.com', hash);

-- 2. Koppel aanmelding
UPDATE aanmeldingen 
SET gebruiker_id = (SELECT id FROM gebruikers WHERE email = 'email@example.com')
WHERE email = 'email@example.com';
```

### Q: Waarom niet alles in één tabel?
**A:** Scheiding van concerns:
- **Gebruikers tabel** = Authenticatie & Authorization (RBAC)
- **Aanmeldingen tabel** = Event participatie & gamification

Dit maakt het systeem flexibeler en schaalbaarder.

---

**Laatst bijgewerkt:** 2026-01-02  
**Migraties:** V1_51, V1_52