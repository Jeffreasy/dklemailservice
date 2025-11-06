-- Gecombineerde Migratie: Geïntegreerde toevoegingen van testdata, aanmeldingen en nieuwe features
-- Beschrijving: Combinatie van migraties V1_11 t/m V1_20 voor een enkelvoudig script met seeds en nieuwe tabellen
-- Versie: 1.1.0 (geconsolideerd)

-- Zorg ervoor dat de pgcrypto extensie beschikbaar is voor gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Maak chat_channels tabel aan (uit V1_16, geïntegreerd met V1_17: toegevoegd is_public)
CREATE TABLE IF NOT EXISTS chat_channels (
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

-- Maak chat_channel_participants tabel aan (uit V1_16, geïntegreerd met V1_18: toegevoegd last_read_at)
CREATE TABLE IF NOT EXISTS chat_channel_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID,
    role TEXT DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    last_read_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(channel_id, user_id)
);

-- Maak chat_messages tabel aan (uit V1_16)
CREATE TABLE IF NOT EXISTS chat_messages (
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

-- Maak chat_message_reactions tabel aan (uit V1_16)
CREATE TABLE IF NOT EXISTS chat_message_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES chat_messages(id) ON DELETE CASCADE,
    user_id UUID,
    emoji TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(message_id, user_id, emoji)
);

-- Maak chat_user_presence tabel aan (uit V1_16)
CREATE TABLE IF NOT EXISTS chat_user_presence (
    user_id UUID PRIMARY KEY,
    status TEXT DEFAULT 'offline' CHECK (status IN ('online', 'away', 'busy', 'offline')),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indices voor chat-tabellen (uit V1_16)
CREATE INDEX IF NOT EXISTS idx_chat_messages_channel_id_created_at ON chat_messages(channel_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_user_id ON chat_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_channel_participants_channel_id ON chat_channel_participants(channel_id);
CREATE INDEX IF NOT EXISTS idx_chat_channel_participants_user_id ON chat_channel_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_message_reactions_message_id ON chat_message_reactions(message_id);

-- Alter gebruikers voor newsletters (uit V1_19: toegevoegd newsletter_subscribed)
ALTER TABLE IF EXISTS gebruikers
    ADD COLUMN IF NOT EXISTS newsletter_subscribed BOOLEAN DEFAULT false;

-- Maak newsletters tabel aan (uit V1_19)
CREATE TABLE IF NOT EXISTS newsletters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    batch_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indices voor newsletters (uit V1_19)
CREATE INDEX IF NOT EXISTS idx_newsletters_sent_at ON newsletters (sent_at);

-- Maak RBAC-tabellen aan (uit V1_20)
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES gebruikers(id),
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system_permission BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(resource, action)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id),
    UNIQUE(role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES gebruikers(id),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, role_id)
);

-- Indices voor RBAC (uit V1_20)
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_active ON user_roles(is_active) WHERE is_active = true;

-- Alter gebruikers voor RBAC (uit V1_20: toegevoegd role_id)
ALTER TABLE gebruikers ADD COLUMN IF NOT EXISTS role_id UUID REFERENCES roles(id);

-- Create view voor user_permissions (uit V1_20)
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

