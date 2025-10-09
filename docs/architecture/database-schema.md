# Database Schema

Complete database schema documentatie voor de DKL Email Service.

## Overzicht

De database gebruikt PostgreSQL met GORM als ORM. Alle tabellen gebruiken UUID's als primary keys en hebben timestamps voor audit trails.

## Core Tables

### Gebruikers

**Table:** `gebruikers`

**Model:** [`models/gebruiker.go:11`](../../models/gebruiker.go:11)

```go
type Gebruiker struct {
    ID                   string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Naam                 string     `gorm:"not null"`
    Email                string     `gorm:"not null;uniqueIndex"`
    WachtwoordHash       string     `gorm:"not null"`
    Rol                  string     `gorm:"default:'gebruiker';index"` // Legacy
    RoleID               *string    `gorm:"type:uuid"` // RBAC role reference
    IsActief             bool       `gorm:"default:true"`
    NewsletterSubscribed bool       `gorm:"default:false;index"`
    LaatsteLogin         *time.Time
    CreatedAt            time.Time  `gorm:"autoCreateTime"`
    UpdatedAt            time.Time  `gorm:"autoUpdateTime"`
    
    Roles []RBACRole `gorm:"many2many:user_roles;"`
}
```

**Schema:**
```sql
CREATE TABLE gebruikers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    wachtwoord_hash VARCHAR(255) NOT NULL,
    rol VARCHAR(50) DEFAULT 'gebruiker',
    role_id UUID REFERENCES roles(id),
    is_actief BOOLEAN DEFAULT true,
    newsletter_subscribed BOOLEAN DEFAULT false,
    laatste_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_gebruikers_email ON gebruikers(email);
CREATE INDEX idx_gebruikers_rol ON gebruikers(rol);
CREATE INDEX idx_gebruikers_newsletter ON gebruikers(newsletter_subscribed);
```

### Contact Formulieren

**Table:** `contact_formulieren`

**Model:** [`models/contact.go:5`](../../models/contact.go:5)

```go
type ContactFormulier struct {
    ID               string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    CreatedAt        time.Time  `gorm:"autoCreateTime"`
    UpdatedAt        time.Time  `gorm:"autoUpdateTime"`
    Naam             string     `gorm:"not null"`
    Email            string     `gorm:"not null;index"`
    Bericht          string     `gorm:"type:text;not null"`
    EmailVerzonden   bool       `gorm:"default:false"`
    EmailVerzondenOp *time.Time
    PrivacyAkkoord   bool       `gorm:"not null"`
    Status           string     `gorm:"default:'nieuw';index"`
    BehandeldDoor    *string
    BehandeldOp      *time.Time
    Notities         *string    `gorm:"type:text"`
    Beantwoord       bool       `gorm:"default:false"`
    AntwoordTekst    string     `gorm:"type:text"`
    AntwoordDatum    *time.Time
    AntwoordDoor     string
    TestMode         bool       `gorm:"type:boolean;not null;default:false"`
    
    Antwoorden []ContactAntwoord `gorm:"foreignKey:ContactID"`
}
```

**Schema:**
```sql
CREATE TABLE contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    email_verzonden BOOLEAN DEFAULT false,
    email_verzonden_op TIMESTAMP WITH TIME ZONE,
    privacy_akkoord BOOLEAN NOT NULL,
    status VARCHAR(50) DEFAULT 'nieuw',
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP WITH TIME ZONE,
    notities TEXT,
    beantwoord BOOLEAN DEFAULT false,
    antwoord_tekst TEXT,
    antwoord_datum TIMESTAMP WITH TIME ZONE,
    antwoord_door VARCHAR(255),
    test_mode BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_contact_email ON contact_formulieren(email);
CREATE INDEX idx_contact_status ON contact_formulieren(status);
```

### Aanmeldingen

**Table:** `aanmeldingen`

**Model:** [`models/aanmelding.go`](../../models/aanmelding.go:1)

```go
type Aanmelding struct {
    ID               string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    CreatedAt        time.Time  `gorm:"autoCreateTime"`
    UpdatedAt        time.Time  `gorm:"autoUpdateTime"`
    Naam             string     `gorm:"not null"`
    Email            string     `gorm:"not null;index"`
    Telefoon         string
    Rol              string     `gorm:"not null;index"`
    Afstand          string
    Ondersteuning    string     `gorm:"type:text"`
    Bijzonderheden   string     `gorm:"type:text"`
    Terms            bool       `gorm:"not null"`
    EmailVerzonden   bool       `gorm:"default:false"`
    EmailVerzondenOp *time.Time
    Status           string     `gorm:"default:'nieuw';index"`
    BehandeldDoor    *string
    BehandeldOp      *time.Time
    Notities         *string    `gorm:"type:text"`
    TestMode         bool       `gorm:"type:boolean;not null;default:false"`
    
    Antwoorden []AanmeldingAntwoord `gorm:"foreignKey:AanmeldingID"`
}
```

