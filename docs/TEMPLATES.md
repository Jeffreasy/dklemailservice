# Email Templates Handleiding

## Overzicht

Deze handleiding beschrijft de email templates die gebruikt worden in de DKL Email Service, inclusief:
- Template structuur
- Template data models
- Template rendering
- Template beveiliging
- Best practices

## Template Structuur

### Basis Layout
```html
<!DOCTYPE html>
<html lang="nl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Subject }}</title>
    <style>
        /* Responsive email styling */
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: #004D40;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .content {
            padding: 20px;
            background: #fff;
        }
        .footer {
            text-align: center;
            padding: 20px;
            font-size: 12px;
            color: #666;
        }
        @media only screen and (max-width: 480px) {
            body {
                padding: 10px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>De Koninklijke Loop</h1>
    </div>
    <div class="content">
        {{ template "content" . }}
    </div>
    <div class="footer">
        {{ template "footer" . }}
    </div>
</body>
</html>
```

### Partials
```html
{{/* templates/partials/footer.html */}}
{{ define "footer" }}
<p>
    De Koninklijke Loop<br>
    <a href="https://www.dekoninklijkeloop.nl">www.dekoninklijkeloop.nl</a><br>
    <small>© {{ .Year }} De Koninklijke Loop. Alle rechten voorbehouden.</small>
</p>
{{ end }}
```

## Template Types

### Contact Email Templates

#### 1. Contact Bevestiging (contact_email.html)
```html
{{ define "content" }}
<h2>Bedankt voor uw bericht</h2>
<p>Beste {{ .Naam }},</p>
<p>Bedankt voor uw bericht. We hebben het volgende ontvangen:</p>
<blockquote style="background: #f9f9f9; padding: 15px; border-left: 5px solid #004D40;">
    {{ .Bericht }}
</blockquote>
<p>We zullen zo spoedig mogelijk contact met u opnemen via {{ .Email }}.</p>
<p>Met vriendelijke groet,<br>Team De Koninklijke Loop</p>
{{ end }}
```

#### 2. Contact Admin Notificatie (contact_admin_email.html)
```html
{{ define "content" }}
<h2>Nieuw Contact Formulier</h2>
<table style="width: 100%; border-collapse: collapse;">
    <tr>
        <th style="text-align: left; padding: 8px;">Naam:</th>
        <td style="padding: 8px;">{{ .Naam }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Email:</th>
        <td style="padding: 8px;">{{ .Email }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Bericht:</th>
        <td style="padding: 8px;">{{ .Bericht }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Tijdstip:</th>
        <td style="padding: 8px;">{{ .Timestamp.Format "02-01-2006 15:04" }}</td>
    </tr>
</table>
{{ end }}
```

### Aanmelding Email Templates

#### 1. Aanmelding Bevestiging (aanmelding_email.html)
```html
{{ define "content" }}
<h2>Bedankt voor uw aanmelding</h2>
<p>Beste {{ .Naam }},</p>
<p>Bedankt voor uw aanmelding voor De Koninklijke Loop. Hieronder vindt u een overzicht van uw gegevens:</p>
<table style="width: 100%; border-collapse: collapse;">
    <tr>
        <th style="text-align: left; padding: 8px;">Rol:</th>
        <td style="padding: 8px;">{{ .Rol }}</td>
    </tr>
    {{ if eq .Rol "loper" }}
    <tr>
        <th style="text-align: left; padding: 8px;">Afstand:</th>
        <td style="padding: 8px;">{{ .Afstand }}</td>
    </tr>
    {{ end }}
    <tr>
        <th style="text-align: left; padding: 8px;">Email:</th>
        <td style="padding: 8px;">{{ .Email }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Telefoon:</th>
        <td style="padding: 8px;">{{ .Telefoon }}</td>
    </tr>
    {{ if .Ondersteuning }}
    <tr>
        <th style="text-align: left; padding: 8px;">Ondersteuning:</th>
        <td style="padding: 8px;">{{ .Ondersteuning }}</td>
    </tr>
    {{ end }}
</table>
<p>We houden u via email op de hoogte van belangrijke informatie over het evenement.</p>
{{ end }}
```

