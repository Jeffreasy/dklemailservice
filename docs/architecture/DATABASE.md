# Database Architecture

PostgreSQL database schema and design documentation for DKL Email Service.

## Overview

The system uses PostgreSQL 17 with:
- Comprehensive RBAC (Role-Based Access Control)
- Optimized indexes for performance
- Foreign key constraints for data integrity
- Timestamps for audit trails
- Proper normalization

## Core Tables

### Users & Authentication

#### `gebruikers` (Users)
Primary user accounts table.

```sql
CREATE TABLE gebruikers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    wachtwoord_hash VARCHAR(255) NOT NULL,
    naam VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    profielfoto_url TEXT,
    actief BOOLEAN DEFAULT true,
    geverifieerd BOOLEAN DEFAULT false,
    laatste_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gebruikers_email ON gebruikers(email);
CREATE INDEX idx_gebruikers_actief ON gebruikers(actief);
```

#### `refresh_tokens`
JWT refresh token storage.

```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(500) UNIQUE NOT NULL,
    gebruiker_id INTEGER NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT false
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_gebruiker ON refresh_tokens(gebruiker_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

### RBAC System

#### `roles` (Roles)
System roles definition.

```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Full system access'),
    ('moderator', 'Content moderation access'),
    ('user', 'Standard user access'),
    ('guest', 'Limited read-only access');
```

#### `permissions`
Granular permissions.

```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);
```

#### `role_permissions`
Many-to-many relationship between roles and permissions.

```sql
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
```

#### `user_roles`
User role assignments.

```sql
CREATE TABLE user_roles (
    gebruiker_id INTEGER NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER REFERENCES gebruikers(id),
    PRIMARY KEY (gebruiker_id, role_id)
);