### Inkomende Emails

**Table:** `incoming_emails`

**Model:** [`models/incoming_email.go`](../../models/incoming_email.go:1)

```go
type IncomingEmail struct {
    ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    MessageID   string     `gorm:"uniqueIndex;not null"`
    From        string     `gorm:"not null;index"`
    To          string     `gorm:"not null"`
    Subject     string
    Body        string     `gorm:"type:text"`
    ContentType string
    ReceivedAt  time.Time  `gorm:"index"`
    UID         string     `gorm:"uniqueIndex;not null"`
    AccountType string     `gorm:"index"`
    IsProcessed bool       `gorm:"default:false;index"`
    ProcessedAt *time.Time
    CreatedAt   time.Time  `gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}
```

**Schema:**
```sql
CREATE TABLE incoming_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id VARCHAR(255) NOT NULL UNIQUE,
    "from" VARCHAR(255) NOT NULL,
    "to" VARCHAR(255) NOT NULL,
    subject TEXT,
    body TEXT,
    content_type VARCHAR(100),
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    uid VARCHAR(255) NOT NULL UNIQUE,
    account_type VARCHAR(50),
    is_processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_incoming_emails_from ON incoming_emails("from");
CREATE INDEX idx_incoming_emails_received_at ON incoming_emails(received_at);
CREATE INDEX idx_incoming_emails_account_type ON incoming_emails(account_type);
CREATE INDEX idx_incoming_emails_is_processed ON incoming_emails(is_processed);
```

## RBAC Tables

### Roles

**Table:** `roles`

**Model:** [`models/role_rbac.go:8`](../../models/role_rbac.go:8)

```go
type RBACRole struct {
    ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name         string    `gorm:"type:varchar(100);not null;uniqueIndex"`
    Description  string    `gorm:"type:text"`
    IsSystemRole bool      `gorm:"default:false"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
    CreatedBy    *string   `gorm:"type:uuid"`
    
    Permissions []Permission `gorm:"many2many:role_permissions"`
    Users       []Gebruiker  `gorm:"many2many:user_roles"`
}
```

**Schema:** [`database/migrations/V1_20__create_rbac_tables.sql:6`](../../database/migrations/V1_20__create_rbac_tables.sql:6)

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

CREATE INDEX idx_roles_name ON roles(name);
```

### Permissions

**Table:** `permissions`

**Model:** [`models/role_rbac.go:27`](../../models/role_rbac.go:27)

```go
type Permission struct {
    ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Resource           string    `gorm:"type:varchar(100);not null"`
    Action             string    `gorm:"type:varchar(50);not null"`
    Description        string    `gorm:"type:text"`
    IsSystemPermission bool      `gorm:"default:false"`
    CreatedAt          time.Time `gorm:"autoCreateTime"`
    UpdatedAt          time.Time `gorm:"autoUpdateTime"`
    
    Roles []RBACRole `gorm:"many2many:role_permissions"`
}
```

**Schema:** [`database/migrations/V1_20__create_rbac_tables.sql:18`](../../database/migrations/V1_20__create_rbac_tables.sql:18)

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

CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
```

### User Roles

**Table:** `user_roles`

**Model:** [`models/role_rbac.go:62`](../../models/role_rbac.go:62)

```go
type UserRole struct {
    ID         string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID     string     `gorm:"type:uuid;not null"`
    RoleID     string     `gorm:"type:uuid;not null"`
    AssignedAt time.Time  `gorm:"autoCreateTime"`
    AssignedBy *string    `gorm:"type:uuid"`
    ExpiresAt  *time.Time
    IsActive   bool       `gorm:"default:true"`
    
    User Gebruiker `gorm:"foreignKey:UserID"`
    Role RBACRole  `gorm:"foreignKey:RoleID"`
}
```

**Schema:** [`database/migrations/V1_20__create_rbac_tables.sql:40`](../../database/migrations/V1_20__create_rbac_tables.sql:40)

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

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_active ON user_roles(is_active) WHERE is_active = true;
```

### Role Permissions

**Table:** `role_permissions`

**Schema:** [`database/migrations/V1_20__create_rbac_tables.sql:30`](../../database/migrations/V1_20__create_rbac_tables.sql:30)

```sql
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id),
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

## Chat Tables

### Chat Channels

**Table:** `chat_channels`

**Model:** [`models/chat_channel.go`](../../models/chat_channel.go:1)

**Schema:** [`database/migrations/V1_16__create_chat_tables.sql:6`](../../database/migrations/V1_16__create_chat_tables.sql:6)

```sql
CREATE TABLE chat_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL CHECK (type IN ('public', 'private', 'direct')),
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    is_public BOOLEAN DEFAULT false
);
```

