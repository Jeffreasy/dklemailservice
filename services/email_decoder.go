package services

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"path/filepath"
	"strings"

	"dklautomationgo/logger" // Je eigen logger

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// Attachment bevat de data voor een bijlage of inline afbeelding
type Attachment struct {
	Filename    string
	ContentType string
	ContentID   string // Voor inline afbeeldingen (cid:)
	Data        []byte
}

// DecodedEmail bevat alle uitgepakte onderdelen van een e-mail
type DecodedEmail struct {
	Subject     string
	From        *mail.Address
	To          []*mail.Address
	Cc          []*mail.Address
	HTMLBody    string
	TextBody    string
	Attachments []Attachment
	Inline      []Attachment
}

// EmailDecoder verwerkt het decoderen van complexe e-mails
type EmailDecoder struct {
	WordDecoder *mime.WordDecoder
}

// NewEmailDecoder maakt een nieuwe decoder
func NewEmailDecoder() *EmailDecoder {
	return &EmailDecoder{
		WordDecoder: new(mime.WordDecoder),
	}
}

// decodeHeader decodeert een enkele header string met RFC 2047
func (d *EmailDecoder) decodeHeader(s string) string {
	decoded, err := d.WordDecoder.DecodeHeader(s)
	if err != nil {
		logger.Warn("Could not decode header, using original", "header", s, "error", err)
		return s // Retourneer origineel bij fout
	}
	return decoded
}

// DecodeEmail is de hoofd-instapfunctie.
// Het parseert de mail.Message en retourneert de DecodedEmail struct.
func (d *EmailDecoder) DecodeEmail(m *mail.Message) (*DecodedEmail, error) {
	result := &DecodedEmail{}

	// --- Header Decoding ---
	// Decodeer Subject (simpele header)
	result.Subject = d.decodeHeader(m.Header.Get("Subject"))

	// Decodeer Adres-headers (complex, lijsten)
	// AddressList handelt RFC 2047 decodering voor namen automatisch af.
	if from, err := m.Header.AddressList("From"); err == nil {
		if len(from) > 0 {
			result.From = from[0]
		}
	} else {
		logger.Warn("Could not parse 'From' header", "value", m.Header.Get("From"), "error", err)
	}

	if to, err := m.Header.AddressList("To"); err == nil {
		result.To = to
	} else {
		logger.Warn("Could not parse 'To' header", "value", m.Header.Get("To"), "error", err)
	}

	if cc, err := m.Header.AddressList("Cc"); err == nil {
		result.Cc = cc
	} else {
		logger.Warn("Could not parse 'Cc' header", "value", m.Header.Get("Cc"), "error", err)
	}
	// --- End Header Decoding ---

	// Start het recursief parsen van de body
	err := d.parsePart(m.Body, m.Header, result)
	if err != nil {
		return nil, err
	}

	// --- Post-Processing ---

	// Voeg inline afbeeldingen in de HTML in als data-URIs
	d.embedInlineImages(result)

	// Fallback: Als we Text hebben maar geen HTML, converteer Text naar simpele HTML
	if result.TextBody != "" && result.HTMLBody == "" {
		result.HTMLBody = d.textToHTML(result.TextBody)
	}

	return result, nil
}

// parsePart is de kern: een recursieve functie die elk MIME-deel verwerkt
func (d *EmailDecoder) parsePart(body io.Reader, header mail.Header, result *DecodedEmail) error {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain" // Veel e-mails zonder header zijn plain text
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		logger.Warn("Could not parse content type, skipping part", "content_type", contentType, "error", err)
		return nil // Ga door met andere parts
	}

	// Is dit een multipart bericht?
	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			logger.Warn("Multipart message without boundary, reading as plain", "media_type", mediaType)
			// Probeer de body te lezen als een simpele tekst, ook al is het 'multipart'
			return d.processNonMultipart(body, header, result)
		}

		mr := multipart.NewReader(body, boundary)
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Warn("Error reading multipart section, skipping part", "error", err)
				continue
			}

			// RECURSIEVE AANROEP: verwerk het onderdeel
			mailHeader := mail.Header(part.Header)
			if err := d.parsePart(part, mailHeader, result); err != nil {
				logger.Warn("Error parsing nested part, skipping", "error", err)
			}
		}
		return nil
	}

	// Geen multipart, dus verwerk het als een 'simpel' onderdeel (tekst, bijlage, etc.)
	return d.processNonMultipart(body, header, result)
}

