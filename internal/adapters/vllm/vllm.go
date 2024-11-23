package vllm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Document struct {
	DocID   int    `json:"doc_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type VLLMRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type VLLMResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type StreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Text         string  `json:"text"`
		LogProbs     *string `json:"logprobs"`
		FinishReason *string `json:"finish_reason"`
		StopReason   *string `json:"stop_reason"`
	} `json:"choices"`
	Usage *string `json:"usage"`
}

type Vllm struct {
	vllmURL string
	apiKey  string
}

func NewVllm(vllmURL, apiKey string) *Vllm {
	return &Vllm{
		vllmURL: vllmURL,
		apiKey:  apiKey,
	}
}

// Обычный запрос для получения индексов
func (s *Vllm) MakeVLLMRequest(messages []Message, temperature float64) (string, error) {
	var prompt string
	for _, msg := range messages {
		prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	vllmReq := VLLMRequest{
		Model:       "Vikhrmodels/Vikhr-Nemo-12B-Instruct-R-21-09-24",
		Prompt:      prompt,
		MaxTokens:   2048,
		Temperature: temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(vllmReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", s.vllmURL+"/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var vllmResp VLLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&vllmResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(vllmResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return vllmResp.Choices[0].Text, nil
}

// Стриминг запрос для получения ответа
func (s *Vllm) MakeVLLMStreamRequest(messages []Message, temperature float64, stream func(string) error) error {
	var prompt string
	for _, msg := range messages {
		prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	vllmReq := VLLMRequest{
		Model:       "Vikhrmodels/Vikhr-Nemo-12B-Instruct-R-21-09-24",
		Prompt:      prompt,
		MaxTokens:   2048,
		Temperature: temperature,
		Stream:      true,
	}

	jsonData, err := json.Marshal(vllmReq)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", s.vllmURL+"/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading stream: %v", err)
		}

		// Skip empty lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		// Remove "data: " prefix
		data := bytes.TrimPrefix(line, []byte("data: "))

		var streamResp StreamResponse
		if err := json.Unmarshal(data, &streamResp); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			continue
		}
		if len(streamResp.Choices) > 0 && streamResp.Choices[0].Text != "" {
			text := streamResp.Choices[0].Text
			fmt.Printf("Got chunk: %s\n", text)
			if err := stream(text); err != nil {
				return fmt.Errorf("error streaming response: %v", err)
			}
		}

		// Проверяем finish_reason
		if len(streamResp.Choices) > 0 &&
			streamResp.Choices[0].FinishReason != nil &&
			*streamResp.Choices[0].FinishReason == "stop" {
			break
		}
	}

	return nil
}
