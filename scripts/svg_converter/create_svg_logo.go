package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

// SVG kleuren voor De Koninklijke Loop logo
const (
	primaryColor = "#ff9328" // Oranje hoofdkleur
)

// Eenvoudige SVG-representatie van het logo
// Dit is een vereenvoudigde versie die goed werkt in e-mailclients
const svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="77" viewBox="0 0 200 77">
  <g fill="none" fill-rule="evenodd">
    <rect width="200" height="77" fill="%s" rx="8"/>
    <text x="100" y="45" font-family="Arial, sans-serif" font-size="18" font-weight="bold" fill="white" text-anchor="middle">
      De Koninklijke Loop
    </text>
  </g>
</svg>`

// Functie om een SVG-string te genereren
func generateSVG() string {
	return fmt.Sprintf(svgTemplate, primaryColor)
}

// Functie om de SVG naar een Data URI te converteren voor gebruik in HTML
func svgToDataURI(svg string) string {
	// SVG moet geëncodeerd worden voor gebruik in een data URI
	return "data:image/svg+xml;charset=utf-8," + svg
}

// Functie om de img tag in HTML-bestanden te vervangen
func replaceImgTagInFile(filePath string, dataURI string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Verbeterde reguliere expressie om alle varianten van de img tag met base64 data te vinden
	// Deze regex is flexibeler en vindt ook tags met verschillende attribuutvolgorde
	re := regexp.MustCompile(`<img[^>]*src="data:image/png;base64,[^"]*"[^>]*>`)

	// Nieuwe img tag met SVG data URI
	newImgTag := fmt.Sprintf(`<img src="%s" alt="De Koninklijke Loop" class="logo" style="max-width: 200px; width: 100%%; height: auto;">`, dataURI)

	// Vervang de oude img tag met de nieuwe
	newContent := re.ReplaceAllString(string(content), newImgTag)

	// Schrijf de nieuwe inhoud terug naar het bestand
	return ioutil.WriteFile(filePath, []byte(newContent), 0644)
}

func main() {
	// Genereer de SVG
	svg := generateSVG()

	// Converteer naar data URI
	dataURI := svgToDataURI(svg)

	fmt.Printf("SVG gegenereerd met grootte: %d bytes\n", len(svg))
	fmt.Printf("Data URI grootte: %d bytes\n", len(dataURI))

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
		err := replaceImgTagInFile(template, dataURI)
		if err != nil {
			log.Fatalf("Fout bij het verwerken van %s: %v", template, err)
		}
		fmt.Printf("Logo succesvol vervangen in %s\n", template)
	}

	fmt.Println("Alle sjablonen zijn bijgewerkt met het SVG-logo!")
}
