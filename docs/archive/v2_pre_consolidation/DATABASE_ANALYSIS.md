# DKL Email Service - PostgreSQL Database Analyse & Optimalisatie

**Datum**: 30 oktober 2025  
**Database Versie**: PostgreSQL 15 (Alpine)  
**Analysetype**: Volledige schema analyse via migraties  
**Docker Container**: dkl-postgres

---

## 📋 Inhoudsopgave

1. [Executive Summary](#executive-summary)
2. [Docker & Database Configuratie](#docker--database-configuratie)
3. [Database Schema Overzicht](#database-schema-overzicht)
4. [Tabel Details & Relaties](#tabel-details--relaties)
5. [Index Analyse](#index-analyse)
6. [Performance Optimalisaties](#performance-optimalisaties)
7. [Security & Best Practices](#security--best-practices)
8. [Aanbevelingen](#aanbevelingen)

---

## 🎯 Executive Summary

De DKL Email Service database bevat **33 tabellen** verdeeld over 6 functionele domeinen:
- Core Email & Gebruikers Management (9 tabellen)
- Chat Systeem (5 tabellen)
- RBAC (Role-Based Access Control) (4 tabellen)
- Content Management (15 tabellen)

**Totaal aantal migraties**: 46 versies (V1.0.0 - V1.46)

### Sterke Punten ✅
- Goede gebruik van UUID primary keys voor distributed systems
- RBAC systeem correct geïmplementeerd met granular permissions
- Adequate indexen op foreign keys en veel-gebruikte zoek kolommen
- Soft delete pattern geïmplementeerd (uploaded_images)
- Proper timestamp tracking (created_at, updated_at)

### Aandachtspunten ⚠️
- Enkele tabellen missen compound indexes voor specifieke query patterns
- Geen expliciet partitioning voor grote tabellen (chat_messages, verzonden_emails)
- Missing indexes op enkele JOIN kolommen
- Geen database-level constraints voor email validatie
- Geen archivering strategie voor oude data

---

## 🐳 Docker & Database Configuratie

### Container Setup (docker-compose.dev.yml)

```yaml
postgres:
  image: postgres:15-alpine
  container_name: dkl-postgres
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: dklemailservice
  ports:
    - "5433:5432"  # Host:Container
  volumes:
    - postgres_data:/var/lib/postgresql/data
```

**Configuratie Analyse:**
- ✅ Gebruik van PostgreSQL 15 (stabiel en performant)
- ✅ Alpine variant (klein footprint)
- ✅ Persistent volume voor data
- ✅ Health check geïmplementeerd
- ⚠️ Default credentials in dev (GOED voor dev, NIET voor prod)
- ⚠️ Port mapping 5433→5432 (let op bij verbindingen)

### Redis Integration
```yaml
redis:
  image: redis:7-alpine
  ports:
    - "6380:6379"
```
- Gebruikt voor caching en rate limiting
- Verbetert performance van frequently accessed data

---

## 📊 Database Schema Overzicht

### Domein 1: Core Email & Gebruikers (9 tabellen)

#### 1. [`migraties`](database/migrations/001_initial_schema.sql:6)
**Doel**: Tracking van database migraties
```sql
CREATE TABLE migraties (
    id SERIAL PRIMARY KEY,
    versie VARCHAR(50) NOT NULL UNIQUE,
    naam VARCHAR(255) NOT NULL,
    toegepast TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
**Kenmerken:**
- SERIAL primary key (auto-increment)
- UNIQUE constraint op versie
- Timestamp tracking

**Optimalisatie**: ✅ Geen optimalisaties nodig - kleine tabel

---

#### 2. [`gebruikers`](database/migrations/001_initial_schema.sql:14)
**Doel**: Gebruikersaccounts met authenticatie en RBAC
```sql
CREATE TABLE gebruikers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    wachtwoord_hash VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL DEFAULT 'gebruiker',  -- Legacy
    role_id UUID REFERENCES roles(id),             -- RBAC
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    laatste_login TIMESTAMP,
    newsletter_subscribed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:**
- PRIMARY KEY op `id`
- UNIQUE index op `email`
- **MISSING**: Index op `role_id` (FK)
- **MISSING**: Index op `is_actief` voor filtered queries

**Relaties:**
- → `roles` (many-to-one via role_id)
- → `user_roles` (many-to-many voor meerdere rollen)
- ← `contact_antwoorden`, `aanmelding_antwoorden` (FK verzonden_door)
- ← `email_templates` (FK created_by)
- ← `refresh_tokens` (one-to-many)
- ← `uploaded_images` (one-to-many)

**Optimalisatie Suggesties:**
```sql
-- Index voor FK lookup
CREATE INDEX idx_gebruikers_role_id ON gebruikers(role_id);

-- Index voor actieve gebruikers filter
CREATE INDEX idx_gebruikers_is_actief ON gebruikers(is_actief) WHERE is_actief = TRUE;

-- Index voor email login lookups (case-insensitive)
CREATE INDEX idx_gebruikers_email_lower ON gebruikers(LOWER(email));

-- Index voor newsletter subscribers
CREATE INDEX idx_gebruikers_newsletter ON gebruikers(newsletter_subscribed) 
WHERE newsletter_subscribed = TRUE;
```

---

#### 3. [`contact_formulieren`](database/migrations/001_initial_schema.sql:27)
**Doel**: Contact form submissions
```sql
CREATE TABLE contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP,
    privacy_akkoord BOOLEAN NOT NULL DEFAULT TRUE,
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP,
    notities TEXT,
    beantwoord BOOLEAN NOT NULL DEFAULT FALSE,
    antwoord_tekst TEXT,
    antwoord_datum TIMESTAMP,
    antwoord_door VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:**
- PRIMARY KEY op `id`
- **MISSING**: Compound index voor status queries
- **MISSING**: Index op email voor lookup
- **MISSING**: Index op created_at voor chronologisch sorteren

**Optimalisatie Suggesties:**
```sql
-- Compound index voor admin dashboard queries
CREATE INDEX idx_contact_formulieren_status_created 
ON contact_formulieren(status, created_at DESC) 
WHERE beantwoord = FALSE;

-- Index voor email lookup
CREATE INDEX idx_contact_formulieren_email ON contact_formulieren(email);

-- Index voor behandeling tracking
CREATE INDEX idx_contact_formulieren_behandeld 
ON contact_formulieren(behandeld_op DESC NULLS LAST);

-- Partial index voor onbeantwoorde forms
CREATE INDEX idx_contact_formulieren_onbeantwoord 
ON contact_formulieren(created_at DESC) 
WHERE beantwoord = FALSE AND status = 'nieuw';
```

---

#### 4. [`contact_antwoorden`](database/migrations/001_initial_schema.sql:40)
**Doel**: Responses naar contact forms
```sql
CREATE TABLE contact_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL REFERENCES contact_formulieren(id) ON DELETE CASCADE,
    bericht TEXT NOT NULL,  -- Oud veld
    tekst TEXT NOT NULL DEFAULT '',
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_door VARCHAR(255) NOT NULL DEFAULT '',
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:**
- PRIMARY KEY op `id`
- **MISSING**: Index op `contact_id` (FK)

**Aandachtspunt:** Duplicatie van velden (bericht vs tekst, verzonden_op vs verzond_op)

**Optimalisatie Suggesties:**
```sql
-- FK index voor joins
CREATE INDEX idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);

-- Compound index voor chronologische queries per contact
CREATE INDEX idx_contact_antwoorden_contact_verzonden 
ON contact_antwoorden(contact_id, verzond_op DESC);

-- Data cleanup: Consolideer dubbele velden
ALTER TABLE contact_antwoorden DROP COLUMN bericht; -- Gebruik alleen 'tekst'
ALTER TABLE contact_antwoorden DROP COLUMN verzonden_op; -- Gebruik alleen 'verzond_op'
```

---

#### 5. [`aanmeldingen`](database/migrations/001_initial_schema.sql:52)
**Doel**: Event registrations
```sql
CREATE TABLE aanmeldingen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    rol VARCHAR(255),
    afstand VARCHAR(255),
    ondersteuning VARCHAR(255),
    bijzonderheden TEXT,
    terms BOOLEAN NOT NULL DEFAULT TRUE,
    steps INTEGER DEFAULT 0,
    gebruiker_id UUID REFERENCES gebruikers(id),  -- Linked user
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP,
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP,
    notities TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:**
- PRIMARY KEY op `id`
- **MISSING**: Index op `gebruiker_id` (FK)
- **MISSING**: Index op email
- **MISSING**: Compound indexes voor queries

**Optimalisatie Suggesties:**
```sql
-- FK index
CREATE INDEX idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);

-- Email lookup
CREATE INDEX idx_aanmeldingen_email ON aanmeldingen(email);

-- Status queries (admin dashboard)
CREATE INDEX idx_aanmeldingen_status_created 
ON aanmeldingen(status, created_at DESC);

-- Afstand filtering (voor rapportage)
CREATE INDEX idx_aanmeldingen_afstand ON aanmeldingen(afstand) WHERE afstand IS NOT NULL;

-- Steps tracking (voor gamification)
CREATE INDEX idx_aanmeldingen_steps ON aanmeldingen(steps DESC) WHERE steps > 0;
```

---

#### 6. [`aanmelding_antwoorden`](database/migrations/001_initial_schema.sql:66)
Structuur identiek aan `contact_antwoorden` maar voor aanmeldingen.

**Dezelfde optimalisaties als contact_antwoorden van toepassing**

---

#### 7. [`email_templates`](database/migrations/001_initial_schema.sql:78)
**Doel**: Herbruikbare email templates
```sql
CREATE TABLE email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL UNIQUE,
    onderwerp VARCHAR(255) NOT NULL,
    inhoud TEXT NOT NULL,
    beschrijving TEXT,
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES gebruikers(id)
);
```

**Indexen:**
- PRIMARY KEY op `id`
- UNIQUE index op `naam`
- **MISSING**: Index op `is_actief`

**Optimalisatie**: ✅ Kleine tabel, minimale optimalisatie nodig
```sql
-- Partial index voor actieve templates
CREATE INDEX idx_email_templates_actief ON email_templates(naam) WHERE is_actief = TRUE;
```

---

#### 8. [`verzonden_emails`](database/migrations/001_initial_schema.sql:91)
**Doel**: Tracking van alle verzonden emails
```sql
CREATE TABLE verzonden_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ontvanger VARCHAR(255) NOT NULL,
    onderwerp VARCHAR(255) NOT NULL,
    inhoud TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'verzonden',
    fout_bericht TEXT,
    contact_id UUID REFERENCES contact_formulieren(id),
    aanmelding_id UUID REFERENCES aanmeldingen(id),
    template_id UUID REFERENCES email_templates(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:**
- PRIMARY KEY op `id`
- **MISSING**: Alle FK indexes
- **MISSING**: Index op status voor error tracking
- **MISSING**: Index op verzonden_op voor chronologie

⚠️ **KRITIEK**: Deze tabel kan zeer groot worden!

**Optimalisatie Suggesties:**
```sql
-- FK indexes
CREATE INDEX idx_verzonden_emails_contact_id ON verzonden_emails(contact_id);
CREATE INDEX idx_verzonden_emails_aanmelding_id ON verzonden_emails(aanmelding_id);
CREATE INDEX idx_verzonden_emails_template_id ON verzonden_emails(template_id);

-- Status & error tracking
CREATE INDEX idx_verzonden_emails_status ON verzonden_emails(status);
CREATE INDEX idx_verzonden_emails_errors 
ON verzonden_emails(verzonden_op DESC) 
WHERE status = 'failed';

-- Time-based queries
CREATE INDEX idx_verzonden_emails_verzonden_op ON verzonden_emails(verzonden_op DESC);

-- Email recipient lookup
CREATE INDEX idx_verzonden_emails_ontvanger ON verzonden_emails(ontvanger);

-- PARTITIONING STRATEGIE (voor grote datasets):
-- Partition by month voor betere performance
CREATE TABLE verzonden_emails_2025_01 PARTITION OF verzonden_emails
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
-- etc...

-- Of gebruik time-based partitioning
SELECT create_hypertable('verzonden_emails', 'verzonden_op'); -- TimescaleDB extensie
```

---

#### 9. [`incoming_emails`](database/migrations/004_create_incoming_emails_table.sql:6)
**Doel**: Incoming email tracking (IMAP fetch)
```sql
CREATE TABLE incoming_emails (
    id VARCHAR(255) PRIMARY KEY,
    message_id VARCHAR(255),
    "from" VARCHAR(255) NOT NULL,
    "to" VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT,
    content_type VARCHAR(255),
    received_at TIMESTAMP NOT NULL,
    uid VARCHAR(255) UNIQUE,
    account_type VARCHAR(50),
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexen:** ✅ Goed geïmplementeerd!
- PRIMARY KEY op `id`
- UNIQUE index op `uid`
- Index op `message_id`
- Index op `account_type`
- Index op `is_processed`

**Optimalisatie Suggestie:**
```sql
-- Compound index voor processing queue
CREATE INDEX idx_incoming_emails_processing 
ON incoming_emails(is_processed, received_at DESC) 
WHERE is_processed = FALSE;

-- Index voor sender queries
CREATE INDEX idx_incoming_emails_from ON incoming_emails("from");
```

---

### Domein 2: Chat Systeem (5 tabellen)

#### 10. [`chat_channels`](database/migrations/V1_16__create_chat_tables.sql:6)
```sql
CREATE TABLE chat_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL CHECK (type IN ('public', 'private', 'direct')),
    is_public BOOLEAN DEFAULT FALSE,  -- Added in V1_17
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE
);
```

**Indexen:**
- PRIMARY KEY op `id`
- **MISSING**: Index op `type` voor filtering
- **MISSING**: Index op `is_public` en `is_active`

**Optimalisatie:**
```sql
-- Type filtering
CREATE INDEX idx_chat_channels_type ON chat_channels(type);

-- Public channels discovery
CREATE INDEX idx_chat_channels_public 
ON chat_channels(name) 
WHERE is_public = TRUE AND is_active = TRUE;
```

---

#### 11. [`chat_channel_participants`](database/migrations/V1_16__create_chat_tables.sql:18)
```sql
CREATE TABLE chat_channel_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID,
    role TEXT DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member')),
    last_read_at TIMESTAMP WITH TIME ZONE,  -- Added in V1_18
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(channel_id, user_id)
);
```

**Indexen:** ✅ Goed!
- PRIMARY KEY op `id`
- UNIQUE constraint op `(channel_id, user_id)`
- Index op `channel_id`
- Index op `user_id`

**Optimalisatie:**
```sql
-- Unread message queries
CREATE INDEX idx_chat_participants_unread 
ON chat_channel_participants(user_id, last_read_at) 
WHERE is_active = TRUE;
```

---

#### 12. [`chat_messages`](database/migrations/V1_16__create_chat_tables.sql:30)
```sql
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID,
    content TEXT,
    message_type TEXT DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    file_url TEXT,
    file_name TEXT,
    file_size INTEGER,
    thumbnail_url TEXT,  -- Added in V1_29
    reply_to_id UUID REFERENCES chat_messages(id) ON DELETE SET NULL,
    edited_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Indexen:** ✅ Compound index aanwezig!
- PRIMARY KEY op `id`
- Index op `(channel_id, created_at DESC)`
- Index op `user_id`

⚠️ **Deze tabel kan zeer groot worden!**

**Optimalisatie:**
```sql
-- Message threads (replies)
CREATE INDEX idx_chat_messages_reply_to ON chat_messages(reply_to_id) 
WHERE reply_to_id IS NOT NULL;

-- File messages
CREATE INDEX idx_chat_messages_files 
ON chat_messages(channel_id, created_at DESC) 
WHERE message_type IN ('image', 'file');

-- Full-text search op content
CREATE INDEX idx_chat_messages_content_fts ON chat_messages USING gin(to_tsvector('dutch', content));

-- PARTITIONING voor oude berichten
-- Maandelijks partitioneren
```

---

#### 13. [`chat_message_reactions`](database/migrations/V1_16__create_chat_tables.sql:46)
```sql
CREATE TABLE chat_message_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES chat_messages(id) ON DELETE CASCADE,
    user_id UUID,
    emoji TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(message_id, user_id, emoji)
);
```

**Indexen:** ✅
- PRIMARY KEY op `id`
- UNIQUE constraint op `(message_id, user_id, emoji)`
- Index op `message_id`

---

#### 14. [`chat_user_presence`](database/migrations/V1_16__create_chat_tables.sql:56)
```sql
CREATE TABLE chat_user_presence (
    user_id UUID PRIMARY KEY,
    status TEXT DEFAULT 'offline' CHECK (status IN ('online', 'away', 'busy', 'offline')),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Optimalisatie:**
```sql
-- Online users query
CREATE INDEX idx_chat_user_presence_online 
ON chat_user_presence(status, last_seen DESC) 
WHERE status != 'offline';
```

---

### Domein 3: RBAC (Role-Based Access Control) (4 tabellen)

#### 15. [`roles`](database/migrations/V1_20__create_rbac_tables.sql:6)
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES gebruikers(id)
);
```

**Indexen:** ✅
- PRIMARY KEY op `id`
- UNIQUE index op `name`
- Index op `name`

---

#### 16. [`permissions`](database/migrations/V1_20__create_rbac_tables.sql:18)
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system_permission BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(resource, action)
);
```

**Indexen:** ✅
- PRIMARY KEY op `id`
- UNIQUE constraint op `(resource, action)`
- Index op `(resource, action)`

---

#### 17. [`role_permissions`](database/migrations/V1_20__create_rbac_tables.sql:30)
```sql
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id),
    UNIQUE(role_id, permission_id)
);
```

**Indexen:** ✅
- PRIMARY KEY op `id`
- UNIQUE constraint op `(role_id, permission_id)`
- Index op `role_id`
- Index op `permission_id`

---

#### 18. [`user_roles`](database/migrations/V1_20__create_rbac_tables.sql:40)
```sql
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, role_id)
);
```

**Indexen:** ✅ Excellent!
- PRIMARY KEY op `id`
- UNIQUE constraint op `(user_id, role_id)`
- Index op `user_id`
- Index op `role_id`
- Partial index op `is_active WHERE is_active = TRUE`

**View:** ✅ `user_permissions` view voor easy querying

---

### Domein 4: Authentication & Tokens (1 tabel)

#### 19. [`refresh_tokens`](database/migrations/V1_28__add_refresh_tokens.sql:4)
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE
);
```

