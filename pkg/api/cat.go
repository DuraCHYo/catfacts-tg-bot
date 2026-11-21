package api

import (
	"encoding/json"
	"io"
	"log/slog"
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
		slog.Error("Ошибка обработка запроса к Catfact API", "status", resp.Status, "response", body)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Ошибка чтения ответа от Catfact API", "error", err)
		return "", nil
	}
	if err := json.Unmarshal(body, &fact); err != nil {
		slog.Error("Ошибка Catfact API", "error", err)
		return "", nil
	}

	return fact.Fact, nil
}
