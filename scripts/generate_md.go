package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	output := "# Database Folder Content\n\n"

	err := filepath.WalkDir("database", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Skip the script itself if it were in database, but it's not
		if strings.HasSuffix(path, "generate_md.go") || strings.HasSuffix(path, "database_full_content.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		output += fmt.Sprintf("## %s\n\n```\n%s\n```\n\n", path, string(content))
		return nil
	})

	if err != nil {
		panic(err)
	}

	err = os.WriteFile("database/database_full_content.md", []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated database/database_full_content.md")
}
