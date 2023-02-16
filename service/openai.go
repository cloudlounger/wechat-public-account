package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ASAK = "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"

var defaultPayload *Payload

func init() {
	defaultPayload = NewPayload()
}

type Payload struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int64   `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int64   `json:"top_p"`
	FrequencyPenalty int64   `json:"frequency_penalty"`
	PresencePenalty  int64   `json:"presence_penalty"`
	Model            string  `json:"model"`
}

func NewPayload() *Payload {
	return &Payload{
		MaxTokens:        2048,
		Temperature:      0.5,
		TopP:             0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "text-davinci-003",
	}
}

func (s *Payload) SendMessage(prompt string) (respWord string, err error) {
	s.Prompt = prompt
	payloadBytes, err := json.Marshal(s)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", ASAK))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	respWord = string(b)
	return
}
