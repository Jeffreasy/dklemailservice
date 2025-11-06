-- GECONSOLIDEERDE V18 - PERFORMANCE OPTIMALISATIES
-- Logica van V1_47.
-- Voegt een grote set indexen toe. Dubbele indexen van V1 worden overgeslagen dankzij 'IF NOT EXISTS'.
-- Oude 'INSERT INTO migraties' is verwijderd.

-- ============================================
-- SECTION 1: MISSING FOREIGN KEY INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_gebruikers_role_id ON gebruikers(role_id);
COMMENT ON INDEX idx_gebruikers_role_id IS 'FK index for RBAC role lookups';

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);
COMMENT ON INDEX idx_aanmeldingen_gebruiker_id IS 'FK index for user registration lookups';

CREATE INDEX IF NOT EXISTS idx_verzonden_emails_contact_id ON verzonden_emails(contact_id);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_aanmelding_id ON verzonden_emails(aanmelding_id);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_template_id ON verzonden_emails(template_id);
COMMENT ON INDEX idx_verzonden_emails_contact_id IS 'FK index for contact form email tracking';
COMMENT ON INDEX idx_verzonden_emails_aanmelding_id IS 'FK index for registration email tracking';
COMMENT ON INDEX idx_verzonden_emails_template_id IS 'FK index for template usage tracking';

CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);
COMMENT ON INDEX idx_contact_antwoorden_contact_id IS 'FK index for contact responses';

CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);
COMMENT ON INDEX idx_aanmelding_antwoorden_aanmelding_id IS 'FK index for registration responses';

-- ============================================
-- SECTION 2: SINGLE COLUMN INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_gebruikers_email ON gebruikers(email);
CREATE INDEX IF NOT EXISTS idx_gebruikers_is_actief ON gebruikers(is_actief) WHERE is_actief = TRUE;
COMMENT ON INDEX idx_gebruikers_is_actief IS 'Partial index for active users only';

CREATE INDEX IF NOT EXISTS idx_contact_formulieren_email ON contact_formulieren(email);
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_email ON aanmeldingen(email);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_afstand ON aanmeldingen(afstand) WHERE afstand IS NOT NULL;
COMMENT ON INDEX idx_aanmeldingen_afstand IS 'Partial index for distance filtering in reports';

CREATE INDEX IF NOT EXISTS idx_verzonden_emails_status ON verzonden_emails(status);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_ontvanger ON verzonden_emails(ontvanger);
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_verzonden_op ON verzonden_emails(verzonden_op DESC);
COMMENT ON INDEX idx_verzonden_emails_status IS 'Status filtering for error tracking';
COMMENT ON INDEX idx_verzonden_emails_verzonden_op IS 'Chronological sorting for email history';

CREATE INDEX IF NOT EXISTS idx_incoming_emails_from ON incoming_emails("from");
COMMENT ON INDEX idx_incoming_emails_from IS 'Sender lookup for email filtering';

-- ============================================
-- SECTION 3: COMPOUND INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status_created 
ON contact_formulieren(status, created_at DESC) 
WHERE beantwoord = FALSE;
COMMENT ON INDEX idx_contact_formulieren_status_created IS 'Compound index for unanswered contact forms dashboard';

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status_created 
ON aanmeldingen(status, created_at DESC);
COMMENT ON INDEX idx_aanmeldingen_status_created IS 'Compound index for registration dashboard';

CREATE INDEX IF NOT EXISTS idx_verzonden_emails_status_tijd 
ON verzonden_emails(status, verzonden_op DESC);
COMMENT ON INDEX idx_verzonden_emails_status_tijd IS 'Compound index for email status tracking over time';

CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_verzonden 
ON contact_antwoorden(contact_id, verzonden_op DESC);
COMMENT ON INDEX idx_contact_antwoorden_contact_verzonden IS 'Chronological responses per contact';

CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_verzonden 
ON aanmelding_antwoorden(aanmelding_id, verzonden_op DESC);
COMMENT ON INDEX idx_aanmelding_antwoorden_aanmelding_verzonden IS 'Chronological responses per registration';

-- ============================================
-- SECTION 4: PARTIAL INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_errors 
ON verzonden_emails(verzonden_op DESC) 
WHERE status = 'failed';
COMMENT ON INDEX idx_verzonden_emails_errors IS 'Partial index for failed email tracking';

CREATE INDEX IF NOT EXISTS idx_incoming_emails_processing 
ON incoming_emails(is_processed, received_at DESC) 
WHERE is_processed = FALSE;
COMMENT ON INDEX idx_incoming_emails_processing IS 'Partial index for unprocessed email queue';

CREATE INDEX IF NOT EXISTS idx_contact_formulieren_nieuw 
ON contact_formulieren(created_at DESC) 
WHERE status = 'nieuw' AND beantwoord = FALSE;
COMMENT ON INDEX idx_contact_formulieren_nieuw IS 'Partial index for new unanswered contact forms';

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_nieuw 
ON aanmeldingen(created_at DESC) 
WHERE status = 'nieuw';
COMMENT ON INDEX idx_aanmeldingen_nieuw IS 'Partial index for new registrations';