**Channel Types:**
- `public` - Openbaar kanaal
- `private` - Privé kanaal (alleen deelnemers)
- `direct` - Direct message tussen twee gebruikers

### Chat Messages

**Table:** `chat_messages`

**Schema:** [`database/migrations/V1_16__create_chat_tables.sql:30`](../../database/migrations/V1_16__create_chat_tables.sql:30)

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
    reply_to_id UUID REFERENCES chat_messages(id) ON DELETE SET NULL,
    edited_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_channel_id_created_at ON chat_messages(channel_id, created_at DESC);
CREATE INDEX idx_chat_messages_user_id ON chat_messages(user_id);
```

**Message Types:**
- `text` - Tekst bericht
- `image` - Afbeelding
- `file` - Bestand
- `system` - Systeem bericht

### Chat Channel Participants

**Table:** `chat_channel_participants`

**Schema:** [`database/migrations/V1_16__create_chat_tables.sql:18`](../../database/migrations/V1_16__create_chat_tables.sql:18)

```sql
CREATE TABLE chat_channel_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID,
    role TEXT DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE,
    last_read_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(channel_id, user_id)
);

CREATE INDEX idx_chat_channel_participants_channel_id ON chat_channel_participants(channel_id);
CREATE INDEX idx_chat_channel_participants_user_id ON chat_channel_participants(user_id);
```

**Participant Roles:**
- `owner` - Kanaal eigenaar
- `admin` - Kanaal beheerder
- `member` - Gewoon lid

## Newsletter Tables

### Newsletters

**Table:** `newsletters`

**Model:** [`models/newsletter.go`](../../models/newsletter.go:1)

```sql
CREATE TABLE newsletters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    html_content TEXT,
    source_url TEXT,
    published_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    recipient_count INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_newsletters_status ON newsletters(status);
CREATE INDEX idx_newsletters_published_at ON newsletters(published_at);
```

**Newsletter Status:**
- `draft` - Concept
- `scheduled` - Ingepland
- `sending` - Wordt verzonden
- `sent` - Verzonden
- `failed` - Mislukt

## Notification Tables

### Notifications

**Table:** `notifications`

**Model:** [`models/notification.go`](../../models/notification.go:1)

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent BOOLEAN DEFAULT false,
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_sent ON notifications(sent);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_priority ON notifications(priority);
```

**Notification Types:**
- `contact` - Contact formulier
- `aanmelding` - Aanmelding
- `auth` - Authenticatie event
- `system` - Systeem event
- `health` - Health check

**Priority Levels:**
- `low` - Lage prioriteit
- `medium` - Normale prioriteit
- `high` - Hoge prioriteit
- `critical` - Kritiek

## Refresh Tokens

**Table:** `refresh_tokens`

**Model:** [`models/refresh_token.go`](../../models/refresh_token.go:1)

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

## Database Views

### User Permissions View

**View:** `user_permissions`

**Schema:** [`database/migrations/V1_20__create_rbac_tables.sql:65`](../../database/migrations/V1_20__create_rbac_tables.sql:65)

```sql
CREATE OR REPLACE VIEW user_permissions AS
SELECT
    ur.user_id,
    u.email,
    r.name as role_name,
    p.resource,
    p.action,
    rp.assigned_at as permission_assigned_at,
    ur.assigned_at as role_assigned_at
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE ur.is_active = true
ORDER BY ur.user_id, r.name, p.resource, p.action;
```

**Gebruik:**
```go
// Repository query
permissions, err := repo.GetUserPermissions(ctx, userID)
```

## Relationships

### Entity Relationship Diagram

```
gebruikers
    ├─── 1:N ──> contact_formulieren (behandeld_door)
    ├─── 1:N ──> aanmeldingen (behandeld_door)
    ├─── 1:N ──> refresh_tokens
    ├─── M:N ──> roles (via user_roles)
    ├─── 1:N ──> chat_messages
    └─── 1:N ──> chat_channel_participants

roles
    ├─── M:N ──> permissions (via role_permissions)
    └─── M:N ──> gebruikers (via user_roles)

chat_channels
    ├─── 1:N ──> chat_messages
    └─── 1:N ──> chat_channel_participants

chat_messages
    ├─── 1:N ──> chat_message_reactions
    └─── 1:1 ──> chat_messages (reply_to)

contact_formulieren
    └─── 1:N ──> contact_antwoorden

aanmeldingen
    └─── 1:N ──> aanmelding_antwoorden
```

## Migrations

### Migration Manager

**Implementatie:** [`database/migrations.go`](../../database/migrations.go:1)

