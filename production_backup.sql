--
-- PostgreSQL database dump
--

\restrict 0jXUPxC9NSPHoCpxmC2c42IQDmqL4d8huoBscL62YaMpTxcWQsnDsF0Jrw5Clge

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg120+1)
-- Dumped by pg_dump version 16.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: refresh_dashboard_stats(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.refresh_dashboard_stats() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY dashboard_stats;
END;
$$;


--
-- Name: FUNCTION refresh_dashboard_stats(); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.refresh_dashboard_stats() IS 'Refresh dashboard statistics view (concurrent safe)';


--
-- Name: update_aanmelding_antwoorden_count(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_aanmelding_antwoorden_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE aanmeldingen 
        SET antwoorden_count = antwoorden_count + 1 
        WHERE id = NEW.aanmelding_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE aanmeldingen 
        SET antwoorden_count = GREATEST(0, antwoorden_count - 1) 
        WHERE id = OLD.aanmelding_id;
    END IF;
    RETURN NULL;
END;
$$;


--
-- Name: update_contact_antwoorden_count(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_contact_antwoorden_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE contact_formulieren 
        SET antwoorden_count = antwoorden_count + 1 
        WHERE id = NEW.contact_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE contact_formulieren 
        SET antwoorden_count = GREATEST(0, antwoorden_count - 1) 
        WHERE id = OLD.contact_id;
    END IF;
    RETURN NULL;
END;
$$;


--
-- Name: update_route_funds_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_route_funds_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


--
-- Name: FUNCTION update_updated_at_column(); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.update_updated_at_column() IS 'Generic trigger function to update updated_at timestamp';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: aanmelding_antwoorden; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.aanmelding_antwoorden (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    aanmelding_id uuid NOT NULL,
    verzonden_op timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    tekst text DEFAULT ''::text NOT NULL,
    verzond_op timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    verzond_door character varying(255) DEFAULT ''::character varying NOT NULL,
    email_verzonden boolean DEFAULT false NOT NULL,
    verzonden_door character varying(255)
);
ALTER TABLE ONLY public.aanmelding_antwoorden ALTER COLUMN aanmelding_id SET STATISTICS 500;


--
-- Name: TABLE aanmelding_antwoorden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.aanmelding_antwoorden IS 'Antwoorden op aanmeldingen';


--
-- Name: aanmeldingen; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.aanmeldingen (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    naam character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    telefoon character varying(50),
    status character varying(50) DEFAULT 'nieuw'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    rol character varying(255),
    afstand character varying(255),
    ondersteuning character varying(255),
    bijzonderheden text,
    terms boolean DEFAULT true NOT NULL,
    email_verzonden boolean DEFAULT false NOT NULL,
    email_verzonden_op timestamp without time zone,
    behandeld_door character varying(255),
    behandeld_op timestamp without time zone,
    notities text,
    test_mode boolean DEFAULT false NOT NULL,
    steps integer DEFAULT 0,
    gebruiker_id uuid,
    antwoorden_count integer DEFAULT 0,
    CONSTRAINT aanmeldingen_email_check CHECK (((email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT aanmeldingen_email_consistency CHECK ((((email_verzonden = false) AND (email_verzonden_op IS NULL)) OR ((email_verzonden = true) AND (email_verzonden_op IS NOT NULL)))),
    CONSTRAINT aanmeldingen_naam_not_empty CHECK ((length(TRIM(BOTH FROM naam)) > 0)),
    CONSTRAINT aanmeldingen_status_check CHECK (((status)::text = ANY ((ARRAY['nieuw'::character varying, 'bevestigd'::character varying, 'geannuleerd'::character varying, 'voltooid'::character varying])::text[]))),
    CONSTRAINT aanmeldingen_steps_check CHECK ((steps >= 0))
);
ALTER TABLE ONLY public.aanmeldingen ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE ONLY public.aanmeldingen ALTER COLUMN status SET STATISTICS 500;


--
-- Name: TABLE aanmeldingen; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.aanmeldingen IS 'Aanmeldingen voor De Koninklijke Loop';


--
-- Name: COLUMN aanmeldingen.rol; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.aanmeldingen.rol IS 'Rol van de deelnemer (deelnemer, vrijwilliger, sponsor)';


--
-- Name: COLUMN aanmeldingen.afstand; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.aanmeldingen.afstand IS 'Gekozen afstand voor hardlopers';


--
-- Name: COLUMN aanmeldingen.test_mode; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.aanmeldingen.test_mode IS 'Geeft aan of dit een testaanmelding is (geen echte email verzenden)';


--
-- Name: COLUMN aanmeldingen.gebruiker_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.aanmeldingen.gebruiker_id IS 'Link naar gebruikersaccount voor authenticatie en step tracking';


--
-- Name: COLUMN aanmeldingen.antwoorden_count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.aanmeldingen.antwoorden_count IS 'Cached count of responses - updated via trigger';


--
-- Name: album_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.album_photos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    album_id uuid NOT NULL,
    photo_id uuid NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: albums; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.albums (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    cover_photo_id uuid,
    visible boolean DEFAULT true NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: chat_channel_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_channel_participants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    channel_id uuid,
    user_id uuid,
    role text DEFAULT 'member'::text,
    joined_at timestamp with time zone DEFAULT now(),
    last_seen_at timestamp with time zone,
    is_active boolean DEFAULT true,
    last_read_at timestamp with time zone,
    CONSTRAINT chat_channel_participants_role_check CHECK ((role = ANY (ARRAY['owner'::text, 'admin'::text, 'member'::text])))
);


--
-- Name: chat_channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_channels (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    type text NOT NULL,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    is_active boolean DEFAULT true,
    is_public boolean DEFAULT false,
    CONSTRAINT chat_channels_type_check CHECK ((type = ANY (ARRAY['public'::text, 'private'::text, 'direct'::text])))
);


--
-- Name: chat_message_reactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_message_reactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message_id uuid,
    user_id uuid,
    emoji text NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    channel_id uuid,
    user_id uuid,
    content text,
    message_type text DEFAULT 'text'::text,
    file_url text,
    file_name text,
    file_size integer,
    reply_to_id uuid,
    edited_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    thumbnail_url text,
    CONSTRAINT chat_messages_message_type_check CHECK ((message_type = ANY (ARRAY['text'::text, 'image'::text, 'file'::text, 'system'::text])))
);


--
-- Name: TABLE chat_messages; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.chat_messages IS 'High-traffic table - monitor size and consider partitioning';


--
-- Name: chat_user_presence; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_user_presence (
    user_id uuid NOT NULL,
    status text DEFAULT 'offline'::text,
    last_seen timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT chat_user_presence_status_check CHECK ((status = ANY (ARRAY['online'::text, 'away'::text, 'busy'::text, 'offline'::text])))
);


--
-- Name: contact_antwoorden; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contact_antwoorden (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contact_id uuid NOT NULL,
    verzonden_op timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    tekst text DEFAULT ''::text NOT NULL,
    verzond_op timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    verzond_door character varying(255) DEFAULT ''::character varying NOT NULL,
    email_verzonden boolean DEFAULT false NOT NULL,
    verzonden_door character varying(255)
);
ALTER TABLE ONLY public.contact_antwoorden ALTER COLUMN contact_id SET STATISTICS 500;


--
-- Name: TABLE contact_antwoorden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.contact_antwoorden IS 'Antwoorden op contactformulieren';


--
-- Name: contact_formulieren; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contact_formulieren (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    naam character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    bericht text NOT NULL,
    status character varying(50) DEFAULT 'nieuw'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    email_verzonden boolean DEFAULT false NOT NULL,
    email_verzonden_op timestamp without time zone,
    privacy_akkoord boolean DEFAULT true NOT NULL,
    behandeld_door character varying(255),
    behandeld_op timestamp without time zone,
    notities text,
    beantwoord boolean DEFAULT false NOT NULL,
    antwoord_tekst text,
    antwoord_datum timestamp without time zone,
    antwoord_door character varying(255),
    test_mode boolean DEFAULT false NOT NULL,
    antwoorden_count integer DEFAULT 0,
    CONSTRAINT contact_formulieren_bericht_not_empty CHECK ((length(TRIM(BOTH FROM bericht)) > 0)),
    CONSTRAINT contact_formulieren_email_check CHECK (((email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT contact_formulieren_email_consistency CHECK ((((email_verzonden = false) AND (email_verzonden_op IS NULL)) OR ((email_verzonden = true) AND (email_verzonden_op IS NOT NULL)))),
    CONSTRAINT contact_formulieren_naam_not_empty CHECK ((length(TRIM(BOTH FROM naam)) > 0)),
    CONSTRAINT contact_formulieren_status_check CHECK (((status)::text = ANY ((ARRAY['nieuw'::character varying, 'in_behandeling'::character varying, 'beantwoord'::character varying, 'gesloten'::character varying])::text[])))
);
ALTER TABLE ONLY public.contact_formulieren ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE ONLY public.contact_formulieren ALTER COLUMN status SET STATISTICS 500;


--
-- Name: TABLE contact_formulieren; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.contact_formulieren IS 'Contactformulieren van de website';


--
-- Name: COLUMN contact_formulieren.email_verzonden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.contact_formulieren.email_verzonden IS 'Geeft aan of er een email is verzonden naar de afzender';


--
-- Name: COLUMN contact_formulieren.test_mode; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.contact_formulieren.test_mode IS 'Geeft aan of dit een testbericht is (geen echte email verzenden)';


--
-- Name: COLUMN contact_formulieren.antwoorden_count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.contact_formulieren.antwoorden_count IS 'Cached count of responses - updated via trigger';


--
-- Name: verzonden_emails; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.verzonden_emails (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ontvanger character varying(255) NOT NULL,
    onderwerp character varying(255) NOT NULL,
    inhoud text NOT NULL,
    verzonden_op timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status character varying(50) DEFAULT 'verzonden'::character varying NOT NULL,
    contact_id uuid,
    aanmelding_id uuid,
    template_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    fout_bericht text
);
ALTER TABLE ONLY public.verzonden_emails ALTER COLUMN status SET STATISTICS 500;
ALTER TABLE ONLY public.verzonden_emails ALTER COLUMN contact_id SET STATISTICS 500;
ALTER TABLE ONLY public.verzonden_emails ALTER COLUMN aanmelding_id SET STATISTICS 500;


--
-- Name: TABLE verzonden_emails; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.verzonden_emails IS 'High-traffic table - monitor size and consider partitioning';


--
-- Name: dashboard_stats; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.dashboard_stats AS
 SELECT 'contact_formulieren'::text AS entity,
    contact_formulieren.status,
    contact_formulieren.beantwoord,
    count(*) AS count,
    max(contact_formulieren.created_at) AS last_created
   FROM public.contact_formulieren
  GROUP BY contact_formulieren.status, contact_formulieren.beantwoord
UNION ALL
 SELECT 'aanmeldingen'::text AS entity,
    aanmeldingen.status,
    NULL::boolean AS beantwoord,
    count(*) AS count,
    max(aanmeldingen.created_at) AS last_created
   FROM public.aanmeldingen
  GROUP BY aanmeldingen.status
UNION ALL
 SELECT 'verzonden_emails'::text AS entity,
    verzonden_emails.status,
    NULL::boolean AS beantwoord,
    count(*) AS count,
    max(verzonden_emails.verzonden_op) AS last_created
   FROM public.verzonden_emails
  WHERE (verzonden_emails.verzonden_op > (now() - '30 days'::interval))
  GROUP BY verzonden_emails.status
  WITH NO DATA;


--
-- Name: MATERIALIZED VIEW dashboard_stats; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON MATERIALIZED VIEW public.dashboard_stats IS 'Cached dashboard statistics - refresh hourly or on demand';


--
-- Name: email_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    naam character varying(255) NOT NULL,
    onderwerp character varying(255) NOT NULL,
    inhoud text NOT NULL,
    beschrijving text,
    is_actief boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by uuid
);


--
-- Name: gebruikers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gebruikers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    naam character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    wachtwoord_hash character varying(255) NOT NULL,
    rol character varying(50) DEFAULT 'gebruiker'::character varying NOT NULL,
    is_actief boolean DEFAULT true NOT NULL,
    laatste_login timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    newsletter_subscribed boolean DEFAULT false,
    role_id uuid,
    CONSTRAINT gebruikers_email_check CHECK (((email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT gebruikers_naam_not_empty CHECK ((length(TRIM(BOTH FROM naam)) > 0))
);
ALTER TABLE ONLY public.gebruikers ALTER COLUMN email SET STATISTICS 1000;


--
-- Name: incoming_emails; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.incoming_emails (
    id character varying(255) NOT NULL,
    message_id character varying(255),
    "from" character varying(255) NOT NULL,
    "to" character varying(255) NOT NULL,
    subject character varying(255) NOT NULL,
    body text,
    content_type character varying(255),
    received_at timestamp without time zone NOT NULL,
    uid character varying(255),
    account_type character varying(50),
    is_processed boolean DEFAULT false NOT NULL,
    processed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: TABLE incoming_emails; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.incoming_emails IS 'Monitor for unprocessed email buildup';


--
-- Name: migraties; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.migraties (
    id bigint NOT NULL,
    versie text NOT NULL,
    naam text NOT NULL,
    toegepast timestamp with time zone
);


--
-- Name: migraties_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.migraties_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: migraties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.migraties_id_seq OWNED BY public.migraties.id;


--
-- Name: newsletters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.newsletters (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    subject text NOT NULL,
    content text NOT NULL,
    sent_at timestamp with time zone,
    batch_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    type character varying(50) NOT NULL,
    priority character varying(20) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    sent boolean DEFAULT false NOT NULL,
    sent_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: TABLE notifications; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.notifications IS 'Stores notifications to be sent via Telegram';


--
-- Name: partners; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.partners (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    logo text,
    website text,
    tier text,
    since date,
    visible boolean DEFAULT true NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permissions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    resource character varying(100) NOT NULL,
    action character varying(50) NOT NULL,
    description text,
    is_system_permission boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.photos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    url text NOT NULL,
    alt_text text,
    visible boolean DEFAULT true NOT NULL,
    thumbnail_url text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    title text,
    description text,
    year integer,
    cloudinary_folder text
);


--
-- Name: program_schedule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.program_schedule (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "time" text NOT NULL,
    event_description text NOT NULL,
    category text,
    icon_name text,
    order_number integer,
    visible boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    latitude numeric(10,8),
    longitude numeric(11,8)
);


--
-- Name: radio_recordings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.radio_recordings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    date text,
    audio_url text,
    thumbnail_url text,
    visible boolean DEFAULT true NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    revoked_at timestamp without time zone,
    is_revoked boolean DEFAULT false
);


--
-- Name: TABLE refresh_tokens; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.refresh_tokens IS 'Stores refresh tokens for JWT authentication with 7-day expiry';


--
-- Name: COLUMN refresh_tokens.token; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.token IS 'Base64 encoded random token (32 bytes)';


--
-- Name: COLUMN refresh_tokens.expires_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.expires_at IS 'Token expiration timestamp (7 days from creation)';


--
-- Name: COLUMN refresh_tokens.is_revoked; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.is_revoked IS 'Whether the token has been revoked (for token rotation)';


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_permissions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    assigned_at timestamp with time zone DEFAULT now(),
    assigned_by uuid
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    is_system_role boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    created_by uuid
);


--
-- Name: route_funds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.route_funds (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    route character varying(50) NOT NULL,
    amount integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT route_funds_amount_check CHECK ((amount >= 0))
);


--
-- Name: social_embeds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.social_embeds (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    platform text NOT NULL,
    embed_code text NOT NULL,
    order_number integer,
    visible boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: social_links; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.social_links (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    platform text NOT NULL,
    url text NOT NULL,
    bg_color_class text,
    icon_color_class text,
    order_number integer,
    visible boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: sponsors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sponsors (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    logo_url text,
    website_url text,
    order_number integer,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    visible boolean DEFAULT true NOT NULL
);


--
-- Name: title_section_content; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.title_section_content (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    event_title text NOT NULL,
    event_subtitle text,
    image_url text,
    image_alt text,
    detail_1_title text,
    detail_1_description text,
    detail_2_title text,
    detail_2_description text,
    detail_3_title text,
    detail_3_description text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    participant_count integer DEFAULT 0
);


--
-- Name: under_construction; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.under_construction (
    id integer NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    title text,
    message text,
    footer_text text,
    logo_url text,
    expected_date timestamp with time zone,
    social_links jsonb,
    progress_percentage integer,
    contact_email text,
    newsletter_enabled boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: under_construction_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.under_construction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: under_construction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.under_construction_id_seq OWNED BY public.under_construction.id;


--
-- Name: uploaded_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uploaded_images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    public_id text NOT NULL,
    url text NOT NULL,
    secure_url text NOT NULL,
    filename text NOT NULL,
    size bigint NOT NULL,
    mime_type text NOT NULL,
    width integer,
    height integer,
    folder text NOT NULL,
    thumbnail_url text,
    deleted_at timestamp without time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    assigned_at timestamp with time zone DEFAULT now(),
    assigned_by uuid,
    expires_at timestamp with time zone,
    is_active boolean DEFAULT true NOT NULL
);


--
-- Name: user_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.user_permissions AS
 SELECT ur.user_id,
    u.email,
    r.name AS role_name,
    p.resource,
    p.action,
    rp.assigned_at AS permission_assigned_at,
    ur.assigned_at AS role_assigned_at
   FROM ((((public.user_roles ur
     JOIN public.roles r ON ((ur.role_id = r.id)))
     JOIN public.role_permissions rp ON ((r.id = rp.role_id)))
     JOIN public.permissions p ON ((rp.permission_id = p.id)))
     JOIN public.gebruikers u ON ((ur.user_id = u.id)))
  WHERE (ur.is_active = true)
  ORDER BY ur.user_id, r.name, p.resource, p.action;


--
-- Name: v_user_role_migration_status; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_user_role_migration_status AS
 SELECT g.id AS user_id,
    g.email,
    g.naam,
    g.rol AS legacy_role,
    COALESCE(string_agg((r.name)::text, ', '::text ORDER BY (r.name)::text), 'GEEN RBAC ROL'::text) AS rbac_roles,
    count(ur.id) AS rbac_role_count,
        CASE
            WHEN (count(ur.id) = 0) THEN 'MISSING RBAC'::text
            WHEN ((g.rol IS NULL) OR ((g.rol)::text = ''::text)) THEN 'NO LEGACY'::text
            WHEN (EXISTS ( SELECT 1
               FROM (public.user_roles ur2
                 JOIN public.roles r2 ON ((ur2.role_id = r2.id)))
              WHERE ((ur2.user_id = g.id) AND (lower((r2.name)::text) = lower((g.rol)::text)) AND (ur2.is_active = true)))) THEN 'MIGRATED'::text
            ELSE 'MISMATCH'::text
        END AS migration_status
   FROM ((public.gebruikers g
     LEFT JOIN public.user_roles ur ON (((g.id = ur.user_id) AND (ur.is_active = true))))
     LEFT JOIN public.roles r ON ((ur.role_id = r.id)))
  GROUP BY g.id, g.email, g.naam, g.rol
  ORDER BY
        CASE
            WHEN (count(ur.id) = 0) THEN 1
            ELSE 2
        END, g.email;


--
-- Name: videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.videos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    video_id text NOT NULL,
    url text NOT NULL,
    title text,
    description text,
    thumbnail_url text,
    visible boolean DEFAULT true NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: migraties id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migraties ALTER COLUMN id SET DEFAULT nextval('public.migraties_id_seq'::regclass);


--
-- Name: under_construction id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.under_construction ALTER COLUMN id SET DEFAULT nextval('public.under_construction_id_seq'::regclass);


--
-- Data for Name: aanmelding_antwoorden; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.aanmelding_antwoorden (id, aanmelding_id, verzonden_op, created_at, updated_at, tekst, verzond_op, verzond_door, email_verzonden, verzonden_door) FROM stdin;
\.


--
-- Data for Name: aanmeldingen; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.aanmeldingen (id, naam, email, telefoon, status, created_at, updated_at, rol, afstand, ondersteuning, bijzonderheden, terms, email_verzonden, email_verzonden_op, behandeld_door, behandeld_op, notities, test_mode, steps, gebruiker_id, antwoorden_count) FROM stdin;
3e62d5d3-070d-47b1-a1ef-30665f982789	TGTest	laventejeffrey@gmail.com	06123456789	nieuw	2025-03-23 17:06:26.132297	2025-11-02 02:09:47.251906	Begeleider	15 KM	Anders	Telegram Test bericht - officiele weg	t	f	\N	\N	\N	\N	f	483	8f333073-8d5a-4202-9627-02863000b822	0
db3ec762-dd54-4ba7-98ab-981235cc316a	Mila Veenendaal	gaminggirlayla@gmail.com	\N	nieuw	2025-03-26 16:26:53.777384	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-03-26 16:30:57.269	\N	\N	\N	f	0	d12f1566-e096-4b80-b717-ca40036ab4a2	0
d17c16c6-c423-43de-a876-d40326b62d9e	Ayla Toprak	gamergirlayla@gmail.com	\N	nieuw	2025-03-26 16:25:43.756211	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-03-26 16:30:56.685	\N	\N	\N	f	0	a64fb32b-7cdc-4279-89d2-a7e9daeb1f05	0
917206a7-a28d-4bad-8b41-ed127eab743a	A. Bistolfi	nedarg@icloud.com	\N	nieuw	2025-03-26 12:27:20.236848	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-03-26 16:30:56.062	\N	\N	\N	f	0	01f63c58-2311-41ff-94ba-b727437ee555	0
61a3f823-82f6-4107-a9a9-b6a18e6b12c3	Henk Rekers 	h.rekers59@kpnmail.nl	\N	nieuw	2025-04-14 20:03:47.361483	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-15 13:23:40.051	\N	\N	\N	f	0	03dff890-db8e-4899-8275-ddbf6421550b	0
8a5c2c59-ebca-411e-9471-a38be47e1192	Hilde Rekers 	h.rekers59@kpnmail.nl	\N	nieuw	2025-04-14 20:00:47.479484	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-15 13:23:39.203	\N	\N	\N	f	0	03dff890-db8e-4899-8275-ddbf6421550b	0
282e6ec6-2c97-4b46-a992-273c826c1f91	Albert 	diesbosje@hotmail.com	\N	nieuw	2025-04-12 07:08:26.565956	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-12 16:52:18.027	\N	\N	\N	f	0	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	0
d10065f9-229d-4a1d-9178-324bdffa57c0	Diesmer 	diesbosje@hotmail.com	0613429612	nieuw	2025-04-12 07:06:59.872906	2025-11-02 02:09:47.251906	Begeleider	6 KM	Nee		t	t	2025-04-12 16:52:17.255	\N	\N	\N	f	0	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	0
613770dd-4733-4b54-964d-c3753cfebfd7	Sylvia Dijkstra	sylvia.dijkstra@sheerenloo.nl	0683081728 	nieuw	2025-04-06 18:48:14.993723	2025-11-02 02:09:47.251906	Begeleider	2.5 KM	Nee		t	t	2025-04-10 11:15:46.282	\N	\N	\N	f	0	3f58d1fa-4472-4d62-a31f-6906732b95f0	0
2aaffc78-4dca-4b06-9a6a-3799cbfbe67b	Noa hiddes	Klaskehiddes@gmail.com	\N	nieuw	2025-03-31 11:57:23.536871	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-04-01 17:11:04.313	\N	\N	\N	f	0	e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	0
200f862d-4d80-4ccd-8f06-bf795be727fb	Anneke van de Glind 	Klaskehiddes@gmail.com	\N	nieuw	2025-03-31 11:56:37.582458	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-04-01 17:11:05.257	\N	\N	\N	f	0	e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	0
990f67e0-6964-45de-8a6b-fdfb78e62907	john	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 15:00:09.409909	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-07 15:00:14.363894	admin@dekoninklijkeloop.nl	2025-10-05 17:36:37.260733	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
3c420d85-5f78-4645-88dd-c9e90eb8b6f7	Henk Rekers 	h.rekers1959@kpnmail.nl	\N	nieuw	2025-04-12 10:21:37.470747	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee	email mislukt	t	f	\N	\N	\N	\N	f	0	cb74dac2-bd74-460f-a83f-04fe02cd9363	0
02f9605b-8b26-4461-9ca7-ed7ebbd12311	Theun 	diesbosje@hotmail.com		nieuw	2025-04-12 07:09:19.884142	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-12 16:52:19.056	\N	\N	\N	f	778	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	0
7420d8c4-d879-4782-9c86-51bc6a91c035	Pieter Streefland	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 17:50:19.857981	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-07 17:50:25.456495	admin@dekoninklijkeloop.nl	2025-10-05 17:36:32.459421	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
dc225be5-5077-4eba-b43b-4a92e7f13141	Yunis 	lidaahmadi99@gmail.com		nieuw	2025-04-29 19:02:00.775439	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-29 19:02:05.494646	\N	\N	\N	f	0	564546ae-3227-471c-a743-9713746567f0	0
a8b3fd40-45fe-44a9-9461-9859ebda9708	Ismail 	lidaahmadi99@gmail.com		nieuw	2025-04-29 19:02:35.68749	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-29 19:02:40.487906	\N	\N	\N	f	0	564546ae-3227-471c-a743-9713746567f0	0
9116f64d-50c0-42c7-bf0e-d000c8504a2f	Mayke Rood	rood1960@gmail.com		nieuw	2025-04-29 21:16:05.431876	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-04-29 21:16:13.611673	\N	\N	\N	f	0	0647c6b8-ec19-4d1f-89de-c95ecd498e9b	0
af9b1795-e728-4848-bfcb-b88ddcc4ad6e	laura	steun94@xs4all.nl		nieuw	2025-04-30 10:03:24.549914	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-04-30 10:03:30.57633	\N	\N	\N	f	0	2f84da50-3f8d-44ce-a245-e7870cbc779f	0
b5e67c64-6bfa-46d4-ae40-98b093c8b720	Salih	topraks@gmail.com	\N	nieuw	2025-03-17 19:55:07.379647	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-03-21 06:47:38.608	\N	\N	\N	f	0	a76943e8-9673-41c0-8e89-b947294882d7	0
b2fd3412-8368-409f-8029-b2cdd581ade1	Manuela van zwam	benjaminlaan.64a@sheerenloo.nl	\N	nieuw	2025-03-11 11:28:05.483427	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-03-21 06:47:39.928	\N	\N	\N	f	0	eadc2084-7bed-4e3f-9c74-22ce7b5e4715	0
391f63c5-f034-466e-8a1f-ba9d06ed1192	Joyce Thielen	Joyce.thielen@sheerenloo.nl		nieuw	2025-03-09 16:52:09.564437	2025-11-02 02:09:47.251906	Begeleider	6 KM	Nee		t	t	2025-03-21 06:47:41.071	\N	\N	\N	f	0	00a26c26-4932-4091-91ba-4385a251e285	0
391f2579-d7cb-4ef3-afbe-14dc4115c519	Dick van Norden	Enckerkamp.27@sheerenloo.nl	\N	nieuw	2025-03-08 10:28:37.053379	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-03-08 10:28:42.2	\N	\N	\N	f	0	2aebf5e7-1cfe-4b43-afac-73a27ee49e00	0
90e477cc-89d1-4524-8edc-b697be8c504d	Angelo van Ingen	Enckerkamp.27@sheerenloo.nl	\N	nieuw	2025-03-08 10:27:31.078013	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-03-08 10:27:36.787	\N	\N	\N	f	0	2aebf5e7-1cfe-4b43-afac-73a27ee49e00	0
f4fc2312-ec8a-4dfc-90b5-a8da317618e6	Martin van der Wal	mjvdwal@hotmail.com		nieuw	2025-01-28 21:58:37.55756	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Ja	loopt samen met Dirk-Jan mee als vrijwilliger	t	t	2025-01-28 21:58:39.243	\N	\N	\N	f	0	ba656c7c-a049-46e6-a1db-7259e722fad6	0
4bfe814b-e0b8-4e60-9f46-fe38852d9ecb	Dirk-Jan Hempe	mjvdwal@hotmail.com		nieuw	2025-01-28 21:56:10.099964	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-01-28 21:56:12.004	\N	\N	\N	f	0	ba656c7c-a049-46e6-a1db-7259e722fad6	0
51855fec-eab9-494a-9321-c40d22da4ffc	Mirjam Kerkvliet	mirjam.kerkvliet@gmail.com	\N	nieuw	2025-03-24 17:27:25.191642	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	f	\N	\N	\N	\N	f	0	322b1f43-6d6f-433b-a1ea-ee6382a80bde	0
47775742-8950-4b94-9dd1-571ff4902688	Arno Kerkvliet	arno.kerkvliet@gmail.com	\N	nieuw	2025-03-24 17:26:12.1724	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	f	\N	\N	\N	\N	f	0	12df00b8-34e9-4768-b010-6a8a4855783d	0
d92ed75c-c275-47a4-88a9-ff7a4106f8ee	Jean-paul Hup	molenkamp.19@sheerenloo.nl	\N	nieuw	2025-03-24 09:17:59.726501	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	f	\N	\N	\N	\N	f	0	24f08c69-5be4-4d25-a21b-5ddcdf45cfc5	0
ecb8332b-ea39-4611-9f58-64921226f2a6	Annerieke Mandemaker-Timmer	annerieketimmer@hotmail.com	06 17 37 28 40 	nieuw	2025-03-24 09:16:42.111808	2025-11-02 02:09:47.251906	Begeleider	6 KM	Nee		t	f	\N	\N	\N	\N	f	0	4b6bffd8-3a13-45d3-b284-41bb9c7d101c	0
1ca80f61-f5c1-431f-b224-e6557150b65b	Han van Doornik	LaanvanGS.26@sheerenloo.nl	\N	nieuw	2025-03-30 08:07:45.334762	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Ja	Ik wil wel graag begeleiding 	t	f	\N	\N	\N	\N	f	0	554e0410-d07d-4b38-bc7e-2b1a110f3802	0
8c10628b-d36a-4f59-ac91-36ec1b045acd	TEST	laventejeffrey@gmail.com		nieuw	2025-04-17 12:30:46.920866	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	f	\N	\N	\N	\N	f	0	8f333073-8d5a-4202-9627-02863000b822	0
85bfe17f-9e09-4b8a-93af-d8eaecc26c9b	JeffTest	laventejeffrey@gmail.com		nieuw	2025-10-03 15:39:33.342306	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-10-03 15:39:39.34666	admin@dekoninklijkeloop.nl	2025-10-03 18:44:49.029869	\N	f	0	8f333073-8d5a-4202-9627-02863000b822	0
84cfd753-d77f-4b98-986b-f19665f6d7b2	jeffrey	laventejeffrey@gmail.com		nieuw	2025-05-18 11:31:00.038079	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-05-18 11:31:05.468884	admin@dekoninklijkeloop.nl	2025-10-05 17:36:27.643568	\N	f	0	8f333073-8d5a-4202-9627-02863000b822	0
e699a38a-b715-4452-afa4-358ac70e4163	yassine	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 18:53:45.366708	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-07 18:53:50.89824	admin@dekoninklijkeloop.nl	2025-10-05 17:36:29.74378	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
275490c0-1021-4bf4-9005-7df9884b0fe6	Bas heijenk 	basheijenk96@gmail.com	\N	nieuw	2025-03-22 16:43:03.19496	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	f	\N	\N	\N	\N	f	0	88e61c0b-b432-48c4-84df-2a7796f2f5a2	0
2499686e-2e62-4827-9079-78b468cb26c9	Bertram tijsma	Klaskehiddes@gmail.com	\N	nieuw	2025-03-29 10:07:12.393466	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-03-29 10:07:57.838	\N	\N	\N	f	0	e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	0
9f75b1df-4c72-4e36-9901-0f74cc26574f	Klaske van de glind	Klaskehiddes@gmail.com	\N	nieuw	2025-03-29 10:05:04.276906	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-03-29 10:07:57.215	\N	\N	\N	f	0	e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	0
721d297d-d670-46b1-a581-bcc095565bbd	Hilde Rekers 	h.rekers1959@kpnmail.nl	\N	nieuw	2025-04-12 10:20:25.399025	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee	email mislukt	t	f	\N	\N	\N	\N	f	0	cb74dac2-bd74-460f-a83f-04fe02cd9363	0
51205069-a20a-4231-bcb9-cb2d6fd042c8	Manuela van Zwam	rik.van-harxen@sheerenloo.nl	\N	nieuw	2025-03-23 14:08:53.328628	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Ja	Vaste begeleider die meeloopt	t	f	\N	\N	\N	\N	f	0	610f724f-69ed-4457-ab31-1c2572417fc7	0
26ea058b-2608-49d0-862a-611e98d7dc61	Janny van de Wall	mjvdwal@hotmail.com	\N	nieuw	2025-02-20 11:55:30.818333	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-02-20 11:55:33.048	\N	\N	\N	f	0	ba656c7c-a049-46e6-a1db-7259e722fad6	0
9f464844-8c93-4190-90f0-e74765c7f09a	Karin de Jong	karin.de.jong82@outlook.com	\N	nieuw	2025-03-24 18:47:05.053651	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	f	\N	\N	\N	\N	f	0	5ae279bf-2d9a-43b9-a3c6-98220b3cf92a	0
75f66868-9cd7-4612-b4e4-c888d74e4805	jeffreyTEST	laventejeffrey@gmail.com		nieuw	2025-04-17 12:43:19.042926	2025-11-02 02:09:47.251906	Deelnemer	2.5 KM	Nee		t	t	2025-04-17 12:43:24.496139	\N	\N	\N	f	0	8f333073-8d5a-4202-9627-02863000b822	0
fe5074fc-bb67-4c01-be49-49ae7b454c04	Sophie Hubers	sophieehubers@gmail.com		nieuw	2025-04-29 09:28:08.562089	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-04-29 09:28:14.444904	\N	\N	\N	f	0	7fb1db03-7d34-4b42-b046-9a24b5b00068	0
23b85646-d492-4972-8204-b3c86122274b	Geer Hubers	hub3008@gmail.com		nieuw	2025-04-29 09:28:53.416007	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-04-29 09:28:58.185455	\N	\N	\N	f	0	604ce292-214d-4b3a-86b4-d6cda2e51af0	0
7cc6bcff-70c9-427e-988f-72993b4477be	Saleem	lidaahmadi99@gmail.com	0649020648	nieuw	2025-04-29 18:55:23.435938	2025-11-02 02:09:47.251906	Begeleider	6 KM	Nee		t	t	2025-04-29 18:55:28.865905	\N	\N	\N	f	0	564546ae-3227-471c-a743-9713746567f0	0
80815a1d-562e-435c-9d92-cf4d67e1da8c	Yussef	lidaahmadi99@gmail.com		nieuw	2025-04-29 18:59:57.785545	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-29 19:00:02.567367	\N	\N	\N	f	0	564546ae-3227-471c-a743-9713746567f0	0
45bc3588-191e-4aaf-a534-dc73e7d792f4	Sulleeman 	lidaahmadi99@gmail.com		nieuw	2025-04-29 19:00:39.16856	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-04-29 19:00:44.727312	\N	\N	\N	f	0	564546ae-3227-471c-a743-9713746567f0	0
6b5493ef-c680-443d-95d5-a8670082337f	Luciënne 	luciennetang@gmail.com		nieuw	2025-04-30 20:55:07.613343	2025-11-02 02:09:47.251906	Deelnemer	15 KM	Nee		t	t	2025-04-30 20:55:13.465866	\N	\N	\N	f	0	c9366f76-ec43-44ce-92a0-ea00a4d44b04	0
4f8aad03-39b7-4fbb-8e6a-0aff29bf6140	Theodora Naus	tgemooi@gmail.com		nieuw	2025-05-01 10:14:27.911394	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-01 10:14:33.699276	\N	\N	\N	f	0	ce9f1f44-48af-4d26-987d-cc726cb4dc22	0
c9c742b9-381c-41eb-ba5e-7da0f8205fcc	Theodora Naus	tgemooi@gmail.com		nieuw	2025-05-01 10:23:07.753249	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-01 10:23:13.628827	\N	\N	\N	f	0	ce9f1f44-48af-4d26-987d-cc726cb4dc22	0
381abc94-0241-4f0f-94ad-cf96506fd352	michel	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 14:58:37.893809	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-07 14:58:43.536661	admin@dekoninklijkeloop.nl	2025-10-05 17:36:43.710127	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
20da5ef6-077a-48f5-ab34-c301ac4fdb6e	peter	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 14:59:16.617602	2025-11-02 02:09:47.251906	Deelnemer	10 KM	Nee		t	t	2025-05-07 14:59:21.398644	admin@dekoninklijkeloop.nl	2025-10-05 17:36:40.418285	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
46c4c9b9-67fa-4698-85c0-ed75ec983535	Danny 	enckerkamp.23@sheerenloo.nl		nieuw	2025-05-07 15:00:53.079124	2025-11-02 02:09:47.251906	Deelnemer	6 KM	Nee		t	t	2025-05-07 15:00:57.977426	admin@dekoninklijkeloop.nl	2025-10-05 17:36:34.746766	\N	f	0	49b72de1-ec0f-4af0-80cd-e6286d320417	0
\.


--
-- Data for Name: album_photos; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.album_photos (id, album_id, photo_id, order_number, created_at) FROM stdin;
10592b5a-9f16-48d2-8e95-150e5c5c53e5	d51cff45-b958-4370-a983-51e650ffa43e	ee20de98-2fa8-4e23-8bf1-0b705b55aa7c	30	2025-05-19 06:42:03.513117+00
aa902fe8-9ef7-44db-bbe1-1a48c203a26c	d51cff45-b958-4370-a983-51e650ffa43e	ec209a88-3f3b-4168-a4c9-1c9275edcffb	9	2025-05-17 20:19:11.664743+00
e3734c7c-a99d-421f-9a0c-ac02f08cbaf8	d51cff45-b958-4370-a983-51e650ffa43e	be16d3ae-c8a4-457b-9806-fee919bc79a2	17	2025-05-17 20:19:11.664743+00
1fa5d9f4-942d-4465-ad93-77ed1098d1dd	d51cff45-b958-4370-a983-51e650ffa43e	ac4fdceb-3598-4ba9-9d61-4f40ef0e44e3	18	2025-05-17 20:19:11.664743+00
c80655d5-0091-40bc-95cb-b1fcdab0112e	d51cff45-b958-4370-a983-51e650ffa43e	c5cc7e20-9e4d-4e09-bcf5-0c2edcb83360	19	2025-05-17 20:19:11.664743+00
c64d8242-6200-4c19-b55e-2eb7ca8d13d9	d51cff45-b958-4370-a983-51e650ffa43e	fded98ba-c7fa-4c42-ad2e-6c7afea5b8fc	20	2025-05-17 20:19:11.664743+00
486dcb08-ba33-47a7-af91-612d56e98639	d51cff45-b958-4370-a983-51e650ffa43e	9db233dc-58ff-4352-be6e-fb0d3b30798a	21	2025-05-17 20:19:11.664743+00
e7de00f4-fc44-4ff8-abd7-4d9c079c97fb	d51cff45-b958-4370-a983-51e650ffa43e	7d525fa5-192a-4b5f-9f03-7b8f9d4eb038	22	2025-05-17 20:19:11.664743+00
74069e01-a97a-45a0-a82e-4c3baf231e23	d51cff45-b958-4370-a983-51e650ffa43e	096c8248-2b05-4337-b256-f857419fa9bd	23	2025-05-17 20:19:11.664743+00
ea79a259-5dee-4f64-8009-00ddeae4c636	d51cff45-b958-4370-a983-51e650ffa43e	2debaea6-d5ab-4dd3-ae1b-15919c80245c	24	2025-05-17 20:19:11.664743+00
09101c7c-e7ac-4f99-af78-63d06f343d94	d51cff45-b958-4370-a983-51e650ffa43e	7eb75484-b2cb-41a7-bdc6-b16a8811ecca	25	2025-05-17 20:19:11.664743+00
65f4c380-e20c-42da-8a06-70d85d3f3ec6	d51cff45-b958-4370-a983-51e650ffa43e	19312aba-29ab-4169-83eb-934a80763ad8	26	2025-05-17 20:19:11.664743+00
bf01850a-e7d1-43c2-a614-7e1a65f3a8fc	d51cff45-b958-4370-a983-51e650ffa43e	db8dbb96-c528-4ae2-adf7-fc0cd165ad3c	27	2025-05-17 20:19:11.664743+00
df7ce81b-0480-4955-ad42-cf17dcef84ab	d51cff45-b958-4370-a983-51e650ffa43e	eafe6902-3d8e-4929-8cf3-c9b49429adac	28	2025-05-17 20:19:11.664743+00
898ea918-4be5-4c7f-8e6b-5d8b73198974	d51cff45-b958-4370-a983-51e650ffa43e	754ceb60-f4d8-4434-9338-065339e636e8	29	2025-05-17 20:19:11.664743+00
ae9666e2-5867-4d92-b25b-b6e665414dc7	d51cff45-b958-4370-a983-51e650ffa43e	7c296b82-d35d-4065-8523-dffc976688f7	2	2025-05-17 20:19:11.664743+00
9ce0e5fc-b7ea-4f63-bf63-c5c0f720ed61	d51cff45-b958-4370-a983-51e650ffa43e	08362e92-340a-432a-b306-153ad27ee686	1	2025-05-17 20:19:11.664743+00
e519317a-c8a1-47ff-a935-fc4e4c8c6afc	d51cff45-b958-4370-a983-51e650ffa43e	1350bb1c-b132-4250-82b8-58efb05fb53b	3	2025-05-17 20:19:11.664743+00
28974549-324c-49ff-b441-82fe90c947cd	d51cff45-b958-4370-a983-51e650ffa43e	0911da4d-6169-4383-962b-9a640b5eac0e	4	2025-05-17 20:19:11.664743+00
d2f53994-5156-4fde-a63c-02e6c39c198e	d51cff45-b958-4370-a983-51e650ffa43e	baf57b95-0ab7-4217-af98-1453a9e6a938	5	2025-05-17 20:19:11.664743+00
0bf73fa0-f48d-4272-87aa-abb53a5afd5a	d51cff45-b958-4370-a983-51e650ffa43e	76278110-be4c-4ffc-bc35-e2d4b14fa42f	6	2025-05-17 20:19:11.664743+00
c7cad470-1808-46ab-99af-9b5681c8ca69	d51cff45-b958-4370-a983-51e650ffa43e	10fff5f8-2701-4f34-86f1-b063b252f35a	7	2025-05-17 20:19:11.664743+00
f5e86512-f99d-4700-bc32-c2f5b235a670	d51cff45-b958-4370-a983-51e650ffa43e	c7ddfbd8-6bf4-4bd9-9285-f6492b98bef4	8	2025-05-17 20:19:11.664743+00
7551cf97-06fe-4651-ab0b-538dab99f5c6	d51cff45-b958-4370-a983-51e650ffa43e	6649c67b-2eb5-4d29-9e1c-0373ea3b1771	10	2025-05-17 20:19:11.664743+00
4e9db229-bcff-4118-b846-3931301fc20c	d51cff45-b958-4370-a983-51e650ffa43e	49f983ab-ec63-4489-93ce-ba9272ba7f49	11	2025-05-17 20:19:11.664743+00
78872e14-ad35-4b12-92a7-8e35f59c448b	d51cff45-b958-4370-a983-51e650ffa43e	83e75346-54b4-40b7-8ba4-2d43b3d8c867	12	2025-05-17 20:19:11.664743+00
6fbfc586-b7ac-4af6-9053-13ee5cdc4f86	d51cff45-b958-4370-a983-51e650ffa43e	18eb6cd7-b9a4-4b21-8363-ef77f420ac09	13	2025-05-17 20:19:11.664743+00
17bf4b95-a2e2-46f9-830d-c86825aaddcc	d51cff45-b958-4370-a983-51e650ffa43e	991bfd41-1365-4bb4-9ac8-8519fc9bfb32	14	2025-05-17 20:19:11.664743+00
9e8ef436-84a7-4b42-b5a0-da23dd2caed6	d51cff45-b958-4370-a983-51e650ffa43e	e971f01a-1d4a-4400-bb54-3428fbc69a98	15	2025-05-17 20:19:11.664743+00
622a17b5-98a4-4655-8c7b-db075c1a7d2f	d51cff45-b958-4370-a983-51e650ffa43e	c199f984-64e5-4405-b23c-e6ff4a3eaed3	16	2025-05-17 20:19:11.664743+00
cb67b8ed-cb18-4266-b5dc-d8a8c7b0128b	72831c18-4c6c-4dc8-9a2a-ace696e8996b	26188a5b-d542-4674-8ae9-b63520fbd4b2	15	2025-04-19 00:31:58.466572+00
5db43327-8974-48c7-accc-33affa2bf999	72831c18-4c6c-4dc8-9a2a-ace696e8996b	dbb68d91-3f6f-46c3-b751-a13f34269b79	16	2025-04-19 00:31:58.466572+00
29bea012-0b00-408d-8ff2-beffcd286730	72831c18-4c6c-4dc8-9a2a-ace696e8996b	e7b84300-7158-475b-a79b-10e97b416d58	14	2025-04-19 00:31:58.466572+00
47bee588-1be9-43c5-afea-80d67e97a4c2	72831c18-4c6c-4dc8-9a2a-ace696e8996b	f4ce7f8e-b573-4602-9e12-69e0585df779	13	2025-04-19 00:31:58.466572+00
fb9d85b1-c9aa-4935-8c1d-848603223d5c	ce8df963-f118-4296-9f5c-e33308dc7bfa	d78d9d0f-29ac-42a1-951a-a0bcb355d227	10	2025-04-18 21:22:50.820523+00
c3035d06-0fff-421a-894b-0b8c804906ae	ce8df963-f118-4296-9f5c-e33308dc7bfa	245ddf1a-61ab-4b64-b8b2-13d5079d6592	5	2025-03-17 19:49:06.359018+00
6f2c18d9-b0aa-45d7-bf24-62aa24b4f309	ce8df963-f118-4296-9f5c-e33308dc7bfa	3acaeca2-aa51-4cd6-9fad-7b20b607cba7	9	2025-03-17 19:49:06.359018+00
d0a158e2-c89e-4866-9ccb-5493e5cdc445	ce8df963-f118-4296-9f5c-e33308dc7bfa	0334c63c-230a-46c7-b87d-3f4a5cc946c0	8	2025-03-17 19:49:06.359018+00
cfecedf5-7b64-4229-80ff-90dc3dc9901a	ce8df963-f118-4296-9f5c-e33308dc7bfa	4059d3cf-0e92-44c2-bbfa-93235dc19eae	7	2025-03-17 19:49:06.359018+00
ec507140-12fe-4d96-9b41-78da5b6d5ba0	ce8df963-f118-4296-9f5c-e33308dc7bfa	6fd90008-f1ca-4e3d-9915-32b914208239	6	2025-03-17 19:49:06.359018+00
c70e5b3f-197f-4fe1-86bb-cd3a2c5f4b65	ce8df963-f118-4296-9f5c-e33308dc7bfa	8a4d5c20-ea73-4336-8c6b-af7a197ef7c2	1	2025-03-17 19:49:06.359018+00
3227ab06-4171-433c-97db-727a0a59b873	ce8df963-f118-4296-9f5c-e33308dc7bfa	af027789-d718-4899-88b2-d8f0814b73ee	2	2025-03-17 19:49:06.359018+00
70055993-7e41-49c1-903e-ff69c8f4fee3	ce8df963-f118-4296-9f5c-e33308dc7bfa	e1350a21-3d42-4252-bc56-5d9ce264ff41	3	2025-03-17 19:49:06.359018+00
8c64e44a-37f8-4fb9-b1af-8195b6d2bf8e	ce8df963-f118-4296-9f5c-e33308dc7bfa	7431c540-2aa5-42db-b5c5-df55a17856ae	4	2025-03-17 19:49:06.359018+00
3d958630-2821-4633-8126-23fc9f881863	72831c18-4c6c-4dc8-9a2a-ace696e8996b	d78d9d0f-29ac-42a1-951a-a0bcb355d227	12	2025-02-14 14:36:41.258168+00
\.


--
-- Data for Name: albums; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.albums (id, title, description, cover_photo_id, visible, order_number, created_at, updated_at) FROM stdin;
72831c18-4c6c-4dc8-9a2a-ace696e8996b	Voorbereidingen 	Voor werk	dbb68d91-3f6f-46c3-b751-a13f34269b79	t	3	2025-02-14 14:35:32.369+00	2025-02-14 14:35:32.37+00
ce8df963-f118-4296-9f5c-e33308dc7bfa	DKL-2024	DKL 2024	8a4d5c20-ea73-4336-8c6b-af7a197ef7c2	t	2	2024-12-26 15:17:54.268+00	2025-02-17 13:32:50.793+00
d51cff45-b958-4370-a983-51e650ffa43e	DKL 2025	De koninklijke Loop 2025!	08362e92-340a-432a-b306-153ad27ee686	t	1	2025-05-17 20:10:00+00	2025-05-17 20:10:06.643082+00
\.


--
-- Data for Name: chat_channel_participants; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_channel_participants (id, channel_id, user_id, role, joined_at, last_seen_at, is_active, last_read_at) FROM stdin;
aa934d58-36ab-45df-b60b-a9e6c50d6e42	2f723afc-aa20-4887-8852-746521ed4c4c	7157f3f6-da85-4058-9d38-19133ec93b03	owner	2025-10-05 13:12:37.775883+00	0001-01-01 00:00:00+00	t	\N
7fb938a1-76d9-4b55-9458-58ad889807bf	598024a1-3d79-4411-9382-45a484152f15	7157f3f6-da85-4058-9d38-19133ec93b03	owner	2025-10-05 13:56:46.633076+00	0001-01-01 00:00:00+00	t	\N
76d26d8f-96b3-48cd-808a-c942adbb471c	9e6a284e-d881-425d-908b-5383e281543c	748320dd-5b5e-4434-ad7e-8a405fd6266f	member	2025-10-05 17:59:58.714365+00	0001-01-01 00:00:00+00	t	\N
f1df609e-f45e-41fc-817c-eae605405816	9e6a284e-d881-425d-908b-5383e281543c	7157f3f6-da85-4058-9d38-19133ec93b03	member	2025-10-05 17:59:58.715891+00	0001-01-01 00:00:00+00	t	\N
ae974248-1aef-4b13-a06d-dc4b64b473da	dca1aa32-e350-4bf9-b022-1de98abf0ba9	748320dd-5b5e-4434-ad7e-8a405fd6266f	owner	2025-10-05 18:01:12.807577+00	0001-01-01 00:00:00+00	t	\N
02ae827c-4625-4ad5-a233-4663c1bf17b0	dca1aa32-e350-4bf9-b022-1de98abf0ba9	7157f3f6-da85-4058-9d38-19133ec93b03	member	2025-10-05 18:27:18.591956+00	0001-01-01 00:00:00+00	t	\N
\.


--
-- Data for Name: chat_channels; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_channels (id, name, description, type, created_by, created_at, updated_at, is_active, is_public) FROM stdin;
2f723afc-aa20-4887-8852-746521ed4c4c	test		public	7157f3f6-da85-4058-9d38-19133ec93b03	2025-10-05 13:12:37.77017+00	2025-10-05 13:12:37.77017+00	t	f
598024a1-3d79-4411-9382-45a484152f15	tests		public	7157f3f6-da85-4058-9d38-19133ec93b03	2025-10-05 13:56:46.63066+00	2025-10-05 13:56:46.63066+00	t	f
9e6a284e-d881-425d-908b-5383e281543c	Direct chat		direct	748320dd-5b5e-4434-ad7e-8a405fd6266f	2025-10-05 17:59:58.712388+00	2025-10-05 17:59:58.712388+00	t	f
dca1aa32-e350-4bf9-b022-1de98abf0ba9	testmain		public	748320dd-5b5e-4434-ad7e-8a405fd6266f	2025-10-05 18:01:12.80599+00	2025-10-05 18:01:12.80599+00	t	t
\.


--
-- Data for Name: chat_message_reactions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_message_reactions (id, message_id, user_id, emoji, created_at) FROM stdin;
\.


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_messages (id, channel_id, user_id, content, message_type, file_url, file_name, file_size, reply_to_id, edited_at, created_at, updated_at, thumbnail_url) FROM stdin;
4106ed28-e046-4d1d-9936-50ccf09c4e10	2f723afc-aa20-4887-8852-746521ed4c4c	7157f3f6-da85-4058-9d38-19133ec93b03	yoo	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 14:22:00.988385+00	2025-10-05 14:22:00.988385+00	\N
518b6dc0-6627-4598-aed0-3ba62aaebd1e	2f723afc-aa20-4887-8852-746521ed4c4c	7157f3f6-da85-4058-9d38-19133ec93b03	test	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 15:06:31.033978+00	2025-10-05 15:06:31.033978+00	\N
5a36f3e2-beed-43d4-bc1b-9ad10e87a66a	598024a1-3d79-4411-9382-45a484152f15	7157f3f6-da85-4058-9d38-19133ec93b03	testw	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 15:44:30.578804+00	2025-10-05 15:44:30.578804+00	\N
a4535e0c-ee05-45a0-8754-5d2a176d0422	9e6a284e-d881-425d-908b-5383e281543c	7157f3f6-da85-4058-9d38-19133ec93b03	hey man	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 18:00:09.324025+00	2025-10-05 18:00:09.324025+00	\N
9e9e1e47-e076-4bef-9b28-cd3abfa15343	9e6a284e-d881-425d-908b-5383e281543c	7157f3f6-da85-4058-9d38-19133ec93b03	alles lekker?	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 18:00:11.586907+00	2025-10-05 18:00:11.586907+00	\N
47a69593-839a-4c87-b784-1105eb966eb7	dca1aa32-e350-4bf9-b022-1de98abf0ba9	7157f3f6-da85-4058-9d38-19133ec93b03	yoo	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 18:27:24.345357+00	2025-10-05 18:27:24.345357+00	\N
86180561-7cb9-413c-b72c-7229e962f349	dca1aa32-e350-4bf9-b022-1de98abf0ba9	748320dd-5b5e-4434-ad7e-8a405fd6266f	yoo	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 18:27:37.599983+00	2025-10-05 18:27:37.599983+00	\N
461605af-da7f-4688-9065-e2bf48f7dd5c	9e6a284e-d881-425d-908b-5383e281543c	7157f3f6-da85-4058-9d38-19133ec93b03	yoo	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 19:17:13.396528+00	2025-10-05 19:17:13.396528+00	\N
f65c48e9-7736-49ca-97e0-d39ea15b70ac	9e6a284e-d881-425d-908b-5383e281543c	748320dd-5b5e-4434-ad7e-8a405fd6266f	hey	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 19:17:21.071146+00	2025-10-05 19:17:21.071146+00	\N
2d1c76b5-6aac-4e37-bce9-dc883e924e34	9e6a284e-d881-425d-908b-5383e281543c	7157f3f6-da85-4058-9d38-19133ec93b03	yo man	text			0	\N	0001-01-01 00:00:00+00	2025-10-05 19:44:21.740962+00	2025-10-05 19:44:21.740962+00	\N
\.


--
-- Data for Name: chat_user_presence; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_user_presence (user_id, status, last_seen, updated_at) FROM stdin;
748320dd-5b5e-4434-ad7e-8a405fd6266f	online	2025-10-08 16:59:43.403442+00	2025-10-08 16:59:43.404032+00
11a3ec93-b159-473b-86d2-d3979b9c9e3a	offline	2025-10-05 17:57:19.031989+00	2025-10-05 17:57:19.033164+00
7157f3f6-da85-4058-9d38-19133ec93b03	online	2025-11-02 02:07:14.725872+00	2025-11-02 02:07:14.739065+00
\.


--
-- Data for Name: contact_antwoorden; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.contact_antwoorden (id, contact_id, verzonden_op, created_at, updated_at, tekst, verzond_op, verzond_door, email_verzonden, verzonden_door) FROM stdin;
\.


--
-- Data for Name: contact_formulieren; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.contact_formulieren (id, naam, email, bericht, status, created_at, updated_at, email_verzonden, email_verzonden_op, privacy_akkoord, behandeld_door, behandeld_op, notities, beantwoord, antwoord_tekst, antwoord_datum, antwoord_door, test_mode, antwoorden_count) FROM stdin;
78910428-f760-485d-ae57-653db478ca35	Bas heijenk 	basheijenk96@gmail.com	Hallo ik heb me op gegeven maar ik kan helaas niet sorry 	nieuw	2025-03-23 07:15:20.8851	2025-11-02 02:09:47.251906	f	\N	t	\N	\N	\N	f	\N	\N	\N	f	0
6ce2e9a3-59fd-4430-aa5c-66df48fbd695	je geheime liefde	de.konining@willem.alexander.nl	Gedeelte doneren doet het niet.	nieuw	2025-01-28 00:03:05.830054	2025-11-02 02:09:47.251906	t	2025-01-28 00:03:08.041	t	marieke@dekoninklijkeloop.nl	2025-02-05 02:23:22.97	\N	f	\N	\N	\N	f	0
\.


--
-- Data for Name: email_templates; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.email_templates (id, naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by) FROM stdin;
0dc14d46-1207-45d0-ac5e-b05ca79fd664	contact_admin_email	Nieuw contactformulier	<p>Er is een nieuw contactformulier ingevuld door {{.Contact.Naam}}.</p><p>Email: {{.Contact.Email}}</p><p>Bericht: {{.Contact.Bericht}}</p>	Email die naar de admin wordt gestuurd bij een nieuw contactformulier	t	2025-03-14 15:22:28.710911	2025-03-14 15:22:28.710911	7157f3f6-da85-4058-9d38-19133ec93b03
7ca5166d-75ae-4474-a76b-460c1eee68fb	contact_email	Bedankt voor je bericht	<p>Beste {{.Contact.Naam}},</p><p>Bedankt voor je bericht. We nemen zo snel mogelijk contact met je op.</p>	Bevestigingsemail die naar de gebruiker wordt gestuurd bij een contactformulier	t	2025-03-14 15:22:28.710911	2025-03-14 15:22:28.710911	7157f3f6-da85-4058-9d38-19133ec93b03
eb568086-205f-44a5-9edd-0680c7f6501b	aanmelding_admin_email	Nieuwe aanmelding ontvangen	<p>Er is een nieuwe aanmelding ontvangen van {{.Aanmelding.Naam}}.</p><p>Email: {{.Aanmelding.Email}}</p>	Email die naar de admin wordt gestuurd bij een nieuwe aanmelding	t	2025-03-14 15:22:28.710911	2025-03-14 15:22:28.710911	7157f3f6-da85-4058-9d38-19133ec93b03
52e2d3c8-32df-4a0b-bcb9-8e28f293da7f	aanmelding_email	Bedankt voor je aanmelding	<p>Beste {{.Aanmelding.Naam}},</p><p>Bedankt voor je aanmelding. We hebben je aanmelding ontvangen en zullen deze zo snel mogelijk verwerken.</p>	Bevestigingsemail die naar de gebruiker wordt gestuurd bij een aanmelding	t	2025-03-14 15:22:28.710911	2025-03-14 15:22:28.710911	7157f3f6-da85-4058-9d38-19133ec93b03
\.


--
-- Data for Name: gebruikers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed, role_id) FROM stdin;
748320dd-5b5e-4434-ad7e-8a405fd6266f	Jeffrey	jeffrey@dekoninklijkeloop.nl	$2a$10$SxGtzqId5ZwGhvhmU4ys0O4tzEhsi3HYljq3ObRaEjppjohtta.2a	staff	t	2025-11-01 06:47:34.468919	2025-03-14 19:40:14.755163	2025-11-01 06:47:34.472282	f	\N
0197cfc3-7ca2-403b-ae4d-32627cd47222	Lida	Lida@dekoninklijkeloop.nl	$2a$10$mrJh4LXjKO6/bcC/ZPJdQ.TQFGLUhiQU4XPJvbjebmnJ/eu5KgjZm	staff	t	2025-11-02 14:52:04.124096	2025-11-02 12:54:43.931013	2025-11-02 14:52:04.123658	f	\N
7157f3f6-da85-4058-9d38-19133ec93b03	SuperAdmin	admin@dekoninklijkeloop.nl	$2a$10$rA7tcRZw3s9uiMnb3lgqOee4gwQnixQDohxjqdnsBQGYBBjQDP2sC	admin	t	2025-11-02 14:53:13.200357	2025-03-14 15:22:28.710911	2025-11-02 14:53:13.199847	f	\N
b1bec8cb-1709-420f-88cb-fc74e0c6eec2	Salih	Salih@dekoninklijkeloop.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	staff	t	\N	2024-12-28 02:43:22	2025-10-05 17:15:03.315274	f	\N
11a3ec93-b159-473b-86d2-d3979b9c9e3a	Marieke	Marieke@dekoninklijkeloop.nl	$2a$10$AJUC49EoFrusB9GOGP1aHucLItd5OZfKMWwgF5dsNwKRvF/lPcxse	staff	t	2025-10-05 17:42:17.59877	2024-12-27 15:40:16.942007	2025-10-05 17:42:17.599415	f	\N
e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	Theun 	diesbosje@hotmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	2025-10-26 10:11:05.061707	2025-10-25 11:58:19.59437	2025-10-26 10:11:05.062215	f	\N
0647c6b8-ec19-4d1f-89de-c95ecd498e9b	Mayke Rood	rood1960@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
4b6bffd8-3a13-45d3-b284-41bb9c7d101c	Annerieke Mandemaker-Timmer	annerieketimmer@hotmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	begeleider	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
12df00b8-34e9-4768-b010-6a8a4855783d	Arno Kerkvliet	arno.kerkvliet@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
88e61c0b-b432-48c4-84df-2a7796f2f5a2	Bas heijenk 	basheijenk96@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
eadc2084-7bed-4e3f-9c74-22ce7b5e4715	Manuela van zwam	benjaminlaan.64a@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
7fb1db03-7d34-4b42-b046-9a24b5b00068	Sophie Hubers	sophieehubers@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
49b72de1-ec0f-4af0-80cd-e6286d320417	yassine	enckerkamp.23@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
2aebf5e7-1cfe-4b43-afac-73a27ee49e00	Dick van Norden	Enckerkamp.27@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
a64fb32b-7cdc-4279-89d2-a7e9daeb1f05	Ayla Toprak	gamergirlayla@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
d12f1566-e096-4b80-b717-ca40036ab4a2	Mila Veenendaal	gaminggirlayla@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
2f84da50-3f8d-44ce-a245-e7870cbc779f	laura	steun94@xs4all.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
03dff890-db8e-4899-8275-ddbf6421550b	Henk Rekers 	h.rekers59@kpnmail.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
604ce292-214d-4b3a-86b4-d6cda2e51af0	Geer Hubers	hub3008@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
00a26c26-4932-4091-91ba-4385a251e285	Joyce Thielen	Joyce.thielen@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	begeleider	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
5ae279bf-2d9a-43b9-a3c6-98220b3cf92a	Karin de Jong	karin.de.jong82@outlook.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	Noa hiddes	Klaskehiddes@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
554e0410-d07d-4b38-bc7e-2b1a110f3802	Han van Doornik	LaanvanGS.26@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
564546ae-3227-471c-a743-9713746567f0	Ismail 	lidaahmadi99@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
c9366f76-ec43-44ce-92a0-ea00a4d44b04	Luciënne 	luciennetang@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
322b1f43-6d6f-433b-a1ea-ee6382a80bde	Mirjam Kerkvliet	mirjam.kerkvliet@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
ba656c7c-a049-46e6-a1db-7259e722fad6	Janny van de Wall	mjvdwal@hotmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
24f08c69-5be4-4d25-a21b-5ddcdf45cfc5	Jean-paul Hup	molenkamp.19@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
01f63c58-2311-41ff-94ba-b727437ee555	A. Bistolfi	nedarg@icloud.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
610f724f-69ed-4457-ab31-1c2572417fc7	Manuela van Zwam	rik.van-harxen@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
3f58d1fa-4472-4d62-a31f-6906732b95f0	Sylvia Dijkstra	sylvia.dijkstra@sheerenloo.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	begeleider	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
ce9f1f44-48af-4d26-987d-cc726cb4dc22	Theodora Naus	tgemooi@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	\N	2025-10-25 11:58:19.59437	2025-10-25 12:09:50.868231	f	\N
a76943e8-9673-41c0-8e89-b947294882d7	Salih	topraks@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	2025-10-25 17:31:42.977083	2025-10-25 11:58:19.59437	2025-10-25 17:31:42.977658	f	\N
cb74dac2-bd74-460f-a83f-04fe02cd9363	Henk Rekers 	h.rekers1959@kpnmail.nl	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	deelnemer	t	2025-10-25 17:33:42.84163	2025-10-25 11:58:19.59437	2025-10-25 17:33:42.842175	f	\N
8f333073-8d5a-4202-9627-02863000b822	JeffTest	laventejeffrey@gmail.com	$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku	socialmedia	t	2025-11-02 12:58:24.606474	2025-10-25 11:58:19.59437	2025-11-02 12:58:24.609733	f	\N
\.


--
-- Data for Name: incoming_emails; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.incoming_emails (id, message_id, "from", "to", subject, body, content_type, received_at, uid, account_type, is_processed, processed_at, created_at, updated_at) FROM stdin;
371a5a8e-4d29-48a5-ac32-d88ae98fdcc1	<AS4PR10MB57982D1680DB871E678DBCC4EBE42@AS4PR10MB5798.EURPRD10.PROD.OUTLOOK.COM>	"Kleverwal, Lindsay" <lindsay.kleverwal@sheerenloo.nl>	info@dekoninklijkeloop.nl	18 mei	--_004_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_\r\nContent-Type: multipart/alternative;\r\n\tboundary="_000_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_"\r\n\r\n--_000_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_\r\nContent-Type: text/plain; charset="iso-8859-1"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nGoedemorgen,\r\n\r\nWij hebben ons opgegeven voor de sponsorloop 18 mei. Klopt het dat je daar =\r\ngeen bevestigingsmail van krijgt? En (electrische) rolstoelen mee in de pen=\r\ndelbus? Of is het handiger als wij de taxi naar het startpunt worden gebrac=\r\nht?\r\n\r\n\r\nMet vriendelijke groet,\r\n\r\nLindsay Kleverwal\r\n\r\nBegeleider E\r\n\r\n06-10073648\r\n\r\n[Logo 's Heeren Loo]\r\n\r\nAchisomoglaan 350,\r\n\r\n7325 BS Apeldoorn.\r\n\r\n\r\n________________________________\r\nDisclaimer\r\n--\r\nDe informatie verzonden met dit e-mailbericht (en bijlagen) is uitsluitend =\r\nbestemd voor de geadresseerde(n) en zij die van de geadresseerde(n) toestem=\r\nming hebben dit bericht te lezen. Gebruik door anderen dan geadresseerde(n)=\r\n is verboden. De informatie in dit e-mailbericht (en de bijlagen) kan vertr=\r\nouwelijk van aard zijn en kan binnen het bereik vallen van een geheimhoudin=\r\ngsplicht.\r\n's Heeren Loo is niet aansprakelijk voor schade ten gevolge van het gebruik=\r\n van elektronische middelen van communicatie, daaronder begrepen -maar niet=\r\n beperkt tot- schade ten gevolge van niet aflevering of vertraging bij de a=\r\nflevering van elektronische berichten, onderschepping of manipulatie van el=\r\nektronische berichten door derden of door programmatuur/apparatuur gebruikt=\r\n voor elektronische communicatie en overbrenging van virussen en andere kwa=\r\nadaardige programmatuur.\r\n\r\n--_000_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_\r\nContent-Type: text/html; charset="iso-8859-1"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html>\r\n<head>\r\n<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Diso-8859-=\r\n1">\r\n<style type=3D"text/css" style=3D"display:none;"> P {margin-top:0;margin-bo=\r\nttom:0;} </style>\r\n</head>\r\n<body dir=3D"ltr">\r\n<div class=3D"elementToProof" style=3D"margin: 0px; font-family: Verdana, G=\r\neneva, sans-serif; font-size: 10pt; color: rgb(0, 0, 0);">\r\nGoedemorgen,</div>\r\n<div style=3D"margin: 0px; font-family: Verdana, Geneva, sans-serif; font-s=\r\nize: 10pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div class=3D"elementToProof" style=3D"margin: 0px; font-family: Verdana, G=\r\neneva, sans-serif; font-size: 10pt; color: rgb(0, 0, 0);">\r\nWij hebben ons opgegeven voor de sponsorloop 18 mei. Klopt het dat je daar =\r\ngeen bevestigingsmail van krijgt? En (electrische) rolstoelen mee in de pen=\r\ndelbus? Of is het handiger als wij de taxi naar het startpunt worden gebrac=\r\nht?</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Verdana, Geneva, sans-s=\r\nerif; font-size: 10pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div id=3D"Signature">\r\n<p style=3D"text-align: left; text-indent: 0px; line-height: 13pt; backgrou=\r\nnd-color: rgb(255, 255, 255); margin-top: 0pt; margin-bottom: 0pt;">\r\n<span style=3D"font-family: Verdana; font-size: 10pt; color: rgb(0, 0, 0);"=\r\n>Met vriendelijke groet,<br>\r\n<br>\r\n</span><span style=3D"font-family: Verdana; font-size: 13.3333px;"><b>Linds=\r\nay </b></span><span style=3D"font-family: Verdana; font-size: 10pt; color: =\r\nrgb(0, 0, 0);"><b>Kleverwal</b></span></p>\r\n<p style=3D"text-align: left; text-indent: 0px; line-height: 13pt; backgrou=\r\nnd-color: rgb(255, 255, 255); margin-top: 0pt; margin-bottom: 0pt;">\r\n<span style=3D"font-family: Verdana; font-size: 10pt; color: rgb(0, 0, 0);"=\r\n>Begeleider E</span></p>\r\n<p style=3D"text-align: left; text-indent: 0px; line-height: 13pt; backgrou=\r\nnd-color: rgb(255, 255, 255); margin-top: 0pt; margin-bottom: 0pt;">\r\n<span style=3D"font-family: Verdana; font-size: 10pt; color: rgb(0, 0, 0);"=\r\n>06-10073648<br>\r\n<br>\r\n<img alt=3D"Logo 's Heeren Loo" id=3D"image_0" style=3D"width: 4.42cm; heig=\r\nht: 2.38cm; max-width: 657px; margin: 0px; vertical-align: top;" data-outlo=\r\nok-trace=3D"F:1|T:1" src=3D"cid:af7ec6a8-1069-4103-8697-3c9072a8bbc4"><br>\r\n<br>\r\nAchisomoglaan 350,</span></p>\r\n<p style=3D"text-align: left; text-indent: 0px; line-height: 10pt; backgrou=\r\nnd-color: white; margin: 0px;">\r\n<span style=3D"font-family: Verdana, Geneva, sans-serif; font-size: 10pt; c=\r\nolor: rgb(32, 31, 30); line-height: normal;">7325 BS Apeldoorn.</span></p>\r\n<p style=3D"text-align: left; text-indent: 0px; line-height: 13pt; backgrou=\r\nnd-color: rgb(255, 255, 255); margin-top: 0pt; margin-bottom: 0pt;">\r\n<span style=3D"font-family: Verdana; font-size: 10pt; color: rgb(0, 0, 0);"=\r\n><br>\r\n</span></p>\r\n</div>\r\n<hr>\r\n<font face=3D"Arial" color=3D"Gray" size=3D"1">Disclaimer<br>\r\n-- <begin disclaimer=3D"" text=3D""><br>\r\nDe informatie verzonden met dit e-mailbericht (en bijlagen) is uitsluitend =\r\nbestemd voor de geadresseerde(n) en zij die van de geadresseerde(n) toestem=\r\nming hebben dit bericht te lezen. Gebruik door anderen dan geadresseerde(n)=\r\n is verboden. De informatie in dit\r\n e-mailbericht (en de bijlagen) kan vertrouwelijk van aard zijn en kan binn=\r\nen het bereik vallen van een geheimhoudingsplicht.<br>\r\n's Heeren Loo is niet aansprakelijk voor schade ten gevolge van het gebruik=\r\n van elektronische middelen van communicatie, daaronder begrepen -maar niet=\r\n beperkt tot- schade ten gevolge van niet aflevering of vertraging bij de a=\r\nflevering van elektronische berichten,\r\n onderschepping of manipulatie van elektronische berichten door derden of d=\r\noor programmatuur/apparatuur gebruikt voor elektronische communicatie en ov=\r\nerbrenging van virussen en andere kwaadaardige programmatuur.<br>\r\n</font>\r\n</body>\r\n</html>\r\n\r\n--_000_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_--\r\n\r\n--_004_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_\r\nContent-Type: image/png; name="Outlook-Logo .png"\r\nContent-Description: Outlook-Logo .png\r\nContent-Disposition: inline; filename="Outlook-Logo .png"; size=21380;\r\n\tcreation-date="Tue, 07 May 2024 07:57:25 GMT";\r\n\tmodification-date="Tue, 07 May 2024 07:57:25 GMT"\r\nContent-ID: <af7ec6a8-1069-4103-8697-3c9072a8bbc4>\r\nContent-Transfer-Encoding: base64\r\n\r\n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\r\nRK5CYII=\r\n\r\n--_004_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_--\r\n	multipart/related; boundary="_004_AS4PR10MB57982D1680DB871E678DBCC4EBE42AS4PR10MB5798EURP_"; type="multipart/alternative"	2024-05-07 07:57:25	9	info	t	2025-04-20 00:29:13.05147	2025-04-20 00:03:01.021095	2025-04-20 00:29:13.05192
8594309f-0caf-4531-9d35-405bba8dcc79	<20240510105455.2AB4D6802B34E@n8nfibz.plsk.shared.prod.hostnetbv.nl>	noreply@hostnet-websitebuilder.nl	info@dekoninklijkeloop.nl	=?UTF-8?b?Q29udGFjdGZvcm11bGllcg==?=	Naam : angelique Wildeman\r\nE-mail : angelique.wildeman@sheerenloo.nl\r\nHeb_je_een_vraag_of_opmerking_neem_dan_contact_op_met_ons : Goedemiddag,\r\n\r\nEen aantal weken geleden heb ik twee clienten opgegeven voor de loop. Wesley Visschers en Wesley van  Schijndel. Maar ik heb tot op heden geen bevestiging van de inschrijving ontvangen. Kan ik er wel gewoon vanuit gaan dat ik met hun de loop kan gaan doen?\r\n\r\nMet vriendelijke groet,\r\n\r\nAngelique Wildeman\r\n	text/plain; charset=utf-8	2024-05-10 12:54:55	10	info	t	2025-04-20 00:29:36.397425	2025-04-20 00:03:01.026267	2025-04-20 00:29:36.399131
506b701f-728f-4679-acd7-75b951cd108f	<wtm.7eed2f23-33e9-46d5-a6b1-4bd9a92457b0@wetransfer.com>	WeTransfer <noreply@wetransfer.com>	info@dekoninklijkeloop.nl	vanvarikg@gmail.com heeft je Matinee gestuurd via WeTransfer	\r\n----==_mimepart_6644b652229b6_b93e4930359\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nvanvarikg@gmail.com heeft je Matinee gestuurd\r\n\r\n1 item, 132 MB in totaal =E3=83=BB Verloopt op 22 mei 2024\r\n\r\n\r\nDownloadlink:\r\nhttps://wetransfer.com/downloads/a15442eb1f5a1fa74910807d7c98073f2024051513=\r\n1815/1ccb70731a2931a080eab0beb226465120240515131837/42de49?trk=3DTRN_TDL_01=\r\n&utm_campaign=3DTRN_TDL_01&utm_medium=3Demail&utm_source=3Dsendgrid\r\n\r\n\r\n\r\nBericht:\r\nHallo,\r\n\r\nHierbij het radio interview.\r\n\r\nMet vriendelijke groet,\r\nRTV Apeldoorn\r\n\r\n\r\n\r\n1 item\r\n\r\nmatinee 1.mp3 - 132 MB\r\n\r\n\r\n[Haal meer uit WeTransfer met Pro](https://wetransfer.com/pricing?trk=3DTRN=\r\n_TDL_01&utm_campaign=3DTRN_TDL_01&utm_medium=3Demail&utm_source=3Dsendgrid)\r\n\r\n\r\n\r\nOver WeTransfer: https://wetransfer.com/about\r\nHulp: https://wetransfer.zendesk.com/hc/en-us\r\nAlgemene voorwaarden: https://wetransfer.com/legal/terms\r\n\r\n\r\nVoeg noreply@wetransfer.com toe aan [je contactpersonen](https://wetransfer=\r\n.zendesk.com/hc/en-us/articles/204909429) zodat je zeker weet dat onze emai=\r\nls aankomen.\r\n\r\n\r\n----==_mimepart_6644b652229b6_b93e4930359\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.=\r\nw3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">\r\n<html xmlns=3D"http://www.w3.org/1999/xhtml" xmlns=3D"http://www.w3.org/199=\r\n9/xhtml" xmlns:v=3D"urn:schemas-microsoft-com:vml" xmlns:o=3D"urn:schemas-m=\r\nicrosoft-com:office:office">\r\n\r\n    <head>\r\n    <meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Dutf-8=\r\n" />\r\n    <meta name=3D"viewport" content=3D"width=3Ddevice-width, initial-scale=\r\n=3D1.0" />\r\n    <!--[if !mso]><!-- -->\r\n    <meta http-equiv=3D"x-ua-compatible" content=3D"ie=3Dedge" />\r\n    <!--<![endif]-->\r\n    <meta name=3D"apple-mobile-web-app-capable" content=3D"yes" />\r\n    <meta name=3D"apple-mobile-web-app-status-bar-style" content=3D"black" =\r\n/>\r\n    <meta name=3D"format-detection" content=3D"telephone=3Dno" />\r\n    <meta http-equiv=3D"x-rim-auto-match" content=3D"none" />\r\n    <!-- Fonts, excluded from MSO -->\r\n    <!--[if !mso]><!-- -->\r\n    <style type=3D"text/css" data-premailer=3D"ignore">\r\n        /* Font-family declarations must not be inlined */\r\n        @font-face {\r\n            font-family: 'Fakt Pro';\r\n            src: url(/fonts/faktpro/FaktProWeb-Normal.woff);\r\n            font-weight: normal;\r\n            font-style: normal;\r\n        }\r\n        @font-face {\r\n            font-family: 'Fakt Pro Medium';\r\n            src: url(/fonts/faktpro/FaktProWeb-Medium.woff);\r\n            font-weight: normal;\r\n            font-style: normal;\r\n        }\r\n    </style>\r\n    <!--<![endif]-->\r\n   =20\r\n   =20\r\n    <style type=3D"text/css" data-premailer=3D"ignore">\r\n        /* client specific styles - DO NOT EDIT */\r\n        @-ms-viewport {\r\n            width: device-width;\r\n        }\r\n        @media only screen and (max-width: 600px) {\r\n\r\n        body {\r\n            padding: 0 !important;\r\n            width: auto !important;\r\n            min-width: auto !important;\r\n            max-width: auto !important;\r\n        }\r\n\r\n        *[class~=3Dinner_wrapper_table],\r\n        *[class~=3Dinner_wrapper_td] {\r\n            height: auto !important;\r\n            max-width: 320px !important;\r\n            min-width: 0 !important;\r\n            width: 100% !important;\r\n        }\r\n\r\n        *[class~=3Dpadded_mobile] {\r\n            padding-left: 30px !important;\r\n            padding-right: 30px !important;\r\n        }\r\n\r\n        *[class~=3Dunpadded_mobile] {\r\n            padding-left: 0 !important;\r\n            padding-right: 0 !important;\r\n        }\r\n\r\n        td[class~=3Dmain_heading_td_profile_picture_visible]{\r\n            font-size: 19px !important;\r\n        }\r\n\r\n        td[class~=3Dlogo_td] {\r\n            padding-top: 18px !important;\r\n        }\r\n\r\n        td[class~=3Dlogo_outer_wrapper_td] {\r\n            padding-top: 2px !important;\r\n            padding-bottom: 1px !important;\r\n        }\r\n\r\n        td[class~=3Dmain_heading_td] {\r\n            font-size: 19px !important;\r\n            line-height: 24px !important;\r\n            padding-top: 40px !important;\r\n        }\r\n\r\n        td[class~=3Davatar_outer_wrapper_td] {\r\n            padding-top: 40px !important;\r\n        }\r\n        td[class~=3Dbutton_outer_wrapper_td] {\r\n            padding-top: 30px !important;\r\n        }\r\n        td[class~=3Dbody_content_td] {\r\n            padding-top: 30px !important;\r\n            padding-bottom: 0px !important;\r\n        }\r\n        td[class~=3Dseparator_20_outer_wrapper_td] {\r\n            padding-top: 30px !important;\r\n            padding-bottom: 0px !important;\r\n        }\r\n        td[class~=3Dcall_to_action_td] {\r\n            padding-top: 0px !important;\r\n            padding-bottom: 0px !important;\r\n        }\r\n        td[class~=3Dseparator_40_outer_wrapper_td] {\r\n            padding-top: 30px !important;\r\n            padding-bottom: 0px !important;\r\n        }\r\n        td[class~=3Dbutton_outer_wrapper_td_keep_it_button] {\r\n            padding-top: 30px !important;\r\n            padding-bottom: 0px !important;\r\n        }\r\n\r\n        *[class~=3Davatar_outer_wrapper_td] img {\r\n            height: 80px !important;\r\n            width: 80px !important;\r\n        }\r\n\r\n        *[class~=3Dmonster] img {\r\n            height: 100px !important;\r\n            width: 100px !important;\r\n        }\r\n\r\n        td[class~=3Dfiles_details_td] {\r\n            padding-top: 10px !important;\r\n            color: #6A6D70;\r\n            font-size: 12px !important;\r\n            mso-line-height-rule: exactly;\r\n            /* this needs to proceed the line-height property */\r\n            line-height: 19px;\r\n        }\r\n\r\n        td[class~=3Dfiles_list] {\r\n            padding-top: 30px !important;\r\n            font-size: 14px !important;\r\n            mso-line-height-rule: exactly;\r\n            /* this needs to proceed the line-height property */\r\n            line-height: 23px !important;\r\n        }\r\n\r\n        span[class~=3Dbody_content_subheading_span] {\r\n            font-size: 16px;\r\n        }\r\n\r\n        table[class~=3Dbutton_table] {\r\n            margin-left: auto !important;\r\n            margin-right: auto !important;\r\n            width: auto !important;\r\n        }\r\n\r\n        td[class~=3Dsignoff_td] {\r\n            font-size: 16px !important;\r\n            line-height: 24px !important;\r\n            padding-top: 30px !important;\r\n            padding-bottom: 40px !important;\r\n        }\r\n\r\n        td[class~=3Dadd_our_email_td] {\r\n            font-size: 12px !important;\r\n            line-height: 19px !important;\r\n            padding-top: 15px !important;\r\n            padding-bottom: 15px !important;\r\n        }\r\n\r\n       td[class~=3Dfooter_td] {\r\n            font-size: 12px !important;\r\n            line-height: 23px !important;\r\n            padding-top: 20px !important;\r\n            padding-bottom: 30px !important;\r\n\r\n        }\r\n        td[class~=3Dbody_content_padding_bottom_td] {\r\n            padding-bottom: 40px !important;\r\n        }\r\n        td[class~=3Dfiles_list_title] {\r\n            padding-bottom: 0 !important;\r\n        }\r\n        td[class~=3Dfiles_list_content] {\r\n            padding-top: 0 !important;\r\n        }\r\n    }\r\n    </style>\r\n\r\n    <!--[if gte mso 9]>\r\n    <xml>\r\n      <o:OfficeDocumentSettings>\r\n        <o:AllowPNG/>\r\n        <o:PixelsPerInch>96</o:PixelsPerInch>\r\n      </o:OfficeDocumentSettings>\r\n    </xml>\r\n    <![endif]-->\r\n\r\n<style>.ReadMsgBody {\r\nwidth: 100%;\r\n}\r\n.ExternalClass {\r\nwidth: 100%;\r\n}\r\nbody {\r\n-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-sm=\r\noothing: antialiased; -moz-osx-font-smoothing: grayscale; margin: 0; outlin=\r\ne: none; padding: 0;\r\n}\r\nimg {\r\n-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-sm=\r\noothing: antialiased; -moz-osx-font-smoothing: grayscale; margin: 0; outlin=\r\ne: none; padding: 0;\r\n}\r\nimg {\r\nborder: none; display: block; height: auto; line-height: 100%; -ms-interpol=\r\nation-mode: bicubic; text-decoration: none;\r\n}\r\nbody {\r\nbackground-color: #f4f4f4;\r\n}\r\n.body_content_td a:active {\r\ncolor: #17181a; text-decoration: underline;\r\n}\r\n.footer_td p a:active {\r\ncolor: #17181a; text-decoration: underline;\r\n}\r\n.recipient_information a:active {\r\ncolor: #5268ff; font-weight: normal; text-decoration: none;\r\n}\r\n.add_our_email_td a:active {\r\ncolor: #797c7f; font-weight: normal; text-decoration: underline;\r\n}\r\n.email_without_default_client_style a:active {\r\ncolor: #797c7f; font-weight: normal; text-decoration: none !important;\r\n}\r\nbody.nu_body {\r\nbackground-color: #ffffff;\r\n}\r\n.nu_footer_td p a:active {\r\ncolor: #919599;\r\n}\r\n.nu_footer_marketing p a:active {\r\nfont-size: 16px; color: #ffffff; line-height: 22px;\r\n}\r\n</style></head>\r\n\r\n\r\n    <body style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: no=\r\nne; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale=\r\n; outline: none; margin: 0; padding: 0;" bgcolor=3D"#f4f4f4">\r\n        <!-- outer wrapper table -->\r\n            <!--[if !mso]><!-- -->\r\n        <table width=3D"100%" cellspacing=3D"0" cellpadding=3D"0" border=3D=\r\n"0" class=3D"outer_wrapper_table" style=3D"-ms-text-size-adjust: none; -web=\r\nkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-f=\r\nont-smoothing: grayscale; outline: none; border-collapse: collapse; border-=\r\nspacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: auto !i=\r\nmportant; margin: 0; padding: 0;" bgcolor=3D"#f4f4f4">\r\n            <tr>\r\n                <!-- separate background image for hotmail, with lots of tr=\r\nansparency added to the bottom -->\r\n                <td style=3D"-ms-text-size-adjust: none; -webkit-text-size-=\r\nadjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing:=\r\n grayscale; outline: none; width: 100%; margin: 0; padding: 0;" align=3D"le=\r\nft" valign=3D"top">\r\n\r\n                    <!--<![endif]-->\r\n\r\n                    <!--[if gte mso 9]>\r\n                    <table style=3D"width: 100%; margin: 0px; padding: 0px;=\r\n border-collapse: collapse; border-spacing: 0px; mso-table-lspace: 0pt; mso=\r\n-table-rspace: 0pt; background-color: transparent;" cellspacing=3D"0" cellp=\r\nadding=3D"0" border=3D"0">\r\n                      <tr>\r\n                        <td style=3D"width: 100%; margin: 0px; padding: 0px=\r\n; vertical-align: top; text-align: left;">\r\n\r\n                    <![endif]-->\r\n\r\n                    <center>\r\n                        <!-- inner wrapper table -->\r\n                        <table width=3D"600" align=3D"center" cellspacing=\r\n=3D"0" cellpadding=3D"0" border=3D"0" class=3D"inner_wrapper_table table_ce=\r\nntered" style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none=\r\n; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; =\r\noutline: none; border-collapse: collapse; border-spacing: 0; mso-table-lspa=\r\nce: 0; mso-table-rspace: 0; table-layout: fixed; width: 600px; min-width: 6=\r\n00px; margin: 0 auto; padding: 0;">\r\n                            <tr>\r\n                                <td width=3D"600" class=3D"inner_wrapper_td=\r\n" style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -web=\r\nkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outlin=\r\ne: none; width: 600px; min-width: 600px; margin: 0; padding: 0;" align=3D"l=\r\neft" valign=3D"top">\r\n\r\n                                    <table cellspacing=3D"0" cellpadding=3D=\r\n"0" border=3D"0" class=3D"hidden_teaser" style=3D"-ms-text-size-adjust: non=\r\ne; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -mo=\r\nz-osx-font-smoothing: grayscale; outline: none; border-collapse: collapse; =\r\nborder-spacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: =\r\nfixed; display: none; font-size: 0; max-height: 0; width: 0; margin: 0; pad=\r\nding: 0;">\r\n    <tr>\r\n        <td class=3D"hidden_teaser" style=3D"-ms-text-size-adjust: none; -w=\r\nebkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx=\r\n-font-smoothing: grayscale; outline: none; width: 0; display: none; font-si=\r\nze: 0; max-height: 0; margin: 0; padding: 0;" align=3D"left" valign=3D"top"=\r\n>\r\n          This transfer expires on 2024-05-22\r\n        </td>\r\n    </tr>\r\n</table>\r\n\r\n\r\n                                    <table cellspacing=3D"0" cellpadding=3D=\r\n"0" border=3D"0" class=3D"table_full_width" style=3D"-ms-text-size-adjust: =\r\nnone; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; =\r\n-moz-osx-font-smoothing: grayscale; outline: none; border-collapse: collaps=\r\ne; border-spacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layou=\r\nt: fixed; width: 100%; margin: 0; padding: 0;">\r\n    <tr>\r\n        <td class=3D"logo_outer_wrapper_td" style=3D"-ms-text-size-adjust: =\r\nnone; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; =\r\n-moz-osx-font-smoothing: grayscale; outline: none; width: 100%; margin: 0; =\r\npadding: 55px 0 0;" align=3D"left" valign=3D"top">\r\n            <table cellspacing=3D"0" cellpadding=3D"0" border=3D"0" class=\r\n=3D"table_full_width" style=3D"-ms-text-size-adjust: none; -webkit-text-siz=\r\ne-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothin=\r\ng: grayscale; outline: none; border-collapse: collapse; border-spacing: 0; =\r\nmso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; width: 100%;=\r\n margin: 0; padding: 0;">\r\n                <tr>\r\n                    <!--[if !mso]><!-- -->\r\n                    <td class=3D"logo_inner_wrapper_td" style=3D"-ms-text-s=\r\nize-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoothing: a=\r\nntialiased; -moz-osx-font-smoothing: grayscale; outline: none; width: 100%;=\r\n font-size: 10px; background: no-repeat center / cover; margin: 0; padding:=\r\n 0;" align=3D"left" bgcolor=3D"#5268ff" valign=3D"top">\r\n                        <!--<![endif]-->\r\n                        <!--[if mso]> <td> <![endif]-->\r\n\r\n                        <!-- background image for MSO -->\r\n                        <!--[if gte mso 9]>  <v:rect xmlns:v=3D"urn:schemas=\r\n-microsoft-com:vml" fill=3D"false" stroke=3D"false" style=3D"position: rela=\r\ntive; top: 0; left: 0; width: 600px; height: 60px;"> <v:textbox inset=3D"0,=\r\n0,0,0"> <![endif]-->\r\n                        <center>\r\n                            <table align=3D"center" cellspacing=3D"0" cellp=\r\nadding=3D"0" border=3D"0" class=3D"table_centered" style=3D"-ms-text-size-a=\r\ndjust: none; -webkit-text-size-adjust: none; -webkit-font-smoothing: antial=\r\niased; -moz-osx-font-smoothing: grayscale; outline: none; border-collapse: =\r\ncollapse; border-spacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; tabl=\r\ne-layout: fixed; width: auto; margin: 0 auto; padding: 0;">\r\n                                <!-- spacing tr/td for Yahoo mail and MSO -=\r\n->\r\n                                <tr>\r\n                                    <td height=3D"16px" style=3D"height: 16=\r\npx; -ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-fon=\r\nt-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none=\r\n; width: 100%; margin: 0; padding: 0;" align=3D"left" valign=3D"top"></td>\r\n                                </tr>\r\n                                <!-- ENDOF spacing tr/td for Yahoo mail and=\r\n MSO -->\r\n                                <tr>\r\n                                    <td style=3D"-ms-text-size-adjust: none=\r\n; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz=\r\n-osx-font-smoothing: grayscale; outline: none; width: 100%; margin: 0; padd=\r\ning: 0;" align=3D"left" valign=3D"top">\r\n                                        <a href=3D"https://wetransfer.com/?=\r\ntrk=3DTRN_TDL_01&amp;utm_campaign=3DTRN_TDL_01&amp;utm_medium=3Demail&amp;u=\r\ntm_source=3Dsendgrid" target=3D"_blank">\r\n                                            <img src=3D"https://email.wetra=\r\nnsfer.net/emails/logos/blue_2x.png" alt=3D"Click 'Download images' to view =\r\nimages" align=3D"center" border=3D"0" width=3D"56" class=3D"logo_blue_img" =\r\nstyle=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webki=\r\nt-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outline:=\r\n none; display: block; height: auto; line-height: 100%; -ms-interpolation-m=\r\node: bicubic; text-decoration: none; width: 56px; margin: 0; padding: 0; bo=\r\nrder-style: none;" />\r\n                                        </a>\r\n                                    </td>\r\n                                </tr>\r\n                                <!-- spacing tr/td for Yahoo mail and MSO -=\r\n->\r\n                                <tr>\r\n                                    <td height=3D"11px" style=3D"height: 11=\r\npx; -ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-fon=\r\nt-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none=\r\n; width: 100%; margin: 0; padding: 0;" align=3D"left" valign=3D"top"></td>\r\n                                </tr>\r\n                               <!-- ENDOF spacing tr/td for Yahoo mail and =\r\nMSO -->\r\n                            </table>\r\n                        </center>\r\n                        <!--[if gte mso 9]> </v:textbox> </v:rect> <![endif=\r\n]-->\r\n                    </td>\r\n                </tr>\r\n            </table>\r\n        </td>\r\n    </tr>\r\n</table>\r\n\r\n\r\n                                    <table cellspacing=3D"0" cellpadding=3D=\r\n"0" border=3D"0" class=3D"main_content_outer_wrapper_table" style=3D"-ms-te=\r\nxt-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoothin=\r\ng: antialiased; -moz-osx-font-smoothing: grayscale; outline: none; border-c=\r\nollapse: collapse; border-spacing: 0; mso-table-lspace: 0; mso-table-rspace=\r\n: 0; table-layout: fixed; width: 100%; margin: 0; padding: 0;" bgcolor=3D"#=\r\nffffff">\r\n    <tr>\r\n        <td style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: =\r\nnone; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: graysca=\r\nle; outline: none; width: 100%; margin: 0; padding: 0;" align=3D"left" vali=\r\ngn=3D"top">\r\n            <table cellspacing=3D"0" cellpadding=3D"0" border=3D"0" class=\r\n=3D"table_full_width" style=3D"-ms-text-size-adjust: none; -webkit-text-siz=\r\ne-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothin=\r\ng: grayscale; outline: none; border-collapse: collapse; border-spacing: 0; =\r\nmso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; width: 100%;=\r\n margin: 0; padding: 0;">\r\n                <tr>\r\n                    <td class=3D"padded_mobile main_content_inner_wrapper_t=\r\nd" style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -we=\r\nbkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outli=\r\nne: none; width: 100%; margin: 0; padding: 0;" align=3D"left" valign=3D"top=\r\n">\r\n                        <table cellspacing=3D"0" cellpadding=3D"0" border=\r\n=3D"0" class=3D"table_full_width" style=3D"-ms-text-size-adjust: none; -web=\r\nkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-f=\r\nont-smoothing: grayscale; outline: none; border-collapse: collapse; border-=\r\nspacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; =\r\nwidth: 100%; margin: 0; padding: 0;">\r\n                           =20\r\n<tr>\r\n<td class=3D"main_heading_td unpadded_mobile main_heading_td_wider" style=\r\n=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-fon=\r\nt-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none=\r\n; width: 100%; color: #17181a; font-size: 26px; font-style: normal; font-we=\r\night: normal; mso-line-height-rule: exactly; line-height: 30px; word-spacin=\r\ng: 0; margin: 0; padding: 60px 80px 0;" align=3D"center" valign=3D"top">\r\n<a class=3D"main_heading_email_link" href=3D"mailto:vanvarikg@gmail.com" st=\r\nyle=3D"color: #17181a; font-weight: normal; text-decoration: none;"><span c=\r\nlass=3D"main_heading_email_link" style=3D"color: #5268ff; font-weight: norm=\r\nal; text-decoration: none;">vanvarikg@gmail.com</span></a>\r\n<br />\r\n<span class=3D"transfer_display_name" style=3D"display: inline-flex;">\r\nheeft je Matinee gestuurd\r\n</span>\r\n</td>\r\n</tr>\r\n\r\n<tr>\r\n<td class=3D"files_details_td unpadded_mobile" style=3D"-ms-text-size-adjus=\r\nt: none; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiase=\r\nd; -moz-osx-font-smoothing: grayscale; outline: none; width: 100%; color: #=\r\n6a6d70; font-family: 'Fakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial,=\r\n sans-serif; font-size: 14px; font-style: normal; font-weight: normal; mso-=\r\nline-height-rule: exactly; line-height: 23px; word-spacing: 0; margin: 0; p=\r\nadding: 20px 80px 0;" align=3D"center" valign=3D"top">\r\n1 item, 132 MB in totaal\r\n=E3=83=BB Verloopt op 22 mei 2024\r\n</td>\r\n</tr>\r\n\r\n<tr>\r\n<td class=3D"body_content_td  unpadded_mobile" style=3D"-ms-text-size-adjus=\r\nt: none; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiase=\r\nd; -moz-osx-font-smoothing: grayscale; outline: none; width: 100%; color: #=\r\n797c7f; font-family: 'Fakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial,=\r\n sans-serif; font-size: 14px; font-style: normal; font-weight: normal; mso-=\r\nline-height-rule: exactly; line-height: 24px; word-spacing: 0; margin: 0; p=\r\nadding: 50px 80px 0;" align=3D"left" valign=3D"top">\r\n<p class=3D"message_content" style=3D"-ms-text-size-adjust: none; -webkit-t=\r\next-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-s=\r\nmoothing: grayscale; outline: none; background-color: #fafafa; color: #1718=\r\n1a; margin: -10px 0 0; padding: 20px;">\r\n<strong class=3D"transfer_display_name_with_message" style=3D"display: flex=\r\n;">Matinee\r\n</strong>\r\nHallo,\r\n<br />\r\n\r\n<br />\r\nHierbij het radio interview.\r\n<br />\r\n\r\n<br />\r\nMet vriendelijke groet,\r\n<br />\r\nRTV Apeldoorn\r\n</p>\r\n</td>\r\n</tr>\r\n\r\n<!--[if !mso]><!-- -->\r\n<tr>\r\n  <td class=3D"button_outer_wrapper_td_get unpadded_mobile" style=3D"-ms-te=\r\nxt-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoothin=\r\ng: antialiased; -moz-osx-font-smoothing: grayscale; outline: none; width: 1=\r\n00%; margin: 0; padding: 40px 160px;" align=3D"left" valign=3D"top">\r\n    <table cellspacing=3D"0" cellpadding=3D"0" border=3D"0" class=3D"table_=\r\nfull_width button_table" style=3D"-ms-text-size-adjust: none; -webkit-text-=\r\nsize-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-smoot=\r\nhing: grayscale; outline: none; border-collapse: collapse; border-spacing: =\r\n0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; width: 10=\r\n0%; margin: 0; padding: 0;">\r\n      <tr>\r\n        <td style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: =\r\nnone; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: graysca=\r\nle; outline: none; width: 100%; margin: 0; padding: 0;" align=3D"left" vali=\r\ngn=3D"top">\r\n          <a target=3D"_blank" class=3D"button_anchor button_2_anchor" href=\r\n=3D"https://wetransfer.com/downloads/a15442eb1f5a1fa74910807d7c98073f202405=\r\n15131815/1ccb70731a2931a080eab0beb226465120240515131837/42de49?trk=3DTRN_TD=\r\nL_01&amp;utm_campaign=3DTRN_TDL_01&amp;utm_medium=3Demail&amp;utm_source=3D=\r\nsendgrid" style=3D"background-color: #5268ff; color: #ffffff; display: bloc=\r\nk; font-size: 14px; font-style: normal; text-align: center; text-decoration=\r\n: none; word-spacing: 0; -mox-border-radius: 25px; -webkit-border-radius: 2=\r\n5px; o-border-radius: 25px; -ms-border-radius: 25px; border-radius: 25px; p=\r\nadding: 15px 20px;">\r\n            <span>\r\n                  Haal je bestanden op\r\n            </span>\r\n</a>        </td>\r\n      </tr>\r\n    </table>\r\n  </td>\r\n</tr>\r\n<!--<![endif]-->\r\n<!--[if mso]>\r\n<tr>\r\n  <td style=3D"padding: 40px 160px 0px 160px;">\r\n    <v:roundrect xmlns:v=3D"urn:schemas-microsoft-com:vml" xmlns:w=3D"urn:s=\r\nchemas-microsoft-com:office:word"\r\n  href=3D"https://wetransfer.com/downloads/a15442eb1f5a1fa74910807d7c98073f=\r\n20240515131815/1ccb70731a2931a080eab0beb226465120240515131837/42de49?trk=3D=\r\nTRN_TDL_01&amp;utm_campaign=3DTRN_TDL_01&amp;utm_medium=3Demail&amp;utm_sou=\r\nrce=3Dsendgrid" style=3D"height: 44px; v-text-anchor: middle; width: 270px;=\r\n" arcsize=3D"50%" stroke=3D"false" fillcolor=3D"#409FFF">\r\n      <w:anchorlock/>\r\n      <center style=3D"color: #FFFFFF; font-family: Arial, sans-serif; font=\r\n-size: 14px; font-weight: normal; text-decoration: none;">\r\n          Haal je bestanden op\r\n      </center>\r\n    </v:roundrect>\r\n  </td>\r\n</tr>\r\n<![endif]-->\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <tr>\r\n    <td class=3D"separator_40_outer_wrapper_td unpadded_mobile" style=3D"-m=\r\ns-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoo=\r\nthing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none; widt=\r\nh: 100%; margin: 0; padding: 40px 80px 0;" align=3D"left" valign=3D"top">\r\n        <table cellspacing=3D"0" cellpadding=3D"0" border=3D"0" class=3D"ta=\r\nble_full_width" style=3D"-ms-text-size-adjust: none; -webkit-text-size-adju=\r\nst: none; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: gra=\r\nyscale; outline: none; border-collapse: collapse; border-spacing: 0; mso-ta=\r\nble-lspace: 0; mso-table-rspace: 0; table-layout: fixed; width: 100%; margi=\r\nn: 0; padding: 0;">\r\n            <tr>\r\n                <td class=3D"separator_td" style=3D"-ms-text-size-adjust: n=\r\none; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -=\r\nmoz-osx-font-smoothing: grayscale; outline: none; width: 100%; border-botto=\r\nm-width: 2px; border-bottom-color: #f4f4f4; border-bottom-style: solid; fon=\r\nt-size: 1px; mso-line-height-rule: exactly; line-height: 0; margin: 0; padd=\r\ning: 0;" align=3D"left" valign=3D"top">=C2=A0</td>\r\n            </tr>\r\n        </table>\r\n    </td>\r\n</tr>\r\n\r\n                              <tr>\r\n<td class=3D"body_content_td unpadded_mobile download_link_container" style=\r\n=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-fon=\r\nt-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none=\r\n; width: 100%; color: #797c7f; font-family: 'Fakt Pro', 'Segoe UI', 'SanFra=\r\nncisco Display', Arial, sans-serif; font-size: 14px; font-style: normal; fo=\r\nnt-weight: normal; mso-line-height-rule: exactly; line-height: 24px; word-s=\r\npacing: 0; word-break: break-all; margin: 0; padding: 50px 80px 0;" align=\r\n=3D"left" valign=3D"top">\r\n<span class=3D"body_content_subheading_span" style=3D"color: #17181a; font-=\r\nsize: 18px; font-weight: 500;">\r\nDownloadlink\r\n</span>\r\n<br />\r\n<a class=3D"download_link_link" href=3D"https://wetransfer.com/downloads/a1=\r\n5442eb1f5a1fa74910807d7c98073f20240515131815/1ccb70731a2931a080eab0beb22646=\r\n5120240515131837/42de49?trk=3DTRN_TDL_01&amp;utm_campaign=3DTRN_TDL_01&amp;=\r\nutm_medium=3Demail&amp;utm_source=3Dsendgrid" style=3D"color: #17181a; text=\r\n-decoration: underline; font-weight: normal; word-wrap: break-word;"><span =\r\nclass=3D"download_link_link" style=3D"color: #5268ff; font-weight: normal; =\r\ntext-decoration: underline; word-wrap: break-word;">https://wetransfer.com/=\r\ndownloads/a15442eb1f5a1fa74910807d7c98073f20240515131815/1ccb70731a2931a080=\r\neab0beb226465120240515131837/42de49</span>\r\n</a>\r\n</td>\r\n</tr>\r\n\r\n  <tr>\r\n<td class=3D"body_content_td body_content_padding_bottom_td files_list file=\r\ns_list_title unpadded_mobile" style=3D"-ms-text-size-adjust: none; -webkit-=\r\ntext-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-font-=\r\nsmoothing: grayscale; outline: none; width: 100%; color: #797c7f; font-fami=\r\nly: 'Fakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial, sans-serif; font=\r\n-size: 14px; font-style: normal; font-weight: normal; mso-line-height-rule:=\r\n exactly; line-height: 24px; word-spacing: 0; margin: 0; padding: 50px 80px=\r\n 0;" align=3D"left" valign=3D"top">\r\n<span class=3D"body_content_subheading_span" style=3D"color: #17181a; font-=\r\nsize: 18px; font-weight: 500;">\r\n1 item\r\n</span>\r\n</td>\r\n</tr>\r\n<tr>\r\n<td class=3D"tp0 body_content_td body_content_padding_bottom_td files_list =\r\nfiles_list_content unpadded_mobile" style=3D"-ms-text-size-adjust: none; -w=\r\nebkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx=\r\n-font-smoothing: grayscale; outline: none; width: 100%; color: #797c7f; fon=\r\nt-family: 'Fakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial, sans-serif=\r\n; font-size: 14px; font-style: normal; font-weight: normal; mso-line-height=\r\n-rule: exactly; line-height: 24px; word-spacing: 0; margin: 0; padding: 0 8=\r\n0px 50px;" align=3D"left" valign=3D"top">\r\n<div class=3D"body_content_subheading_span" style=3D"color: #17181a; font-s=\r\nize: 18px; font-weight: 500;"></div>\r\n<div class=3D"transfer_item transfer_item_last" style=3D"border-bottom-widt=\r\nh: 1px; border-bottom-color: #f4f4f4; border-bottom-style: none; padding: 9=\r\npx 0 7px;">\r\n<div class=3D"transfer_item_title" style=3D"color: #17181a; font-family: 'F=\r\nakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial, sans-serif; font-size:=\r\n 14px; font-style: normal; font-weight: normal; mso-line-height-rule: exact=\r\nly; line-height: 16px; word-spacing: 0;">\r\nmatinee 1.mp3\r\n</div>\r\n<div class=3D"transfer_item_description" style=3D"color: #6a6d70; font-size=\r\n: 12px; mso-line-height-rule: exactly; line-height: 16px;">\r\n\r\n132 MB\r\n</div>\r\n</div>\r\n</td>\r\n</tr>\r\n\r\n\r\n\r\n                        </table>\r\n                    </td>\r\n                </tr>\r\n            </table>\r\n        </td>\r\n    </tr>\r\n</table>\r\n\r\n\r\n                                    <table cellspacing=3D"0" cellpadding=3D=\r\n"0" border=3D"0" class=3D"table_full_width" style=3D"-ms-text-size-adjust: =\r\nnone; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; =\r\n-moz-osx-font-smoothing: grayscale; outline: none; border-collapse: collaps=\r\ne; border-spacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layou=\r\nt: fixed; width: 100%; margin: 0; padding: 0;">\r\n    <tr>\r\n        <td class=3D"add_our_email_outer_wrapper_td" style=3D"-ms-text-size=\r\n-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoothing: anti=\r\naliased; -moz-osx-font-smoothing: grayscale; outline: none; width: 100%; ma=\r\nrgin: 0; padding: 2px 0 0;" align=3D"left" valign=3D"top">\r\n\r\n            <table cellspacing=3D"0" cellpadding=3D"0" border=3D"0" class=\r\n=3D"add_our_email_wrapper_table" style=3D"-ms-text-size-adjust: none; -webk=\r\nit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-fo=\r\nnt-smoothing: grayscale; outline: none; border-collapse: collapse; border-s=\r\npacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; w=\r\nidth: 100%; margin: 0; padding: 0;" bgcolor=3D"#ffffff">\r\n                <tr>\r\n                    <td class=3D"padded_mobile add_our_email_inner_wrapper_=\r\ntd" style=3D"-ms-text-size-adjust: none; -webkit-text-size-adjust: none; -w=\r\nebkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; outl=\r\nine: none; width: 100%; margin: 0; padding: 0 20px;" align=3D"left" valign=\r\n=3D"top">\r\n\r\n                        <table cellspacing=3D"0" cellpadding=3D"0" border=\r\n=3D"0" class=3D"table_full_width" style=3D"-ms-text-size-adjust: none; -web=\r\nkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-f=\r\nont-smoothing: grayscale; outline: none; border-collapse: collapse; border-=\r\nspacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layout: fixed; =\r\nwidth: 100%; margin: 0; padding: 0;">\r\n                            <tr>\r\n                                <td class=3D"add_our_email_td" style=3D"-ms=\r\n-text-size-adjust: none; -webkit-text-size-adjust: none; -webkit-font-smoot=\r\nhing: antialiased; -moz-osx-font-smoothing: grayscale; outline: none; width=\r\n: 100%; color: #797c7f; font-family: 'Fakt Pro', 'Segoe UI', 'SanFrancisco =\r\nDisplay', Arial, sans-serif; font-size: 12px; font-style: normal; font-weig=\r\nht: normal; mso-line-height-rule: exactly; line-height: 24px; word-spacing:=\r\n 0; margin: 0; padding: 13px 0;" align=3D"center" valign=3D"top">\r\n                                <p style=3D"-ms-text-size-adjust: none; -we=\r\nbkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-=\r\nfont-smoothing: grayscale; outline: none; margin: 0; padding: 0;">Voeg <a c=\r\nlass=3D"email_without_default_client_style" href=3D"mailto:noreply@wetransf=\r\ner.com" style=3D"color: #797c7f; font-weight: normal; text-decoration: none=\r\n !important;"><span class=3D"email_without_default_client_style" style=3D"c=\r\nolor: #797c7f; font-weight: normal; text-decoration: none !important;">nore=\r\nply@wetransfer.com</span></a> toe aan <a href=3D"https://wetransfer.zendesk=\r\n.com/hc/en-us/articles/204909429?utm_campaign=3DTRN_TDL_01&amp;utm_source=\r\n=3Dsendgrid&amp;utm_medium=3Demail&amp;trk=3DTRN_TDL_01" style=3D"color: #7=\r\n97c7f; font-weight: normal; text-decoration: underline;">je contactpersonen=\r\n</a> zodat je zeker weet dat onze emails aankomen.</p>\r\n\r\n                                </td>\r\n                            </tr>\r\n                        </table>\r\n\r\n                    </td>\r\n                </tr>\r\n            </table>\r\n\r\n        </td>\r\n    </tr>\r\n</table>\r\n\r\n\r\n                                    <table cellspacing=3D"0" cellpadding=3D=\r\n"0" border=3D"0" class=3D"table_full_width" style=3D"-ms-text-size-adjust: =\r\nnone; -webkit-text-size-adjust: none; -webkit-font-smoothing: antialiased; =\r\n-moz-osx-font-smoothing: grayscale; outline: none; border-collapse: collaps=\r\ne; border-spacing: 0; mso-table-lspace: 0; mso-table-rspace: 0; table-layou=\r\nt: fixed; width: 100%; margin: 0; padding: 0;">\r\n    <tr>\r\n        <td class=3D"footer_td" style=3D"-ms-text-size-adjust: none; -webki=\r\nt-text-size-adjust: none; -webkit-font-smoothing: antialiased; -moz-osx-fon=\r\nt-smoothing: grayscale; outline: none; width: 100%; color: #797c7f; font-fa=\r\nmily: 'Fakt Pro', 'Segoe UI', 'SanFrancisco Display', Arial, sans-serif; fo=\r\nnt-size: 12px; font-style: normal; font-weight: normal; mso-line-height-rul=\r\ne: exactly; line-height: 23px; word-spacing: 0; margin: 0; padding: 30px 20=\r\npx;" align=3D"center" valign=3D"top">\r\n\r\n\r\n            <a href=3D"https://wetransfer.com/about?trk=3DTRN_TDL_01&amp;ut=\r\nm_campaign=3DTRN_TDL_01&amp;utm_medium=3Demail&amp;utm_source=3Dsendgrid" c=\r\nlass=3D"footer_link" style=3D"color: #797c7f; font-weight: normal; text-dec=\r\noration: underline;">\r\n              <span class=3D"footer_link" style=3D"color: #797c7f; font-wei=\r\nght: normal; text-decoration: underline;">Over WeTransfer</span></a>\r\n\r\n            <span class=3D"footer_link_separator" style=3D"color: #797c7f;"=\r\n>=C2=A0=C2=A0=E3=83=BB=C2=A0=C2=A0</span>\r\n\r\n            <a href=3D"https://wetransfer.zendesk.com/hc/en-us?utm_campaign=\r\n=3DTRN_TDL_01&amp;utm_source=3Dsendgrid&amp;utm_medium=3Demail&amp;trk=3DTR=\r\nN_TDL_01" class=3D"footer_link" style=3D"color: #797c7f; font-weight: norma=\r\nl; text-decoration: underline;">\r\n                <span class=3D"footer_link" style=3D"color: #797c7f; font-w=\r\neight: normal; text-decoration: underline;">Hulp</span></a>\r\n\r\n            <span class=3D"footer_link_separator" style=3D"color: #797c7f;"=\r\n>=C2=A0=C2=A0=E3=83=BB=C2=A0=C2=A0</span>\r\n\r\n            <a href=3D"https://wetransfer.com/legal/terms?trk=3DTRN_TDL_01&=\r\namp;utm_campaign=3DTRN_TDL_01&amp;utm_medium=3Demail&amp;utm_source=3Dsendg=\r\nrid" class=3D"footer_link" style=3D"color: #797c7f; font-weight: normal; te=\r\nxt-decoration: underline;">\r\n                <span class=3D"footer_link" style=3D"color: #797c7f; font-w=\r\neight: normal; text-decoration: underline;">Algemene voorwaarden</span></a>\r\n\r\n                <span class=3D"footer_link_separator" style=3D"color: #797c=\r\n7f;">=C2=A0=C2=A0=E3=83=BB=C2=A0=C2=A0</span>\r\n  <a href=3D"https://safety.wetransfer.com/report?productUrl=3Dhttps://wetr=\r\nansfer.com/downloads/a15442eb1f5a1fa74910807d7c98073f20240515131815/cf9c55"=\r\n rel=3D"external" target=3D"_blank" class=3D"footer_link" style=3D"color: #=\r\n797c7f; font-weight: normal; text-decoration: underline;">\r\n                <span class=3D"footer_link" style=3D"color: #797c7f; font-w=\r\neight: normal; text-decoration: underline;">Meld deze transfer</span></a>\r\n\r\n\r\n\r\n        </td>\r\n    </tr>\r\n</table>\r\n\r\n                                </td>\r\n                            </tr>\r\n                    </table>\r\n                    <!-- ENDOF inner wrapper table -->\r\n                </center>\r\n            </td>\r\n        </tr>\r\n    </table>\r\n    <!-- ENDOF outer wrapper table -->\r\n\r\n    <!-- tracking start -->\r\n    <!-- tracking end -->\r\n\r\n    <!-- Removed GA tracking so we won't hit their limit -->\r\n\r\n    <!-- crazy gmail app fix -->\r\n    <div style=3D"display: none; white-space: nowrap; color: #FFFFFF; font:=\r\n 20px courier;">=C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=\r\n=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =\r\n=C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=A0 =C2=\r\n=A0</div>\r\n    <!-- ENDOF crazy gmail app fix -->\r\n\r\n<img src=3D"http://email.wetransfer.com/wf/open?upn=3Du001.-2FXZ9wMKpvdHn9M=\r\nhPB6zy3j10YvWaGVj68QiIU-2BSU0xiC8VtBTT-2BwCPZGkgJ6DxCoX4lqp3mE7S-2FO1Lfy5As=\r\nXs1LRV664aZJXTXAPdaPCVFqdbbjfeoMI2iUTWTEpf5J1VjQFZwG-2FmKflHEJPxddSzpq491AL=\r\nh1u4uRE9g4qoHn1gsqhIcQu7ImpnZE60lT0ziPuxgnhGRejesgk9HKr2kwtGmKuUxxoJW-2B0v-=\r\n2FzfC47H1ZoghnysNbzcs-2BIvO8hnsl7S7-2BbsSncBIU0jd8qsbvY7ewhUK-2BUaN-2Feq0uT=\r\nsMkk-2F2eG9OFQchRRowjO5B1yRmq5fhM0STawpWmuA8Ljal3BeKuwS-2F4a6LNXTyiuZ5LxvDc=\r\nTAlfEEkLdwwN3X7A5F0VLkUdpKG36AzF7f1gBnAXQ-3D-3D" alt=3D"" width=3D"1" heigh=\r\nt=3D"1" border=3D"0" style=3D"height:1px !important;width:1px !important;bo=\r\nrder-width:0 !important;margin-top:0 !important;margin-bottom:0 !important;=\r\nmargin-right:0 !important;margin-left:0 !important;padding-top:0 !important=\r\n;padding-bottom:0 !important;padding-right:0 !important;padding-left:0 !imp=\r\nortant;"/></body>\r\n</html>\r\n\r\n----==_mimepart_6644b652229b6_b93e4930359--\r\n	multipart/alternative; boundary="--==_mimepart_6644b652229b6_b93e4930359"; charset=UTF-8	2024-05-15 13:19:14	11	info	t	2025-04-20 00:41:19.079013	2025-04-20 00:03:01.028305	2025-04-20 00:41:19.079515
994f6414-fb55-44b4-8bcf-41bd93d90b8d	957dd920-07b6-11f0-9cea-91f375301d21@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	Nieuw contactformulier	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuw contactformulier bericht</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuw bericht binnengekomen via het contactformulier.</p>\n                \n                <div>\n                    <strong>Contactgegevens:</strong><br>\n                    Naam: Bas heijenk <br>\n                    Email: basheijenk96@gmail.com<br>\n                    <br>\n                    <strong>Bericht:</strong><br>\n                    Hallo ik heb me op gegeven maar ik kan helaas niet sorry \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-23 07:15:24	106	info	t	2025-04-20 00:41:19.687087	2025-04-20 00:03:01.031336	2025-04-20 00:41:19.687514
1a1a038e-352e-43b9-9545-8d0831d09f97	bb8c00cc-0369-11f0-bbd9-91f375301d21@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Salih<br>\n                    Email: topraks@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-17 19:55:10	38	inschrijving	f	\N	2025-04-20 00:03:01.042461	2025-04-20 00:03:01.042461
9ee6c78c-1cd0-4f41-b894-3b3c57a99073	<977803562.1683661.1742262446100@s.appsuite.hostnet.nl>	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Fwd: Nieuwe aanmelding ontvangen	------=_Part_1683659_1029948377.1742262446100\r\nContent-Type: multipart/alternative; \r\n\tboundary="----=_Part_1683660_580753535.1742262446100"\r\n\r\n------=_Part_1683660_580753535.1742262446100\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 7bit\r\n\r\n\r\n------=_Part_1683660_580753535.1742262446100\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: 7bit\r\n\r\n<!doctype html>\r\n<html>\r\n <head> \r\n  <meta charset="UTF-8"> \r\n </head>\r\n <body>\r\n  &nbsp;\r\n </body>\r\n</html>\r\n------=_Part_1683660_580753535.1742262446100--\r\n\r\n------=_Part_1683659_1029948377.1742262446100\r\nContent-Type: message/rfc822; name=Nieuwe_aanmelding_ontvangen.eml\r\nContent-Disposition: attachment; filename=Nieuwe_aanmelding_ontvangen.eml\r\n\r\nReturn-Path: <info@dekoninklijkeloop.nl>\r\nDelivered-To: info@dekoninklijkeloop.nl\r\nReceived: from mx4.cst.mailpod13-cph3.one.com ([10.27.16.14])\r\n\tby mailstorage27.cst.mailpod13-cph3.one.com with LMTP\r\n\tid QFEhB0oe0Gd4JAAAu0OFfw\r\n\t(envelope-from <info@dekoninklijkeloop.nl>)\r\n\tfor <info@dekoninklijkeloop.nl>; Tue, 11 Mar 2025 11:28:10 +0000\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\r\n\td=custmx.one.com; s=20201015;\r\n\th=content-transfer-encoding:content-type:subject:to:from:date:mime-version:\r\n\t message-id:x-halone-refid:x-halone-sa:from:x-halone-sa:x-halone-refid;\r\n\tbh=x49PwcnwHVCY1Ht+4iXo0xe13YdV7b7C8kOQz3HMSKg=;\r\n\tb=BAMC7Q7b5dH/FHT3xSBFJ0QTlIAFaO+4Yf5x8LD9CVQ6e7bdW9MEtMaUNdijc5DwceCg1jn/vrE94\r\n\t ErILMxOYm6N7ia980OBOXR4rxFaqnItkikrvtB7BfwQR9KzmwKDHxzgNfs1a3lUtDSh1Xp5FjuwsDW\r\n\t zF2SGOvpd5PzJerWqU+HCE4wEc42L0iuo8AjFj90E+7eLg8MZS3tU9BJTGpHgeZPcFw2T1+aFutBqh\r\n\t PSE/j2s69BWgEXvxlKOJtDHlRtvHu1bOYeztte5djlJ+ZyRlXQkdRdMmCu4Mb6mW1nDgfHeX19qjWX\r\n\t s9stP8mPazu7cmAw5tWci7mHj9ENLdA==\r\nX-HalOne-SA: 1.1\r\nX-HalOne-RefID: 155866::1741692489-0E8994CB-2C4EA744/0/0\r\nX-HalOne-Spam-Probability: 0.008\r\nAuthentication-Results: mx4.pub.mailpod13-cph3.one.com;\r\n\tspf=pass smtp.mailfrom=dekoninklijkeloop.nl smtp.remote-ip=46.30.211.185;\r\n\tdkim=pass header.d=dekoninklijkeloop.nl header.s=ed1 header.a=ed25519-sha256 header.b=wf1RSWMQ;\r\n\tdmarc=pass header.from=dekoninklijkeloop.nl;\r\nReceived: from mailrelay4-2.pub.mailoutpod2-cph3.one.com (mailrelay4-2.pub.mailoutpod2-cph3.one.com [46.30.211.185])\r\n\tby mx4.pub.mailpod13-cph3.one.com (Halon) with ESMTPS\r\n\tid e91c1c06-fe6b-11ef-8050-0892049f56cc;\r\n\tTue, 11 Mar 2025 11:28:09 +0000 (UTC)\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; t=1741692489; x=1742297289;\r\n\td=dekoninklijkeloop.nl; s=rsa1;\r\n\th=content-transfer-encoding:content-type:subject:to:from:date:mime-version:\r\n\t message-id:from;\r\n\tbh=x49PwcnwHVCY1Ht+4iXo0xe13YdV7b7C8kOQz3HMSKg=;\r\n\tb=5/23sckFR0Tgy6yVuL7nlzxwQ7K2H03jZLKIfSs76UXccnQB66IdZUB7+6vj4LIhTJYNYydgQxJxl\r\n\t McFA/aVKVhFr8iE08A3/l1HkejQyafEf/6zTow6eK6trVt2iHI0MIEKvCdj6DIXbnI2IOPf963Pw8N\r\n\t LIZqCeBL8blQAa4jtfCQymTM7eHkvzNUmi/hS3spD9znuLAJe8qyWxMywsn1lJLPoJR41jNpRLBq8X\r\n\t fb8pBNHgczlh7o4zw3WuOvb25R7qw8y766r/nRhlsMuCQFRwYv5tU9bK3KXvaLAul9RVsfNnO1W24y\r\n\t a3CmsUduAwOtJOT/LgeblHSmPbUwkoA==\r\nDKIM-Signature: v=1; a=ed25519-sha256; c=relaxed/relaxed; t=1741692489; x=1742297289;\r\n\td=dekoninklijkeloop.nl; s=ed1;\r\n\th=content-transfer-encoding:content-type:subject:to:from:date:mime-version:\r\n\t message-id:from;\r\n\tbh=x49PwcnwHVCY1Ht+4iXo0xe13YdV7b7C8kOQz3HMSKg=;\r\n\tb=wf1RSWMQGiwQ8/hazfUVoy/zuyX5sqo0rP+3Tm0u03tTB0vkDB/EJswsjvDmRhhhQuXjva7B+q0cV\r\n\t nEuowSCCw==\r\nMessage-ID: e812d0a9-fe6b-11ef-b624-e77cec7da75b@dekoninklijkeloop.nl\r\nX-HalOne-ID: e91c1c06-fe6b-11ef-8050-0892049f56cc\r\nReceived: from mailproxy2.cst.dirpod3-cph3.one.com (n2xm8yy.mp.shared.prod.hostnet.nl [185.107.112.245])\r\n\tby mailrelay4.pub.mailoutpod2-cph3.one.com (Halon) with ESMTPSA\r\n\tid e812d0a9-fe6b-11ef-b624-e77cec7da75b;\r\n\tTue, 11 Mar 2025 11:28:09 +0000 (UTC)\r\nMime-Version: 1.0\r\nDate: Tue, 11 Mar 2025 11:28:09 +0000\r\nFrom: info@dekoninklijkeloop.nl\r\nTo: info@dekoninklijkeloop.nl\r\nSubject: Nieuwe aanmelding ontvangen\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<!DOCTYPE html>\r\n<html>\r\n<head>\r\n    <meta charset=3D"UTF-8">\r\n    <title>Nieuwe Aanmelding - De Koninklijke Loop</title>\r\n</head>\r\n<body>\r\n    <h2>Nieuwe Aanmelding Ontvangen</h2>\r\n    <p>Er is een nieuwe aanmelding binnengekomen via het aanmeldformulier:<=\r\n/p>\r\n\r\n    <h3>Deelnemer Gegevens:</h3>\r\n    <ul>\r\n        <li><strong>Naam:</strong> Manuela van zwam</li>\r\n        <li><strong>E-mail:</strong> benjaminlaan.64a@sheerenloo.nl</li>\r\n        <li><strong>Telefoonnummer:</strong> </li>\r\n        <li><strong>Rol:</strong> Deelnemer</li>\r\n        <li><strong>Gekozen Afstand:</strong> 2.5 KM</li>\r\n        <li><strong>Ondersteuning nodig:</strong> Nee</li>\r\n       =20\r\n    </ul>\r\n\r\n    <p>De deelnemer heeft de algemene voorwaarden geaccepteerd: Ja</p>\r\n</body>\r\n</html>=20\r\n\r\n------=_Part_1683659_1029948377.1742262446100\r\nContent-Type: message/rfc822; name=Nieuwe_aanmelding_ontvangen.eml\r\nContent-Disposition: attachment; filename=Nieuwe_aanmelding_ontvangen.eml\r\n\r\nReturn-Path: <info@dekoninklijkeloop.nl>\r\nDelivered-To: info@dekoninklijkeloop.nl\r\nReceived: from mx4.cst.mailpod13-cph3.one.com ([10.27.16.14])\r\n\tby mailstorage27.cst.mailpod13-cph3.one.com with LMTP\r\n\tid YBxWKD3HzWeGfgAAu0OFfw\r\n\t(envelope-from <info@dekoninklijkeloop.nl>)\r\n\tfor <info@dekoninklijkeloop.nl>; Sun, 09 Mar 2025 16:52:13 +0000\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;\r\n\td=custmx.one.com; s=20201015;\r\n\th=content-transfer-encoding:content-type:to:from:subject:date:mime-version:\r\n\t message-id:x-halone-refid:x-halone-sa:from:x-halone-sa:x-halone-refid;\r\n\tbh=He8zFugeCEFsZ7Z5zq5WbFbUKJnNQK+4W2cnAKmSlpI=;\r\n\tb=WTik6/5kpH80U8qwLEyVXtxr0BToOkcWmKZePX/A2kIcL6xsCKT+Uitw9nvgyJZ9FFsstuP8o5k6Q\r\n\t S80YLQkfEgFZU+Ci0d1m99gGvLn20G0YHLE8gkOoAnL0cBgFepKLPv46JL30dH4vv7bqPNm8/sOA1z\r\n\t wlDIxyIJMCLrAWNK5khWVLumxOMGfauLRjH+s7ay2br26MSU7AGUHGCjMWWZDCL4WU/IWXdc3gOY0a\r\n\t /AIAu9cfV5aFe9C+kdjBMGjuX/49nF1jRMQ/Pmu3nWsv85iNrclfOIpjzD+796U8DnJCe5JJQznWva\r\n\t 53Vu67LKo3vkfK2FoAzXncXlLaiESlA==\r\nX-HalOne-SA: 1.1\r\nX-HalOne-RefID: 155866::1741539133-FD7A74CB-0C903F88/0/0\r\nX-HalOne-Spam-Probability: 0.008\r\nAuthentication-Results: mx4.pub.mailpod13-cph3.one.com;\r\n\tspf=pass smtp.mailfrom=dekoninklijkeloop.nl smtp.remote-ip=46.30.211.247;\r\n\tdkim=pass header.d=dekoninklijkeloop.nl header.s=ed1 header.a=ed25519-sha256 header.b=G2VLHIIB;\r\n\tdmarc=pass header.from=dekoninklijkeloop.nl;\r\nReceived: from mailrelay2-2.pub.mailoutpod3-cph3.one.com (mailrelay2-2.pub.mailoutpod3-cph3.one.com [46.30.211.247])\r\n\tby mx4.pub.mailpod13-cph3.one.com (Halon) with ESMTPS\r\n\tid d990b94e-fd06-11ef-804c-0892049f56cc;\r\n\tSun, 09 Mar 2025 16:52:13 +0000 (UTC)\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed; t=1741539133; x=1742143933;\r\n\td=dekoninklijkeloop.nl; s=rsa1;\r\n\th=content-transfer-encoding:content-type:to:from:subject:date:mime-version:\r\n\t message-id:from;\r\n\tbh=He8zFugeCEFsZ7Z5zq5WbFbUKJnNQK+4W2cnAKmSlpI=;\r\n\tb=eSaqBbBl4683eon6Rk3YpYXz0+E2hY3yyXnNqTcR4sQZevbA8Rb3ioi5ZLoopCeWuB5d2y+GidgdZ\r\n\t 4aboR7pZh24w8a74AfM0xriLkEG2BA808lQINmizPHrJIBKDXX5lhNstlhgCXslYFK4MR0CuhPvkYM\r\n\t wfapgPonjA4l3r3smmEtOOxsjcT9YyhyVJn5g5uSkp7vMu8s+Cl8c+/axgy73Nrhag6eHFTTSQAieM\r\n\t kXnsOUTsSZKdjf8Rq2tt64+aU+rbKu1UH88+PpRSSvgBJ8GzPPEN0AyBtJd5iOBE5jscrl3T/U1Q2+\r\n\t YFI66c75k0+cC+XEfnLKAfGC41FZ44A==\r\nDKIM-Signature: v=1; a=ed25519-sha256; c=relaxed/relaxed; t=1741539133; x=1742143933;\r\n\td=dekoninklijkeloop.nl; s=ed1;\r\n\th=content-transfer-encoding:content-type:to:from:subject:date:mime-version:\r\n\t message-id:from;\r\n\tbh=He8zFugeCEFsZ7Z5zq5WbFbUKJnNQK+4W2cnAKmSlpI=;\r\n\tb=G2VLHIIBWphCT7qen1R2qTa84rZpM4NlQXym/iYvJ5KkHEAJOYrtYo5ynWkZiX7JLycPZIUUeK8R4\r\n\t jULt3VjBA==\r\nMessage-ID: d8e22b52-fd06-11ef-8ac9-b37c246f863f@dekoninklijkeloop.nl\r\nX-HalOne-ID: d990b94e-fd06-11ef-804c-0892049f56cc\r\nReceived: from mailproxy2.cst.dirpod4-cph3.one.com (n2xm8yy.mp.shared.prod.hostnet.nl [185.107.112.245])\r\n\tby mailrelay2.pub.mailoutpod3-cph3.one.com (Halon) with ESMTPSA\r\n\tid d8e22b52-fd06-11ef-8ac9-b37c246f863f;\r\n\tSun, 09 Mar 2025 16:52:12 +0000 (UTC)\r\nMime-Version: 1.0\r\nDate: Sun, 09 Mar 2025 16:52:13 +0000\r\nSubject: Nieuwe aanmelding ontvangen\r\nFrom: info@dekoninklijkeloop.nl\r\nTo: info@dekoninklijkeloop.nl\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<!DOCTYPE html>\r\n<html>\r\n<head>\r\n    <meta charset=3D"UTF-8">\r\n    <title>Nieuwe Aanmelding - De Koninklijke Loop</title>\r\n</head>\r\n<body>\r\n    <h2>Nieuwe Aanmelding Ontvangen</h2>\r\n    <p>Er is een nieuwe aanmelding binnengekomen via het aanmeldformulier:<=\r\n/p>\r\n\r\n    <h3>Deelnemer Gegevens:</h3>\r\n    <ul>\r\n        <li><strong>Naam:</strong> Joyce Thielen</li>\r\n        <li><strong>E-mail:</strong> Joyce.thielen@sheerenloo.nl</li>\r\n        <li><strong>Telefoonnummer:</strong> </li>\r\n        <li><strong>Rol:</strong> Begeleider</li>\r\n        <li><strong>Gekozen Afstand:</strong> 6 KM</li>\r\n        <li><strong>Ondersteuning nodig:</strong> Nee</li>\r\n       =20\r\n    </ul>\r\n\r\n    <p>De deelnemer heeft de algemene voorwaarden geaccepteerd: Ja</p>\r\n</body>\r\n</html>=20\r\n\r\n------=_Part_1683659_1029948377.1742262446100--\r\n	multipart/mixed; boundary="----=_Part_1683659_1029948377.1742262446100"	2025-03-18 02:47:26	40	inschrijving	f	\N	2025-04-20 00:03:01.044592	2025-04-20 00:03:01.044592
b0fac1c6-772e-4162-8b5b-0bdeb00f7331	e65ddf03-05a5-11f0-a00f-471e5e6299d9@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: TEST*<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Anders<br>\n                    \n                    <br>\n                    Bijzonderheden: test<br>\n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-20 16:10:55	43	inschrijving	f	\N	2025-04-20 00:03:01.047198	2025-04-20 00:03:01.047198
209358bb-6580-4d5c-b3fa-a978f1676905	ba1cc684-073c-11f0-93c3-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Bas heijenk <br>\n                    Email: basheijenk96@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-22 16:43:06	44	inschrijving	f	\N	2025-04-20 00:03:01.049118	2025-04-20 00:03:01.049118
66c851a7-2bb5-4db7-9741-147ab28928eb	5aa31c0c-07f0-11f0-9970-471e5e6299d9@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Manuela van Zwam<br>\n                    Email: rik.van-harxen@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Ja<br>\n                    \n                    <br>\n                    Bijzonderheden: Vaste begeleider die meeloopt<br>\n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-23 14:08:56	45	inschrijving	f	\N	2025-04-20 00:03:01.051074	2025-04-20 00:03:01.051074
3b2fc44d-6a7a-4932-8bb0-e34008692aa6	a1317501-07ff-11f0-93d3-15d104443858@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: jeffrey<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Anders<br>\n                    \n                    <br>\n                    Bijzonderheden: Telegram test<br>\n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-23 15:58:17	46	inschrijving	f	\N	2025-04-20 00:03:01.053081	2025-04-20 00:03:01.053081
6319f977-9150-44e3-9759-a6608f819c69	e6e74642-0800-11f0-9359-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Telegram<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-23 16:07:22	47	inschrijving	f	\N	2025-04-20 00:03:01.055327	2025-04-20 00:03:01.055327
2a794d8e-fbb0-4aeb-a3f4-a0930154e66b	2866f9ee-0809-11f0-94a3-15d104443858@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: TGTest<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    <br>\n                    Telefoon: 06123456789<br>\n                    \n                    Rol: Begeleider<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Anders<br>\n                    \n                    <br>\n                    Bijzonderheden: Telegram Test bericht - officiele weg<br>\n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-23 17:06:29	48	inschrijving	f	\N	2025-04-20 00:03:01.057468	2025-04-20 00:03:01.057468
1421bfbd-c040-4f86-99e8-b182ba51d495	b441f0a7-0890-11f0-9729-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Annerieke Mandemaker-Timmer<br>\n                    Email: annerieketimmer@hotmail.com<br>\n                    <br>\n                    Telefoon: 06 17 37 28 40 <br>\n                    \n                    Rol: Begeleider<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-24 09:16:45	49	inschrijving	f	\N	2025-04-20 00:03:01.059488	2025-04-20 00:03:01.059488
2cbda9d0-9e47-4178-9c7a-b4a53a3346c7	e2bda172-0890-11f0-bce6-2b8368a4d5c5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Jean-paul Hup<br>\n                    Email: molenkamp.19@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-24 09:18:03	50	inschrijving	f	\N	2025-04-20 00:03:01.061571	2025-04-20 00:03:01.061571
faf8a5ec-2e4b-4824-b890-748e3d00abfe	166c7368-08d5-11f0-8427-2b8368a4d5c5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Arno Kerkvliet<br>\n                    Email: arno.kerkvliet@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-24 17:26:15	51	inschrijving	f	\N	2025-04-20 00:03:01.063625	2025-04-20 00:03:01.063625
60fc3dc7-0b81-4337-b263-cbc4c9571586	4105f9b1-08d5-11f0-9e39-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Mirjam Kerkvliet<br>\n                    Email: mirjam.kerkvliet@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-24 17:27:27	52	inschrijving	f	\N	2025-04-20 00:03:01.065665	2025-04-20 00:03:01.065665
8f74d1a3-a7a3-457f-80df-d791e9d91a51	622c2e0a-08e0-11f0-a599-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Karin de Jong<br>\n                    Email: karin.de.jong82@outlook.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-24 18:47:08	53	inschrijving	f	\N	2025-04-20 00:03:01.067894	2025-04-20 00:03:01.067894
e63afea6-9e75-40af-b6a1-bbbebaeb86ee	ab3ba39b-0a3d-11f0-8bc3-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: A. Bistolfi<br>\n                    Email: nedarg@icloud.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-26 12:27:24	54	inschrijving	f	\N	2025-04-20 00:03:01.069988	2025-04-20 00:03:01.069988
454f3d76-d2dd-4889-9bea-0af27d67a0ef	f7db975e-0a5e-11f0-8e31-e77cec7da75b@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Ayla Toprak<br>\n                    Email: gamergirlayla@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 10 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-26 16:25:46	55	inschrijving	f	\N	2025-04-20 00:03:01.072073	2025-04-20 00:03:01.072073
fd8a91c8-13aa-4bb6-a559-a4b6a31115bf	2105bae4-0a5f-11f0-9088-91f375301d21@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Mila Veenendaal<br>\n                    Email: gaminggirlayla@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 10 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-26 16:26:56	56	inschrijving	f	\N	2025-04-20 00:03:01.109077	2025-04-20 00:03:01.109077
af2c4a73-8e92-4dfa-bd40-2341dfbe7520	5276dd15-0c15-11f0-bc67-91f375301d21@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: jeffrey<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 10 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-28 20:43:37	57	inschrijving	f	\N	2025-04-20 00:03:01.111695	2025-04-20 00:03:01.111695
e957afcf-b51d-4a7d-8676-9dfb8718a396	498ecf8e-0c85-11f0-a987-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Klaske van de glind<br>\n                    Email: Klaskehiddes@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-29 10:05:07	58	inschrijving	f	\N	2025-04-20 00:03:01.113679	2025-04-20 00:03:01.113679
f2fdc142-3d6d-4573-a775-eb4f78ce9466	96558240-0c85-11f0-92da-4d2191f5f3b5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Bertram tijsma<br>\n                    Email: Klaskehiddes@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-29 10:07:14	59	inschrijving	f	\N	2025-04-20 00:03:01.115666	2025-04-20 00:03:01.115666
d488dac5-0e68-4cc0-b849-8d38a5c22382	10f1fbd1-0d3e-11f0-81b9-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Han van Doornik<br>\n                    Email: LaanvanGS.26@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Ja<br>\n                    \n                    <br>\n                    Bijzonderheden: Ik wil wel graag begeleiding <br>\n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-30 08:07:48	60	inschrijving	f	\N	2025-04-20 00:03:01.117646	2025-04-20 00:03:01.117646
3ec70d0c-9683-4cad-a6c2-09e03424579f	34de60b0-0e27-11f0-926e-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Anneke van de Glind <br>\n                    Email: Klaskehiddes@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-31 11:56:40	61	inschrijving	f	\N	2025-04-20 00:03:01.119812	2025-04-20 00:03:01.119812
cd394b47-aa2b-4282-a0bd-39a3658c39f9	4fb452db-0e27-11f0-b68b-91f375301d21@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Noa hiddes<br>\n                    Email: Klaskehiddes@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-31 11:57:26	62	inschrijving	f	\N	2025-04-20 00:03:01.121919	2025-04-20 00:03:01.121919
d25cebb7-4353-495c-bdbc-06b55bf273a1	b3e281cd-1317-11f0-8f6a-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Sylvia Dijkstra<br>\n                    Email: sylvia.dijkstra@sheerenloo.nl<br>\n                    <br>\n                    Telefoon: 0683081728 <br>\n                    \n                    Rol: Begeleider<br>\n                    Afstand: 2.5 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-06 18:48:18	63	inschrijving	f	\N	2025-04-20 00:03:01.124115	2025-04-20 00:03:01.124115
a5d4c694-7bac-4148-baee-92e5a550fcc4	bc3b4654-176c-11f0-8963-15d104443858@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Diesmer <br>\n                    Email: diesbosje@hotmail.com<br>\n                    <br>\n                    Telefoon: 0613429612<br>\n                    \n                    Rol: Begeleider<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-12 07:07:05	64	inschrijving	f	\N	2025-04-20 00:03:01.127411	2025-04-20 00:03:01.127411
4de9da89-1f7f-4a13-a9c6-9a1d5fe858ff	f43e811b-176c-11f0-91f9-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Albert <br>\n                    Email: diesbosje@hotmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-12 07:08:39	65	inschrijving	f	\N	2025-04-20 00:03:01.129804	2025-04-20 00:03:01.129804
9b9d77be-e5bd-4d76-97f4-5075b93357cd	0e9da5b4-176d-11f0-91fa-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Theun <br>\n                    Email: diesbosje@hotmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-12 07:09:23	66	inschrijving	f	\N	2025-04-20 00:03:01.140875	2025-04-20 00:03:01.140875
2f1bf558-97b1-40c8-b190-9c2aa40dd1b4	3ec275c2-252c-11f0-a36e-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Sulleeman <br>\n                    Email: lidaahmadi99@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 19:00:41	79	inschrijving	t	2025-04-29 19:02:00.934741	2025-04-29 19:01:55.79811	2025-04-29 19:02:00.935239
3dee63e6-f8af-41eb-ac22-b82424dabbc0	ebafa007-1787-11f0-a863-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Henk Rekers <br>\n                    Email: h.rekers1959@kpnmail.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-12 10:21:40	68	inschrijving	t	2025-04-20 12:51:03.439133	2025-04-20 00:03:01.158165	2025-04-20 12:51:03.439737
b674ce7c-496d-4570-b398-8df78621f1fd	c0c78ba6-1787-11f0-a253-e77cec7da75b@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Hilde Rekers <br>\n                    Email: h.rekers1959@kpnmail.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-12 10:20:29	67	inschrijving	t	2025-04-23 15:47:34.074482	2025-04-20 00:03:01.153194	2025-04-23 15:47:34.078198
be86cbaf-a385-4ef3-9da9-21d53f575b98	8696387b-1b6e-11f0-9230-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Bram Aarnoudse<br>\n                    Email: bram.aarnoudse@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-17 09:29:57	71	inschrijving	t	2025-04-20 00:28:09.95285	2025-04-20 00:03:01.165264	2025-04-20 00:28:09.953224
c03871f1-c892-41ac-9e2b-d93b15a6bf11	557058a8-0d46-11f0-a454-2b8368a4d5c5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	Nieuw contactformulier	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuw contactformulier bericht</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuw bericht binnengekomen via het contactformulier.</p>\n                \n                <div>\n                    <strong>Contactgegevens:</strong><br>\n                    Naam: Mirjam Rijswijk<br>\n                    Email: mirjamrijswijk@icloud.com<br>\n                    <br>\n                    <strong>Bericht:</strong><br>\n                    Ik wil graag helpen met dit evenement. Ik ben evenementen eerste hulpverlener bij het Rode Kruis en daarnaast wandelbegeleider bij de wandeluitdaging Apeldoorn. Ik hoor graag van jullie . Hartelijke groet, Mirjam Rijswijk \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-03-30 09:06:59	115	info	t	2025-04-20 00:41:21.945158	2025-04-20 00:03:01.033582	2025-04-20 00:41:21.945559
a6b298d9-4a94-45ca-84f3-a9fbacaf3bba	6e71645e-252c-11f0-b0b8-4d2191f5f3b5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Yunis <br>\n                    Email: lidaahmadi99@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 19:02:02	80	inschrijving	t	2025-10-05 18:24:36.982171	2025-04-29 19:02:13.752553	2025-10-05 18:24:36.982479
1ee91c4e-a045-4543-bcb8-a3287ca543ae	e7594aa7-29a4-11f0-b1f5-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	Nieuw contactformulier	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuw contactformulier bericht</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuw bericht binnengekomen via het contactformulier.</p>\n                \n                <div>\n                    <strong>Contactgegevens:</strong><br>\n                    Naam: Theodora Naus<br>\n                    Email: tgemooi@gmail.com<br>\n                    <br>\n                    <strong>Bericht:</strong><br>\n                    Ik had me aangemeld voor de afstand van 6 km maar kan helaas niet meedoen. Gelukkig had ik al wel gedoneerd.\nSucces verder!\n\nVriendelijke groet,\n\nTheodora Naus\n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-05 11:34:30	126	info	t	2025-05-06 08:51:42.524992	2025-05-05 11:47:14.81734	2025-05-06 08:51:42.528623
5e59095f-3a84-477c-90e8-ee9a0fdccdfd	<3b7190100a08e4fc4cd41026c491e720@mta1.fredericshensch.de>	Hostnet bv <web@mta1.fredericshensch.de>	info@dekoninklijkeloop.nl	Een belangrijke opmerking over het verlengen van uw domein	This is a multi-part message in MIME format.\r\n\r\n--12e7e7af0d370b4a7e6f0ac9d77a9a797\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n\r\n\r\nGeachte heer/mevrouw eigenaar van het landgoed,\r\nuw domein bereikt de vervaldatum en moet die dag worden verwijderd.Vervaldatum: 05/26/2025 \r\nOm onderbreking van e-mail- of websiteservices op de vervaldatum en verlies van het door u gekozen domein te voorkomen,\r\nNeem dan contact op met de hostnet Customer Service en regel de verlenging van je domein, of verleng het rechtstreeks via volgende link \r\nToegang krijgen\r\nAls je je domein zelf hebt laten verwijderen en je hebt het niet meer nodig, hoef je niets te doen.\r\n------------------------------------------------------------------------jjVorsi­tzender d­es Aufsic­htsratesj©jhostnetjAGDit is een automatische systeem mail. Gelieve deze e-mail niet te beantwoorden.\r\n\r\n\r\n--12e7e7af0d370b4a7e6f0ac9d77a9a797\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n<html><head></head><body><font color="#ffffff">\r\n<p><span style="FONT-SIZE: small"><span style="FONT-FAMILY: Arial,Helvetica,sans-serif"><span style="COLOR: #555555"><img alt="" src="https://images.squarespace-cdn.com/content/v1/598c1a6b46c3c4a7b18a8908/c29854d4-5d30-4ac9-8874-42ccbe21dad9/hostnet.png" width="129" height="73" data-cke-saved-src="https://images.squarespace-cdn.com/content/v1/598c1a6b46c3c4a7b18a8908/c29854d4-5d30-4ac9-8874-42ccbe21dad9/hostnet.png"></span></span></span></p></font>\r\n<p>Geachte heer/mevrouw eigenaar van het landgoed,</p>\r\n<p>uw domein bereikt de vervaldatum en moet die dag worden verwijderd.<br>Vervaldatum: 05/26/2025 </p>\r\n<p>Om onderbreking van e-mail- of websiteservices op de vervaldatum en verlies van het door u gekozen domein te voorkomen,</p>\r\n<p>Neem dan contact op met de hostnet Customer Service en regel de verlenging van je domein, of verleng het rechtstreeks via volgende link </p>\r\n<p><span style="COLOR: #ffffff"><a href="https://appespacemanager.piscinadelgolf.it?pwd=noc" data-cke-saved-href="http://google.com"><span style="FONT-SIZE: 16px"><strong>Toegang krijgen</strong></span></a><wbr></span></p>\r\n<p>Als je je domein zelf hebt laten verwijderen en je hebt het niet meer nodig, hoef je niets te doen.</p>\r\n<p>------------------------------<wbr>------------------------------<wbr>------------<br><span style="COLOR: #ffffff">jj</span>Vorsi­tzender d­es Aufsic­htsrates<span style="COLOR: #ffffff">j</span><span style="COLOR: #ffffff"><span style="FONT-SIZE: small"><span style="FONT-FAMILY: Arial,Helvetica,sans-serif"><span style="COLOR: #666666"><br><br>©<font color="#ffffff" face="Times New Roman">j</font>hostnet<font color="#ffffff" face="Times New Roman">j</font>AG<br><br><font size="2">Dit is een automatische systeem mail. Gelieve deze e-mail niet te beantwoorden.</font></span></span></span></span></p></body></html>\r\n\r\n\r\n\r\n--12e7e7af0d370b4a7e6f0ac9d77a9a797--\r\n\r\n	multipart/alternative; boundary="12e7e7af0d370b4a7e6f0ac9d77a9a797"	2025-05-21 10:17:21	129	info	t	2025-09-10 18:39:53.648414	2025-05-21 10:32:14.470405	2025-09-10 18:39:53.649316
af81d62f-a08b-412d-be06-76b929d1142f	<494d79d784bc8078fbee6b52e2db89cf9b308b77@export-afrika.com>	Arena Autoinkoop - Rob <contact@export-afrika.com>	info@dekoninklijkeloop.nl	Snel en betrouwbaar uw auto verkopen	\r\n--_=_swift_1748628155_c1bb87ceae3b07d332bcc2d85d9f1818_=_\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nGeachte mevrouw / mijnheer,\r\n\r\nBent u van plan om uw bedrijfsauto, busjes=\r\n of andere voertuigen te\r\nverkopen?\r\nAutoinkoop Arena kan u daarbij helpe=\r\nn. Wij zijn gespecialiseerd in de\r\nin- en verkoop van diverse typen auto=\r\n=E2=80=99s, busjes en bedrijfswagens,\r\nen wij kopen graag uw voertuigen op=\r\n voor een eerlijke en passende\r\nprijs.\r\n\r\nWIJ KOPEN VOERTUIGEN AAN:\r\n=\r\n\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Ongeacht de staat\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0 Ongeacht de kilometerstand\r\n*=C2=A0=C2=\r\n=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0 Ongeacht he=\r\nt type brandstof (diesel,\r\nbenzine, LPG)\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=\r\n=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0 Ongeacht eventuele schade, de=\r\nfecten of een\r\nverlopen APK\r\n\r\nWij begrijpen dat voertuigen, vooral uit =\r\noudere bouwjaren, mogelijk\r\ngebruikssporen vertonen en veel kilometers heb=\r\nben gereden. Dit is geen\r\nenkel probleem voor ons =E2=80=93 zelfs een geld=\r\nige APK is niet vereist.\r\n\r\nWij overhandigen u na verkoop een officieel=\r\n\r\nRDW-vrijwaringsbewijs.Daarnaast hebben wij een landelijke\r\nophaaldienst=\r\n, zodat u uw voertuig op kunt laten halen wanneer het u\r\nhet beste uitkomt=\r\n =E2=80=93 geheel kosteloos en met een snelle, nette\r\nafhandeling.\r\n\r\nAU=\r\nTOINKOOP ARENA ZORGT ERVOOR DAT U:\r\n\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0 Uw voertuig vlot en zonder gedoe=\r\n kunt\r\nverkopen\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=\r\n=A0=C2=A0=C2=A0=C2=A0 Vrijwaring en transport naar ons bedrijf\r\nzonder kos=\r\nten geregeld worden\r\n*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0=C2=A0=C2=A0=C2=A0 De betaling ontvangt zoals u dat wenst\r\n\r\nU kunt=\r\n uw voertuig geheel kosteloos en vrijblijvend aanmelden via onze\r\nwebsite =\r\nwww.autoinkoop-arena.nl\r\n[https://export-afrika.com/latest/index.php/campa=\r\nigns/bw153mozoce4d/track-url/st150wbz5k503/1ec27dca20c442265160b213625615bf=\r\ne9bba620].\r\nReageren op deze e-mail kan natuurlijk ook.\r\n\r\nBent u moment=\r\neel nog niet van plan om uw voertuig te verkopen, maar\r\nmisschien in de to=\r\nekomst? Bewaar dan dit bericht, wij blijven altijd\r\nge=C3=AFnteresseerd in=\r\n het opkopen van voertuigen.\r\n\r\nWIJ HOPEN SNEL VAN U TE HOREN EN DANKEN U=\r\n ALVAST VOOR UW AANDACHT.\r\n\r\n=C2=A0\r\nMet vriendelijke groet,\r\n\r\nAutoin=\r\nkoop Arena , David\r\n\r\nUnsubscribe\r\n[https://export-afrika.com/latest/ind=\r\nex.php/campaigns/bw153mozoce4d/track-url/st150wbz5k503/d163aac23675c9af4b7d=\r\n88a61e287cd252d2fc24]\r\n=C2=A0\r\n\r\n--_=_swift_1748628155_c1bb87ceae3b07d332bcc2d85d9f1818_=_\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<!DOCTYPE html>\r\n<html>\r\n<head><meta charset=3D"utf-8"/>\r\n=09<title> Snel en betrouwbaar uw auto verkopen </title>\r\n</head>\r\n<body><span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,san=\r\ns-serif">Geachte mevrouw / mijnheer,</span></span><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">Bent u van plan om uw bedrijfsauto, busjes of andere voertuigen te verko=\r\npen?<br />\r\nAutoinkoop Arena kan u daarbij helpen. Wij zijn gespecialiseerd in de in- e=\r\nn verkoop van diverse typen auto=E2=80=99s, busjes en bedrijfswagens, en wi=\r\nj kopen graag uw voertuigen op voor een eerlijke en passende prijs.</span><=\r\n/span><br />\r\n<br />\r\n<strong><span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,s=\r\nans-serif">Wij kopen voertuigen aan:</span></span></strong><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Ongeacht de staat</span></span><br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Ongeacht de kilometerstand</span></span><br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Ongeacht het type brandstof (diesel, benzine, LPG)</span></span><br =\r\n/>\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Ongeacht eventuele schade, defecten of een verlopen APK</span></span=\r\n><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">Wij begrijpen dat voertuigen, vooral uit oudere bouwjaren, mogelijk gebr=\r\nuikssporen vertonen en veel kilometers hebben gereden. Dit is geen enkel pr=\r\nobleem voor ons =E2=80=93 zelfs een geldige APK is niet vereist.</span></sp=\r\nan><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">Wij overhandigen u na verkoop een officieel RDW-vrijwaringsbewijs.Daarna=\r\nast hebben wij een landelijke ophaaldienst, zodat u uw voertuig op kunt lat=\r\nen halen wanneer het u het beste uitkomt =E2=80=93 geheel kosteloos en met =\r\neen snelle, nette afhandeling.</span></span><br />\r\n<br />\r\n<strong><span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,s=\r\nans-serif">Autoinkoop Arena zorgt ervoor dat u:</span></span></strong><br /=\r\n>\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Uw voertuig vlot en zonder gedoe kunt verkopen</span></span><br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 Vrijwaring en transport naar ons bedrijf zonder kosten geregeld word=\r\nen</span></span><br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">*=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=C2=A0=\r\n=C2=A0 De betaling ontvangt zoals u dat wenst</span></span><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf"><strong>U kunt uw voertuig geheel kosteloos en vrijblijvend aanmelden vi=\r\na onze website <a href=3D"https://export-afrika.com/latest/index.php/campai=\r\ngns/bw153mozoce4d/track-url/st150wbz5k503/1ec27dca20c442265160b213625615bfe=\r\n9bba620">www.autoinkoop-arena.nl</a>.<br />\r\nReageren op deze <a href=3D"mailto:info@autoinkoop-arena.nl?subject=3DIk%20=\r\nwil%20graag%20mijn%20voertuig%20vrijblijvend%20aanbieden&body=3DKenteken%3A=\r\n%0AMerk%3A%0AModel%3A%0ABouwjaar%3A%0AKm-stand%3A%0ABijzonderheden%3A%0AGew=\r\nenste%20prijs%3A">e-mail</a> kan natuurlijk ook.</strong><br />\r\n<br />\r\nBent u momenteel nog niet van plan om uw voertuig te verkopen, maar misschi=\r\nen in de toekomst? Bewaar dan dit bericht, wij blijven altijd ge=C3=AFntere=\r\nsseerd in het opkopen van voertuigen.</span></span><br />\r\n<br />\r\n<strong><span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,s=\r\nans-serif">Wij hopen snel van u te horen en danken u alvast voor uw aandach=\r\nt.</span></span></strong><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">=C2=A0</span></span><br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">Met vriendelijke groet,</span></span><br />\r\n<br />\r\n<span style=3D"font-size:11pt"><span style=3D"font-family:Calibri,sans-seri=\r\nf">Autoinkoop Arena , David</span></span><br />\r\n<br />\r\n<br />\r\n<span style=3D"font-size:12px;"><a data-unsubtag=3D"_DIRECT_UNSUBSCRIBE_URL=\r\n_" href=3D"https://export-afrika.com/latest/index.php/campaigns/bw153mozoce=\r\n4d/track-url/st150wbz5k503/d163aac23675c9af4b7d88a61e287cd252d2fc24">Unsubs=\r\ncribe</a></span><br />\r\n=C2=A0<img width=3D"1" height=3D"1" src=3D"https://export-afrika.com/latest=\r\n/index.php/campaigns/bw153mozoce4d/track-opening/st150wbz5k503" alt=3D"" />=\r\n\r\n</body>\r\n</html>\r\n\r\n--_=_swift_1748628155_c1bb87ceae3b07d332bcc2d85d9f1818_=_--\r\n\r\n	multipart/alternative; boundary="_=_swift_1748628155_c1bb87ceae3b07d332bcc2d85d9f1818_=_"	2025-05-30 18:02:35	130	info	t	2025-09-10 18:40:05.961476	2025-05-30 18:17:14.460808	2025-09-10 18:40:05.963872
8ee872d8-b239-48a4-9a4f-40ca82c13020	92b11aea-33db-11f0-812f-f3c0f7fef5ee@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: jeffrey<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-18 11:31:02	93	inschrijving	t	2025-09-10 18:40:19.938094	2025-05-18 11:32:14.329127	2025-09-10 18:40:19.938558
aedd6601-c533-43a3-b8d2-a8b081800270	2a5ebdb9-253f-11f0-b287-4d2191f5f3b5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Mayke Rood<br>\n                    Email: rood1960@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 21:16:09	82	inschrijving	t	2025-10-05 18:24:38.008221	2025-04-29 21:17:16.307117	2025-10-05 18:24:38.008577
2369ca8d-f68c-4997-b3ac-7f67bcc86d54	<012301dbbf42$0319f0a0$094dd1e0$@gmail.com>	<rood1960@gmail.com>	info@dekoninklijkeloop.nl	AFMELDING loop 17 mei	This is a multipart message in MIME format.\r\n\r\n------=_NextPart_000_0124_01DBBF52.C6A3AB00\r\nContent-Type: text/plain;\r\n\tcharset="UTF-8"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nGoedemiddag,\r\n\r\n=20\r\n\r\nHelaas moet ik mij afmelden voor de wandeltocht op 17 mei. Desondanks =\r\nwens ik jullie een mooie, succesvolle loop toe.\r\n\r\n=20\r\n\r\nMet vriendelijke groet,\r\n\r\nMayke Rood\r\n\r\n=20\r\n\r\nVan: info@dekoninklijkeloop.nl <info@dekoninklijkeloop.nl>=20\r\nVerzonden: dinsdag 29 april 2025 23:16\r\nAan: rood1960@gmail.com\r\nOnderwerp: Bedankt voor je aanmelding\r\n\r\n=20\r\n\r\n  =\r\n<https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e59=\r\n3a1e81556b4238_0760849fb8_yn6vdm.png>=20\r\n\r\n\r\nBedankt voor je aanmelding!\r\n\r\n\r\nBeste Mayke Rood,\r\n\r\nBedankt voor je inschrijving voor De Koninklijke Loop. We hebben je =\r\naanmelding in goede orde ontvangen.\r\n\r\nJe inschrijfgegevens:\r\nNaam: Mayke Rood\r\nEmail: rood1960@gmail.com <mailto:rood1960@gmail.com>=20\r\nRol: Deelnemer\r\nAfstand: 15 KM\r\n\r\nOndersteuning: Nee\r\n\r\nHeb je vragen over je inschrijving? Bekijk dan onze website voor meer =\r\ninformatie:\r\n\r\n*\t <https://dekoninklijkeloop.nl/faq> Veelgestelde vragen\r\n*\t <https://dekoninklijkeloop.nl/over-ons> Over De Koninklijke Loop\r\n\r\n <https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/> =\r\nFacebook  <https://www.instagram.com/koninklijkeloop/> Instagram=20\r\n\r\nMet sportieve groet,\r\nTeam De Koninklijke Loop\r\n\r\n=C2=A9 2025 De Koninklijke Loop. Alle rechten voorbehouden.\r\n\r\n\r\n------=_NextPart_000_0124_01DBBF52.C6A3AB00\r\nContent-Type: text/html;\r\n\tcharset="UTF-8"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html xmlns:v=3D"urn:schemas-microsoft-com:vml" =\r\nxmlns:o=3D"urn:schemas-microsoft-com:office:office" =\r\nxmlns:w=3D"urn:schemas-microsoft-com:office:word" =\r\nxmlns:m=3D"http://schemas.microsoft.com/office/2004/12/omml" =\r\nxmlns=3D"http://www.w3.org/TR/REC-html40"><head><meta =\r\nhttp-equiv=3DContent-Type content=3D"text/html; charset=3Dutf-8"><meta =\r\nname=3DGenerator content=3D"Microsoft Word 15 (filtered medium)"><!--[if =\r\n!mso]><style>v\\:* {behavior:url(#default#VML);}\r\no\\:* {behavior:url(#default#VML);}\r\nw\\:* {behavior:url(#default#VML);}\r\n.shape {behavior:url(#default#VML);}\r\n</style><![endif]--><title>Bedankt voor je inschrijving - De Koninklijke =\r\nLoop</title><style><!--\r\n/* Font Definitions */\r\n@font-face\r\n\t{font-family:Wingdings;\r\n\tpanose-1:5 0 0 0 0 0 0 0 0 0;}\r\n@font-face\r\n\t{font-family:"Cambria Math";\r\n\tpanose-1:2 4 5 3 5 4 6 3 2 4;}\r\n@font-face\r\n\t{font-family:Calibri;\r\n\tpanose-1:2 15 5 2 2 2 4 3 2 4;}\r\n@font-face\r\n\t{font-family:Aptos;}\r\n@font-face\r\n\t{font-family:"Segoe UI";\r\n\tpanose-1:2 11 5 2 4 2 4 2 2 3;}\r\n/* Style Definitions */\r\np.MsoNormal, li.MsoNormal, div.MsoNormal\r\n\t{margin:0cm;\r\n\tfont-size:12.0pt;\r\n\tfont-family:"Aptos",sans-serif;}\r\nh1\r\n\t{mso-style-priority:9;\r\n\tmso-style-link:"Kop 1 Char";\r\n\tmso-margin-top-alt:auto;\r\n\tmargin-right:0cm;\r\n\tmso-margin-bottom-alt:auto;\r\n\tmargin-left:0cm;\r\n\tfont-size:24.0pt;\r\n\tfont-family:"Aptos",sans-serif;\r\n\tfont-weight:bold;}\r\na:link, span.MsoHyperlink\r\n\t{mso-style-priority:99;\r\n\tcolor:blue;\r\n\ttext-decoration:underline;}\r\nspan.Kop1Char\r\n\t{mso-style-name:"Kop 1 Char";\r\n\tmso-style-priority:9;\r\n\tmso-style-link:"Kop 1";\r\n\tfont-family:"Calibri Light",sans-serif;\r\n\tcolor:#2F5496;}\r\nspan.E-mailStijl31\r\n\t{mso-style-type:personal-reply;\r\n\tfont-family:"Calibri",sans-serif;\r\n\tcolor:windowtext;}\r\n.MsoChpDefault\r\n\t{mso-style-type:export-only;\r\n\tfont-size:10.0pt;\r\n\tmso-ligatures:none;}\r\n@page WordSection1\r\n\t{size:612.0pt 792.0pt;\r\n\tmargin:70.85pt 70.85pt 70.85pt 70.85pt;}\r\ndiv.WordSection1\r\n\t{page:WordSection1;}\r\n/* List Definitions */\r\n@list l0\r\n\t{mso-list-id:232274059;\r\n\tmso-list-template-ids:-425793744;}\r\n@list l0:level1\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=B7;\r\n\tmso-level-tab-stop:36.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Symbol;}\r\n@list l0:level2\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:o;\r\n\tmso-level-tab-stop:72.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:"Courier New";\r\n\tmso-bidi-font-family:"Times New Roman";}\r\n@list l0:level3\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:108.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level4\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:144.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level5\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:180.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level6\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:216.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level7\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:252.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level8\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:288.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\n@list l0:level9\r\n\t{mso-level-number-format:bullet;\r\n\tmso-level-text:=EF=82=A7;\r\n\tmso-level-tab-stop:324.0pt;\r\n\tmso-level-number-position:left;\r\n\ttext-indent:-18.0pt;\r\n\tmso-ansi-font-size:10.0pt;\r\n\tfont-family:Wingdings;}\r\nol\r\n\t{margin-bottom:0cm;}\r\nul\r\n\t{margin-bottom:0cm;}\r\n--></style><!--[if gte mso 9]><xml>\r\n<o:shapedefaults v:ext=3D"edit" spidmax=3D"1026" />\r\n</xml><![endif]--><!--[if gte mso 9]><xml>\r\n<o:shapelayout v:ext=3D"edit">\r\n<o:idmap v:ext=3D"edit" data=3D"1" />\r\n</o:shapelayout></xml><![endif]--></head><body bgcolor=3D"#F3F4F6" =\r\nlang=3DNL link=3Dblue vlink=3Dpurple style=3D'word-wrap:break-word'><div =\r\nclass=3DWordSection1><p class=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'>Goedemiddag,<o:p></o:p></span></p><p =\r\nclass=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'><o:p>&nbsp;</o:p></span></p><p class=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'>Helaas moet ik mij afmelden voor de wandeltocht op 17 mei. =\r\nDesondanks wens ik jullie een mooie, succesvolle loop =\r\ntoe.<o:p></o:p></span></p><p class=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'><o:p>&nbsp;</o:p></span></p><p class=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'>Met vriendelijke groet,<o:p></o:p></span></p><p =\r\nclass=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'>Mayke Rood<o:p></o:p></span></p><p class=3DMsoNormal><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif;mso-fareast-la=\r\nnguage:EN-US'><o:p>&nbsp;</o:p></span></p><div><div =\r\nstyle=3D'border:none;border-top:solid #E1E1E1 1.0pt;padding:3.0pt 0cm =\r\n0cm 0cm'><p class=3DMsoNormal><b><span =\r\nstyle=3D'font-size:11.0pt;font-family:"Calibri",sans-serif'>Van:</span></=\r\nb><span style=3D'font-size:11.0pt;font-family:"Calibri",sans-serif'> =\r\ninfo@dekoninklijkeloop.nl &lt;info@dekoninklijkeloop.nl&gt; =\r\n<br><b>Verzonden:</b> dinsdag 29 april 2025 23:16<br><b>Aan:</b> =\r\nrood1960@gmail.com<br><b>Onderwerp:</b> Bedankt voor je =\r\naanmelding<o:p></o:p></span></p></div></div><p =\r\nclass=3DMsoNormal><o:p>&nbsp;</o:p></p><div><div><p class=3DMsoNormal =\r\nalign=3Dcenter style=3D'text-align:center;background:#FF9328'><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif;color:white'><img width=3D554 =\r\nheight=3D214 style=3D'width:5.775in;height:2.2333in' id=3D"_x0000_i1025" =\r\nsrc=3D"https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b=\r\n8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt=3D"De Koninklijke =\r\nLoop"></span><span style=3D'font-family:"Segoe =\r\nUI",sans-serif;color:white'><o:p></o:p></span></p><h1 align=3Dcenter =\r\nstyle=3D'margin:0cm;text-align:center;background:#FF9328'><span =\r\nstyle=3D'font-size:18.0pt;font-family:"Segoe =\r\nUI",sans-serif;color:white'>Bedankt voor je =\r\naanmelding!<o:p></o:p></span></h1><div><p =\r\nstyle=3D'background:white'><span style=3D'font-family:"Segoe =\r\nUI",sans-serif;color:#374151'>Beste Mayke Rood,<o:p></o:p></span></p><p =\r\nstyle=3D'background:white'><span style=3D'font-family:"Segoe =\r\nUI",sans-serif;color:#374151'>Bedankt voor je inschrijving voor De =\r\nKoninklijke Loop. We hebben je aanmelding in goede orde =\r\nontvangen.<o:p></o:p></span></p><div style=3D'border:solid #FFEDD5 =\r\n1.0pt;padding:12.0pt 12.0pt 12.0pt =\r\n12.0pt;margin-top:12.0pt;margin-bottom:12.0pt'><p class=3DMsoNormal =\r\nstyle=3D'background:#FFF7ED'><strong><span style=3D'font-family:"Segoe =\r\nUI",sans-serif;color:#9A3412'>Je inschrijfgegevens:</span></strong><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif;color:#9A3412'><br>Naam: =\r\nMayke Rood<br>Email: <a =\r\nhref=3D"mailto:rood1960@gmail.com">rood1960@gmail.com</a><br>Rol: =\r\nDeelnemer<br>Afstand: 15 KM<br><br>Ondersteuning: =\r\nNee<o:p></o:p></span></p></div><p style=3D'background:white'><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif;color:#374151'>Heb je vragen =\r\nover je inschrijving? Bekijk dan onze website voor meer =\r\ninformatie:<o:p></o:p></span></p><ul type=3Ddisc><li class=3DMsoNormal =\r\nstyle=3D'color:#374151;mso-margin-top-alt:auto;mso-margin-bottom-alt:auto=\r\n;mso-list:l0 level1 lfo1;background:white'><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif'><a =\r\nhref=3D"https://dekoninklijkeloop.nl/faq"><span =\r\nstyle=3D'color:#FF9328'>Veelgestelde =\r\nvragen</span></a><o:p></o:p></span></li><li class=3DMsoNormal =\r\nstyle=3D'color:#374151;mso-margin-top-alt:auto;mso-margin-bottom-alt:auto=\r\n;mso-list:l0 level1 lfo1;background:white'><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif'><a =\r\nhref=3D"https://dekoninklijkeloop.nl/over-ons"><span =\r\nstyle=3D'color:#FF9328'>Over De Koninklijke =\r\nLoop</span></a><o:p></o:p></span></li></ul><div =\r\nstyle=3D'margin-top:12.0pt'><p class=3DMsoNormal align=3Dcenter =\r\nstyle=3D'text-align:center;background:white'><span =\r\nstyle=3D'font-family:"Segoe UI",sans-serif;color:#374151'><a =\r\nhref=3D"https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/"><=\r\nspan style=3D'color:#FF9328;text-decoration:none'>Facebook</span></a> <a =\r\nhref=3D"https://www.instagram.com/koninklijkeloop/"><span =\r\nstyle=3D'color:#FF9328;text-decoration:none'>Instagram</span></a> =\r\n<o:p></o:p></span></p></div></div><p align=3Dcenter =\r\nstyle=3D'text-align:center;background:white'><span =\r\nstyle=3D'font-size:10.5pt;font-family:"Segoe =\r\nUI",sans-serif;color:#6B7280'>Met sportieve groet,<br>Team De =\r\nKoninklijke Loop<o:p></o:p></span></p><p align=3Dcenter =\r\nstyle=3D'text-align:center;background:white'><span =\r\nstyle=3D'font-size:10.5pt;font-family:"Segoe =\r\nUI",sans-serif;color:#6B7280'>=C2=A9 2025 De Koninklijke Loop. Alle =\r\nrechten =\r\nvoorbehouden.<o:p></o:p></span></p></div></div></div></body></html>\r\n------=_NextPart_000_0124_01DBBF52.C6A3AB00--\r\n\r\n	multipart/alternative; boundary="----=_NextPart_000_0124_01DBBF52.C6A3AB00"	2025-05-07 13:20:17	127	info	t	2025-05-07 18:34:14.782289	2025-05-07 11:32:14.559632	2025-05-07 18:34:14.782776
a243ccb0-9955-4870-a7c0-e32462f5a9bf	<PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FA@PSAPR06MB4198.apcprd06.prod.outlook.com>	Jose butler <josebutlerSEOexpert22@hotmail.com>	info@dekoninklijkeloop.nl	Re:	--_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_\r\nContent-Type: text/plain; charset="Windows-1252"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nHello,\r\n\r\nI=92m writing to follow up. Since I haven't received a reply, I assume\r\n\r\nyou are either busy and haven't had a chance to reply.\r\n\r\nCan I send you our whole proposal, price list, and best SEO quote?\r\n\r\nThank You=85!!\r\n\r\n________________________________\r\nFrom: Jose butler <josebutlerSEOexpert22@hotmail.com>\r\nSent: 27, May 2025 04:00 PM\r\nTo: Jose butler <josebutlerSEOexpert22@hotmail.com>\r\nSubject: Re: Price list.?\r\n\r\nHi,\r\n\r\nI found your details on Google and I have looked at your website.\r\n\r\nI can help bring your website to the first page of Google.\r\n\r\nIf you were on page #1, you'd get so many new more customers.\r\n\r\nCan I send you a price & Quote? If you are interested.\r\n\r\nThanks,\r\nJose butler.\r\nBusiness Executive || SEO Expert\r\n\r\n\r\n--_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_\r\nContent-Type: text/html; charset="Windows-1252"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html>\r\n<head>\r\n<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3DWindows-1=\r\n252">\r\n<style type=3D"text/css" style=3D"display:none;"> P {margin-top:0;margin-bo=\r\nttom:0;} </style>\r\n</head>\r\n<body dir=3D"ltr">\r\n<div class=3D"elementToProof" style=3D"font-family: Calibri, Helvetica, san=\r\ns-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nHello,</div>\r\n<div id=3D"appendonsend"></div>\r\n<div id=3D"divRplyFwdMsg"></div>\r\n<div id=3D"x_appendonsend"></div>\r\n<div id=3D"x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_appendonsend"></div>\r\n<div id=3D"x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div=\r\n>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></di=\r\nv>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></d=\r\niv>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></=\r\ndiv>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"><=\r\n/div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI=92m writing to follow up. Since I haven't received a reply, I assume</div=\r\n>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nyou are either busy and haven't had a chance to reply.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nCan I send you our whole proposal, price list, and best SEO quote?</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nThank You=85!!</div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"background-color: rgb(255, 255, 255); margin: 0px;"></div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, Helvetica, =\r\nsans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<hr style=3D"direction: ltr; text-align: left; text-indent: 0px; background=\r\n-color: rgb(255, 255, 255); display: inline-block; width: 587.986px;">\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, sans-serif;=\r\n font-size: 14.6667px; color: rgb(0, 0, 0);">\r\n<span style=3D"background-color: rgb(255, 255, 255);"><b>From:</b>&nbsp;Jos=\r\ne butler &lt;josebutlerSEOexpert22@hotmail.com&gt;</span></div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, sans-serif;=\r\n font-size: 14.6667px; color: rgb(0, 0, 0);">\r\n<span style=3D"background-color: rgb(255, 255, 255);"><b>Sent:</b>&nbsp;27,=\r\n May 2025 04:00 PM&nbsp;</span></div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, sans-serif;=\r\n font-size: 14.6667px; color: rgb(0, 0, 0);">\r\n<span style=3D"background-color: rgb(255, 255, 255);"><b>To:</b>&nbsp;Jose =\r\nbutler &lt;josebutlerSEOexpert22@hotmail.com&gt;</span></div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, sans-serif;=\r\n font-size: 14.6667px; color: rgb(0, 0, 0);">\r\n<span style=3D"background-color: rgb(255, 255, 255);"><b>Subject:</b>&nbsp;=\r\nRe: Price list.?</span></div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, Helvetica, =\r\nsans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, Helvetica, =\r\nsans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nHi,</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI found your details on Google and I have looked at your website.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI can help bring your website to the first page of Google.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nIf you were on page #1, you'd get so many new more customers.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nCan I send you a price &amp; Quote? If you are interested.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nThanks,</div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, Helvetica, =\r\nsans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nJose butler.</div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; backgroun=\r\nd-color: rgb(255, 255, 255); margin: 0px; font-family: Calibri, Helvetica, =\r\nsans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nBusiness Executive || SEO Expert</div>\r\n<div style=3D"direction: ltr; font-family: Aptos, Aptos_EmbeddedFont, Aptos=\r\n_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; color: rgb=\r\n(0, 0, 0);">\r\n<br>\r\n</div>\r\n</body>\r\n</html>\r\n\r\n--_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_--\r\n	multipart/alternative; boundary="_000_PSAPR06MB41989C59EA54CF3BE8FCE0BCB36FAPSAPR06MB4198apcp_"	2025-06-05 10:31:25	131	info	t	2025-09-10 18:40:14.337879	2025-06-05 10:32:14.500151	2025-09-10 18:40:14.338289
b3bc62fd-189f-4fe6-a501-421a55650bcd	124005b2-2b54-11f0-a697-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Danny <br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 15:00:54	90	inschrijving	t	2025-09-10 18:40:23.308332	2025-05-07 15:02:14.118959	2025-09-10 18:40:23.308789
5d4df3db-eede-4e27-b9f0-e4f778ec4268	82ba22d2-252c-11f0-a371-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Ismail <br>\n                    Email: lidaahmadi99@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 19:02:37	81	inschrijving	t	2025-10-05 18:24:37.415999	2025-04-29 19:17:14.586288	2025-10-05 18:24:37.416281
9a197a84-9a30-41c2-9347-933239a02c3a	f8edecca-2b53-11f0-94ad-471e5e6299d9@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: john<br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 15:00:11	89	inschrijving	t	2025-05-07 18:34:50.497515	2025-05-07 15:02:14.116658	2025-05-07 18:34:50.501371
2af54bd6-74ad-4ddf-99f4-eb05767008e7	<SL2P216MB1803400813FC81C04D237850BF70A@SL2P216MB1803.KORP216.PROD.OUTLOOK.COM>	laccy lark <laccylark3462@outlook.com>	info@dekoninklijkeloop.nl	Re: Price list.?	--_000_SL2P216MB1803400813FC81C04D237850BF70ASL2P216MB1803KORP_\r\nContent-Type: text/plain; charset="Windows-1252"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nHello,\r\n\r\nI=92m writing to follow up. Since I haven't received a reply, I assume\r\n\r\nyou are either busy and haven't had a chance to reply.\r\n\r\nCan I send you our whole proposal, price list, and best SEO quote?\r\n\r\nThank You=85!!\r\n\r\n________________________________\r\nFrom: laccy lark <laccylark3462@outlook.com>\r\nSent: Tuesday, June 03, 2025, 03:55 PM\r\nTo: laccy lark <laccylark3462@outlook.com>\r\nSubject: Re: YES Costing..?\r\n\r\nHi,\r\n\r\nI found your details on Google and I have looked at your website.\r\n\r\nI can help bring your website to the first page of Google.\r\n\r\nIf you were on page #1, you'd get so many new more customers.\r\n\r\nCan I send you a SEO price & Quote? If you are interested.\r\n\r\nThanks,\r\n[laccy lark]\r\n\r\n--_000_SL2P216MB1803400813FC81C04D237850BF70ASL2P216MB1803KORP_\r\nContent-Type: text/html; charset="Windows-1252"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html>\r\n<head>\r\n<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3DWindows-1=\r\n252">\r\n<style type=3D"text/css" style=3D"display:none;"> P {margin-top:0;margin-bo=\r\nttom:0;} </style>\r\n</head>\r\n<body dir=3D"ltr">\r\n<div class=3D"elementToProof" style=3D"font-family: Calibri, Helvetica, san=\r\ns-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nHello,</div>\r\n<div id=3D"appendonsend"></div>\r\n<div id=3D"divRplyFwdMsg"></div>\r\n<div id=3D"x_appendonsend"></div>\r\n<div id=3D"x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_appendonsend"></div>\r\n<div id=3D"x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_appendonsend"></div>\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_divRplyFwdMsg"></div=\r\n>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI=92m writing to follow up. Since I haven't received a reply, I assume</div=\r\n>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nyou are either busy and haven't had a chance to reply.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nCan I send you our whole proposal, price list, and best SEO quote?</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n&nbsp;</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nThank You=85!!</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<hr style=3D"direction: ltr; display: inline-block; width: 98%;">\r\n<div id=3D"x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_=\r\nx_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x_x=\r\n_x_x_x_x_x_x_x_x_x_divRplyFwdMsg">\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<b>From:</b>&nbsp;laccy lark &lt;laccylark3462@outlook.com&gt;<br>\r\n<b>Sent:</b>&nbsp;Tuesday, June 03, 2025, 03:55 PM<br>\r\n<b>To:</b>&nbsp;laccy lark &lt;laccylark3462@outlook.com&gt;<br>\r\n<b>Subject:</b>&nbsp;Re: YES Costing..?</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt;">\r\n&nbsp;</div>\r\n</div>\r\n<div style=3D"direction: ltr; text-align: left; text-indent: 0px; font-fami=\r\nly: Calibri, Helvetica, sans-serif; font-size: 11pt; color: rgb(0, 0, 0);">\r\nHi,</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI found your details on Google and I have looked at your website.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nI can help bring your website to the first page of Google.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nIf you were on page #1, you'd get so many new more customers.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nCan I send you a SEO price &amp; Quote? If you are interested.</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\nThanks,</div>\r\n<div style=3D"direction: ltr; font-family: Calibri, Helvetica, sans-serif; =\r\nfont-size: 11pt; color: rgb(0, 0, 0);">\r\n[laccy lark]</div>\r\n</body>\r\n</html>\r\n\r\n--_000_SL2P216MB1803400813FC81C04D237850BF70ASL2P216MB1803KORP_--\r\n	multipart/alternative; boundary="_000_SL2P216MB1803400813FC81C04D237850BF70ASL2P216MB1803KORP_"	2025-06-16 10:34:06	132	info	t	2025-09-10 18:40:15.322602	2025-06-16 10:47:14.39863	2025-09-10 18:40:15.323017
bff9d198-6a77-4b55-9175-bdad81f6c209	c244ed4c-2b53-11f0-9175-4d2191f5f3b5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: michel<br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 14:58:40	87	inschrijving	t	2025-09-30 18:21:33.141308	2025-05-07 15:02:14.110806	2025-09-30 18:21:33.142447
3eadecb4-5be4-4462-8938-ebf81ea122fc	5b146f63-25aa-11f0-a210-152d8afab6bc@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: laura<br>\n                    Email: steun94@xs4all.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 10 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-30 10:03:27	83	inschrijving	t	2025-10-05 18:24:38.395433	2025-04-30 10:17:14.88272	2025-10-05 18:24:38.396559
da242ef4-7e39-456e-a179-a87c6c78b4b4	bf2cdb2c-2b6b-11f0-96ec-2b8368a4d5c5@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Pieter Streefland<br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 17:50:22	91	inschrijving	t	2025-05-07 18:34:48.021165	2025-05-07 18:02:14.35081	2025-05-07 18:34:48.022027
57513be8-6331-4598-a210-23c583aac09c	826069c8-252b-11f0-8663-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Saleem<br>\n                    Email: lidaahmadi99@gmail.com<br>\n                    <br>\n                    Telefoon: 0649020648<br>\n                    \n                    Rol: Begeleider<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 18:55:25	77	inschrijving	t	2025-10-05 18:24:35.538279	2025-04-29 19:01:55.793745	2025-10-05 18:24:35.538554
4ada6cc2-715c-4451-bea8-1fc2a49a877c	9acc96b7-2b74-11f0-b99c-db70b48e0ecf@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: yassine<br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 18:53:48	92	inschrijving	t	2025-09-30 18:21:26.935136	2025-05-07 19:00:27.146555	2025-09-30 18:21:26.939557
79020c31-3e4b-43c5-af82-011d6fb9dd35	25bb9762-252c-11f0-9150-db70b48e0ecf@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Yussef<br>\n                    Email: lidaahmadi99@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 18:59:59	78	inschrijving	t	2025-10-05 18:24:36.603537	2025-04-29 19:01:55.795972	2025-10-05 18:24:36.603851
95ab6929-0326-407a-a293-bfd02278f0c0	93dde194-196b-11f0-ad41-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Henk Rekers <br>\n                    Email: h.rekers59@kpnmail.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-14 20:03:49	70	inschrijving	t	2025-04-20 12:51:02.157161	2025-04-20 00:03:01.163228	2025-04-20 12:51:02.161598
70cd0e5f-445e-4109-8bbf-7248934a7c61	<DB9PR08MB6811E02FC35DC425302C42E9A38BA@DB9PR08MB6811.eurprd08.prod.outlook.com>	Geer Hubers <hub3008@gmail.com>	info@dekoninklijkeloop.nl	Re: Bedankt voor je aanmelding	--_000_DB9PR08MB6811E02FC35DC425302C42E9A38BADB9PR08MB6811eurp_\r\nContent-Type: text/plain; charset="iso-8859-1"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nGoedemiddag,\r\n\r\ndoor onvoorziene omstandigheden moet ik mij helaas afmelden voor de wandelt=\r\nocht op 17 mei aanstaande.\r\nIk wens u allen een succesvol evenement toe.\r\n\r\nMet vriendelijke groet,\r\nGeer Hubers.\r\n________________________________\r\nVan: info@dekoninklijkeloop.nl <info@dekoninklijkeloop.nl>\r\nVerzonden: dinsdag 29 april 2025 11:28\r\nAan: hub3008@gmail.com <hub3008@gmail.com>\r\nOnderwerp: Bedankt voor je aanmelding\r\n\r\n[De Koninklijke Loop]\r\nBedankt voor je aanmelding!\r\n\r\nBeste Geer Hubers,\r\n\r\nBedankt voor je inschrijving voor De Koninklijke Loop. We hebben je aanmeld=\r\ning in goede orde ontvangen.\r\n\r\nJe inschrijfgegevens:\r\nNaam: Geer Hubers\r\nEmail: hub3008@gmail.com\r\nRol: Deelnemer\r\nAfstand: 15 KM\r\n\r\nOndersteuning: Nee\r\n\r\nHeb je vragen over je inschrijving? Bekijk dan onze website voor meer infor=\r\nmatie:\r\n\r\n  *   Veelgestelde vragen<https://dekoninklijkeloop.nl/faq>\r\n  *   Over De Koninklijke Loop<https://dekoninklijkeloop.nl/over-ons>\r\n\r\nFacebook<https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/> In=\r\nstagram<https://www.instagram.com/koninklijkeloop/>\r\n\r\nMet sportieve groet,\r\nTeam De Koninklijke Loop\r\n\r\n=A9 2025 De Koninklijke Loop. Alle rechten voorbehouden.\r\n\r\n--_000_DB9PR08MB6811E02FC35DC425302C42E9A38BADB9PR08MB6811eurp_\r\nContent-Type: text/html; charset="iso-8859-1"\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html>\r\n<head>\r\n<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Diso-8859-=\r\n1">\r\n<style type=3D"text/css" style=3D"display:none;"> P {margin-top:0;margin-bo=\r\nttom:0;} </style>\r\n</head>\r\n<body dir=3D"ltr">\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\nGoedemiddag,</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\ndoor onvoorziene omstandigheden moet ik mij helaas afmelden voor de wandelt=\r\nocht op 17 mei aanstaande.<br>\r\nIk wens u allen een succesvol evenement toe.</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\nMet vriendelijke groet,</div>\r\n<div class=3D"elementToProof" style=3D"font-family: Aptos, Aptos_EmbeddedFo=\r\nnt, Aptos_MSFontService, Calibri, Helvetica, sans-serif; font-size: 12pt; c=\r\nolor: rgb(0, 0, 0);">\r\nGeer Hubers.</div>\r\n<div id=3D"appendonsend"></div>\r\n<hr style=3D"display:inline-block;width:98%" tabindex=3D"-1">\r\n<div id=3D"divRplyFwdMsg" dir=3D"ltr"><font face=3D"Calibri, sans-serif" st=\r\nyle=3D"font-size:11pt" color=3D"#000000"><b>Van:</b> info@dekoninklijkeloop=\r\n.nl &lt;info@dekoninklijkeloop.nl&gt;<br>\r\n<b>Verzonden:</b> dinsdag 29 april 2025 11:28<br>\r\n<b>Aan:</b> hub3008@gmail.com &lt;hub3008@gmail.com&gt;<br>\r\n<b>Onderwerp:</b> Bedankt voor je aanmelding</font>\r\n<div>&nbsp;</div>\r\n</div>\r\n<style>\r\n<!--\r\ndiv\r\n\t{font-family:'Inter',-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,He=\r\nlvetica,Arial,sans-serif;\r\n\tline-height:1.5;\r\n\tcolor:#374151;\r\n\tmargin:0;\r\n\tpadding:0;\r\n\tbackground-color:#f3f4f6}\r\n.x_container\r\n\t{max-width:600px;\r\n\tmargin:0 auto;\r\n\tpadding:20px}\r\n.x_card\r\n\t{background-color:#ffffff;\r\n\tborder-radius:12px;\r\n\tbox-shadow:0 4px 6px -1px rgba(0,0,0,0.1),0 2px 4px -1px rgba(0,0,0,0.06);\r\n\toverflow:hidden}\r\n.x_header\r\n\t{background-color:#ff9328;\r\n\tcolor:#ffffff;\r\n\tpadding:24px;\r\n\ttext-align:center}\r\n.x_logo\r\n\t{height:60px;\r\n\tmargin-bottom:16px;\r\n\tdisplay:block;\r\n\tmargin-left:auto;\r\n\tmargin-right:auto}\r\n.x_content\r\n\t{padding:24px}\r\n.x_message-box\r\n\t{background-color:#fff7ed;\r\n\tborder:1px solid #ffedd5;\r\n\tborder-radius:8px;\r\n\tpadding:16px;\r\n\tmargin:16px 0;\r\n\tcolor:#9a3412}\r\n.x_footer\r\n\t{text-align:center;\r\n\tpadding:24px;\r\n\tcolor:#6b7280;\r\n\tfont-size:14px}\r\n.x_social-links\r\n\t{margin-top:16px;\r\n\ttext-align:center}\r\n.x_social-link\r\n\t{display:inline-block;\r\n\tmargin:0 8px;\r\n\tcolor:#ff9328;\r\n\ttext-decoration:none}\r\n-->\r\n</style>\r\n<div>\r\n<div class=3D"x_container">\r\n<div class=3D"x_card">\r\n<div class=3D"x_header"><img alt=3D"De Koninklijke Loop" class=3D"x_logo" s=\r\ntyle=3D"max-width:200px; width:100%; height:auto" src=3D"https://res.cloudi=\r\nnary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_076084=\r\n9fb8_yn6vdm.png">\r\n<h1 style=3D"margin:0; font-size:24px; font-weight:700">Bedankt voor je aan=\r\nmelding!</h1>\r\n</div>\r\n<div class=3D"x_content">\r\n<p>Beste Geer Hubers,</p>\r\n<p>Bedankt voor je inschrijving voor De Koninklijke Loop. We hebben je aanm=\r\nelding in goede orde ontvangen.</p>\r\n<div class=3D"x_message-box"><strong>Je inschrijfgegevens:</strong><br>\r\nNaam: Geer Hubers<br>\r\nEmail: hub3008@gmail.com<br>\r\nRol: Deelnemer<br>\r\nAfstand: 15 KM<br>\r\n<br>\r\nOndersteuning: Nee<br>\r\n</div>\r\n<p>Heb je vragen over je inschrijving? Bekijk dan onze website voor meer in=\r\nformatie:</p>\r\n<ul>\r\n<li><a href=3D"https://dekoninklijkeloop.nl/faq" style=3D"color:#ff9328">Ve=\r\nelgestelde vragen</a></li><li><a href=3D"https://dekoninklijkeloop.nl/over-=\r\nons" style=3D"color:#ff9328">Over De Koninklijke Loop</a></li></ul>\r\n<div class=3D"x_social-links"><a href=3D"https://www.facebook.com/p/De-Koni=\r\nnklijke-Loop-61556315443279/" class=3D"x_social-link">Facebook</a>\r\n<a href=3D"https://www.instagram.com/koninklijkeloop/" class=3D"x_social-li=\r\nnk">Instagram</a>\r\n</div>\r\n</div>\r\n<div class=3D"x_footer">\r\n<p>Met sportieve groet,<br>\r\nTeam De Koninklijke Loop</p>\r\n<p>=A9 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\r\n</div>\r\n</div>\r\n</div>\r\n</div>\r\n</body>\r\n</html>\r\n\r\n--_000_DB9PR08MB6811E02FC35DC425302C42E9A38BADB9PR08MB6811eurp_--\r\n	multipart/alternative; boundary="_000_DB9PR08MB6811E02FC35DC425302C42E9A38BADB9PR08MB6811eurp_"	2025-05-08 12:12:55	128	info	t	2025-05-09 13:24:32.3847	2025-05-08 12:17:14.648742	2025-05-09 13:24:32.3886
ae7d9324-c4db-4c1a-8010-598a13e973f3	<992062149.1173885.1750671302160@mail.yahoo.com>	Sem Karan <karan.sem@aol.com>	info@dekoninklijkeloop.nl	Re: SEO Service & Plans	------=_Part_1173884_1593284664.1750671302159\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 7bit\r\n\r\n Hallo,\r\nIk wacht nog steeds op uw reactie.\r\nKan ik u een offerte en prijsopgave sturen?\r\nMet vriendelijke groet"\r\n\r\n    On Thursday, April 17, 2025 at 12:47:52 PM GMT+5:30, Sem Karan <karan.sem@aol.com> wrote:   \r\n\r\n Hallo,\r\nIk heb je website geanalyseerd en vind hem geweldig, maar hij scoort niet goed in de zoekmachines.\r\nWil je je website optimaliseren voor meer zichtbaarheid en meer klanten? Je zult al in de eerste maand goede resultaten zien.\r\nKan ik je een website-auditrapport, offerte en prijslijst voor onze werkzaamheden sturen?\r\nBedankt.\r\n    \r\n------=_Part_1173884_1593284664.1750671302159\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<html><head></head><body><div class=3D"ydpe04a3113yahoo-style-wrap" style=\r\n=3D"font-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-size:16px=\r\n;"><div></div>\r\n        </div><div id=3D"ydp6555334fyahoo_quoted_1544112677" class=3D"ydp65=\r\n55334fyahoo_quoted"><div class=3D"ydp6555334fyahoo-style-wrap" style=3D"fon=\r\nt-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-size:16px;"><div=\r\n style=3D"font-family:'Helvetica Neue', Helvetica, Arial, sans-serif;font-s=\r\nize:13px;color:#26282a;border-left: 1px solid #ccc;padding-left: 8px;margin=\r\n: 0px 0px 0px 8px" class=3D"ydp6555334finline_reply_quote_container" data-s=\r\nplit-quote-node=3D"true"><div><div id=3D"ydp6555334fyiv0717243290"><div><di=\r\nv style=3D"font-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-si=\r\nze:16px;" class=3D"ydp6555334fyiv0717243290ydp6f843dcfyahoo-style-wrap"><di=\r\nv dir=3D"ltr"><div><div style=3D"color:rgb(0, 0, 0);font-family:Helvetica N=\r\neue, Helvetica, Arial, sans-serif;font-size:16px;letter-spacing:-0.32px;">H=\r\nallo,</div><div style=3D"color:rgb(0, 0, 0);font-family:Helvetica Neue, Hel=\r\nvetica, Arial, sans-serif;font-size:16px;letter-spacing:-0.32px;"><br clear=\r\n=3D"none"></div><div style=3D"color:rgb(0, 0, 0);font-family:Helvetica Neue=\r\n, Helvetica, Arial, sans-serif;font-size:16px;letter-spacing:-0.32px;">Ik w=\r\nacht nog steeds op uw reactie.</div><div style=3D"color:rgb(0, 0, 0);font-f=\r\namily:Helvetica Neue, Helvetica, Arial, sans-serif;font-size:16px;letter-sp=\r\nacing:-0.32px;"><br clear=3D"none"></div><div style=3D"color:rgb(0, 0, 0);f=\r\nont-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-size:16px;lett=\r\ner-spacing:-0.32px;"><b style=3D"background-color:rgb(253, 239, 43);">Kan i=\r\nk u een offerte en prijsopgave sturen?</b></div><div style=3D"color:rgb(0, =\r\n0, 0);font-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-size:16=\r\npx;letter-spacing:-0.32px;"><br clear=3D"none"></div><div style=3D"color:rg=\r\nb(0, 0, 0);font-family:Helvetica Neue, Helvetica, Arial, sans-serif;font-si=\r\nze:16px;letter-spacing:-0.32px;">Met vriendelijke groet"</div></div><br cle=\r\nar=3D"none"></div><div><br clear=3D"none"></div>\r\n       =20\r\n        </div><div id=3D"ydp6555334fyiv0717243290yqt95010" class=3D"ydp6555=\r\n334fyiv0717243290yqt0193128104"><div id=3D"ydp6555334fyiv0717243290ydp2543f=\r\n8c5yahoo_quoted_1024091124" class=3D"ydp6555334fyiv0717243290ydp2543f8c5yah=\r\noo_quoted"><div style=3D"font-family:Helvetica Neue, Helvetica, Arial, sans=\r\n-serif;font-size:16px;" class=3D"ydp6555334fyiv0717243290ydp2543f8c5yahoo-s=\r\ntyle-wrap">\r\n            <div style=3D"font-family:'Helvetica Neue', Helvetica, Arial, s=\r\nans-serif;font-size:13px;color:#26282a;">\r\n               =20\r\n                <div class=3D"ydp6555334fyiv0717243290ydp2543f8c5quoted-tex=\r\nt-header">\r\n                        On Thursday, April 17, 2025 at 12:47:52 PM GMT+5:30=\r\n, Sem Karan &lt;karan.sem@aol.com&gt; wrote:\r\n                    </div>\r\n                </div><div style=3D"font-family:'Helvetica Neue', Helvetica=\r\n, Arial, sans-serif;font-size:13px;color:#26282a;border-left:1px solid #ccc=\r\n;padding-left:8px;margin:0px 0px 0px 8px;" class=3D"ydp6555334fyiv071724329=\r\n0ydp2543f8c5inline_reply_quote_container">\r\n                <div><br clear=3D"none"></div><div><br clear=3D"none"></div=\r\n>\r\n                <div><div id=3D"ydp6555334fyiv0717243290ydp2543f8c5yiv98029=\r\n02998"><div><div style=3D"font-family:Helvetica Neue, Helvetica, Arial, san=\r\ns-serif;font-size:16px;" class=3D"ydp6555334fyiv0717243290ydp2543f8c5yiv980=\r\n2902998yahoo-style-wrap"><div dir=3D"ltr"><div><div style=3D"font-size:16px=\r\n;letter-spacing:-0.32px;font-family:Helvetica, Arial, sans-serif;">Hallo,</=\r\ndiv><div style=3D"font-size:16px;letter-spacing:-0.32px;font-family:Helveti=\r\nca, Arial, sans-serif;"><br clear=3D"none"></div><div style=3D"font-size:16=\r\npx;letter-spacing:-0.32px;font-family:Helvetica, Arial, sans-serif;">Ik heb=\r\n je website geanalyseerd en vind hem geweldig, maar hij scoort niet goed in=\r\n de zoekmachines.</div><div style=3D"font-size:16px;letter-spacing:-0.32px;=\r\nfont-family:Helvetica, Arial, sans-serif;"><br clear=3D"none"></div><div st=\r\nyle=3D"font-size:16px;letter-spacing:-0.32px;font-family:Helvetica, Arial, =\r\nsans-serif;">Wil je je website optimaliseren voor meer zichtbaarheid en mee=\r\nr klanten? Je zult al in de eerste maand goede resultaten zien.</div><div s=\r\ntyle=3D"font-size:16px;letter-spacing:-0.32px;font-family:Helvetica, Arial,=\r\n sans-serif;"><br clear=3D"none"></div><div style=3D"font-size:16px;letter-=\r\nspacing:-0.32px;font-family:Helvetica, Arial, sans-serif;">Kan ik je een we=\r\nbsite-auditrapport, offerte en prijslijst voor onze werkzaamheden sturen?</=\r\ndiv><div style=3D"font-size:16px;letter-spacing:-0.32px;font-family:Helveti=\r\nca, Arial, sans-serif;"><br clear=3D"none"></div><div style=3D"font-size:16=\r\npx;letter-spacing:-0.32px;font-family:Helvetica, Arial, sans-serif;">Bedank=\r\nt.</div></div><br clear=3D"none"></div></div></div></div></div>\r\n            </div>\r\n        </div></div></div></div></div></div>\r\n            </div>\r\n        </div></div></body></html>\r\n------=_Part_1173884_1593284664.1750671302159--\r\n	multipart/alternative; boundary="----=_Part_1173884_1593284664.1750671302159"	2025-06-23 09:35:02	133	info	t	2025-09-10 18:39:39.245788	2025-06-23 09:37:40.960615	2025-09-10 18:39:39.248079
391f3d11-2199-42a9-8d56-4aa6d0cbfdac	11d2d338-2675-11f0-be40-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Theodora Naus<br>\n                    Email: tgemooi@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-01 10:14:30	85	inschrijving	t	2025-10-05 18:24:39.050767	2025-05-01 10:17:14.22501	2025-10-05 18:24:39.051051
3bcdfe35-b0bf-4705-b5f3-458fc0c855b6	469bd6d1-2676-11f0-b122-db70b48e0ecf@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Theodora Naus<br>\n                    Email: tgemooi@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-01 10:23:10	86	inschrijving	t	2025-10-05 18:24:40.428432	2025-05-01 10:32:14.365326	2025-10-05 18:24:40.428767
177f5e53-906c-4268-a766-4629d86fc3e0	29197208-196b-11f0-978b-417246ffdc90@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Hilde Rekers <br>\n                    Email: h.rekers59@kpnmail.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 6 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-14 20:00:50	69	inschrijving	t	2025-04-20 23:35:27.990681	2025-04-20 00:03:01.161129	2025-04-20 23:35:27.993393
8e24c875-d50e-4c01-9ccb-2fdf43dd194a	858b0110-28b3-11f0-8cb5-d7c209f8bd06@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	Nieuw contactformulier	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuw contactformulier bericht</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuw bericht binnengekomen via het contactformulier.</p>\n                \n                <div>\n                    <strong>Contactgegevens:</strong><br>\n                    Naam: Luciēnne Tang<br>\n                    Email: luciennetang@gmail.com<br>\n                    <br>\n                    <strong>Bericht:</strong><br>\n                    Annulering inschrijving.\n\nMijn inschrijfgegevens zijn:\nNaam : Luciēnne \nEmail: luciennetang@gmail.com \nRol: deelnemer\nAfstand: 15 KM\n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-04 06:46:36	125	info	t	2025-05-06 08:54:08.518699	2025-05-04 06:47:13.979529	2025-05-06 08:54:08.519197
d8faceca-4604-49aa-b615-dfbcaa30f22c	435e219c-24dc-11f0-9d26-b37c246f863f@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Sophie Hubers<br>\n                    Email: sophieehubers@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 09:28:11	75	inschrijving	t	2025-04-29 16:11:58.63309	2025-04-29 09:32:14.544255	2025-04-29 16:11:58.633649
55da93f8-5b3c-4b2a-b70f-b753a1841e2f	5d9fbad1-24dc-11f0-bf45-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Geer Hubers<br>\n                    Email: hub3008@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-29 09:28:55	76	inschrijving	t	2025-05-07 18:36:05.689569	2025-04-29 09:32:14.547863	2025-05-07 18:36:05.690062
999c3c5d-2b1b-407f-af63-f3ddbeec963b	66d336f1-2605-11f0-968c-29b2d794c87d@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: Luciënne <br>\n                    Email: luciennetang@gmail.com<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 15 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-04-30 20:55:10	84	inschrijving	t	2025-10-05 18:24:38.682725	2025-04-30 21:02:14.884409	2025-10-05 18:24:38.683071
93066ddf-c682-483a-909b-276a5dd00501	d85ac82e-2b53-11f0-b5e4-db70b48e0ecf@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	inschrijving@dekoninklijkeloop.nl	Nieuwe aanmelding ontvangen	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuwe aanmelding ontvangen</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuwe inschrijving binnengekomen voor De Koninklijke Loop.</p>\n                \n                <div>\n                    <strong>Inschrijfgegevens:</strong><br>\n                    Naam: peter<br>\n                    Email: enckerkamp.23@sheerenloo.nl<br>\n                    \n                    Rol: Deelnemer<br>\n                    Afstand: 10 KM<br>\n                    <br>\n                    Ondersteuning: Nee<br>\n                    \n                    \n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-05-07 14:59:18	88	inschrijving	t	2025-10-05 18:24:41.049282	2025-05-07 15:02:14.114048	2025-10-05 18:24:41.049691
6d66df18-c14a-4b01-bf26-b1834166e01c	<913de533-a855-11f0-bc17-d510462faafc@dekoninklijkeloop.nl>	info@dekoninklijkeloop.nl	info@dekoninklijkeloop.nl	Nieuw contactformulier	\n\n\n    \n    \n    \n    \n\n\n    <div>\n        <div>\n            <div>\n                <img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop">\n                <h1>Nieuw contactformulier bericht</h1>\n            </div>\n            \n            <div>\n                <p>Er is een nieuw bericht binnengekomen via het contactformulier.</p>\n                \n                <div>\n                    <strong>Contactgegevens:</strong><br>\n                    Naam: JeffTest<br>\n                    Email: laventejeffrey@gmail.com<br>\n                    <br>\n                    <strong>Bericht:</strong><br>\n                    Test bericht\n                </div>\n            </div>\n            \n            <div>\n                <p>Dit is een automatisch gegenereerd bericht van het DKL Email Systeem.</p>\n                <p>© 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n            </div>\n        </div>\n    </div>\n\n \n	text/html; charset=UTF-8	2025-10-13 16:56:33	134	info	t	2025-10-31 11:52:30.838748	2025-10-13 16:57:36.735328	2025-10-31 11:52:30.838615
\.


--
-- Data for Name: migraties; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.migraties (id, versie, naam, toegepast) FROM stdin;
1	1.0.0	Initiële database setup	2025-03-14 15:22:28.513509+00
2	1.0.1	Initiële data seeding	2025-03-14 15:22:28.710911+00
5	1.0.2	Update schema om overeen te komen met Go models	2025-03-14 17:23:21.010937+00
51	1.0.3	Create incoming_emails table	2025-03-18 00:52:08.913979+00
264	1.14	Toevoegen van extra aanmeldingen maart 2025	2025-03-30 11:06:33.518624+00
290	1.15	Toevoegen van aanmeldingen april 2025	2025-04-16 15:18:23.726407+00
429	1.16.0	Add chat tables	2025-10-04 17:05:57.249006+00
521	1.17.0	Add is_public to chat_channels	2025-10-05 17:56:36.356503+00
530	1.18.0	Add last_read_at to chat_channel_participants	2025-10-05 18:26:14.351692+00
648	1.20.0	Create RBAC tables for flexible permission management	2025-10-07 17:15:41.648064+00
649	1.21.0	Seed initial RBAC data based on current system	2025-10-07 17:15:41.774011+00
672	1.22.0	Assign admin role to existing admin user	2025-10-07 17:38:28.549894+00
685	1.23.0	Add staff access permission and assign to admin/staff roles	2025-10-07 18:06:57.054356+00
764	1.24.0	Add missing permissions for Photos, Albums, Partners, Sponsors, Videos	2025-10-07 20:00:54.648925+00
779	1.25.0	Assign new permissions for Photos, Albums, Partners, Sponsors, Videos to admin role	2025-10-07 20:10:40.949057+00
795	1.26.0	Assign read permissions for Photos, Albums, Partners, Sponsors, Videos to staff role	2025-10-07 20:20:03.555599+00
828	1.27.0	Assign staff role to jeffrey@dekoninklijkeloop.nl	2025-10-07 20:35:03.355284+00
1050	1.28.0	Add RBAC management permissions	2025-10-08 12:51:57.512242+00
1204	1.29.0	Add thumbnail_url to chat_messages	2025-10-10 18:39:18.649997+00
1205	1.30.0	Create uploaded_images table	2025-10-10 18:39:18.651275+00
1700	1.45.0	Add steps permissions for RBAC system	2025-10-24 19:24:48.264458+00
1936	1.34.0	Add steps permissions for deelnemer/begeleider/vrijwilliger roles	2025-10-25 15:20:36.749699+00
2043	1.47.0	Performance optimizations: FK indexes, compound indexes, partial indexes, and FTS	2025-10-30 23:56:33.646736+00
2154	1.48.0	Advanced optimizations: triggers, constraints, materialized views, data quality	2025-10-31 01:02:31.251438+00
\.


--
-- Data for Name: newsletters; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.newsletters (id, subject, content, sent_at, batch_id, created_at, updated_at) FROM stdin;
cf0c8d61-5f24-4f34-98e5-3bbb3c377592	Nieuwsbrief DKL25	<div style="text-align: center; margin-bottom: 24px;">\n<img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop" style="max-width: 200px; width: 100%; height: auto; margin-bottom: 16px;">\n</div>\n\n<h2 style="color: #ff9328; margin-bottom: 16px;">Welkom bij de nieuwsbrief van De Koninklijke Loop 2025</h2>\n\n<p>Beste leden en geïnteresseerden,</p>\n\n<p>Hier vindt u het laatste nieuws over De Koninklijke Loop 2025.</p>\n\n<div style="background-color: #fff7ed; border: 1px solid #ffedd5; border-radius: 8px; padding: 16px; margin: 20px 0; color: #9a3412;">\n<strong>Wat kunt u verwachten:</strong><br>\n• Actuele informatie over het evenement<br>\n• Inschrijvingsupdates<br>\n• Route-informatie<br>\n• Sponsormogelijkheden\n</div>\n\n<p>Blijf op de hoogte van alle ontwikkelingen!</p>\n\n<p>Met sportieve groet,<br><strong>Team De Koninklijke Loop</strong></p>\n\n<div style="text-align: center; margin-top: 24px; padding-top: 16px; border-top: 1px solid #e5e7eb;">\n<p style="color: #6b7280; font-size: 14px;">\n<a href="https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Facebook</a> |\n<a href="https://www.instagram.com/koninklijkeloop/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Instagram</a> |\n<a href="https://dekoninklijkeloop.nl" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Website</a>\n</p>\n<p style="color: #9ca3af; font-size: 12px; margin-top: 8px;">&copy; 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n</div>	2025-10-07 15:13:00.650778+00	newsletter_manual_cf0c8d61-5f24-4f34-98e5-3bbb3c377592	2025-10-07 12:53:55.459869+00	2025-10-07 15:13:00.651678+00
ec0fdb08-3dbe-425f-846c-cd9df59f1b30	Nieuwsbrief DKL25	<div style="text-align: center; margin-bottom: 24px;">\n<img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop" style="max-width: 200px; width: 100%; height: auto; margin-bottom: 16px;">\n</div>\n\n<h2 style="color: #ff9328; margin-bottom: 16px;">Welkom bij de nieuwsbrief van De Koninklijke Loop 2025</h2>\n\n<p>Beste leden en geïnteresseerden,</p>\n\n<p>Hier vindt u het laatste nieuws over De Koninklijke Loop 2025.</p>\n\n<div style="background-color: #fff7ed; border: 1px solid #ffedd5; border-radius: 8px; padding: 16px; margin: 20px 0; color: #9a3412;">\n<strong>Wat kunt u verwachten:</strong><br>\n• Actuele informatie over het evenement<br>\n• Inschrijvingsupdates<br>\n• Route-informatie<br>\n• Sponsormogelijkheden\n</div>\n\n<p>Blijf op de hoogte van alle ontwikkelingen!</p>\n\n<p>Met sportieve groet,<br><strong>Team De Koninklijke Loop</strong></p>\n\n<div style="text-align: center; margin-top: 24px; padding-top: 16px; border-top: 1px solid #e5e7eb;">\n<p style="color: #6b7280; font-size: 14px;">\n<a href="https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Facebook</a> |\n<a href="https://www.instagram.com/koninklijkeloop/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Instagram</a> |\n<a href="https://dekoninklijkeloop.nl" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Website</a>\n</p>\n<p style="color: #9ca3af; font-size: 12px; margin-top: 8px;">&copy; 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n</div>	2025-10-07 15:21:54.328914+00	newsletter_manual_ec0fdb08-3dbe-425f-846c-cd9df59f1b30	2025-10-07 15:16:43.779633+00	2025-10-07 15:21:54.329701+00
b2144362-495e-43f5-8761-f5378645fc86	Nieuwsbrief DKL25	<div style="text-align: center; margin-bottom: 24px;">\n<img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop" style="max-width: 200px; width: 100%; height: auto; margin-bottom: 16px;">\n</div>\n\n<h2 style="color: #ff9328; margin-bottom: 16px;">Welkom bij de nieuwsbrief van De Koninklijke Loop 2025</h2>\n\n<p>Beste leden en geïnteresseerden,</p>\n\n<p>Hier vindt u het laatste nieuws over De Koninklijke Loop 2025.</p>\n\n<div style="background-color: #fff7ed; border: 1px solid #ffedd5; border-radius: 8px; padding: 16px; margin: 20px 0; color: #9a3412;">\n<strong>Wat kunt u verwachten:</strong><br>\n• Actuele informatie over het evenement<br>\n• Inschrijvingsupdates<br>\n• Route-informatie<br>\n• Sponsormogelijkheden\n</div>\n\n<p>Blijf op de hoogte van alle ontwikkelingen!</p>\n\n<p>Met sportieve groet,<br><strong>Team De Koninklijke Loop</strong></p>\n\n<div style="text-align: center; margin-top: 24px; padding-top: 16px; border-top: 1px solid #e5e7eb;">\n<p style="color: #6b7280; font-size: 14px;">\n<a href="https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Facebook</a> |\n<a href="https://www.instagram.com/koninklijkeloop/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Instagram</a> |\n<a href="https://dekoninklijkeloop.nl" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Website</a>\n</p>\n<p style="color: #9ca3af; font-size: 12px; margin-top: 8px;">&copy; 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n</div>	2025-10-07 15:29:06.318357+00	newsletter_manual_b2144362-495e-43f5-8761-f5378645fc86	2025-10-07 15:28:55.783019+00	2025-10-07 15:29:06.318609+00
21736eac-dc9b-4160-9fb7-330df12ae528	Nieuwsbrief DKL25	<div style="text-align: center; margin-bottom: 24px;">\n<img src="https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png" alt="De Koninklijke Loop" style="max-width: 200px; width: 100%; height: auto; margin-bottom: 16px;">\n</div>\n\n<h2 style="color: #ff9328; margin-bottom: 16px;">Welkom bij de nieuwsbrief van De Koninklijke Loop 2025</h2>\n\n<p>Beste leden en geïnteresseerden,</p>\n\n<p>Hier vindt u het laatste nieuws over De Koninklijke Loop 2025.</p>\n\n<div style="background-color: #fff7ed; border: 1px solid #ffedd5; border-radius: 8px; padding: 16px; margin: 20px 0; color: #9a3412;">\n<strong>Wat kunt u verwachten:</strong><br>\n• Actuele informatie over het evenement<br>\n• Inschrijvingsupdates<br>\n• Route-informatie<br>\n• Sponsormogelijkheden\n</div>\n\n<p>Blijf op de hoogte van alle ontwikkelingen!</p>\n\n<p>Met sportieve groet,<br><strong>Team De Koninklijke Loop</strong></p>\n\n<div style="text-align: center; margin-top: 24px; padding-top: 16px; border-top: 1px solid #e5e7eb;">\n<p style="color: #6b7280; font-size: 14px;">\n<a href="https://www.facebook.com/p/De-Koninklijke-Loop-61556315443279/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Facebook</a> |\n<a href="https://www.instagram.com/koninklijkeloop/" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Instagram</a> |\n<a href="https://dekoninklijkeloop.nl" style="color: #ff9328; text-decoration: none; margin: 0 8px;">Website</a>\n</p>\n<p style="color: #9ca3af; font-size: 12px; margin-top: 8px;">&copy; 2025 De Koninklijke Loop. Alle rechten voorbehouden.</p>\n</div>	2025-10-07 15:38:57.003813+00	newsletter_manual_21736eac-dc9b-4160-9fb7-330df12ae528	2025-10-07 15:38:48.226746+00	2025-10-07 15:38:57.004162+00
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.notifications (id, type, priority, title, message, sent, sent_at, created_at, updated_at) FROM stdin;
34e5e110-ec72-4be3-ab47-89fc77f5468b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6f86df8f4-gpgwl in omgeving unknown.	f	\N	2025-03-18 12:58:20.619769+00	2025-03-18 12:58:20.619769+00
d9653f9f-1513-43e5-b298-dfe5817f0b49	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7df587d7f7-nrq2g in omgeving unknown.	f	\N	2025-03-18 13:13:05.419368+00	2025-03-18 13:13:05.419368+00
1576ec87-8d30-4408-90f1-da247c5fd5fc	system	medium	Test Notificatie	Dit is een testnotificatie vanuit het API test script.	t	2025-03-18 13:15:04.2558+00	2025-03-18 13:15:03.640272+00	2025-03-18 13:15:04.25765+00
e91a71b8-07b9-4866-91b0-f22df76b8cab	contact	medium	Nieuw Contactverzoek	<b>test</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\ntestvoortelegram	t	2025-03-18 13:16:15.256951+00	2025-03-18 13:16:14.860279+00	2025-03-18 13:16:15.25801+00
60e3db10-5a28-4d97-baf5-a9a12f1e0f76	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-b76867b9c-kcvdn in omgeving unknown.	f	\N	2025-03-18 15:41:01.7162+00	2025-03-18 15:41:01.7162+00
57f57b40-e39d-4007-9150-2417a33f179b	aanmelding	medium	Nieuwe Aanmelding	<b>jefftest</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-03-18 15:41:34.416292+00	2025-03-18 15:41:33.812526+00	2025-03-18 15:41:34.416819+00
946cc65e-b3a2-4835-aba3-026ebed8ec42	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-54f57dcf86-jhp52 in omgeving unknown.	f	\N	2025-03-18 16:26:15.914681+00	2025-03-18 16:26:15.914681+00
2d217375-d81e-4064-87d6-e41c8f65d113	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7cb77f686d-fdrvv in omgeving unknown.	f	\N	2025-03-18 18:51:49.319493+00	2025-03-18 18:51:49.319493+00
b412267e-4361-449d-9f3f-fb3084e04680	system	medium	Test Notificatie	Dit is een testnotificatie vanuit het API test script.	t	2025-03-18 19:04:12.609047+00	2025-03-18 19:04:12.072222+00	2025-03-18 19:04:12.609685+00
13c41f40-3766-49e5-9526-66f976728dee	system	medium	Test Notificatie	Dit is een testnotificatie vanuit het API test script.	t	2025-03-18 20:54:32.879951+00	2025-03-18 20:54:32.217675+00	2025-03-18 20:54:32.880511+00
f73fcf28-8f7d-40ca-a12e-9d04a5777a2d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-75b4bf6bf4-s5dws in omgeving unknown.	f	\N	2025-03-18 22:04:09.930125+00	2025-03-18 22:04:09.930125+00
2a71755c-7b69-471c-b284-d10aa182c81a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c8f979768-qpdf7 in omgeving unknown.	f	\N	2025-03-18 22:52:29.71475+00	2025-03-18 22:52:29.71475+00
459f29d7-af63-45db-a35c-d772cb97e8e4	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-75688c7f5b-hlbkf in omgeving unknown.	f	\N	2025-03-19 00:37:58.325305+00	2025-03-19 00:37:58.325305+00
726cbeb8-c655-456e-8e15-b03e1a01fc44	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c59c7cb-kb8gz in omgeving unknown.	f	\N	2025-03-19 00:50:03.813384+00	2025-03-19 00:50:03.813384+00
4062a88d-0924-4eed-b75b-02ba4a01ad7d	aanmelding	medium	Nieuwe Aanmelding	<b>TEST*</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n<b>Bijzonderheden:</b>\ntest	t	2025-03-20 16:10:59.136266+00	2025-03-20 16:10:58.469588+00	2025-03-20 16:10:59.137652+00
52818077-70d6-4047-be48-fafe1338f56c	contact	medium	Nieuw Contactverzoek	<b>test</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nLaatste test bericht. via officiele route	t	2025-03-20 16:11:27.482037+00	2025-03-20 16:11:27.222769+00	2025-03-20 16:11:27.482607+00
4ffed255-cfb0-458b-b628-aaa366f25ed7	aanmelding	medium	Nieuwe Aanmelding	<b>Bas heijenk </b> heeft zich aangemeld.\n\n<b>Email:</b> basheijenk96@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-03-22 16:43:10.021458+00	2025-03-22 16:43:09.434134+00	2025-03-22 16:43:10.022624+00
e4c03960-d8d8-4173-abb9-dca850616713	contact	medium	Nieuw Contactverzoek	<b>Bas heijenk </b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> basheijenk96@gmail.com\n\n<b>Bericht:</b>\nHallo ik heb me op gegeven maar ik kan helaas niet sorry 	t	2025-03-23 07:15:27.597333+00	2025-03-23 07:15:27.063018+00	2025-03-23 07:15:27.598429+00
bd5c54d8-febe-42a3-8e3f-9dbc1b097fb2	aanmelding	medium	Nieuwe Aanmelding	<b>Manuela van Zwam</b> heeft zich aangemeld.\n\n<b>Email:</b> rik.van-harxen@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n<b>Bijzonderheden:</b>\nVaste begeleider die meeloopt	t	2025-03-23 14:09:00.10471+00	2025-03-23 14:08:59.482077+00	2025-03-23 14:09:00.106329+00
6507c42b-ca4c-4f39-b6ec-fb2bde4be9c7	aanmelding	medium	Nieuwe Aanmelding	<b>jeffrey</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n<b>Bijzonderheden:</b>\nTelegram test	t	2025-03-23 15:58:21.157452+00	2025-03-23 15:58:20.540789+00	2025-03-23 15:58:21.158569+00
176878a6-5a6f-49eb-9206-02d568765be1	aanmelding	medium	Nieuwe Aanmelding	<b>Telegram</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-03-23 16:07:26.553854+00	2025-03-23 16:07:25.953365+00	2025-03-23 16:07:26.554405+00
e123c4a3-757a-45f2-bed2-11db8d5c6277	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-76969b7bdc-7ctwn in omgeving unknown.	f	\N	2025-03-23 16:38:37.12226+00	2025-03-23 16:38:37.12226+00
235d8fb3-5ba3-4405-9698-264204ae76e4	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5fdd4b7794-6jjb6 in omgeving unknown.	f	\N	2025-03-23 16:42:47.11464+00	2025-03-23 16:42:47.11464+00
a73152e7-e86f-4603-9e7d-03c2b7e53249	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-778f74fc7-ktv9k in omgeving unknown.	f	\N	2025-03-23 17:04:29.220807+00	2025-03-23 17:04:29.220807+00
262ee31d-3535-4f98-8787-15741d5e1ce5	system	medium	Test Notificatie	Dit is een testnotificatie vanuit het API test script.	t	2025-03-23 17:05:21.426644+00	2025-03-23 17:05:20.888168+00	2025-03-23 17:05:21.427272+00
6233f4ab-c49d-4d03-97c6-0e5fb0f87c1d	aanmelding	medium	Nieuwe Aanmelding	<b>TGTest</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Begeleider\n<b>Afstand:</b> 15 KM\n<b>Telefoon:</b> 06123456789\n\n<b>Bijzonderheden:</b>\nTelegram Test bericht - officiele weg	t	2025-03-23 17:06:32.335936+00	2025-03-23 17:06:32.107197+00	2025-03-23 17:06:32.336494+00
56861229-bee1-4cfa-9e54-4fe25b06a034	contact	medium	Nieuw Contactverzoek	<b>test</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nTelegram Officiele weg	t	2025-03-23 17:06:54.689137+00	2025-03-23 17:06:54.487177+00	2025-03-23 17:06:54.68971+00
d4b09f12-4716-48b3-8969-5abb90596c2e	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-b9bfb84d9-vb2kt in omgeving unknown.	f	\N	2025-03-23 17:18:19.823013+00	2025-03-23 17:18:19.823013+00
8570995f-14b9-4085-9654-594031c01a08	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c6d8c5c6c-tdkff in omgeving unknown.	f	\N	2025-03-23 17:24:58.62311+00	2025-03-23 17:24:58.62311+00
6f7dcb10-ab8b-4141-a827-7a0d431fc092	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7b57f44fcf-plpt2 in omgeving unknown.	f	\N	2025-03-23 17:41:37.914053+00	2025-03-23 17:41:37.914053+00
8b3ea5b3-46cf-47ad-93d0-d906889d83b7	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79bd7fbb9f-8ww7l in omgeving unknown.	f	\N	2025-03-23 17:49:56.819609+00	2025-03-23 17:49:56.819609+00
7c92faa6-c22f-4151-b940-5472312521a6	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-bd444899f-km4cw in omgeving unknown.	f	\N	2025-03-23 17:57:14.922456+00	2025-03-23 17:57:14.922456+00
6081e21c-d548-46e3-a307-3ffa87eb922d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7d664b54b8-vl5z9 in omgeving unknown.	f	\N	2025-03-23 18:19:14.113221+00	2025-03-23 18:19:14.113221+00
8133fd37-6520-40bc-985f-1ed59ec8cdfa	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-566d9fc454-k9glf in omgeving unknown.	f	\N	2025-03-23 18:25:26.921177+00	2025-03-23 18:25:26.921177+00
0b3010a5-7e15-497d-b3d1-452e36d7d654	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-c9d9c6856-j4q6r in omgeving unknown.	f	\N	2025-03-23 18:27:19.92034+00	2025-03-23 18:27:19.92034+00
bf69ef25-8fc4-421f-a948-8c376aea6311	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7999754c4d-h9kqr in omgeving unknown.	f	\N	2025-03-23 18:57:09.830767+00	2025-03-23 18:57:09.830767+00
2d739314-f700-4d03-aa2e-6dcff3351a04	aanmelding	medium	Nieuwe Aanmelding	<b>Annerieke Mandemaker-Timmer</b> heeft zich aangemeld.\n\n<b>Email:</b> annerieketimmer@hotmail.com\n\n<b>Rol:</b> Begeleider\n<b>Afstand:</b> 6 KM\n<b>Telefoon:</b> 06 17 37 28 40 \n\n	t	2025-03-24 09:16:48.570098+00	2025-03-24 09:16:48.343152+00	2025-03-24 09:16:48.571314+00
c8496a82-5301-442b-a806-0a0334854dbd	aanmelding	medium	Nieuwe Aanmelding	<b>Pieter Streefland</b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-07 17:50:25.669019+00	2025-05-07 17:50:25.462634+00	2025-05-07 17:50:25.669713+00
501fb888-4af6-4ef6-9455-e3a72ac86ce8	aanmelding	medium	Nieuwe Aanmelding	<b>Jean-paul Hup</b> heeft zich aangemeld.\n\n<b>Email:</b> molenkamp.19@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-03-24 09:18:06.804946+00	2025-03-24 09:18:06.502832+00	2025-03-24 09:18:06.805601+00
33871660-b6cb-42a9-959f-d564ce71630f	aanmelding	medium	Nieuwe Aanmelding	<b>Arno Kerkvliet</b> heeft zich aangemeld.\n\n<b>Email:</b> arno.kerkvliet@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-03-24 17:26:18.868101+00	2025-03-24 17:26:18.628209+00	2025-03-24 17:26:18.868645+00
22d60281-09cc-45dc-9a8c-83f850f93ee6	aanmelding	medium	Nieuwe Aanmelding	<b>Mirjam Kerkvliet</b> heeft zich aangemeld.\n\n<b>Email:</b> mirjam.kerkvliet@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-03-24 17:27:30.782205+00	2025-03-24 17:27:30.561496+00	2025-03-24 17:27:30.783262+00
1ef54881-1960-44f7-a326-9206f8e4bed0	aanmelding	medium	Nieuwe Aanmelding	<b>Karin de Jong</b> heeft zich aangemeld.\n\n<b>Email:</b> karin.de.jong82@outlook.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-03-24 18:47:11.362691+00	2025-03-24 18:47:11.151679+00	2025-03-24 18:47:11.364287+00
7f3701fb-a73a-4526-9057-3c89f01c5d73	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-c8dd5c95f-hd6q5 in omgeving unknown.	f	\N	2025-03-24 20:27:19.121317+00	2025-03-24 20:27:19.121317+00
449852eb-52a3-4772-8a67-6857a5db4d61	aanmelding	medium	Nieuwe Aanmelding	<b>A. Bistolfi</b> heeft zich aangemeld.\n\n<b>Email:</b> nedarg@icloud.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-03-26 12:27:27.447178+00	2025-03-26 12:27:27.17636+00	2025-03-26 12:27:27.448372+00
5ed39d7e-2a97-4d25-b205-c3e7d8052007	aanmelding	medium	Nieuwe Aanmelding	<b>Ayla Toprak</b> heeft zich aangemeld.\n\n<b>Email:</b> gamergirlayla@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 10 KM\n	t	2025-03-26 16:25:50.392164+00	2025-03-26 16:25:50.141281+00	2025-03-26 16:25:50.392665+00
460c2939-0560-4cf1-a84b-2b714d1bc258	aanmelding	medium	Nieuwe Aanmelding	<b>Mila Veenendaal</b> heeft zich aangemeld.\n\n<b>Email:</b> gaminggirlayla@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 10 KM\n	t	2025-03-26 16:26:59.286978+00	2025-03-26 16:26:59.053071+00	2025-03-26 16:26:59.287547+00
ded8489e-e699-4f7e-b7a1-ebf25ed0a551	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7cf4f44579-rz4t9 in omgeving unknown.	f	\N	2025-03-28 15:41:37.022145+00	2025-03-28 15:41:37.022145+00
06a0b656-2ed1-43ca-9c48-fbde84f830f6	contact	medium	Nieuw Contactverzoek	<b>Test User from Whisky For Charity</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nThis is a test email from Whisky For Charity via DKL Email Service. Testing integration with your service.	t	2025-03-28 15:49:16.996319+00	2025-03-28 15:49:16.738573+00	2025-03-28 15:49:16.997017+00
06df0313-d04f-4e49-9ac2-a00b953b1930	contact	medium	Nieuw Contactverzoek	<b>Test User from Whisky For Charity</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nThis is a test email from Whisky For Charity via DKL Email Service. Testing integration with your service.	t	2025-03-28 16:08:27.673923+00	2025-03-28 16:08:27.432678+00	2025-03-28 16:08:27.674701+00
2700db07-6a56-49f0-aaac-916fa3da291d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7445b5b859-kq84t in omgeving unknown.	f	\N	2025-03-28 16:14:24.324091+00	2025-03-28 16:14:24.324091+00
ddc7ab4b-6be1-481b-8369-f1a36ae51659	contact	medium	Nieuw Contactverzoek	<b>Test User from Whisky For Charity</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nThis is a test email from Whisky For Charity via DKL Email Service. Testing integration with your service.	t	2025-03-28 16:22:56.737232+00	2025-03-28 16:22:56.508798+00	2025-03-28 16:22:56.737877+00
084ca7eb-0e17-478c-80c6-d062fea3b997	contact	medium	Nieuw Contactverzoek	<b>Test User from Whisky For Charity</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nThis is a test email from Whisky For Charity via DKL Email Service. Testing integration with your service.	t	2025-03-28 16:38:24.512777+00	2025-03-28 16:31:02.852773+00	2025-03-28 16:38:24.513357+00
e127aeeb-db7c-48ff-870d-d5a5d21ac2c2	contact	medium	Nieuw Contactverzoek	<b>Test User from Whisky For Charity</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nThis is a test email from Whisky For Charity via DKL Email Service. Testing integration with your service.	t	2025-03-28 16:53:24.535157+00	2025-03-28 16:36:55.284376+00	2025-03-28 16:53:24.535755+00
0fc26c68-49c0-4da2-b238-7dbad42b014d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7489c4746f-w7dzp in omgeving unknown.	f	\N	2025-03-28 18:58:38.014201+00	2025-03-28 18:58:38.014201+00
45d7b5dc-b401-44c3-931f-8e3a668a5f5a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-564c7dccc5-6ssw7 in omgeving unknown.	f	\N	2025-03-28 19:54:49.322193+00	2025-03-28 19:54:49.322193+00
a554c424-d2da-4ee9-8dd2-d7e25855fc64	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-84fdd474f9-r62xr in omgeving unknown.	f	\N	2025-03-28 20:02:54.823238+00	2025-03-28 20:02:54.823238+00
b27771da-b714-42e0-bf5e-fd5c2e3f4c1d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-ff5cdcd84-rts8z in omgeving unknown.	f	\N	2025-03-28 20:18:53.608365+00	2025-03-28 20:18:53.608365+00
7d556d41-6380-4422-b197-146890c37146	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-56d9b6558-942bp in omgeving unknown.	f	\N	2025-03-28 20:28:50.915286+00	2025-03-28 20:28:50.915286+00
f489b9ce-793d-4262-8001-2c008620c5c8	aanmelding	medium	Nieuwe Aanmelding	<b>jeffrey</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 10 KM\n	t	2025-03-28 20:43:40.676004+00	2025-03-28 20:43:40.454931+00	2025-03-28 20:43:40.676705+00
c924724f-2bc3-4529-b284-e77a4159c7d4	contact	medium	Nieuw Contactverzoek	<b>test</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nTEST hybride systeem	t	2025-03-28 20:44:13.953122+00	2025-03-28 20:44:13.73327+00	2025-03-28 20:44:13.953802+00
4de770c1-fe3d-42f3-b1a0-5c3c6bbe890c	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-598c977d5f-424w2 in omgeving unknown.	f	\N	2025-03-28 21:00:02.917783+00	2025-03-28 21:00:02.917783+00
72b6f888-c649-4ba4-b928-a364131daf35	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-f7bd67956-jbsls in omgeving unknown.	f	\N	2025-03-28 21:10:21.921066+00	2025-03-28 21:10:21.921066+00
6e7af7e2-3b52-40a8-8c36-50737ce24cbd	aanmelding	medium	Nieuwe Aanmelding	<b>Klaske van de glind</b> heeft zich aangemeld.\n\n<b>Email:</b> Klaskehiddes@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-03-29 10:05:10.739644+00	2025-03-29 10:05:10.510755+00	2025-03-29 10:05:10.740267+00
4003ae74-eaa7-4867-8aa5-759be1cb0174	aanmelding	medium	Nieuwe Aanmelding	<b>Bertram tijsma</b> heeft zich aangemeld.\n\n<b>Email:</b> Klaskehiddes@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-03-29 10:07:17.829437+00	2025-03-29 10:07:17.625243+00	2025-03-29 10:07:17.830548+00
e6edfca0-d728-4263-8aad-c44ed013d722	aanmelding	medium	Nieuwe Aanmelding	<b>Han van Doornik</b> heeft zich aangemeld.\n\n<b>Email:</b> LaanvanGS.26@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n<b>Bijzonderheden:</b>\nIk wil wel graag begeleiding 	t	2025-03-30 08:07:52.183647+00	2025-03-30 08:07:51.959973+00	2025-03-30 08:07:52.184165+00
981093f6-aba8-489e-af77-09fdfe0faf4f	contact	medium	Nieuw Contactverzoek	<b>Mirjam Rijswijk</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> mirjamrijswijk@icloud.com\n\n<b>Bericht:</b>\nIk wil graag helpen met dit evenement. Ik ben evenementen eerste hulpverlener bij het Rode Kruis en daarnaast wandelbegeleider bij de wandeluitdaging Apeldoorn. Ik hoor graag van jullie . Hartelijke groet, Mirjam Rijswijk 	t	2025-03-30 09:07:02.463054+00	2025-03-30 09:07:02.234551+00	2025-03-30 09:07:02.463605+00
9b7f780a-54cd-45c1-8039-175605f2e04b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6488fc4cc8-bxgbg in omgeving unknown.	f	\N	2025-03-30 11:06:33.721439+00	2025-03-30 11:06:33.721439+00
d4dae691-488a-4abe-ab99-d69a22cbd8c5	aanmelding	medium	Nieuwe Aanmelding	<b>Anneke van de Glind </b> heeft zich aangemeld.\n\n<b>Email:</b> Klaskehiddes@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-03-31 11:56:44.006021+00	2025-03-31 11:56:43.7781+00	2025-03-31 11:56:44.006531+00
0fbec407-11cc-413d-8d2f-48a2d061950a	aanmelding	medium	Nieuwe Aanmelding	<b>Noa hiddes</b> heeft zich aangemeld.\n\n<b>Email:</b> Klaskehiddes@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-03-31 11:57:29.663624+00	2025-03-31 11:57:29.466009+00	2025-03-31 11:57:29.664164+00
d6856ce1-35ba-43f5-84eb-cbccd16d4763	aanmelding	medium	Nieuwe Aanmelding	<b>Sylvia Dijkstra</b> heeft zich aangemeld.\n\n<b>Email:</b> sylvia.dijkstra@sheerenloo.nl\n\n<b>Rol:</b> Begeleider\n<b>Afstand:</b> 2.5 KM\n<b>Telefoon:</b> 0683081728 \n\n	t	2025-04-06 18:48:21.960115+00	2025-04-06 18:48:21.708765+00	2025-04-06 18:48:21.960642+00
591b7ed9-ab7b-412f-b738-b716aea34be2	aanmelding	medium	Nieuwe Aanmelding	<b>Diesmer </b> heeft zich aangemeld.\n\n<b>Email:</b> diesbosje@hotmail.com\n\n<b>Rol:</b> Begeleider\n<b>Afstand:</b> 6 KM\n<b>Telefoon:</b> 0613429612\n\n	t	2025-04-12 07:07:08.274327+00	2025-04-12 07:07:08.055343+00	2025-04-12 07:07:08.275514+00
71ebe4f9-6d4f-47ca-9e4b-0a0653012cf7	aanmelding	medium	Nieuwe Aanmelding	<b>Albert </b> heeft zich aangemeld.\n\n<b>Email:</b> diesbosje@hotmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-12 07:08:42.322707+00	2025-04-12 07:08:42.126194+00	2025-04-12 07:08:42.323162+00
12dd1a50-f167-4639-919e-eab484da7010	aanmelding	medium	Nieuwe Aanmelding	<b>Theun </b> heeft zich aangemeld.\n\n<b>Email:</b> diesbosje@hotmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-12 07:09:26.206901+00	2025-04-12 07:09:26.019694+00	2025-04-12 07:09:26.20747+00
93c67a31-13fb-4773-a34b-10bd4f125458	aanmelding	medium	Nieuwe Aanmelding	<b>Hilde Rekers </b> heeft zich aangemeld.\n\n<b>Email:</b> h.rekers1959@kpnmail.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-12 10:20:32.584396+00	2025-04-12 10:20:32.286569+00	2025-04-12 10:20:32.584934+00
add126f2-d069-4b2b-a9db-a2335394478b	aanmelding	medium	Nieuwe Aanmelding	<b>Henk Rekers </b> heeft zich aangemeld.\n\n<b>Email:</b> h.rekers1959@kpnmail.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-12 10:21:43.425462+00	2025-04-12 10:21:43.146833+00	2025-04-12 10:21:43.425936+00
70440d88-3971-4872-8c15-080c508a37c9	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-648c6d85cc-w4jh6 in omgeving unknown.	f	\N	2025-04-12 19:56:41.118096+00	2025-04-12 19:56:41.118096+00
d3bfb39b-ec6e-4a67-b137-35d8b7901c9a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-596c97b984-wr9sj in omgeving unknown.	f	\N	2025-04-12 20:02:20.423579+00	2025-04-12 20:02:20.423579+00
6c19f243-75be-4078-a1e8-c9a0766f4c75	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79668b887d-8n2sh in omgeving unknown.	f	\N	2025-04-12 20:16:45.520234+00	2025-04-12 20:16:45.520234+00
d41eb4a6-410a-4139-acba-bc54fd98bae6	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-64f7cb68c5-gkqws in omgeving unknown.	f	\N	2025-04-12 22:44:15.315746+00	2025-04-12 22:44:15.315746+00
7641ca10-cc5f-477e-b2ba-4a0cafb93910	aanmelding	medium	Nieuwe Aanmelding	<b>Hilde Rekers </b> heeft zich aangemeld.\n\n<b>Email:</b> h.rekers59@kpnmail.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-14 20:00:53.682779+00	2025-04-14 20:00:53.445643+00	2025-04-14 20:00:53.683416+00
90c783d7-3cac-4995-8de4-56720ce2b84c	aanmelding	medium	Nieuwe Aanmelding	<b>Henk Rekers </b> heeft zich aangemeld.\n\n<b>Email:</b> h.rekers59@kpnmail.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-14 20:03:52.936231+00	2025-04-14 20:03:52.718349+00	2025-04-14 20:03:52.940936+00
8ee2f8f6-fdd0-47b6-9ccc-3ef7c7fcf8be	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-76f88dcf6-sfkv7 in omgeving unknown.	f	\N	2025-04-16 15:18:23.919452+00	2025-04-16 15:18:23.919452+00
0077f557-5b0e-4d57-8489-d86e6ffc6027	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-cb59747c6-kf5tl in omgeving unknown.	f	\N	2025-04-16 16:00:48.819124+00	2025-04-16 16:00:48.819124+00
3b33fcc7-a6e9-4954-80d2-87b62d9085b8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c7467d96-7c5pv in omgeving unknown.	f	\N	2025-04-16 16:28:56.512053+00	2025-04-16 16:28:56.512053+00
d58ef111-b36d-4dbd-adf7-f7b483c6d6ed	aanmelding	medium	Nieuwe Aanmelding	<b>Bram Aarnoudse</b> heeft zich aangemeld.\n\n<b>Email:</b> bram.aarnoudse@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-04-17 09:30:00.788975+00	2025-04-17 09:30:00.55972+00	2025-04-17 09:30:00.790278+00
92b4e39f-afed-4eba-820d-172d3da71618	aanmelding	medium	Nieuwe Aanmelding	<b>jeff</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-04-17 11:20:34.64949+00	2025-04-17 11:20:34.410597+00	2025-04-17 11:20:34.651233+00
b7d15fe4-97e7-4b79-90fb-b312eacdda6b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-686d4c87d6-nxgjw in omgeving unknown.	f	\N	2025-04-17 12:06:51.221192+00	2025-04-17 12:06:51.221192+00
03db4995-2263-4063-9b06-3a1b3963287d	aanmelding	medium	Nieuwe Aanmelding	<b>TEST</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-04-17 12:30:53.565191+00	2025-04-17 12:30:53.318018+00	2025-04-17 12:30:53.565975+00
1981a16c-3b26-43b1-9158-8da75d386c7b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7b8cf68554-s47jk in omgeving unknown.	f	\N	2025-04-17 12:42:38.531762+00	2025-04-17 12:42:38.531762+00
331e9934-9ddb-4a9c-9879-99205b7bfd10	aanmelding	medium	Nieuwe Aanmelding	<b>jeffreyTEST</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-04-17 12:43:24.696357+00	2025-04-17 12:43:24.500542+00	2025-04-17 12:43:24.697059+00
9a1542b7-77d4-4ef1-b658-1300cc0a85a8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-c96968967-4q7v6 in omgeving unknown.	f	\N	2025-04-18 14:20:24.920119+00	2025-04-18 14:20:24.920119+00
f84f5d28-b01b-4d73-a513-a62fe404853a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-9f9f9775c-l2fbl in omgeving unknown.	f	\N	2025-04-18 14:41:53.618073+00	2025-04-18 14:41:53.618073+00
29e63e30-9bae-4705-9798-9f1a87a51ba5	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-66b4b44949-nnv44 in omgeving unknown.	f	\N	2025-04-18 18:58:57.021484+00	2025-04-18 18:58:57.021484+00
0c531ce0-5ecf-4823-8dd9-17eb38345176	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-78cf97869b-589g6 in omgeving unknown.	f	\N	2025-04-18 19:16:33.315552+00	2025-04-18 19:16:33.315552+00
4f9fe98a-fc52-476f-ba9b-746d0048f91d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-bb49bc48d-8dvjq in omgeving unknown.	f	\N	2025-04-18 19:30:29.511593+00	2025-04-18 19:30:29.511593+00
64cc9022-5056-4109-868d-132e344d741b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-77bf495b67-f7k7r in omgeving unknown.	f	\N	2025-04-18 19:46:12.216392+00	2025-04-18 19:46:12.216392+00
ef0cb0a1-bcf3-4840-862d-51183364c39f	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6cc48f47cc-ss85d in omgeving unknown.	f	\N	2025-04-18 19:53:48.515885+00	2025-04-18 19:53:48.515885+00
39518bce-ff86-47f6-91ab-3a8e2f744b85	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-68976bc86-pzlcf in omgeving unknown.	f	\N	2025-04-18 19:59:02.118797+00	2025-04-18 19:59:02.118797+00
ae0a6b74-f90a-4070-9d24-2070a749b445	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-9797cb8b9-mdsr6 in omgeving unknown.	f	\N	2025-04-19 20:32:13.123737+00	2025-04-19 20:32:13.123737+00
6a9245cb-bbd5-46cc-9b05-7d49739a35fd	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6756f76b5f-4gwlm in omgeving unknown.	f	\N	2025-04-19 22:23:37.513327+00	2025-04-19 22:23:37.513327+00
fbedb8da-0007-4543-b815-d38c83f0e04b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7cd4f8f78f-4cr49 in omgeving unknown.	f	\N	2025-04-19 22:33:57.214796+00	2025-04-19 22:33:57.214796+00
9098438e-af41-491a-b3d6-2c453a13766a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-64dcc8747b-zt22m in omgeving unknown.	f	\N	2025-04-19 23:50:41.517187+00	2025-04-19 23:50:41.517187+00
31fba7f7-d199-495e-b46b-ea40cb984e81	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6d48d4c9cc-b88gm in omgeving unknown.	f	\N	2025-04-20 00:02:57.723733+00	2025-04-20 00:02:57.723733+00
6736a193-4085-457f-987f-7482146d5f0b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-b7c4699bf-qbm6v in omgeving unknown.	f	\N	2025-04-20 13:02:11.720155+00	2025-04-20 13:02:11.720155+00
227af31e-dc78-422f-b609-acded58ce292	aanmelding	medium	Nieuwe Aanmelding	<b>Sophie Hubers</b> heeft zich aangemeld.\n\n<b>Email:</b> sophieehubers@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-04-29 09:28:14.682526+00	2025-04-29 09:28:14.452575+00	2025-04-29 09:28:14.683314+00
d273ab63-a6a1-4333-9b1c-fc22a0bde99d	aanmelding	medium	Nieuwe Aanmelding	<b>Geer Hubers</b> heeft zich aangemeld.\n\n<b>Email:</b> hub3008@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-04-29 09:28:58.42277+00	2025-04-29 09:28:58.187731+00	2025-04-29 09:28:58.423458+00
38939be1-5400-4297-a9ef-daa4e849d8de	aanmelding	medium	Nieuwe Aanmelding	<b>yassine</b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-07 18:53:51.116871+00	2025-05-07 18:53:50.905345+00	2025-05-07 18:53:51.117551+00
7cda9262-2469-4fe2-9489-eed74ac02ea2	aanmelding	medium	Nieuwe Aanmelding	<b>Saleem</b> heeft zich aangemeld.\n\n<b>Email:</b> lidaahmadi99@gmail.com\n\n<b>Rol:</b> Begeleider\n<b>Afstand:</b> 6 KM\n<b>Telefoon:</b> 0649020648\n\n	t	2025-04-29 18:55:29.068586+00	2025-04-29 18:55:28.873163+00	2025-04-29 18:55:29.069358+00
755361a8-2686-4c3b-bf1d-6d8d166c444f	aanmelding	medium	Nieuwe Aanmelding	<b>jeffrey</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-05-18 11:31:05.787691+00	2025-05-18 11:31:05.474121+00	2025-05-18 11:31:05.788444+00
9b93488e-bea7-457a-82ee-e99402e9d63c	aanmelding	medium	Nieuwe Aanmelding	<b>Yussef</b> heeft zich aangemeld.\n\n<b>Email:</b> lidaahmadi99@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-29 19:00:02.773013+00	2025-04-29 19:00:02.576289+00	2025-04-29 19:00:02.773725+00
e36aee99-904f-4644-a1f3-14782e745426	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-66749fc578-s9gkb in omgeving unknown.	f	\N	2025-06-23 05:22:38.172283+00	2025-06-23 05:22:38.172283+00
b0cda113-b8e6-4314-ae9d-562d9180afcb	aanmelding	medium	Nieuwe Aanmelding	<b>Sulleeman </b> heeft zich aangemeld.\n\n<b>Email:</b> lidaahmadi99@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-29 19:00:45.000933+00	2025-04-29 19:00:44.732378+00	2025-04-29 19:00:45.00161+00
cc40fd05-db88-4cab-b614-94814d05c0a1	aanmelding	medium	Nieuwe Aanmelding	<b>Yunis </b> heeft zich aangemeld.\n\n<b>Email:</b> lidaahmadi99@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-29 19:02:05.748344+00	2025-04-29 19:02:05.50951+00	2025-04-29 19:02:05.749009+00
c42711c4-8358-4d3c-8498-c0ebd6af6a9e	aanmelding	medium	Nieuwe Aanmelding	<b>JeffTest</b> heeft zich aangemeld.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 2.5 KM\n	t	2025-10-03 15:39:39.690675+00	2025-10-03 15:39:39.351353+00	2025-10-03 15:39:39.691298+00
59993c48-8a30-4302-8fd6-905d4df6765a	aanmelding	medium	Nieuwe Aanmelding	<b>Ismail </b> heeft zich aangemeld.\n\n<b>Email:</b> lidaahmadi99@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-04-29 19:02:40.712803+00	2025-04-29 19:02:40.497176+00	2025-04-29 19:02:40.71349+00
a0db5902-ae5e-4dcd-845d-4aef8a6cdbfe	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-545798dbcc-frc67 in omgeving unknown.	f	\N	2025-10-03 19:13:20.15951+00	2025-10-03 19:13:20.15951+00
51a5d1c9-1caa-49b7-a1e3-275256dbca7a	aanmelding	medium	Nieuwe Aanmelding	<b>Mayke Rood</b> heeft zich aangemeld.\n\n<b>Email:</b> rood1960@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-04-29 21:16:13.819082+00	2025-04-29 21:16:13.618692+00	2025-04-29 21:16:13.81977+00
71aa140a-39ab-4692-8205-d6e2dc2a33a9	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-85655cb77d-bgdd4 in omgeving unknown.	f	\N	2025-10-04 17:05:57.359983+00	2025-10-04 17:05:57.359983+00
e9df2ab8-df8f-4ca8-b65a-ad6a820aa4a0	aanmelding	medium	Nieuwe Aanmelding	<b>laura</b> heeft zich aangemeld.\n\n<b>Email:</b> steun94@xs4all.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 10 KM\n	t	2025-04-30 10:03:30.807248+00	2025-04-30 10:03:30.583613+00	2025-04-30 10:03:30.808061+00
099fef2a-6a55-4363-b25f-784a62c2a687	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7fc68d49ff-clbf4 in omgeving unknown.	f	\N	2025-10-05 12:56:02.063713+00	2025-10-05 12:56:02.063713+00
705a09bf-4f8e-4023-8fea-bcf50e44bec3	aanmelding	medium	Nieuwe Aanmelding	<b>Luciënne </b> heeft zich aangemeld.\n\n<b>Email:</b> luciennetang@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 15 KM\n	t	2025-04-30 20:55:13.912084+00	2025-04-30 20:55:13.561101+00	2025-04-30 20:55:13.912745+00
039b49a0-d18c-4d4c-b903-6d2af896f8b7	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-fdddb49c9-btgsb in omgeving unknown.	f	\N	2025-10-05 13:11:16.361641+00	2025-10-05 13:11:16.361641+00
16ce6042-a150-4e73-8b44-c5eab417a8e5	aanmelding	medium	Nieuwe Aanmelding	<b>Theodora Naus</b> heeft zich aangemeld.\n\n<b>Email:</b> tgemooi@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-01 10:14:33.97848+00	2025-05-01 10:14:33.70683+00	2025-05-01 10:14:33.979151+00
ddf46521-f564-45cc-a401-310b65c0b324	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7f9494895c-cfnsj in omgeving unknown.	f	\N	2025-10-05 13:24:21.860648+00	2025-10-05 13:24:21.860648+00
52d093dd-c648-4d7c-8f8e-618d422b0cb0	aanmelding	medium	Nieuwe Aanmelding	<b>Theodora Naus</b> heeft zich aangemeld.\n\n<b>Email:</b> tgemooi@gmail.com\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-01 10:30:11.946256+00	2025-05-01 10:23:13.633343+00	2025-05-01 10:30:11.946953+00
ae6d8129-c778-40dc-9323-ec8caf573371	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-649b558b9f-z22vl in omgeving unknown.	f	\N	2025-10-05 13:41:43.458885+00	2025-10-05 13:41:43.458885+00
fb67495d-948e-4d13-bf1a-153595c7e0e9	contact	medium	Nieuw Contactverzoek	<b>Luciēnne Tang</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> luciennetang@gmail.com\n\n<b>Bericht:</b>\nAnnulering inschrijving.\n\nMijn inschrijfgegevens zijn:\nNaam : Luciēnne \nEmail: luciennetang@gmail.com \nRol: deelnemer\nAfstand: 15 KM	t	2025-05-04 06:46:39.307636+00	2025-05-04 06:46:39.108385+00	2025-05-04 06:46:39.308563+00
65d45a16-1752-47ae-8a83-e0a403767d49	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-769779bd55-4c2mc in omgeving unknown.	f	\N	2025-10-05 14:19:28.054147+00	2025-10-05 14:19:28.054147+00
2660ffcd-d342-4f70-a769-5f34de916efb	contact	medium	Nieuw Contactverzoek	<b>Theodora Naus</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> tgemooi@gmail.com\n\n<b>Bericht:</b>\nIk had me aangemeld voor de afstand van 6 km maar kan helaas niet meedoen. Gelukkig had ik al wel gedoneerd.\nSucces verder!\n\nVriendelijke groet,\n\nTheodora Naus	t	2025-05-05 11:34:33.751816+00	2025-05-05 11:34:33.420991+00	2025-05-05 11:34:33.752458+00
fa797eb8-0779-491f-85cd-20a73eac7769	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-9cfdc894b-spr82 in omgeving unknown.	f	\N	2025-10-05 14:33:27.352759+00	2025-10-05 14:33:27.352759+00
e77a5c53-ad35-4e28-bfce-45c36dce4775	aanmelding	medium	Nieuwe Aanmelding	<b>michel</b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-07 14:58:43.747116+00	2025-05-07 14:58:43.544072+00	2025-05-07 14:58:43.747826+00
4a4633a9-4846-4631-983f-e939ce2adf89	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-998d4c7bc-ldgnx in omgeving unknown.	f	\N	2025-10-05 15:05:53.555242+00	2025-10-05 15:05:53.555242+00
d29441f4-2cdc-4f4b-94a0-f09bb18fd1e1	aanmelding	medium	Nieuwe Aanmelding	<b>peter</b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 10 KM\n	t	2025-05-07 14:59:21.614261+00	2025-05-07 14:59:21.403694+00	2025-05-07 14:59:21.614926+00
633f7234-886c-44f5-a0c2-dc01b4ba3803	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6db9d6c4bf-zchtg in omgeving unknown.	f	\N	2025-10-05 15:35:05.368827+00	2025-10-05 15:35:05.368827+00
98351597-a2f7-4c31-935f-dd778a980441	aanmelding	medium	Nieuwe Aanmelding	<b>john</b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-07 15:00:14.560042+00	2025-05-07 15:00:14.367912+00	2025-05-07 15:00:14.560775+00
d0e06945-1ffc-403e-85d5-eb0fc05df201	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-dc5948f67-fbz65 in omgeving unknown.	f	\N	2025-10-05 15:43:17.762389+00	2025-10-05 15:43:17.762389+00
aa0b8704-4574-4e8e-94f6-14e916d69607	aanmelding	medium	Nieuwe Aanmelding	<b>Danny </b> heeft zich aangemeld.\n\n<b>Email:</b> enckerkamp.23@sheerenloo.nl\n\n<b>Rol:</b> Deelnemer\n<b>Afstand:</b> 6 KM\n	t	2025-05-07 15:00:58.223611+00	2025-05-07 15:00:57.982082+00	2025-05-07 15:00:58.224312+00
3b65c880-7b6e-44da-9ee9-70503aee7bdd	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-688df8746f-5x4v6 in omgeving unknown.	f	\N	2025-10-05 16:27:29.151098+00	2025-10-05 16:27:29.151098+00
cc375610-f788-42d1-b3b9-2b90df00a8be	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-65d646bcd-cx466 in omgeving unknown.	f	\N	2025-10-05 17:00:59.362334+00	2025-10-05 17:00:59.362334+00
249e13c5-02c8-4f3e-9053-c0a7a4aeafae	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6774bf65c9-r9msw in omgeving unknown.	f	\N	2025-10-05 17:09:26.248532+00	2025-10-05 17:09:26.248532+00
1c63f3b3-558c-4de0-bacb-a8a71141e68b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5874767f86-b582g in omgeving unknown.	f	\N	2025-10-05 17:56:36.45734+00	2025-10-05 17:56:36.45734+00
95cf16d6-e2cb-4164-b265-be72f16bd473	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-9974c5b88-b6l8p in omgeving unknown.	f	\N	2025-10-05 18:26:14.455436+00	2025-10-05 18:26:14.455436+00
e8433672-4f25-4784-a4d5-32dddcba413f	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c85c8d895-x5l55 in omgeving unknown.	f	\N	2025-10-05 18:35:40.156527+00	2025-10-05 18:35:40.156527+00
f0146c7d-68fa-4bef-af47-f24e441f2b72	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-666876855b-ssxjf in omgeving unknown.	f	\N	2025-10-05 18:44:13.563249+00	2025-10-05 18:44:13.563249+00
542ee875-6f67-4415-a1de-7bf9100537c8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6559c977b7-qcd6b in omgeving unknown.	f	\N	2025-10-05 19:04:05.35154+00	2025-10-05 19:04:05.35154+00
927e4adb-1079-42d8-8556-d1e5abfa9055	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5bf75f4d7f-22lp8 in omgeving unknown.	f	\N	2025-10-05 19:12:17.458932+00	2025-10-05 19:12:17.458932+00
ac0e1759-e63c-432d-8fbf-2691592ae8b5	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c5bb959d7-ww2lw in omgeving unknown.	f	\N	2025-10-05 19:58:21.858263+00	2025-10-05 19:58:21.858263+00
abb00e76-806b-4bf6-b9aa-f1dce052813d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5d887bc945-jnv66 in omgeving unknown.	f	\N	2025-10-07 12:15:57.964405+00	2025-10-07 12:15:57.964405+00
48cab1b3-fe51-4d76-97e9-08915cea6880	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-58dbd85bfd-r2cmn in omgeving unknown.	f	\N	2025-10-07 12:39:37.05986+00	2025-10-07 12:39:37.05986+00
affdf172-bc10-4f24-ae5b-5a0ddd42b94e	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-77bc9df446-dg6gd in omgeving unknown.	f	\N	2025-10-07 13:11:58.66167+00	2025-10-07 13:11:58.66167+00
8d0ac417-80fc-46f0-862f-92ee88b4ca5f	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7959db56b9-f9sxp in omgeving unknown.	f	\N	2025-10-07 15:03:48.849781+00	2025-10-07 15:03:48.849781+00
c826659f-c966-40be-a1ee-fe1aea3a8625	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6d98fb4c97-pgd6q in omgeving unknown.	f	\N	2025-10-07 15:15:49.45611+00	2025-10-07 15:15:49.45611+00
6f798d28-1988-40f8-8752-f2117a5a139a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5b6f64bb9f-hd5cq in omgeving unknown.	f	\N	2025-10-07 15:28:38.052489+00	2025-10-07 15:28:38.052489+00
8f419f1d-0e52-40ce-ab0c-18c0f37219ef	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5c968b8845-kgzgs in omgeving unknown.	f	\N	2025-10-07 15:36:30.953553+00	2025-10-07 15:36:30.953553+00
9299e698-e06d-40f0-a5d6-ca9bd467ecd7	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6db76f9f9f-cxfnt in omgeving unknown.	f	\N	2025-10-07 17:15:41.960658+00	2025-10-07 17:15:41.960658+00
ba0a1372-3e7e-4610-b0d1-1266d23461dd	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-fd5744c4c-wpvh4 in omgeving unknown.	f	\N	2025-10-07 17:29:13.521134+00	2025-10-07 17:29:13.521134+00
145d4a65-c5c8-4ae6-bf00-562eae467e8a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-67b97fbb6c-pl5tc in omgeving unknown.	f	\N	2025-10-07 17:38:28.733771+00	2025-10-07 17:38:28.733771+00
87ac69da-23c0-49b1-8d37-81eaa6f4b221	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-785ddffc89-z6bnk in omgeving unknown.	f	\N	2025-10-07 18:06:57.331187+00	2025-10-07 18:06:57.331187+00
e7cba290-ef76-49f3-8794-3ccffcede71c	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-65b4c9b6b6-bspc5 in omgeving unknown.	f	\N	2025-10-07 18:19:27.627535+00	2025-10-07 18:19:27.627535+00
93140f40-ddf4-4149-91db-d42a44bbab0a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-856f998ff4-xt9nl in omgeving unknown.	f	\N	2025-10-07 18:45:53.723432+00	2025-10-07 18:45:53.723432+00
58f9794a-98db-4d47-be34-c077d2c7f399	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5874c76484-85jmd in omgeving unknown.	f	\N	2025-10-07 18:51:15.733931+00	2025-10-07 18:51:15.733931+00
aa8b4247-1c0f-4fa8-90aa-0637ac588340	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5458c76bd4-kd9zd in omgeving unknown.	f	\N	2025-10-07 19:10:04.830662+00	2025-10-07 19:10:04.830662+00
bf1c2037-604c-4c5e-bd38-b9175dbd7c09	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5767dd58c8-42m66 in omgeving unknown.	f	\N	2025-10-07 19:27:58.435611+00	2025-10-07 19:27:58.435611+00
7d30ae1f-f76b-4d1c-9eef-f940b5678960	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-76c4df44c7-szkdq in omgeving unknown.	f	\N	2025-10-07 20:00:55.246067+00	2025-10-07 20:00:55.246067+00
52fa8f6f-76dc-418e-99d5-1cf2f33d094a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6fdf4b7b68-mb96d in omgeving unknown.	f	\N	2025-10-07 20:10:41.134536+00	2025-10-07 20:10:41.134536+00
235ca3ee-bc18-4ad9-bc07-2f1a72935f35	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-84f4b6457b-g75ll in omgeving unknown.	f	\N	2025-10-07 20:20:03.729981+00	2025-10-07 20:20:03.729981+00
13326620-d71b-47e6-8883-ee7984252fa5	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-58b9794659-pdd55 in omgeving unknown.	f	\N	2025-10-07 20:28:02.025496+00	2025-10-07 20:28:02.025496+00
b146be44-b932-4ba7-aaee-81f0d51c05c0	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7c79f5ff7f-n4v62 in omgeving unknown.	f	\N	2025-10-07 20:35:03.625059+00	2025-10-07 20:35:03.625059+00
c7a3476d-ca71-4d73-8efa-2e1b409357d3	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-564b44d669-tw6vv in omgeving unknown.	f	\N	2025-10-07 20:45:40.219862+00	2025-10-07 20:45:40.219862+00
16dbf341-8429-4965-8ffd-64dc5f552f64	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-84d5b97f54-v5vnp in omgeving unknown.	f	\N	2025-10-07 21:02:35.932012+00	2025-10-07 21:02:35.932012+00
7fb906ea-4937-47df-8a44-5523a8076e0b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-65d96ff647-qnb65 in omgeving unknown.	f	\N	2025-10-07 21:17:57.731358+00	2025-10-07 21:17:57.731358+00
d28eed70-f819-4ed3-a322-db814373a3d8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c578d5d45-f99xb in omgeving unknown.	f	\N	2025-10-07 21:49:32.234453+00	2025-10-07 21:49:32.234453+00
ce6dc094-310c-4948-bb3a-9e9d2aeaa281	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-ff7c94fb4-cpj9g in omgeving unknown.	f	\N	2025-10-07 22:07:14.424563+00	2025-10-07 22:07:14.424563+00
84d30a4e-9d1b-4500-85bc-ce7a690ef743	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-55fddd5ffd-d2tks in omgeving unknown.	f	\N	2025-10-07 22:15:34.323792+00	2025-10-07 22:15:34.323792+00
effd4ce3-9747-41ad-bc01-e12eda82e3f3	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-cc9c7bf85-k5lpz in omgeving unknown.	f	\N	2025-10-07 23:24:20.221039+00	2025-10-07 23:24:20.221039+00
551427cf-a678-4242-85b0-a05150ec89fa	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-77544b595-v7vbq in omgeving unknown.	f	\N	2025-10-08 10:18:01.125842+00	2025-10-08 10:18:01.125842+00
97813004-88e0-42b8-b30e-b94aa4bddb15	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6f7d48b87f-d8mw4 in omgeving unknown.	f	\N	2025-10-08 10:29:05.331357+00	2025-10-08 10:29:05.331357+00
f97da915-99e5-4cdb-95ee-1bdd2967cca8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-695d8c5c8b-28hzj in omgeving unknown.	f	\N	2025-10-08 10:36:56.326672+00	2025-10-08 10:36:56.326672+00
9ed86468-d319-4444-81c9-6a4f45f211b2	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79fbc85746-t7rsc in omgeving unknown.	f	\N	2025-10-08 11:02:17.624495+00	2025-10-08 11:02:17.624495+00
957765f2-5a22-4c5b-b3a6-240bccf1b6d0	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-56d69b4b54-9g2f8 in omgeving unknown.	f	\N	2025-10-08 11:22:18.323858+00	2025-10-08 11:22:18.323858+00
39df7eb8-5dcc-4bf7-96b9-a13af81c907f	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-55bb56cb7d-cm5xf in omgeving unknown.	f	\N	2025-10-08 11:46:36.031624+00	2025-10-08 11:46:36.031624+00
680ceaec-a8a4-48a4-91f6-54376f7ae0ce	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7b5b7ddbcc-jflsm in omgeving unknown.	f	\N	2025-10-08 13:00:41.730482+00	2025-10-08 13:00:41.730482+00
db209e9e-0f18-41aa-92c5-d5dedbf9fc73	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-785cf585f-dcdfj in omgeving unknown.	f	\N	2025-10-08 16:50:49.122869+00	2025-10-08 16:50:49.122869+00
b4a88dbb-e7c6-4c27-928b-7eca8b2de190	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7587dcb7b8-m2kpl in omgeving unknown.	f	\N	2025-10-09 12:07:03.523202+00	2025-10-09 12:07:03.523202+00
7830e8a4-8c39-4a43-9ac4-359275f08cfa	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7fc998d8f9-mg7cg in omgeving unknown.	f	\N	2025-10-09 15:31:09.12933+00	2025-10-09 15:31:09.12933+00
1e2df6a1-7ab1-4306-8ee6-ac9c97ee8cf4	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-dd98889c9-tkhbg in omgeving unknown.	f	\N	2025-10-09 15:37:16.345979+00	2025-10-09 15:37:16.345979+00
fcb305e6-da43-491e-afa6-3167f268b3a1	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79d5f949fc-ttvg5 in omgeving unknown.	f	\N	2025-10-09 15:52:40.030605+00	2025-10-09 15:52:40.030605+00
dba1ac91-5113-473b-b2d5-c4f4dcb885b8	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-569c65f65-ggp52 in omgeving unknown.	f	\N	2025-10-09 15:59:35.130619+00	2025-10-09 15:59:35.130619+00
b98ab4c5-8d2f-4062-a031-7fc63e343391	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-85f847888-lfk22 in omgeving unknown.	f	\N	2025-10-09 16:07:02.128682+00	2025-10-09 16:07:02.128682+00
3b4e0419-afa1-41a3-8798-6bf460f4ff67	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79c6fdc87f-sb9gm in omgeving unknown.	f	\N	2025-10-10 18:39:18.922957+00	2025-10-10 18:39:18.922957+00
93e6930b-e9e5-4f90-86ef-55c5f16fc7cf	system	low	Service Gestart	DKL Email Service is gestart op DESKTOP-FQOB5OC in omgeving unknown.	f	\N	2025-10-10 21:49:28.692919+00	2025-10-10 21:49:28.692919+00
905e94bf-8ffe-47ea-9915-333224987132	system	low	Service Gestart	DKL Email Service is gestart op DESKTOP-FQOB5OC in omgeving unknown.	f	\N	2025-10-10 21:53:29.693862+00	2025-10-10 21:53:29.693862+00
88198afa-9035-4396-999c-c9ad4ffdf9fd	system	low	Service Gestart	DKL Email Service is gestart op DESKTOP-FQOB5OC in omgeving unknown.	f	\N	2025-10-10 21:55:03.763437+00	2025-10-10 21:55:03.763437+00
555525cc-beb2-4116-93b1-f4a2a412d08c	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-866b75f565-rlj6m in omgeving unknown.	f	\N	2025-10-11 01:40:42.824653+00	2025-10-11 01:40:42.824653+00
25fe1020-501e-45b6-8f38-dd710a700df5	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c9f4c59cd-4xcht in omgeving unknown.	f	\N	2025-10-11 02:21:13.432651+00	2025-10-11 02:21:13.432651+00
9d71327c-53aa-4e03-8366-979c69210a5a	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-55c9654879-8kvsk in omgeving unknown.	f	\N	2025-10-11 16:43:55.633495+00	2025-10-11 16:43:55.633495+00
51228597-b2fb-44b4-af3d-253c4266a563	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5d87df668f-prffv in omgeving unknown.	f	\N	2025-10-11 17:25:26.81996+00	2025-10-11 17:25:26.81996+00
1f3fe283-95af-4458-868c-a92265844f53	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6cccd556c7-t2lbq in omgeving unknown.	f	\N	2025-10-11 19:06:18.219903+00	2025-10-11 19:06:18.219903+00
ad6b9679-a458-4c71-81a8-93e683313383	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-764bf66487-hrbdf in omgeving unknown.	f	\N	2025-10-12 18:51:14.331671+00	2025-10-12 18:51:14.331671+00
122ad898-0bc2-40fc-8781-337e03a2cce3	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-55c6488fb5-42rq9 in omgeving unknown.	f	\N	2025-10-12 19:12:56.131448+00	2025-10-12 19:12:56.131448+00
2b4b49ce-78e6-477f-88d4-24a6f5951805	system	low	Service Gestart	DKL Email Service is gestart op DESKTOP-FQOB5OC in omgeving unknown.	f	\N	2025-10-12 19:22:45.439186+00	2025-10-12 19:22:45.439186+00
58a8c6fd-24d8-4e79-a264-ee6b0f9bdb0d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-566f7b657f-29pd7 in omgeving unknown.	f	\N	2025-10-12 19:30:47.64586+00	2025-10-12 19:30:47.64586+00
9f83ccb8-0ff1-4439-954d-e725a0a986ec	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-685698d8d8-fzpn4 in omgeving unknown.	f	\N	2025-10-12 20:02:44.226204+00	2025-10-12 20:02:44.226204+00
4a7103a7-3b19-4b4d-a9c6-853a9a52fa1b	system	low	Service Gestart	DKL Email Service is gestart op DESKTOP-FQOB5OC in omgeving unknown.	f	\N	2025-10-12 21:02:36.543218+00	2025-10-12 21:02:36.543218+00
a206d8b1-d5e9-4735-b465-da60b4ec2bd1	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-c4d4f6b78-jmqk9 in omgeving unknown.	f	\N	2025-10-12 21:11:31.130047+00	2025-10-12 21:11:31.130047+00
4e4c2869-7fea-4ce7-9434-b6eb6a40326c	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5fc68dc5d7-cjm5x in omgeving unknown.	f	\N	2025-10-12 22:42:33.92738+00	2025-10-12 22:42:33.92738+00
074ee1bd-881d-4980-abe6-19c11877c57a	contact	medium	Nieuw Contactverzoek	<b>JeffTest</b> heeft contact opgenomen via het contactformulier.\n\n<b>Email:</b> laventejeffrey@gmail.com\n\n<b>Bericht:</b>\nTest bericht	f	\N	2025-10-13 16:56:36.739922+00	2025-10-13 16:56:36.739922+00
5eae4aa1-bbb2-4e02-a71d-01f7c6b31f63	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5cddcf5b4c-9x5h9 in omgeving unknown.	f	\N	2025-10-24 19:24:48.549272+00	2025-10-24 19:24:48.549272+00
8b59c817-36e4-4017-825a-fe86787854f5	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-bcd8f8d55-mvjs7 in omgeving unknown.	f	\N	2025-10-24 20:27:05.328236+00	2025-10-24 20:27:05.328236+00
f92161a2-63a3-4134-ae0a-88ae1eb632c7	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-75764dddb7-xqg8m in omgeving unknown.	f	\N	2025-10-24 21:48:00.821027+00	2025-10-24 21:48:00.821027+00
99d6ed78-d078-44eb-be0d-66949a9feebf	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6d556df78c-r5ft5 in omgeving unknown.	f	\N	2025-10-25 12:04:17.822408+00	2025-10-25 12:04:17.822408+00
d78e60d2-8e72-4906-9940-65cfaf85c8ff	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-56dd796dd8-mmcl2 in omgeving unknown.	f	\N	2025-10-25 12:15:05.142094+00	2025-10-25 12:15:05.142094+00
689ae161-61d5-4d22-a628-fffd5521903c	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-57cf96d7c4-s6x7x in omgeving unknown.	f	\N	2025-10-25 12:21:27.526198+00	2025-10-25 12:21:27.526198+00
1f057e49-16d4-49a9-978e-93cc807e0891	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-b5c56f459-8zkqq in omgeving unknown.	f	\N	2025-10-25 15:20:37.026599+00	2025-10-25 15:20:37.026599+00
fe12837f-c5d4-4cfe-add5-3768d2a29f3d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-d749cf94d-kq85c in omgeving unknown.	f	\N	2025-10-25 15:26:40.245859+00	2025-10-25 15:26:40.245859+00
f21d2398-861e-481f-aabd-fd2b25b76616	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-644ffc74fc-2hd99 in omgeving unknown.	f	\N	2025-10-25 16:20:09.622583+00	2025-10-25 16:20:09.622583+00
5da249f9-f849-4e54-9734-d60f2e2f37d1	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-56bdfd44d7-6n8j5 in omgeving unknown.	f	\N	2025-10-30 23:56:33.930266+00	2025-10-30 23:56:33.930266+00
61736f15-514c-49ff-a558-6a1e85b49678	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-6c7cfbf9b7-2vvrh in omgeving unknown.	f	\N	2025-10-31 01:02:31.617888+00	2025-10-31 01:02:31.617888+00
9e6eace9-3c7f-45d6-81fe-64459f5323bb	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-7f9b5d6699-p5t9b in omgeving unknown.	f	\N	2025-11-01 15:44:02.934455+00	2025-11-01 15:44:02.934455+00
858fa078-2361-431f-9e94-66025ff7cc76	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-744bb66db5-b684x in omgeving unknown.	f	\N	2025-11-01 19:34:06.122445+00	2025-11-01 19:34:06.122445+00
1d02b284-f347-46b6-a728-9e4aef050066	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-54979d89d9-k44sr in omgeving unknown.	f	\N	2025-11-01 22:55:29.225671+00	2025-11-01 22:55:29.225671+00
4138a30f-a747-4225-9584-9b58744a68f2	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-78fcfc499d-vv6n2 in omgeving unknown.	f	\N	2025-11-02 00:09:48.727838+00	2025-11-02 00:09:48.727838+00
f83f6544-6c85-45cf-8200-5d532e2a22ce	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-d87757c77-n9zcq in omgeving unknown.	f	\N	2025-11-02 00:22:01.826336+00	2025-11-02 00:22:01.826336+00
b7919178-4f33-4c1a-b2cb-5191bd56c03d	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-8f455fcdf-wjv9h in omgeving unknown.	f	\N	2025-11-02 00:40:36.630212+00	2025-11-02 00:40:36.630212+00
2da385cf-d19f-4ade-9dce-b6e1e12c1b0b	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79cbf5fd8-l54cg in omgeving unknown.	f	\N	2025-11-02 01:03:53.126515+00	2025-11-02 01:03:53.126515+00
34e3e125-c625-405a-b83b-80c524a820ca	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-5d4cd486dd-xjvzb in omgeving unknown.	f	\N	2025-11-02 02:00:35.828337+00	2025-11-02 02:00:35.828337+00
803dffa6-cd09-482d-9c96-79de2a1465b3	system	low	Service Gestart	DKL Email Service is gestart op srv-cv6p735umphs738cd8d0-79bf7cdc66-95jmx in omgeving unknown.	f	\N	2025-11-02 02:09:47.526861+00	2025-11-02 02:09:47.526861+00
\.


--
-- Data for Name: partners; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.partners (id, name, description, logo, website, tier, since, visible, order_number, created_at, updated_at) FROM stdin;
26f11d04-c0da-4755-b2e1-5fcb5e9887d4	Accress	Accres beheert in Apeldoorn ruim 60 locaties, waaronder sporthallen, wijkcentra, zwembaden, kinderboerderijen en een stadspark.\r\n\r\nZij helpen ond bij het halen van ons doel!	https://res.cloudinary.com/dgfuv7wif/image/upload/v1744388421/accres_logo_ochsmg.jpg	https://www.accres.nl/	bronze	2025-04-01	t	5	2025-04-11 16:21:40+00	2025-04-11 16:45:10.227784+00
510d8e4d-6a7d-4ab3-b311-17b32df7f01f	Apeldoorn	Apeldoorn ondersteunt ons in ons doel.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895194/nw1qxouzupzsshckzkab.png	https://www.apeldoorn.nl/	bronze	2024-12-22	t	0	2024-12-22 19:19:55.414207+00	2025-02-26 21:10:59.040907+00
5e8b8390-e637-4c6d-80d8-25ff4155d5a9	Sheeren Loo	Samen met bewoners van SheerenLoo wordt deze loop georganiseerd.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734894570/mtvucaouruenat2cllsi.png	https://www.sheerenloo.nl/	silver	2024-12-22	t	0	2024-12-22 19:09:31.099813+00	2025-01-07 14:11:39.051238+00
9c311d0f-4db7-4da9-9d87-a836e48078cd	Liliane Fonds	Samen maken we ons sterk	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734893709/qsygajx2tdxxbqbfyurr.png	https://www.lilianefonds.nl/	bronze	2024-12-22	t	0	2024-12-22 18:55:09.823207+00	2025-01-08 19:43:58.464461+00
eda49448-3db9-4db5-9d01-a7331879f5b5	De Grote Kerk	De grote kerk ondersteund ons al vanaf het begin. Hier is het allemaal begonnen.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895146/ri4vclttn4nn2wh53wj0.jpg	https://www.grotekerkapeldoorn.nl/	gold	2024-12-22	t	0	2024-12-22 19:19:07.18416+00	2025-04-11 16:46:10.191766+00
\.


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.permissions (id, resource, action, description, is_system_permission, created_at, updated_at) FROM stdin;
138b5f63-5ed3-4679-8d47-a26d5db26e9f	contact	read	Contactformulieren bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
5e209000-9ec7-400d-b26c-87cf351e7069	contact	write	Contactformulieren bewerken (status, notities, antwoorden)	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
d7caf95f-3e98-4458-859a-70246b6df5d6	contact	delete	Contactformulieren verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
38bad1db-f9f8-4050-bb28-26859775db73	aanmelding	read	Aanmeldingen bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
4b46d3ca-f207-40ee-986c-644aaca424d4	aanmelding	write	Aanmeldingen bewerken (status, notities, antwoorden)	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
98d16b03-78e5-484c-8f32-f5875e7bdb90	aanmelding	delete	Aanmeldingen verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
6d83994d-9c33-49ef-b6d0-54325ef73cc7	newsletter	read	Nieuwsbrieven bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
70c71044-4361-409c-868b-5c078a4f1b85	newsletter	write	Nieuwsbrieven aanmaken/bewerken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
15d60692-9cfb-4bd5-b30d-a2bb0a3462cf	newsletter	send	Nieuwsbrieven verzenden	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
f6dfa1f3-1be8-4f8d-a486-51d5d1c7971d	newsletter	delete	Nieuwsbrieven verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
ec7368db-48f0-406b-be17-76b050e9ea1f	email	read	Inkomende emails bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
4baab258-54fd-46f3-b7af-6112109d9d1e	email	write	Emails bewerken (markeren als verwerkt)	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
2b110b0f-c8eb-412e-b9d6-26933c373099	email	delete	Emails verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
ec205c8b-1e9e-427b-96bc-3af2179ae513	email	fetch	Nieuwe emails ophalen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
92b3327c-6a62-47d2-98f1-e0d294d202b6	admin_email	send	Emails verzenden namens admin	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
9ab65b7b-3a0b-4a28-80c8-3e4301dfbfa5	user	read	Gebruikers bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
dd485ead-8246-4f97-a91d-7818efc801e4	user	write	Gebruikers aanmaken/bewerken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
ff7ce90d-1ec0-4139-a7cd-36ce8711495f	user	delete	Gebruikers verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
3de90dd9-6a0d-4b00-a8f8-c605c0338a62	user	manage_roles	Gebruikersrollen beheren	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
78e5ad68-1209-4662-9abf-6e9206d39d0b	chat	read	Chat kanalen en berichten bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
4e71123a-14d5-49c6-92df-c6a76e01ecac	chat	write	Berichten verzenden	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
7c6478a7-b34a-47cf-ab33-050d7041787b	chat	manage_channel	Kanalen aanmaken/beheren	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
31ab0eff-d7ac-45eb-94af-ef97123e3c56	chat	moderate	Berichten modereren (bewerken/verwijderen)	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
ee5be0e4-9818-4d0d-aa65-53454c7c08c1	notification	read	Notificaties bekijken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
96c38e72-1cbc-4761-b45f-c8cae777d26a	notification	write	Notificaties aanmaken	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
7771de13-161e-4962-ad7c-1f8c55c0b988	notification	delete	Notificaties verwijderen	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
665c5885-40d6-4480-bbc9-37e22e4315a2	system	admin	Volledige systeemtoegang	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00
fe009820-67d5-4e7a-a14b-8c89b0e40f42	staff	access	Toegang tot staff functies	t	2025-10-07 18:06:57.049357+00	2025-10-07 18:06:57.049357+00
ffb331b4-8682-41aa-9e24-1ee10de6932f	admin	access	Volledige admin toegang	t	2025-10-07 19:10:04.648759+00	2025-10-07 19:10:04.648759+00
20cae3c5-9649-4533-81f8-0c8a223a1efd	photo	read	Foto's bekijken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
a17c0e6e-132e-442c-bfc1-499119c862de	photo	write	Foto's uploaden/bewerken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
13a0ce40-f06a-4006-8e47-e848a843a175	photo	delete	Foto's verwijderen	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
e10a5b01-f824-4236-abc7-bb9ed0d01bd0	album	read	Albums bekijken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
93e95fb8-91e5-4e5c-96df-8594c93588fd	album	write	Albums aanmaken/bewerken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
1cdd5c92-40ac-4f3c-9699-8f37d40c47a6	album	delete	Albums verwijderen	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
d8cf9638-51fe-4639-8c31-ca6b813fdaf2	partner	read	Partners bekijken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
3e489e46-1b0d-4c3c-a3ff-ee3a035a2222	partner	write	Partners aanmaken/bewerken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
274746c2-93fe-4c9a-9257-521710a505c5	partner	delete	Partners verwijderen	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
f7f09f53-d4f4-4fae-9f49-c0f88797aa80	sponsor	read	Sponsors bekijken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
32632262-5c8f-426b-9dfd-d1974ade1fb2	sponsor	write	Sponsors aanmaken/bewerken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
48353034-3afc-407d-8f33-a87d69923bac	sponsor	delete	Sponsors verwijderen	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
cdca2b10-8719-44f1-a131-6f29a661ffd1	video	read	Video's bekijken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
ed2deabf-a1d0-442d-a8e3-0af563e03d99	video	write	Video's uploaden/bewerken	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
264dc5ea-f523-42a4-ae30-21c60221d779	video	delete	Video's verwijderen	t	2025-10-07 20:00:54.648925+00	2025-10-07 20:00:54.648925+00
d8a4f321-e951-4ed5-b095-b96923e4d50d	role	read	Rollen bekijken	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
2a5a5b7c-cd86-4523-b91d-57b665ab2db9	role	write	Rollen aanmaken/bewerken	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
1f46e8d5-73ee-42d3-8691-939cb83f8d79	role	delete	Rollen verwijderen	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
d8e96391-c815-484f-be94-74912ae89616	role	manage_permissions	Rol permissies beheren	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
348f8695-5d3c-4c22-be33-41c3a1fb1d58	permission	read	Permissies bekijken	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
716e3672-b648-4ecd-a97b-b6cbbf3bb4e9	permission	write	Permissies aanmaken/bewerken	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
d3fb8dea-aa6d-47eb-b9e4-92b9d41f394f	permission	delete	Permissies verwijderen	t	2025-10-08 12:51:57.512242+00	2025-10-08 12:51:57.512242+00
d0161ac5-6ab8-4659-9b29-5d4b99e54a4f	radio_recording	read	Radio opnames bekijken	t	2025-10-11 01:40:42.354802+00	2025-10-11 01:40:42.354802+00
00b888b5-23c1-4d0b-b1a6-44f00ddb1fbc	radio_recording	write	Radio opnames aanmaken/bewerken	t	2025-10-11 01:40:42.354802+00	2025-10-11 01:40:42.354802+00
fe3a4df1-571f-4998-836f-ed8a06d917b6	radio_recording	delete	Radio opnames verwijderen	t	2025-10-11 01:40:42.354802+00	2025-10-11 01:40:42.354802+00
557ef223-08f6-489c-ba18-9b36527f4d03	program_schedule	read	Programma bekijken	t	2025-10-11 01:40:42.360115+00	2025-10-11 01:40:42.360115+00
8f2bae59-58f7-4d3d-a90d-9532df4e2d19	program_schedule	write	Programma aanmaken/bewerken	t	2025-10-11 01:40:42.360115+00	2025-10-11 01:40:42.360115+00
167d7532-e67d-42c3-8330-d729fa877c40	program_schedule	delete	Programma verwijderen	t	2025-10-11 01:40:42.360115+00	2025-10-11 01:40:42.360115+00
06ba39c0-1cfd-4e19-87b2-9ec00b1e705e	social_embed	read	Social embeds bekijken	t	2025-10-11 01:40:42.445833+00	2025-10-11 01:40:42.445833+00
38ae19a5-834a-464a-82d0-96d6d4c584fe	social_embed	write	Social embeds aanmaken/bewerken	t	2025-10-11 01:40:42.445833+00	2025-10-11 01:40:42.445833+00
b3d294e1-28bf-4aba-83eb-556ee191b5bb	social_embed	delete	Social embeds verwijderen	t	2025-10-11 01:40:42.445833+00	2025-10-11 01:40:42.445833+00
4fd44556-4e94-4378-8b36-d3ae23f51b3e	social_link	read	Social links bekijken	t	2025-10-11 01:40:42.447813+00	2025-10-11 01:40:42.447813+00
cb8c4930-6a0b-4479-a521-e4b2d3468955	social_link	write	Social links aanmaken/bewerken	t	2025-10-11 01:40:42.447813+00	2025-10-11 01:40:42.447813+00
5a4a8eaf-df25-4ee2-a057-9a329cb8c9ce	social_link	delete	Social links verwijderen	t	2025-10-11 01:40:42.447813+00	2025-10-11 01:40:42.447813+00
7fff7adf-f84c-44b6-bc6b-7a03249330cf	under_construction	read	Under construction bekijken	t	2025-10-11 01:40:42.449408+00	2025-10-11 01:40:42.449408+00
f49a50eb-1b4a-40de-8a86-11fbd44bdb58	under_construction	write	Under construction aanmaken/bewerken	t	2025-10-11 01:40:42.449408+00	2025-10-11 01:40:42.449408+00
4f845d71-7e33-4a2a-95a8-89b70fda46f0	under_construction	delete	Under construction verwijderen	t	2025-10-11 01:40:42.449408+00	2025-10-11 01:40:42.449408+00
a24c3aee-10a4-437f-8772-0db10e02dbc0	title_section	read	Title section bekijken	t	2025-10-12 18:51:14.154327+00	2025-10-12 18:51:14.154327+00
2746de28-baac-47b4-a14a-c381983ea512	title_section	write	Title section aanmaken/bewerken	t	2025-10-12 18:51:14.154327+00	2025-10-12 18:51:14.154327+00
2f06aeab-e4ad-4f44-80ba-4d9643c3a790	title_section	delete	Title section verwijderen	t	2025-10-12 18:51:14.154327+00	2025-10-12 18:51:14.154327+00
8cf79304-7204-4df8-9317-5c420d6fb7bc	steps	read	Stappen data bekijken (dashboard, totaal, verdeling)	t	2025-10-24 19:24:48.264458+00	2025-10-24 19:24:48.264458+00
7181ad0f-19b9-41c0-a1b6-b1c06c223292	steps	write	Stappen bijwerken voor deelnemers	t	2025-10-24 19:24:48.264458+00	2025-10-24 19:24:48.264458+00
94e18583-8d38-4db6-91e3-95d26d5a4c0d	steps	read_all	Alle deelnemers stappen bekijken (admin/staff)	t	2025-10-25 15:20:36.749699+00	2025-10-25 15:20:36.749699+00
abb03e0b-d835-405c-bdb5-4d5111b0a209	steps	write_all	Alle deelnemers stappen bijwerken (admin/staff)	t	2025-10-25 15:20:36.749699+00	2025-10-25 15:20:36.749699+00
9d1ea9ec-4042-45fa-9d88-69d2b557535c	steps	manage	Volledige steps beheer (route funds, etc.)	t	2025-10-25 15:20:36.749699+00	2025-10-25 15:20:36.749699+00
556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	steps	read_total	Totaal aantal stappen van alle deelnemers bekijken	t	2025-10-25 15:26:39.960312+00	2025-10-25 15:26:39.960312+00
\.


--
-- Data for Name: photos; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.photos (id, url, alt_text, visible, thumbnail_url, created_at, updated_at, title, description, year, cloudinary_folder) FROM stdin;
ee20de98-2fa8-4e23-8bf1-0b705b55aa7c	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747636922/vv5v84gadrf02rl3iiji.jpg	1000016660	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747636922/vv5v84gadrf02rl3iiji.jpg	2025-05-19 06:42:03.254692+00	2025-05-19 06:42:03.254692+00	1000016660	\N	\N	\N
754ceb60-f4d8-4434-9338-065339e636e8	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513150/d23b6xefqsaxekgnkpeq.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513150/d23b6xefqsaxekgnkpeq.jpg	2025-05-17 20:19:11.491775+00	2025-05-17 20:19:11.491775+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
eafe6902-3d8e-4929-8cf3-c9b49429adac	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513150/s82ykrgnv8zuwxrraixd.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513150/s82ykrgnv8zuwxrraixd.jpg	2025-05-17 20:19:10.549969+00	2025-05-17 20:19:10.549969+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
db8dbb96-c528-4ae2-adf7-fc0cd165ad3c	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513149/t2sanzejw8lztqbhbu0h.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513149/t2sanzejw8lztqbhbu0h.jpg	2025-05-17 20:19:09.809387+00	2025-05-17 20:19:09.809387+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
19312aba-29ab-4169-83eb-934a80763ad8	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513148/c6bzsdf9osgub9cundss.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513148/c6bzsdf9osgub9cundss.jpg	2025-05-17 20:19:08.91355+00	2025-05-17 20:19:08.91355+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
7eb75484-b2cb-41a7-bdc6-b16a8811ecca	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513147/jwkot927pq1nb8kiwg4h.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513147/jwkot927pq1nb8kiwg4h.jpg	2025-05-17 20:19:08.151916+00	2025-05-17 20:19:08.151916+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
2debaea6-d5ab-4dd3-ae1b-15919c80245c	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513146/k9o1so9g7jh98dpigakj.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513146/k9o1so9g7jh98dpigakj.jpg	2025-05-17 20:19:07.245565+00	2025-05-17 20:19:07.245565+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
096c8248-2b05-4337-b256-f857419fa9bd	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513145/cg1my6knmyme8b9ocpre.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513145/cg1my6knmyme8b9ocpre.jpg	2025-05-17 20:19:06.156387+00	2025-05-17 20:19:06.156387+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
7d525fa5-192a-4b5f-9f03-7b8f9d4eb038	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513144/tlhlt2etlp0wd3hhh1e1.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513144/tlhlt2etlp0wd3hhh1e1.jpg	2025-05-17 20:19:05.098892+00	2025-05-17 20:19:05.098892+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
9db233dc-58ff-4352-be6e-fb0d3b30798a	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513143/wmtpajgjjvlf7gdpromp.jpg	WhatsApp Image 2025-05-17 at 20	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513143/wmtpajgjjvlf7gdpromp.jpg	2025-05-17 20:19:04.357247+00	2025-05-17 20:19:04.357247+00	WhatsApp Image 2025-05-17 at 20	\N	\N	\N
fded98ba-c7fa-4c42-ad2e-6c7afea5b8fc	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513143/jtcxza8j43cwwyynrq5g.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513143/jtcxza8j43cwwyynrq5g.jpg	2025-05-17 20:19:03.527768+00	2025-05-17 20:19:03.527768+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
c5cc7e20-9e4d-4e09-bcf5-0c2edcb83360	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513142/w9nvohxsntxfiy4cevbf.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513142/w9nvohxsntxfiy4cevbf.jpg	2025-05-17 20:19:02.830463+00	2025-05-17 20:19:02.830463+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
ac4fdceb-3598-4ba9-9d61-4f40ef0e44e3	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513141/dol0bmbwzaamhvf7zeni.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513141/dol0bmbwzaamhvf7zeni.jpg	2025-05-17 20:19:02.129184+00	2025-05-17 20:19:02.129184+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
be16d3ae-c8a4-457b-9806-fee919bc79a2	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513140/grdni6fzojt466urmkya.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513140/grdni6fzojt466urmkya.jpg	2025-05-17 20:19:01.404594+00	2025-05-17 20:19:01.404594+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
c199f984-64e5-4405-b23c-e6ff4a3eaed3	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513140/men0m6mk5f505uoaclhf.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513140/men0m6mk5f505uoaclhf.jpg	2025-05-17 20:19:00.628702+00	2025-05-17 20:19:00.628702+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
e971f01a-1d4a-4400-bb54-3428fbc69a98	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513139/j19kt4rorb9ybtpx0x1r.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513139/j19kt4rorb9ybtpx0x1r.jpg	2025-05-17 20:18:59.94824+00	2025-05-17 20:18:59.94824+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
991bfd41-1365-4bb4-9ac8-8519fc9bfb32	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513138/uelmhlfiuccmv2slqbaw.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513138/uelmhlfiuccmv2slqbaw.jpg	2025-05-17 20:18:59.143222+00	2025-05-17 20:18:59.143222+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
18eb6cd7-b9a4-4b21-8363-ef77f420ac09	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513137/zgflyhmixak9ci0warv4.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513137/zgflyhmixak9ci0warv4.jpg	2025-05-17 20:18:58.33864+00	2025-05-17 20:18:58.33864+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
83e75346-54b4-40b7-8ba4-2d43b3d8c867	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513136/iimq27dhkyimotercqan.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513136/iimq27dhkyimotercqan.jpg	2025-05-17 20:18:57.33427+00	2025-05-17 20:18:57.33427+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
49f983ab-ec63-4489-93ce-ba9272ba7f49	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513135/bwpkhyzltxncxkza2lsz.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513135/bwpkhyzltxncxkza2lsz.jpg	2025-05-17 20:18:56.358544+00	2025-05-17 20:18:56.358544+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
6649c67b-2eb5-4d29-9e1c-0373ea3b1771	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513134/l6o1ibr7lzzbn7glpjcu.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513134/l6o1ibr7lzzbn7glpjcu.jpg	2025-05-17 20:18:55.523962+00	2025-05-17 20:18:55.523962+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
ec209a88-3f3b-4168-a4c9-1c9275edcffb	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513133/khtzc08kc7wgkta5rh7s.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513133/khtzc08kc7wgkta5rh7s.jpg	2025-05-17 20:18:54.564013+00	2025-05-17 20:18:54.564013+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
c7ddfbd8-6bf4-4bd9-9285-f6492b98bef4	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513133/unhly8fepi83vegupc6a.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513133/unhly8fepi83vegupc6a.jpg	2025-05-17 20:18:53.737493+00	2025-05-17 20:18:53.737493+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
10fff5f8-2701-4f34-86f1-b063b252f35a	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513132/zoqmk50gcxuqkda0sqoe.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513132/zoqmk50gcxuqkda0sqoe.jpg	2025-05-17 20:18:53.033148+00	2025-05-17 20:18:53.033148+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
76278110-be4c-4ffc-bc35-e2d4b14fa42f	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513131/hfo3n8mzetzeqefr418d.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513131/hfo3n8mzetzeqefr418d.jpg	2025-05-17 20:18:52.092844+00	2025-05-17 20:18:52.092844+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
baf57b95-0ab7-4217-af98-1453a9e6a938	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513130/tbsbwsnimr5fiv1i7pkn.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513130/tbsbwsnimr5fiv1i7pkn.jpg	2025-05-17 20:18:51.341415+00	2025-05-17 20:18:51.341415+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
0911da4d-6169-4383-962b-9a640b5eac0e	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513129/p8aoklpqxch3jlbukl3r.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513129/p8aoklpqxch3jlbukl3r.jpg	2025-05-17 20:18:50.351733+00	2025-05-17 20:18:50.351733+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
1350bb1c-b132-4250-82b8-58efb05fb53b	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513128/thtibyrsflsuen2lotv5.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513128/thtibyrsflsuen2lotv5.jpg	2025-05-17 20:18:49.478846+00	2025-05-17 20:18:49.478846+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
08362e92-340a-432a-b306-153ad27ee686	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513127/thhu8mxkqhfhxjydi5ai.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513127/thhu8mxkqhfhxjydi5ai.jpg	2025-05-17 20:18:48.433401+00	2025-05-17 20:18:48.433401+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
7c296b82-d35d-4065-8523-dffc976688f7	https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513126/j5l0h7pakhadhwcrxo8o.jpg	WhatsApp Image 2025-05-17 at 19	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513126/j5l0h7pakhadhwcrxo8o.jpg	2025-05-17 20:18:47.544035+00	2025-05-17 20:18:47.544035+00	WhatsApp Image 2025-05-17 at 19	\N	\N	\N
5b38433f-3c33-49cb-950a-e0bc7eb21a80	https://res.cloudinary.com/dgfuv7wif/image/upload/v1745793210/rirkdj5bav7k0pvvtfoq.jpg	1000014905	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745793210/rirkdj5bav7k0pvvtfoq.jpg	2025-04-27 22:33:31.229221+00	2025-04-27 22:33:31.229221+00	1000014905	\N	\N	\N
dbb68d91-3f6f-46c3-b751-a13f34269b79	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734808286/ubjiz1fal82jh42nzjmg.jpg	ubjiz1fal82jh42nzjmg	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1734808286/ubjiz1fal82jh42nzjmg.jpg	2025-04-19 00:30:16.928614+00	2025-04-20 02:08:06.401666+00	321	\N	2025	\N
f4ce7f8e-b573-4602-9e12-69e0585df779	https://res.cloudinary.com/dgfuv7wif/image/upload/v1745022279/dlrhmdl4gcddunkqzkei.jpg	dlrhmdl4gcddunkqzkei	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745022279/dlrhmdl4gcddunkqzkei.jpg	2025-04-19 00:30:16.928614+00	2025-04-20 02:07:54.136061+00	123	\N	2025	\N
e7b84300-7158-475b-a79b-10e97b416d58	https://res.cloudinary.com/dgfuv7wif/image/upload/v1745022320/k7qobfrinlsjrqrxxdfy.jpg	k7qobfrinlsjrqrxxdfy	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745022320/k7qobfrinlsjrqrxxdfy.jpg	2025-04-19 00:30:16.928614+00	2025-04-20 02:08:27.195419+00	231	\N	2025	\N
26188a5b-d542-4674-8ae9-b63520fbd4b2	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734808313/bmkf9wcfwrseamcgdu9t.jpg	bmkf9wcfwrseamcgdu9t	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1734808313/bmkf9wcfwrseamcgdu9t.jpg	2025-04-19 00:30:16.928614+00	2025-04-20 02:08:16.199794+00	213	\N	2025	\N
d78d9d0f-29ac-42a1-951a-a0bcb355d227	https://res.cloudinary.com/dgfuv7wif/image/upload/v1739543787/vdohoeldmm6iiv9cwikm.jpg	1000011474	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1739543787/vdohoeldmm6iiv9cwikm.jpg	2025-02-14 14:36:28.503034+00	2025-04-20 02:09:25.344705+00	011474	\N	2025	\N
af027789-d718-4899-88b2-d8f0814b73ee	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163484/671433bb6de927107d736c2a_DKL19-p-800_slbxvx.webp	Koninklijke Loop 2023 - Sfeerimpressie	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163484/671433bb6de927107d736c2a_DKL19-p-800_slbxvx.webp	2024-12-22 20:03:33.896605+00	2025-04-17 19:05:18.330672+00	Koninklijke Loop 2023 - Sfeerimpressie	\N	2024	\N
6fd90008-f1ca-4e3d-9915-32b914208239	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714330a3210f9f512505de3_dkl2-p-800_pvztxb.webp	Koninklijke Loop 2023 - Parcours impressie	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714330a3210f9f512505de3_dkl2-p-800_pvztxb.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:37.910515+00	Koninklijke Loop 2023 - Parcours impressie	\N	2024	\N
3acaeca2-aa51-4cd6-9fad-7b20b607cba7	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163486/6714336722ccd9c2f9c73177_DKL12-p-800_hguzw1.webp	Koninklijke Loop 2023 - Evenement overzicht	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163486/6714336722ccd9c2f9c73177_DKL12-p-800_hguzw1.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:27.251872+00	Koninklijke Loop 2023 - Evenement overzicht	\N	2024	\N
4059d3cf-0e92-44c2-bbfa-93235dc19eae	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163486/67143385825a5682baca0273_DKL13-p-800_m9sc65.webp	Koninklijke Loop 2023 - Sfeerbeeld	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163486/67143385825a5682baca0273_DKL13-p-800_m9sc65.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:30.097262+00	Koninklijke Loop 2023 - Sfeerbeeld	\N	2024	\N
245ddf1a-61ab-4b64-b8b2-13d5079d6592	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714335d14c39bfea035fa1b_DKL9-p-800_umfakb.webp	Koninklijke Loop 2023 - Finish moment	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714335d14c39bfea035fa1b_DKL9-p-800_umfakb.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:24.427491+00	Koninklijke Loop 2023 - Finish moment	\N	2024	\N
0334c63c-230a-46c7-b87d-3f4a5cc946c0	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714333d9d3ac37ffaa1d539_DKL6-p-800_o16v1m.webp	Koninklijke Loop 2023 - Deelnemers samen	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714333d9d3ac37ffaa1d539_DKL6-p-800_o16v1m.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:19.507503+00	Koninklijke Loop 2023 - Deelnemers samen		2024	\N
e1350a21-3d42-4252-bc56-5d9ce264ff41	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/67143326f367de40be79ebfa_DKL4-p-800_ceagh7.webp	Koninklijke Loop 2023 - Lopers onderweg	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/67143326f367de40be79ebfa_DKL4-p-800_ceagh7.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:56.31221+00	Koninklijke Loop 2023 - Lopers onderweg	\N	2024	\N
8a4d5c20-ea73-4336-8c6b-af7a197ef7c2	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163484/671433481981522775f532d1_DKL7-p-800_zdihdc.webp	Koninklijke Loop 2023 - Deelnemers in actie	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163484/671433481981522775f532d1_DKL7-p-800_zdihdc.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:44.13556+00	Koninklijke Loop 2023 - Deelnemers in actie	\N	2024	\N
7431c540-2aa5-42db-b5c5-df55a17856ae	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/67143316c21bdbd9d9800bcf_DKL3-p-800_kcx0b3.webp	Koninklijke Loop 2023 - Groepsfoto	t	https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/67143316c21bdbd9d9800bcf_DKL3-p-800_kcx0b3.webp	2024-12-22 20:03:33.896605+00	2025-03-17 19:50:41.386987+00	Koninklijke Loop 2023 - Groepsfoto	\N	2024	\N
\.


--
-- Data for Name: program_schedule; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.program_schedule (id, "time", event_description, category, icon_name, order_number, visible, created_at, updated_at, latitude, longitude) FROM stdin;
075095c7-925d-411e-bebc-a7fc96a3000a	12:00u	Aanvang Deelnemers 10km bij het coördinatiepunt	Aanvang	aanvang	50	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
1302ac11-1b7a-4738-a2dc-535990a09e69	10:15u	Aanvang Deelnemers 15km bij het coördinatiepunt	Aanvang	aanvang	10	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
1572e571-5f90-4cd9-b930-173d31df0124	11:05u	Deelnemers 15km aanwezig startpunt (Kootwijk)	Aanwezig	aanwezig	30	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.18474064	5.77074940
16f07558-27cc-4e70-8d2f-4093d5e47009	15:35u	START 2,5KM, Hervatting 6km, 10km en 15km	Start	start	190	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.22044762	5.92889575
2eb4673c-6bde-465a-b4e8-27425bc32d54	15:15u	Verwachte aankomst 15, 10, 6 km lopers bij rustpunt (Berg & Bos - 15 min pauze)	Rustpunt	rustpunt	180	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
425c461a-cee7-480c-a399-7d469cc7efbe	12:50u	Deelnemers 10km aanwezig bij het startpunt (Halte Assel)	Aanwezig	aanwezig	80	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.20071362	5.83602324
460625b7-ec90-446d-a895-2ddaefb98335	14:15u	START 6KM, Hervatting 10km en 15km	Start	start	140	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.21916438	5.87255921
4f201162-2524-4b03-b0ea-0fb5ccaf29c3	15:00u	Vertrek deelnemers 2,5 km met de pendelbussen naar het startpunt 2,5km	Vertrek	vertrek	160	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
5af9c5b7-9b73-44aa-9153-61e2277b9233	17:00u – 18:00u	INHULDIGINGSFEEST	Feest	feest	220	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
652b1d86-b062-4b63-bbf0-5294c71979d0	12:45u	Verwachte aankomst 15 km lopers bij rustpunt (Halte Assel - 15 min pauze)	Rustpunt	rustpunt	70	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
65dd7eb0-4de3-4d54-99b7-2f278a08ece4	12:30u	Vertrek deelnemers 10km met de pendelbussen naar het startpunt 10km	Vertrek	vertrek	60	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
6730216e-9a4c-4df5-aaca-119da8595eef	10:45u	Vertrek pendelbussen naar startpunt 15km	Vertrek	vertrek	20	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
7ae60214-dc2d-4b23-a4d1-0993fbb56e46	14:00u	Deelnemers 6km aanwezig bij het startpunt (Hoog Soeren)	Aanwezig	aanwezig	130	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.21916438	5.87255921
91666bd3-32fd-4054-ab31-4d9f7532c0ce	14:30u	Aanvang Deelnemers 2,5km bij het coördinatiepunt	Aanvang	aanvang	150	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
9564d363-17d5-45b4-b868-45d165a82c72	16:10u – 16:30u	FINISH	Finish	finish	210	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
9ff53726-0e9a-48de-b4f9-35cfa7666756	15:55u	Aankomst bij De Naald / START INHULDIGINGSLOOP	Aankomst	aankomst	200	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
b9c046ea-de08-48e5-995e-9d4a98d76b6e	14:00u	Verwachte aankomst 15, 10 km lopers bij rustpunt (Hoog Soeren - 15 min pauze)	Rustpunt	rustpunt	120	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
bc30a0ea-0f22-443a-8124-0cd52e10a2b3	11:15u	START 15KM	Start	start	40	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.18474064	5.77074940
cf1514db-2e06-45ab-891b-76736eb308d3	13:00u	START 10KM, Hervatting 15km	Start	start	90	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.20071362	5.83602324
de0cd4c1-2fe0-4f13-9b94-00c03ce38527	15:05u	Deelnemers 2,5km aanwezig bij het startpunt (Berg & Bos)	Aanwezig	aanwezig	170	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	52.22044762	5.92889575
fd80912e-f112-4bee-a2dd-570c4cf88c89	13:15u	Aanvang Deelnemers 6km bij het coördinatiepunt	Aanvang	aanvang	100	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
fec9b64c-4cdd-4fc1-a507-9c222fdeb958	13:45u	Vertrek deelnemers 6 km met de pendelbussen naar het startpunt 6km	Vertrek	vertrek	110	t	2025-04-15 21:57:08.444242+00	2025-04-15 21:57:08.444242+00	\N	\N
\.


--
-- Data for Name: radio_recordings; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.radio_recordings (id, title, description, date, audio_url, thumbnail_url, visible, order_number, created_at, updated_at) FROM stdin;
a6e73425-6af2-4bbd-84f8-67b16d195f99	De koninklijke Loop 2025 uitzending!	Luister naar het live radioverslag van De Koninklijke Loop 2025, uitgezonden op RTV Apeldoorn, met interviews met de organisatie!\r\n\r\n	14 mei 2025	https://res.cloudinary.com/dgfuv7wif/video/upload/v1747733438/DKLRTV2025_dpdydc.wav	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	t	1	2025-05-20 09:32:01+00	2025-05-20 09:33:27.540519+00
c2d84ca9-97a9-4990-9ee4-1fe1718a8c5b	Radioverslag Koninklijke Loop 2024 (RTV Apeldoorn)	Luister naar het live radioverslag van De Koninklijke Loop 2024, uitgezonden op RTV Apeldoorn, met interviews en sfeerimpressies.	15 mei 2024	https://res.cloudinary.com/dgfuv7wif/video/upload/v1714042357/matinee_1_nbm0ph.wav	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	t	2	2025-04-08 19:50:35.125804+00	2025-05-20 09:32:30.549188+00
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.refresh_tokens (id, user_id, token, expires_at, created_at, revoked_at, is_revoked) FROM stdin;
5589c6df-6964-4df2-98f7-d13137d4c0d2	7157f3f6-da85-4058-9d38-19133ec93b03	R1eBX10a-u3L0K2sTSHfyEi3MuolygDDs-Q1e-8JPVo=	2025-10-15 16:51:50.795866	2025-10-08 16:51:50.79701	\N	f
6156ec3d-b3f0-4303-bc31-da5fbab76815	7157f3f6-da85-4058-9d38-19133ec93b03	akDVNeAuDJf5nV_peayOkBNPNlw_xAC8Fh_13g3-oTg=	2025-10-15 16:54:39.608723	2025-10-08 16:54:39.609091	\N	f
804f329a-a075-4af1-919d-6f4ec4876994	7157f3f6-da85-4058-9d38-19133ec93b03	DzVTIlNFApJixeI-lOd5LVXrAB3uJVb2ablU1eNIn40=	2025-10-15 16:56:38.401387	2025-10-08 16:56:38.40183	\N	f
a18e1b85-1c97-4f5f-bd5c-1ea72b79cd40	748320dd-5b5e-4434-ad7e-8a405fd6266f	AxVKn-EVCQnF6SuKyeV7Cja8DnCjWnlxuZD4UAe78b0=	2025-10-15 17:25:10.828121	2025-10-08 17:25:10.828534	\N	f
48d3c44a-f910-40cd-894d-7a9bad95ea5b	748320dd-5b5e-4434-ad7e-8a405fd6266f	waC4nEA_GUqXN5g7YzHCQKVSfIx0tdnd5hTJKMClbAY=	2025-10-15 16:59:23.885932	2025-10-08 16:59:23.886296	2025-10-08 17:25:10.829767	t
95d9d08f-830f-4f34-977c-576d89f32aac	7157f3f6-da85-4058-9d38-19133ec93b03	sl_bE2yLLmyjLmRwYybmJwdB6kQ943Q1n0EiO3EdyCc=	2025-10-15 17:25:20.617541	2025-10-08 17:25:20.617851	\N	f
128d2286-76c3-4823-a161-bdcb862b0ef4	7157f3f6-da85-4058-9d38-19133ec93b03	dXSZu9KazNjlrzdy83-GZSM0a6fB2xumhbYE_JZMliw=	2025-10-15 17:30:30.665717	2025-10-08 17:30:30.666033	\N	f
ac560ec3-f504-4211-8dd0-cf73fb9b9750	7157f3f6-da85-4058-9d38-19133ec93b03	vfyGzfbvqRM4Bdcw96JBMsQ6hQcIHZ-hsHjp0j9rNeU=	2025-10-15 17:51:13.898601	2025-10-08 17:51:13.8996	\N	f
177e7ac2-ef23-4bf7-ae68-1817cc008754	748320dd-5b5e-4434-ad7e-8a405fd6266f	SS59wJKU-MfkzN_vNMFRUFCxib4-FvMJJcjhu1PG6C0=	2025-10-15 18:12:47.460465	2025-10-08 18:12:47.46085	\N	f
342029b6-f506-4ac5-a651-1c8fbb4d8dea	7157f3f6-da85-4058-9d38-19133ec93b03	tw4o1rwfmWOP5HT0lSN61agAdRo4blPIBkecTgv_7AU=	2025-10-15 18:22:50.299995	2025-10-08 18:22:50.30033	\N	f
85a3b9ed-1921-48f1-8cf8-2fcd24030eac	7157f3f6-da85-4058-9d38-19133ec93b03	iP0z98ooBScaOxqsRFcahUImsysTaNo78p65ooRX5RA=	2025-10-15 19:27:31.099281	2025-10-08 19:27:31.099648	\N	f
b190ffa8-5742-4ff6-8252-50ea0fd7710f	7157f3f6-da85-4058-9d38-19133ec93b03	dsXq0LnR6UH8suJ_bYgOBEwdc9C0ZU2RFATSPkgGsIY=	2025-10-15 18:48:14.080694	2025-10-08 18:48:14.081046	2025-10-08 19:27:31.101898	t
613f2338-eda6-4425-aea2-050c0d4b3fe0	7157f3f6-da85-4058-9d38-19133ec93b03	LAMyAfwPzBoyM_rwZr14eofs_npMgnOVhdT-P65FMwI=	2025-10-15 19:29:06.698337	2025-10-08 19:29:06.698617	\N	f
869e5def-9383-4f71-9702-5cdaf32ee309	7157f3f6-da85-4058-9d38-19133ec93b03	WNqf15EeyfrNzyY4zTbuIimEnnBlWe-RERc8TLgh8L0=	2025-10-15 19:52:14.463641	2025-10-08 19:52:14.463993	\N	f
26f481b4-0e36-4d8c-b4bd-1f52e9e760db	7157f3f6-da85-4058-9d38-19133ec93b03	iWTrihKs2-dLQmCjnC0R2yWntEmIDQWYMWebp6KavE0=	2025-10-15 19:52:15.460275	2025-10-08 19:52:15.460606	\N	f
edb94647-d1af-474d-9b33-ddb6249bb620	7157f3f6-da85-4058-9d38-19133ec93b03	IU5jyJu6cFL9DM9UeVvHnslrL0qaNZiv5WyexGjSnpg=	2025-10-15 20:07:59.165205	2025-10-08 20:07:59.165523	\N	f
5c8459fe-6e91-442c-9c30-df060d69d5cb	7157f3f6-da85-4058-9d38-19133ec93b03	yUzf0aiLzbCV4kn9D-qQZYfQCnTEzm3ROmE5-qyi9d8=	2025-10-15 21:26:57.937712	2025-10-08 21:26:57.938022	\N	f
ff5262db-6ebc-4ac0-bb11-60e8cd042d63	7157f3f6-da85-4058-9d38-19133ec93b03	IQs7TRSe-0c7okNTcv_ttrHKggom7Vg1jocMo_TbkNg=	2025-10-15 18:09:52.567472	2025-10-08 18:09:52.568609	2025-10-08 21:26:57.939822	t
cc8ab132-b4f8-4e8a-a8a8-b818cdf2a598	7157f3f6-da85-4058-9d38-19133ec93b03	XIGAR012W0LnumU4pRVQZWHc9APlGp8ATpOsufRi_0U=	2025-10-16 01:21:46.275719	2025-10-09 01:21:46.276043	\N	f
3e15308a-ac25-44d5-8f6a-e68aa0102d9a	7157f3f6-da85-4058-9d38-19133ec93b03	KHwnY_bnaLYYFsw7vkoL3LPVPJ7KmA6Nxx-QgYeBkcM=	2025-10-16 01:38:19.68964	2025-10-09 01:38:19.689958	2025-10-09 10:12:39.035454	t
9e13af5b-cc61-4206-b966-256e0b18862b	7157f3f6-da85-4058-9d38-19133ec93b03	okqYDDW3sQcu6Fum5-Q1mSWPpmnwN8q_fqs54dB95og=	2025-10-16 14:10:20.56308	2025-10-09 14:10:20.563546	\N	f
7ab98e6e-0e00-4b7d-ac7c-d6c6a5528157	7157f3f6-da85-4058-9d38-19133ec93b03	E4dpLzowVwid9UHaNI5IS8gc2uCH-J5I255ka0H2nEA=	2025-10-16 01:27:30.771974	2025-10-09 01:27:30.772917	2025-10-09 14:10:20.565668	t
aaa37ab2-a8cb-459e-a261-455f486fb75a	748320dd-5b5e-4434-ad7e-8a405fd6266f	Zjmr1BSPRXlue0ULlhpnAGdtS7DstI8PtQ5k8LXXYno=	2025-10-16 14:15:48.273479	2025-10-09 14:15:48.273933	\N	f
78e425d5-f831-4149-b5ba-fb65a290375e	748320dd-5b5e-4434-ad7e-8a405fd6266f	qg4omLfo2rbEjw_dP-ZO3-x9VZ3UqAYmpo5rLSf6yPk=	2025-10-16 14:19:28.392882	2025-10-09 14:19:28.393375	\N	f
79364cf4-fd03-4e5e-8294-b9444a4161d4	748320dd-5b5e-4434-ad7e-8a405fd6266f	RpQ2YIXTjKnABC86VUXoMyMx1bte7nYHRQDyZ6QTzc4=	2025-10-16 14:23:16.461402	2025-10-09 14:23:16.461845	\N	f
0ddc55d8-22fd-4f7a-a466-bb712687c787	748320dd-5b5e-4434-ad7e-8a405fd6266f	r8EykJmAW5NS9wSmAOHit1ipENrr2KfR6Z9t6lcyzeQ=	2025-10-16 14:25:38.177469	2025-10-09 14:25:38.178003	\N	f
82c9661f-944a-406d-b5b6-9b164aca6a69	7157f3f6-da85-4058-9d38-19133ec93b03	mGWucDtM9QcSYRNqDLr0JQyCernBmcAuKFLdgycSSkU=	2025-10-16 14:25:49.105812	2025-10-09 14:25:49.106268	\N	f
3da4ecc9-e9ff-4d27-9c75-c44e39527e88	7157f3f6-da85-4058-9d38-19133ec93b03	yf0239pW7E-gBFVvnQ5tbS4VOFyN3E0uskhmrhILKbM=	2025-10-16 14:26:12.692662	2025-10-09 14:26:12.69312	\N	f
cd91182a-8711-4f20-9b61-a03f1079baf6	7157f3f6-da85-4058-9d38-19133ec93b03	lDA13V6K6gZFNUdVn2oAsQWVzhO6gkVDFpTj1gf98hg=	2025-10-16 14:26:31.697682	2025-10-09 14:26:31.700974	\N	f
0432b0b2-c229-4ab0-b232-1cdc27876a9f	7157f3f6-da85-4058-9d38-19133ec93b03	PdEF57nyQjkwtrSUz-wvpn6N3PW08CNxLf_C6KbQfl0=	2025-10-16 14:37:16.576767	2025-10-09 14:37:16.577313	\N	f
07127ef1-d1df-4e84-9b4d-4644abbd6a65	7157f3f6-da85-4058-9d38-19133ec93b03	yn2Yj8tcwjneznIej17fCNWy2ShDZfdjfVsyakWNei0=	2025-10-16 15:17:01.49328	2025-10-09 15:17:01.496846	\N	f
88910567-1876-4763-bd2f-d2d74e3367ad	7157f3f6-da85-4058-9d38-19133ec93b03	SyLFc4NC26WaYhDCJfSbd_MXTHOlQcTohS6FvFrs5eQ=	2025-10-16 15:32:11.752388	2025-10-09 15:32:11.753517	\N	f
9936abc5-d999-470b-b704-c21dd9d2f086	7157f3f6-da85-4058-9d38-19133ec93b03	cJjbn1OaAvXYdfN-302hLXq5Vznxnmm6zQrpWu5AZY4=	2025-10-16 16:00:09.942312	2025-10-09 16:00:09.942768	\N	f
192a4aad-1070-4ea8-a42b-c81c18c89b87	748320dd-5b5e-4434-ad7e-8a405fd6266f	-KMBTeUkd_snAEibi3WHA4tioWb2Nu5_2dZ3nnegxrA=	2025-10-16 14:36:02.772809	2025-10-09 14:36:02.773332	2025-10-09 16:41:27.117766	t
c56f10a0-7067-4581-89e2-20e03a8f43e6	748320dd-5b5e-4434-ad7e-8a405fd6266f	OxkDtPholIpJ0exHJsPV4TYaa0TgM4wdwZm6xq7qvl8=	2025-10-16 20:59:54.323615	2025-10-09 20:59:54.32588	\N	f
dec6bcd8-e8ac-4904-bac6-fa34aaf86f4c	748320dd-5b5e-4434-ad7e-8a405fd6266f	cpjOD-evyjLje7XKo7i2AYb7TGMvYknB9Qo6D_UozUs=	2025-10-16 16:41:27.113421	2025-10-09 16:41:27.11396	2025-10-09 20:59:54.328789	t
f013c4ae-9d22-410a-a205-78bbb4533b9e	7157f3f6-da85-4058-9d38-19133ec93b03	Iy8WIS1MQRxJ52NjOrFl4eaJkHw22V5YZDPubSgu3so=	2025-10-17 03:47:33.622272	2025-10-10 03:47:33.622689	\N	f
3af3be91-420d-4d65-8422-a76da00db95c	7157f3f6-da85-4058-9d38-19133ec93b03	Kum6VNtgApBLpQFsLlNp32g042qepKbhck_2jD6D9LI=	2025-10-17 18:56:02.118006	2025-10-10 18:56:02.118832	\N	f
1455a1ee-071c-4016-b631-8cd39d61e8e0	7157f3f6-da85-4058-9d38-19133ec93b03	rAsd4a3wbIuGVfgtb4byqjqSGwt0ifS3wtazNXoZ3n4=	2025-10-17 21:29:13.63884	2025-10-10 21:29:13.639125	\N	f
606fbeb8-c716-40de-8957-d662eedbee11	7157f3f6-da85-4058-9d38-19133ec93b03	2BkXYG4uUMSldV1MfcCqshq3sdqVfElzbxVExAcMZpA=	2025-10-18 10:39:36.316676	2025-10-11 10:39:36.317295	\N	f
eeb38829-0a76-4557-9041-ac4dd180209a	7157f3f6-da85-4058-9d38-19133ec93b03	_U3MocJeG_hFufSS4_2VxqTyLd9Xh-lPjKQ26RZ2r48=	2025-10-18 16:18:58.37912	2025-10-11 16:18:58.379686	\N	f
360e236d-8291-4ca8-9a84-cfa6294b6a7b	7157f3f6-da85-4058-9d38-19133ec93b03	-JDdfIsKxZ50T4Ry5BJPvAGLy2FKt-eLdpEE4kUk6Zo=	2025-10-18 15:57:11.889617	2025-10-11 15:57:11.889992	2025-10-11 16:18:58.382446	t
efedf9d5-f106-4c11-9f77-3ca92a886683	7157f3f6-da85-4058-9d38-19133ec93b03	2W61s7uSp9fAi1ApxolpkyFEB4aR0XA3XcRrObsmeG4=	2025-10-18 16:19:06.596006	2025-10-11 16:19:06.596408	\N	f
0f4e2312-6f84-4620-a26d-c512f6ecf1ba	7157f3f6-da85-4058-9d38-19133ec93b03	3tm5gFTtS9DeLjiXLHbIsyDDcP5vKpYZnzxHOu5P3o4=	2025-10-18 17:07:13.939979	2025-10-11 17:07:13.941299	\N	f
2fc39d0f-ddd0-4069-9057-ca4df32983fb	7157f3f6-da85-4058-9d38-19133ec93b03	1FYvHkEmUDNAJjg4EXVolRHyLjMEzOOoxO_JV5PAmqw=	2025-10-18 16:45:16.856993	2025-10-11 16:45:16.857677	2025-10-11 17:07:13.945387	t
86165ead-cdf8-477e-8d0d-cfbed9724ed2	7157f3f6-da85-4058-9d38-19133ec93b03	DRcZ4nR9H7qUiXzjYaBlUEGQk0Zmxeiz7uSyPC0LZRE=	2025-10-18 17:28:56.927044	2025-10-11 17:28:56.927344	\N	f
ebf7a9e0-edb9-4fb6-954f-8b3aedc42959	7157f3f6-da85-4058-9d38-19133ec93b03	-6Ok5HyW9Bj2EEJZM8VNjV-NmNz5Yv7D4zMGaiG42uI=	2025-10-18 17:07:13.941481	2025-10-11 17:07:13.942047	2025-10-11 17:28:56.929265	t
1e005f48-3df5-498c-a898-e73004adfbde	7157f3f6-da85-4058-9d38-19133ec93b03	2Qg5eb3o8uiI_zMPXH8WLjzp-jlOWU_fWM3UdKHu3P4=	2025-10-18 17:29:03.852157	2025-10-11 17:29:03.852472	\N	f
18b2dcd5-dfc2-4e2c-91f2-7d762076d1f2	7157f3f6-da85-4058-9d38-19133ec93b03	e70QIiVNadil1s3JsWRTL2jvZAzFEwjfurOc7Yhmsw4=	2025-10-18 17:32:25.056596	2025-10-11 17:32:25.056886	\N	f
004f4b77-ce6b-4cd8-a141-7c5a7767f029	7157f3f6-da85-4058-9d38-19133ec93b03	JpKGTtxIGubSyRDzgTfztWKKjMvgQojtKADiVUrU8UQ=	2025-10-18 19:06:57.048863	2025-10-11 19:06:57.049214	\N	f
196530db-36ed-4527-8f71-fb2559d5db55	7157f3f6-da85-4058-9d38-19133ec93b03	os7PrhynsIkGgo48V4Z9bbB-uWHEi1t8ZNrCn_ULxxI=	2025-10-18 18:45:30.061276	2025-10-11 18:45:30.061518	2025-10-11 19:06:57.050974	t
64afc5b7-6388-4b9a-acbd-ac4f5d317796	7157f3f6-da85-4058-9d38-19133ec93b03	pKbMg2IjNdtXCB-ZePJAe-d6LejhK_FylsEN4oIm_Gg=	2025-10-18 19:07:04.049065	2025-10-11 19:07:04.049393	\N	f
0d41e9d1-2715-4b86-911b-aed9d1c7a12e	7157f3f6-da85-4058-9d38-19133ec93b03	OrJQuD1eAf7IpcxyMGaLK5-A-lZuVSlcZ1okXAWVvZU=	2025-10-18 19:14:35.272876	2025-10-11 19:14:35.273186	\N	f
a3343be3-4a73-46d8-82e5-8d71d415161e	7157f3f6-da85-4058-9d38-19133ec93b03	rESgRyiKiJ-5x3S4C_dcTXE5_cJUFk7z82OPcfjzMKo=	2025-10-19 21:11:53.606899	2025-10-12 21:11:53.607199	\N	f
a9c9ee99-692b-4a10-8b08-97410b96dd36	7157f3f6-da85-4058-9d38-19133ec93b03	fpL_QmMYYsRUoTc3P2do6iUsBE-bXTUQk_BTf0T00VA=	2025-10-19 20:45:32.347809	2025-10-12 20:45:32.348253	2025-10-12 21:11:53.610568	t
9b3c536b-d499-4f85-b407-edbf0362ff96	7157f3f6-da85-4058-9d38-19133ec93b03	x152hCMakJUbzlbZ-1yYypliVjZ4-2D58gyYdbG9Mho=	2025-10-19 20:43:35.473455	2025-10-12 20:43:35.474019	2025-10-12 21:48:12.750929	t
d5c14d8a-0234-459c-a8cc-3e28d06a6ff2	7157f3f6-da85-4058-9d38-19133ec93b03	aj2dzO0xs16O61MLbJz7JcsugVq0vK2eOf3yV0vqIPk=	2025-10-19 22:02:51.428517	2025-10-12 22:02:51.428836	\N	f
3e63505e-343e-42a9-8a98-174c8c253ca5	7157f3f6-da85-4058-9d38-19133ec93b03	s3wS0fVlp5_eWucbbc5U2jz8kPpeNGmNDnV1HH3i1dM=	2025-10-19 21:48:12.748317	2025-10-12 21:48:12.748786	2025-10-12 22:13:48.210923	t
cd7fe467-241f-4629-ab99-718d346064aa	7157f3f6-da85-4058-9d38-19133ec93b03	tsX0DUUvQHnaqI-v2j8PrIRS2YR80SRlVzuqXORoTqc=	2025-10-16 10:12:39.032974	2025-10-09 10:12:39.033368	2025-10-12 23:59:52.429844	t
798dc42f-3a22-482e-80a1-f383f8b94451	7157f3f6-da85-4058-9d38-19133ec93b03	98Fyfp7T-EDcGegToG9QGKQrWA_4HEQwi3g31zwI7JM=	2025-10-19 22:02:51.431551	2025-10-12 22:02:51.431829	\N	f
af141f30-26f2-46af-b668-f136de03cd93	7157f3f6-da85-4058-9d38-19133ec93b03	GDJWkxdVvmxjw3rbxtXLkkCmdTB2xZkMjyZDWTtNEz4=	2025-10-19 21:12:01.208578	2025-10-12 21:12:01.208893	2025-10-12 22:02:51.433529	t
241acc65-d2e5-44a5-b07c-04c3849ec848	7157f3f6-da85-4058-9d38-19133ec93b03	pc2In7hYIRQ32eIdqoewr2fFg5Cu6eACQRFSvxY6FzA=	2025-10-19 22:43:57.89943	2025-10-12 22:43:57.899933	\N	f
f856436f-8ffb-4a2a-a99d-0ff757ebd566	7157f3f6-da85-4058-9d38-19133ec93b03	uTi1CaJzrLL389aKGCqBKfLrmjK_vOt10xEPFBGPRwU=	2025-10-19 22:21:42.648863	2025-10-12 22:21:42.649193	2025-10-12 22:43:57.902273	t
af675638-e826-4184-9206-0ccd9c7ecb27	7157f3f6-da85-4058-9d38-19133ec93b03	nDprsKmFdFZygReNJ2Qn-Wj7Qysh5TOQ9JlwqpMnohc=	2025-10-19 22:44:05.495932	2025-10-12 22:44:05.496438	\N	f
dccf14c8-622b-48a5-8461-6cbc11bcad5f	7157f3f6-da85-4058-9d38-19133ec93b03	9rD1nzJGTV2Dxx0MaslnYt3ssqyLkLMRTm2aEynG-uU=	2025-10-20 13:44:39.449796	2025-10-13 13:44:39.450267	\N	f
172ba9b7-0422-40ef-93cd-c3ee831d7b4d	7157f3f6-da85-4058-9d38-19133ec93b03	9DehOQwbbM3nYRrAABEaiXX-Tx1S_6aezLGEdYg2vhI=	2025-10-19 23:59:52.427083	2025-10-12 23:59:52.427549	2025-10-13 13:44:39.452576	t
652bc03a-5ae5-4b7b-a39d-87755b8654fe	748320dd-5b5e-4434-ad7e-8a405fd6266f	jck1iEyB1_ajtVNUzyPF_cJshcCG_Yps--LC5HSS4Fc=	2025-10-20 14:55:48.491094	2025-10-13 14:55:48.491523	\N	f
a2a6ea8e-4736-48d1-a1c3-97a84eb225d3	7157f3f6-da85-4058-9d38-19133ec93b03	-iwEL71BswLde55WhrSXc4GxQXFWfekH3ftUpJqGaO4=	2025-10-20 16:39:57.277856	2025-10-13 16:39:57.28078	\N	f
8aeba7b4-778f-4428-a904-173d692a9b43	7157f3f6-da85-4058-9d38-19133ec93b03	mYjwo6t8THui-54F507pj-PQ-iuqpaIdI1Cb1EjEVnc=	2025-10-20 16:49:56.425021	2025-10-13 16:49:56.425463	\N	f
95177e31-c8fc-4b7c-a425-67683e348d24	7157f3f6-da85-4058-9d38-19133ec93b03	7kn0doAPszUoYUpmBqYCR_gaj00Tw2hmIEt0iOE7Llg=	2025-10-19 22:13:48.208936	2025-10-12 22:13:48.209324	2025-10-14 12:27:54.936392	t
671b9775-03ed-42db-be85-a0e01dedabe5	7157f3f6-da85-4058-9d38-19133ec93b03	om_ncMwaoAuqXO9AOAALjWqrgNX3kKI-09jLjfsKKxg=	2025-10-22 19:07:23.812728	2025-10-15 19:07:23.815388	\N	f
b190a328-787d-42ef-9a36-505a016d8086	7157f3f6-da85-4058-9d38-19133ec93b03	X59bN4ZuXIKz2EorqCUw6r8mv1YUCmBNe_hrE70aLvc=	2025-10-21 12:27:54.920076	2025-10-14 12:27:54.931899	2025-10-15 19:07:23.81811	t
45bd49bf-f0f6-44bb-82d6-585663ef92d2	748320dd-5b5e-4434-ad7e-8a405fd6266f	_l7sjJeoW2h_5a_Hfey46tMEki2YRGaruq6p0HqR324=	2025-10-31 21:20:12.442633	2025-10-24 21:20:12.442983	\N	f
76bae90c-d86b-456b-9802-83bfb410d65b	7157f3f6-da85-4058-9d38-19133ec93b03	9F7ZHqlEYVSX_SLOnQxe5QC1FC2JHgn8GCZGkDh9Ipk=	2025-10-31 21:23:48.768239	2025-10-24 21:23:48.768484	\N	f
ab4b67fd-6be0-4f87-ab12-02c95c4c6ad5	7157f3f6-da85-4058-9d38-19133ec93b03	uwom6YDw5KFApvRxrA-_Tba_SkktBNId9LUbdKANvqE=	2025-10-31 21:53:26.607648	2025-10-24 21:53:26.608128	\N	f
2e188e17-340a-46e4-bd24-4ab3c578e5e2	7157f3f6-da85-4058-9d38-19133ec93b03	m_aa5xYd4xx-a-0qWMeaT_az5VfbVIPQdUE4KQnXymw=	2025-11-01 11:14:03.728394	2025-10-25 11:14:03.729053	\N	f
31286e9b-d5c8-482e-9696-b4feddd603e5	7157f3f6-da85-4058-9d38-19133ec93b03	_d5lADscKWtXnmZ2_MOH0JHWwj8Isp1jaHj-qNFM80o=	2025-11-01 11:32:42.419178	2025-10-25 11:32:42.419546	\N	f
69ac0f0c-4bb7-498f-86d3-6349a52e00fa	7157f3f6-da85-4058-9d38-19133ec93b03	9mOSN9MSXBvSEpssMXr7CqFeMbFbpPgjnuGn11qDrqY=	2025-11-01 11:40:36.207036	2025-10-25 11:40:36.207389	\N	f
3172aace-8ec1-4080-9990-9e02080c2ecc	748320dd-5b5e-4434-ad7e-8a405fd6266f	GmyMDW_349xtedRoAXQwzmiBAOxgjMlsUPZf4-jvVkQ=	2025-11-01 11:44:13.010141	2025-10-25 11:44:13.010497	\N	f
ba90fbde-6674-406e-9c2b-d187cb974709	8f333073-8d5a-4202-9627-02863000b822	cVk1ZWgwUw_foxG9Jc2N7D1bkeCyOjE6lYwf0g6ew_c=	2025-11-01 12:10:17.980268	2025-10-25 12:10:17.980699	\N	f
3c854533-0788-46b8-b3f8-373c178d954a	8f333073-8d5a-4202-9627-02863000b822	9IEu9xPWMHphxAP0TYCTXnZVQ6wvtn3FbWgecalgY8g=	2025-11-01 12:22:14.654747	2025-10-25 12:22:14.655191	\N	f
5b393654-48fc-45c5-87e4-9da80bd63592	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	CT_pmvRnE_7hjJYbIWWRhYtfTtd968IB8ReHVImw2Zg=	2025-11-01 12:25:31.858025	2025-10-25 12:25:31.858393	\N	f
d296d0ce-a8bb-45d1-a4c2-75807c03a453	8f333073-8d5a-4202-9627-02863000b822	Cgm_jtv21LB568aft89nuJZm9B95KYdPebVMvqPGSF0=	2025-11-01 12:37:08.997207	2025-10-25 12:37:08.997485	\N	f
9202d302-11d5-4e4c-ad59-dfd7ebf40835	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	JT0rEDaRimoNcVtyPJWzbZFnVZcH2y9AMt4cA-o51yc=	2025-11-01 12:37:32.152202	2025-10-25 12:37:32.152455	\N	f
c78d98bf-9e64-4152-b9ef-b1052a24f86f	8f333073-8d5a-4202-9627-02863000b822	Lac2Nj2pj48JedHhyrx7R3k_l1pdhKZ8gbjoFsPnmDQ=	2025-11-01 15:00:11.754042	2025-10-25 15:00:11.754312	\N	f
80f58ec7-6c9d-411f-b8f3-8f2a33362afd	8f333073-8d5a-4202-9627-02863000b822	0_bQMIPtNbJqik3afSeHyeD1J5rCqzKPmsiQScTXiBA=	2025-11-01 15:07:16.381046	2025-10-25 15:07:16.38135	\N	f
bd22b91c-01dc-4cc3-b02b-39b7881a93a0	8f333073-8d5a-4202-9627-02863000b822	phuARrRjiLF5A3WUBNkHKt-t8YW0sonPfyNLkKj-vpA=	2025-11-01 15:09:04.038525	2025-10-25 15:09:04.038873	\N	f
57cd9150-dcae-48df-93a9-4c3d6d2bf023	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	S1QWAvYm4X_kLI6hvc_3PfRlLntpx3rPJKw86U1tZPo=	2025-11-01 15:16:18.060555	2025-10-25 15:16:18.060812	\N	f
e1dee871-bd75-4fb4-b4a9-201408ea5a9b	7157f3f6-da85-4058-9d38-19133ec93b03	t3CPOaslYnLDZkH-ZLrYRLpqQc4tNRnQuEkAAnEUa_c=	2025-11-01 15:17:18.559795	2025-10-25 15:17:18.560092	\N	f
10bfb46e-47c0-4f93-9835-faabc87a7443	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	3AzE5AMe4hnHxZDox1OD9qJa0eX8B-Xm3BEaGuAJzy4=	2025-11-01 15:19:36.56429	2025-10-25 15:19:36.56452	\N	f
7600f87e-ba39-417f-8667-b0c946d22f68	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	iBb1nFbK4JwjGGeRbnfqkHrItd2T9UAd8PCbWtKw-ag=	2025-11-01 15:54:42.371233	2025-10-25 15:54:42.37171	\N	f
dfeb7945-9221-4d11-865b-c523ba470d10	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	qZm9AHsWNrb5dFPhYEL9uMTGKV4lONPddTXGHhRPp2g=	2025-11-01 15:57:14.891254	2025-10-25 15:57:14.891582	\N	f
4cf1b810-eae1-4e69-89bc-3a56b97107dc	7157f3f6-da85-4058-9d38-19133ec93b03	CEmernOwDUIBWkaqNnk2TeELCzEnjuWl-mbWGAWiq5Y=	2025-11-01 16:00:13.874265	2025-10-25 16:00:13.874602	\N	f
a618e698-2b09-4355-ad0c-eb7015ff524e	7157f3f6-da85-4058-9d38-19133ec93b03	EaL8Yhgn6eOqbwYN9D06EXOp8qQHEoL-wKKqCQZ_8Tw=	2025-11-01 16:02:07.472964	2025-10-25 16:02:07.473361	\N	f
026615e9-3f64-4fe5-9e36-5598210517f4	7157f3f6-da85-4058-9d38-19133ec93b03	ISpcoC9W-KnfTRzIKnpJgSzvVh7u5xQtWVA5OyDkKmg=	2025-11-01 16:04:13.486283	2025-10-25 16:04:13.486652	\N	f
38de5d03-3202-4557-904d-d151c479d358	7157f3f6-da85-4058-9d38-19133ec93b03	tQze8ktPTLNh81hxFMag002pPdXiMZe6JHU6excNZ9A=	2025-11-01 16:07:26.790167	2025-10-25 16:07:26.790482	\N	f
307d8f9f-c72f-4455-8a95-5685b65adaf7	7157f3f6-da85-4058-9d38-19133ec93b03	TIFMWbGNI-QJHmuf5LuVVxLOutf3UgoLP1AS6A_EHP4=	2025-11-01 16:19:00.402096	2025-10-25 16:19:00.402418	\N	f
f1ec4564-204b-42ba-b4f0-a363c49f2f2a	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	f0sYMOqO2LxRRitzigzIZlE20cm8eBTqRzfLw4b8N3g=	2025-11-01 16:19:38.507968	2025-10-25 16:19:38.508369	\N	f
fd26d3fc-7f8d-4d76-af28-6c818c5196e9	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	Sqy6Wyas8lGps9W6G6py967YWQ5R9pIK117XhZy-NbQ=	2025-11-01 16:21:05.348787	2025-10-25 16:21:05.349291	\N	f
7a48eb83-3424-40e4-8ca5-08bf49be6343	7157f3f6-da85-4058-9d38-19133ec93b03	N1BY4oT4BV9PkT6lGKbuuRRYU53ypeKzno-TMK9gsis=	2025-11-01 16:22:18.780668	2025-10-25 16:22:18.781164	\N	f
9c29bf7f-b2be-4ff4-a82c-f7c8442a05b0	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	NoKlwrDc_QkgKMPloXWVGuhfciJTOp8azgoQuFXYblg=	2025-11-01 17:07:41.669946	2025-10-25 17:07:41.670228	\N	f
4b206ea3-ac14-4d4a-baed-0e1b689a47c0	7157f3f6-da85-4058-9d38-19133ec93b03	xpI5tAFLI3qTMCweQ4tPJ1XKoT27AJVWS2JQpV465_w=	2025-11-01 17:21:05.751342	2025-10-25 17:21:05.751627	\N	f
ff31ee64-c811-4fec-a3f0-4e3fceb088ae	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	ZMqEquyqMXYIUB-m4uKt0e1S3nW2XmWr9tN99yZwkuQ=	2025-11-01 17:23:17.444262	2025-10-25 17:23:17.444608	\N	f
3ee12d0e-6a13-4eca-9a76-a177d46bde89	a76943e8-9673-41c0-8e89-b947294882d7	k3CWnYj-MyOXWd8TMCAVzheQeWIG59qXs1mEq1leLNg=	2025-11-01 17:31:42.979027	2025-10-25 17:31:42.979316	\N	f
3fcb4eb1-6a3d-4d06-83c1-4b843ea7f057	cb74dac2-bd74-460f-a83f-04fe02cd9363	XXVphCKKYd9BMlE04MT5rkbqFWd5kJ_kyJO4qHSEjoo=	2025-11-01 17:33:42.843532	2025-10-25 17:33:42.843825	\N	f
32abe936-e24e-46f0-a0d4-ab3fc96bcde6	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	ixyQRYHMDVvzTwclljbpEyxPdSA7RkyCCq217oJJHZ0=	2025-11-01 18:19:39.043445	2025-10-25 18:19:39.043756	\N	f
0b189184-6b18-44b4-ade4-3a22a2ab5795	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	g0rSU9f7Wu8eQYIJqJdaFUMdYXyClQR_xznggato20U=	2025-11-01 19:07:35.04417	2025-10-25 19:07:35.044452	\N	f
28b7a6fd-62b6-45cb-9bc2-3a20ff4584af	7157f3f6-da85-4058-9d38-19133ec93b03	Is3_OdbQy7O5hC9h5CxAj2Goj6owIi2oco9--eTyQao=	2025-11-01 19:11:13.555758	2025-10-25 19:11:13.556028	\N	f
18c2b9f0-ca66-40d8-b3b6-1dc5ff183eb7	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	TuSAdZONaEfh47IoB-P2-HGNLnJhHCuvbTBnmOhb_i8=	2025-11-01 19:11:57.249841	2025-10-25 19:11:57.250171	\N	f
c8156b16-c123-4e8c-beae-69efc72c0f8a	7157f3f6-da85-4058-9d38-19133ec93b03	U2t46xykAEieudVX1qvgwSoIRDPX_8Ga03Ogx9d0trY=	2025-11-01 19:14:58.842614	2025-10-25 19:14:58.842918	\N	f
6987e386-da05-4b4b-83df-c45a77c9d134	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	bZAyi253EZ85yR7nJ8WGUjmSlQ6tpKIEo7I0YWbh8wU=	2025-11-01 19:22:43.046021	2025-10-25 19:22:43.046609	\N	f
8dca74c9-712a-4601-bf6a-ddfc6f9f2d7e	7157f3f6-da85-4058-9d38-19133ec93b03	QEVd67SSpXWvgOilGQYNZcMh_QD6WYZxvrB8LH5lgJk=	2025-11-01 19:23:46.351826	2025-10-25 19:23:46.352144	\N	f
95f11338-0985-49a9-9401-35d7a9737aee	7157f3f6-da85-4058-9d38-19133ec93b03	dUYPIqRcvSUnlZNZzkKA3K5jiJd9fBukldCXcTqy0i4=	2025-11-01 19:26:21.645181	2025-10-25 19:26:21.645519	\N	f
a802151e-0dfa-4966-95e6-4055e19c03a9	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	dxY-F7ziJcdCO9VQtNz2KyU5O6PfKtG6uwXX9PRoxnE=	2025-11-01 19:27:57.658759	2025-10-25 19:27:57.659077	\N	f
7cc6ad7e-1e13-4228-a766-99f32ff767e5	7157f3f6-da85-4058-9d38-19133ec93b03	T1jGiDROQ-73X83wuWw9oELYJIEaehno3Pln40sMngg=	2025-11-01 19:30:08.543368	2025-10-25 19:30:08.543697	\N	f
3d3f1ffb-6e16-4b85-9653-f1a1a1521a9d	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	dHn_Xi3-LNn9mjVrvPO6tNjWWkqSO-QHoySvWjE3Dkk=	2025-11-01 19:47:35.68396	2025-10-25 19:47:35.68474	\N	f
57aa0ce4-038f-4477-8db2-9b8a758007ef	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	BEHxlNSBXhMppggKitA0u5IzouFQoQeSm-BLEIO94jc=	2025-11-01 19:49:30.843805	2025-10-25 19:49:30.844249	\N	f
214f7f32-9bb1-4027-9102-242a371cc13e	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	oepNO_Tv1DtR-SK_T9TPdQXNJO2fqEQTr2oKSlwSxI0=	2025-11-01 19:50:41.282894	2025-10-25 19:50:41.283401	\N	f
438a2e09-28bc-43fe-9b26-736d6b931cdd	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	3Xxvnz0OF-KfwKJWd8TSeZmtXsXlF4y6HUvBUbiNKIQ=	2025-11-01 20:20:49.833073	2025-10-25 20:20:49.833366	\N	f
9ba8f58f-28f6-4101-8134-970186ad0b82	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	D0obo4gLXKJIpGtbb68cPwr2GlcHjb1Lr9bxy_QJuSs=	2025-11-01 20:23:07.538018	2025-10-25 20:23:07.538342	\N	f
73ca621d-2e87-47ad-bbe8-d5a66be5f8fd	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	VTZBpzKxF8bWCozd4j00g9DGOLBWyKuGPpe6cEa3m9U=	2025-11-01 20:39:07.969162	2025-10-25 20:39:07.969467	\N	f
fbfa3460-2880-490e-9b4d-616d39004b7e	7157f3f6-da85-4058-9d38-19133ec93b03	JX9Ok6IBzzupzjFLyM6IaiTBR_cR9K5lixPKrfokF68=	2025-11-01 20:39:42.943803	2025-10-25 20:39:42.944086	\N	f
4014a17c-ec9a-4ce0-bad7-1a7f96df2ec9	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	xyh1tIsJL3ImRI2jpDX2eGxQ7tlWQBRm5kn0LEZOcuY=	2025-11-01 20:59:51.459008	2025-10-25 20:59:51.459257	\N	f
103c9b6f-b13f-4931-955f-d2c9f86a4526	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	P91lPiCLqbWfDMWrzXMxO8p9DN1p86QypeQW5fRseug=	2025-11-01 21:03:10.44639	2025-10-25 21:03:10.446692	\N	f
f324582c-6d1b-4080-a92a-2b6cf88e4765	7157f3f6-da85-4058-9d38-19133ec93b03	AF17HfHxo8myc7H6tseo_N8m0tuZgiiaTnpVKB5xgkk=	2025-11-01 21:27:51.132946	2025-10-25 21:27:51.133207	\N	f
43f91005-7b91-460a-b21c-6f8fe3bb5efc	7157f3f6-da85-4058-9d38-19133ec93b03	b8ryoj2kAjO6NP8ba9Fzh_WkS7HR_tpZrAFdTGcxYB0=	2025-11-01 21:29:12.643345	2025-10-25 21:29:12.643641	\N	f
54762781-6ce5-4c4e-8e59-f4c1bdbfa7ff	7157f3f6-da85-4058-9d38-19133ec93b03	ea9eGE5_jAmcuTs37NLO5emR6kz8fL5_63uCTwleX7g=	2025-11-01 21:30:35.236052	2025-10-25 21:30:35.236381	\N	f
690c3370-f60b-4b0c-8127-1558ad7fd1d1	7157f3f6-da85-4058-9d38-19133ec93b03	7EYcGIUb-LEEjdqGe7ot2x0R2tBj265zO1NbLNYrfMw=	2025-11-01 21:46:12.459581	2025-10-25 21:46:12.46363	\N	f
99487a39-56f3-4865-8241-4dcf7f770b03	7157f3f6-da85-4058-9d38-19133ec93b03	I08fGK2f4uWClvulyvho4k3pjZ98LvS9t-0sGwiwViU=	2025-11-01 21:46:44.349106	2025-10-25 21:46:44.349405	\N	f
13d8a128-a3cb-4360-93b5-d3c948fbbf95	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	i5nxuj5qUrnkAlE8iDCKku7fjqAqPGXD-Jc9qSO3BVM=	2025-11-01 21:47:36.371421	2025-10-25 21:47:36.371721	\N	f
1a556923-80fa-47ed-ba59-91008ab87de9	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	9RI4WyuXcMqIO4-H1y2a89890an9JM9RhzejZ5nKUoE=	2025-11-01 22:39:11.61161	2025-10-25 22:39:11.611947	\N	f
d6e7a325-cabf-4c37-912d-977bc824ffef	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	P1lJcedeJTu6_vPmpscXEXYL82k9l7tZGQAYaX-T7Cs=	2025-11-01 22:58:09.374451	2025-10-25 22:58:09.378562	\N	f
394df2e7-97ac-4148-b50a-d6d24a3f9e52	7157f3f6-da85-4058-9d38-19133ec93b03	o39Xyv5P8x0gqCy6oNHcjgDHCqF_kqcrbhPJrgpVdX4=	2025-11-02 00:42:48.748532	2025-10-26 00:42:48.748814	\N	f
3c908539-5e3f-49fb-9c58-cabb9618b09e	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	pv56bOTyd-PPOQguX-9PBAIjJwIfytFVykDmpvM3FxM=	2025-11-02 00:43:51.061112	2025-10-26 00:43:51.061422	\N	f
2d93273e-a8cc-417c-9d68-98408d25d7a1	7157f3f6-da85-4058-9d38-19133ec93b03	U-zRR6lx6tpqI9dEILW9br9AGfTa9yE8oc6lrCBbZZA=	2025-11-02 00:48:01.075738	2025-10-26 00:48:01.0761	\N	f
2c3fc2fb-0a29-4ac4-a87c-ee622c7036c9	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	26FiGwF8Eb81vCNAI39kMthQnbvkYrKqWLnjRwZvgNk=	2025-11-02 00:49:20.751475	2025-10-26 00:49:20.751801	\N	f
08b32b79-1719-4019-9e44-e622528d64be	748320dd-5b5e-4434-ad7e-8a405fd6266f	X4FTJowLmA-QidIDGch886jxEX6VRev2syYybZLbDRo=	2025-11-02 10:08:00.151273	2025-10-26 10:08:00.152356	\N	f
590259e6-0b1d-4922-9299-f2c6af0f91a2	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	RsmONenzyOJbKHfquoKZ_pL9iaVI9BZ-19zJ0DAu7ws=	2025-11-02 10:11:05.064008	2025-10-26 10:11:05.064272	\N	f
f0e9c2d9-75c1-4e69-ba77-9955b3e7de01	748320dd-5b5e-4434-ad7e-8a405fd6266f	zzgcO0HdG6YEmwlsLUC1KEW_uwAsgCy-bttFwlS-TpQ=	2025-11-03 10:07:24.257932	2025-10-27 10:07:24.258268	\N	f
4f4315df-7af4-4956-a9a5-aa1771b4df54	8f333073-8d5a-4202-9627-02863000b822	eBy38ZLLYES3OTQADdNpSkYCSlxXBKk-oUJb3CNpmq8=	2025-11-03 10:07:55.149622	2025-10-27 10:07:55.14998	\N	f
65bde963-9454-4541-a495-4966f793a9e7	7157f3f6-da85-4058-9d38-19133ec93b03	y3xrDeePWSIolBK9oxB-rDGqlVSeSnGkm6f--9_Azz4=	2025-11-03 10:12:57.955862	2025-10-27 10:12:57.956206	\N	f
48f0aac6-fc98-4418-80f8-a1166c950907	7157f3f6-da85-4058-9d38-19133ec93b03	boshYQpsmQkHsfZcPki96LsUypkwwi9085k1m55Cm0Y=	2025-11-03 11:59:27.947659	2025-10-27 11:59:27.947998	\N	f
0d62987b-1c55-48b9-87c7-3793425c7613	8f333073-8d5a-4202-9627-02863000b822	KE49423MB-Aie6q6y7O3V8avwfe0MiQdr0HBbW6CulI=	2025-11-05 11:58:01.140592	2025-10-29 11:58:01.140925	\N	f
29eed939-f0b0-47fc-a5fa-ce123ea5d9ce	748320dd-5b5e-4434-ad7e-8a405fd6266f	aYZ0yw7CwbSk8B0g_khglMUP39tbNDr_bBnu2YgCXZE=	2025-11-05 20:08:38.970558	2025-10-29 20:08:38.970855	\N	f
e47d6b25-d896-4ff1-90a2-872d571ad33e	7157f3f6-da85-4058-9d38-19133ec93b03	xsCEG06_c8LliUMWPMPEFYWZfghZ94isnTThZUwJyfI=	2025-11-05 20:17:06.648534	2025-10-29 20:17:06.649113	\N	f
e7944512-80eb-4af7-99b1-f6161273131c	8f333073-8d5a-4202-9627-02863000b822	PZ9OEmyUSzpEHzVkiMLT6vKfbLyWTw6b9nXEI7COQ7A=	2025-11-05 20:18:26.745238	2025-10-29 20:18:26.745553	\N	f
7a4af321-a1ec-4f9b-a591-99b4212c9fdd	8f333073-8d5a-4202-9627-02863000b822	ktTKUfIjkihvQ3Oef-vErKUQDIKOlCum_zUcq95WRaQ=	2025-11-05 20:19:49.548523	2025-10-29 20:19:49.548945	\N	f
28693aa5-5ef0-4e04-a4a8-93523bfa5847	8f333073-8d5a-4202-9627-02863000b822	mwpyBg4ilvGbIWVSTx-UvSljrX3CsnCM5C5hGJ_qFUo=	2025-11-05 20:20:33.744414	2025-10-29 20:20:33.744721	\N	f
ec2d247f-83e9-41c6-ac5b-b94a35694448	748320dd-5b5e-4434-ad7e-8a405fd6266f	KdEROFUE8C6BY80d7zG722ysQxeRNUPoHD-MWq6L0vo=	2025-11-06 19:11:48.349903	2025-10-30 19:11:48.350325	\N	f
1b651df4-4090-4d6e-a839-f642d2611304	8f333073-8d5a-4202-9627-02863000b822	t2Yefh_Ki20ZsLvMEIRIQhpNy9i1O7yy-genjCmUVAo=	2025-11-06 19:12:40.746537	2025-10-30 19:12:40.746813	\N	f
f3fe7044-48f6-468a-89bf-f4f3036c547b	748320dd-5b5e-4434-ad7e-8a405fd6266f	3eVYKhokzvZlTS-6O3ujyF7gu73YeQnuskN5R7zaxBs=	2025-11-05 20:08:40.255469	2025-10-29 20:08:40.255827	2025-10-30 20:27:37.849475	t
d36ab99e-6545-47de-a49d-c5b21562e12c	748320dd-5b5e-4434-ad7e-8a405fd6266f	6thpSiPVDgrl_W0kMK_2sIrCcdmF_trlPmfOKdsthaE=	2025-11-07 11:51:44.679756	2025-10-31 11:51:44.711303	\N	f
e07a4db1-d5fb-4ac8-8e94-e287c2aa58b1	748320dd-5b5e-4434-ad7e-8a405fd6266f	Ffh97cMOTi_YMeEFAafb55t8rLStbQBbGRa-YLlos9I=	2025-11-06 20:27:37.753745	2025-10-30 20:27:37.75414	2025-10-31 11:51:44.714252	t
b5867851-5b18-4b7b-ab1f-11fa4351cad5	748320dd-5b5e-4434-ad7e-8a405fd6266f	1lyy411yVh_dUgkB3uNnTX5UhJOouSsM21p_LERxN1g=	2025-11-08 06:47:34.480439	2025-11-01 06:47:34.480807	\N	f
15a4dc01-7890-4f85-baa5-752f00518a47	8f333073-8d5a-4202-9627-02863000b822	Yx9AFIYVqDEgFVGtyG8MYiBNqW2-1DUahFF0jDdcip0=	2025-11-08 11:11:04.496502	2025-11-01 11:11:04.4969	\N	f
55d5d078-5788-490d-bf38-b6483c27188f	8f333073-8d5a-4202-9627-02863000b822	Wa5rNgAx--HipzBmFs3hhlKPQNf8iL-eaJuixZ2afKk=	2025-11-08 18:49:18.476957	2025-11-01 18:49:18.47799	\N	f
cc415958-074d-4e71-bc93-0c8dd364a9a3	7157f3f6-da85-4058-9d38-19133ec93b03	3uXZE5kQ3IdtKr8-kCvnfwEwLuCAr8urSbKn2yfUMsc=	2025-11-08 19:05:17.968709	2025-11-01 19:05:17.969101	\N	f
8ce904f6-0238-4f04-85d5-f352b069f45e	8f333073-8d5a-4202-9627-02863000b822	-atTZMtWuQ1SABM2W2W7dzqggyvZHOyv19r1dlQBccE=	2025-11-08 19:07:19.483198	2025-11-01 19:07:19.483608	\N	f
b47d3c03-52fb-4f6b-a073-0fab448c7def	7157f3f6-da85-4058-9d38-19133ec93b03	q4232sZYHjj4a85NjEULxyzrH0hI1g_RfJHbZfqV8xU=	2025-11-08 19:55:32.392273	2025-11-01 19:55:32.392693	\N	f
81427438-2010-461b-ab34-5a653b5fccd1	7157f3f6-da85-4058-9d38-19133ec93b03	rraZ4o0Upi4KrzuLAK3rvkzXTHAzyt5cp0QtxWffF8U=	2025-11-08 19:31:05.391188	2025-11-01 19:31:05.39159	2025-11-01 19:55:32.395569	t
702061ca-be8a-4ec8-9207-9087e100a082	8f333073-8d5a-4202-9627-02863000b822	xmBzVRemJV707MA0618pEXIeGoLeKW04364PbofoCLY=	2025-11-08 19:55:41.419223	2025-11-01 19:55:41.419565	\N	f
38824aa6-ba3b-4175-800d-8c2f19cc7164	7157f3f6-da85-4058-9d38-19133ec93b03	5crmBG-RAV7lLywM3kGYZALwUPX4DTS33bpoPO_AI-E=	2025-11-08 20:02:44.432674	2025-11-01 20:02:44.433028	\N	f
3cecad02-aa80-4fb1-920d-d330bcd01e2e	7157f3f6-da85-4058-9d38-19133ec93b03	QR-shBYpsm080Zxc5F969eFc87AitT8xjIndsg1jyNo=	2025-11-08 20:29:52.430525	2025-11-01 20:29:52.430882	\N	f
e1923ad6-e066-411d-9b10-27eaa617bf36	7157f3f6-da85-4058-9d38-19133ec93b03	RRxx-NP9c00P69eJQTaZoRwFU4HvfCR7xOzDG-TxWeE=	2025-11-08 21:08:03.64496	2025-11-01 21:08:03.645401	\N	f
9a4e5103-441c-4c28-84a7-c408083d0ccb	7157f3f6-da85-4058-9d38-19133ec93b03	wdYZ47VVho5kxfBdlhNekRhWnWRMBqY42H9thwjROFk=	2025-11-08 21:51:20.095974	2025-11-01 21:51:20.096406	\N	f
0f5e1d64-67d6-49a1-a832-d94a94dd20ba	7157f3f6-da85-4058-9d38-19133ec93b03	4bW-P1c7wI1fUET78Ht3YmkfNMraIbDYlc8M3FJxtrI=	2025-11-08 21:28:10.336871	2025-11-01 21:28:10.337256	2025-11-01 21:51:20.101611	t
9d35e1f8-1ff8-4f5f-b454-e551e58d296b	7157f3f6-da85-4058-9d38-19133ec93b03	Naj0SX0BwMCdTSNgI92w1zhdc-ZWhjrm3HfS5sSGuJ8=	2025-11-08 22:13:00.487659	2025-11-01 22:13:00.488006	\N	f
0826b69b-28b0-45f2-b730-b84e287b7b31	7157f3f6-da85-4058-9d38-19133ec93b03	bLopkaKx5lJWYY4O6LpalT0fnfDFVfXdso8nrs4uwRo=	2025-11-08 21:51:27.516554	2025-11-01 21:51:27.516901	2025-11-01 22:13:00.546477	t
dfb6943a-262d-45ce-a472-4e97c0e04fc8	7157f3f6-da85-4058-9d38-19133ec93b03	1zBeJbeT6HLGV-ijupU4Rs8UOiBsn5reWkscp3APXQM=	2025-11-08 22:15:35.434192	2025-11-01 22:15:35.434636	\N	f
dc82171a-0de2-420a-83c2-1e518bbc6cf5	7157f3f6-da85-4058-9d38-19133ec93b03	8cv1P0D0hFV9HNu5HGBGA345Ynw7S-xxGFi70BWV3ik=	2025-11-08 22:45:51.230828	2025-11-01 22:45:51.231215	\N	f
3eb3a324-e65f-4770-bef6-bd15fe8a1857	7157f3f6-da85-4058-9d38-19133ec93b03	w85B1Bpo_M03hDpgzSpEyW0OyJ2M5fDZYYFRMfHaMTw=	2025-11-08 23:11:47.910536	2025-11-01 23:11:47.910939	\N	f
32cd0475-e5a1-4063-9ccc-f6eaf50dc49e	7157f3f6-da85-4058-9d38-19133ec93b03	rNtzxSorn2dANRNYx1jijYOxf9MBN3dKxyZuJYunJik=	2025-11-08 23:47:58.144016	2025-11-01 23:47:58.144402	\N	f
fe22f3d5-ebae-4617-a6ef-9eeb82501726	7157f3f6-da85-4058-9d38-19133ec93b03	H4tl88gDVRjVpwbUCkV1qp4v5OR1Xxz0n4X5kicuaGs=	2025-11-08 23:27:46.213677	2025-11-01 23:27:46.213902	2025-11-01 23:47:58.146757	t
0738aac4-5517-4ba2-8fe0-79a50639d372	7157f3f6-da85-4058-9d38-19133ec93b03	LzArCMc19WRnshLRHX2KB4CWQUYEuLszoguccqoyrPY=	2025-11-09 00:10:06.219516	2025-11-02 00:10:06.219771	\N	f
834faf85-20dd-4a85-8ec4-2d315e4ac194	7157f3f6-da85-4058-9d38-19133ec93b03	3Zshw5dm48oR8ysJGk4cYydgkq8vLIeEcTlnZGThoDo=	2025-11-08 20:52:07.769006	2025-11-01 20:52:07.769448	2025-11-02 02:04:21.530526	t
3a9d70ea-df88-48ef-9afd-418a2013ec68	7157f3f6-da85-4058-9d38-19133ec93b03	JIcqQ_PAN508g9L40e3YXj5KTOLJCpUFRAZ1WKwsqyE=	2025-11-08 23:48:04.014843	2025-11-01 23:48:04.015079	2025-11-02 00:10:06.221399	t
26a16ce7-73f0-4592-b9a5-f6fd47d95a52	7157f3f6-da85-4058-9d38-19133ec93b03	y11WyU3GatGwGT36rc29VpwT-L4_dLGgrlIkhuIhn84=	2025-11-09 00:10:11.586533	2025-11-02 00:10:11.586781	\N	f
2142e2de-acff-4ce7-9138-d690b82f6d7b	7157f3f6-da85-4058-9d38-19133ec93b03	2S2BSAuznZBRM57G_fTgP_-wzirNGhqNN4o32Wp1ois=	2025-11-09 00:58:43.752694	2025-11-02 00:58:43.753043	\N	f
170dc9da-d745-4f30-b668-736ffbed0c9f	7157f3f6-da85-4058-9d38-19133ec93b03	D2498o5IqdXxwEF1Q8U-Q60bLJbkhChc1gZkw-FSIa4=	2025-11-09 02:04:21.52592	2025-11-02 02:04:21.526829	\N	f
c137978c-ff71-4cbe-b7fb-a834695b24f4	7157f3f6-da85-4058-9d38-19133ec93b03	KVZ-KlArKrfWgOg7XgmJCrgQrnH_fODJxl0xkGD3aXY=	2025-11-09 02:24:09.701937	2025-11-02 02:24:09.702381	\N	f
a08e0bcd-e458-44b5-8774-83038ddce1e1	7157f3f6-da85-4058-9d38-19133ec93b03	Z4DmLeOxAnSxjT4vX3SkDQhoMhZRvYR_upyiy9Fw1jE=	2025-11-09 12:54:08.804568	2025-11-02 12:54:08.804977	\N	f
0ea1f4f5-9559-4b63-973b-32831ee61270	0197cfc3-7ca2-403b-ae4d-32627cd47222	ssmnmwr7OIEOdyZSRy3iR8nlzNqFDDOXN7xgOa3Qa5o=	2025-11-09 12:54:55.799967	2025-11-02 12:54:55.800456	\N	f
dc9746c7-f16b-42ce-bd98-0a38d0f2077f	7157f3f6-da85-4058-9d38-19133ec93b03	WcOMVnGOuauEawghokpxw6H2rQDmUxSLwWuFmfT-DCs=	2025-11-09 12:55:23.02281	2025-11-02 12:55:23.023262	\N	f
2e2a753a-a127-48c3-b53b-e82230ec1f75	8f333073-8d5a-4202-9627-02863000b822	k01JE_MS6r2Jfe1Zd-oPq5sPK8K0uRl9vyjSQ2pOxmQ=	2025-11-09 12:58:24.609546	2025-11-02 12:58:24.609897	\N	f
925ce7b3-a21e-4185-ad5b-3ce3f762ab06	7157f3f6-da85-4058-9d38-19133ec93b03	gmpRfWk10QePpNHuFk-7O7T0EyW34WHVhFO1cgptAPw=	2025-11-09 13:15:17.396051	2025-11-02 13:15:17.396341	\N	f
acc818e7-73af-430e-92b8-7f37188d4faa	7157f3f6-da85-4058-9d38-19133ec93b03	qQURrqm-vOlymxP13mGcIGib2W3Xp1qvhU4Zc_i3woY=	2025-11-09 13:15:50.4122	2025-11-02 13:15:50.413983	\N	f
a70432ee-2882-439b-883e-3ac7c3f787c2	7157f3f6-da85-4058-9d38-19133ec93b03	z5jbZpCDTOSqUtPkTK7YjkSvc8xLdSdlpJHZECp_LM4=	2025-11-09 14:16:18.599049	2025-11-02 14:16:18.599382	\N	f
b1886a53-fbf0-4554-917b-79074bc3f29b	0197cfc3-7ca2-403b-ae4d-32627cd47222	pFH8L-O7LN0jB7K-eIXNFp5DWaVzCn7zxTxQy6_qWpo=	2025-11-09 14:17:09.103657	2025-11-02 14:17:09.104043	\N	f
6b15e929-942d-44eb-af34-590588749653	7157f3f6-da85-4058-9d38-19133ec93b03	tC0xxOOfAyJ-GzFz3-1y8KOCMVmn7hwXBCF7QnfrhRI=	2025-11-09 14:17:42.651453	2025-11-02 14:17:42.651794	\N	f
4e0a43a5-629e-41e7-9ba8-059a859be44e	0197cfc3-7ca2-403b-ae4d-32627cd47222	qZm1QT65H7Q3ByMeSBqhDLJEOlkb8nkZx-fnySS-Nzc=	2025-11-09 14:23:04.022362	2025-11-02 14:23:04.022766	\N	f
a32e3375-db87-466c-bc0b-2d9bc9f714f4	7157f3f6-da85-4058-9d38-19133ec93b03	6HtCDM-rm7XDFo7aZAwuThne02qlBj8ikEAn8BrrtFE=	2025-11-09 14:23:24.185853	2025-11-02 14:23:24.186164	\N	f
4c49ae10-8206-4da2-98b4-a3af91f6a337	0197cfc3-7ca2-403b-ae4d-32627cd47222	7Bg7Q1-95QS2RvG1xD2b8U7eFt3nF5rTVgUZ_yVZuKE=	2025-11-09 14:34:39.517625	2025-11-02 14:34:39.51796	\N	f
90e60b5b-88c2-4509-9634-625473b25fb8	7157f3f6-da85-4058-9d38-19133ec93b03	PWLV2TilUYu7wbX37e7KprXn70jN84-9HO1DOh_w8Jo=	2025-11-09 14:35:54.696311	2025-11-02 14:35:54.696577	\N	f
25739ed8-23ba-4a9f-8543-2890566bec11	0197cfc3-7ca2-403b-ae4d-32627cd47222	L9gViuUDn0olHgQuy_lwOEUrSc-70NNfSMI-IUkkHbU=	2025-11-09 14:36:08.296362	2025-11-02 14:36:08.29664	\N	f
5a992af7-6d0c-49c8-a52c-69139e4bfa7f	7157f3f6-da85-4058-9d38-19133ec93b03	GY48ZB8sxxQUWrDzbDGRLRQenE2zP9USS5EtbLIbHCI=	2025-11-09 14:36:47.417514	2025-11-02 14:36:47.41784	\N	f
26ebdf4e-1812-4d44-a0d4-349d574e8eaa	0197cfc3-7ca2-403b-ae4d-32627cd47222	5mEtDj3MBTx8WTxUuHlQFq5LT-sh4MscbevtopjG0ZQ=	2025-11-09 14:52:04.128534	2025-11-02 14:52:04.128816	\N	f
5431485e-9e43-4da1-86d3-d9338e645f9d	7157f3f6-da85-4058-9d38-19133ec93b03	1akkqwtf4dyh4pX9N2z_GjTLAB0CrqbPUC_OTKq2N2Y=	2025-11-09 14:53:13.203943	2025-11-02 14:53:13.204291	\N	f
\.


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.role_permissions (id, role_id, permission_id, assigned_at, assigned_by) FROM stdin;
2c769a3e-b71d-40e7-8a40-6d8d0a877449	56dc7812-254e-4498-88d9-b03393f3992b	138b5f63-5ed3-4679-8d47-a26d5db26e9f	2025-10-07 17:15:41.774011+00	\N
4bba8e52-721f-4286-b7c4-95c24d295494	56dc7812-254e-4498-88d9-b03393f3992b	5e209000-9ec7-400d-b26c-87cf351e7069	2025-10-07 17:15:41.774011+00	\N
e3c39994-a296-400b-8d87-59ea26cecbe9	56dc7812-254e-4498-88d9-b03393f3992b	d7caf95f-3e98-4458-859a-70246b6df5d6	2025-10-07 17:15:41.774011+00	\N
6f860803-652e-414d-9437-a4a4cb112e73	56dc7812-254e-4498-88d9-b03393f3992b	38bad1db-f9f8-4050-bb28-26859775db73	2025-10-07 17:15:41.774011+00	\N
63b94af3-2bc2-4a9e-b13e-5f4553e7ca85	56dc7812-254e-4498-88d9-b03393f3992b	6d83994d-9c33-49ef-b6d0-54325ef73cc7	2025-10-07 17:15:41.774011+00	\N
7c32f824-2046-478b-99b1-e73374880498	56dc7812-254e-4498-88d9-b03393f3992b	70c71044-4361-409c-868b-5c078a4f1b85	2025-10-07 17:15:41.774011+00	\N
e55734c2-237f-4a9b-84d5-4f649e0cf192	56dc7812-254e-4498-88d9-b03393f3992b	15d60692-9cfb-4bd5-b30d-a2bb0a3462cf	2025-10-07 17:15:41.774011+00	\N
ccf9e338-c353-4d21-9ca6-978f61e9e519	56dc7812-254e-4498-88d9-b03393f3992b	f6dfa1f3-1be8-4f8d-a486-51d5d1c7971d	2025-10-07 17:15:41.774011+00	\N
536c8d52-1881-4def-a00c-a1224f2be35c	56dc7812-254e-4498-88d9-b03393f3992b	ec7368db-48f0-406b-be17-76b050e9ea1f	2025-10-07 17:15:41.774011+00	\N
d978746e-66fc-460e-817a-c3f70c7f87f1	56dc7812-254e-4498-88d9-b03393f3992b	4baab258-54fd-46f3-b7af-6112109d9d1e	2025-10-07 17:15:41.774011+00	\N
563e3d8f-2495-44ca-9650-e8e05f84f9a0	56dc7812-254e-4498-88d9-b03393f3992b	2b110b0f-c8eb-412e-b9d6-26933c373099	2025-10-07 17:15:41.774011+00	\N
74faaea7-ff71-40a5-8937-b26e28be9217	56dc7812-254e-4498-88d9-b03393f3992b	ec205c8b-1e9e-427b-96bc-3af2179ae513	2025-10-07 17:15:41.774011+00	\N
5c1ff963-c507-48e0-aa84-f12fbe408a3b	56dc7812-254e-4498-88d9-b03393f3992b	92b3327c-6a62-47d2-98f1-e0d294d202b6	2025-10-07 17:15:41.774011+00	\N
3f591444-e909-4944-a51b-e212a05c4270	56dc7812-254e-4498-88d9-b03393f3992b	9ab65b7b-3a0b-4a28-80c8-3e4301dfbfa5	2025-10-07 17:15:41.774011+00	\N
fceb451d-e31c-4885-98b2-bb1f54895d17	56dc7812-254e-4498-88d9-b03393f3992b	dd485ead-8246-4f97-a91d-7818efc801e4	2025-10-07 17:15:41.774011+00	\N
bf987ef3-04f8-4398-a558-cbabc9626fc4	56dc7812-254e-4498-88d9-b03393f3992b	ff7ce90d-1ec0-4139-a7cd-36ce8711495f	2025-10-07 17:15:41.774011+00	\N
b47db85d-3961-4b85-aecc-6a0d1f3caab3	56dc7812-254e-4498-88d9-b03393f3992b	3de90dd9-6a0d-4b00-a8f8-c605c0338a62	2025-10-07 17:15:41.774011+00	\N
6009173a-599b-46f0-8eb9-dbd3cd1247b4	56dc7812-254e-4498-88d9-b03393f3992b	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-07 17:15:41.774011+00	\N
b6221b2a-7799-49c8-823f-7f8f1c6d3608	56dc7812-254e-4498-88d9-b03393f3992b	4e71123a-14d5-49c6-92df-c6a76e01ecac	2025-10-07 17:15:41.774011+00	\N
85e3dd28-69e2-430b-b660-2771ef5b3585	56dc7812-254e-4498-88d9-b03393f3992b	7c6478a7-b34a-47cf-ab33-050d7041787b	2025-10-07 17:15:41.774011+00	\N
c332ec84-a7d2-458a-934d-ec1a35341879	56dc7812-254e-4498-88d9-b03393f3992b	31ab0eff-d7ac-45eb-94af-ef97123e3c56	2025-10-07 17:15:41.774011+00	\N
1c7f120f-275d-480b-ae0b-c7acfe2127e2	56dc7812-254e-4498-88d9-b03393f3992b	ee5be0e4-9818-4d0d-aa65-53454c7c08c1	2025-10-07 17:15:41.774011+00	\N
246ab401-1cef-46ae-be21-119fda7dc5e0	56dc7812-254e-4498-88d9-b03393f3992b	96c38e72-1cbc-4761-b45f-c8cae777d26a	2025-10-07 17:15:41.774011+00	\N
894dfcc7-fb1c-460d-999f-187ff32f5906	56dc7812-254e-4498-88d9-b03393f3992b	7771de13-161e-4962-ad7c-1f8c55c0b988	2025-10-07 17:15:41.774011+00	\N
78011cd9-e8c9-48d4-956c-124db64f374c	56dc7812-254e-4498-88d9-b03393f3992b	665c5885-40d6-4480-bbc9-37e22e4315a2	2025-10-07 17:15:41.774011+00	\N
feeda1e3-7386-4179-afcd-2c32b08a90c6	d61c2052-b520-4a31-a699-b839daed32f8	7c6478a7-b34a-47cf-ab33-050d7041787b	2025-10-07 17:15:41.774011+00	\N
4ec0f861-cd5e-4d4d-a575-064f668f2730	d61c2052-b520-4a31-a699-b839daed32f8	31ab0eff-d7ac-45eb-94af-ef97123e3c56	2025-10-07 17:15:41.774011+00	\N
70212fb3-a900-4bec-ad12-bc4f5b27fac8	d61c2052-b520-4a31-a699-b839daed32f8	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-07 17:15:41.774011+00	\N
a92d0597-789d-4d0d-8f6b-212ccf470bdc	d61c2052-b520-4a31-a699-b839daed32f8	4e71123a-14d5-49c6-92df-c6a76e01ecac	2025-10-07 17:15:41.774011+00	\N
9924fda3-150c-419f-a2b5-4a3d1d898343	082cec3d-0e4e-47fb-bb1b-dd2f4c780c6a	31ab0eff-d7ac-45eb-94af-ef97123e3c56	2025-10-07 17:15:41.774011+00	\N
f8d8b809-d0fc-44d5-b6ee-e6fdab90c55f	082cec3d-0e4e-47fb-bb1b-dd2f4c780c6a	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-07 17:15:41.774011+00	\N
033cb018-c6ff-46d7-a16a-cbfabaa1cb99	082cec3d-0e4e-47fb-bb1b-dd2f4c780c6a	4e71123a-14d5-49c6-92df-c6a76e01ecac	2025-10-07 17:15:41.774011+00	\N
cea58c19-1cc1-4c9b-9939-8f6dbdc040ab	e1f16e84-d3a1-44ee-a8b2-d5ce76ea8e8d	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-07 17:15:41.774011+00	\N
df6302be-69a8-4f0d-ad8f-08f1a6a389bb	e1f16e84-d3a1-44ee-a8b2-d5ce76ea8e8d	4e71123a-14d5-49c6-92df-c6a76e01ecac	2025-10-07 17:15:41.774011+00	\N
0c9eb6b1-cac4-4da0-b20b-076965c8f696	6910eb4d-68df-4ef5-8fe0-ab4c8e823b22	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-07 17:15:41.774011+00	\N
7001ea90-6d6f-422e-a406-5645c4bfd633	6910eb4d-68df-4ef5-8fe0-ab4c8e823b22	4e71123a-14d5-49c6-92df-c6a76e01ecac	2025-10-07 17:15:41.774011+00	\N
12a1eb38-3c3b-405f-a13a-89f8bc6008d5	56dc7812-254e-4498-88d9-b03393f3992b	fe009820-67d5-4e7a-a14b-8c89b0e40f42	2025-10-07 18:06:57.049357+00	\N
e10be2fa-997e-42e2-bc11-303008ada919	56dc7812-254e-4498-88d9-b03393f3992b	ffb331b4-8682-41aa-9e24-1ee10de6932f	2025-10-07 19:10:04.648759+00	\N
d04d093c-4df8-4ca4-90b4-b41fef416958	56dc7812-254e-4498-88d9-b03393f3992b	20cae3c5-9649-4533-81f8-0c8a223a1efd	2025-10-07 20:10:40.858169+00	\N
0b12510c-8f6f-403b-8b1b-f4f50c935805	56dc7812-254e-4498-88d9-b03393f3992b	a17c0e6e-132e-442c-bfc1-499119c862de	2025-10-07 20:10:40.858169+00	\N
b01190b0-fd4d-469b-9351-d4f196d7f1af	56dc7812-254e-4498-88d9-b03393f3992b	13a0ce40-f06a-4006-8e47-e848a843a175	2025-10-07 20:10:40.858169+00	\N
606fecf7-035f-47b5-b22a-daca1c71629f	56dc7812-254e-4498-88d9-b03393f3992b	e10a5b01-f824-4236-abc7-bb9ed0d01bd0	2025-10-07 20:10:40.858169+00	\N
ab2c15f4-108b-44a6-b1ee-b29356919848	56dc7812-254e-4498-88d9-b03393f3992b	93e95fb8-91e5-4e5c-96df-8594c93588fd	2025-10-07 20:10:40.858169+00	\N
1b267ed3-a3a7-43b9-8cc1-a2a9755e975d	56dc7812-254e-4498-88d9-b03393f3992b	1cdd5c92-40ac-4f3c-9699-8f37d40c47a6	2025-10-07 20:10:40.858169+00	\N
3f0c7a1f-3400-49af-815d-329a0bb8ee33	56dc7812-254e-4498-88d9-b03393f3992b	d8cf9638-51fe-4639-8c31-ca6b813fdaf2	2025-10-07 20:10:40.858169+00	\N
03cad201-6501-435b-9687-2e1cb43d182d	56dc7812-254e-4498-88d9-b03393f3992b	3e489e46-1b0d-4c3c-a3ff-ee3a035a2222	2025-10-07 20:10:40.858169+00	\N
28fc9ed0-2208-42a3-81c8-4f6324765a48	56dc7812-254e-4498-88d9-b03393f3992b	274746c2-93fe-4c9a-9257-521710a505c5	2025-10-07 20:10:40.858169+00	\N
5cf69205-57b5-40d2-9c4c-b62b43bc7719	56dc7812-254e-4498-88d9-b03393f3992b	f7f09f53-d4f4-4fae-9f49-c0f88797aa80	2025-10-07 20:10:40.858169+00	\N
2bb9c133-a37a-4b36-9254-85547905d29a	56dc7812-254e-4498-88d9-b03393f3992b	32632262-5c8f-426b-9dfd-d1974ade1fb2	2025-10-07 20:10:40.858169+00	\N
1022ce67-f130-4132-844d-b7576b29a435	56dc7812-254e-4498-88d9-b03393f3992b	48353034-3afc-407d-8f33-a87d69923bac	2025-10-07 20:10:40.858169+00	\N
130255a5-c568-4b88-bb30-7028592ea62e	56dc7812-254e-4498-88d9-b03393f3992b	cdca2b10-8719-44f1-a131-6f29a661ffd1	2025-10-07 20:10:40.858169+00	\N
93819f24-2655-456d-81cf-d7ec06c9d760	56dc7812-254e-4498-88d9-b03393f3992b	ed2deabf-a1d0-442d-a8e3-0af563e03d99	2025-10-07 20:10:40.858169+00	\N
f7d4a7a3-8b89-4e69-bd11-10ed3df3b8d0	56dc7812-254e-4498-88d9-b03393f3992b	264dc5ea-f523-42a4-ae30-21c60221d779	2025-10-07 20:10:40.858169+00	\N
e7a43f51-4191-4a7e-9fb5-cb24c1d5d47b	d72e8ca5-eb37-4557-89db-0a7efe3717b9	138b5f63-5ed3-4679-8d47-a26d5db26e9f	2025-10-08 11:02:17.347674+00	\N
dc2fd97b-467c-4ba9-960c-e070b404b3a6	d72e8ca5-eb37-4557-89db-0a7efe3717b9	6d83994d-9c33-49ef-b6d0-54325ef73cc7	2025-10-08 11:02:17.347674+00	\N
e48fb104-97c8-440d-a9ff-00d014bf8a06	d72e8ca5-eb37-4557-89db-0a7efe3717b9	ec7368db-48f0-406b-be17-76b050e9ea1f	2025-10-08 11:02:17.347674+00	\N
97cfbaaa-40ce-43f9-9af3-110055b7d59e	d72e8ca5-eb37-4557-89db-0a7efe3717b9	9ab65b7b-3a0b-4a28-80c8-3e4301dfbfa5	2025-10-08 11:02:17.347674+00	\N
2bbc10af-24b6-4d26-806d-115a7cbb3fb2	d72e8ca5-eb37-4557-89db-0a7efe3717b9	78e5ad68-1209-4662-9abf-6e9206d39d0b	2025-10-08 11:02:17.347674+00	\N
7e6583b4-53c5-42b3-8179-dab1d74af87a	d72e8ca5-eb37-4557-89db-0a7efe3717b9	ee5be0e4-9818-4d0d-aa65-53454c7c08c1	2025-10-08 11:02:17.347674+00	\N
e83af692-4249-4ecb-b6b2-7e8f36a7a546	d72e8ca5-eb37-4557-89db-0a7efe3717b9	fe009820-67d5-4e7a-a14b-8c89b0e40f42	2025-10-08 11:02:17.347674+00	\N
f0d2cb0c-6c6e-4e72-b5af-450224ee4896	d72e8ca5-eb37-4557-89db-0a7efe3717b9	e10a5b01-f824-4236-abc7-bb9ed0d01bd0	2025-10-08 11:02:17.3542+00	\N
d14d3df6-ef18-401b-9895-495062bb1ec0	d72e8ca5-eb37-4557-89db-0a7efe3717b9	d8cf9638-51fe-4639-8c31-ca6b813fdaf2	2025-10-08 11:02:17.3542+00	\N
b5211cc3-d857-4355-b885-49acf2f6759f	d72e8ca5-eb37-4557-89db-0a7efe3717b9	20cae3c5-9649-4533-81f8-0c8a223a1efd	2025-10-08 11:02:17.3542+00	\N
0a877b25-0c73-4ba1-a276-f1311b01f14c	d72e8ca5-eb37-4557-89db-0a7efe3717b9	f7f09f53-d4f4-4fae-9f49-c0f88797aa80	2025-10-08 11:02:17.3542+00	\N
aac5c855-8578-4c58-b32a-981e4dfb22a1	d72e8ca5-eb37-4557-89db-0a7efe3717b9	cdca2b10-8719-44f1-a131-6f29a661ffd1	2025-10-08 11:02:17.3542+00	\N
b3aabc65-cd0f-466c-a0d3-ab42db0d2b92	56dc7812-254e-4498-88d9-b03393f3992b	d8a4f321-e951-4ed5-b095-b96923e4d50d	2025-10-08 12:51:57.512242+00	\N
3440e808-4283-4890-a25a-17aba8a8b435	56dc7812-254e-4498-88d9-b03393f3992b	2a5a5b7c-cd86-4523-b91d-57b665ab2db9	2025-10-08 12:51:57.512242+00	\N
54d58a55-a82c-4e47-a8eb-79872eca8e96	56dc7812-254e-4498-88d9-b03393f3992b	1f46e8d5-73ee-42d3-8691-939cb83f8d79	2025-10-08 12:51:57.512242+00	\N
cf9291f6-7f6e-4c42-aaea-400d2488dc49	56dc7812-254e-4498-88d9-b03393f3992b	d8e96391-c815-484f-be94-74912ae89616	2025-10-08 12:51:57.512242+00	\N
4e71b0a5-0cac-420a-b8a9-1d9ae1faa604	56dc7812-254e-4498-88d9-b03393f3992b	348f8695-5d3c-4c22-be33-41c3a1fb1d58	2025-10-08 12:51:57.512242+00	\N
09f776ac-45e1-416c-811d-7cc527e91a53	56dc7812-254e-4498-88d9-b03393f3992b	716e3672-b648-4ecd-a97b-b6cbbf3bb4e9	2025-10-08 12:51:57.512242+00	\N
e2921659-769c-4d34-9690-37165c0a6a60	56dc7812-254e-4498-88d9-b03393f3992b	d3fb8dea-aa6d-47eb-b9e4-92b9d41f394f	2025-10-08 12:51:57.512242+00	\N
88aaac3c-c338-4e61-8eb0-cf701c7e0b03	56dc7812-254e-4498-88d9-b03393f3992b	4b46d3ca-f207-40ee-986c-644aaca424d4	2025-10-08 13:00:41.554534+00	\N
6f38eed7-c60f-4eda-b5df-1ddaef15bd15	d72e8ca5-eb37-4557-89db-0a7efe3717b9	38bad1db-f9f8-4050-bb28-26859775db73	2025-10-08 13:00:41.554534+00	\N
a0585f36-420a-4626-aad4-b91875916e39	56dc7812-254e-4498-88d9-b03393f3992b	98d16b03-78e5-484c-8f32-f5875e7bdb90	2025-10-08 16:50:48.852878+00	\N
19b33bcc-a78a-460c-a0c2-3ce7a7cb5520	56dc7812-254e-4498-88d9-b03393f3992b	d0161ac5-6ab8-4659-9b29-5d4b99e54a4f	2025-10-11 01:40:42.354802+00	\N
f55aafc6-0a86-4203-a43f-e02766173c9b	56dc7812-254e-4498-88d9-b03393f3992b	00b888b5-23c1-4d0b-b1a6-44f00ddb1fbc	2025-10-11 01:40:42.354802+00	\N
207669f4-41ac-4730-9f62-8d0fdbc4308e	56dc7812-254e-4498-88d9-b03393f3992b	fe3a4df1-571f-4998-836f-ed8a06d917b6	2025-10-11 01:40:42.354802+00	\N
63b52c59-ba2e-47b9-9fa5-11d28e64caae	d72e8ca5-eb37-4557-89db-0a7efe3717b9	d0161ac5-6ab8-4659-9b29-5d4b99e54a4f	2025-10-11 01:40:42.354802+00	\N
601464ff-e765-4864-9e1d-8054c948fd5d	56dc7812-254e-4498-88d9-b03393f3992b	557ef223-08f6-489c-ba18-9b36527f4d03	2025-10-11 01:40:42.360115+00	\N
ba71f0e5-9c83-4ee4-9075-2b2f5d2f355f	56dc7812-254e-4498-88d9-b03393f3992b	8f2bae59-58f7-4d3d-a90d-9532df4e2d19	2025-10-11 01:40:42.360115+00	\N
be42be17-1371-4b2b-b4f6-84665f6fc5a8	56dc7812-254e-4498-88d9-b03393f3992b	167d7532-e67d-42c3-8330-d729fa877c40	2025-10-11 01:40:42.360115+00	\N
69681fe7-aec4-432f-b5b1-980901c99cbc	d72e8ca5-eb37-4557-89db-0a7efe3717b9	557ef223-08f6-489c-ba18-9b36527f4d03	2025-10-11 01:40:42.360115+00	\N
e6e91a5f-2549-46b8-bdc8-afd1036480f3	56dc7812-254e-4498-88d9-b03393f3992b	06ba39c0-1cfd-4e19-87b2-9ec00b1e705e	2025-10-11 01:40:42.445833+00	\N
6dc70a16-dac3-48ba-bb9b-73cda01bc6f3	56dc7812-254e-4498-88d9-b03393f3992b	38ae19a5-834a-464a-82d0-96d6d4c584fe	2025-10-11 01:40:42.445833+00	\N
63ba914d-67ff-4336-8cc0-d90ea867826a	56dc7812-254e-4498-88d9-b03393f3992b	b3d294e1-28bf-4aba-83eb-556ee191b5bb	2025-10-11 01:40:42.445833+00	\N
afe4711f-0367-49ab-93c6-093ab580be08	d72e8ca5-eb37-4557-89db-0a7efe3717b9	06ba39c0-1cfd-4e19-87b2-9ec00b1e705e	2025-10-11 01:40:42.445833+00	\N
6c4328e2-c4c8-4f8b-bf7e-c029546c5a58	56dc7812-254e-4498-88d9-b03393f3992b	4fd44556-4e94-4378-8b36-d3ae23f51b3e	2025-10-11 01:40:42.447813+00	\N
01757b8d-0146-4215-88d3-3cfd1278454b	56dc7812-254e-4498-88d9-b03393f3992b	cb8c4930-6a0b-4479-a521-e4b2d3468955	2025-10-11 01:40:42.447813+00	\N
b44e37f5-6cd7-411c-8b97-9f98a65d85f1	56dc7812-254e-4498-88d9-b03393f3992b	5a4a8eaf-df25-4ee2-a057-9a329cb8c9ce	2025-10-11 01:40:42.447813+00	\N
ac70f6fb-81aa-49e5-9b78-a7002bd8fd20	d72e8ca5-eb37-4557-89db-0a7efe3717b9	4fd44556-4e94-4378-8b36-d3ae23f51b3e	2025-10-11 01:40:42.447813+00	\N
89ef2c8f-a442-4788-8063-85b8e5137bf6	56dc7812-254e-4498-88d9-b03393f3992b	7fff7adf-f84c-44b6-bc6b-7a03249330cf	2025-10-11 01:40:42.449408+00	\N
7a88d929-f3b8-4b69-9035-8557cc6d4c8d	56dc7812-254e-4498-88d9-b03393f3992b	f49a50eb-1b4a-40de-8a86-11fbd44bdb58	2025-10-11 01:40:42.449408+00	\N
153fa850-c3dd-4db2-8bde-c5596fc2a9cd	56dc7812-254e-4498-88d9-b03393f3992b	4f845d71-7e33-4a2a-95a8-89b70fda46f0	2025-10-11 01:40:42.449408+00	\N
f233f0cc-6aaf-4840-8ff2-9984c0231864	d72e8ca5-eb37-4557-89db-0a7efe3717b9	7fff7adf-f84c-44b6-bc6b-7a03249330cf	2025-10-11 01:40:42.449408+00	\N
e500202d-847d-477d-b6d4-0f8847e14066	56dc7812-254e-4498-88d9-b03393f3992b	a24c3aee-10a4-437f-8772-0db10e02dbc0	2025-10-12 18:51:14.154327+00	\N
cfe556eb-b5c2-4029-83ff-145abcd1f2be	56dc7812-254e-4498-88d9-b03393f3992b	2746de28-baac-47b4-a14a-c381983ea512	2025-10-12 18:51:14.154327+00	\N
39c8bd54-4fbf-4ceb-a0ef-6d16b5d125ce	56dc7812-254e-4498-88d9-b03393f3992b	2f06aeab-e4ad-4f44-80ba-4d9643c3a790	2025-10-12 18:51:14.154327+00	\N
faad0f7d-1e5a-41f7-baaa-a6db7c1b806a	d72e8ca5-eb37-4557-89db-0a7efe3717b9	a24c3aee-10a4-437f-8772-0db10e02dbc0	2025-10-12 18:51:14.154327+00	\N
bcb20852-9c70-40dc-918f-1dfd9ca4f983	56dc7812-254e-4498-88d9-b03393f3992b	8cf79304-7204-4df8-9317-5c420d6fb7bc	2025-10-24 19:24:48.264458+00	\N
fa0b8884-8b5c-403a-97e8-d2e38ace8a7b	56dc7812-254e-4498-88d9-b03393f3992b	7181ad0f-19b9-41c0-a1b6-b1c06c223292	2025-10-24 19:24:48.264458+00	\N
2ce58efa-7f5e-4a00-bcdc-6e33e07da19f	d72e8ca5-eb37-4557-89db-0a7efe3717b9	8cf79304-7204-4df8-9317-5c420d6fb7bc	2025-10-24 19:24:48.264458+00	\N
0b38b300-db5b-44a8-80d6-a78e4878d37c	d72e8ca5-eb37-4557-89db-0a7efe3717b9	7181ad0f-19b9-41c0-a1b6-b1c06c223292	2025-10-24 19:24:48.264458+00	\N
3c02fe92-29ba-4a73-9de1-47bf232f3309	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	8cf79304-7204-4df8-9317-5c420d6fb7bc	2025-10-24 19:24:48.264458+00	\N
05aa4b03-82e8-4d1f-b65b-aadb126e744e	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	7181ad0f-19b9-41c0-a1b6-b1c06c223292	2025-10-25 15:20:36.749699+00	\N
d966b015-acf6-49b7-bcad-ab2866ff20bb	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	8cf79304-7204-4df8-9317-5c420d6fb7bc	2025-10-25 15:20:36.749699+00	\N
e80f6ee5-a582-4d6e-a873-cb39d866338e	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	7181ad0f-19b9-41c0-a1b6-b1c06c223292	2025-10-25 15:20:36.749699+00	\N
3b987035-ba8d-4e54-8fb2-d3805248ea52	79e10d6a-4620-4416-bb74-ff9e40639bb1	8cf79304-7204-4df8-9317-5c420d6fb7bc	2025-10-25 15:20:36.749699+00	\N
a9f6af8a-43f4-46a8-bd3d-72c540740057	79e10d6a-4620-4416-bb74-ff9e40639bb1	7181ad0f-19b9-41c0-a1b6-b1c06c223292	2025-10-25 15:20:36.749699+00	\N
1513f78c-33d3-4d4c-9edc-4ce9d2431ed4	d72e8ca5-eb37-4557-89db-0a7efe3717b9	94e18583-8d38-4db6-91e3-95d26d5a4c0d	2025-10-25 15:20:36.749699+00	\N
c3a4719a-45d1-4bfd-b026-2d32c61898ee	d72e8ca5-eb37-4557-89db-0a7efe3717b9	abb03e0b-d835-405c-bdb5-4d5111b0a209	2025-10-25 15:20:36.749699+00	\N
a529ebbf-0760-4cf9-a771-8983765b231a	56dc7812-254e-4498-88d9-b03393f3992b	94e18583-8d38-4db6-91e3-95d26d5a4c0d	2025-10-25 15:20:36.749699+00	\N
0f83e288-9fa7-49d9-a6dc-1baccb49b3ce	56dc7812-254e-4498-88d9-b03393f3992b	abb03e0b-d835-405c-bdb5-4d5111b0a209	2025-10-25 15:20:36.749699+00	\N
c8663a80-683b-46bb-9a02-c18428d6e71e	56dc7812-254e-4498-88d9-b03393f3992b	9d1ea9ec-4042-45fa-9d88-69d2b557535c	2025-10-25 15:20:36.749699+00	\N
7da45279-191e-48e1-a846-26e49fe15ac4	d72e8ca5-eb37-4557-89db-0a7efe3717b9	9d1ea9ec-4042-45fa-9d88-69d2b557535c	2025-10-25 15:20:36.854367+00	\N
9ae1f0b5-b702-4bb3-aeeb-e5c623a589fc	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	2025-10-25 15:26:39.960312+00	\N
fb9e3536-7542-4bc1-a60f-8ef6db9c97b9	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	2025-10-25 15:26:39.960312+00	\N
c37de986-21af-4c29-8106-deddb4c71bfb	79e10d6a-4620-4416-bb74-ff9e40639bb1	556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	2025-10-25 15:26:39.960312+00	\N
be8fc263-79a7-4b00-ba21-8516f7920480	56dc7812-254e-4498-88d9-b03393f3992b	556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	2025-10-25 15:26:39.960312+00	\N
b83273bd-18a5-45d5-91ff-85bf8bba89a8	d72e8ca5-eb37-4557-89db-0a7efe3717b9	556a6ae0-a089-4edc-a5d8-fb62ae8bc3cf	2025-10-25 15:26:40.059367+00	\N
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.roles (id, name, description, is_system_role, created_at, updated_at, created_by) FROM stdin;
56dc7812-254e-4498-88d9-b03393f3992b	admin	Volledige beheerder met toegang tot alle functies	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
d72e8ca5-eb37-4557-89db-0a7efe3717b9	staff	Ondersteunend personeel met beperkte beheerrechten	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
6910eb4d-68df-4ef5-8fe0-ab4c8e823b22	user	Standaard gebruiker	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
d61c2052-b520-4a31-a699-b839daed32f8	owner	Chat kanaal eigenaar	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
082cec3d-0e4e-47fb-bb1b-dd2f4c780c6a	chat_admin	Chat kanaal beheerder	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
e1f16e84-d3a1-44ee-a8b2-d5ce76ea8e8d	member	Chat kanaal lid	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	deelnemer	Evenement deelnemer	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
33dfe3fe-496c-402a-903b-a2cdaba5cc1a	begeleider	Evenement begeleider	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
79e10d6a-4620-4416-bb74-ff9e40639bb1	vrijwilliger	Evenement vrijwilliger	t	2025-10-07 17:15:41.774011+00	2025-10-07 17:15:41.774011+00	\N
c72bda1f-4205-4ffa-ac3e-379de5d857bc	socialmediamanager	Onderhoudt Social Media's.	f	2025-11-01 20:06:01.145206+00	2025-11-01 20:06:01.145206+00	\N
\.


--
-- Data for Name: route_funds; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.route_funds (id, route, amount, created_at, updated_at) FROM stdin;
df1df98b-896e-4247-ba54-76bc93f55650	6 KM	50	2025-10-24 20:27:05.146547+00	2025-10-24 20:27:05.146547+00
bf393c11-14ff-47e4-b007-b789c26f339f	10 KM	75	2025-10-24 20:27:05.146547+00	2025-10-24 20:27:05.146547+00
0965ca21-8207-4b53-9be4-00a9e33d17f2	15 KM	100	2025-10-24 20:27:05.146547+00	2025-10-24 20:27:05.146547+00
0415d2bb-1c51-4057-82a2-bee688670ef1	20 KM	125	2025-10-24 20:27:05.146547+00	2025-10-24 20:27:05.146547+00
\.


--
-- Data for Name: social_embeds; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.social_embeds (id, platform, embed_code, order_number, visible, created_at, updated_at) FROM stdin;
5709a899-ee12-4883-8b57-3a4d0e8543a6	instagram	<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" data-instgrm-version="14" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;">Dit bericht op Instagram bekijken</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center; margin-bottom: 24px;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 224px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 144px;"></div></div></a><p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;"><a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none;" target="_blank">Een bericht gedeeld door Koninklijke Loop (@koninklijkeloop)</a></p></div></blockquote>\r\n<script async src="//www.instagram.com/embed.js"></script>	1	t	2024-12-22 19:30:53.366424+00	2024-12-22 19:30:53.366424+00
ee8d1152-2fd4-464d-82a8-76c02ad56ed9	facebook	<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fpermalink.php%3Fstory_fbid%3Dpfbid02XNU75Y2gMxWhvsVQar7oaM98GvMLLryXQVMTjxnBkEg6e6imJ8ecgoEF9SrTVJDpl%26id%3D61556315443279&show_text=true&width=500" width="500" height="737" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true" allow="autoplay; clipboard-write; encrypted-media; picture-in-picture; web-share"></iframe>	2	t	2024-12-22 19:30:53.366424+00	2024-12-22 19:30:53.366424+00
\.


--
-- Data for Name: social_links; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.social_links (id, platform, url, bg_color_class, icon_color_class, order_number, visible, created_at, updated_at) FROM stdin;
1de3ed0c-9bf3-4924-8d76-72045cc1c0ec	instagram	https://www.instagram.com/koninklijkeloop	\N	\N	2	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
29988672-0d98-4083-9908-18f6bbe34f3f	linkedin	https://www.linkedin.com/company/koninklijkeloop	\N	\N	4	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
dc917a65-6bb4-45cc-a1dc-890b1cdf1f5b	facebook	https://www.facebook.com/koninklijkeloop	\N	\N	1	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
f713d59e-e8a8-40d2-bb11-fc64f89b9ae9	youtube	https://www.youtube.com/@koninklijkeloop	\N	\N	3	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
\.


--
-- Data for Name: sponsors; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.sponsors (id, name, description, logo_url, website_url, order_number, is_active, created_at, updated_at, visible) FROM stdin;
484576a1-2a60-4201-b582-e1f3754ab12a	3x3 Anders	3x3 Anders is een zorgbemiddelingsbureau gespecialiseerd in het matchen van zorgaanbieders met gekwalificeerde zorgprofessionals.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166671/3x3anderslogo_itwm3g.webp	https://3x3anders.nl/	4	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
6408a640-cca1-4aaa-b845-04888f62ccec	Sterk In Vloeren	De website van Sterk In Vloeren biedt een uitgebreid assortiment aan vloeren, waaronder laminaat, PVC-vloeren en tapijt. Ze benadrukken heldere afspraken en hanteren all-in prijzen.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/SterkinVloerenLOGO_zrdofb.webp	https://sterkinvloeren.nl/	1	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
6acd1b1e-8fed-4c8b-89cb-85eee9053536	Beeldpakker	Johan Groot Jebbink, een fotograaf met meer dan tien jaar ervaring, gespecialiseerd in portretfotografie. Actief in Ermelo en internationaal.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166670/BeeldpakkerLogo_wijjmq.webp	https://beeldpakker.nl/	2	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
88cec1c1-3f08-4a6f-9623-c71605fe35b5	Bas Visual Story Telling	BAS Visual Storytelling, heeft een passie voor content en verhalen. Mijn hobby is uitgegroeid tot een eigen onderneming in het vastleggen van verhalen. Bij BAS Visual Storytelling laten we verhalen niet verstoffen op de plank, maar brengen ze tot leven! Waar ik ga of sta, mijn camera's gaan met mij mee, leg de mooiste beelden haarscherp vast en breng jouw verhaal tot leven. Dus vertel eens, 'wat is jouw verhaal?'\r\n\r\n	https://res.cloudinary.com/dgfuv7wif/image/upload/v1746017513/krqjbwwerv9hs6hyrhcy.png	https://basvisualstorytelling.nl/	5	t	2025-04-30 12:51:53.997413+00	2025-05-01 09:30:16.725422+00	t
caa59f1f-65b4-442f-84d2-22cc52212dea	Mojo Dojo	Mojo Dojo is een veelzijdige studio in Rotterdam die diensten aanbiedt voor creatieve producties, waaronder muziekopnames, podcasts en livestreams.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/LogoLayout_1_iphclc.webp	https://mojodojo.studio/	3	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
\.


--
-- Data for Name: title_section_content; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.title_section_content (id, event_title, event_subtitle, image_url, image_alt, detail_1_title, detail_1_description, detail_2_title, detail_2_description, detail_3_title, detail_3_description, created_at, updated_at, participant_count) FROM stdin;
550e8400-e29b-41d4-a716-446655440001	De Koninklijke Loop (DKL) 2026	Op de koninklijke weg in Apeldoorn kunnen mensen met een beperking samen wandelen tijdens dit unieke, rolstoelvriendelijke sponsorloop (DKL), samen met hun verwanten, vrijwilligers of begeleiders.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1760112848/Wij_gaan_17_mei_lopen_voor_hen_3_zllxno_zoqd7z.webp	Promotiebanner De Koninklijke Loop (DKL) 2026: Wij gaan 16 mei lopen voor hen	16 mei 2026	Starttijden variëren per afstand. Zie programma.	Voor iedereen	wandelaars met of zonder beperking (rolstoelvriendelijk).	Lopen voor een goed doel	Steun het goede doel via dit unieke wandelevenement.	2025-04-16 01:31:29.48241+00	2025-10-10 16:21:36.786249+00	75
\.


--
-- Data for Name: under_construction; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.under_construction (id, is_active, title, message, footer_text, logo_url, expected_date, social_links, progress_percentage, contact_email, newsletter_enabled, created_at, updated_at) FROM stdin;
1	f	Website in onderhoud	We stomen ons klaar voor De Koninklijke Loop 2026, op dit moment is de website helaas niet bereikbaar	Bedankt voor uw geduld!	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	2026-01-31 17:00:00+00	[{"url": "https://twitter.com/koninklijkeloop", "platform": "Twitter"}, {"url": "https://instagram.com/koninklijkeloop", "platform": "Instagram"}, {"url": "https://www.youtube.com/@DeKoninklijkeLoop", "platform": "YouTube"}]	85	info@koninklijkeloop.nl	f	2025-09-26 17:37:22.197854+00	2025-10-12 22:02:57.277628+00
7	f	Onder Constructie	Deze website is momenteel onder constructie...	Bedankt voor uw geduld!		\N	[]	0		f	2025-10-12 22:07:19.114274+00	2025-10-12 22:07:39.88183+00
\.


--
-- Data for Name: uploaded_images; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.uploaded_images (id, user_id, public_id, url, secure_url, filename, size, mime_type, width, height, folder, thumbnail_url, deleted_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.user_roles (id, user_id, role_id, assigned_at, assigned_by, expires_at, is_active) FROM stdin;
062b823b-a533-45ed-9df6-031fe0d9313c	7157f3f6-da85-4058-9d38-19133ec93b03	56dc7812-254e-4498-88d9-b03393f3992b	2025-10-07 17:38:28.549894+00	\N	\N	t
505c687d-de01-487d-83aa-37da2206847f	748320dd-5b5e-4434-ad7e-8a405fd6266f	d72e8ca5-eb37-4557-89db-0a7efe3717b9	2025-10-07 20:35:03.355284+00	\N	\N	t
4e892b8d-c63f-433c-bc4c-683877858670	0647c6b8-ec19-4d1f-89de-c95ecd498e9b	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
f723561d-82f7-4638-a296-883e7ae9d5a6	4b6bffd8-3a13-45d3-b284-41bb9c7d101c	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	2025-10-25 11:58:19.59437+00	\N	\N	t
607aaacb-809b-4fca-a3f0-f9b2049487bd	12df00b8-34e9-4768-b010-6a8a4855783d	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
c43268dd-4d9d-4f5b-9d7c-2952752e3dfb	88e61c0b-b432-48c4-84df-2a7796f2f5a2	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
8310ad49-1d0f-4776-8079-eafed62db14a	eadc2084-7bed-4e3f-9c74-22ce7b5e4715	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
aa1c2a97-8b73-4291-8bb6-1b576d7ee56d	e7609af5-f3d1-44c8-b6ff-0c59cbf3f881	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
27498fb7-7ab3-4cd8-aa55-cb5f308f7d00	7fb1db03-7d34-4b42-b046-9a24b5b00068	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
63451f01-aeca-4582-ae1f-5988ec76b7ce	49b72de1-ec0f-4af0-80cd-e6286d320417	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
ce86d748-92b7-427a-8fe7-5876807a7bb9	2aebf5e7-1cfe-4b43-afac-73a27ee49e00	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
2f0c1803-4741-460e-b99d-70381f9dcbe4	a64fb32b-7cdc-4279-89d2-a7e9daeb1f05	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
52352c8c-30b6-4bf4-ab09-1a262dc5d3d6	d12f1566-e096-4b80-b717-ca40036ab4a2	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
e8e0ee7a-8fe0-45d0-8136-9c9a90446a18	cb74dac2-bd74-460f-a83f-04fe02cd9363	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
b0b729a3-0a02-4031-b5ec-965bd3d42f53	2f84da50-3f8d-44ce-a245-e7870cbc779f	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
1171c754-5d02-43a4-ae22-e9af4cec0258	03dff890-db8e-4899-8275-ddbf6421550b	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
2ddace92-f14c-4b48-9181-3be1b74713ea	604ce292-214d-4b3a-86b4-d6cda2e51af0	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
d04e9fd6-2602-4caa-b053-9039779b63be	00a26c26-4932-4091-91ba-4385a251e285	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	2025-10-25 11:58:19.59437+00	\N	\N	t
34eefc0e-e3fd-49df-afb2-72a9eab99e36	5ae279bf-2d9a-43b9-a3c6-98220b3cf92a	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
4386ceb0-6a80-40de-8880-b805a6c8fbc0	e13d7fe0-4a42-4d7d-8a2b-1853f87e005a	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
25b6a1e4-5012-4cf2-8ccf-27f9721cd4c8	554e0410-d07d-4b38-bc7e-2b1a110f3802	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
6daf7ef5-8b4f-4301-b7a6-26bc045cd7e4	8f333073-8d5a-4202-9627-02863000b822	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
7b5df3a1-948a-4fab-98e3-7d045f6fbc7d	564546ae-3227-471c-a743-9713746567f0	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
c664477d-71f5-4466-878f-14c6f9c4d9cb	c9366f76-ec43-44ce-92a0-ea00a4d44b04	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
068ba995-13fa-4393-9523-30d3f98fbc3e	322b1f43-6d6f-433b-a1ea-ee6382a80bde	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
e1eb2f99-03f3-42ae-8d12-ea9b68b8dd37	ba656c7c-a049-46e6-a1db-7259e722fad6	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
7d1a5c83-09af-4d4d-a0d5-6eaffc931b2a	24f08c69-5be4-4d25-a21b-5ddcdf45cfc5	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
47fd9cbf-ad3a-42f2-b49e-34c0fa37617b	01f63c58-2311-41ff-94ba-b727437ee555	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
550721ba-1dc4-4d8a-89a9-299fa0d7f3a3	610f724f-69ed-4457-ab31-1c2572417fc7	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
b9f8e530-5758-4859-bc89-65c3c4f34f05	3f58d1fa-4472-4d62-a31f-6906732b95f0	33dfe3fe-496c-402a-903b-a2cdaba5cc1a	2025-10-25 11:58:19.59437+00	\N	\N	t
ddfce99b-51fb-4335-8815-be5ce6a1f945	ce9f1f44-48af-4d26-987d-cc726cb4dc22	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
ceeee82f-a2b9-4f41-835d-86859c21df81	a76943e8-9673-41c0-8e89-b947294882d7	af3df91f-79ab-4f6a-ac5b-100fe8b0ca91	2025-10-25 11:58:19.59437+00	\N	\N	t
c700abff-77ee-4df3-8193-3d8fd3a01299	b1bec8cb-1709-420f-88cb-fc74e0c6eec2	d72e8ca5-eb37-4557-89db-0a7efe3717b9	2024-12-28 02:43:22+00	\N	\N	t
5dc7bffc-60f1-44b4-b430-a3757e77a8a2	11a3ec93-b159-473b-86d2-d3979b9c9e3a	d72e8ca5-eb37-4557-89db-0a7efe3717b9	2024-12-27 15:40:16.942007+00	\N	\N	t
\.


--
-- Data for Name: verzonden_emails; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.verzonden_emails (id, ontvanger, onderwerp, inhoud, verzonden_op, status, contact_id, aanmelding_id, template_id, created_at, updated_at, fout_bericht) FROM stdin;
\.


--
-- Data for Name: videos; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.videos (id, video_id, url, title, description, thumbnail_url, visible, order_number, created_at, updated_at) FROM stdin;
14ee164b-50e3-4f59-a3e9-3a54312af9cd	q9ngqu	https://streamable.com/e/q9ngqu	De koninklijkeloop!	Preview!!!	\N	t	1	2025-03-28 21:45:53+00	2025-04-21 10:49:03.744+00
18d951d2-f5d1-4b6e-95af-a72ac5ff18ff	x8zj4k	https://streamable.com/e/x8zj4k	Promotie De Koninklijke Loop: Flyers verspreiden	Bekijk hoe vrijwilligers flyers uitdelen om mensen uit te nodigen voor het DKL wandelevenement.	\N	t	4	2025-03-03 20:25:54.251108+00	2025-03-03 20:25:54.251108+00
87502f84-91db-419f-9766-4071ade3e94f	tt6k80	https://streamable.com/e/0o2qf9	Highlights Koninklijke Loop 2024 (Wandelevenement Apeldoorn)	Herbeleef de mooiste momenten en de sfeer van De Koninklijke Loop 2024 in deze highlight video.	\N	t	2	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:42.699+00
99bbe55b-32ef-46ab-bb59-860ce92f1d58	cvfrpi	https://streamable.com/e/cvfrpi	De spannende start van de Koninklijke Loop 2024	Bekijk de start van de deelnemers aan de sponsorloop De Koninklijke Loop 2024.	\N	t	3	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:43.979+00
ac987839-edc1-468f-9080-064f894b3e5d	tt6k80	https://streamable.com/e/tt6k80	Koninklijke Loop 2024 - Hoofdevenement	Een sfeerimpressie van het hoofdevenement van de Koninklijke Loop 2024, met deelnemers, vrijwilligers en muziek.	\N	t	5	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:41.341+00
\.


--
-- Name: migraties_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.migraties_id_seq', 2369, true);


--
-- Name: under_construction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.under_construction_id_seq', 7, true);


--
-- Name: aanmelding_antwoorden aanmelding_antwoorden_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aanmelding_antwoorden
    ADD CONSTRAINT aanmelding_antwoorden_pkey PRIMARY KEY (id);


--
-- Name: aanmeldingen aanmeldingen_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aanmeldingen
    ADD CONSTRAINT aanmeldingen_pkey PRIMARY KEY (id);


--
-- Name: album_photos album_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_pkey PRIMARY KEY (id);


--
-- Name: albums albums_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_pkey PRIMARY KEY (id);


--
-- Name: chat_channel_participants chat_channel_participants_channel_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_channel_id_user_id_key UNIQUE (channel_id, user_id);


--
-- Name: chat_channel_participants chat_channel_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_pkey PRIMARY KEY (id);


--
-- Name: chat_channels chat_channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_channels
    ADD CONSTRAINT chat_channels_pkey PRIMARY KEY (id);


--
-- Name: chat_message_reactions chat_message_reactions_message_id_user_id_emoji_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_message_id_user_id_emoji_key UNIQUE (message_id, user_id, emoji);


--
-- Name: chat_message_reactions chat_message_reactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: chat_user_presence chat_user_presence_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_user_presence
    ADD CONSTRAINT chat_user_presence_pkey PRIMARY KEY (user_id);


--
-- Name: contact_antwoorden contact_antwoorden_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact_antwoorden
    ADD CONSTRAINT contact_antwoorden_pkey PRIMARY KEY (id);


--
-- Name: contact_formulieren contact_formulieren_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact_formulieren
    ADD CONSTRAINT contact_formulieren_pkey PRIMARY KEY (id);


--
-- Name: email_templates email_templates_naam_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_naam_key UNIQUE (naam);


--
-- Name: email_templates email_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_pkey PRIMARY KEY (id);


--
-- Name: gebruikers gebruikers_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_email_key UNIQUE (email);


--
-- Name: gebruikers gebruikers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_pkey PRIMARY KEY (id);


--
-- Name: incoming_emails incoming_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.incoming_emails
    ADD CONSTRAINT incoming_emails_pkey PRIMARY KEY (id);


--
-- Name: incoming_emails incoming_emails_uid_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.incoming_emails
    ADD CONSTRAINT incoming_emails_uid_key UNIQUE (uid);


--
-- Name: migraties migraties_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migraties
    ADD CONSTRAINT migraties_pkey PRIMARY KEY (id);


--
-- Name: newsletters newsletters_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.newsletters
    ADD CONSTRAINT newsletters_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: partners partners_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.partners
    ADD CONSTRAINT partners_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_resource_action_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_resource_action_key UNIQUE (resource, action);


--
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: program_schedule program_schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.program_schedule
    ADD CONSTRAINT program_schedule_pkey PRIMARY KEY (id);


--
-- Name: radio_recordings radio_recordings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.radio_recordings
    ADD CONSTRAINT radio_recordings_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_role_id_permission_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_permission_id_key UNIQUE (role_id, permission_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: route_funds route_funds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.route_funds
    ADD CONSTRAINT route_funds_pkey PRIMARY KEY (id);


--
-- Name: route_funds route_funds_route_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.route_funds
    ADD CONSTRAINT route_funds_route_key UNIQUE (route);


--
-- Name: social_embeds social_embeds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.social_embeds
    ADD CONSTRAINT social_embeds_pkey PRIMARY KEY (id);


--
-- Name: social_links social_links_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.social_links
    ADD CONSTRAINT social_links_pkey PRIMARY KEY (id);


--
-- Name: sponsors sponsors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sponsors
    ADD CONSTRAINT sponsors_pkey PRIMARY KEY (id);


--
-- Name: title_section_content title_section_content_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.title_section_content
    ADD CONSTRAINT title_section_content_pkey PRIMARY KEY (id);


--
-- Name: under_construction under_construction_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.under_construction
    ADD CONSTRAINT under_construction_pkey PRIMARY KEY (id);


--
-- Name: uploaded_images uploaded_images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_pkey PRIMARY KEY (id);


--
-- Name: uploaded_images uploaded_images_public_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_public_id_key UNIQUE (public_id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_user_id_role_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_role_id_key UNIQUE (user_id, role_id);


--
-- Name: verzonden_emails verzonden_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_pkey PRIMARY KEY (id);


--
-- Name: videos videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY (id);


--
-- Name: idx_aanmelding_antwoorden_aanmelding_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmelding_antwoorden_aanmelding_id ON public.aanmelding_antwoorden USING btree (aanmelding_id);


--
-- Name: INDEX idx_aanmelding_antwoorden_aanmelding_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmelding_antwoorden_aanmelding_id IS 'FK index for registration responses';


--
-- Name: idx_aanmelding_antwoorden_aanmelding_verzonden; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmelding_antwoorden_aanmelding_verzonden ON public.aanmelding_antwoorden USING btree (aanmelding_id, verzond_op DESC);


--
-- Name: INDEX idx_aanmelding_antwoorden_aanmelding_verzonden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmelding_antwoorden_aanmelding_verzonden IS 'Chronological responses per registration';


--
-- Name: idx_aanmelding_antwoorden_verzonden_door; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmelding_antwoorden_verzonden_door ON public.aanmelding_antwoorden USING btree (verzonden_door);


--
-- Name: idx_aanmeldingen_afstand; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_afstand ON public.aanmeldingen USING btree (afstand) WHERE (afstand IS NOT NULL);


--
-- Name: INDEX idx_aanmeldingen_afstand; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmeldingen_afstand IS 'Partial index for distance filtering in reports';


--
-- Name: idx_aanmeldingen_antwoorden_count; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_antwoorden_count ON public.aanmeldingen USING btree (antwoorden_count) WHERE (antwoorden_count > 0);


--
-- Name: idx_aanmeldingen_behandeld; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_behandeld ON public.aanmeldingen USING btree (behandeld_op DESC NULLS LAST, status);


--
-- Name: idx_aanmeldingen_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_email ON public.aanmeldingen USING btree (email);


--
-- Name: idx_aanmeldingen_fts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_fts ON public.aanmeldingen USING gin (to_tsvector('dutch'::regconfig, (((((COALESCE(naam, ''::character varying))::text || ' '::text) || (COALESCE(email, ''::character varying))::text) || ' '::text) || COALESCE(bijzonderheden, ''::text))));


--
-- Name: INDEX idx_aanmeldingen_fts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmeldingen_fts IS 'Full-text search on registration details';


--
-- Name: idx_aanmeldingen_gebruiker_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_gebruiker_id ON public.aanmeldingen USING btree (gebruiker_id);


--
-- Name: INDEX idx_aanmeldingen_gebruiker_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmeldingen_gebruiker_id IS 'FK index for user registration lookups';


--
-- Name: idx_aanmeldingen_nieuw; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_nieuw ON public.aanmeldingen USING btree (created_at DESC) WHERE ((status)::text = 'nieuw'::text);


--
-- Name: INDEX idx_aanmeldingen_nieuw; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmeldingen_nieuw IS 'Partial index for new registrations';


--
-- Name: idx_aanmeldingen_rol; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_rol ON public.aanmeldingen USING btree (rol);


--
-- Name: idx_aanmeldingen_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_status ON public.aanmeldingen USING btree (status);


--
-- Name: idx_aanmeldingen_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aanmeldingen_status_created ON public.aanmeldingen USING btree (status, created_at DESC);


--
-- Name: INDEX idx_aanmeldingen_status_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_aanmeldingen_status_created IS 'Compound index for registration dashboard';


--
-- Name: idx_chat_channel_participants_channel_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_channel_participants_channel_id ON public.chat_channel_participants USING btree (channel_id);


--
-- Name: idx_chat_channel_participants_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_channel_participants_user_id ON public.chat_channel_participants USING btree (user_id);


--
-- Name: idx_chat_channels_public; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_channels_public ON public.chat_channels USING btree (name) WHERE ((is_public = true) AND (is_active = true));


--
-- Name: INDEX idx_chat_channels_public; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_channels_public IS 'Public channel discovery';


--
-- Name: idx_chat_channels_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_channels_type ON public.chat_channels USING btree (type);


--
-- Name: INDEX idx_chat_channels_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_channels_type IS 'Channel type filtering';


--
-- Name: idx_chat_message_reactions_message_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_message_reactions_message_id ON public.chat_message_reactions USING btree (message_id);


--
-- Name: idx_chat_messages_channel_id_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_channel_id_created_at ON public.chat_messages USING btree (channel_id, created_at DESC);


--
-- Name: idx_chat_messages_files; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_files ON public.chat_messages USING btree (channel_id, created_at DESC) WHERE (message_type = ANY (ARRAY['image'::text, 'file'::text]));


--
-- Name: INDEX idx_chat_messages_files; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_messages_files IS 'File and image messages';


--
-- Name: idx_chat_messages_fts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_fts ON public.chat_messages USING gin (to_tsvector('dutch'::regconfig, COALESCE(content, ''::text)));


--
-- Name: INDEX idx_chat_messages_fts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_messages_fts IS 'Full-text search on chat message content';


--
-- Name: idx_chat_messages_reply_to; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_reply_to ON public.chat_messages USING btree (reply_to_id, created_at DESC) WHERE (reply_to_id IS NOT NULL);


--
-- Name: INDEX idx_chat_messages_reply_to; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_messages_reply_to IS 'Message reply threads';


--
-- Name: idx_chat_messages_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_user_id ON public.chat_messages USING btree (user_id);


--
-- Name: idx_chat_participants_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_participants_active ON public.chat_channel_participants USING btree (channel_id, user_id) WHERE (is_active = true);


--
-- Name: INDEX idx_chat_participants_active; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_participants_active IS 'Partial index for active channel participants';


--
-- Name: idx_chat_participants_unread; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_participants_unread ON public.chat_channel_participants USING btree (user_id, last_read_at) WHERE (is_active = true);


--
-- Name: INDEX idx_chat_participants_unread; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_participants_unread IS 'Unread message tracking per user';


--
-- Name: idx_chat_user_presence_online; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_user_presence_online ON public.chat_user_presence USING btree (status, last_seen DESC) WHERE (status <> 'offline'::text);


--
-- Name: INDEX idx_chat_user_presence_online; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_chat_user_presence_online IS 'Online and away users';


--
-- Name: idx_contact_antwoorden_contact_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_antwoorden_contact_id ON public.contact_antwoorden USING btree (contact_id);


--
-- Name: INDEX idx_contact_antwoorden_contact_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_contact_antwoorden_contact_id IS 'FK index for contact responses';


--
-- Name: idx_contact_antwoorden_contact_verzonden; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_antwoorden_contact_verzonden ON public.contact_antwoorden USING btree (contact_id, verzond_op DESC);


--
-- Name: INDEX idx_contact_antwoorden_contact_verzonden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_contact_antwoorden_contact_verzonden IS 'Chronological responses per contact';


--
-- Name: idx_contact_antwoorden_verzonden_door; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_antwoorden_verzonden_door ON public.contact_antwoorden USING btree (verzonden_door);


--
-- Name: idx_contact_formulieren_antwoorden_count; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_antwoorden_count ON public.contact_formulieren USING btree (antwoorden_count) WHERE (antwoorden_count > 0);


--
-- Name: idx_contact_formulieren_behandeld; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_behandeld ON public.contact_formulieren USING btree (behandeld_op DESC NULLS LAST, status);


--
-- Name: idx_contact_formulieren_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_email ON public.contact_formulieren USING btree (email);


--
-- Name: idx_contact_formulieren_fts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_fts ON public.contact_formulieren USING gin (to_tsvector('dutch'::regconfig, (((((COALESCE(naam, ''::character varying))::text || ' '::text) || (COALESCE(email, ''::character varying))::text) || ' '::text) || COALESCE(bericht, ''::text))));


--
-- Name: INDEX idx_contact_formulieren_fts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_contact_formulieren_fts IS 'Full-text search on name, email, and message';


--
-- Name: idx_contact_formulieren_nieuw; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_nieuw ON public.contact_formulieren USING btree (created_at DESC) WHERE (((status)::text = 'nieuw'::text) AND (beantwoord = false));


--
-- Name: INDEX idx_contact_formulieren_nieuw; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_contact_formulieren_nieuw IS 'Partial index for new unanswered contact forms';


--
-- Name: idx_contact_formulieren_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_status ON public.contact_formulieren USING btree (status);


--
-- Name: idx_contact_formulieren_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_formulieren_status_created ON public.contact_formulieren USING btree (status, created_at DESC) WHERE (beantwoord = false);


--
-- Name: INDEX idx_contact_formulieren_status_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_contact_formulieren_status_created IS 'Compound index for unanswered contact forms dashboard';


--
-- Name: idx_dashboard_stats_entity_status; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_dashboard_stats_entity_status ON public.dashboard_stats USING btree (entity, status);


--
-- Name: idx_gebruikers_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gebruikers_email ON public.gebruikers USING btree (email);


--
-- Name: idx_gebruikers_is_actief; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gebruikers_is_actief ON public.gebruikers USING btree (is_actief) WHERE (is_actief = true);


--
-- Name: INDEX idx_gebruikers_is_actief; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_gebruikers_is_actief IS 'Partial index for active users only';


--
-- Name: idx_gebruikers_newsletter; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gebruikers_newsletter ON public.gebruikers USING btree (email) WHERE ((newsletter_subscribed = true) AND (is_actief = true));


--
-- Name: INDEX idx_gebruikers_newsletter; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_gebruikers_newsletter IS 'Partial index for active newsletter subscribers';


--
-- Name: idx_gebruikers_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gebruikers_role_id ON public.gebruikers USING btree (role_id);


--
-- Name: INDEX idx_gebruikers_role_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_gebruikers_role_id IS 'FK index for RBAC role lookups';


--
-- Name: idx_incoming_emails_account_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incoming_emails_account_type ON public.incoming_emails USING btree (account_type);


--
-- Name: idx_incoming_emails_from; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incoming_emails_from ON public.incoming_emails USING btree ("from");


--
-- Name: INDEX idx_incoming_emails_from; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_incoming_emails_from IS 'Sender lookup for email filtering';


--
-- Name: idx_incoming_emails_is_processed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incoming_emails_is_processed ON public.incoming_emails USING btree (is_processed);


--
-- Name: idx_incoming_emails_message_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incoming_emails_message_id ON public.incoming_emails USING btree (message_id);


--
-- Name: idx_incoming_emails_processing; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_incoming_emails_processing ON public.incoming_emails USING btree (is_processed, received_at DESC) WHERE (is_processed = false);


--
-- Name: INDEX idx_incoming_emails_processing; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_incoming_emails_processing IS 'Partial index for unprocessed email queue';


--
-- Name: idx_migraties_versie; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_migraties_versie ON public.migraties USING btree (versie);


--
-- Name: idx_newsletters_sent_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_newsletters_sent_at ON public.newsletters USING btree (sent_at);


--
-- Name: idx_newsletters_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_newsletters_status ON public.newsletters USING btree (sent_at DESC);


--
-- Name: INDEX idx_newsletters_status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_newsletters_status IS 'Draft newsletters first, then sent in reverse chronological order';


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at);


--
-- Name: idx_notifications_priority; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_priority ON public.notifications USING btree (priority);


--
-- Name: idx_notifications_sent; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_sent ON public.notifications USING btree (sent);


--
-- Name: idx_notifications_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);


--
-- Name: idx_permissions_resource_action; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_resource_action ON public.permissions USING btree (resource, action);


--
-- Name: idx_refresh_tokens_cleanup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_cleanup ON public.refresh_tokens USING btree (expires_at) WHERE (is_revoked = false);


--
-- Name: INDEX idx_refresh_tokens_cleanup; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_refresh_tokens_cleanup IS 'Expired token cleanup for scheduled jobs';


--
-- Name: idx_refresh_tokens_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--
-- Name: idx_refresh_tokens_is_revoked; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_is_revoked ON public.refresh_tokens USING btree (is_revoked);


--
-- Name: idx_refresh_tokens_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_role_permissions_permission_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_permission_id ON public.role_permissions USING btree (permission_id);


--
-- Name: idx_role_permissions_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_role_id ON public.role_permissions USING btree (role_id);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_route_funds_route; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_route_funds_route ON public.route_funds USING btree (route);


--
-- Name: idx_uploaded_images_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_active ON public.uploaded_images USING btree (user_id, created_at DESC) WHERE (deleted_at IS NULL);


--
-- Name: INDEX idx_uploaded_images_active; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_uploaded_images_active IS 'Active (non-deleted) images per user';


--
-- Name: idx_uploaded_images_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_created_at ON public.uploaded_images USING btree (created_at DESC);


--
-- Name: idx_uploaded_images_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_deleted_at ON public.uploaded_images USING btree (deleted_at);


--
-- Name: idx_uploaded_images_folder; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_folder ON public.uploaded_images USING btree (folder);


--
-- Name: idx_uploaded_images_public_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_public_id ON public.uploaded_images USING btree (public_id);


--
-- Name: idx_uploaded_images_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_uploaded_images_user_id ON public.uploaded_images USING btree (user_id);


--
-- Name: idx_user_roles_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_active ON public.user_roles USING btree (is_active) WHERE (is_active = true);


--
-- Name: idx_user_roles_active_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_active_lookup ON public.user_roles USING btree (user_id, role_id) WHERE (is_active = true);


--
-- Name: INDEX idx_user_roles_active_lookup; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_user_roles_active_lookup IS 'Active user role assignments';


--
-- Name: idx_user_roles_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_role_id ON public.user_roles USING btree (role_id);


--
-- Name: idx_user_roles_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_user_id ON public.user_roles USING btree (user_id);


--
-- Name: idx_verzonden_emails_aanmelding_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_aanmelding_id ON public.verzonden_emails USING btree (aanmelding_id);


--
-- Name: INDEX idx_verzonden_emails_aanmelding_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_aanmelding_id IS 'FK index for registration email tracking';


--
-- Name: idx_verzonden_emails_contact_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_contact_id ON public.verzonden_emails USING btree (contact_id);


--
-- Name: INDEX idx_verzonden_emails_contact_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_contact_id IS 'FK index for contact form email tracking';


--
-- Name: idx_verzonden_emails_errors; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_errors ON public.verzonden_emails USING btree (verzonden_op DESC) WHERE ((status)::text = 'failed'::text);


--
-- Name: INDEX idx_verzonden_emails_errors; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_errors IS 'Partial index for failed email tracking';


--
-- Name: idx_verzonden_emails_ontvanger; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_ontvanger ON public.verzonden_emails USING btree (ontvanger);


--
-- Name: idx_verzonden_emails_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_status ON public.verzonden_emails USING btree (status);


--
-- Name: INDEX idx_verzonden_emails_status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_status IS 'Status filtering for error tracking';


--
-- Name: idx_verzonden_emails_status_tijd; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_status_tijd ON public.verzonden_emails USING btree (status, verzonden_op DESC);


--
-- Name: INDEX idx_verzonden_emails_status_tijd; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_status_tijd IS 'Compound index for email status tracking over time';


--
-- Name: idx_verzonden_emails_template_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_template_id ON public.verzonden_emails USING btree (template_id);


--
-- Name: INDEX idx_verzonden_emails_template_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_template_id IS 'FK index for template usage tracking';


--
-- Name: idx_verzonden_emails_verzonden_op; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_verzonden_emails_verzonden_op ON public.verzonden_emails USING btree (verzonden_op DESC);


--
-- Name: INDEX idx_verzonden_emails_verzonden_op; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_verzonden_emails_verzonden_op IS 'Chronological sorting for email history';


--
-- Name: aanmelding_antwoorden trigger_aanmelding_antwoorden_count; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_aanmelding_antwoorden_count AFTER INSERT OR DELETE ON public.aanmelding_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_aanmelding_antwoorden_count();


--
-- Name: aanmelding_antwoorden trigger_aanmelding_antwoorden_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_aanmelding_antwoorden_updated_at BEFORE UPDATE ON public.aanmelding_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: aanmeldingen trigger_aanmeldingen_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_aanmeldingen_updated_at BEFORE UPDATE ON public.aanmeldingen FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: albums trigger_albums_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_albums_updated_at BEFORE UPDATE ON public.albums FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_channels trigger_chat_channels_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_chat_channels_updated_at BEFORE UPDATE ON public.chat_channels FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_messages trigger_chat_messages_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_chat_messages_updated_at BEFORE UPDATE ON public.chat_messages FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_user_presence trigger_chat_user_presence_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_chat_user_presence_updated_at BEFORE UPDATE ON public.chat_user_presence FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: contact_antwoorden trigger_contact_antwoorden_count; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_contact_antwoorden_count AFTER INSERT OR DELETE ON public.contact_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_contact_antwoorden_count();


--
-- Name: contact_antwoorden trigger_contact_antwoorden_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_contact_antwoorden_updated_at BEFORE UPDATE ON public.contact_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: contact_formulieren trigger_contact_formulieren_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_contact_formulieren_updated_at BEFORE UPDATE ON public.contact_formulieren FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: email_templates trigger_email_templates_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_email_templates_updated_at BEFORE UPDATE ON public.email_templates FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: gebruikers trigger_gebruikers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_gebruikers_updated_at BEFORE UPDATE ON public.gebruikers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: incoming_emails trigger_incoming_emails_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_incoming_emails_updated_at BEFORE UPDATE ON public.incoming_emails FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: newsletters trigger_newsletters_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_newsletters_updated_at BEFORE UPDATE ON public.newsletters FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: photos trigger_photos_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_photos_updated_at BEFORE UPDATE ON public.photos FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: route_funds trigger_route_funds_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_route_funds_updated_at BEFORE UPDATE ON public.route_funds FOR EACH ROW EXECUTE FUNCTION public.update_route_funds_updated_at();


--
-- Name: sponsors trigger_sponsors_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_sponsors_updated_at BEFORE UPDATE ON public.sponsors FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: uploaded_images trigger_uploaded_images_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_uploaded_images_updated_at BEFORE UPDATE ON public.uploaded_images FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: verzonden_emails trigger_verzonden_emails_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_verzonden_emails_updated_at BEFORE UPDATE ON public.verzonden_emails FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: videos trigger_videos_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_videos_updated_at BEFORE UPDATE ON public.videos FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: aanmelding_antwoorden aanmelding_antwoorden_aanmelding_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aanmelding_antwoorden
    ADD CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey FOREIGN KEY (aanmelding_id) REFERENCES public.aanmeldingen(id) ON DELETE CASCADE;


--
-- Name: CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey ON aanmelding_antwoorden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey ON public.aanmelding_antwoorden IS 'FK to aanmeldingen';


--
-- Name: aanmeldingen aanmeldingen_gebruiker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aanmeldingen
    ADD CONSTRAINT aanmeldingen_gebruiker_id_fkey FOREIGN KEY (gebruiker_id) REFERENCES public.gebruikers(id);


--
-- Name: album_photos album_photos_album_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_album_id_fkey FOREIGN KEY (album_id) REFERENCES public.albums(id) ON DELETE CASCADE;


--
-- Name: album_photos album_photos_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: albums albums_cover_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_cover_photo_id_fkey FOREIGN KEY (cover_photo_id) REFERENCES public.photos(id);


--
-- Name: chat_channel_participants chat_channel_participants_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.chat_channels(id) ON DELETE CASCADE;


--
-- Name: chat_message_reactions chat_message_reactions_message_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_message_id_fkey FOREIGN KEY (message_id) REFERENCES public.chat_messages(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.chat_channels(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_reply_to_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_reply_to_id_fkey FOREIGN KEY (reply_to_id) REFERENCES public.chat_messages(id) ON DELETE SET NULL;


--
-- Name: contact_antwoorden contact_antwoorden_contact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact_antwoorden
    ADD CONSTRAINT contact_antwoorden_contact_id_fkey FOREIGN KEY (contact_id) REFERENCES public.contact_formulieren(id) ON DELETE CASCADE;


--
-- Name: CONSTRAINT contact_antwoorden_contact_id_fkey ON contact_antwoorden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON CONSTRAINT contact_antwoorden_contact_id_fkey ON public.contact_antwoorden IS 'FK to contact_formulieren';


--
-- Name: email_templates email_templates_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.gebruikers(id);


--
-- Name: gebruikers gebruikers_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.gebruikers(id);


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: roles roles_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.gebruikers(id);


--
-- Name: uploaded_images uploaded_images_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.gebruikers(id);


--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: verzonden_emails verzonden_emails_aanmelding_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_aanmelding_id_fkey FOREIGN KEY (aanmelding_id) REFERENCES public.aanmeldingen(id);


--
-- Name: verzonden_emails verzonden_emails_contact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_contact_id_fkey FOREIGN KEY (contact_id) REFERENCES public.contact_formulieren(id);


--
-- Name: verzonden_emails verzonden_emails_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_template_id_fkey FOREIGN KEY (template_id) REFERENCES public.email_templates(id);


--
-- Name: dashboard_stats; Type: MATERIALIZED VIEW DATA; Schema: public; Owner: -
--

REFRESH MATERIALIZED VIEW public.dashboard_stats;


--
-- PostgreSQL database dump complete
--

\unrestrict 0jXUPxC9NSPHoCpxmC2c42IQDmqL4d8huoBscL62YaMpTxcWQsnDsF0Jrw5Clge