**Indexen:** ✅ Goed!
- PRIMARY KEY op `id`
- UNIQUE index op `token`
- Index op `user_id`
- Index op `token`
- Index op `expires_at`
- Index op `is_revoked`

**Optimalisatie:**
```sql
-- Cleanup van verlopen tokens (scheduled job)
DELETE FROM refresh_tokens 
WHERE expires_at < NOW() - INTERVAL '30 days';

-- Of met partitioning
CREATE INDEX idx_refresh_tokens_cleanup 
ON refresh_tokens(expires_at) 
WHERE is_revoked = FALSE;
```

---

### Domein 5: Content Management (13 tabellen)

#### 20. [`newsletters`](database/migrations/V1_19__add_newsletter_subscribed_and_newsletters.sql:6)
```sql
CREATE TABLE newsletters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    batch_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

**Indexen:**
- PRIMARY KEY op `id`
- Index op `sent_at`

**Optimalisatie:**
```sql
-- Draft vs sent newsletters
CREATE INDEX idx_newsletters_status ON newsletters(sent_at DESC NULLS FIRST);
```

---

#### 21. [`uploaded_images`](database/migrations/V1_30__create_uploaded_images_table.sql:6)
**Doel**: Cloudinary image metadata tracking
```sql
CREATE TABLE uploaded_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    public_id TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    folder TEXT NOT NULL,
    thumbnail_url TEXT,
    deleted_at TIMESTAMP,  -- Soft delete
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Indexen:** ✅ Excellent!
- PRIMARY KEY op `id`
- UNIQUE index op `public_id`
- Index op `user_id`
- Index op `public_id`
- Index op `folder`
- Index op `deleted_at`
- Index op `created_at DESC`

