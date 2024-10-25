
# SRT Translator

## Descripción del Proyecto

SRT Translator es una herramienta avanzada diseñada para traducir archivos de subtítulos (SRT) a múltiples idiomas utilizando la potente API de OpenAI. Este proyecto está orientado a facilitar la traducción precisa y eficiente de subtítulos, manteniendo el formato original del archivo SRT.

## Resumen del Funcionamiento

El proceso de traducción de SRT Translator se realiza en varios pasos clave:

1. **Lectura del Archivo SRT**: El archivo SRT se lee línea por línea, identificando y separando las líneas de tiempo y las líneas de texto.
2. **Traducción del Texto**: Solo las líneas de texto se envían a la API de OpenAI para su traducción, mientras que las líneas de tiempo se mantienen intactas. Esta forma ayuda a minimizar el gasto de token en cada peticion.
3. **Reconstrucción del Archivo**: Las líneas de tiempo y las líneas de texto traducidas se combinan para reconstruir el archivo SRT en el idioma de destino.

Este enfoque asegura que el formato del archivo SRT se mantenga intacto, proporcionando una traducción precisa y bien estructurada.

## Instalación

1. Clona el repositorio:
    
    ```sh
    git clone https://github.com/tu-usuario/srt-translator.git
    cd srt-translator
    ```

2. Instala las dependencias:
    ```go
    go mod tidy
    ```

## Configuración

Asegúrate de configurar tu clave API de OpenAI en el archivo `config/config.go`:
```go
const (
    APIKeyEnvVar = "tu-clave-api"
    APIURL       = "https://api.openai.com/v1/chat/completions"
)
```

## Uso

1. Ejecuta el programa:
    ```sh
    go run main.go
    ```

2. Sigue las instrucciones en la consola para seleccionar los archivos SRT que deseas traducir y el idioma de destino.

## Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o un pull request para discutir cualquier cambio que te gustaría realizar.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Consulta el archivo LICENSE para más detalles.