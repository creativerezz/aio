package ai

import (
	"context"

	"github.com/creativerezz/aio/internal/chat"
	"github.com/creativerezz/aio/internal/plugins"

	"github.com/creativerezz/aio/internal/domain"
)

type Vendor interface {
	plugins.Plugin
	ListModels() ([]string, error)
	SendStream([]*chat.ChatCompletionMessage, *domain.ChatOptions, chan domain.StreamUpdate) error
	Send(context.Context, []*chat.ChatCompletionMessage, *domain.ChatOptions) (string, error)
	NeedsRawMode(modelName string) bool
}
