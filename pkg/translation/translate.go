package translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Translate struct {
	Text string `json:"translatedText"`
}

func TranslateText(text string) (string, error) {
	var translated Translate
	payload := map[string]string{
		"q":      text,
		"source": "en",
		"target": "ru",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://localhost:5000/translate", "application/json", bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		result := "Ошибка обработка запроса: " + resp.Status + " " + string(body)
		return result, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
	}
	if err := json.Unmarshal(body, &translated); err != nil {
		fmt.Println("Ошибка:", err)
		return "", nil
	}
	return translated.Text, nil
}
