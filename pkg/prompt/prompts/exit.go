package prompts

import (
	"context"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
)

type PromptExit struct{}

func NewPromptsExit() *PromptExit {
	return &PromptExit{}
}

func (p *PromptExit) Complete(_ context.Context, in []string, level int) []prompt.Suggest {
	return []prompt.Suggest{}
}
func (p *PromptExit) Execute(_ context.Context, _ []string, _ int) {
	fmt.Println("Bye!")
	os.Exit(0)
}
