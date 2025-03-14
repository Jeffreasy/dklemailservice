# Authenticatie in DKL Email Service

Dit document beschrijft de authenticatie functionaliteit in de DKL Email Service.

## Overzicht

De authenticatie is gebaseerd op JWT (JSON Web Tokens) en biedt de volgende functionaliteit:

- Login met email en wachtwoord
- Logout
- Wachtwoord reset
- Gebruikersprofiel ophalen
- Rol-gebaseerde toegangscontrole (admin vs. gebruiker)

## Endpoints

### Login

```
POST /api/auth/login
```

**Request body:**

```json
{
  "email": "gebruiker@example.com",
  "wachtwoord": "wachtwoord123"
}
```

**Response:**

```json
{
  "token": "jwt_token_hier",
  "message": "Login succesvol"
}
```

### Logout

```
POST /api/auth/logout
```

**Response:**

```json
{
  "message": "Logout succesvol"
}
```

### Gebruikersprofiel ophalen

```
GET /api/auth/profile
```

**Headers:**

```
Authorization: Bearer jwt_token_hier
```

**Response:**

```json
{
  "id": "7157f3f6-da85-4058-9d38-19133ec93b03",
  "naam": "Admin",
  "email": "admin@dekoninklijkeloop.nl",
  "rol": "admin",
  "is_actief": true,
  "laatste_login": "2023-03-14T15:22:28.710911Z",
  "created_at": "2023-03-14T15:22:28.710911Z"
}
```

### Wachtwoord reset

```
POST /api/auth/reset-password
```

**Headers:**

```
Authorization: Bearer jwt_token_hier
```

**Request body:**

```json
{
  "huidig_wachtwoord": "wachtwoord123",
  "nieuw_wachtwoord": "nieuw_wachtwoord456"
}
```

**Response:**

```json
{
  "message": "Wachtwoord succesvol gewijzigd"
}
```

## Middleware

### Auth Middleware

De `AuthMiddleware` controleert of de gebruiker is ingelogd door het JWT token te valideren.

### Admin Middleware

De `AdminMiddleware` controleert of de ingelogde gebruiker een admin is.

### Rate Limit Middleware

De `RateLimitMiddleware` beperkt het aantal verzoeken dat een gebruiker kan doen binnen een bepaalde periode.

## Beveiligingsmaatregelen

- JWT tokens worden opgeslagen in HTTP-only cookies
- Rate limiting voor login pogingen
- Wachtwoorden worden gehashed opgeslagen met bcrypt
- CSRF bescherming
- Secure cookies voor token opslag

## Configuratie

De volgende omgevingsvariabelen kunnen worden gebruikt om de authenticatie te configureren:

```
# JWT configuratie
JWT_SECRET=change_this_in_production
JWT_TOKEN_EXPIRY=24h

# Rate limiting voor login
LOGIN_LIMIT_COUNT=5
LOGIN_LIMIT_PERIOD=300
LOGIN_LIMIT_PER_IP=true
```

## Gebruikers

De database bevat standaard één admin gebruiker:

- Email: admin@dekoninklijkeloop.nl
- Wachtwoord: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy (gehashed)
- Rol: admin

## Testen

Je kunt de authenticatie testen met de volgende commando's:

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"admin123"}'

# Gebruikersprofiel ophalen
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer jwt_token_hier"

# Wachtwoord wijzigen
curl -X POST http://localhost:8080/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer jwt_token_hier" \
  -d '{"huidig_wachtwoord":"admin123","nieuw_wachtwoord":"nieuw_wachtwoord"}'

# Uitloggen
curl -X POST http://localhost:8080/api/auth/logout
``` 