package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func GetRandomFact() (string, error) {
	var fact Fact
	resp, err := http.Get("https://catfact.ninja/fact")
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
	if err := json.Unmarshal(body, &fact); err != nil {
		fmt.Println("Ошибка:", err)
		return "", nil
	}

	return fact.Fact, nil
}
