package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"srt-translator/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Modelos y sus límites de tokens
var models = map[string]int{
	"gpt-3.5-turbo":     4096,
	"gpt-4":             8192,
	"gpt-3.5-turbo-16k": 16384,
}

// Cliente HTTP reutilizable
var httpClient = &http.Client{}

// Depuración habilitada
var debug = true

func TranslateSRT(srtContent, targetLang string) (string, error) {
	apiKey := config.APIKeyEnvVar
	if apiKey == "" {
		return "", fmt.Errorf("API key is not set")
	}

	lines := strings.Split(srtContent, "\n")
	var translatedLines []string
	var chunk []string
	var chunkSize int

	debugLog("Iniciando la traducción del archivo SRT...")

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if isTimecode(line) || line == "" || isIndex(line) {
			// Agregar líneas de tiempo y líneas vacías directamente
			if len(chunk) > 0 {
				model, _ := selectModel(chunkSize)
				debugLog("Enviando chunk para traducción con el modelo:", model)
				translatedChunk, err := translateChunk(strings.Join(chunk, "\n"), targetLang, apiKey, model)
				if err != nil {
					return "", err
				}
				translatedLines = append(translatedLines, strings.Split(translatedChunk, "\n")...)
				chunk = []string{}
				chunkSize = 0
			}
			translatedLines = append(translatedLines, line)
		} else {
			model, maxTokens := selectModel(chunkSize + len(line))
			if chunkSize+len(line) > maxTokens {
				debugLog("Enviando chunk para traducción con el modelo:", model)
				translatedChunk, err := translateChunk(strings.Join(chunk, "\n"), targetLang, apiKey, model)
				if err != nil {
					return "", err
				}
				translatedLines = append(translatedLines, strings.Split(translatedChunk, "\n")...)
				chunk = []string{}
				chunkSize = 0
			}
			chunk = append(chunk, line)
			chunkSize += len(line)
		}
	}

	if len(chunk) > 0 {
		model, _ := selectModel(chunkSize)
		debugLog("Enviando último chunk para traducción con el modelo:", model)
		translatedChunk, err := translateChunk(strings.Join(chunk, "\n"), targetLang, apiKey, model)
		if err != nil {
			return "", err
		}
		translatedLines = append(translatedLines, strings.Split(translatedChunk, "\n")...)
	}

	debugLog("Traducción completada.")
	return strings.Join(translatedLines, "\n"), nil
}

func isTimecode(line string) bool {
	return strings.Contains(line, "-->")
}

func isIndex(line string) bool {
	_, err := strconv.Atoi(line)
	return err == nil
}

func selectModel(tokenCount int) (string, int) {
	for model, maxTokens := range models {
		if tokenCount <= maxTokens {
			return model, maxTokens
		}
	}
	// Si el tokenCount excede todos los modelos disponibles, usar el modelo con el mayor límite de tokens
	var bestModel string
	var bestMaxTokens int
	for model, maxTokens := range models {
		if maxTokens > bestMaxTokens {
			bestModel = model
			bestMaxTokens = maxTokens
		}
	}
	return bestModel, bestMaxTokens
}

func translateChunk(chunk, targetLang, apiKey, model string) (string, error) {
	debugLog("Creando solicitud HTTP para chunk con el modelo:", model)

	requestBody, err := json.Marshal(OpenAIRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a translator that translates text to %s.", targetLang),
			},
			{
				Role:    "user",
				Content: chunk,
			},
		},
	})
	if err != nil {
		debugLog("Error al crear el cuerpo de la solicitud:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", config.APIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		debugLog("Error al crear la solicitud HTTP:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	debugLog("Enviando solicitud HTTP...")

	resp, err := httpClient.Do(req)
	if err != nil {
		debugLog("Error al enviar la solicitud HTTP:", err)
		return "", err
	}
	defer resp.Body.Close()

	debugLog("Solicitud HTTP enviada. Código de estado:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		debugLog("Error en la respuesta de la API:", fmt.Sprintf("unexpected status code: %d, response: %s", resp.StatusCode, responseBody))
		return "", fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, responseBody)
	}

	debugLog("Procesando respuesta de la API...")

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		debugLog("Error al decodificar la respuesta de la API:", err)
		return "", err
	}

	if len(response.Choices) == 0 {
		debugLog("No se encontró traducción en la respuesta de la API.")
		return "", fmt.Errorf("no translation found in response")
	}

	debugLog("Traducción obtenida exitosamente.")
	return response.Choices[0].Message.Content, nil
}

func debugLog(v ...interface{}) {
	if debug {
		fmt.Println(v...)
	}
}
