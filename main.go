package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"srt-translator/config"
	"srt-translator/process"
	"srt-translator/utils"
	"sync"

	"github.com/sqweek/dialog"
)

func main() {
	// Solicitar al usuario que seleccione la ruta del directorio de entrada
	inputDir, err := dialog.Directory().Title("Seleccione el directorio de entrada").Browse()
	if err != nil {
		log.Fatalf(string(config.Bold)+string(config.Red)+"Error al seleccionar el directorio de entrada: %v"+string(config.Reset), err)
	}

	// Verificar si la ruta del directorio de entrada no está vacía y si el directorio existe
	if inputDir == "" {
		log.Fatalf(string(config.Bold) + string(config.Yellow) + "La ruta del directorio de entrada no puede estar vacía" + string(config.Reset))
	}
	inputDir, err = filepath.Abs(inputDir)
	if err != nil {
		log.Fatalf(string(config.Bold)+string(config.Red)+"Error al obtener la ruta absoluta del directorio de entrada: %v"+string(config.Reset), err)
	}
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		log.Fatalf(string(config.Bold)+string(config.Red)+"El directorio de entrada no existe: %s"+string(config.Reset), inputDir)
	}
	config.InputDir = inputDir
	fmt.Printf(string(config.Blue)+"Directorio de entrada: %s\n"+string(config.Reset), config.InputDir)

	// Solicitar al usuario que seleccione la ruta del directorio de salida
	outputDir, err := dialog.Directory().Title("Seleccione el directorio de salida").Browse()
	if err != nil {
		log.Fatalf("Error al seleccionar el directorio de salida: %v", err)
	}

	// Verificar si la ruta del directorio de salida no está vacía y si el directorio existe
	if outputDir == "" {
		log.Fatalf(string(config.Bold) + string(config.Red) + "La ruta del directorio de salida no puede estar vacía")
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("Error al obtener la ruta absoluta del directorio de salida: %v", err)
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Fatalf("El directorio de salida no existe: %s", outputDir)
	}
	config.OutputDir = outputDir
	fmt.Printf(string(config.Blue)+"Directorio de salida: %s\n"+string(config.Reset), config.OutputDir)

	// Mostrar la lista de idiomas permitidos
	fmt.Println(string(config.Bold) + string(config.Yellow) + "Idiomas disponibles:")
	langCodes := utils.GetValidLangCodes()
	for code, lang := range langCodes {
		fmt.Printf(string(config.Blue)+"%s: %s\n"+string(config.Reset), code, lang)
	}

	// Solicitar al usuario que seleccione un idioma
	var targetLang string
	for {
		fmt.Print(string(config.Bold) + string(config.Yellow) + "Ingrese el código del idioma de destino: " + string(config.Reset))
		fmt.Scanln(&targetLang)
		if _, valid := langCodes[targetLang]; valid {
			fmt.Println(string(config.Blue) + "Idioma seleccionado: " + string(targetLang))
			break
		} else {
			fmt.Println("Código de idioma no válido. Por favor, intente de nuevo.")
		}
	}

	// Listar archivos en el directorio de entrada
	files, err := utils.ListFilesInDir(config.InputDir)
	if err != nil {
		log.Fatalf("Error al listar archivos en el directorio de entrada: %v", err)
	}

	// Mostrar archivos a traducir
	fmt.Println("Archivos a traducir:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	// Solicitar al usuario que confirme la traducción
	fmt.Print("Ingrese 1 continuar o 0 para salir: ")
	var choice int
	fmt.Scanln(&choice)
	if choice == 0 {
		fmt.Println("Saliendo...")
		return
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))

	for _, inputFile := range files {
		wg.Add(1)
		go func(inputFile string) {
			defer wg.Done()
			process.ProcessFile(inputFile, targetLang)
		}(inputFile)
	}

	wg.Wait()
	close(fileChan)
}

func init() {
	// Limpiar la consola
	clearConsole()

	// Mostrar mensajes de bienvenida con colores y estilos
	fmt.Println(string(config.Bold) + string(config.Blue) + "╔═══════════════════════════════════════════════════════════════════════════╗" + string(config.Reset))
	fmt.Println(string(config.Bold) + string(config.Green) + "║ Bienvenido al programa de traducción de archivos SRT. 		   ║" + string(config.Reset))
	fmt.Println(string(config.Bold) + string(config.Yellow) + "  ------------------------------------------------------------------------" + string(config.Reset))
	fmt.Println(string(config.Bold) + string(config.Blue) + "║ Este programa le ayudará a traducir archivos SRT a diferentes idiomas.    ║" + string(config.Reset))
	fmt.Println(string(config.Bold) + string(config.Blue) + "╚═══════════════════════════════════════════════════════════════════════════╝" + string(config.Reset))
	fmt.Println(string(config.Bold) + string(config.Red) + "Presione Enter para continuar..." + string(config.Reset))
	fmt.Scanln()
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
