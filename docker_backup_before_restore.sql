--
-- PostgreSQL database dump
--

\restrict ynkt5EaPkd5cffaI1hn4JWbqquN1EXCdrwsFHR8nB3jEQXvwxjdADKAZaZd6Fbx

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

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
-- Name: refresh_dashboard_stats(); Type: FUNCTION; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE FUNCTION public.refresh_dashboard_stats() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY dashboard_stats;
END;
$$;


ALTER FUNCTION public.refresh_dashboard_stats() OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: FUNCTION refresh_dashboard_stats(); Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON FUNCTION public.refresh_dashboard_stats() IS 'Refresh dashboard statistics view (concurrent safe)';


--
-- Name: update_aanmelding_antwoorden_count(); Type: FUNCTION; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER FUNCTION public.update_aanmelding_antwoorden_count() OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: update_contact_antwoorden_count(); Type: FUNCTION; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER FUNCTION public.update_contact_antwoorden_count() OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: update_route_funds_updated_at(); Type: FUNCTION; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE FUNCTION public.update_route_funds_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_route_funds_updated_at() OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: FUNCTION update_updated_at_column(); Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON FUNCTION public.update_updated_at_column() IS 'Generic trigger function to update updated_at timestamp';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: aanmelding_antwoorden; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.aanmelding_antwoorden OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE aanmelding_antwoorden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.aanmelding_antwoorden IS 'Antwoorden op aanmeldingen';


--
-- Name: aanmeldingen; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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
    gebruiker_id uuid,
    steps integer DEFAULT 0,
    antwoorden_count integer DEFAULT 0,
    CONSTRAINT aanmeldingen_email_check CHECK (((email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT aanmeldingen_naam_not_empty CHECK ((length(TRIM(BOTH FROM naam)) > 0)),
    CONSTRAINT aanmeldingen_steps_check CHECK ((steps >= 0))
);
ALTER TABLE ONLY public.aanmeldingen ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE ONLY public.aanmeldingen ALTER COLUMN status SET STATISTICS 500;


ALTER TABLE public.aanmeldingen OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE aanmeldingen; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.aanmeldingen IS 'Aanmeldingen voor De Koninklijke Loop';


--
-- Name: COLUMN aanmeldingen.rol; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.aanmeldingen.rol IS 'Rol van de deelnemer (deelnemer, vrijwilliger, sponsor)';


--
-- Name: COLUMN aanmeldingen.afstand; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.aanmeldingen.afstand IS 'Gekozen afstand voor hardlopers';


--
-- Name: COLUMN aanmeldingen.test_mode; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.aanmeldingen.test_mode IS 'Geeft aan of dit een testaanmelding is (geen echte email verzenden)';


--
-- Name: COLUMN aanmeldingen.gebruiker_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.aanmeldingen.gebruiker_id IS 'Link naar gebruikersaccount voor authenticatie en step tracking';


--
-- Name: COLUMN aanmeldingen.antwoorden_count; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.aanmeldingen.antwoorden_count IS 'Cached count of responses - updated via trigger';


--
-- Name: album_photos; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.album_photos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    album_id uuid NOT NULL,
    photo_id uuid NOT NULL,
    order_number integer,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.album_photos OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: albums; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.albums OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: chat_channel_participants; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.chat_channel_participants OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: chat_channels; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.chat_channels OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: chat_message_reactions; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.chat_message_reactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message_id uuid,
    user_id uuid,
    emoji text NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.chat_message_reactions OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.chat_messages OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE chat_messages; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.chat_messages IS 'High-traffic table - monitor size and consider partitioning';


--
-- Name: chat_user_presence; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.chat_user_presence (
    user_id uuid NOT NULL,
    status text DEFAULT 'offline'::text,
    last_seen timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT chat_user_presence_status_check CHECK ((status = ANY (ARRAY['online'::text, 'away'::text, 'busy'::text, 'offline'::text])))
);


ALTER TABLE public.chat_user_presence OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: contact_antwoorden; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.contact_antwoorden OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE contact_antwoorden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.contact_antwoorden IS 'Antwoorden op contactformulieren';


--
-- Name: contact_formulieren; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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
    CONSTRAINT contact_formulieren_naam_not_empty CHECK ((length(TRIM(BOTH FROM naam)) > 0))
);
ALTER TABLE ONLY public.contact_formulieren ALTER COLUMN email SET STATISTICS 1000;
ALTER TABLE ONLY public.contact_formulieren ALTER COLUMN status SET STATISTICS 500;


ALTER TABLE public.contact_formulieren OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE contact_formulieren; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.contact_formulieren IS 'Contactformulieren van de website';


--
-- Name: COLUMN contact_formulieren.email_verzonden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.contact_formulieren.email_verzonden IS 'Geeft aan of er een email is verzonden naar de afzender';


--
-- Name: COLUMN contact_formulieren.test_mode; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.contact_formulieren.test_mode IS 'Geeft aan of dit een testbericht is (geen echte email verzenden)';


--
-- Name: COLUMN contact_formulieren.antwoorden_count; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.contact_formulieren.antwoorden_count IS 'Cached count of responses - updated via trigger';


--
-- Name: verzonden_emails; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.verzonden_emails OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE verzonden_emails; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.verzonden_emails IS 'High-traffic table - monitor size and consider partitioning';


--
-- Name: dashboard_stats; Type: MATERIALIZED VIEW; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.dashboard_stats OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: MATERIALIZED VIEW dashboard_stats; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON MATERIALIZED VIEW public.dashboard_stats IS 'Cached dashboard statistics - refresh hourly or on demand';


--
-- Name: email_templates; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.email_templates OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: gebruikers; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.gebruikers OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: incoming_emails; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.incoming_emails OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE incoming_emails; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.incoming_emails IS 'Monitor for unprocessed email buildup';


--
-- Name: migraties; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.migraties (
    id bigint NOT NULL,
    versie text NOT NULL,
    naam text NOT NULL,
    toegepast timestamp with time zone
);


ALTER TABLE public.migraties OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: migraties_id_seq; Type: SEQUENCE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE SEQUENCE public.migraties_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.migraties_id_seq OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: migraties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER SEQUENCE public.migraties_id_seq OWNED BY public.migraties.id;


--
-- Name: newsletters; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.newsletters OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: notifications; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.notifications OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE notifications; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.notifications IS 'Stores notifications to be sent via Telegram';


--
-- Name: partners; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.partners OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: permissions; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.permissions OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: photos; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.photos OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: program_schedule; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.program_schedule OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: radio_recordings; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.radio_recordings OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.refresh_tokens OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: TABLE refresh_tokens; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON TABLE public.refresh_tokens IS 'Stores refresh tokens for JWT authentication with 7-day expiry';


--
-- Name: COLUMN refresh_tokens.token; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.refresh_tokens.token IS 'Base64 encoded random token (32 bytes)';


--
-- Name: COLUMN refresh_tokens.expires_at; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.refresh_tokens.expires_at IS 'Token expiration timestamp (7 days from creation)';


--
-- Name: COLUMN refresh_tokens.is_revoked; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON COLUMN public.refresh_tokens.is_revoked IS 'Whether the token has been revoked (for token rotation)';


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.role_permissions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    assigned_at timestamp with time zone DEFAULT now(),
    assigned_by uuid
);


ALTER TABLE public.role_permissions OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.roles OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: route_funds; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TABLE public.route_funds (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    route character varying(50) NOT NULL,
    amount integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT route_funds_amount_check CHECK ((amount >= 0))
);


ALTER TABLE public.route_funds OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: social_embeds; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.social_embeds OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: social_links; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.social_links OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: sponsors; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.sponsors OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: title_section_content; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.title_section_content OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: under_construction; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.under_construction OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: under_construction_id_seq; Type: SEQUENCE; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE SEQUENCE public.under_construction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.under_construction_id_seq OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: under_construction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER SEQUENCE public.under_construction_id_seq OWNED BY public.under_construction.id;


--
-- Name: uploaded_images; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.uploaded_images OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.user_roles OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: user_permissions; Type: VIEW; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.user_permissions OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: v_user_role_migration_status; Type: VIEW; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.v_user_role_migration_status OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: videos; Type: TABLE; Schema: public; Owner: dekoninklijkeloopdatabase_user
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


ALTER TABLE public.videos OWNER TO dekoninklijkeloopdatabase_user;

--
-- Name: migraties id; Type: DEFAULT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.migraties ALTER COLUMN id SET DEFAULT nextval('public.migraties_id_seq'::regclass);


--
-- Name: under_construction id; Type: DEFAULT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.under_construction ALTER COLUMN id SET DEFAULT nextval('public.under_construction_id_seq'::regclass);


--
-- Data for Name: aanmelding_antwoorden; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.aanmelding_antwoorden (id, aanmelding_id, verzonden_op, created_at, updated_at, tekst, verzond_op, verzond_door, email_verzonden, verzonden_door) FROM stdin;
\.


