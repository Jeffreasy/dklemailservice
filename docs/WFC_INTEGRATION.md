# Whisky for Charity Email Integration

Dit document beschrijft hoe je de Whisky for Charity email configuratie kunt gebruiken om emails te versturen via de arg-plplcl14.argewebhosting.nl server.

## Configuratie

Voeg de volgende omgevingsvariabelen toe aan je `.env` bestand:

```
# SMTP configuratie voor Whisky for Charity
WFC_SMTP_HOST=arg-plplcl14.argewebhosting.nl
WFC_SMTP_PORT=465
WFC_SMTP_USER=noreply@whiskyforcharity.com
WFC_SMTP_PASSWORD=your_password_here
WFC_SMTP_FROM=noreply@whiskyforcharity.com
WFC_SMTP_SSL=true
```

## Gebruik in code

Je kunt emails versturen met deze configuratie via de nieuwe methoden in de SMTP client:

```go
import "dklautomationgo/services"

// ...

func SendWFCExampleEmail(emailService *services.EmailService, smtpClient services.SMTPClient) error {
    // Directe toegang tot SMTP client
    err := smtpClient.SendWFCEmail(
        "recipient@example.com",
        "Test email van Whisky for Charity",
        "<p>Dit is een test email via de Whisky for Charity SMTP server.</p>",
    )
    if err != nil {
        return err
    }
    
    // Of met volledige message
    msg := &services.EmailMessage{
        To:       "recipient@example.com",
        Subject:  "Nog een test email",
        Body:     "<p>Dit is nog een test email via Whisky for Charity.</p>",
        TestMode: false,
    }
    
    return smtpClient.SendWFC(msg)
}
```

## Voorbeeld: Uitbreiden EmailService

Je kunt ook de `EmailService` uitbreiden om een specifieke methode toe te voegen voor WFC emails:

```go
// SendWhiskyForCharityEmail verzendt een email via de WFC configuratie
func (s *EmailService) SendWhiskyForCharityEmail(to, subject, body string) error {
    // Check rate limits
    if !s.rateLimiter.AllowEmail("email_generic", "") {
        return fmt.Errorf("rate limit exceeded")
    }
    
    msg := &EmailMessage{
        To:      to,
        Subject: subject,
        Body:    body,
    }
    
    // Gebruik WFC SMTP configuratie
    err := s.smtpClient.SendWFC(msg)
    if err != nil {
        if s.metrics != nil {
            s.metrics.RecordEmailFailed("wfc_email")
        }
        return err
    }
    
    if s.metrics != nil {
        s.metrics.RecordEmailSent("wfc_email")
    }
    return nil
}
```

## Technische details

De implementatie gebruikt directe SSL verbinding op poort 465 in plaats van STARTTLS op poort 587. Dit is geconfigureerd via de `SSL` property in de gomail dialer.

## Troubleshooting

Als er problemen zijn met het verzenden van emails, controleer het volgende:

1. Controleer of het wachtwoord correct is
2. Zorg dat poort 465 niet geblokkeerd wordt door firewalls
3. Controleer de logs op SMTP fouten
4. Test de verbinding met een eenvoudig telnet commando:
   ```
   telnet arg-plplcl14.argewebhosting.nl 465
   ```

## Overige instellingen

Voor inkomende mail kun je de volgende instellingen gebruiken:

- POP3: arg-plplcl14.argewebhosting.nl:995 (SSL)
- IMAP: arg-plplcl14.argewebhosting.nl:993 (SSL) 