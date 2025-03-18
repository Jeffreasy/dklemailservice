# Notificaties voor DKL Email Service

De DKL Email Service ondersteunt real-time notificaties via Telegram voor belangrijke gebeurtenissen binnen het systeem. Deze notificaties houden je op de hoogte van nieuwe contactverzoeken, aanmeldingen, systeemgebeurtenissen en potentiële beveiligingsproblemen.

## Configuratie

Om notificaties te activeren, voeg je de volgende instellingen toe aan je `.env` bestand:

```
ENABLE_NOTIFICATIONS=true
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
NOTIFICATION_THROTTLE=15m
NOTIFICATION_MIN_PRIORITY=medium
```

### Telegram Bot aanmaken

1. Open Telegram en zoek naar `@BotFather`
2. Start een chat met BotFather en stuur `/newbot`
3. Volg de instructies om een naam en gebruikersnaam voor je bot in te stellen
4. Noteer de bot token die je ontvangt (bijvoorbeeld: `123456789:ABCdefGhIjKlmnOPQrsTUVwxyZ`)
5. Open een chat met je nieuwe bot door op de link te klikken die BotFather stuurt

### Chat ID vinden

Methode 1:
1. Start een chat met je bot en stuur een bericht
2. Bezoek `https://api.telegram.org/bot<YourBOTToken>/getUpdates` in je browser (vervang `<YourBOTToken>` met je eigen token)
3. Zoek in de JSON-respons naar `"chat":{"id":123456789}` - dit nummer is je chat ID

Methode 2 (voor groepschats):
1. Voeg je bot toe aan een groep
2. Stuur een bericht in de groep, waarin je de bot tagt (@jebot)
3. Bezoek `https://api.telegram.org/bot<YourBOTToken>/getUpdates`
4. Zoek in de JSON-respons naar `"chat":{"id":-123456789}` - let op het minteken voor groepschats

## Notificatietypen

De service ondersteunt de volgende typen notificaties:

| Type | Beschrijving | Emoji |
|------|-------------|-------|
| contact | Nieuwe contactformulieren | ⚠️ |
| aanmelding | Nieuwe aanmeldingen | 🔴 |
| auth | Beveiligingsmeldingen | 🚨 |
| system | Systeemgebeurtenissen | ℹ️ |
| health | Gezondheidscheck | ℹ️ |

## Prioriteitsniveaus

Notificaties hebben verschillende prioriteitsniveaus:

| Prioriteit | Emoji | Beschrijving |
|------------|-------|-------------|
| low | ℹ️ | Informatieve meldingen, zoals opstart- en shutdownberichten |
| medium | ⚠️ | Standaardmeldingen, zoals nieuwe contactverzoeken |
| high | 🔴 | Belangrijke meldingen, zoals VIP-aanmeldingen |
| critical | 🚨 | Kritieke meldingen, zoals beveiligingsincidenten |

Je kunt het minimale prioriteitsniveau instellen met `NOTIFICATION_MIN_PRIORITY`. Alleen meldingen met die prioriteit of hoger worden verzonden.

## Throttling

Om te voorkomen dat je overspoeld wordt met meldingen, heeft het systeem een 'throttling' mechanisme. Identieke meldingen binnen een bepaalde periode worden samengevoegd. Je kunt de duur instellen met `NOTIFICATION_THROTTLE`.

Geldige formaten zijn: `10s`, `5m`, `2h`, `24h`, etc.

## Voorbeeld Notificaties

### Contact Notificatie
```
⚠️ Nieuw Contactverzoek

Jan Jansen heeft contact opgenomen.

Email: jan@example.com

Bericht:
Ik heb een vraag over het inschrijven voor de loop. Kunnen jullie mij helpen?
```

### Aanmelding Notificatie
```
🔴 Nieuwe VIP Aanmelding

Piet Pietersen heeft zich aangemeld als VIP deelnemer.

Email: piet@example.com
Telefoonnummer: 06-12345678
```

### Beveiligings Notificatie
```
🚨 Verdachte Inlogpogingen

Er zijn meerdere mislukte inlogpogingen gedetecteerd voor gebruiker:
admin@dekoninklijkeloop.nl

IP-adres: 192.168.1.1
Aantal pogingen: 5
```

### Systeem Notificatie
```
ℹ️ Service Gestart

DKL Email Service is gestart op server-01 in omgeving production.
```

## Implementatiedetails

De notificatieservice is geïmplementeerd als een centrale component die:

1. Verschillende gebeurtenissen in de applicatie vastlegt
2. Prioriteiten toewijst aan meldingen
3. Throttling toepast om spam te voorkomen
4. Berichten formatteert met duidelijke emoji's
5. Meldingen verzendt via Telegram

De service draait in de achtergrond en controleert periodiek of er niet-verzonden notificaties zijn die alsnog verzonden moeten worden. 