```go
type MigrationManager struct {
    db              *gorm.DB
    migratieRepo    repository.MigratieRepository
    migrationsDir   string
}

func (m *MigrationManager) MigrateDatabase() error {
    // Lees alle migratie bestanden
    files, err := filepath.Glob(filepath.Join(m.migrationsDir, "*.sql"))
    
    for _, file := range files {
        // Check of migratie al is uitgevoerd
        version := extractVersion(file)
        if m.isMigrationApplied(version) {
            continue
        }
        
        // Voer migratie uit
        if err := m.executeMigration(file); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Migration Files

**Locatie:** [`database/migrations/`](../../database/migrations/)

**Naming Convention:**
```
V{major}_{minor}__{description}.sql
```

**Voorbeelden:**
- `V1_16__create_chat_tables.sql`
- `V1_20__create_rbac_tables.sql`
- `V1_28__add_refresh_tokens.sql`

**Volgorde:**
Migraties worden alfabetisch uitgevoerd op basis van versienummer.

## Indexes

### Performance Indexes

**Email Lookups:**
```sql
CREATE INDEX idx_contact_email ON contact_formulieren(email);
CREATE INDEX idx_aanmelding_email ON aanmeldingen(email);
CREATE INDEX idx_incoming_emails_from ON incoming_emails("from");
```

**Status Filtering:**
```sql
CREATE INDEX idx_contact_status ON contact_formulieren(status);
CREATE INDEX idx_aanmelding_status ON aanmeldingen(status);
CREATE INDEX idx_incoming_emails_is_processed ON incoming_emails(is_processed);
```

**Date Sorting:**
```sql
CREATE INDEX idx_chat_messages_channel_id_created_at ON chat_messages(channel_id, created_at DESC);
CREATE INDEX idx_incoming_emails_received_at ON incoming_emails(received_at);
```

**RBAC Performance:**
```sql
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_active ON user_roles(is_active) WHERE is_active = true;
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
```

## Database Configuration

### Connection Pool

**Implementatie:** [`config/database.go:210`](../../config/database.go:210)

```go
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}

sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Settings:**
- **MaxIdleConns:** 10 - Aantal idle connecties
- **MaxOpenConns:** 100 - Maximum aantal open connecties
- **ConnMaxLifetime:** 1 uur - Maximum levensduur van connectie

### SSL Mode

**Development:**
```bash
DB_SSL_MODE=disable
```

**Production:**
```bash
DB_SSL_MODE=require
```

## Backup & Restore

### Backup

**Full Backup:**
```bash
pg_dump -h localhost -U postgres dklemailservice > backup_$(date +%Y%m%d).sql
```

**Schema Only:**
```bash
pg_dump -h localhost -U postgres --schema-only dklemailservice > schema.sql
```

**Data Only:**
```bash
pg_dump -h localhost -U postgres --data-only dklemailservice > data.sql
```

### Restore

**Full Restore:**
```bash
psql -h localhost -U postgres dklemailservice < backup_20240320.sql
```

**Schema Only:**
```bash
psql -h localhost -U postgres dklemailservice < schema.sql
```

## Query Optimization

### Common Queries

**Get User with Permissions:**
```sql
SELECT 
    u.*,
    json_agg(DISTINCT jsonb_build_object(
        'resource', p.resource,
        'action', p.action
    )) as permissions
FROM gebruikers u
LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.is_active = true
LEFT JOIN role_permissions rp ON ur.role_id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = $1
GROUP BY u.id;
```

**Get Unprocessed Emails:**
```sql
SELECT *
FROM incoming_emails
WHERE is_processed = false
ORDER BY received_at DESC
LIMIT 50;
```

### Query Performance

**EXPLAIN ANALYZE:**
```sql
EXPLAIN ANALYZE
SELECT * FROM contact_formulieren
WHERE status = 'nieuw'
ORDER BY created_at DESC
LIMIT 50;
```

## Maintenance

### Vacuum

**Auto Vacuum:**
PostgreSQL auto vacuum is standaard enabled.

**Manual Vacuum:**
```sql
VACUUM ANALYZE contact_formulieren;
VACUUM ANALYZE aanmeldingen;
```

### Statistics

**Update Statistics:**
```sql
ANALYZE contact_formulieren;
ANALYZE aanmeldingen;
```

### Cleanup Old Data

**Archive Old Records:**
```sql
-- Archive contact formulieren ouder dan 1 jaar
INSERT INTO contact_formulieren_archive
SELECT * FROM contact_formulieren
WHERE created_at < NOW() - INTERVAL '1 year'
AND status = 'afgehandeld';

DELETE FROM contact_formulieren
WHERE created_at < NOW() - INTERVAL '1 year'
AND status = 'afgehandeld';
```

## Zie Ook

- [Components](./components.md) - System components
- [Authentication](./authentication-and-authorization.md) - Auth system
- [Deployment Guide](../guides/deployment.md) - Production deployment