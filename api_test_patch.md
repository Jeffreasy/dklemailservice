# API Testmodus Patch

Om de API te laten werken met ons testscript, kun je de volgende wijzigingen aanbrengen aan de `handlers/email_handler.go` file. Deze patch voegt ondersteuning toe voor een testmodus, waarbij geen echte emails worden verstuurd.

```go
// In de HandleContactEmail functie, voeg dit toe net na het parsen van het request
// rond regel 33
var testMode bool
if testModeValue, ok := c.GetReqHeaders()["X-Test-Mode"]; ok && testModeValue == "true" {
    testMode = true
    logger.Info("Test modus gedetecteerd via header", "remote_ip", c.IP())
}

if c.Locals("test_mode") != nil {
    testMode = true
    logger.Info("Test modus gedetecteerd via locals", "remote_ip", c.IP())
}

var requestMap map[string]interface{}
if err := json.Unmarshal(c.Body(), &requestMap); err == nil {
    if val, ok := requestMap["test_mode"]; ok && val.(bool) {
        testMode = true
        logger.Info("Test modus gedetecteerd via body parameter", "remote_ip", c.IP())
    }
}

// Nu rond regel 80, voor het verzenden van de admin email
if testMode {
    logger.Info("Test modus: Geen admin email verzonden", "admin_email", adminEmail)
} else {
    logger.Info("Admin email wordt verzonden", "admin_email", adminEmail, "contact_naam", request.Naam)
    if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
        logger.Error("Fout bij verzenden admin email", "error", err, "admin_email", adminEmail, "elapsed", time.Since(start))
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Fout bij het verzenden van de email: " + err.Error(),
        })
    }
    logger.Info("Admin email verzonden", "admin_email", adminEmail, "elapsed", time.Since(start))
}

// En rond regel 95, voor het verzenden van de gebruiker email
if testMode {
    logger.Info("Test modus: Geen gebruiker email verzonden", "user_email", request.Email)
} else {
    logger.Info("Bevestigingsemail wordt verzonden", "user_email", request.Email, "naam", request.Naam)
    if err := h.emailService.SendContactEmail(userEmailData); err != nil {
        logger.Error("Fout bij verzenden bevestigingsemail", "error", err, "user_email", request.Email, "elapsed", time.Since(start))
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Fout bij het verzenden van de bevestigingsemail: " + err.Error(),
        })
    }
    logger.Info("Bevestigingsemail verzonden", "user_email", request.Email, "elapsed", time.Since(start))
}

// Pas de return waarde aan op regel 111 om test modus aan te geven
if testMode {
    logger.Info("Contact formulier succesvol verwerkt in test modus", "naam", request.Naam, "email", request.Email, "total_elapsed", time.Since(start))
    return c.JSON(fiber.Map{
        "success": true,
        "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
        "test_mode": true,
    })
} else {
    logger.Info("Contact formulier succesvol verwerkt", "naam", request.Naam, "email", request.Email, "total_elapsed", time.Since(start))
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Je bericht is verzonden! Je ontvangt ook een bevestiging per email.",
    })
}
```

Doe hetzelfde voor de `HandleAanmeldingEmail` functie met vergelijkbare aanpassingen.

## Implementatie in code

Om deze aanpassingen door te voeren, kun je:

1. De `email_handler.go` file aanpassen met de bovenstaande wijzigingen
2. Hercompileren en opnieuw deployen van de API

Of als alternatief, kun je een middleware toevoegen die testmodus detecteert. Dit kun je toevoegen aan `handlers/middleware.go`:

```go
// TestModeMiddleware controleert of de request in testmodus moet worden uitgevoerd
func TestModeMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Controleer op test mode header
        if testMode := c.Get("X-Test-Mode"); testMode == "true" {
            c.Locals("test_mode", true)
            logger.Debug("Test modus geactiveerd via header", "path", c.Path(), "ip", c.IP())
        }
        
        // Ga verder met de request
        return c.Next()
    }
}
```

En voeg deze middleware toe in `main.go`:

```go
// Voeg deze regel toe waar de andere middleware wordt toegevoegd
app.Use(handlers.TestModeMiddleware())
```

Dit biedt een elegante manier om testmodus te ondersteunen zonder de handlers zelf te hoeven aanpassen.

## Uitleg

Deze aanpassingen zorgen ervoor dat:

1. De API testmodus detecteert via een header (`X-Test-Mode: true`) of een request parameter (`test_mode: true`)
2. In testmodus geen echte emails verstuurt
3. Een speciale respons teruggeeft die aangeeft dat de request succesvol was verwerkt in testmodus

Voordelen van deze aanpak:
- Je kunt de API testen zonder echte emails te versturen
- Het vermindert onnodige belasting van de SMTP server
- Het voorkomt dat testberichten naar echte gebruikers worden gestuurd
- Het maakt geautomatiseerd testen mogelijk 