# DKL Automation Go Service

Een email automation service geschreven in Go voor het DKL project.

## Vereisten

- Go 1.21 of hoger
- SMTP server toegang (bijvoorbeeld Gmail)

## Setup

1. Clone de repository
2. Kopieer `.env.example` naar `.env` en vul de juiste waarden in:
   ```
   PORT=8080
   ALLOWED_ORIGINS=http://localhost:3000,https://your-vercel-app.vercel.app
   SMTP_HOST=smtp.gmail.com
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-specific-password
   SMTP_FROM=your-email@gmail.com
   ```

3. Installeer dependencies:
   ```bash
   go mod download
   ```

4. Start de server:
   ```bash
   go run main.go
   ```

## API Endpoints

### Health Check
```
GET /health
```

### Send Email
```
POST /send-email
Content-Type: application/json

{
    "to": "recipient@example.com",
    "subject": "Email onderwerp",
    "body": "Email inhoud (HTML ondersteund)"
}
```

## Deployment op Render

1. Maak een nieuwe Web Service aan op Render
2. Verbind met je GitHub repository
3. Kies "Docker" als runtime
4. Voeg de environment variables toe uit je `.env` bestand
5. Deploy! 