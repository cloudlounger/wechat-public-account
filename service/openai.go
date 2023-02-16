package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"wxcloudrun-golang/db/model"
)

var defaultPayload *Payload
var cache *sync.Map
var result *sync.Map
var queue chan *model.WXMessage

func init() {
	defaultPayload = NewPayload()
	cache = new(sync.Map)
	result = new(sync.Map)
	queue = make(chan *model.WXMessage, 2)
	go func() {
		for {
			msg := <-queue
			key := getKey(msg)
			word := SendAsync(msg)
			result.Store(key, word)
			cache.Delete(key)
		}
	}()
}

func pushQueue(msg *model.WXMessage) {
	queue <- msg
}

func loopCheck(key string, cancelC <-chan struct{}) (quit bool) {
	for {
		select {
		case <-cancelC:
			quit = true
			return
		case <-time.After(5 * time.Second):
			quit = true
			return
		default:
		}
		if _, ok := cache.Load(key); ok {
			time.Sleep(200 * time.Millisecond)
		} else {
			return
		}
	}
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
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", "sk-"+"dr6u6AP1XgyCXo0Ss7kZT3BlbkFJrooo7cpexDmS317DjLQe"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("not 200 status code")
		fmt.Println("[debug] error status code, code", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	aiResp := new(AIResp)
	err = json.Unmarshal(b, aiResp)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	if len(aiResp.Choices) <= 0 {
		err = errors.New("nothing get")
		fmt.Println("[debug] error", err)
		return
	}
	return aiResp.Choices[0].Text, nil
}

type AIResp struct {
	ID          string       `json:"id"`
	Object      string       `json:"object"`
	CreatedTime int          `json:"created"`
	Choices     []*AIContent `json:"choices"`
}

type AIContent struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}
