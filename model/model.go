package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"newshub/config"
	"strings"
	"time"
)

func Init() {

}

func Chat(messages []map[string]string) (string, error) {
	var client *http.Client = &http.Client{
		Timeout: 15 * 60 * time.Second,
	}
	request := map[string]any{
		"model":    config.Model,
		"messages": messages,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.ModelUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ModelKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return strings.TrimLeft(result.Choices[0].Message.Content, "\n"), nil
}