---

#### 22-33. Content Tables (photos, albums, videos, sponsors, etc.)

Alle content tabellen volgen vergelijkbaar patroon:
- UUID primary key
- `visible` boolean voor publication control
- `order_number` voor custom ordering
- Timestamps (created_at, updated_at)
- Adequate indexing

**Standaard optimalisatie voor alle:**
```sql
-- Voor admin interfaces
CREATE INDEX idx_{table}_visible_order ON {table}(order_number) WHERE visible = TRUE;
```

---

## 🚀 Performance Optimalisaties

### High Priority Optimalisaties

#### 1. Missing Foreign Key Indexes
```sql
-- gebruikers
CREATE INDEX idx_gebruikers_role_id ON gebruikers(role_id);

-- aanmeldingen
CREATE INDEX idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);

-- verzonden_emails (KRITIEK!)
CREATE INDEX idx_verzonden_emails_contact_id ON verzonden_emails(contact_id);
CREATE INDEX idx_verzonden_emails_aanmelding_id ON verzonden_emails(aanmelding_id);
CREATE INDEX idx_verzonden_emails_template_id ON verzonden_emails(template_id);

-- contact_antwoorden
CREATE INDEX idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);

-- aanmelding_antwoorden  
CREATE INDEX idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);
```

#### 2. Compound Indexes voor Common Queries
```sql
-- Admin dashboards
CREATE INDEX idx_contact_formulieren_dashboard 
ON contact_formulieren(status, created_at DESC) 
WHERE beantwoord = FALSE;

CREATE INDEX idx_aanmeldingen_dashboard 
ON aanmeldingen(status, created_at DESC);

-- Chat message history
-- Al aanwezig: idx_chat_messages_channel_id_created_at

-- Email tracking
CREATE INDEX idx_verzonden_emails_recent 
ON verzonden_emails(status, verzonden_op DESC);
```