#### 2. Aanmelding Admin Notificatie (aanmelding_admin_email.html)
```html
{{ define "content" }}
<h2>Nieuwe Aanmelding</h2>
<table style="width: 100%; border-collapse: collapse;">
    <tr>
        <th style="text-align: left; padding: 8px;">Naam:</th>
        <td style="padding: 8px;">{{ .Naam }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Rol:</th>
        <td style="padding: 8px;">{{ .Rol }}</td>
    </tr>
    {{ if eq .Rol "loper" }}
    <tr>
        <th style="text-align: left; padding: 8px;">Afstand:</th>
        <td style="padding: 8px;">{{ .Afstand }}</td>
    </tr>
    {{ end }}
    <tr>
        <th style="text-align: left; padding: 8px;">Email:</th>
        <td style="padding: 8px;">{{ .Email }}</td>
    </tr>
    <tr>
        <th style="text-align: left; padding: 8px;">Telefoon:</th>
        <td style="padding: 8px;">{{ .Telefoon }}</td>
    </tr>
    {{ if .Ondersteuning }}
    <tr>
        <th style="text-align: left; padding: 8px;">Ondersteuning:</th>
        <td style="padding: 8px;">{{ .Ondersteuning }}</td>
    </tr>
    {{ end }}
    {{ if .Bijzonderheden }}
    <tr>
        <th style="text-align: left; padding: 8px;">Bijzonderheden:</th>
        <td style="padding: 8px;">{{ .Bijzonderheden }}</td>
    </tr>
    {{ end }}
    <tr>
        <th style="text-align: left; padding: 8px;">Tijdstip:</th>
        <td style="padding: 8px;">{{ .Timestamp.Format "02-01-2006 15:04" }}</td>
    </tr>
</table>
{{ end }}
```

## Template Data Models

### Contact Email Data
```go
type ContactEmailData struct {
    Naam      string    `json:"naam"`
    Email     string    `json:"email"`
    Bericht   string    `json:"bericht"`
    Timestamp time.Time `json:"timestamp"`
}
```

### Aanmelding Email Data
```go
type AanmeldingEmailData struct {
    Naam           string    `json:"naam"`
    Email          string    `json:"email"`
    Telefoon       string    `json:"telefoon"`
    Rol            string    `json:"rol"`
    Afstand        string    `json:"afstand,omitempty"`
    Ondersteuning  string    `json:"ondersteuning,omitempty"`
    Bijzonderheden string    `json:"bijzonderheden,omitempty"`
    Timestamp      time.Time `json:"timestamp"`
}
```

## Template Management

### Template Loading
```go
// Template manager
type TemplateManager struct {
    templates  *template.Template
    directory string
    reloadMux sync.RWMutex
}

// Load templates
func (tm *TemplateManager) LoadTemplates() error {
    tm.reloadMux.Lock()
    defer tm.reloadMux.Unlock()

    pattern := filepath.Join(tm.directory, "*.html")
    templates, err := template.New("").Funcs(tm.templateFuncs()).ParseGlob(pattern)
    if err != nil {
        return fmt.Errorf("failed to load templates: %w", err)
    }

    tm.templates = templates
    return nil
}

// Template functions
func (tm *TemplateManager) templateFuncs() template.FuncMap {
    return template.FuncMap{
        "formatDate": func(t time.Time) string {
            return t.Format("02-01-2006")
        },
        "formatTime": func(t time.Time) string {
            return t.Format("15:04")
        },
        "safeHTML": func(s string) template.HTML {
            return template.HTML(bluemonday.UGCPolicy().Sanitize(s))
        },
    }
}
```