--
-- Data for Name: aanmeldingen; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.aanmeldingen (id, naam, email, telefoon, status, created_at, updated_at, rol, afstand, ondersteuning, bijzonderheden, terms, email_verzonden, email_verzonden_op, behandeld_door, behandeld_op, notities, test_mode, gebruiker_id, steps, antwoorden_count) FROM stdin;
3e62d5d3-070d-47b1-a1ef-30665f982789	TGTest	laventejeffrey@gmail.com	06123456789	nieuw	2025-03-23 17:06:26.132297	2025-11-02 15:04:11.346757	Begeleider	15 KM	Anders	Telegram Test bericht - officiele weg	t	f	\N	\N	\N	\N	f	\N	0	0
51205069-a20a-4231-bcb9-cb2d6fd042c8	Manuela van Zwam	rik.van-harxen@sheerenloo.nl	\N	nieuw	2025-03-23 14:08:53.328628	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Ja	Vaste begeleider die meeloopt	t	f	\N	\N	\N	\N	f	\N	0	0
275490c0-1021-4bf4-9005-7df9884b0fe6	Bas heijenk 	basheijenk96@gmail.com	\N	nieuw	2025-03-22 16:43:03.19496	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
b5e67c64-6bfa-46d4-ae40-98b093c8b720	Salih	topraks@gmail.com	\N	nieuw	2025-03-17 19:55:07.379647	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Nee		t	t	2025-03-21 06:47:38.608	\N	\N	\N	f	\N	0	0
b2fd3412-8368-409f-8029-b2cdd581ade1	Manuela van zwam	benjaminlaan.64a@sheerenloo.nl	\N	nieuw	2025-03-11 11:28:05.483427	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Nee		t	t	2025-03-21 06:47:39.928	\N	\N	\N	f	\N	0	0
391f63c5-f034-466e-8a1f-ba9d06ed1192	Joyce Thielen	Joyce.thielen@sheerenloo.nl		nieuw	2025-03-09 16:52:09.564437	2025-11-02 15:04:11.346757	Begeleider	6 KM	Nee		t	t	2025-03-21 06:47:41.071	\N	\N	\N	f	\N	0	0
391f2579-d7cb-4ef3-afbe-14dc4115c519	Dick van Norden	Enckerkamp.27@sheerenloo.nl	\N	nieuw	2025-03-08 10:28:37.053379	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-03-08 10:28:42.2	\N	\N	\N	f	\N	0	0
90e477cc-89d1-4524-8edc-b697be8c504d	Angelo van Ingen	Enckerkamp.27@sheerenloo.nl	\N	nieuw	2025-03-08 10:27:31.078013	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-03-08 10:27:36.787	\N	\N	\N	f	\N	0	0
26ea058b-2608-49d0-862a-611e98d7dc61	Janny van de Wall	mjvdwal@hotmail.com	\N	nieuw	2025-02-20 11:55:30.818333	2025-11-02 15:04:11.346757	Deelnemer	10 KM	Nee		t	t	2025-02-20 11:55:33.048	\N	\N	\N	f	\N	0	0
f4fc2312-ec8a-4dfc-90b5-a8da317618e6	Martin van der Wal	mjvdwal@hotmail.com		nieuw	2025-01-28 21:58:37.55756	2025-11-02 15:04:11.346757	Deelnemer	10 KM	Ja	loopt samen met Dirk-Jan mee als vrijwilliger	t	t	2025-01-28 21:58:39.243	\N	\N	\N	f	\N	0	0
4bfe814b-e0b8-4e60-9f46-fe38852d9ecb	Dirk-Jan Hempe	mjvdwal@hotmail.com		nieuw	2025-01-28 21:56:10.099964	2025-11-02 15:04:11.346757	Deelnemer	10 KM	Nee		t	t	2025-01-28 21:56:12.004	\N	\N	\N	f	\N	0	0
9f464844-8c93-4190-90f0-e74765c7f09a	Karin de Jong	karin.de.jong82@outlook.com	\N	nieuw	2025-03-24 18:47:05.053651	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
51855fec-eab9-494a-9321-c40d22da4ffc	Mirjam Kerkvliet	mirjam.kerkvliet@gmail.com	\N	nieuw	2025-03-24 17:27:25.191642	2025-11-02 15:04:11.346757	Deelnemer	15 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
47775742-8950-4b94-9dd1-571ff4902688	Arno Kerkvliet	arno.kerkvliet@gmail.com	\N	nieuw	2025-03-24 17:26:12.1724	2025-11-02 15:04:11.346757	Deelnemer	15 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
d92ed75c-c275-47a4-88a9-ff7a4106f8ee	Jean-paul Hup	molenkamp.19@sheerenloo.nl	\N	nieuw	2025-03-24 09:17:59.726501	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
ecb8332b-ea39-4611-9f58-64921226f2a6	Annerieke Mandemaker-Timmer	annerieketimmer@hotmail.com	06 17 37 28 40 	nieuw	2025-03-24 09:16:42.111808	2025-11-02 15:04:11.346757	Begeleider	6 KM	Nee		t	f	\N	\N	\N	\N	f	\N	0	0
1ca80f61-f5c1-431f-b224-e6557150b65b	Han van Doornik	LaanvanGS.26@sheerenloo.nl	\N	nieuw	2025-03-30 08:07:45.334762	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Ja	Ik wil wel graag begeleiding 	t	f	\N	\N	\N	\N	f	\N	0	0
3c420d85-5f78-4645-88dd-c9e90eb8b6f7	Henk Rekers 	h.rekers1959@kpnmail.nl	\N	nieuw	2025-04-12 10:21:37.470747	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee	email mislukt	t	f	2025-04-12 16:52:21.796	\N	\N	\N	f	\N	0	0
721d297d-d670-46b1-a581-bcc095565bbd	Hilde Rekers 	h.rekers1959@kpnmail.nl	\N	nieuw	2025-04-12 10:20:25.399025	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee	email mislukt	t	f	2025-04-12 16:52:20.444	\N	\N	\N	f	\N	0	0
2499686e-2e62-4827-9079-78b468cb26c9	Bertram tijsma	Klaskehiddes@gmail.com	\N	nieuw	2025-03-29 10:07:12.393466	2025-11-02 15:04:11.346757	Deelnemer	15 KM	Nee		t	t	2025-03-29 10:07:57.838	\N	\N	\N	f	\N	0	0
9f75b1df-4c72-4e36-9901-0f74cc26574f	Klaske van de glind	Klaskehiddes@gmail.com	\N	nieuw	2025-03-29 10:05:04.276906	2025-11-02 15:04:11.346757	Deelnemer	15 KM	Nee		t	t	2025-03-29 10:07:57.215	\N	\N	\N	f	\N	0	0
db3ec762-dd54-4ba7-98ab-981235cc316a	Mila Veenendaal	gaminggirlayla@gmail.com	\N	nieuw	2025-03-26 16:26:53.777384	2025-11-02 15:04:11.346757	Deelnemer	10 KM	Nee		t	t	2025-03-26 16:30:57.269	\N	\N	\N	f	\N	0	0
d17c16c6-c423-43de-a876-d40326b62d9e	Ayla Toprak	gamergirlayla@gmail.com	\N	nieuw	2025-03-26 16:25:43.756211	2025-11-02 15:04:11.346757	Deelnemer	10 KM	Nee		t	t	2025-03-26 16:30:56.685	\N	\N	\N	f	\N	0	0
917206a7-a28d-4bad-8b41-ed127eab743a	A. Bistolfi	nedarg@icloud.com	\N	nieuw	2025-03-26 12:27:20.236848	2025-11-02 15:04:11.346757	Deelnemer	15 KM	Nee		t	t	2025-03-26 16:30:56.062	\N	\N	\N	f	\N	0	0
61a3f823-82f6-4107-a9a9-b6a18e6b12c3	Henk Rekers 	h.rekers59@kpnmail.nl	\N	nieuw	2025-04-14 20:03:47.361483	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-04-15 13:23:40.051	\N	\N	\N	f	\N	0	0
8a5c2c59-ebca-411e-9471-a38be47e1192	Hilde Rekers 	h.rekers59@kpnmail.nl	\N	nieuw	2025-04-14 20:00:47.479484	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-04-15 13:23:39.203	\N	\N	\N	f	\N	0	0
02f9605b-8b26-4461-9ca7-ed7ebbd12311	Theun 	diesbosje@hotmail.com	\N	nieuw	2025-04-12 07:09:19.884142	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-04-12 16:52:19.056	\N	\N	\N	f	\N	0	0
282e6ec6-2c97-4b46-a992-273c826c1f91	Albert 	diesbosje@hotmail.com	\N	nieuw	2025-04-12 07:08:26.565956	2025-11-02 15:04:11.346757	Deelnemer	6 KM	Nee		t	t	2025-04-12 16:52:18.027	\N	\N	\N	f	\N	0	0
d10065f9-229d-4a1d-9178-324bdffa57c0	Diesmer 	diesbosje@hotmail.com	0613429612	nieuw	2025-04-12 07:06:59.872906	2025-11-02 15:04:11.346757	Begeleider	6 KM	Nee		t	t	2025-04-12 16:52:17.255	\N	\N	\N	f	\N	0	0
613770dd-4733-4b54-964d-c3753cfebfd7	Sylvia Dijkstra	sylvia.dijkstra@sheerenloo.nl	0683081728 	nieuw	2025-04-06 18:48:14.993723	2025-11-02 15:04:11.346757	Begeleider	2.5 KM	Nee		t	t	2025-04-10 11:15:46.282	\N	\N	\N	f	\N	0	0
2aaffc78-4dca-4b06-9a6a-3799cbfbe67b	Noa hiddes	Klaskehiddes@gmail.com	\N	nieuw	2025-03-31 11:57:23.536871	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Nee		t	t	2025-04-01 17:11:04.313	\N	\N	\N	f	\N	0	0
200f862d-4d80-4ccd-8f06-bf795be727fb	Anneke van de Glind 	Klaskehiddes@gmail.com	\N	nieuw	2025-03-31 11:56:37.582458	2025-11-02 15:04:11.346757	Deelnemer	2.5 KM	Nee		t	t	2025-04-01 17:11:05.257	\N	\N	\N	f	\N	0	0
\.


--
-- Data for Name: album_photos; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
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
-- Data for Name: albums; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.albums (id, title, description, cover_photo_id, visible, order_number, created_at, updated_at) FROM stdin;
72831c18-4c6c-4dc8-9a2a-ace696e8996b	Voorbereidingen 	Voor werk	dbb68d91-3f6f-46c3-b751-a13f34269b79	t	3	2025-02-14 14:35:32.369+00	2025-02-14 14:35:32.37+00
ce8df963-f118-4296-9f5c-e33308dc7bfa	DKL-2024	DKL 2024	8a4d5c20-ea73-4336-8c6b-af7a197ef7c2	t	2	2024-12-26 15:17:54.268+00	2025-02-17 13:32:50.793+00
d51cff45-b958-4370-a983-51e650ffa43e	DKL 2025	De koninklijke Loop 2025!	08362e92-340a-432a-b306-153ad27ee686	t	1	2025-05-17 20:10:00+00	2025-05-17 20:10:06.643082+00
\.


--
-- Data for Name: chat_channel_participants; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.chat_channel_participants (id, channel_id, user_id, role, joined_at, last_seen_at, is_active, last_read_at) FROM stdin;
\.


--
-- Data for Name: chat_channels; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.chat_channels (id, name, description, type, created_by, created_at, updated_at, is_active, is_public) FROM stdin;
\.


--
-- Data for Name: chat_message_reactions; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.chat_message_reactions (id, message_id, user_id, emoji, created_at) FROM stdin;
\.


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.chat_messages (id, channel_id, user_id, content, message_type, file_url, file_name, file_size, reply_to_id, edited_at, created_at, updated_at, thumbnail_url) FROM stdin;
\.


--
-- Data for Name: chat_user_presence; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.chat_user_presence (user_id, status, last_seen, updated_at) FROM stdin;
\.


--
-- Data for Name: contact_antwoorden; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.contact_antwoorden (id, contact_id, verzonden_op, created_at, updated_at, tekst, verzond_op, verzond_door, email_verzonden, verzonden_door) FROM stdin;
\.


--
-- Data for Name: contact_formulieren; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.contact_formulieren (id, naam, email, bericht, status, created_at, updated_at, email_verzonden, email_verzonden_op, privacy_akkoord, behandeld_door, behandeld_op, notities, beantwoord, antwoord_tekst, antwoord_datum, antwoord_door, test_mode, antwoorden_count) FROM stdin;
78910428-f760-485d-ae57-653db478ca35	Bas heijenk 	basheijenk96@gmail.com	Hallo ik heb me op gegeven maar ik kan helaas niet sorry 	nieuw	2025-03-23 07:15:20.8851	2025-11-02 15:04:11.346757	f	\N	t	\N	\N	\N	f	\N	\N	\N	f	0
6ce2e9a3-59fd-4430-aa5c-66df48fbd695	je geheime liefde	de.konining@willem.alexander.nl	Gedeelte doneren doet het niet.	nieuw	2025-01-28 00:03:05.830054	2025-11-02 15:04:11.346757	t	2025-01-28 00:03:08.041	t	marieke@dekoninklijkeloop.nl	2025-02-05 02:23:22.97	\N	f	\N	\N	\N	f	0
\.


--
-- Data for Name: email_templates; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.email_templates (id, naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by) FROM stdin;
af67c5cb-b20b-4953-a509-3a966897a1e8	contact_admin_email	Nieuw contactformulier	<p>Er is een nieuw contactformulier ingevuld door {{.Contact.Naam}}.</p><p>Email: {{.Contact.Email}}</p><p>Bericht: {{.Contact.Bericht}}</p>	Email die naar de admin wordt gestuurd bij een nieuw contactformulier	t	2025-10-30 23:00:55.71226	2025-10-30 23:00:55.71226	d6e1fc60-4bea-46b4-8095-4ef57220bb79
f2d90198-0e9e-4b35-ac8f-df3c829eedbf	contact_email	Bedankt voor je bericht	<p>Beste {{.Contact.Naam}},</p><p>Bedankt voor je bericht. We nemen zo snel mogelijk contact met je op.</p>	Bevestigingsemail die naar de gebruiker wordt gestuurd bij een contactformulier	t	2025-10-30 23:00:55.71226	2025-10-30 23:00:55.71226	d6e1fc60-4bea-46b4-8095-4ef57220bb79
9b1adeba-577a-4e2a-8dbf-7a3744d7a9ad	aanmelding_admin_email	Nieuwe aanmelding ontvangen	<p>Er is een nieuwe aanmelding ontvangen van {{.Aanmelding.Naam}}.</p><p>Email: {{.Aanmelding.Email}}</p>	Email die naar de admin wordt gestuurd bij een nieuwe aanmelding	t	2025-10-30 23:00:55.71226	2025-10-30 23:00:55.71226	d6e1fc60-4bea-46b4-8095-4ef57220bb79
a14e9a83-dedd-4169-958b-bbbee3d1f8c6	aanmelding_email	Bedankt voor je aanmelding	<p>Beste {{.Aanmelding.Naam}},</p><p>Bedankt voor je aanmelding. We hebben je aanmelding ontvangen en zullen deze zo snel mogelijk verwerken.</p>	Bevestigingsemail die naar de gebruiker wordt gestuurd bij een aanmelding	t	2025-10-30 23:00:55.71226	2025-10-30 23:00:55.71226	d6e1fc60-4bea-46b4-8095-4ef57220bb79
\.


--
-- Data for Name: gebruikers; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, laatste_login, created_at, updated_at, newsletter_subscribed, role_id) FROM stdin;
d6e1fc60-4bea-46b4-8095-4ef57220bb79	Admin	admin@dekoninklijkeloop.nl	$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy	admin	t	\N	2025-10-30 23:00:55.71226	2025-10-30 23:00:55.71226	f	\N
\.


