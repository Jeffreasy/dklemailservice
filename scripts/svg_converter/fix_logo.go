package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

func main() {
	// Pad naar het originele logo
	logoPath := "../../scripts/logo_converter/dkllogo.png"

	// Lees het logo bestand
	logoData, err := ioutil.ReadFile(logoPath)
	if err != nil {
		log.Fatalf("Fout bij het lezen van het logo bestand: %v", err)
	}

	// Converteer naar base64
	logoBase64 := base64.StdEncoding.EncodeToString(logoData)

	// Maak de volledige img tag
	imgTag := fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="De Koninklijke Loop" class="logo" style="max-width: 200px; width: 100%%; height: auto;">`, logoBase64)

	fmt.Printf("Logo geladen, grootte: %d bytes\n", len(logoData))
	fmt.Printf("Base64 grootte: %d bytes\n", len(logoBase64))

	// Pad naar de templates directory
	templatesDir := "../../templates"

	// Lijst van alle e-mailsjablonen
	templates := []string{
		filepath.Join(templatesDir, "aanmelding_email.html"),
		filepath.Join(templatesDir, "aanmelding_admin_email.html"),
		filepath.Join(templatesDir, "contact_admin_email.html"),
		filepath.Join(templatesDir, "contact_email.html"),
	}

	// Vervang de img tag in elk sjabloon
	for _, template := range templates {
		// Lees het sjabloon bestand
		content, err := ioutil.ReadFile(template)
		if err != nil {
			log.Fatalf("Fout bij het lezen van %s: %v", template, err)
		}

		// Reguliere expressie om de img tag te vinden
		re := regexp.MustCompile(`<img[^>]*src="data:image/[^"]*"[^>]*>`)

		// Vervang de oude img tag met de nieuwe
		newContent := re.ReplaceAllString(string(content), imgTag)

		// Schrijf de nieuwe inhoud terug naar het bestand
		err = ioutil.WriteFile(template, []byte(newContent), 0644)
		if err != nil {
			log.Fatalf("Fout bij het schrijven naar %s: %v", template, err)
		}

		fmt.Printf("Logo succesvol vervangen in %s\n", template)
	}

	fmt.Println("Alle sjablonen zijn bijgewerkt met het originele logo!")
}