### Template Rendering
```go
// Render template
func (tm *TemplateManager) RenderTemplate(name string, data interface{}) (string, error) {
    tm.reloadMux.RLock()
    defer tm.reloadMux.RUnlock()

    var buf bytes.Buffer
    if err := tm.templates.ExecuteTemplate(&buf, name, data); err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", name, err)
    }

    return buf.String(), nil
}
```

## Template Beveiliging

### XSS Preventie
```go
// Template sanitization
func sanitizeTemplateData(data interface{}) interface{} {
    policy := bluemonday.UGCPolicy()
    
    switch v := data.(type) {
    case string:
        return policy.Sanitize(v)
    case map[string]interface{}:
        result := make(map[string]interface{})
        for key, value := range v {
            result[key] = sanitizeTemplateData(value)
        }
        return result
    case []interface{}:
        result := make([]interface{}, len(v))
        for i, value := range v {
            result[i] = sanitizeTemplateData(value)
        }
        return result
    default:
        return v
    }
}
```

### Template Validatie
```go
// Validate template syntax
func validateTemplate(content string) error {
    _, err := template.New("test").Parse(content)
    if err != nil {
        return fmt.Errorf("invalid template syntax: %w", err)
    }
    return nil
}

// Validate template variables
func validateTemplateData(tmpl *template.Template, data interface{}) error {
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return fmt.Errorf("template data validation failed: %w", err)
    }
    return nil
}
```

## Best Practices

### Email Client Compatibiliteit
1. Gebruik tabellen voor layout
2. Inline CSS styles
3. Eenvoudige HTML structuur
4. Fallback fonts
5. Alt tekst voor afbeeldingen

### Performance
1. Minimaliseer template grootte
2. Cache gecompileerde templates
3. Gebruik efficiënte template functies
4. Voorkom complexe template logic
5. Batch template updates

### Onderhoud
1. Gebruik version control voor templates
2. Documenteer template variabelen
3. Test templates regelmatig
4. Monitor render performance
5. Review template gebruik

### Testing
1. Unit tests voor template rendering
2. Validatie van template syntax
3. Preview in verschillende email clients
4. Test met verschillende data sets
5. Verificatie van email layout

## Troubleshooting

### Common Issues

1. Template Not Found
```go
// Check template existence
if _, err := tm.templates.Lookup(name); err != nil {
    log.Printf("Template %s not found in directory %s", name, tm.directory)
}
```

2. Missing Variables
```go
// Validate required variables
func validateRequiredVars(data interface{}, required []string) error {
    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    for _, field := range required {
        if v.FieldByName(field).IsZero() {
            return fmt.Errorf("required variable %s is missing", field)
        }
    }
    return nil
}
```

3. Rendering Errors
```go
// Debug template rendering
func debugTemplateRender(name string, data interface{}) {
    log.Printf("Rendering template: %s", name)
    log.Printf("Template data: %+v", data)
    
    if result, err := tm.RenderTemplate(name, data); err != nil {
        log.Printf("Render error: %v", err)
    } else {
        log.Printf("Render result length: %d", len(result))
    }
}
``` 
   - Inline CSS styles 
```

## Inkomende Email Verwerking

De EmailAutoFetcher haalt automatisch emails op en slaat deze op in de database. Voor het verwerken en weergeven van deze inkomende emails zijn een aantal templates beschikbaar.

### Email Overzicht Template

Deze template wordt gebruikt voor het weergeven van een overzicht van alle opgehaalde emails in het admin dashboard.

**Template Locatie:** `templates/admin/email_list.gohtml`

**Template Functionaliteiten:**
- Lijst weergave van alle opgehaalde emails
- Sortering op datum, afzender, onderwerp
- Filter functies (verwerkt/onverwerkt, account type)
- Paginering

**Voorbeeld:**

```html
{{define "admin/email_list"}}
<!DOCTYPE html>
<html>
<head>
    <title>Email Beheer - De Koninklijke Loop</title>
    {{template "common/head" .}}