--
-- Data for Name: incoming_emails; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.incoming_emails (id, message_id, "from", "to", subject, body, content_type, received_at, uid, account_type, is_processed, processed_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: migraties; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.migraties (id, versie, naam, toegepast) FROM stdin;
1	1.0.0	Initiële database setup	2025-10-30 23:00:55.512742+00
2	1.0.1	Initiële data seeding	2025-10-30 23:00:55.71226+00
3	1.0.2	Update schema om overeen te komen met Go models	2025-10-30 23:00:55.723982+00
4	1.0.3	Create incoming_emails table	2025-10-30 23:00:55.760846+00
45	1.14	Toevoegen van extra aanmeldingen maart 2025	2025-10-30 23:03:36.679575+00
46	1.15	Toevoegen van aanmeldingen april 2025	2025-10-30 23:03:36.684619+00
47	1.16.0	Add chat tables	2025-10-30 23:03:36.690332+00
48	1.17.0	Add is_public to chat_channels	2025-10-30 23:03:36.869092+00
49	1.18.0	Add last_read_at to chat_channel_participants	2025-10-30 23:03:36.874597+00
50	1.20.0	Create RBAC tables for flexible permission management	2025-10-30 23:03:36.91198+00
51	1.21.0	Seed initial RBAC data based on current system	2025-10-30 23:03:37.080141+00
52	1.22.0	Assign admin role to existing admin user	2025-10-30 23:03:37.095624+00
53	1.23.0	Add staff access permission and assign to admin/staff roles	2025-10-30 23:03:37.102+00
54	1.24.0	Add missing permissions for Photos, Albums, Partners, Sponsors, Videos	2025-10-30 23:03:37.106674+00
55	1.25.0	Assign new permissions for Photos, Albums, Partners, Sponsors, Videos to admin role	2025-10-30 23:03:37.11141+00
56	1.26.0	Assign read permissions for Photos, Albums, Partners, Sponsors, Videos to staff role	2025-10-30 23:03:37.118101+00
57	1.27.0	Assign staff role to jeffrey@dekoninklijkeloop.nl	2025-10-30 23:03:37.123261+00
58	1.29.0	Add thumbnail_url to chat_messages	2025-10-30 23:03:37.196549+00
59	1.30.0	Create uploaded_images table	2025-10-30 23:03:37.201888+00
60	1.34.0	Add steps permissions for deelnemer/begeleider/vrijwilliger roles	2025-10-30 23:03:37.607828+00
61	1.45.0	Add steps permissions for RBAC system	2025-10-30 23:03:37.702061+00
357	1.47.0	Performance optimizations: FK indexes, compound indexes, partial indexes, and FTS	2025-11-01 22:52:54.893135+00
358	1.48.0	Advanced optimizations: triggers, constraints, materialized views, data quality	2025-11-01 22:52:55.146137+00
\.


--
-- Data for Name: newsletters; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.newsletters (id, subject, content, sent_at, batch_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.notifications (id, type, priority, title, message, sent, sent_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: partners; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.partners (id, name, description, logo, website, tier, since, visible, order_number, created_at, updated_at) FROM stdin;
26f11d04-c0da-4755-b2e1-5fcb5e9887d4	Accress	Accres beheert in Apeldoorn ruim 60 locaties, waaronder sporthallen, wijkcentra, zwembaden, kinderboerderijen en een stadspark.\r\n\r\nZij helpen ond bij het halen van ons doel!	https://res.cloudinary.com/dgfuv7wif/image/upload/v1744388421/accres_logo_ochsmg.jpg	https://www.accres.nl/	bronze	2025-04-01	t	5	2025-04-11 16:21:40+00	2025-04-11 16:45:10.227784+00
510d8e4d-6a7d-4ab3-b311-17b32df7f01f	Apeldoorn	Apeldoorn ondersteunt ons in ons doel.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895194/nw1qxouzupzsshckzkab.png	https://www.apeldoorn.nl/	bronze	2024-12-22	t	0	2024-12-22 19:19:55.414207+00	2025-02-26 21:10:59.040907+00
5e8b8390-e637-4c6d-80d8-25ff4155d5a9	Sheeren Loo	Samen met bewoners van SheerenLoo wordt deze loop georganiseerd.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734894570/mtvucaouruenat2cllsi.png	https://www.sheerenloo.nl/	silver	2024-12-22	t	0	2024-12-22 19:09:31.099813+00	2025-01-07 14:11:39.051238+00
9c311d0f-4db7-4da9-9d87-a836e48078cd	Liliane Fonds	Samen maken we ons sterk	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734893709/qsygajx2tdxxbqbfyurr.png	https://www.lilianefonds.nl/	bronze	2024-12-22	t	0	2024-12-22 18:55:09.823207+00	2025-01-08 19:43:58.464461+00
eda49448-3db9-4db5-9d01-a7331879f5b5	De Grote Kerk	De grote kerk ondersteund ons al vanaf het begin. Hier is het allemaal begonnen.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895146/ri4vclttn4nn2wh53wj0.jpg	https://www.grotekerkapeldoorn.nl/	gold	2024-12-22	t	0	2024-12-22 19:19:07.18416+00	2025-04-11 16:46:10.191766+00
\.


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.permissions (id, resource, action, description, is_system_permission, created_at, updated_at) FROM stdin;
219f4be7-7f3a-4d5f-b31b-9f2ebf38132b	contact	read	Contactformulieren bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
319e0b90-825e-4124-8495-079a3d37cc72	contact	write	Contactformulieren bewerken (status, notities, antwoorden)	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
9a9567fa-8a47-4368-9b94-a71aa037cf1a	contact	delete	Contactformulieren verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
ab3a99d3-98fa-4f13-a97f-5862a2d29a2d	aanmelding	read	Aanmeldingen bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
439d5def-b842-472f-9d85-30a0d20cd2e4	aanmelding	write	Aanmeldingen bewerken (status, notities, antwoorden)	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
70396b04-38ac-4b99-98de-b937d6157206	aanmelding	delete	Aanmeldingen verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
1f296357-8404-49f0-a79f-d2030da5a61e	newsletter	read	Nieuwsbrieven bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
0cab6305-bb09-47b4-aefa-ccb261bc6767	newsletter	write	Nieuwsbrieven aanmaken/bewerken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
7f6300a1-61b2-4660-990d-491ff255ab0c	newsletter	send	Nieuwsbrieven verzenden	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
0b876cf8-46bb-46a7-94b1-4f3f34397d1b	newsletter	delete	Nieuwsbrieven verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
33dd39df-4027-4eb3-a929-514d63736dd8	email	read	Inkomende emails bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
6bad89ff-57d5-430f-b062-0081feadef3d	email	write	Emails bewerken (markeren als verwerkt)	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
3ab1f863-4fee-4967-9167-781c884e02af	email	delete	Emails verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
fc67abc3-f0bb-4318-805d-e7f24258f27f	email	fetch	Nieuwe emails ophalen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
55ab421d-43ea-4aac-a040-9011b136f3cc	admin_email	send	Emails verzenden namens admin	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
3984003d-9061-44a9-be9f-67bef9471476	user	read	Gebruikers bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
c45bd8ae-7e7b-4823-80e3-8d7e2006cb29	user	write	Gebruikers aanmaken/bewerken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
961bd12d-8f1d-4a68-9bd6-c0aa17f04ddb	user	delete	Gebruikers verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
b5e94f0d-361f-417a-8895-270f3ba65d23	user	manage_roles	Gebruikersrollen beheren	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
a2ea86f2-f47b-4fa5-8c20-6850c77c472b	chat	read	Chat kanalen en berichten bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
8a4679c4-f899-4ae7-9822-5536cee24939	chat	write	Berichten verzenden	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
e256169e-3a87-43aa-979b-4bcf08f71783	chat	manage_channel	Kanalen aanmaken/beheren	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
70584198-796b-43d9-a99c-286fba67a0cf	chat	moderate	Berichten modereren (bewerken/verwijderen)	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
890de49f-a488-45b1-84c5-20e7ab3e04b8	notification	read	Notificaties bekijken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
d1b40ed2-4eda-49c2-af61-4a694e08e25a	notification	write	Notificaties aanmaken	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
91ff2e12-0cc5-47a3-ba1a-0b5db88690b1	notification	delete	Notificaties verwijderen	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
580538c5-43e1-4e76-bc62-37194d883eab	system	admin	Volledige systeemtoegang	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
9c27afaf-ef5e-4ddd-9ae8-69550d5a1e0c	admin	access	Volledige admin toegang	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
2465d366-5230-4084-828d-b8375c208e8f	staff	access	Toegang tot staff functies	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00
d004f277-329c-4f1c-bfe2-031ceba2fda7	photo	read	Foto's bekijken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
9ba5d323-365c-49a4-8def-7366842e952f	photo	write	Foto's uploaden/bewerken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
03a7641e-4235-4d3c-af82-d387b66bd9df	photo	delete	Foto's verwijderen	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
6320b222-bcd5-4e51-b967-1671ffa665be	album	read	Albums bekijken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
b026fdd4-7fc1-40ea-8706-228a21be189a	album	write	Albums aanmaken/bewerken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
a6cc7ce4-f9a1-4230-8e11-3d4ce59f7bdc	album	delete	Albums verwijderen	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
0a68fbdf-9eb5-4c52-a0b5-7653d50b01e5	partner	read	Partners bekijken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
7e0bb2c0-d381-48e0-9c5d-cbae83ed6514	partner	write	Partners aanmaken/bewerken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
4069d078-524c-4866-a0c1-789adcaea177	partner	delete	Partners verwijderen	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
6ae06535-94b9-4b3d-b145-4c88610a159a	sponsor	read	Sponsors bekijken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
fdfb4016-7688-4b59-b0e2-1545bdca957b	sponsor	write	Sponsors aanmaken/bewerken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
1562fece-15c0-447f-8a0b-c42ad1cc3e98	sponsor	delete	Sponsors verwijderen	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
8af3b606-781d-4739-9660-8fb4d35e9a5a	video	read	Video's bekijken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
eceae009-c048-4ee1-9243-b0ece6755256	video	write	Video's uploaden/bewerken	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
7767b791-c46c-4dfc-b6d4-9f3796a39aaa	video	delete	Video's verwijderen	t	2025-10-30 23:03:37.106674+00	2025-10-30 23:03:37.106674+00
36e849b8-89f3-4b6a-a19c-0b4fb6a1f993	radio_recording	read	Radio opnames bekijken	t	2025-10-30 23:03:37.600191+00	2025-10-30 23:03:37.600191+00
26d94e66-9df2-43de-8e02-de00cfa11ace	radio_recording	write	Radio opnames aanmaken/bewerken	t	2025-10-30 23:03:37.600191+00	2025-10-30 23:03:37.600191+00
6a3b21dc-5d6d-4763-9f10-d46ced7c6072	radio_recording	delete	Radio opnames verwijderen	t	2025-10-30 23:03:37.600191+00	2025-10-30 23:03:37.600191+00
3ca971f3-58f0-48b5-ba52-d14fa3ee1172	steps	read	Eigen stappen en dashboard bekijken	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
1a218ec2-abcd-43f9-809b-0d0802dd8cef	steps	write	Eigen stappen bijwerken	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
06326093-4711-421c-8b92-3b09f1dbdbda	steps	read_total	Totaal aantal stappen van alle deelnemers bekijken	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
c991528d-e8d7-48bc-8fd5-5a46bf9d2aff	steps	read_all	Alle deelnemers stappen bekijken (admin/staff)	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
f0c856d4-f9f6-4b4a-bdb1-4c1d841d6071	steps	write_all	Alle deelnemers stappen bijwerken (admin/staff)	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
ed1777fa-ca93-40cc-bc78-9dd4b6d66eb9	steps	manage	Volledige steps beheer (route funds, etc.)	t	2025-10-30 23:03:37.607828+00	2025-10-30 23:03:37.607828+00
bf2ddfa1-c5d0-4cc1-85d0-1d0d2b8f0720	program_schedule	read	Programma bekijken	t	2025-10-30 23:03:37.620183+00	2025-10-30 23:03:37.620183+00
1dc47a93-d6b7-409a-8c0f-3cd2af99deaa	program_schedule	write	Programma aanmaken/bewerken	t	2025-10-30 23:03:37.620183+00	2025-10-30 23:03:37.620183+00
87bf9f61-067e-43f9-bb91-96b027639c67	program_schedule	delete	Programma verwijderen	t	2025-10-30 23:03:37.620183+00	2025-10-30 23:03:37.620183+00
6a6c915e-36b7-4ba7-9340-489ec48d11a5	social_embed	read	Social embeds bekijken	t	2025-10-30 23:03:37.625759+00	2025-10-30 23:03:37.625759+00
512ba706-647c-40ae-8260-6cfba44cc298	social_embed	write	Social embeds aanmaken/bewerken	t	2025-10-30 23:03:37.625759+00	2025-10-30 23:03:37.625759+00
bb6d9185-a4d3-4d99-b249-3f682b081f0a	social_embed	delete	Social embeds verwijderen	t	2025-10-30 23:03:37.625759+00	2025-10-30 23:03:37.625759+00
361cfb50-ded3-4a9a-8fed-8a2b4fd9dd7c	social_link	read	Social links bekijken	t	2025-10-30 23:03:37.630867+00	2025-10-30 23:03:37.630867+00
67699394-54c1-4d88-9916-c5bc6225e396	social_link	write	Social links aanmaken/bewerken	t	2025-10-30 23:03:37.630867+00	2025-10-30 23:03:37.630867+00
84292fef-e612-41e6-a567-ba15c7adb1e8	social_link	delete	Social links verwijderen	t	2025-10-30 23:03:37.630867+00	2025-10-30 23:03:37.630867+00
347379c7-3461-468f-a821-e70e913f05e6	under_construction	read	Under construction bekijken	t	2025-10-30 23:03:37.635707+00	2025-10-30 23:03:37.635707+00
38ac3e82-1bc2-4f92-aa71-fd53aee74d4a	under_construction	write	Under construction aanmaken/bewerken	t	2025-10-30 23:03:37.635707+00	2025-10-30 23:03:37.635707+00
5183e3db-d61b-4376-8681-6502912ff7e1	under_construction	delete	Under construction verwijderen	t	2025-10-30 23:03:37.635707+00	2025-10-30 23:03:37.635707+00
3dda4703-1938-4ff6-932f-daceb622f6e9	title_section	read	Title section bekijken	t	2025-10-30 23:03:37.661975+00	2025-10-30 23:03:37.661975+00
1a2fea42-0aad-4621-8f6d-edc7c7894c45	title_section	write	Title section aanmaken/bewerken	t	2025-10-30 23:03:37.661975+00	2025-10-30 23:03:37.661975+00
1ea9ca57-5b1d-4e63-a052-5c58b5482489	title_section	delete	Title section verwijderen	t	2025-10-30 23:03:37.661975+00	2025-10-30 23:03:37.661975+00
\.


--
-- Data for Name: photos; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
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
-- Data for Name: program_schedule; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
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
-- Data for Name: radio_recordings; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.radio_recordings (id, title, description, date, audio_url, thumbnail_url, visible, order_number, created_at, updated_at) FROM stdin;
a6e73425-6af2-4bbd-84f8-67b16d195f99	De koninklijke Loop 2025 uitzending!	Luister naar het live radioverslag van De Koninklijke Loop 2025, uitgezonden op RTV Apeldoorn, met interviews met de organisatie!\r\n\r\n	14 mei 2025	https://res.cloudinary.com/dgfuv7wif/video/upload/v1747733438/DKLRTV2025_dpdydc.wav	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	t	1	2025-05-20 09:32:01+00	2025-05-20 09:33:27.540519+00
c2d84ca9-97a9-4990-9ee4-1fe1718a8c5b	Radioverslag Koninklijke Loop 2024 (RTV Apeldoorn)	Luister naar het live radioverslag van De Koninklijke Loop 2024, uitgezonden op RTV Apeldoorn, met interviews en sfeerimpressies.	15 mei 2024	https://res.cloudinary.com/dgfuv7wif/video/upload/v1714042357/matinee_1_nbm0ph.wav	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	t	2	2025-04-08 19:50:35.125804+00	2025-05-20 09:32:30.549188+00
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.refresh_tokens (id, user_id, token, expires_at, created_at, revoked_at, is_revoked) FROM stdin;
\.


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.role_permissions (id, role_id, permission_id, assigned_at, assigned_by) FROM stdin;
a09d349c-aefc-43fe-8a1c-4843892759f3	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	219f4be7-7f3a-4d5f-b31b-9f2ebf38132b	2025-10-30 23:03:37.080141+00	\N
5dc37610-aec4-4bb4-ad05-c655c49f0dfd	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	319e0b90-825e-4124-8495-079a3d37cc72	2025-10-30 23:03:37.080141+00	\N
a4716a10-a62f-4d2d-aedb-b79ce0fef2c6	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	9a9567fa-8a47-4368-9b94-a71aa037cf1a	2025-10-30 23:03:37.080141+00	\N
94b0faba-e756-4fa2-a50e-9901c776ef3f	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	ab3a99d3-98fa-4f13-a97f-5862a2d29a2d	2025-10-30 23:03:37.080141+00	\N
d42eb209-589d-4e5f-a75d-d9c2432a8d44	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	439d5def-b842-472f-9d85-30a0d20cd2e4	2025-10-30 23:03:37.080141+00	\N
3af84f40-f951-4766-ac39-6ff6e672bf5e	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	70396b04-38ac-4b99-98de-b937d6157206	2025-10-30 23:03:37.080141+00	\N
69d26118-1228-4948-87bf-b847e3774976	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1f296357-8404-49f0-a79f-d2030da5a61e	2025-10-30 23:03:37.080141+00	\N
d6bab5b5-af12-474d-b800-00336b76bb29	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	0cab6305-bb09-47b4-aefa-ccb261bc6767	2025-10-30 23:03:37.080141+00	\N
b2fab645-b654-470b-b7fb-467e3506634b	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	7f6300a1-61b2-4660-990d-491ff255ab0c	2025-10-30 23:03:37.080141+00	\N
48ac954d-d01f-426d-ac94-2533ea949717	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	0b876cf8-46bb-46a7-94b1-4f3f34397d1b	2025-10-30 23:03:37.080141+00	\N
14c3320c-b482-45c9-9c1a-8bb0fc4c5aee	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	33dd39df-4027-4eb3-a929-514d63736dd8	2025-10-30 23:03:37.080141+00	\N
d410a664-bb8e-4d34-9eef-6945d9be5708	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	6bad89ff-57d5-430f-b062-0081feadef3d	2025-10-30 23:03:37.080141+00	\N
c68197d5-0cbc-4e30-b3ab-5c5e79092eda	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	3ab1f863-4fee-4967-9167-781c884e02af	2025-10-30 23:03:37.080141+00	\N
3e20a304-ae86-447b-b9fc-b63f389323b6	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	fc67abc3-f0bb-4318-805d-e7f24258f27f	2025-10-30 23:03:37.080141+00	\N
7a1e9dc3-81d4-4cae-853b-2e1d7bddaf94	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	55ab421d-43ea-4aac-a040-9011b136f3cc	2025-10-30 23:03:37.080141+00	\N
c84ef6ad-b917-460d-b8bf-3996814609cc	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	3984003d-9061-44a9-be9f-67bef9471476	2025-10-30 23:03:37.080141+00	\N
ef171753-1cc4-4661-bf28-79c75e1e730c	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	c45bd8ae-7e7b-4823-80e3-8d7e2006cb29	2025-10-30 23:03:37.080141+00	\N
e5091cb0-243e-4715-a219-893d09899d85	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	961bd12d-8f1d-4a68-9bd6-c0aa17f04ddb	2025-10-30 23:03:37.080141+00	\N
056c0355-855a-41cc-9f32-b4e0b83815fd	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	b5e94f0d-361f-417a-8895-270f3ba65d23	2025-10-30 23:03:37.080141+00	\N
ec6f3ffa-6324-415e-9be7-f48325763fa7	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
853d730e-1c7b-4733-b84e-9948d57f9a61	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	8a4679c4-f899-4ae7-9822-5536cee24939	2025-10-30 23:03:37.080141+00	\N
5a8c5c3d-2440-4c71-84b2-2ad7d5987ef3	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	e256169e-3a87-43aa-979b-4bcf08f71783	2025-10-30 23:03:37.080141+00	\N
86c3084c-85dc-4e9d-a291-7b8d08f173cf	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	70584198-796b-43d9-a99c-286fba67a0cf	2025-10-30 23:03:37.080141+00	\N
7d9b4ce3-8501-45a8-8899-dca4de57d57f	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	890de49f-a488-45b1-84c5-20e7ab3e04b8	2025-10-30 23:03:37.080141+00	\N
d6c4e4c9-e2e1-4478-8e84-3155fd0c0211	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	d1b40ed2-4eda-49c2-af61-4a694e08e25a	2025-10-30 23:03:37.080141+00	\N
b3b6e994-2acd-40ac-81cb-3c88ab03ce5a	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	91ff2e12-0cc5-47a3-ba1a-0b5db88690b1	2025-10-30 23:03:37.080141+00	\N
703b971d-f22b-4f2b-ac34-69650f2e7802	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	580538c5-43e1-4e76-bc62-37194d883eab	2025-10-30 23:03:37.080141+00	\N
d3337d28-6af5-453a-8c3c-6cf36014c657	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	9c27afaf-ef5e-4ddd-9ae8-69550d5a1e0c	2025-10-30 23:03:37.080141+00	\N
99ee014d-ef6b-474d-800d-57334cf753b0	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	2465d366-5230-4084-828d-b8375c208e8f	2025-10-30 23:03:37.080141+00	\N
4d522af7-846e-4e2c-8562-13cd1a9b466e	81b071ca-9637-4983-988c-02c108b55e71	219f4be7-7f3a-4d5f-b31b-9f2ebf38132b	2025-10-30 23:03:37.080141+00	\N
78a0d48d-4698-4291-84e4-aaa261fd6294	81b071ca-9637-4983-988c-02c108b55e71	ab3a99d3-98fa-4f13-a97f-5862a2d29a2d	2025-10-30 23:03:37.080141+00	\N
2c703b17-311d-4a1d-983f-361e070b1052	81b071ca-9637-4983-988c-02c108b55e71	1f296357-8404-49f0-a79f-d2030da5a61e	2025-10-30 23:03:37.080141+00	\N
d7c61db0-0b95-455a-9564-d4ab9be692f0	81b071ca-9637-4983-988c-02c108b55e71	33dd39df-4027-4eb3-a929-514d63736dd8	2025-10-30 23:03:37.080141+00	\N
d5250920-3bfb-42df-b9d5-2353c4d72253	81b071ca-9637-4983-988c-02c108b55e71	3984003d-9061-44a9-be9f-67bef9471476	2025-10-30 23:03:37.080141+00	\N
f61f7ec4-49f9-484f-8831-889c88ac4625	81b071ca-9637-4983-988c-02c108b55e71	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
c68c7893-72d4-4a4b-aafd-af9678745d85	81b071ca-9637-4983-988c-02c108b55e71	890de49f-a488-45b1-84c5-20e7ab3e04b8	2025-10-30 23:03:37.080141+00	\N
722709ba-d932-490c-beaf-999f9e2f30bf	81b071ca-9637-4983-988c-02c108b55e71	2465d366-5230-4084-828d-b8375c208e8f	2025-10-30 23:03:37.080141+00	\N
cac166ec-bf87-4eb1-afb8-febcf40e5cab	fddb1fd2-14c1-4c8e-b1ce-ca56524199a4	e256169e-3a87-43aa-979b-4bcf08f71783	2025-10-30 23:03:37.080141+00	\N
734446d8-8d79-49d1-8b7e-3b25ce6f9aba	fddb1fd2-14c1-4c8e-b1ce-ca56524199a4	70584198-796b-43d9-a99c-286fba67a0cf	2025-10-30 23:03:37.080141+00	\N
3fe620e3-4184-4183-b3e3-2ad6b2dd7348	fddb1fd2-14c1-4c8e-b1ce-ca56524199a4	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
98c9ac40-ec09-4c34-9915-3e08590dc209	fddb1fd2-14c1-4c8e-b1ce-ca56524199a4	8a4679c4-f899-4ae7-9822-5536cee24939	2025-10-30 23:03:37.080141+00	\N
25eadd32-7aae-402a-8d4f-5de69d26cc33	9c90cb3a-9a28-4c3c-9495-9118c4d30364	70584198-796b-43d9-a99c-286fba67a0cf	2025-10-30 23:03:37.080141+00	\N
73f64c01-5a10-4bef-8962-cefc4a2de5d5	9c90cb3a-9a28-4c3c-9495-9118c4d30364	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
62070c75-9699-4549-9d39-0dced750af07	9c90cb3a-9a28-4c3c-9495-9118c4d30364	8a4679c4-f899-4ae7-9822-5536cee24939	2025-10-30 23:03:37.080141+00	\N
6a7ef4cd-b8eb-4e23-8d06-ce972aacdcb4	af58c072-9fe0-4fce-a797-001ffbf52955	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
64c00193-6264-405e-b3cf-3305dfbb4bf3	af58c072-9fe0-4fce-a797-001ffbf52955	8a4679c4-f899-4ae7-9822-5536cee24939	2025-10-30 23:03:37.080141+00	\N
ebcb4674-00ba-4e24-9fb2-7172a57198d9	5d3eedbb-f1fc-4ef5-8e91-03abcc977bc8	a2ea86f2-f47b-4fa5-8c20-6850c77c472b	2025-10-30 23:03:37.080141+00	\N
bf520c86-5b92-4ff8-ba83-9d064a560195	5d3eedbb-f1fc-4ef5-8e91-03abcc977bc8	8a4679c4-f899-4ae7-9822-5536cee24939	2025-10-30 23:03:37.080141+00	\N
19cfa8a9-ba5d-4231-9ea2-e1fa40617d45	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	d004f277-329c-4f1c-bfe2-031ceba2fda7	2025-10-30 23:03:37.11141+00	\N
f8bf7351-d386-4f96-af35-41102305ef13	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	9ba5d323-365c-49a4-8def-7366842e952f	2025-10-30 23:03:37.11141+00	\N
a47a4fa9-757b-4fe3-86a3-e2e615da47fa	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	03a7641e-4235-4d3c-af82-d387b66bd9df	2025-10-30 23:03:37.11141+00	\N
675ad3e2-620b-42b8-816b-ad2dbf27dbc0	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	6320b222-bcd5-4e51-b967-1671ffa665be	2025-10-30 23:03:37.11141+00	\N
e8c3f9f8-564b-4d06-93f9-a824675dd5e6	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	b026fdd4-7fc1-40ea-8706-228a21be189a	2025-10-30 23:03:37.11141+00	\N
3d71061e-5ee0-42d1-9869-9c05a2ff58d6	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	a6cc7ce4-f9a1-4230-8e11-3d4ce59f7bdc	2025-10-30 23:03:37.11141+00	\N
9b6a384c-313d-4746-acd3-b2d1c7c78ff9	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	0a68fbdf-9eb5-4c52-a0b5-7653d50b01e5	2025-10-30 23:03:37.11141+00	\N
6a7b3d68-0d05-4fdc-b12a-d8f7ff4c1248	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	7e0bb2c0-d381-48e0-9c5d-cbae83ed6514	2025-10-30 23:03:37.11141+00	\N
cd9801ff-55f9-4ffb-a193-f0472da64e2a	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	4069d078-524c-4866-a0c1-789adcaea177	2025-10-30 23:03:37.11141+00	\N
32b2b2cd-5c44-468e-a8c8-30b7db5cc8a5	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	6ae06535-94b9-4b3d-b145-4c88610a159a	2025-10-30 23:03:37.11141+00	\N
348f292e-3f39-43eb-96cc-57689ebd1d92	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	fdfb4016-7688-4b59-b0e2-1545bdca957b	2025-10-30 23:03:37.11141+00	\N
7c95fd3f-45e4-4c7b-815a-66405cf52994	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1562fece-15c0-447f-8a0b-c42ad1cc3e98	2025-10-30 23:03:37.11141+00	\N
f4a86463-c5c9-42ac-a105-7f4fc3639ed5	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	8af3b606-781d-4739-9660-8fb4d35e9a5a	2025-10-30 23:03:37.11141+00	\N
ea961e9b-5249-433a-a3b1-9b45ae8df5f5	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	eceae009-c048-4ee1-9243-b0ece6755256	2025-10-30 23:03:37.11141+00	\N
33519ee5-acef-4b8c-ae19-59671b30f8a4	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	7767b791-c46c-4dfc-b6d4-9f3796a39aaa	2025-10-30 23:03:37.11141+00	\N
5b052a37-f8e3-4819-8951-73f3b49f72ac	81b071ca-9637-4983-988c-02c108b55e71	6320b222-bcd5-4e51-b967-1671ffa665be	2025-10-30 23:03:37.118101+00	\N
3176791b-40e2-4424-9ad9-3d09a5072f84	81b071ca-9637-4983-988c-02c108b55e71	0a68fbdf-9eb5-4c52-a0b5-7653d50b01e5	2025-10-30 23:03:37.118101+00	\N
953c480f-4123-47d8-b5d9-a1ce312f5782	81b071ca-9637-4983-988c-02c108b55e71	d004f277-329c-4f1c-bfe2-031ceba2fda7	2025-10-30 23:03:37.118101+00	\N
2fbceb8c-7b8d-4e6f-84dc-bf0305dbb2b9	81b071ca-9637-4983-988c-02c108b55e71	6ae06535-94b9-4b3d-b145-4c88610a159a	2025-10-30 23:03:37.118101+00	\N
a2d12973-3072-4339-ae99-56351bbd1b26	81b071ca-9637-4983-988c-02c108b55e71	8af3b606-781d-4739-9660-8fb4d35e9a5a	2025-10-30 23:03:37.118101+00	\N
9acc3f12-ad4e-4d20-a87e-a2ed81faa311	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	36e849b8-89f3-4b6a-a19c-0b4fb6a1f993	2025-10-30 23:03:37.600191+00	\N
b44a0b5d-b06c-49b2-a40f-6cb2022af2d4	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	26d94e66-9df2-43de-8e02-de00cfa11ace	2025-10-30 23:03:37.600191+00	\N
4b54421f-72e0-4888-a38b-7b34d8f95a19	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	6a3b21dc-5d6d-4763-9f10-d46ced7c6072	2025-10-30 23:03:37.600191+00	\N
18f52283-70de-4f09-abce-6ce633e5b1bf	81b071ca-9637-4983-988c-02c108b55e71	36e849b8-89f3-4b6a-a19c-0b4fb6a1f993	2025-10-30 23:03:37.600191+00	\N
02795734-59f3-49a9-b344-f90ffa2db9d3	70ce44dc-5b00-453d-902d-241258b03474	3ca971f3-58f0-48b5-ba52-d14fa3ee1172	2025-10-30 23:03:37.607828+00	\N
fe1efcec-5831-4139-a291-8d460be58484	70ce44dc-5b00-453d-902d-241258b03474	06326093-4711-421c-8b92-3b09f1dbdbda	2025-10-30 23:03:37.607828+00	\N
fdd6f322-6e18-4319-a802-dff4e16bdf8b	70ce44dc-5b00-453d-902d-241258b03474	1a218ec2-abcd-43f9-809b-0d0802dd8cef	2025-10-30 23:03:37.607828+00	\N
0078244a-13b5-4b9f-83ae-6817a83982e9	e6b4f9ec-78bb-4f9c-8da9-dfdb648fbe2d	3ca971f3-58f0-48b5-ba52-d14fa3ee1172	2025-10-30 23:03:37.607828+00	\N
cd241d91-6eb8-42b1-80af-e201b4d4032a	e6b4f9ec-78bb-4f9c-8da9-dfdb648fbe2d	06326093-4711-421c-8b92-3b09f1dbdbda	2025-10-30 23:03:37.607828+00	\N
9c67b561-2cd2-43ac-9578-4e3ea431dee3	e6b4f9ec-78bb-4f9c-8da9-dfdb648fbe2d	1a218ec2-abcd-43f9-809b-0d0802dd8cef	2025-10-30 23:03:37.607828+00	\N
265ce23d-f69f-4bff-b29b-01801290114e	c23f6919-2728-4504-8d70-c2bc6d93eb73	3ca971f3-58f0-48b5-ba52-d14fa3ee1172	2025-10-30 23:03:37.607828+00	\N
852b3c31-9484-45d1-971f-f0506f270d75	c23f6919-2728-4504-8d70-c2bc6d93eb73	06326093-4711-421c-8b92-3b09f1dbdbda	2025-10-30 23:03:37.607828+00	\N
9269863a-301f-4bc5-8792-3c24dd605089	c23f6919-2728-4504-8d70-c2bc6d93eb73	1a218ec2-abcd-43f9-809b-0d0802dd8cef	2025-10-30 23:03:37.607828+00	\N
0cfe48a1-7563-4911-ac9b-c15f482db1f4	81b071ca-9637-4983-988c-02c108b55e71	3ca971f3-58f0-48b5-ba52-d14fa3ee1172	2025-10-30 23:03:37.607828+00	\N
5d241e77-1b74-4306-8251-e3e19f263398	81b071ca-9637-4983-988c-02c108b55e71	c991528d-e8d7-48bc-8fd5-5a46bf9d2aff	2025-10-30 23:03:37.607828+00	\N
13770c37-2dc9-4b9b-a24d-66348fcd5f34	81b071ca-9637-4983-988c-02c108b55e71	1a218ec2-abcd-43f9-809b-0d0802dd8cef	2025-10-30 23:03:37.607828+00	\N
27b99538-c453-4510-ba56-fceb19d6ee0e	81b071ca-9637-4983-988c-02c108b55e71	f0c856d4-f9f6-4b4a-bdb1-4c1d841d6071	2025-10-30 23:03:37.607828+00	\N
ea4a8c74-8549-4f54-b811-7320f69740e2	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	ed1777fa-ca93-40cc-bc78-9dd4b6d66eb9	2025-10-30 23:03:37.607828+00	\N
d610ad57-66b5-49eb-b1a4-f2a561132888	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	3ca971f3-58f0-48b5-ba52-d14fa3ee1172	2025-10-30 23:03:37.607828+00	\N
ec9e1270-70b5-4dc0-bde3-d26a2692fc0e	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	c991528d-e8d7-48bc-8fd5-5a46bf9d2aff	2025-10-30 23:03:37.607828+00	\N
e33d1c60-119b-4149-ab3f-734427fc829b	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	06326093-4711-421c-8b92-3b09f1dbdbda	2025-10-30 23:03:37.607828+00	\N
18f32831-6bb8-4025-9512-d031764b4113	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1a218ec2-abcd-43f9-809b-0d0802dd8cef	2025-10-30 23:03:37.607828+00	\N
d7dd972e-fad8-44ed-b24b-f5c5d1073e75	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	f0c856d4-f9f6-4b4a-bdb1-4c1d841d6071	2025-10-30 23:03:37.607828+00	\N
f37ffbbb-b05d-46bd-b2b5-759bd630e84c	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	87bf9f61-067e-43f9-bb91-96b027639c67	2025-10-30 23:03:37.620183+00	\N
8c740bfb-6ec3-4d09-b34e-1f110fca7019	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	bf2ddfa1-c5d0-4cc1-85d0-1d0d2b8f0720	2025-10-30 23:03:37.620183+00	\N
4bd5c30d-1d09-4df9-868d-36fa82a2048c	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1dc47a93-d6b7-409a-8c0f-3cd2af99deaa	2025-10-30 23:03:37.620183+00	\N
03137bf8-93b2-40dc-a4ab-9bcaae58c3a7	81b071ca-9637-4983-988c-02c108b55e71	bf2ddfa1-c5d0-4cc1-85d0-1d0d2b8f0720	2025-10-30 23:03:37.620183+00	\N
b012a794-d6c4-418f-911c-89a3f9abca35	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	bb6d9185-a4d3-4d99-b249-3f682b081f0a	2025-10-30 23:03:37.625759+00	\N
67757331-85e8-4b28-91ce-da31163fe7ed	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	6a6c915e-36b7-4ba7-9340-489ec48d11a5	2025-10-30 23:03:37.625759+00	\N
96da4413-cf2b-43e5-8580-418fd9e480fa	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	512ba706-647c-40ae-8260-6cfba44cc298	2025-10-30 23:03:37.625759+00	\N
5ca9f0b2-6735-4640-8844-8d392c0459de	81b071ca-9637-4983-988c-02c108b55e71	6a6c915e-36b7-4ba7-9340-489ec48d11a5	2025-10-30 23:03:37.625759+00	\N
b0cc057c-fb96-41c6-bcfd-a44a31b5d0bd	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	84292fef-e612-41e6-a567-ba15c7adb1e8	2025-10-30 23:03:37.630867+00	\N
113db998-33c0-43c5-815c-b4b0b01255f7	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	361cfb50-ded3-4a9a-8fed-8a2b4fd9dd7c	2025-10-30 23:03:37.630867+00	\N
4c16ba5c-9bee-4255-b828-1529fb96006a	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	67699394-54c1-4d88-9916-c5bc6225e396	2025-10-30 23:03:37.630867+00	\N
fe51fde8-d37d-4ea6-978e-763c72e8a009	81b071ca-9637-4983-988c-02c108b55e71	361cfb50-ded3-4a9a-8fed-8a2b4fd9dd7c	2025-10-30 23:03:37.630867+00	\N
8f7ac4b4-4096-4111-9e12-e00a7fb29eb4	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	5183e3db-d61b-4376-8681-6502912ff7e1	2025-10-30 23:03:37.635707+00	\N
283cb5bb-0a09-4150-9967-3655cb60ea56	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	347379c7-3461-468f-a821-e70e913f05e6	2025-10-30 23:03:37.635707+00	\N
e87eaf3e-753e-4092-827d-56db821c1491	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	38ac3e82-1bc2-4f92-aa71-fd53aee74d4a	2025-10-30 23:03:37.635707+00	\N
336b9c37-42ee-46f4-be3d-284ffa6d37a8	81b071ca-9637-4983-988c-02c108b55e71	347379c7-3461-468f-a821-e70e913f05e6	2025-10-30 23:03:37.635707+00	\N
55c0e4d2-fccc-411a-80cf-3d19d1dbb565	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1ea9ca57-5b1d-4e63-a052-5c58b5482489	2025-10-30 23:03:37.661975+00	\N
e8249b08-3884-4675-8187-3720514f845e	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	3dda4703-1938-4ff6-932f-daceb622f6e9	2025-10-30 23:03:37.661975+00	\N
4cdbee75-a2b6-4f73-a6f7-2d75ad11b28a	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	1a2fea42-0aad-4621-8f6d-edc7c7894c45	2025-10-30 23:03:37.661975+00	\N
062ba14d-ea4a-434b-81af-6f33b177c714	81b071ca-9637-4983-988c-02c108b55e71	3dda4703-1938-4ff6-932f-daceb622f6e9	2025-10-30 23:03:37.661975+00	\N
41b16452-3d55-4571-a108-101a748b461f	81b071ca-9637-4983-988c-02c108b55e71	ed1777fa-ca93-40cc-bc78-9dd4b6d66eb9	2025-10-30 23:03:37.702061+00	\N
ee3474e1-da8c-4a97-8e75-f502b1cc9d65	81b071ca-9637-4983-988c-02c108b55e71	06326093-4711-421c-8b92-3b09f1dbdbda	2025-10-30 23:03:37.702061+00	\N
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.roles (id, name, description, is_system_role, created_at, updated_at, created_by) FROM stdin;
ffe8d8f6-1e97-4845-bf77-62ec57d859fe	admin	Volledige beheerder met toegang tot alle functies	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
81b071ca-9637-4983-988c-02c108b55e71	staff	Ondersteunend personeel met beperkte beheerrechten	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
5d3eedbb-f1fc-4ef5-8e91-03abcc977bc8	user	Standaard gebruiker	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
fddb1fd2-14c1-4c8e-b1ce-ca56524199a4	owner	Chat kanaal eigenaar	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
9c90cb3a-9a28-4c3c-9495-9118c4d30364	chat_admin	Chat kanaal beheerder	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
af58c072-9fe0-4fce-a797-001ffbf52955	member	Chat kanaal lid	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
70ce44dc-5b00-453d-902d-241258b03474	deelnemer	Evenement deelnemer	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
e6b4f9ec-78bb-4f9c-8da9-dfdb648fbe2d	begeleider	Evenement begeleider	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
c23f6919-2728-4504-8d70-c2bc6d93eb73	vrijwilliger	Evenement vrijwilliger	t	2025-10-30 23:03:37.080141+00	2025-10-30 23:03:37.080141+00	\N
\.


--
-- Data for Name: route_funds; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.route_funds (id, route, amount, created_at, updated_at) FROM stdin;
5b408dd3-2705-4f76-aeac-e99167859862	6 KM	50	2025-10-30 23:03:37.710233+00	2025-10-30 23:03:37.710233+00
5a9e2407-8b4c-431e-99f5-a1b2e55f9934	10 KM	75	2025-10-30 23:03:37.710233+00	2025-10-30 23:03:37.710233+00
7c2bfcbe-3ba3-4df5-9aed-ad9f5fcc8ffd	15 KM	100	2025-10-30 23:03:37.710233+00	2025-10-30 23:03:37.710233+00
fec468ed-6d6d-47b2-a8ce-0f468480ec11	20 KM	125	2025-10-30 23:03:37.710233+00	2025-10-30 23:03:37.710233+00
\.


--
-- Data for Name: social_embeds; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.social_embeds (id, platform, embed_code, order_number, visible, created_at, updated_at) FROM stdin;
5709a899-ee12-4883-8b57-3a4d0e8543a6	instagram	<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" data-instgrm-version="14" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;">Dit bericht op Instagram bekijken</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center; margin-bottom: 24px;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 224px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 144px;"></div></div></a><p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;"><a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none;" target="_blank">Een bericht gedeeld door Koninklijke Loop (@koninklijkeloop)</a></p></div></blockquote>\r\n<script async src="//www.instagram.com/embed.js"></script>	1	t	2024-12-22 19:30:53.366424+00	2024-12-22 19:30:53.366424+00
ee8d1152-2fd4-464d-82a8-76c02ad56ed9	facebook	<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fpermalink.php%3Fstory_fbid%3Dpfbid02XNU75Y2gMxWhvsVQar7oaM98GvMLLryXQVMTjxnBkEg6e6imJ8ecgoEF9SrTVJDpl%26id%3D61556315443279&show_text=true&width=500" width="500" height="737" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true" allow="autoplay; clipboard-write; encrypted-media; picture-in-picture; web-share"></iframe>	2	t	2024-12-22 19:30:53.366424+00	2024-12-22 19:30:53.366424+00
\.


--
-- Data for Name: social_links; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.social_links (id, platform, url, bg_color_class, icon_color_class, order_number, visible, created_at, updated_at) FROM stdin;
1de3ed0c-9bf3-4924-8d76-72045cc1c0ec	instagram	https://www.instagram.com/koninklijkeloop	\N	\N	2	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
29988672-0d98-4083-9908-18f6bbe34f3f	linkedin	https://www.linkedin.com/company/koninklijkeloop	\N	\N	4	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
dc917a65-6bb4-45cc-a1dc-890b1cdf1f5b	facebook	https://www.facebook.com/koninklijkeloop	\N	\N	1	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
f713d59e-e8a8-40d2-bb11-fc64f89b9ae9	youtube	https://www.youtube.com/@koninklijkeloop	\N	\N	3	t	2024-12-22 19:49:18.04807+00	2024-12-22 19:49:18.04807+00
\.


--
-- Data for Name: sponsors; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.sponsors (id, name, description, logo_url, website_url, order_number, is_active, created_at, updated_at, visible) FROM stdin;
484576a1-2a60-4201-b582-e1f3754ab12a	3x3 Anders	3x3 Anders is een zorgbemiddelingsbureau gespecialiseerd in het matchen van zorgaanbieders met gekwalificeerde zorgprofessionals.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166671/3x3anderslogo_itwm3g.webp	https://3x3anders.nl/	4	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
6408a640-cca1-4aaa-b845-04888f62ccec	Sterk In Vloeren	De website van Sterk In Vloeren biedt een uitgebreid assortiment aan vloeren, waaronder laminaat, PVC-vloeren en tapijt. Ze benadrukken heldere afspraken en hanteren all-in prijzen.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/SterkinVloerenLOGO_zrdofb.webp	https://sterkinvloeren.nl/	1	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
6acd1b1e-8fed-4c8b-89cb-85eee9053536	Beeldpakker	Johan Groot Jebbink, een fotograaf met meer dan tien jaar ervaring, gespecialiseerd in portretfotografie. Actief in Ermelo en internationaal.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166670/BeeldpakkerLogo_wijjmq.webp	https://beeldpakker.nl/	2	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
88cec1c1-3f08-4a6f-9623-c71605fe35b5	Bas Visual Story Telling	BAS Visual Storytelling, heeft een passie voor content en verhalen. Mijn hobby is uitgegroeid tot een eigen onderneming in het vastleggen van verhalen. Bij BAS Visual Storytelling laten we verhalen niet verstoffen op de plank, maar brengen ze tot leven! Waar ik ga of sta, mijn camera's gaan met mij mee, leg de mooiste beelden haarscherp vast en breng jouw verhaal tot leven. Dus vertel eens, 'wat is jouw verhaal?'\r\n\r\n	https://res.cloudinary.com/dgfuv7wif/image/upload/v1746017513/krqjbwwerv9hs6hyrhcy.png	https://basvisualstorytelling.nl/	5	t	2025-04-30 12:51:53.997413+00	2025-05-01 09:30:16.725422+00	t
caa59f1f-65b4-442f-84d2-22cc52212dea	Mojo Dojo	Mojo Dojo is een veelzijdige studio in Rotterdam die diensten aanbiedt voor creatieve producties, waaronder muziekopnames, podcasts en livestreams.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/LogoLayout_1_iphclc.webp	https://mojodojo.studio/	3	t	2024-11-29 09:49:35+00	2025-04-30 13:07:58.478643+00	t
\.


--
-- Data for Name: title_section_content; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.title_section_content (id, event_title, event_subtitle, image_url, image_alt, detail_1_title, detail_1_description, detail_2_title, detail_2_description, detail_3_title, detail_3_description, created_at, updated_at, participant_count) FROM stdin;
550e8400-e29b-41d4-a716-446655440001	De Koninklijke Loop (DKL) 2025	Op de koninklijke weg in Apeldoorn kunnen mensen met een beperking samen wandelen tijdens dit unieke, rolstoelvriendelijke sponsorloop (DKL), samen met hun verwanten, vrijwilligers of begeleiders.	https://res.cloudinary.com/dgfuv7wif/image/upload/v1760112848/Wij_gaan_17_mei_lopen_voor_hen_3_zllxno_zoqd7z.webp	Promotiebanner De Koninklijke Loop (DKL) 2025: Wij gaan 17 mei lopen voor hen	17 mei 2025	Starttijden variëren per afstand. Zie programma.	Voor iedereen	wandelaars met of zonder beperking (rolstoelvriendelijk).	Lopen voor een goed doel	Steun het goede doel via dit unieke wandelevenement.	2025-04-16 01:31:29.48241+00	2025-10-10 16:21:36.786249+00	69
\.


--
-- Data for Name: under_construction; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.under_construction (id, is_active, title, message, footer_text, logo_url, expected_date, social_links, progress_percentage, contact_email, newsletter_enabled, created_at, updated_at) FROM stdin;
1	f	Website in onderhoud	We stomen ons klaar voor De Koninklijke Loop 2026, op dit moment is de website helaas niet bereikbaar	Bedankt voor uw geduld!	https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png	2026-01-31 18:00:00+00	[{"url": "https://twitter.com/koninklijkeloop", "platform": "Twitter"}, {"url": "https://instagram.com/koninklijkeloop", "platform": "Instagram"}, {"url": "https://www.youtube.com/@DeKoninklijkeLoop", "platform": "YouTube"}]	85	info@koninklijkeloop.nl	f	2025-09-26 17:37:22.197854+00	2025-10-09 21:00:29.391392+00
\.


--
-- Data for Name: uploaded_images; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.uploaded_images (id, user_id, public_id, url, secure_url, filename, size, mime_type, width, height, folder, thumbnail_url, deleted_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.user_roles (id, user_id, role_id, assigned_at, assigned_by, expires_at, is_active) FROM stdin;
04cc7d32-dc6e-4db0-9521-18f03632740d	d6e1fc60-4bea-46b4-8095-4ef57220bb79	ffe8d8f6-1e97-4845-bf77-62ec57d859fe	2025-10-30 23:03:37.095624+00	\N	\N	t
\.


--
-- Data for Name: verzonden_emails; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.verzonden_emails (id, ontvanger, onderwerp, inhoud, verzonden_op, status, contact_id, aanmelding_id, template_id, created_at, updated_at, fout_bericht) FROM stdin;
\.


--
-- Data for Name: videos; Type: TABLE DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COPY public.videos (id, video_id, url, title, description, thumbnail_url, visible, order_number, created_at, updated_at) FROM stdin;
14ee164b-50e3-4f59-a3e9-3a54312af9cd	q9ngqu	https://streamable.com/e/q9ngqu	De koninklijkeloop!	Preview!!!	\N	t	1	2025-03-28 21:45:53+00	2025-04-21 10:49:03.744+00
18d951d2-f5d1-4b6e-95af-a72ac5ff18ff	x8zj4k	https://streamable.com/e/x8zj4k	Promotie De Koninklijke Loop: Flyers verspreiden	Bekijk hoe vrijwilligers flyers uitdelen om mensen uit te nodigen voor het DKL wandelevenement.	\N	t	4	2025-03-03 20:25:54.251108+00	2025-03-03 20:25:54.251108+00
87502f84-91db-419f-9766-4071ade3e94f	tt6k80	https://streamable.com/e/0o2qf9	Highlights Koninklijke Loop 2024 (Wandelevenement Apeldoorn)	Herbeleef de mooiste momenten en de sfeer van De Koninklijke Loop 2024 in deze highlight video.	\N	t	2	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:42.699+00
99bbe55b-32ef-46ab-bb59-860ce92f1d58	cvfrpi	https://streamable.com/e/cvfrpi	De spannende start van de Koninklijke Loop 2024	Bekijk de start van de deelnemers aan de sponsorloop De Koninklijke Loop 2024.	\N	t	3	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:43.979+00
ac987839-edc1-468f-9080-064f894b3e5d	tt6k80	https://streamable.com/e/tt6k80	Koninklijke Loop 2024 - Hoofdevenement	Een sfeerimpressie van het hoofdevenement van de Koninklijke Loop 2024, met deelnemers, vrijwilligers en muziek.	\N	t	5	2024-12-22 20:23:06.654419+00	2024-12-26 23:25:41.341+00
\.


--
-- Name: migraties_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

SELECT pg_catalog.setval('public.migraties_id_seq', 382, true);


--
-- Name: under_construction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

SELECT pg_catalog.setval('public.under_construction_id_seq', 1, false);


--
-- Name: aanmelding_antwoorden aanmelding_antwoorden_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.aanmelding_antwoorden
    ADD CONSTRAINT aanmelding_antwoorden_pkey PRIMARY KEY (id);


--
-- Name: aanmeldingen aanmeldingen_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.aanmeldingen
    ADD CONSTRAINT aanmeldingen_pkey PRIMARY KEY (id);


--
-- Name: album_photos album_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_pkey PRIMARY KEY (id);


--
-- Name: albums albums_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_pkey PRIMARY KEY (id);


--
-- Name: chat_channel_participants chat_channel_participants_channel_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_channel_id_user_id_key UNIQUE (channel_id, user_id);


--
-- Name: chat_channel_participants chat_channel_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_pkey PRIMARY KEY (id);


--
-- Name: chat_channels chat_channels_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_channels
    ADD CONSTRAINT chat_channels_pkey PRIMARY KEY (id);


--
-- Name: chat_message_reactions chat_message_reactions_message_id_user_id_emoji_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_message_id_user_id_emoji_key UNIQUE (message_id, user_id, emoji);


--
-- Name: chat_message_reactions chat_message_reactions_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: chat_user_presence chat_user_presence_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_user_presence
    ADD CONSTRAINT chat_user_presence_pkey PRIMARY KEY (user_id);


--
-- Name: contact_antwoorden contact_antwoorden_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.contact_antwoorden
    ADD CONSTRAINT contact_antwoorden_pkey PRIMARY KEY (id);


--
-- Name: contact_formulieren contact_formulieren_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.contact_formulieren
    ADD CONSTRAINT contact_formulieren_pkey PRIMARY KEY (id);


--
-- Name: email_templates email_templates_naam_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_naam_key UNIQUE (naam);


--
-- Name: email_templates email_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_pkey PRIMARY KEY (id);


--
-- Name: gebruikers gebruikers_email_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_email_key UNIQUE (email);


--
-- Name: gebruikers gebruikers_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_pkey PRIMARY KEY (id);


--
-- Name: incoming_emails incoming_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.incoming_emails
    ADD CONSTRAINT incoming_emails_pkey PRIMARY KEY (id);


--
-- Name: incoming_emails incoming_emails_uid_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.incoming_emails
    ADD CONSTRAINT incoming_emails_uid_key UNIQUE (uid);


--
-- Name: migraties migraties_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.migraties
    ADD CONSTRAINT migraties_pkey PRIMARY KEY (id);


--
-- Name: newsletters newsletters_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.newsletters
    ADD CONSTRAINT newsletters_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: partners partners_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.partners
    ADD CONSTRAINT partners_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_resource_action_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_resource_action_key UNIQUE (resource, action);


--
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: program_schedule program_schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.program_schedule
    ADD CONSTRAINT program_schedule_pkey PRIMARY KEY (id);


--
-- Name: radio_recordings radio_recordings_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.radio_recordings
    ADD CONSTRAINT radio_recordings_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_role_id_permission_id_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_permission_id_key UNIQUE (role_id, permission_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: route_funds route_funds_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.route_funds
    ADD CONSTRAINT route_funds_pkey PRIMARY KEY (id);


--
-- Name: route_funds route_funds_route_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.route_funds
    ADD CONSTRAINT route_funds_route_key UNIQUE (route);


--
-- Name: social_embeds social_embeds_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.social_embeds
    ADD CONSTRAINT social_embeds_pkey PRIMARY KEY (id);


--
-- Name: social_links social_links_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.social_links
    ADD CONSTRAINT social_links_pkey PRIMARY KEY (id);


--
-- Name: sponsors sponsors_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.sponsors
    ADD CONSTRAINT sponsors_pkey PRIMARY KEY (id);


--
-- Name: title_section_content title_section_content_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.title_section_content
    ADD CONSTRAINT title_section_content_pkey PRIMARY KEY (id);


--
-- Name: under_construction under_construction_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.under_construction
    ADD CONSTRAINT under_construction_pkey PRIMARY KEY (id);


--
-- Name: uploaded_images uploaded_images_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_pkey PRIMARY KEY (id);


--
-- Name: uploaded_images uploaded_images_public_id_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_public_id_key UNIQUE (public_id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_user_id_role_id_key; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_role_id_key UNIQUE (user_id, role_id);


--
-- Name: verzonden_emails verzonden_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_pkey PRIMARY KEY (id);


--
-- Name: videos videos_pkey; Type: CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY (id);


--
-- Name: idx_aanmelding_antwoorden_aanmelding_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmelding_antwoorden_aanmelding_id ON public.aanmelding_antwoorden USING btree (aanmelding_id);


--
-- Name: INDEX idx_aanmelding_antwoorden_aanmelding_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmelding_antwoorden_aanmelding_id IS 'FK index for registration responses';


--
-- Name: idx_aanmelding_antwoorden_aanmelding_verzonden; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmelding_antwoorden_aanmelding_verzonden ON public.aanmelding_antwoorden USING btree (aanmelding_id, verzond_op DESC);


--
-- Name: INDEX idx_aanmelding_antwoorden_aanmelding_verzonden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmelding_antwoorden_aanmelding_verzonden IS 'Chronological responses per registration';


--
-- Name: idx_aanmelding_antwoorden_verzonden_door; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmelding_antwoorden_verzonden_door ON public.aanmelding_antwoorden USING btree (verzonden_door);


--
-- Name: idx_aanmeldingen_afstand; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_afstand ON public.aanmeldingen USING btree (afstand) WHERE (afstand IS NOT NULL);


--
-- Name: INDEX idx_aanmeldingen_afstand; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmeldingen_afstand IS 'Partial index for distance filtering in reports';


--
-- Name: idx_aanmeldingen_antwoorden_count; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_antwoorden_count ON public.aanmeldingen USING btree (antwoorden_count) WHERE (antwoorden_count > 0);


--
-- Name: idx_aanmeldingen_behandeld; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_behandeld ON public.aanmeldingen USING btree (behandeld_op DESC NULLS LAST, status);


--
-- Name: idx_aanmeldingen_email; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_email ON public.aanmeldingen USING btree (email);


--
-- Name: idx_aanmeldingen_fts; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_fts ON public.aanmeldingen USING gin (to_tsvector('dutch'::regconfig, (((((COALESCE(naam, ''::character varying))::text || ' '::text) || (COALESCE(email, ''::character varying))::text) || ' '::text) || COALESCE(bijzonderheden, ''::text))));


--
-- Name: INDEX idx_aanmeldingen_fts; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmeldingen_fts IS 'Full-text search on registration details';


--
-- Name: idx_aanmeldingen_gebruiker_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_gebruiker_id ON public.aanmeldingen USING btree (gebruiker_id);


--
-- Name: INDEX idx_aanmeldingen_gebruiker_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmeldingen_gebruiker_id IS 'FK index for user registration lookups';


--
-- Name: idx_aanmeldingen_nieuw; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_nieuw ON public.aanmeldingen USING btree (created_at DESC) WHERE ((status)::text = 'nieuw'::text);


--
-- Name: INDEX idx_aanmeldingen_nieuw; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmeldingen_nieuw IS 'Partial index for new registrations';


--
-- Name: idx_aanmeldingen_rol; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_rol ON public.aanmeldingen USING btree (rol);


--
-- Name: idx_aanmeldingen_status; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_status ON public.aanmeldingen USING btree (status);


--
-- Name: idx_aanmeldingen_status_created; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_aanmeldingen_status_created ON public.aanmeldingen USING btree (status, created_at DESC);


--
-- Name: INDEX idx_aanmeldingen_status_created; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_aanmeldingen_status_created IS 'Compound index for registration dashboard';


--
-- Name: idx_chat_channel_participants_channel_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_channel_participants_channel_id ON public.chat_channel_participants USING btree (channel_id);


--
-- Name: idx_chat_channel_participants_user_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_channel_participants_user_id ON public.chat_channel_participants USING btree (user_id);


--
-- Name: idx_chat_channels_public; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_channels_public ON public.chat_channels USING btree (name) WHERE ((is_public = true) AND (is_active = true));


--
-- Name: INDEX idx_chat_channels_public; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_channels_public IS 'Public channel discovery';


--
-- Name: idx_chat_channels_type; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_channels_type ON public.chat_channels USING btree (type);


--
-- Name: INDEX idx_chat_channels_type; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_channels_type IS 'Channel type filtering';


--
-- Name: idx_chat_message_reactions_message_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_message_reactions_message_id ON public.chat_message_reactions USING btree (message_id);


--
-- Name: idx_chat_messages_channel_id_created_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_messages_channel_id_created_at ON public.chat_messages USING btree (channel_id, created_at DESC);


--
-- Name: idx_chat_messages_files; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_messages_files ON public.chat_messages USING btree (channel_id, created_at DESC) WHERE (message_type = ANY (ARRAY['image'::text, 'file'::text]));


--
-- Name: INDEX idx_chat_messages_files; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_messages_files IS 'File and image messages';


--
-- Name: idx_chat_messages_fts; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_messages_fts ON public.chat_messages USING gin (to_tsvector('dutch'::regconfig, COALESCE(content, ''::text)));


--
-- Name: INDEX idx_chat_messages_fts; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_messages_fts IS 'Full-text search on chat message content';


--
-- Name: idx_chat_messages_reply_to; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_messages_reply_to ON public.chat_messages USING btree (reply_to_id, created_at DESC) WHERE (reply_to_id IS NOT NULL);


--
-- Name: INDEX idx_chat_messages_reply_to; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_messages_reply_to IS 'Message reply threads';


--
-- Name: idx_chat_messages_user_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_messages_user_id ON public.chat_messages USING btree (user_id);


--
-- Name: idx_chat_participants_active; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_participants_active ON public.chat_channel_participants USING btree (channel_id, user_id) WHERE (is_active = true);


--
-- Name: INDEX idx_chat_participants_active; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_participants_active IS 'Partial index for active channel participants';


--
-- Name: idx_chat_participants_unread; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_participants_unread ON public.chat_channel_participants USING btree (user_id, last_read_at) WHERE (is_active = true);


--
-- Name: INDEX idx_chat_participants_unread; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_participants_unread IS 'Unread message tracking per user';


--
-- Name: idx_chat_user_presence_online; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_chat_user_presence_online ON public.chat_user_presence USING btree (status, last_seen DESC) WHERE (status <> 'offline'::text);


--
-- Name: INDEX idx_chat_user_presence_online; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_chat_user_presence_online IS 'Online and away users';


--
-- Name: idx_contact_antwoorden_contact_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_antwoorden_contact_id ON public.contact_antwoorden USING btree (contact_id);


--
-- Name: INDEX idx_contact_antwoorden_contact_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_contact_antwoorden_contact_id IS 'FK index for contact responses';


--
-- Name: idx_contact_antwoorden_contact_verzonden; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_antwoorden_contact_verzonden ON public.contact_antwoorden USING btree (contact_id, verzond_op DESC);


--
-- Name: INDEX idx_contact_antwoorden_contact_verzonden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_contact_antwoorden_contact_verzonden IS 'Chronological responses per contact';


--
-- Name: idx_contact_antwoorden_verzonden_door; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_antwoorden_verzonden_door ON public.contact_antwoorden USING btree (verzonden_door);


--
-- Name: idx_contact_formulieren_antwoorden_count; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_antwoorden_count ON public.contact_formulieren USING btree (antwoorden_count) WHERE (antwoorden_count > 0);


--
-- Name: idx_contact_formulieren_behandeld; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_behandeld ON public.contact_formulieren USING btree (behandeld_op DESC NULLS LAST, status);


--
-- Name: idx_contact_formulieren_email; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_email ON public.contact_formulieren USING btree (email);


--
-- Name: idx_contact_formulieren_fts; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_fts ON public.contact_formulieren USING gin (to_tsvector('dutch'::regconfig, (((((COALESCE(naam, ''::character varying))::text || ' '::text) || (COALESCE(email, ''::character varying))::text) || ' '::text) || COALESCE(bericht, ''::text))));


--
-- Name: INDEX idx_contact_formulieren_fts; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_contact_formulieren_fts IS 'Full-text search on name, email, and message';


--
-- Name: idx_contact_formulieren_nieuw; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_nieuw ON public.contact_formulieren USING btree (created_at DESC) WHERE (((status)::text = 'nieuw'::text) AND (beantwoord = false));


--
-- Name: INDEX idx_contact_formulieren_nieuw; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_contact_formulieren_nieuw IS 'Partial index for new unanswered contact forms';


--
-- Name: idx_contact_formulieren_status; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_status ON public.contact_formulieren USING btree (status);


--
-- Name: idx_contact_formulieren_status_created; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_contact_formulieren_status_created ON public.contact_formulieren USING btree (status, created_at DESC) WHERE (beantwoord = false);


--
-- Name: INDEX idx_contact_formulieren_status_created; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_contact_formulieren_status_created IS 'Compound index for unanswered contact forms dashboard';


--
-- Name: idx_dashboard_stats_entity_status; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE UNIQUE INDEX idx_dashboard_stats_entity_status ON public.dashboard_stats USING btree (entity, status);


--
-- Name: idx_gebruikers_email; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_gebruikers_email ON public.gebruikers USING btree (email);


--
-- Name: idx_gebruikers_is_actief; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_gebruikers_is_actief ON public.gebruikers USING btree (is_actief) WHERE (is_actief = true);


--
-- Name: INDEX idx_gebruikers_is_actief; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_gebruikers_is_actief IS 'Partial index for active users only';


--
-- Name: idx_gebruikers_newsletter; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_gebruikers_newsletter ON public.gebruikers USING btree (email) WHERE ((newsletter_subscribed = true) AND (is_actief = true));


--
-- Name: INDEX idx_gebruikers_newsletter; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_gebruikers_newsletter IS 'Partial index for active newsletter subscribers';


--
-- Name: idx_gebruikers_role_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_gebruikers_role_id ON public.gebruikers USING btree (role_id);


--
-- Name: INDEX idx_gebruikers_role_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_gebruikers_role_id IS 'FK index for RBAC role lookups';


--
-- Name: idx_incoming_emails_account_type; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_incoming_emails_account_type ON public.incoming_emails USING btree (account_type);


--
-- Name: idx_incoming_emails_from; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_incoming_emails_from ON public.incoming_emails USING btree ("from");


--
-- Name: INDEX idx_incoming_emails_from; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_incoming_emails_from IS 'Sender lookup for email filtering';


--
-- Name: idx_incoming_emails_is_processed; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_incoming_emails_is_processed ON public.incoming_emails USING btree (is_processed);


--
-- Name: idx_incoming_emails_message_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_incoming_emails_message_id ON public.incoming_emails USING btree (message_id);


--
-- Name: idx_incoming_emails_processing; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_incoming_emails_processing ON public.incoming_emails USING btree (is_processed, received_at DESC) WHERE (is_processed = false);


--
-- Name: INDEX idx_incoming_emails_processing; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_incoming_emails_processing IS 'Partial index for unprocessed email queue';


--
-- Name: idx_migraties_versie; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE UNIQUE INDEX idx_migraties_versie ON public.migraties USING btree (versie);


--
-- Name: idx_newsletters_sent_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_newsletters_sent_at ON public.newsletters USING btree (sent_at);


--
-- Name: idx_newsletters_status; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_newsletters_status ON public.newsletters USING btree (sent_at DESC);


--
-- Name: INDEX idx_newsletters_status; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_newsletters_status IS 'Draft newsletters first, then sent in reverse chronological order';


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at);


--
-- Name: idx_notifications_priority; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_notifications_priority ON public.notifications USING btree (priority);


--
-- Name: idx_notifications_sent; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_notifications_sent ON public.notifications USING btree (sent);


--
-- Name: idx_notifications_type; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);


--
-- Name: idx_permissions_resource_action; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_permissions_resource_action ON public.permissions USING btree (resource, action);


--
-- Name: idx_refresh_tokens_cleanup; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_refresh_tokens_cleanup ON public.refresh_tokens USING btree (expires_at) WHERE (is_revoked = false);


--
-- Name: INDEX idx_refresh_tokens_cleanup; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_refresh_tokens_cleanup IS 'Expired token cleanup for scheduled jobs';


--
-- Name: idx_refresh_tokens_expires_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--
-- Name: idx_refresh_tokens_is_revoked; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_refresh_tokens_is_revoked ON public.refresh_tokens USING btree (is_revoked);


--
-- Name: idx_refresh_tokens_token; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_role_permissions_permission_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_role_permissions_permission_id ON public.role_permissions USING btree (permission_id);


--
-- Name: idx_role_permissions_role_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_role_permissions_role_id ON public.role_permissions USING btree (role_id);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_route_funds_route; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_route_funds_route ON public.route_funds USING btree (route);


--
-- Name: idx_uploaded_images_active; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_active ON public.uploaded_images USING btree (user_id, created_at DESC) WHERE (deleted_at IS NULL);


--
-- Name: INDEX idx_uploaded_images_active; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_uploaded_images_active IS 'Active (non-deleted) images per user';


--
-- Name: idx_uploaded_images_created_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_created_at ON public.uploaded_images USING btree (created_at DESC);


--
-- Name: idx_uploaded_images_deleted_at; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_deleted_at ON public.uploaded_images USING btree (deleted_at);


--
-- Name: idx_uploaded_images_folder; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_folder ON public.uploaded_images USING btree (folder);


--
-- Name: idx_uploaded_images_public_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_public_id ON public.uploaded_images USING btree (public_id);


--
-- Name: idx_uploaded_images_user_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_uploaded_images_user_id ON public.uploaded_images USING btree (user_id);


--
-- Name: idx_user_roles_active; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_user_roles_active ON public.user_roles USING btree (is_active) WHERE (is_active = true);


--
-- Name: idx_user_roles_active_lookup; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_user_roles_active_lookup ON public.user_roles USING btree (user_id, role_id) WHERE (is_active = true);


--
-- Name: INDEX idx_user_roles_active_lookup; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_user_roles_active_lookup IS 'Active user role assignments';


--
-- Name: idx_user_roles_role_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_user_roles_role_id ON public.user_roles USING btree (role_id);


--
-- Name: idx_user_roles_user_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_user_roles_user_id ON public.user_roles USING btree (user_id);


--
-- Name: idx_verzonden_emails_aanmelding_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_aanmelding_id ON public.verzonden_emails USING btree (aanmelding_id);


--
-- Name: INDEX idx_verzonden_emails_aanmelding_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_aanmelding_id IS 'FK index for registration email tracking';


--
-- Name: idx_verzonden_emails_contact_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_contact_id ON public.verzonden_emails USING btree (contact_id);


--
-- Name: INDEX idx_verzonden_emails_contact_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_contact_id IS 'FK index for contact form email tracking';


--
-- Name: idx_verzonden_emails_errors; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_errors ON public.verzonden_emails USING btree (verzonden_op DESC) WHERE ((status)::text = 'failed'::text);


--
-- Name: INDEX idx_verzonden_emails_errors; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_errors IS 'Partial index for failed email tracking';


--
-- Name: idx_verzonden_emails_ontvanger; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_ontvanger ON public.verzonden_emails USING btree (ontvanger);


--
-- Name: idx_verzonden_emails_status; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_status ON public.verzonden_emails USING btree (status);


--
-- Name: INDEX idx_verzonden_emails_status; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_status IS 'Status filtering for error tracking';


--
-- Name: idx_verzonden_emails_status_tijd; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_status_tijd ON public.verzonden_emails USING btree (status, verzonden_op DESC);


--
-- Name: INDEX idx_verzonden_emails_status_tijd; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_status_tijd IS 'Compound index for email status tracking over time';


--
-- Name: idx_verzonden_emails_template_id; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_template_id ON public.verzonden_emails USING btree (template_id);


--
-- Name: INDEX idx_verzonden_emails_template_id; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_template_id IS 'FK index for template usage tracking';


--
-- Name: idx_verzonden_emails_verzonden_op; Type: INDEX; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE INDEX idx_verzonden_emails_verzonden_op ON public.verzonden_emails USING btree (verzonden_op DESC);


--
-- Name: INDEX idx_verzonden_emails_verzonden_op; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON INDEX public.idx_verzonden_emails_verzonden_op IS 'Chronological sorting for email history';


--
-- Name: aanmelding_antwoorden trigger_aanmelding_antwoorden_count; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_aanmelding_antwoorden_count AFTER INSERT OR DELETE ON public.aanmelding_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_aanmelding_antwoorden_count();


--
-- Name: aanmelding_antwoorden trigger_aanmelding_antwoorden_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_aanmelding_antwoorden_updated_at BEFORE UPDATE ON public.aanmelding_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: aanmeldingen trigger_aanmeldingen_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_aanmeldingen_updated_at BEFORE UPDATE ON public.aanmeldingen FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: albums trigger_albums_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_albums_updated_at BEFORE UPDATE ON public.albums FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_channels trigger_chat_channels_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_chat_channels_updated_at BEFORE UPDATE ON public.chat_channels FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_messages trigger_chat_messages_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_chat_messages_updated_at BEFORE UPDATE ON public.chat_messages FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: chat_user_presence trigger_chat_user_presence_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_chat_user_presence_updated_at BEFORE UPDATE ON public.chat_user_presence FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: contact_antwoorden trigger_contact_antwoorden_count; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_contact_antwoorden_count AFTER INSERT OR DELETE ON public.contact_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_contact_antwoorden_count();


--
-- Name: contact_antwoorden trigger_contact_antwoorden_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_contact_antwoorden_updated_at BEFORE UPDATE ON public.contact_antwoorden FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: contact_formulieren trigger_contact_formulieren_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_contact_formulieren_updated_at BEFORE UPDATE ON public.contact_formulieren FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: email_templates trigger_email_templates_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_email_templates_updated_at BEFORE UPDATE ON public.email_templates FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: gebruikers trigger_gebruikers_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_gebruikers_updated_at BEFORE UPDATE ON public.gebruikers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: incoming_emails trigger_incoming_emails_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_incoming_emails_updated_at BEFORE UPDATE ON public.incoming_emails FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: newsletters trigger_newsletters_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_newsletters_updated_at BEFORE UPDATE ON public.newsletters FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: photos trigger_photos_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_photos_updated_at BEFORE UPDATE ON public.photos FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: route_funds trigger_route_funds_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_route_funds_updated_at BEFORE UPDATE ON public.route_funds FOR EACH ROW EXECUTE FUNCTION public.update_route_funds_updated_at();


--
-- Name: sponsors trigger_sponsors_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_sponsors_updated_at BEFORE UPDATE ON public.sponsors FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: uploaded_images trigger_uploaded_images_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_uploaded_images_updated_at BEFORE UPDATE ON public.uploaded_images FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: verzonden_emails trigger_verzonden_emails_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_verzonden_emails_updated_at BEFORE UPDATE ON public.verzonden_emails FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: videos trigger_videos_updated_at; Type: TRIGGER; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

CREATE TRIGGER trigger_videos_updated_at BEFORE UPDATE ON public.videos FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: aanmelding_antwoorden aanmelding_antwoorden_aanmelding_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.aanmelding_antwoorden
    ADD CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey FOREIGN KEY (aanmelding_id) REFERENCES public.aanmeldingen(id) ON DELETE CASCADE;


--
-- Name: CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey ON aanmelding_antwoorden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON CONSTRAINT aanmelding_antwoorden_aanmelding_id_fkey ON public.aanmelding_antwoorden IS 'FK to aanmeldingen';


--
-- Name: aanmeldingen aanmeldingen_gebruiker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.aanmeldingen
    ADD CONSTRAINT aanmeldingen_gebruiker_id_fkey FOREIGN KEY (gebruiker_id) REFERENCES public.gebruikers(id);


--
-- Name: album_photos album_photos_album_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_album_id_fkey FOREIGN KEY (album_id) REFERENCES public.albums(id) ON DELETE CASCADE;


--
-- Name: album_photos album_photos_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.album_photos
    ADD CONSTRAINT album_photos_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES public.photos(id) ON DELETE CASCADE;


--
-- Name: albums albums_cover_photo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_cover_photo_id_fkey FOREIGN KEY (cover_photo_id) REFERENCES public.photos(id);


--
-- Name: chat_channel_participants chat_channel_participants_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_channel_participants
    ADD CONSTRAINT chat_channel_participants_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.chat_channels(id) ON DELETE CASCADE;


--
-- Name: chat_message_reactions chat_message_reactions_message_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_message_reactions
    ADD CONSTRAINT chat_message_reactions_message_id_fkey FOREIGN KEY (message_id) REFERENCES public.chat_messages(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.chat_channels(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_reply_to_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_reply_to_id_fkey FOREIGN KEY (reply_to_id) REFERENCES public.chat_messages(id) ON DELETE SET NULL;


--
-- Name: contact_antwoorden contact_antwoorden_contact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.contact_antwoorden
    ADD CONSTRAINT contact_antwoorden_contact_id_fkey FOREIGN KEY (contact_id) REFERENCES public.contact_formulieren(id) ON DELETE CASCADE;


--
-- Name: CONSTRAINT contact_antwoorden_contact_id_fkey ON contact_antwoorden; Type: COMMENT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

COMMENT ON CONSTRAINT contact_antwoorden_contact_id_fkey ON public.contact_antwoorden IS 'FK to contact_formulieren';


--
-- Name: email_templates email_templates_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.email_templates
    ADD CONSTRAINT email_templates_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.gebruikers(id);


--
-- Name: gebruikers gebruikers_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.gebruikers
    ADD CONSTRAINT gebruikers_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.gebruikers(id);


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: roles roles_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.gebruikers(id);


--
-- Name: uploaded_images uploaded_images_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.uploaded_images
    ADD CONSTRAINT uploaded_images_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.gebruikers(id);


--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.gebruikers(id) ON DELETE CASCADE;


--
-- Name: verzonden_emails verzonden_emails_aanmelding_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_aanmelding_id_fkey FOREIGN KEY (aanmelding_id) REFERENCES public.aanmeldingen(id);


--
-- Name: verzonden_emails verzonden_emails_contact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_contact_id_fkey FOREIGN KEY (contact_id) REFERENCES public.contact_formulieren(id);


--
-- Name: verzonden_emails verzonden_emails_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

ALTER TABLE ONLY public.verzonden_emails
    ADD CONSTRAINT verzonden_emails_template_id_fkey FOREIGN KEY (template_id) REFERENCES public.email_templates(id);


--
-- Name: dashboard_stats; Type: MATERIALIZED VIEW DATA; Schema: public; Owner: dekoninklijkeloopdatabase_user
--

REFRESH MATERIALIZED VIEW public.dashboard_stats;


--
-- PostgreSQL database dump complete
--

\unrestrict ynkt5EaPkd5cffaI1hn4JWbqquN1EXCdrwsFHR8nB3jEQXvwxjdADKAZaZd6Fbx

