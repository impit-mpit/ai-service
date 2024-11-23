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
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
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

		// Remove "data: " prefix if present
		data := bytes.TrimPrefix(line, []byte("data: "))

		var streamResp StreamResponse
		if err := json.Unmarshal(data, &streamResp); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\nRaw data: %s\n", err, string(data))
			continue
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				fmt.Printf("Got chunk from VLLM: %s\n", content)
				if err := stream(content); err != nil {
					return fmt.Errorf("error streaming response: %v", err)
				}
			}

			if streamResp.Choices[0].FinishReason == "stop" {
				break
			}
		}
	}

	return nil
}