</head>
<body>
    {{template "admin/header" .}}
    
    <div class="container mx-auto py-8">
        <h1 class="text-2xl font-bold mb-4">Inkomende Emails</h1>
        
        {{template "admin/email_filters" .}}
        
        <table class="w-full mt-4">
            <thead>
                <tr>
                    <th class="px-4 py-2 text-left">Datum</th>
                    <th class="px-4 py-2 text-left">Van</th>
                    <th class="px-4 py-2 text-left">Onderwerp</th>
                    <th class="px-4 py-2 text-left">Account</th>
                    <th class="px-4 py-2 text-left">Status</th>
                    <th class="px-4 py-2 text-left">Acties</th>
                </tr>
            </thead>
            <tbody>
                {{range .Emails}}
                <tr class="{{if .IsProcessed}}bg-gray-100{{else}}bg-white{{end}}">
                    <td class="px-4 py-2">{{formatDate .ReceivedAt}}</td>
                    <td class="px-4 py-2">{{.From}}</td>
                    <td class="px-4 py-2">{{.Subject}}</td>
                    <td class="px-4 py-2">{{.AccountType}}</td>
                    <td class="px-4 py-2">
                        {{if .IsProcessed}}
                            <span class="px-2 py-1 bg-green-100 text-green-800 rounded">Verwerkt</span>
                        {{else}}
                            <span class="px-2 py-1 bg-yellow-100 text-yellow-800 rounded">Nieuw</span>
                        {{end}}
                    </td>
                    <td class="px-4 py-2">
                        <a href="/admin/emails/{{.ID}}" class="text-blue-500">Bekijken</a>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
        
        {{template "common/pagination" .Pagination}}
    </div>
    
    {{template "common/footer" .}}
</body>
</html>
{{end}}
```

### Email Detail Template

Deze template toont de details van een specifieke email, inclusief de inhoud.

**Template Locatie:** `templates/admin/email_detail.gohtml`

**Template Functionaliteiten:**
- Volledige email details weergave
- Email content weergave (plaintext of HTML)
- Acties (markeren als verwerkt, verwijderen)
- Koppeling naar reply functionaliteit

**Voorbeeld:**

```html
{{define "admin/email_detail"}}
<!DOCTYPE html>
<html>
<head>
    <title>Email Detail - De Koninklijke Loop</title>
    {{template "common/head" .}}