#### 3. Partial Indexes voor Filtered Queries
```sql
-- Actieve gebruikers
CREATE INDEX idx_gebruikers_active ON gebruikers(email) WHERE is_actief = TRUE;

-- Ongelezen berichten
CREATE INDEX idx_chat_messages_unread 
ON chat_channel_participants(user_id, channel_id) 
WHERE last_read_at < NOW();

-- Failed emails
CREATE INDEX idx_verzonden_emails_failed 
ON verzonden_emails(created_at DESC) 
WHERE status = 'failed';
```

#### 4. Full-Text Search Indexes
```sql
-- Contact formulieren zoeken
CREATE INDEX idx_contact_formulieren_fts 
ON contact_formulieren 
USING gin(to_tsvector('dutch', naam || ' ' || email || ' ' || bericht));

-- Chat berichten zoeken
CREATE INDEX idx_chat_messages_fts 
ON chat_messages 
USING gin(to_tsvector('dutch', content));
```

### Partitioning Strategie

Voor tabellen die snel groeien:

#### `verzonden_emails` (Time-based Partitioning)
```sql
-- Maak parent table partitioned
CREATE TABLE verzonden_emails_new (
    LIKE verzonden_emails INCLUDING ALL
) PARTITION BY RANGE (verzonden_op);

-- Maak partities per maand
CREATE TABLE verzonden_emails_2025_01 
PARTITION OF verzonden_emails_new
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE verzonden_emails_2025_02 
PARTITION OF verzonden_emails_new
FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

-- Etc...

-- Migratie van oude data
INSERT INTO verzonden_emails_new SELECT * FROM verzonden_emails;

-- Rename tables
ALTER TABLE verzonden_emails RENAME TO verzonden_emails_old;
ALTER TABLE verzonden_emails_new RENAME TO verzonden_emails;
```

