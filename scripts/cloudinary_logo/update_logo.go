package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

func main() {
	// Cloudinary URL voor het logo
	logoURL := "https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png"

	// Maak de volledige img tag
	imgTag := fmt.Sprintf(`<img src="%s" alt="De Koninklijke Loop" class="logo" style="max-width: 200px; width: 100%%; height: auto;">`, logoURL)

	fmt.Printf("Cloudinary logo URL: %s\n", logoURL)

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
		re := regexp.MustCompile(`<img[^>]*src="[^"]*"[^>]*>`)

		// Vervang de oude img tag met de nieuwe
		newContent := re.ReplaceAllString(string(content), imgTag)

		// Schrijf de nieuwe inhoud terug naar het bestand
		err = ioutil.WriteFile(template, []byte(newContent), 0644)
		if err != nil {
			log.Fatalf("Fout bij het schrijven naar %s: %v", template, err)
		}

		fmt.Printf("Logo succesvol vervangen in %s\n", template)
	}

	fmt.Println("Alle sjablonen zijn bijgewerkt met het Cloudinary logo!")
}
