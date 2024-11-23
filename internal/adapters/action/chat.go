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

func (a *ChatAction) Execute(input *aiv1.CreateChatRequest, stream aiv1.AIService_ChatServer) error {
	var usecaseInput usecase.ChatInput
	usecaseInput.Message = input.Message

	return a.uc.Execute(context.TODO(), usecaseInput, func(text string) error {
		return stream.Send(&aiv1.ChatResponse{
			Message: text,
		})
	})
}