#### `chat_messages` (Time-based Partitioning)
Zelfde aanpak als `verzonden_emails`

### Query Optimization Tips

#### 1. Use EXPLAIN ANALYZE
```sql
EXPLAIN ANALYZE
SELECT * FROM contact_formulieren 
WHERE status = 'nieuw' 
ORDER BY created_at DESC 
LIMIT 20;
```

#### 2. Avoid SELECT *
```sql
-- Slecht
SELECT * FROM gebruikers WHERE email = 'test@example.com';

-- Goed
SELECT id, naam, email, rol FROM gebruikers WHERE email = 'test@example.com';
```

#### 3. Use Prepared Statements
```go
// In Go code
stmt, err := db.Prepare("SELECT * FROM gebruikers WHERE email = $1")
result, err := stmt.Query(email)
```

#### 4. Batch Operations
```sql
-- Slecht: N queries
INSERT INTO verzonden_emails (...) VALUES (...);
INSERT INTO verzonden_emails (...) VALUES (...);

-- Goed: 1 query
INSERT INTO verzonden_emails (...) VALUES
    (...),
    (...),
    (...);
```

---

## 🔒 Security & Best Practices

### 1. SQL Injection Prevention
✅ Al geïmplementeerd via prepared statements in Go

### 2. Password Hashing
✅ Gebruikt `wachtwoord_hash` kolom (bcrypt aangenomen)