-- Seed data: Inserts voor aanmeldingen en contact_formulieren (gecombineerd uit V1_11, V1_12, V1_13, V1_14, V1_15)
DO $$
BEGIN
    -- Inserts voor aanmeldingen (uit V1_11, V1_13, V1_14, V1_15)
    IF NOT EXISTS (SELECT 1 FROM aanmeldingen WHERE id = '3e62d5d3-070d-47b1-a1ef-30665f982789') THEN
        INSERT INTO aanmeldingen (id, created_at, updated_at, naam, email, telefoon, rol, afstand, ondersteuning, bijzonderheden, terms, email_verzonden, email_verzonden_op, status, test_mode) 
        VALUES 
            ('3e62d5d3-070d-47b1-a1ef-30665f982789', '2025-03-23 17:06:26.132297+00', '2025-03-23 17:06:26.132297+00', 'TGTest', 'laventejeffrey@gmail.com', '06123456789', 'Begeleider', '15 KM', 'Anders', 'Telegram Test bericht - officiele weg', TRUE, FALSE, NULL, 'nieuw', TRUE),
            ('51205069-a20a-4231-bcb9-cb2d6fd042c8', '2025-03-23 14:08:53.328628+00', '2025-03-23 14:08:53.328628+00', 'Manuela van Zwam', 'rik.van-harxen@sheerenloo.nl', NULL, 'Deelnemer', '2.5 KM', 'Ja', 'Vaste begeleider die meeloopt', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('275490c0-1021-4bf4-9005-7df9884b0fe6', '2025-03-22 16:43:03.19496+00', '2025-03-22 16:43:03.19496+00', 'Bas heijenk ', 'basheijenk96@gmail.com', NULL, 'Deelnemer', '2.5 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('b5e67c64-6bfa-46d4-ae40-98b093c8b720', '2025-03-17 19:55:07.379647+00', '2025-03-21 06:47:40.512037+00', 'Salih', 'topraks@gmail.com', NULL, 'Deelnemer', '2.5 KM', 'Nee', '', TRUE, TRUE, '2025-03-21 06:47:38.608+00', 'nieuw', FALSE),
            ('b2fd3412-8368-409f-8029-b2cdd581ade1', '2025-03-11 11:28:05.483427+00', '2025-03-21 06:47:41.806005+00', 'Manuela van zwam', 'benjaminlaan.64a@sheerenloo.nl', NULL, 'Deelnemer', '2.5 KM', 'Nee', '', TRUE, TRUE, '2025-03-21 06:47:39.928+00', 'nieuw', FALSE),
            ('391f63c5-f034-466e-8a1f-ba9d06ed1192', '2025-03-09 16:52:09.564437+00', '2025-03-21 06:47:42.987352+00', 'Joyce Thielen', 'Joyce.thielen@sheerenloo.nl', '', 'Begeleider', '6 KM', 'Nee', '', TRUE, TRUE, '2025-03-21 06:47:41.071+00', 'nieuw', FALSE),
            ('391f2579-d7cb-4ef3-afbe-14dc4115c519', '2025-03-08 10:28:37.053379+00', '2025-03-08 10:28:38.042991+00', 'Dick van Norden', 'Enckerkamp.27@sheerenloo.nl', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-03-08 10:28:42.2+00', 'nieuw', FALSE),
            ('90e477cc-89d1-4524-8edc-b697be8c504d', '2025-03-08 10:27:31.078013+00', '2025-03-08 10:27:32.642032+00', 'Angelo van Ingen', 'Enckerkamp.27@sheerenloo.nl', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-03-08 10:27:36.787+00', 'nieuw', FALSE),
            ('26ea058b-2608-49d0-862a-611e98d7dc61', '2025-02-20 11:55:30.818333+00', '2025-02-20 11:55:32.209292+00', 'Janny van de Wall', 'mjvdwal@hotmail.com', NULL, 'Deelnemer', '10 KM', 'Nee', '', TRUE, TRUE, '2025-02-20 11:55:33.048+00', 'nieuw', FALSE),
            ('f4fc2312-ec8a-4dfc-90b5-a8da317618e6', '2025-01-28 21:58:37.55756+00', '2025-01-28 21:58:38.527642+00', 'Martin van der Wal', 'mjvdwal@hotmail.com', '', 'Deelnemer', '10 KM', 'Ja', 'loopt samen met Dirk-Jan mee als vrijwilliger', TRUE, TRUE, '2025-01-28 21:58:39.243+00', 'nieuw', FALSE),
            ('4bfe814b-e0b8-4e60-9f46-fe38852d9ecb', '2025-01-28 21:56:10.099964+00', '2025-01-28 21:56:11.309572+00', 'Dirk-Jan Hempe', 'mjvdwal@hotmail.com', '', 'Deelnemer', '10 KM', 'Nee', '', TRUE, TRUE, '2025-01-28 21:56:12.004+00', 'nieuw', FALSE),
            ('9f464844-8c93-4190-90f0-e74765c7f09a', '2025-03-24 18:47:05.053651+00', '2025-03-24 18:47:05.053651+00', 'Karin de Jong', 'karin.de.jong82@outlook.com', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('51855fec-eab9-494a-9321-c40d22da4ffc', '2025-03-24 17:27:25.191642+00', '2025-03-24 17:27:25.191642+00', 'Mirjam Kerkvliet', 'mirjam.kerkvliet@gmail.com', NULL, 'Deelnemer', '15 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('47775742-8950-4b94-9dd1-571ff4902688', '2025-03-24 17:26:12.1724+00', '2025-03-24 17:26:12.1724+00', 'Arno Kerkvliet', 'arno.kerkvliet@gmail.com', NULL, 'Deelnemer', '15 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('d92ed75c-c275-47a4-88a9-ff7a4106f8ee', '2025-03-24 09:17:59.726501+00', '2025-03-24 09:17:59.726501+00', 'Jean-paul Hup', 'molenkamp.19@sheerenloo.nl', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('ecb8332b-ea39-4611-9f58-64921226f2a6', '2025-03-24 09:16:42.111808+00', '2025-03-24 09:16:42.111808+00', 'Annerieke Mandemaker-Timmer', 'annerieketimmer@hotmail.com', '06 17 37 28 40 ', 'Begeleider', '6 KM', 'Nee', '', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('1ca80f61-f5c1-431f-b224-e6557150b65b', '2025-03-30 08:07:45.334762+00', '2025-03-30 08:07:45.334762+00', 'Han van Doornik', 'LaanvanGS.26@sheerenloo.nl', NULL, 'Deelnemer', '2.5 KM', 'Ja', 'Ik wil wel graag begeleiding ', TRUE, FALSE, NULL, 'nieuw', FALSE),
            ('2499686e-2e62-4827-9079-78b468cb26c9', '2025-03-29 10:07:12.393466+00', '2025-03-29 10:07:56.169434+00', 'Bertram tijsma', 'Klaskehiddes@gmail.com', NULL, 'Deelnemer', '15 KM', 'Nee', '', TRUE, TRUE, '2025-03-29 10:07:57.838+00', 'verwerkt', FALSE),
            ('9f75b1df-4c72-4e36-9901-0f74cc26574f', '2025-03-29 10:05:04.276906+00', '2025-03-29 10:07:55.533152+00', 'Klaske van de glind', 'Klaskehiddes@gmail.com', NULL, 'Deelnemer', '15 KM', 'Nee', '', TRUE, TRUE, '2025-03-29 10:07:57.215+00', 'verwerkt', FALSE),
            ('db3ec762-dd54-4ba7-98ab-981235cc316a', '2025-03-26 16:26:53.777384+00', '2025-03-26 16:30:58.285233+00', 'Mila Veenendaal', 'gaminggirlayla@gmail.com', NULL, 'Deelnemer', '10 KM', 'Nee', '', TRUE, TRUE, '2025-03-26 16:30:57.269+00', 'verwerkt', FALSE),
            ('d17c16c6-c423-43de-a876-d40326b62d9e', '2025-03-26 16:25:43.756211+00', '2025-03-26 16:30:57.680904+00', 'Ayla Toprak', 'gamergirlayla@gmail.com', NULL, 'Deelnemer', '10 KM', 'Nee', '', TRUE, TRUE, '2025-03-26 16:30:56.685+00', 'verwerkt', FALSE),
            ('917206a7-a28d-4bad-8b41-ed127eab743a', '2025-03-26 12:27:20.236848+00', '2025-03-26 16:30:57.073558+00', 'A. Bistolfi', 'nedarg@icloud.com', NULL, 'Deelnemer', '15 KM', 'Nee', '', TRUE, TRUE, '2025-03-26 16:30:56.062+00', 'verwerkt', FALSE),
            ('61a3f823-82f6-4107-a9a9-b6a18e6b12c3', '2025-04-14 20:03:47.361483+00', '2025-04-15 13:23:43.293548+00', 'Henk Rekers ', 'h.rekers59@kpnmail.nl', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-04-15 13:23:40.051+00', 'verwerkt', FALSE),
            ('8a5c2c59-ebca-411e-9471-a38be47e1192', '2025-04-14 20:00:47.479484+00', '2025-04-15 13:23:42.628968+00', 'Hilde Rekers ', 'h.rekers59@kpnmail.nl', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-04-15 13:23:39.203+00', 'verwerkt', FALSE),
            ('3c420d85-5f78-4645-88dd-c9e90eb8b6f7', '2025-04-12 10:21:37.470747+00', '2025-04-15 15:38:46.288552+00', 'Henk Rekers ', 'h.rekers1959@kpnmail.nl', NULL, 'Deelnemer', '6 KM', 'Nee', 'email mislukt', TRUE, FALSE, '2025-04-12 16:52:21.796+00', 'nieuw', FALSE),
            ('721d297d-d670-46b1-a581-bcc095565bbd', '2025-04-12 10:20:25.399025+00', '2025-04-15 15:38:52.374917+00', 'Hilde Rekers ', 'h.rekers1959@kpnmail.nl', NULL, 'Deelnemer', '6 KM', 'Nee', 'email mislukt', TRUE, FALSE, '2025-04-12 16:52:20.444+00', 'nieuw', FALSE),
            ('02f9605b-8b26-4461-9ca7-ed7ebbd12311', '2025-04-12 07:09:19.884142+00', '2025-04-12 16:52:21.459029+00', 'Theun ', 'diesbosje@hotmail.com', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-04-12 16:52:19.056+00', 'verwerkt', FALSE),
            ('282e6ec6-2c97-4b46-a992-273c826c1f91', '2025-04-12 07:08:26.565956+00', '2025-04-12 16:52:20.449671+00', 'Albert ', 'diesbosje@hotmail.com', NULL, 'Deelnemer', '6 KM', 'Nee', '', TRUE, TRUE, '2025-04-12 16:52:18.027+00', 'verwerkt', FALSE),
            ('d10065f9-229d-4a1d-9178-324bdffa57c0', '2025-04-12 07:06:59.872906+00', '2025-04-12 16:52:19.826576+00', 'Diesmer ', 'diesbosje@hotmail.com', '0613429612', 'Begeleider', '6 KM', 'Nee', '', TRUE, TRUE, '2025-04-12 16:52:17.255+00', 'verwerkt', FALSE),
            ('613770dd-4733-4b54-964d-c3753cfebfd7', '2025-04-06 18:48:14.993723+00', '2025-04-10 11:15:45.930907+00', 'Sylvia Dijkstra', 'sylvia.dijkstra@sheerenloo.nl', '0683081728 ', 'Begeleider', '2.5 KM', 'Nee', '', TRUE, TRUE, '2025-04-10 11:15:46.282+00', 'verwerkt', FALSE),
            ('2aaffc78-4dca-4b06-9a6a-3799cbfbe67b', '2025-03-31 11:57:23.536871+00', '2025-04-01 17:11:03.413137+00', 'Noa hiddes', 'Klaskehiddes@gmail.com', NULL, 'Deelnemer', '2.5 KM', 'Nee', '', TRUE, TRUE, '2025-04-01 17:11:04.313+00', 'verwerkt', FALSE),
            ('200f862d-4d80-4ccd-8f06-bf795be727fb', '2025-03-31 11:56:37.582458+00', '2025-04-01 17:11:04.334788+00', 'Anneke van de Glind ', 'Klaskehiddes@gmail.com', NULL, 'Deelnemer', '2.5 KM', 'Nee', '', TRUE, TRUE, '2025-04-01 17:11:05.257+00', 'verwerkt', FALSE);
    END IF;

    -- Inserts voor contact_formulieren (uit V1_12)
    IF NOT EXISTS (SELECT 1 FROM contact_formulieren WHERE id = '6ce2e9a3-59fd-4430-aa5c-66df48fbd695') THEN
        INSERT INTO contact_formulieren (id, created_at, updated_at, naam, email, bericht, email_verzonden, email_verzonden_op, privacy_akkoord, status, behandeld_door, behandeld_op, notities) 
        VALUES 
            ('6ce2e9a3-59fd-4430-aa5c-66df48fbd695', '2025-01-28 00:03:05.830054+00', '2025-02-05 02:23:23.029501+00', 'je geheime liefde', 'de.konining@willem.alexander.nl', 'Gedeelte doneren doet het niet.', TRUE, '2025-01-28 00:03:08.041+00', TRUE, 'afgehandeld', 'marieke@dekoninklijkeloop.nl', '2025-02-05 02:23:22.97+00', NULL),
            ('78910428-f760-485d-ae57-653db478ca35', '2025-03-23 07:15:20.8851+00', '2025-03-23 07:15:20.8851+00', 'Bas heijenk ', 'basheijenk96@gmail.com', 'Hallo ik heb me op gegeven maar ik kan helaas niet sorry ', FALSE, NULL, TRUE, 'nieuw', NULL, NULL, NULL);
    END IF;

    -- Log voor monitoring
    RAISE NOTICE 'Gecombineerde seeds uitgevoerd: aanmeldingen en contact_formulieren toegevoegd indien niet bestaand.';
END $$;

-- Registreer de gecombineerde migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.1.0', 'Geïntegreerde toevoegingen van testdata, aanmeldingen en nieuwe features (chat, newsletters, RBAC)', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;