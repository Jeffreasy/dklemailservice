package services

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"dklautomationgo/logger"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// EmailDecoder handles advanced email decoding including multipart, charset conversion, and encoding
type EmailDecoder struct{}

// NewEmailDecoder creates a new email decoder
func NewEmailDecoder() *EmailDecoder {
	return &EmailDecoder{}
}

// DecodeEmailBody decodes an email body with proper MIME parsing, charset conversion, and transfer decoding
func (d *EmailDecoder) DecodeEmailBody(m *mail.Message) (string, error) {
	contentType := m.Header.Get("Content-Type")

	// Parse content type to get media type and parameters
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		logger.Warn("Could not parse content type, treating as plain text", "content_type", contentType, "error", err)
		return d.readSimpleBody(m)
	}

	// Handle multipart messages
	if strings.HasPrefix(mediaType, "multipart/") {
		return d.decodeMultipartBody(m, params)
	}

	// Handle single part messages
	return d.decodeSinglePart(m.Body, m.Header.Get("Content-Transfer-Encoding"), params)
}

// decodeMultipartBody handles multipart MIME messages
func (d *EmailDecoder) decodeMultipartBody(m *mail.Message, params map[string]string) (string, error) {
	boundary := params["boundary"]
	if boundary == "" {
		logger.Warn("Multipart message without boundary")
		return d.readSimpleBody(m)
	}

	mr := multipart.NewReader(m.Body, boundary)
	var htmlPart string
	var textPart string

	// Read all parts
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Warn("Error reading multipart section", "error", err)
			continue
		}

		// Get part content type
		partContentType := part.Header.Get("Content-Type")
		partMediaType, partParams, err := mime.ParseMediaType(partContentType)
		if err != nil {
			logger.Warn("Could not parse part content type", "error", err)
			continue
		}

		// Get transfer encoding for this part
		partEncoding := part.Header.Get("Content-Transfer-Encoding")

		// Read and decode part body
		partBody, err := d.decodePartBody(part, partEncoding, partParams)
		if err != nil {
			logger.Warn("Error decoding part body", "error", err, "content_type", partMediaType)
			continue
		}

		// Store HTML and text parts
		if strings.HasPrefix(partMediaType, "text/html") {
			htmlPart = partBody
		} else if strings.HasPrefix(partMediaType, "text/plain") {
			textPart = partBody
		}
	}

	// Prefer HTML over plain text
	if htmlPart != "" {
		return htmlPart, nil
	}
	if textPart != "" {
		// Convert plain text to HTML
		return d.textToHTML(textPart), nil
	}

	return "", nil
}

// decodeSinglePart decodes a single part message
func (d *EmailDecoder) decodeSinglePart(body io.Reader, encoding string, params map[string]string) (string, error) {
	// Apply transfer encoding
	decodedBody, err := d.decodeTransferEncoding(body, encoding)
	if err != nil {
		return "", err
	}

	// Apply charset conversion
	charset := params["charset"]
	if charset != "" && !strings.EqualFold(charset, "utf-8") && !strings.EqualFold(charset, "us-ascii") {
		converted, err := d.convertCharset(decodedBody, charset)
		if err != nil {
			logger.Warn("Charset conversion failed, using original", "charset", charset, "error", err)
			return decodedBody, nil
		}
		return converted, nil
	}

	return decodedBody, nil
}

// decodePartBody decodes a multipart section
func (d *EmailDecoder) decodePartBody(part *multipart.Part, encoding string, params map[string]string) (string, error) {
	// First decode the transfer encoding
	decoded, err := d.decodeTransferEncoding(part, encoding)
	if err != nil {
		return "", err
	}

	// Then convert charset if needed
	charset := params["charset"]
	if charset != "" && !strings.EqualFold(charset, "utf-8") && !strings.EqualFold(charset, "us-ascii") {
		converted, err := d.convertCharset(decoded, charset)
		if err != nil {
			logger.Warn("Charset conversion failed for part", "charset", charset, "error", err)
			return decoded, nil
		}
		return converted, nil
	}

	return decoded, nil
}

// decodeTransferEncoding decodes the transfer encoding (quoted-printable, base64, etc.)
func (d *EmailDecoder) decodeTransferEncoding(reader io.Reader, encoding string) (string, error) {
	var decodedReader io.Reader

	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "quoted-printable":
		decodedReader = quotedprintable.NewReader(reader)
		logger.Debug("Applying quoted-printable decoding")
	case "base64":
		decodedReader = base64.NewDecoder(base64.StdEncoding, reader)
		logger.Debug("Applying base64 decoding")
	case "7bit", "8bit", "binary", "":
		decodedReader = reader
	default:
		logger.Warn("Unknown transfer encoding, treating as plain", "encoding", encoding)
		decodedReader = reader
	}

	bodyBytes, err := io.ReadAll(decodedReader)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// convertCharset converts text from one charset to UTF-8
func (d *EmailDecoder) convertCharset(text string, charset string) (string, error) {
	charset = strings.ToLower(strings.TrimSpace(charset))

	var decoder *charmap.Charmap
	switch charset {
	case "windows-1252", "cp1252":
		decoder = charmap.Windows1252
	case "iso-8859-1", "latin1":
		decoder = charmap.ISO8859_1
	case "iso-8859-15", "latin9":
		decoder = charmap.ISO8859_15
	case "windows-1250":
		decoder = charmap.Windows1250
	case "windows-1251":
		decoder = charmap.Windows1251
	default:
		logger.Debug("Unsupported charset, skipping conversion", "charset", charset)
		return text, nil
	}

	// Convert to UTF-8
	reader := transform.NewReader(strings.NewReader(text), decoder.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// readSimpleBody reads body without MIME parsing (fallback)
func (d *EmailDecoder) readSimpleBody(m *mail.Message) (string, error) {
	bodyBytes, err := io.ReadAll(m.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

// textToHTML converts plain text to simple HTML
func (d *EmailDecoder) textToHTML(text string) string {
	var buf bytes.Buffer

	buf.WriteString("<div style='white-space: pre-wrap; font-family: monospace;'>")

	// Escape HTML characters
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")

	// Convert line breaks to <br>
	text = strings.ReplaceAll(text, "\r\n", "<br>")
	text = strings.ReplaceAll(text, "\n", "<br>")

	buf.WriteString(text)
	buf.WriteString("</div>")

	return buf.String()
}

// DecodeSubject decodes email subject with RFC 2047 encoding
func (d *EmailDecoder) DecodeSubject(subject string) string {
	dec := new(mime.WordDecoder)
	decoded, err := dec.DecodeHeader(subject)
	if err != nil {
		logger.Warn("Could not decode subject", "subject", subject, "error", err)
		return subject
	}
	return decoded
}

// DecodeFrom decodes and parses the From header
func (d *EmailDecoder) DecodeFrom(from string) string {
	dec := new(mime.WordDecoder)
	decoded, err := dec.DecodeHeader(from)
	if err != nil {
		logger.Warn("Could not decode from header", "from", from, "error", err)
		return from
	}
	return decoded
}