// processNonMultipart verwerkt een 'blad' onderdeel (geen multipart)
func (d *EmailDecoder) processNonMultipart(body io.Reader, header mail.Header, result *DecodedEmail) error {
	// Check voor bijlagen of inline afbeeldingen
	contentDisposition := header.Get("Content-Disposition")
	disposition, dParams, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		disposition = "inline" // Standaard
	}

	// Decodeer de body (transfer encoding + charset)
	partData, err := d.decodePartBody(body, header)
	if err != nil {
		logger.Warn("Failed to decode part body", "error", err)
		return nil // Sla dit onderdeel over
	}

	contentType := header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contentType) // We weten dat dit werkt van parsePart

	contentID := header.Get("Content-ID")
	contentID = strings.Trim(contentID, "<>") // Verwijder '<' en '>'

	// Is het een bijlage?
	if disposition == "attachment" {
		filename := dParams["filename"]
		if filename == "" {
			filename = "attachment" // Fallback
		}
		filename = d.decodeHeader(filename) // Decodeer bestandsnaam (RFC 2047)

		result.Attachments = append(result.Attachments, Attachment{
			Filename:    filepath.Base(filename), // Clean up paden
			ContentType: mediaType,
			Data:        partData,
		})
		return nil
	}

	// Is het een inline afbeelding?
	if contentID != "" && strings.HasPrefix(mediaType, "image/") {
		filename := dParams["filename"] // Kan ook een bestandsnaam hebben
		filename = d.decodeHeader(filename)

		result.Inline = append(result.Inline, Attachment{
			Filename:    filepath.Base(filename),
			ContentType: mediaType,
			ContentID:   contentID,
			Data:        partData,
		})
		return nil
	}

	// Is het de HTML body?
	if strings.HasPrefix(mediaType, "text/html") {
		// We slaan de 'beste' (meestal laatste in alternative) op
		result.HTMLBody = string(partData)
		return nil
	}

	// Is het de platte tekst body?
	if strings.HasPrefix(mediaType, "text/plain") {
		result.TextBody = string(partData)
		return nil
	}

	// Anders... negeer het (bijv. vCard, kalender-items die niet als bijlage zijn gemarkeerd)
	logger.Debug("Skipping unhandled inline part", "media_type", mediaType)
	return nil
}

// decodePartBody past transfer encoding en charset conversie toe
func (d *EmailDecoder) decodePartBody(reader io.Reader, header mail.Header) ([]byte, error) {
	// Stap 1: Decodeer transfer encoding (base64, quoted-printable)
	transferEncoding := header.Get("Content-Transfer-Encoding")
	decodedReader, err := d.decodeTransferEncoding(reader, transferEncoding)
	if err != nil {
		return nil, err
	}

	// Stap 2: Decodeer charset (naar UTF-8)
	contentType := header.Get("Content-Type")
	_, params, _ := mime.ParseMediaType(contentType)
	charset := "us-ascii" // Default
	if params["charset"] != "" {
		charset = params["charset"]
	}

	utf8Reader, err := d.convertCharset(decodedReader, charset)
	if err != nil {
		logger.Warn("Charset conversion failed, reading as raw", "charset", charset, "error", err)
		// Fallback: lees de 'raw' (alleen transfer-decoded) body
		bodyBytes, readErr := io.ReadAll(decodedReader) // Moet originele decodedReader lezen
		if readErr != nil {
			return nil, readErr // Fout bij het lezen van fallback
		}
		return bodyBytes, nil // Retourneer raw data, niet de error van convertCharset
	}

	return io.ReadAll(utf8Reader)
}

// decodeTransferEncoding retourneert een io.Reader die de data on-the-fly decodeert
func (d *EmailDecoder) decodeTransferEncoding(reader io.Reader, encoding string) (io.Reader, error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "quoted-printable":
		logger.Debug("Applying quoted-printable decoding")
		return quotedprintable.NewReader(reader), nil
	case "base64":
		logger.Debug("Applying base64 decoding")
		return base64.NewDecoder(base64.StdEncoding, reader), nil
	case "7bit", "8bit", "binary", "":
		return reader, nil // Geen decodering nodig
	default:
		logger.Warn("Unknown transfer encoding, treating as plain", "encoding", encoding)
		return reader, nil
	}
}