### 3. Soft Deletes
⚠️ Alleen geïmplementeerd in `uploaded_images`

**Aanbeveling:** Implementeer soft deletes voor belangrijke tabellen:
```sql
ALTER TABLE contact_formulieren ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE aanmeldingen ADD COLUMN deleted_at TIMESTAMP;
```

### 4. Row-Level Security (RLS)
Niet geïmplementeerd. Overweeg voor multi-tenant scenarios:
```sql
ALTER TABLE chat_messages ENABLE ROW LEVEL SECURITY;

CREATE POLICY chat_messages_policy ON chat_messages
FOR SELECT
USING (
    user_id = current_user_id() OR
    channel_id IN (
        SELECT channel_id FROM chat_channel_participants 
        WHERE user_id = current_user_id()
    )
);
```

### 5. Audit Logging
Niet geïmplementeerd. Overweeg audit trail voor:
- gebruikers wijzigingen
- role assignments
- permission changes

```sql
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name TEXT NOT NULL,
    record_id UUID NOT NULL,
    action TEXT NOT NULL,  -- INSERT, UPDATE, DELETE
    old_values JSONB,
    new_values JSONB,
    changed_by UUID REFERENCES gebruikers(id),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 📈 Maintenance & Monitoring

### Database Maintenance Tasks

#### 1. VACUUM & ANALYZE
```sql
-- Weekly maintenance
VACUUM ANALYZE;

