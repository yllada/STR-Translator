package utils

import (
	"fmt"
	"os"
	"strings"
)

func ListFilesInDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".srt") {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func GetValidLangCodes() map[string]string {
	// Lista de códigos de idioma válidos
	return map[string]string{
		"es": "Español",
		"fr": "Francés",
		"de": "Alemán",
		"it": "Italiano",
		"pt": "Portugués",
		"en": "Inglés",
		// Agregar más códigos de idioma según sea necesario
	}
}
