package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// ICO file format structures
type iconDir struct {
	Reserved  uint16
	Type      uint16
	Count     uint16
	Directory []iconDirEntry
}

type iconDirEntry struct {
	Width       byte
	Height      byte
	ColorCount  byte
	Reserved    byte
	Planes      uint16
	BitCount    uint16
	BytesInRes  uint32
	ImageOffset uint32
}

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

	// Converteer naar PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, resized); err != nil {
		log.Fatal("Kon resized image niet encoderen:", err)
	}
	pngData := buf.Bytes()

	// Maak ICO header
	dir := iconDir{
		Reserved: 0,
		Type:     1,
		Count:    1,
		Directory: []iconDirEntry{{
			Width:       32,
			Height:      32,
			ColorCount:  0,
			Reserved:    0,
			Planes:      1,
			BitCount:    32,
			BytesInRes:  uint32(len(pngData)),
			ImageOffset: 22, // 6 + 16 (header size + directory size)
		}},
	}

	// Maak public directory als die nog niet bestaat
	err = os.MkdirAll("public", 0755)
	if err != nil {
		log.Fatal("Kon public directory niet maken:", err)
	}

	// Schrijf ICO file
	out, err := os.Create(filepath.Join("public", "favicon.ico"))
	if err != nil {
		log.Fatal("Kon favicon.ico niet maken:", err)
	}
	defer out.Close()

	// Schrijf header
	binary.Write(out, binary.LittleEndian, dir.Reserved)
	binary.Write(out, binary.LittleEndian, dir.Type)
	binary.Write(out, binary.LittleEndian, dir.Count)

	// Schrijf directory entry
	entry := dir.Directory[0]
	binary.Write(out, binary.LittleEndian, entry.Width)
	binary.Write(out, binary.LittleEndian, entry.Height)
	binary.Write(out, binary.LittleEndian, entry.ColorCount)
	binary.Write(out, binary.LittleEndian, entry.Reserved)
	binary.Write(out, binary.LittleEndian, entry.Planes)
	binary.Write(out, binary.LittleEndian, entry.BitCount)
	binary.Write(out, binary.LittleEndian, entry.BytesInRes)
	binary.Write(out, binary.LittleEndian, entry.ImageOffset)

	// Schrijf PNG data
	out.Write(pngData)

	log.Println("Favicon succesvol gedownload en opgeslagen als ICO")
}