CREATE INDEX idx_user_roles_gebruiker ON user_roles(gebruiker_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
```

### Content Management

#### `albums`
Photo album management.

```sql
CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    titel VARCHAR(255) NOT NULL,
    beschrijving TEXT,
    cover_foto_url TEXT,
    zichtbaar BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_albums_zichtbaar ON albums(zichtbaar);
CREATE INDEX idx_albums_created_by ON albums(created_by);
```

#### `fotos`
Individual photos.

```sql
CREATE TABLE fotos (
    id SERIAL PRIMARY KEY,
    album_id INTEGER REFERENCES albums(id) ON DELETE CASCADE,
    afbeelding_url TEXT NOT NULL,
    cloudinary_public_id VARCHAR(255),
    titel VARCHAR(255),
    beschrijving TEXT,
    volgorde INTEGER DEFAULT 0,
    uploaded_by INTEGER REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fotos_album ON fotos(album_id);
CREATE INDEX idx_fotos_cloudinary ON fotos(cloudinary_public_id);
CREATE INDEX idx_fotos_volgorde ON fotos(album_id, volgorde);
```

#### `videos`
Video content management.

```sql
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    titel VARCHAR(255) NOT NULL,
    beschrijving TEXT,
    youtube_url TEXT NOT NULL,
    youtube_id VARCHAR(50),
    thumbnail_url TEXT,
    zichtbaar BOOLEAN DEFAULT true,
    created_by INTEGER REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_videos_zichtbaar ON videos(zichtbaar);
CREATE INDEX idx_videos_youtube_id ON videos(youtube_id);
```

### Events & Participants

#### `evenementen` (Events)
Event management.

```sql
CREATE TABLE evenementen (
    id SERIAL PRIMARY KEY,
    naam VARCHAR(255) NOT NULL,
    beschrijving TEXT,
    datum DATE NOT NULL,
    locatie VARCHAR(255),
    max_deelnemers INTEGER,
    inschrijvingen_open BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_evenementen_datum ON evenementen(datum);
CREATE INDEX idx_evenementen_inschrijvingen ON evenementen(inschrijvingen_open);
```

#### `aanmeldingen` (Registrations)
Event registrations.

```sql
CREATE TABLE aanmeldingen (
    id SERIAL PRIMARY KEY,
    evenement_id INTEGER NOT NULL REFERENCES evenementen(id) ON DELETE CASCADE,
    gebruiker_id INTEGER REFERENCES gebruikers(id) ON DELETE SET NULL,
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_aanmeldingen_evenement ON aanmeldingen(evenement_id);
CREATE INDEX idx_aanmeldingen_gebruiker ON aanmeldingen(gebruiker_id);
CREATE INDEX idx_aanmeldingen_status ON aanmeldingen(status);
```

### Notulen (Meeting Notes)

#### `notulen`
Meeting notes storage.

```sql
CREATE TABLE notulen (
    id SERIAL PRIMARY KEY,
    titel VARCHAR(255) NOT NULL,
    datum DATE NOT NULL,
    aanwezig TEXT[],
    afwezig TEXT[],
    verslag TEXT NOT NULL,
    actiepunten TEXT[],
    bijlagen TEXT[],
    created_by INTEGER REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notulen_datum ON notulen(datum);
CREATE INDEX idx_notulen_created_by ON notulen(created_by);
```

### Communication

#### `chat_channels`
Chat channel management.

```sql
CREATE TABLE chat_channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) DEFAULT 'public',
    created_by INTEGER REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_channels_type ON chat_channels(type);
```

#### `chat_messages`
Chat messages.

```sql
CREATE TABLE chat_messages (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    edited BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_messages_channel ON chat_messages(channel_id, created_at);
CREATE INDEX idx_chat_messages_user ON chat_messages(user_id);
```

#### `chat_channel_participants`
Channel membership.

```sql
CREATE TABLE chat_channel_participants (
    channel_id INTEGER NOT NULL REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_read_at TIMESTAMP,
    PRIMARY KEY (channel_id, user_id)
);

CREATE INDEX idx_channel_participants_user ON chat_channel_participants(user_id);
```

### Notifications

#### `notifications`
User notifications.

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, read, created_at);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_data ON notifications USING gin(data);
```

### Newsletter

#### `newsletters`
Newsletter content storage.

```sql
CREATE TABLE newsletters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    html_content TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    scheduled_send_time TIMESTAMP,
    sent_at TIMESTAMP,
    created_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_newsletters_status ON newsletters(status);
CREATE INDEX idx_newsletters_created_by ON newsletters(created_by);
```

---

### Contact Management

#### `contact_formulieren`
Contact form submissions.

```sql
CREATE TABLE contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'nieuw',
    notities TEXT,
    verwerkt_op TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contact_status ON contact_formulieren(status);
CREATE INDEX idx_contact_email ON contact_formulieren(email);
CREATE INDEX idx_contact_created ON contact_formulieren(created_at);
```

#### `contact_antwoorden`
Responses to contact forms.

```sql
CREATE TABLE contact_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL REFERENCES contact_formulieren(id) ON DELETE CASCADE,
    antwoord TEXT NOT NULL,
    beantwoord_door UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contact_antwoorden_contact ON contact_antwoorden(contact_id);
```

---

### Participants & Registrations

#### `participants`
Event participants (persons).

```sql
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    rol VARCHAR(100),
    afstand VARCHAR(50),
    ondersteuning TEXT,
    bijzonderheden TEXT,
    steps INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'actief',
    notities TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_participants_email ON participants(email);
CREATE INDEX idx_participants_rol ON participants(rol);
CREATE INDEX idx_participants_status ON participants(status);
```

#### `participant_antwoorden`
Responses to participants.

```sql
CREATE TABLE participant_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    antwoord TEXT NOT NULL,
    beantwoord_door UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_participant_antwoorden_participant ON participant_antwoorden(participant_id);
```

#### `event_registrations`
Event registrations linking participants to events.

```sql
CREATE TABLE event_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    tracking_status VARCHAR(50) DEFAULT 'registered',
    distance_route VARCHAR(50),
    participant_role VARCHAR(100),
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    check_in_time TIMESTAMP,
    start_time TIMESTAMP,
    finish_time TIMESTAMP,
    steps INTEGER DEFAULT 0,
    total_distance DECIMAL(10,2) DEFAULT 0,
    last_location_update TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, participant_id)
);

CREATE INDEX idx_event_registrations_event ON event_registrations(event_id);
CREATE INDEX idx_event_registrations_participant ON event_registrations(participant_id);
CREATE INDEX idx_event_registrations_status ON event_registrations(status);
CREATE INDEX idx_event_registrations_tracking ON event_registrations(tracking_status);
```

#### `participant_roles`
Role definitions for participants.

```sql
CREATE TABLE participant_roles (
    name VARCHAR(100) PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO participant_roles (name, description) VALUES
    ('deelnemer', 'Evenement deelnemer'),
    ('vrijwilliger', 'Evenement vrijwilliger'),
    ('begeleider', 'Evenement begeleider');
```

#### `distances`
Distance routes available for events.

```sql
CREATE TABLE distances (
    name VARCHAR(50) PRIMARY KEY,
    distance_km DECIMAL(5,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO distances (name, distance_km, description) VALUES
    ('5km', 5.00, '5 kilometer route'),
    ('10km', 10.00, '10 kilometer route'),
    ('15km', 15.00, '15 kilometer route'),
    ('21km', 21.00, 'Halve marathon');
```

---

### Gamification

#### `achievements`
Achievement definitions.

```sql
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url TEXT,
    required_steps INTEGER NOT NULL,
    category VARCHAR(100),
    points INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_achievements_category ON achievements(category);
CREATE INDEX idx_achievements_active ON achievements(is_active);
```

#### `badges`
Badge definitions and tiers.

```sql
CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    tier VARCHAR(50) DEFAULT 'bronze',
    required_achievement_id UUID REFERENCES achievements(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_badges_tier ON badges(tier);
CREATE INDEX idx_badges_achievement ON badges(required_achievement_id);
```

#### `user_achievements`
Tracks which users have unlocked which achievements.

```sql
CREATE TABLE user_achievements (
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_achievement ON user_achievements(achievement_id);
```

#### `user_badges`
Tracks badges awarded to users.

```sql
CREATE TABLE user_badges (
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    badge_id UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    awarded_by UUID REFERENCES gebruikers(id),
    PRIMARY KEY (user_id, badge_id)
);

CREATE INDEX idx_user_badges_user ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge ON user_badges(badge_id);
```

#### `leaderboards`
Leaderboard rankings.

```sql
CREATE TABLE leaderboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period VARCHAR(50) NOT NULL,
    participant_id UUID NOT NULL REFERENCES participants(id),
    steps INTEGER DEFAULT 0,
    distance DECIMAL(10,2) DEFAULT 0,
    rank INTEGER,
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_leaderboards_period ON leaderboards(period);
CREATE INDEX idx_leaderboards_participant ON leaderboards(participant_id);
CREATE INDEX idx_leaderboards_rank ON leaderboards(period, rank);
```

#### `route_funds`
Fund allocation per route.

```sql
CREATE TABLE route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_name VARCHAR(50) NOT NULL,
    allocated_amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    year INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_route_funds_year ON route_funds(year);
CREATE INDEX idx_route_funds_active ON route_funds(is_active);
```

---

### CMS Tables

#### `partners`
Partner organizations.

```sql
CREATE TABLE partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo TEXT,
    website TEXT,
    tier VARCHAR(50),
    since TIMESTAMP,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_partners_visible ON partners(visible);
CREATE INDEX idx_partners_order ON partners(order_number);
```

#### `sponsors`
Event sponsors.

```sql
CREATE TABLE sponsors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_url TEXT,
    website_url TEXT,
    order_number INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sponsors_visible ON sponsors(visible);
CREATE INDEX idx_sponsors_active ON sponsors(is_active);
CREATE INDEX idx_sponsors_order ON sponsors(order_number);
```

#### `videos`
Video content management.

```sql
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id VARCHAR(255),
    url TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_videos_visible ON videos(visible);
CREATE INDEX idx_videos_order ON videos(order_number);
```

#### `fotos` (Photos)
Individual photo management.

```sql
CREATE TABLE fotos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    afbeelding_url TEXT NOT NULL,
    cloudinary_public_id VARCHAR(255),
    title VARCHAR(255),
    description TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fotos_visible ON fotos(visible);
CREATE INDEX idx_fotos_cloudinary ON fotos(cloudinary_public_id);
```

#### `albums`
Photo albums.

```sql
CREATE TABLE albums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titel VARCHAR(255) NOT NULL,
    beschrijving TEXT,
    cover_foto_url TEXT,
    zichtbaar BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_albums_zichtbaar ON albums(zichtbaar);
CREATE INDEX idx_albums_order ON albums(order_number);
```

#### `album_photos`
Many-to-many relationship between albums and photos.

```sql
CREATE TABLE album_photos (
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    photo_id UUID NOT NULL REFERENCES fotos(id) ON DELETE CASCADE,
    order_number INTEGER DEFAULT 0,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (album_id, photo_id)
);

CREATE INDEX idx_album_photos_album ON album_photos(album_id, order_number);
```

#### `radio_recordings`
Radio recordings and podcasts.

```sql
CREATE TABLE radio_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    audio_url TEXT NOT NULL,
    thumbnail_url TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_radio_recordings_visible ON radio_recordings(visible);
CREATE INDEX idx_radio_recordings_date ON radio_recordings(date);
```

#### `program_schedule`
Event program schedule.

```sql
CREATE TABLE program_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    location VARCHAR(255),
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_program_schedule_visible ON program_schedule(visible);
CREATE INDEX idx_program_schedule_start ON program_schedule(start_time);
```

#### `social_links`
Social media links.

```sql
CREATE TABLE social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    icon_url TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_social_links_visible ON social_links(visible);
CREATE INDEX idx_social_links_platform ON social_links(platform);
```

#### `social_embeds`
Embedded social media content.

```sql
CREATE TABLE social_embeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(100) NOT NULL,
    embed_code TEXT NOT NULL,
    url TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_social_embeds_visible ON social_embeds(visible);
```

#### `title_sections`
Page title sections.

```sql
CREATE TABLE title_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_key VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle TEXT,
    background_image TEXT,
    visible BOOLEAN DEFAULT true,
    order_number INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_title_sections_key ON title_sections(section_key);
CREATE INDEX idx_title_sections_visible ON title_sections(visible);
```

#### `under_construction`
Maintenance mode settings.

```sql
CREATE TABLE under_construction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_active BOOLEAN DEFAULT false,
    title VARCHAR(255),
    message TEXT,
    estimated_end TIMESTAMP,
    show_countdown BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

### Email Management

#### `incoming_emails`
Fetched incoming emails.

```sql
CREATE TABLE incoming_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_type VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    sender_name VARCHAR(255),
    subject TEXT,
    body TEXT,
    html_body TEXT,
    received_at TIMESTAMP NOT NULL,
    processed BOOLEAN DEFAULT false,
    uid VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_incoming_emails_account ON incoming_emails(account_type);
CREATE INDEX idx_incoming_emails_processed ON incoming_emails(processed);
CREATE INDEX idx_incoming_emails_uid ON incoming_emails(uid);
CREATE INDEX idx_incoming_emails_received ON incoming_emails(received_at);
```

#### `email_templates`
Reusable email templates.

```sql
CREATE TABLE email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    subject VARCHAR(255) NOT NULL,
    html_content TEXT NOT NULL,
    text_content TEXT,
    variables JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_templates_name ON email_templates(name);
CREATE INDEX idx_email_templates_active ON email_templates(is_active);
```

#### `verzonden_emails`
Sent email log.

```sql
CREATE TABLE verzonden_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_email VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    template_name VARCHAR(100),
    status VARCHAR(50) DEFAULT 'sent',
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0
);

CREATE INDEX idx_verzonden_emails_recipient ON verzonden_emails(recipient_email);
CREATE INDEX idx_verzonden_emails_template ON verzonden_emails(template_name);
CREATE INDEX idx_verzonden_emails_status ON verzonden_emails(status);
CREATE INDEX idx_verzonden_emails_sent ON verzonden_emails(sent_at);
```

---

### Image Tracking

#### `uploaded_images`
Track uploaded images to Cloudinary.

```sql
CREATE TABLE uploaded_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(255) UNIQUE NOT NULL,
    url TEXT NOT NULL,
    secure_url TEXT,
    format VARCHAR(50),
    width INTEGER,
    height INTEGER,
    bytes INTEGER,
    folder VARCHAR(255),
    uploaded_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_uploaded_images_public_id ON uploaded_images(public_id);
CREATE INDEX idx_uploaded_images_folder ON uploaded_images(folder);
CREATE INDEX idx_uploaded_images_uploaded_by ON uploaded_images(uploaded_by);
```

---

### Whisky for Charity

#### `wfc_orders`
Whisky for Charity orders.

```sql
CREATE TABLE wfc_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(255) UNIQUE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_address TEXT,
    customer_city VARCHAR(255),
    customer_postal VARCHAR(20),
    customer_country VARCHAR(100),
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wfc_orders_order_id ON wfc_orders(order_id);
CREATE INDEX idx_wfc_orders_email ON wfc_orders(customer_email);
CREATE INDEX idx_wfc_orders_status ON wfc_orders(status);
```

#### `wfc_order_items`
Individual items in WFC orders.

```sql
CREATE TABLE wfc_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255),
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wfc_order_items_order ON wfc_order_items(order_id);
```

---

### Lookup Tables

#### `event_status_types`
Event status definitions.

```sql
CREATE TABLE event_status_types (
    status VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    display_name VARCHAR(100),
    color VARCHAR(7),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO event_status_types (status, description, display_name, color) VALUES
    ('upcoming', 'Event is gepland maar nog niet gestart', 'Binnenkort', '#3B82F6'),
    ('active', 'Event is momenteel gaande', 'Actief', '#10B981'),
    ('completed', 'Event is afgelopen', 'Voltooid', '#6B7280'),
    ('cancelled', 'Event is geannuleerd', 'Geannuleerd', '#EF4444');
```

#### `registration_status_types`
Registration status definitions.

```sql
CREATE TABLE registration_status_types (
    status VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    display_name VARCHAR(100),
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO registration_status_types (status, description, display_name, color) VALUES
    ('pending', 'Registratie in afwachting van bevestiging', 'In Afwachting', '#F59E0B'),
    ('confirmed', 'Registratie is bevestigd', 'Bevestigd', '#10B981'),
    ('cancelled', 'Registratie is geannuleerd', 'Geannuleerd', '#EF4444'),
    ('completed', 'Deelnemer heeft event voltooid', 'Voltooid', '#6B7280');
```

## Database Relationships

### Complete Entity Relationship Diagram

```
gebruikers (Users)
    ├── 1:N → refresh_tokens
    ├── N:M → roles (via user_roles)
    ├── 1:N → notulen (created_by)
    ├── 1:N → albums (created_by)
    ├── 1:N → uploaded_images (uploaded_by)
    ├── 1:N → chat_messages
    ├── N:M → chat_channels (via chat_channel_participants)
    ├── 1:N → notifications
    ├── N:M → achievements (via user_achievements)
    ├── N:M → badges (via user_badges)
    ├── 1:N → newsletters (created_by)
    ├── 1:N → contact_antwoorden (beantwoord_door)
    └── 1:N → participant_antwoorden (beantwoord_door)

roles
    ├── N:M → permissions (via role_permissions)
    └── N:M → gebruikers (via user_roles)

events
    ├── 1:N → event_registrations
    └── 1:1 → event_status_types (status FK)

participants
    ├── 1:N → event_registrations
    ├── 1:N → participant_antwoorden
    ├── 1:N → leaderboards
    └── 1:1 → participant_roles (rol FK)

event_registrations
    ├── N:1 → events
    ├── N:1 → participants
    ├── 1:1 → registration_status_types (status FK)
    └── 1:1 → distances (distance_route FK)

albums
    └── N:M → fotos (via album_photos)

chat_channels
    ├── 1:N → chat_messages
    ├── N:M → gebruikers (via chat_channel_participants)
    └── 1:N → chat_message_reactions

chat_messages
    └── 1:N → chat_message_reactions

contact_formulieren
    └── 1:N → contact_antwoorden

achievements
    ├── N:M → gebruikers (via user_achievements)
    └── 1:N → badges (required_achievement_id FK)

badges
    ├── N:M → gebruikers (via user_badges)
    └── N:1 → achievements (required_achievement_id FK)

wfc_orders
    └── 1:N → wfc_order_items (order_id)
```

### Key Relationships Explained

**Users & Authentication:**
- Each user can have multiple refresh tokens
- Users have multiple roles (many-to-many)
- Users create various content (albums, notulen, newsletters)

**RBAC System:**
- Roles contain multiple permissions
- Users assigned multiple roles
- Permissions cached in Redis for performance

**Events & Participants:**
- Participants register for events via event_registrations
- Each registration tracks steps, status, and location
- Foreign keys to lookup tables for status descriptions

**Gamification:**
- Achievements unlock badges
- Users collect achievements and badges
- Leaderboards track participant rankings

**Content Management:**
- Albums contain multiple photos (many-to-many)
- All CMS content has visibility and ordering
- Cloudinary integration for media storage

## Indexing Strategy

### Primary Indexes
All tables have primary key indexes automatically created.

### Foreign Key Indexes
Indexes on all foreign key columns for optimal JOIN performance.

### Search Indexes
- Email lookups: `idx_gebruikers_email`
- Active users: `idx_gebruikers_actief`
- Token validation: `idx_refresh_tokens_token`
- Date ranges: `idx_evenementen_datum`, `idx_notulen_datum`

### Composite Indexes
- Chat messages by channel and time: `(channel_id, created_at)`
- User notifications: `(user_id, read, created_at)`

### JSONB Indexes (GIN)
- Notification data: `notifications.data` using GIN index

## Performance Optimization

### Connection Pooling
```go
// Database configuration
MaxOpenConns:    25
MaxIdleConns:    5
ConnMaxLifetime: 5 * time.Minute
ConnMaxIdleTime: 10 * time.Minute
```

### Query Optimization
- Use prepared statements
- Limit result sets with pagination
- Avoid N+1 queries with proper JOINs
- Use EXPLAIN ANALYZE for slow queries

### Partitioning Strategy
Consider partitioning for high-volume tables:
- `chat_messages` by month
- `notifications` by month
- `audit_logs` by month

### Maintenance Tasks
```sql
-- Regular vacuum and analyze
VACUUM ANALYZE;

-- Reindex if needed
REINDEX TABLE table_name;

-- Update statistics
ANALYZE table_name;
```

## Backup Strategy

### Automated Backups
- Daily full backups
- Transaction log archiving
- 30-day retention period
- Off-site backup storage

### Point-in-Time Recovery
Enable WAL archiving for PITR:
```sql
-- postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'cp %p /backup/archive/%f'
```

### Backup Script Example
```bash
#!/bin/bash
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
pg_dump -U postgres -d dklemailservice > backup_${TIMESTAMP}.sql
```

## Migration Management

### Migration Files
Located in `database/migrations/`:
- `V01__initial_schema.sql` - Core tables
- `V02__seed_data.sql` - Default data
- `V03__add_test_data.sql` - Test records
- `V04__create_chat_tables.sql` - Chat system
- `V05__add_newsletter_tables.sql` - Newsletter
- `V06__create_rbac_tables.sql` - RBAC system
- `V07__seed_rbac_tables.sql` - RBAC data
- `V08__migrate_and_assign_user_roles.sql` - User roles
- `V09__create_refresh_tokens_table.sql` - Token storage
- `V10__create_uploaded_images_table.sql` - Image tracking
- And more...

### Running Migrations
```bash
# Apply all migrations
go run database/migrations/run_migrations.go

# Or use PowerShell script
.\database\apply_all_migrations.ps1
```

## Security Considerations

### Password Storage
- Passwords hashed with bcrypt (cost factor 12)
- Never store plain text passwords
- Enforce password complexity

### SQL Injection Prevention
- Use parameterized queries exclusively
- Validate all user input
- Escape special characters

### Access Control
- Database user with minimal required permissions
- Separate read-only user for analytics
- SSL/TLS for database connections

### Audit Logging
Track security-relevant operations:
- Login attempts
- Permission changes
- Data modifications by admins

## Database Monitoring

### Key Metrics
- Connection pool usage
- Query execution time
- Lock wait times
- Cache hit ratio
- Disk I/O

### Slow Query Log
```sql
-- Enable slow query logging
-- postgresql.conf
log_min_duration_statement = 1000  -- Log queries > 1 second
```

### Health Checks
```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check database size
SELECT pg_size_pretty(pg_database_size('dklemailservice'));

-- Check table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Development vs Production

### Development
- Less strict constraints
- Test data included
- Verbose logging
- Lower connection pool limits

### Production
- Strict foreign key enforcement
- No test data
- Optimized logging
- Higher connection pool limits
- SSL required
- Regular backups
- Monitoring enabled

---

For migration details, see [Migration Guide](../guides/MIGRATIONS.md).
For RBAC implementation, see [RBAC Architecture](./AUTHENTICATION.md).