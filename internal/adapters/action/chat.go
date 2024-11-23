package action

import (
	"context"
	aiv1 "neuro-most/ai-service/gen/go/ai/v1"
	"neuro-most/ai-service/internal/usecase"
)

type ChatAction struct {
	uc usecase.ChatUseCase
}

func NewChatAction(uc usecase.ChatUseCase) *ChatAction {
	return &ChatAction{
		uc: uc,
	}
}

func (a *ChatAction) Execute(ctx context.Context, input *aiv1.CreateChatRequest) (*aiv1.ChatResponse, error) {
	var usecaseInput usecase.ChatInput
	usecaseInput.Message = input.Message
	text, err := a.uc.Execute(ctx, usecaseInput)
	if err != nil {
		return nil, err
	}
	return &aiv1.ChatResponse{
		Message: text.Response,
	}, nil
}
