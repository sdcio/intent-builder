package prompts

import (
	"context"

	"github.com/c-bata/go-prompt"
)

type PromptElement interface {
	Complete(ctx context.Context, in []string, level int) []prompt.Suggest
	Execute(ctx context.Context, in []string, level int)
}
