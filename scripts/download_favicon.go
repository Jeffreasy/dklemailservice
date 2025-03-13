package main

import (
	"bytes"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func main() {
	// URL van de Cloudinary afbeelding
	imageURL := "https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png"

	// Download de afbeelding
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Fatal("Kon afbeelding niet downloaden:", err)
	}
	defer resp.Body.Close()

	// Lees de afbeelding
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Kon afbeelding data niet lezen:", err)
	}

	// Decode de PNG
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal("Kon PNG niet decoderen:", err)
	}

	// Resize naar favicon formaat (32x32)
	resized := resize.Resize(32, 32, img, resize.Lanczos3)

	// Maak public directory als die nog niet bestaat
	err = os.MkdirAll("public", 0755)
	if err != nil {
		log.Fatal("Kon public directory niet maken:", err)
	}

	// Sla op als favicon.ico
	out, err := os.Create(filepath.Join("public", "favicon.ico"))
	if err != nil {
		log.Fatal("Kon favicon.ico niet maken:", err)
	}
	defer out.Close()

	// Encode als PNG
	err = png.Encode(out, resized)
	if err != nil {
		log.Fatal("Kon favicon niet opslaan:", err)
	}

	log.Println("Favicon succesvol gedownload en opgeslagen")
}