-- Per tabel
VACUUM ANALYZE verzonden_emails;
VACUUM ANALYZE chat_messages;
```

#### 2. Reindex
```sql
-- Maandelijks
REINDEX DATABASE dklemailservice;

-- Of per tabel
REINDEX TABLE verzonden_emails;
```

#### 3. Update Statistics
```sql
ANALYZE gebruikers;
ANALYZE contact_formulieren;
```

### Monitoring Queries

#### Check Table Sizes
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY bytes DESC;
```

#### Check Index Usage
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan ASC;
```

#### Check Missing Indexes (Potential)
```sql
SELECT 
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    seq_tup_read / seq_scan as avg_seq_read
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 20;
```

#### Check Slow Queries (Enable pg_stat_statements)
```sql
-- In postgresql.conf: shared_preload_libraries = 'pg_stat_statements'
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;
```

---

## 🎯 Prioritized Recommendations

### 🔴 Critical (Immediate Action Required)

1. **Add Missing FK Indexes** (Performance Impact: HIGH)
   - `gebruikers.role_id`
   - `aanmeldingen.gebruiker_id`
   - `verzonden_emails`: contact_id, aanmelding_id, template_id
   - `contact_antwoorden.contact_id`
   - `aanmelding_antwoorden.aanmelding_id`

2. **Implement Partitioning for Large Tables** (Performance Impact: HIGH)
   - `verzonden_emails` (monthly partitions)
   - `chat_messages` (monthly partitions)

3. **Add Compound Indexes for Dashboard Queries** (Performance Impact: MEDIUM-HIGH)
   - Contact formulieren: (status, created_at)
   - Aanmeldingen: (status, created_at)

### 🟡 High Priority (Within 1 Month)

4. **Data Cleanup**
   - Remove duplicate columns in `contact_antwoorden` and `aanmelding_antwoorden`
   - Archive old `verzonden_emails` (older than 1 year)
   - Cleanup expired `refresh_tokens`

5. **Add Full-Text Search Indexes**
   - Contact formulieren (naam, email, bericht)
   - Chat messages (content)

6. **Implement Soft Deletes**
   - Add `deleted_at` to belangrijke tabellen

### 🟢 Medium Priority (Within 3 Months)

7. **Monitoring & Alerting**
   - Enable `pg_stat_statements`
   - Setup slow query monitoring
   - Table size monitoring
   - Index usage monitoring

8. **Backup & Recovery**
   - Implement automated backups (pg_dump)
   - Test restore procedures
   - Document recovery process

9. **Performance Testing**
   - Load testing voor belangrijke endpoints
   - Query performance benchmarks
   - Connection pooling optimization

### 🔵 Low Priority (Future Enhancement)

10. **Advanced Features**
    - Row-Level Security voor multi-tenancy
    - Audit logging systeem
    - Time-series optimization (TimescaleDB)
    - Replication setup (read replicas)

---

## 📝 Migration Script Generator

### Complete Optimization Migration

```sql
-- File: database/migrations/V1_47__performance_optimizations.sql

-- ============================================
-- MISSING FOREIGN KEY INDEXES
-- ============================================

-- gebruikers
CREATE INDEX IF NOT EXISTS idx_gebruikers_role_id ON gebruikers(role_id);
CREATE INDEX IF NOT EXISTS idx_gebruikers_is_actief ON gebruikers(is_actief) WHERE is_actief = TRUE;

-- aanmeldingen  
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_email ON aanmeldingen(email);

