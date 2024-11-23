package usecase

import (
	"context"
	"encoding/json"
	"neuro-most/ai-service/internal/adapters/vllm"
	"strings"
)

type (
	ChatUseCase interface {
		Execute(ctx context.Context, input ChatInput) (ChatOutput, error)
	}

	ChatInput struct {
		Message string
	}

	ChatOutput struct {
		Response string
	}

	chatInteractor struct {
		vllm vllm.Vllm
	}
)

func NewChatUseCase(vllm vllm.Vllm) ChatUseCase {
	return &chatInteractor{vllm: vllm}
}

func (uc chatInteractor) Execute(ctx context.Context, input ChatInput) (ChatOutput, error) {
	systemPrompt := "Your task is to answer the user's questions using only the information from the provided documents. Give two answers to each question: one with a list of relevant document identifiers and the second with the answer to the question itself, using documents with these identifiers."
	documents := []vllm.Document{
		{
			DocID:   0,
			Title:   "Глобальное потепление: ледники",
			Content: "За последние 50 лет объем ледников в мире уменьшился на 30%",
		},
		{
			DocID:   1,
			Title:   "Глобальное потепление: Уровень моря",
			Content: "Уровень мирового океана повысился на 20 см с 1880 года и продолжает расти на 3,3 мм в год",
		},
	}
	docsJson, err := json.Marshal(documents)
	if err != nil {
		return ChatOutput{}, err
	}

	messages := []vllm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "documents", Content: string(docsJson)},
		{Role: "user", Content: input.Message},
	}

	relevantIndexes, err := uc.vllm.MakeVLLMRequest(messages, 0.0)
	if err != nil {
		return ChatOutput{}, err
	}

	messages = append(messages, vllm.Message{
		Role:    "assistant",
		Content: relevantIndexes,
	})

	finalAnswer, err := uc.vllm.MakeVLLMRequest(messages, 0.3)
	if err != nil {
		return ChatOutput{}, err
	}

	return ChatOutput{Response: strings.ReplaceAll(finalAnswer, "assistant", "")}, nil
}