</head>
<body>
    {{template "admin/header" .}}
    
    <div class="container mx-auto py-8">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-bold">Email Detail</h1>
            <div>
                <a href="/admin/emails" class="bg-gray-200 hover:bg-gray-300 px-4 py-2 rounded mr-2">Terug naar overzicht</a>
                {{if not .Email.IsProcessed}}
                    <button id="mark-processed" data-id="{{.Email.ID}}" class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded mr-2">Markeer als verwerkt</button>
                {{end}}
                <button id="delete-email" data-id="{{.Email.ID}}" class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded">Verwijderen</button>
            </div>
        </div>
        
        <div class="bg-white shadow-md rounded p-6 mb-6">
            <div class="mb-4">
                <p class="text-sm text-gray-600">Van</p>
                <p class="font-medium">{{.Email.From}}</p>
            </div>
            <div class="mb-4">
                <p class="text-sm text-gray-600">Aan</p>
                <p class="font-medium">{{.Email.To}}</p>
            </div>
            <div class="mb-4">
                <p class="text-sm text-gray-600">Onderwerp</p>
                <p class="font-medium">{{.Email.Subject}}</p>
            </div>
            <div class="mb-4">
                <p class="text-sm text-gray-600">Ontvangen op</p>
                <p class="font-medium">{{formatDateTime .Email.ReceivedAt}}</p>
            </div>
            <div class="mb-4">
                <p class="text-sm text-gray-600">Account</p>
                <p class="font-medium">{{.Email.AccountType}}</p>
            </div>
            <div class="mb-4">
                <p class="text-sm text-gray-600">Status</p>
                <p class="font-medium">
                    {{if .Email.IsProcessed}}
                        <span class="px-2 py-1 bg-green-100 text-green-800 rounded">Verwerkt op {{formatDateTime .Email.ProcessedAt}}</span>
                    {{else}}
                        <span class="px-2 py-1 bg-yellow-100 text-yellow-800 rounded">Nieuw</span>
                    {{end}}
                </p>
            </div>
        </div>
        
        <div class="bg-white shadow-md rounded p-6">
            <h2 class="text-xl font-bold mb-4">Inhoud</h2>
            {{if eq .Email.ContentType "text/html"}}
                <div class="content-html p-4 bg-gray-50 rounded">
                    {{.EmailHTML}}
                </div>
            {{else}}
                <pre class="whitespace-pre-wrap p-4 bg-gray-50 rounded">{{.Email.Body}}</pre>
            {{end}}
        </div>
    </div>
    
    <script>
        // JavaScript voor mark as processed en delete functies
        document.addEventListener('DOMContentLoaded', function() {
            const markProcessedBtn = document.getElementById('mark-processed');
            if (markProcessedBtn) {
                markProcessedBtn.addEventListener('click', function() {
                    const emailId = this.getAttribute('data-id');
                    fetch(`/api/mail/${emailId}/processed`, {
                        method: 'PUT',
                        headers: {
                            'Authorization': `Bearer ${getAuthToken()}`
                        }
                    })
                    .then(response => response.json())
                    .then(data => {
                        if (data.is_processed) {
                            window.location.reload();
                        }
                    })
                    .catch(error => console.error('Error:', error));
                });
            }
            
            const deleteEmailBtn = document.getElementById('delete-email');
            if (deleteEmailBtn) {
                deleteEmailBtn.addEventListener('click', function() {
                    if (confirm('Weet je zeker dat je deze email wilt verwijderen?')) {
                        const emailId = this.getAttribute('data-id');
                        fetch(`/api/mail/${emailId}`, {
                            method: 'DELETE',
                            headers: {
                                'Authorization': `Bearer ${getAuthToken()}`
                            }
                        })
                        .then(response => response.json())
                        .then(data => {
                            if (data.success) {
                                window.location.href = '/admin/emails';
                            }
                        })
                        .catch(error => console.error('Error:', error));
                    }
                });
            }
            
            function getAuthToken() {
                // Functie om JWT token op te halen uit localStorage of cookies
                return localStorage.getItem('auth_token') || '';
            }
        });
    </script>
    
    {{template "common/footer" .}}
</body>
</html>
{{end}}
```

### Email Sanitization Helpers

Voor het veilig weergeven van inkomende emails worden helper functies gebruikt om de content te sanitiseren en veilig weer te geven.

**Helper Functies:**
- `sanitizeHTML`: Verwijdert potentieel gevaarlijke HTML elementen
- `renderEmailContent`: Geeft email content weer als HTML of plaintext
- `formatEmailAddress`: Formatteert email adressen

**Implementatie:**

```go
// Functie om HTML content te sanitiseren
func sanitizeHTML(htmlContent string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(htmlContent)
}

// Template functie om email content veilig weer te geven
func renderEmailContent(e *models.Email) template.HTML {
    if e.ContentType == "text/html" {
        sanitized := sanitizeHTML(e.Body)
        return template.HTML(sanitized)
    }
    // Voor plaintext, converteer newlines naar <br> tags
    escaped := html.EscapeString(e.Body)
    withBrs := strings.ReplaceAll(escaped, "\n", "<br>")
    return template.HTML(withBrs)
}

// Functie om email adres te formatteren voor privacy
func formatEmailAddress(address string) string {
    parts := strings.Split(address, "@")
    if len(parts) != 2 {
        return address
    }
    username := parts[0]
    domain := parts[1]
    
    // Toon eerste 3 karakters + sterretjes voor de rest
    if len(username) > 3 {
        masked := username[:3] + strings.Repeat("*", len(username)-3)
        return masked + "@" + domain
    }
    return address
}
``` 