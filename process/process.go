package process

import (
	"fmt"
	"log"
	"path/filepath"
	"srt-translator/config"
	"srt-translator/openai"
)

func ProcessFile(inputFile, targetLang string) {
	fmt.Printf("Procesando archivo: %s\n", inputFile)
	outputFile := filepath.Join(config.OutputDir, "translated_"+filepath.Base(inputFile))

	fmt.Printf("Leyendo archivo SRT: %s\n", inputFile)
	content, err := readSRTFile(filepath.Join(config.InputDir, inputFile))
	if err != nil {
		log.Printf("Error al leer el archivo SRT %s: %v\n", inputFile, err)
		return
	}

	fmt.Printf("Traduciendo contenido del archivo: %s\n", inputFile)
	translatedContent, err := openai.TranslateSRT(content, targetLang)
	if err != nil {
		log.Printf("Error al traducir el texto para el archivo %s: %v\n", inputFile, err)
		return
	}

	fmt.Printf("Escribiendo archivo SRT traducido: %s\n", outputFile)
	err = writeSRTFile(outputFile, translatedContent)
	if err != nil {
		log.Printf("Error al escribir el archivo SRT traducido %s: %v\n", outputFile, err)
		return
	}

	fmt.Printf("Archivo procesado exitosamente: %s\n", outputFile)

}