-- verzonden_emails (KRITIEK!)
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_contact_id ON verzonden_emails(contact_id);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_aanmelding_id ON verzonden_emails(aanmelding_id);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_template_id ON verzonden_emails(template_id);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_status ON verzonden_emails(status);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_ontvanger ON verzonden_emails(ontvanger);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_verzonden_op ON verzonden_emails(verzonden_op DESC);

-- contact_antwoorden
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);

-- aanmelding_antwoorden
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);

-- ============================================
-- COMPOUND INDEXES
-- ============================================

-- contact_formulieren
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status_created 
ON contact_formulieren(status, created_at DESC) 
WHERE beantwoord = FALSE;

CREATE INDEX IF NOT EXISTS idx_contact_formulieren_email ON contact_formulieren(email);

-- aanmeldingen
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status_created 
ON aanmeldingen(status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_afstand 
ON aanmeldingen(afstand) 
WHERE afstand IS NOT NULL;

-- ============================================
-- PARTIAL INDEXES
-- ============================================

-- verzonden_emails errors
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_errors 
ON verzonden_emails(verzonden_op DESC) 
WHERE status = 'failed';

-- incoming_emails processing queue
CREATE INDEX IF NOT EXISTS idx_incoming_emails_processing 
ON incoming_emails(is_processed, received_at DESC) 
WHERE is_processed = FALSE;

-- ============================================
-- FULL-TEXT SEARCH INDEXES
-- ============================================

-- contact_formulieren zoeken
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_fts 
ON contact_formulieren 
USING gin(to_tsvector('dutch', COALESCE(naam, '') || ' ' || COALESCE(email, '') || ' ' || COALESCE(bericht, '')));

-- chat_messages zoeken
CREATE INDEX IF NOT EXISTS idx_chat_messages_fts 
ON chat_messages 
USING gin(to_tsvector('dutch', COALESCE(content, '')));

-- ============================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================

COMMENT ON INDEX idx_verzonden_emails_contact_id IS 'FK index for JOIN performance';
COMMENT ON INDEX idx_verzonden_emails_status IS 'Status filtering for admin dashboard';
COMMENT ON INDEX idx_contact_formulieren_status_created IS 'Compound index for dashboard queries';
COMMENT ON INDEX idx_contact_formulieren_fts IS 'Full-text search on name, email, message';

-- ============================================
-- REGISTER MIGRATION
-- ============================================

INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.47.0', 'Performance optimizations: Missing FK indexes, compound indexes, FTS', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
```

---

## 📊 Summary Statistics

### Database Metrics

| Categorie | Count |
|-----------|-------|
| **Totaal Tabellen** | 33 |
| **Core Tabellen** | 9 |
| **Chat Tabellen** | 5 |
| **RBAC Tabellen** | 4 |
| **Content Tabellen** | 15 |
| **Views** | 1 (user_permissions) |
| **Migraties** | 46 |
| **Indexes (bestaand)** | ~50 |
| **Indexes (voorgesteld)** | +30 |

### Performance Impact Schatting

| Optimalisatie | Impact | Effort |
|---------------|--------|--------|
| FK Indexes | 🔴 HIGH | ⚡ LOW |
| Compound Indexes | 🟡 MEDIUM | ⚡ LOW |
| Partitioning | 🔴 HIGH | 🔨 MEDIUM |
| FTS Indexes | 🟢 LOW | ⚡ LOW |
| Data Cleanup | 🟡 MEDIUM | 🔨 MEDIUM |

---

## 🎓 Best Practices Checklist

- ✅ UUID primary keys (distributed-ready)
- ✅ Timestamps op alle tabellen
- ✅ CASCADE deletes waar logisch
- ✅ UNIQUE constraints waar nodig
- ✅ CHECK constraints voor enums
- ✅ RBAC correct geïmplementeerd
- ⚠️ Foreign key indexes (missing)
- ⚠️ Partitioning voor grote tabellen
- ⚠️ Soft deletes (partially implemented)
- ❌ Row-Level Security
- ❌ Audit logging
- ❌ Automated backups documentation

---

## 📚 Resources

- [PostgreSQL Performance Tips](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Index Strategies](https://www.postgresql.org/docs/current/indexes.html)
- [Partitioning](https://www.postgresql.org/docs/current/ddl-partitioning.html)
- [pg_stat_statements](https://www.postgresql.org/docs/current/pgstatstatements.html)

---

**Document Versie**: 1.0  
**Laatst Bijgewerkt**: 30 oktober 2025  
**Auteur**: Database Analysis Tool  
**Review Status**: Ready for Implementation