CREATE INDEX IF NOT EXISTS idx_chat_participants_active 
ON chat_channel_participants(channel_id, user_id) 
WHERE is_active = TRUE;
COMMENT ON INDEX idx_chat_participants_active IS 'Partial index for active channel participants';

CREATE INDEX IF NOT EXISTS idx_gebruikers_newsletter 
ON gebruikers(email) 
WHERE newsletter_subscribed = TRUE AND is_actief = TRUE;
COMMENT ON INDEX idx_gebruikers_newsletter IS 'Partial index for active newsletter subscribers';

-- ============================================
-- SECTION 5: FULL-TEXT SEARCH INDEXES
-- ============================================
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_fts 
ON contact_formulieren 
USING gin(to_tsvector('dutch', COALESCE(naam, '') || ' ' || COALESCE(email, '') || ' ' || COALESCE(bericht, '')));
COMMENT ON INDEX idx_contact_formulieren_fts IS 'Full-text search on name, email, and message';

CREATE INDEX IF NOT EXISTS idx_aanmeldingen_fts 
ON aanmeldingen 
USING gin(to_tsvector('dutch', COALESCE(naam, '') || ' ' || COALESCE(email, '') || ' ' || COALESCE(bijzonderheden, '')));
COMMENT ON INDEX idx_aanmeldingen_fts IS 'Full-text search on registration details';

CREATE INDEX IF NOT EXISTS idx_chat_messages_fts 
ON chat_messages 
USING gin(to_tsvector('dutch', COALESCE(content, '')));
COMMENT ON INDEX idx_chat_messages_fts IS 'Full-text search on chat message content';

-- ============================================
-- SECTION 6: CHAT SYSTEM OPTIMIZATIONS
-- ============================================
CREATE INDEX IF NOT EXISTS idx_chat_channels_type ON chat_channels(type);
COMMENT ON INDEX idx_chat_channels_type IS 'Channel type filtering';

CREATE INDEX IF NOT EXISTS idx_chat_channels_public 
ON chat_channels(name) 
WHERE is_public = TRUE AND is_active = TRUE;
COMMENT ON INDEX idx_chat_channels_public IS 'Public channel discovery';

CREATE INDEX IF NOT EXISTS idx_chat_messages_reply_to 
ON chat_messages(reply_to_id, created_at DESC) 
WHERE reply_to_id IS NOT NULL;
COMMENT ON INDEX idx_chat_messages_reply_to IS 'Message reply threads';

CREATE INDEX IF NOT EXISTS idx_chat_messages_files 
ON chat_messages(channel_id, created_at DESC) 
WHERE message_type IN ('image', 'file');
COMMENT ON INDEX idx_chat_messages_files IS 'File and image messages';

CREATE INDEX IF NOT EXISTS idx_chat_user_presence_online 
ON chat_user_presence(status, last_seen DESC) 
WHERE status != 'offline';
COMMENT ON INDEX idx_chat_user_presence_online IS 'Online and away users';

CREATE INDEX IF NOT EXISTS idx_chat_participants_unread 
ON chat_channel_participants(user_id, last_read_at) 
WHERE is_active = TRUE;
COMMENT ON INDEX idx_chat_participants_unread IS 'Unread message tracking per user';

-- ============================================
-- SECTION 7: RBAC OPTIMIZATIONS
-- ============================================
CREATE INDEX IF NOT EXISTS idx_user_roles_active_lookup 
ON user_roles(user_id, role_id) 
WHERE is_active = TRUE;
COMMENT ON INDEX idx_user_roles_active_lookup IS 'Active user role assignments';

-- ============================================
-- SECTION 8: NEWSLETTER OPTIMIZATIONS
-- ============================================
CREATE INDEX IF NOT EXISTS idx_newsletters_status 
ON newsletters(sent_at DESC NULLS FIRST);
COMMENT ON INDEX idx_newsletters_status IS 'Draft newsletters first, then sent in reverse chronological order';

-- ============================================
-- SECTION 9: TOKEN CLEANUP OPTIMIZATION
-- ============================================
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_cleanup 
ON refresh_tokens(expires_at) 
WHERE is_revoked = FALSE;
COMMENT ON INDEX idx_refresh_tokens_cleanup IS 'Expired token cleanup for scheduled jobs';

-- ============================================
-- SECTION 10: UPLOADED IMAGES OPTIMIZATIONS
-- ============================================
CREATE INDEX IF NOT EXISTS idx_uploaded_images_active 
ON uploaded_images(user_id, created_at DESC) 
WHERE deleted_at IS NULL;
COMMENT ON INDEX idx_uploaded_images_active IS 'Active (non-deleted) images per user';

-- ============================================
-- POST-MIGRATION RECOMMENDATIONS
-- ============================================
-- 1. ANALYZE;
-- 2. Monitor slow query log.