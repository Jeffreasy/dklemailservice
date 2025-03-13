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