// convertCharset (ROBUUST) converteert een reader van een specifieke charset naar UTF-8
func (d *EmailDecoder) convertCharset(reader io.Reader, charsetName string) (io.Reader, error) {
	charsetName = strings.ToLower(strings.TrimSpace(charsetName))
	if charsetName == "" || strings.EqualFold(charsetName, "utf-8") || strings.EqualFold(charsetName, "us-ascii") {
		return reader, nil // Geen conversie nodig
	}

	// Gebruik de robuuste charset library
	e, name := charset.Lookup(charsetName)
	if e == nil {
		// Probeer IANA-naam (soms weet de ene het wel en de andere niet)
		// Skip IANA lookup for now to avoid import issues
		logger.Warn("Unsupported charset, skipping conversion", "charset", charsetName)
		return reader, nil // Retourneer de originele reader
	} else {
		logger.Debug("Charset matched via x/net/html/charset", "original", charsetName, "matched", name)
	}

	// Converteer naar UTF-8
	return transform.NewReader(reader, e.NewDecoder()), nil
}

// embedInlineImages (NIEUWE FUNCTIE) herschrijft cid: links in HTML naar data: URIs
func (d *EmailDecoder) embedInlineImages(email *DecodedEmail) {
	if email.HTMLBody == "" || len(email.Inline) == 0 {
		return // Geen HTML of geen inline afbeeldingen om te verwerken
	}

	html := email.HTMLBody
	for _, inline := range email.Inline {
		if inline.ContentID == "" {
			continue
		}

		// Maak data URI
		b64data := base64.StdEncoding.EncodeToString(inline.Data)
		dataURI := "data:" + inline.ContentType + ";base64," + b64data

		// Maak cid string
		cid := "cid:" + inline.ContentID

		// Vervang in HTML. We vervangen zowel met dubbele als enkele quotes
		// voor robuustheid tegen slordige HTML.
		html = strings.ReplaceAll(html, `src="`+cid+`"`, `src="`+dataURI+`"`)
		html = strings.ReplaceAll(html, `src='`+cid+`'`, `src='`+dataURI+`'`)
	}
	email.HTMLBody = html
}

// textToHTML (GECORRIGEERD) converteert platte tekst naar simpele HTML
func (d *EmailDecoder) textToHTML(text string) string {
	var buf bytes.Buffer
	buf.WriteString("<div style='white-space: pre-wrap; font-family: monospace;'>")

	// Escape HTML characters
	text = strings.ReplaceAll(text, "&", "&amp;")   // GECORRIGEERD
	text = strings.ReplaceAll(text, "<", "&lt;")    // GECORRIGEERD
	text = strings.ReplaceAll(text, ">", "&gt;")    // GECORRIGEERD
	text = strings.ReplaceAll(text, "\"", "&quot;") // GECORRIGEERD

	// Convert line breaks to <br>
	text = strings.ReplaceAll(text, "\r\n", "<br>")
	text = strings.ReplaceAll(text, "\n", "<br>")

	buf.WriteString(text)
	buf.WriteString("</div>")
	return buf.String()
}

// DecodeSubject decodeert een subject header met RFC 2047 encoding
func (d *EmailDecoder) DecodeSubject(subject string) string {
	return d.decodeHeader(subject)
}

// DecodeFrom decodeert een from header naar een leesbare string
func (d *EmailDecoder) DecodeFrom(from string) string {
	// Parse de from header als address list
	if addresses, err := mail.ParseAddressList(from); err == nil && len(addresses) > 0 {
		// Gebruik de eerste address en formatteer als "Name <email>" of alleen email
		addr := addresses[0]
		if addr.Name != "" {
			return addr.Name + " <" + addr.Address + ">"
		}
		return addr.Address
	}

	// Fallback: probeer als enkele address
	if addr, err := mail.ParseAddress(from); err == nil {
		if addr.Name != "" {
			return addr.Name + " <" + addr.Address + ">"
		}
		return addr.Address
	}

	// Laatste fallback: retourneer origineel
	logger.Warn("Could not parse 'From' header, using original", "from", from)
	return from
}

// DecodeEmailBody is de backward-compatible methode die alleen de body string retourneert
func (d *EmailDecoder) DecodeEmailBody(m *mail.Message) (string, error) {
	// Gebruik de volledige, robuuste decoder
	decoded, err := d.DecodeEmail(m)
	if err != nil {
		return "", err
	}

	// `DecodeEmail` heeft de HTMLBody al voorbereid,
	// inclusief de text-to-HTML fallback.
	return decoded.HTMLBody, nil
}
