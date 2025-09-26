package prompt

import (
	"context"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/sdcio/intent-builder/pkg/prompt/prompts"
)

var (
	history = []string{}
)

type Prompter struct {
	ctx            context.Context
	lastSuggestion []prompt.Suggest

	promptTreeRoot PromptTreeElement
}

func NewPrompter(ctx context.Context) *Prompter {
	return &Prompter{
		ctx:            ctx,
		promptTreeRoot: NewPromptTreeRoot(),
	}
}

func (pc *Prompter) AddPrompt(p prompts.PromptElement, paths [][]string) error {
	for _, path := range paths {
		pc.promptTreeRoot.AddPromptElement(p, path, 0)
	}
	return nil
}

func (pc *Prompter) Run() {
	// create new prompt with history
	p := prompt.New(pc.executor, pc.completerFunc, prompt.OptionHistory(history), prompt.OptionShowCompletionAtStart())
	// run the prompt
	p.Run()
}

func (pc *Prompter) executor(in string) {

	if in == "" {
		// return on no input
		return
	}

	// reset last suggestions
	pc.lastSuggestion = nil

	splitIn := strings.Split(in, " ")
	pte, level := pc.promptTreeRoot.Navigate(splitIn, 0)
	pte.Execute(pc.ctx, splitIn, level)
}

func (pc *Prompter) completerFunc(in prompt.Document) []prompt.Suggest {
	if pc.lastSuggestion != nil {
		switch in.LastKeyStroke() {
		case prompt.Tab, prompt.Escape:
			return pc.lastSuggestion
		}
	}

	splitIn := strings.Split(in.Text, " ")
	pte, level := pc.promptTreeRoot.Navigate(splitIn, 0)

	suggestions := pte.Complete(pc.ctx, splitIn, level)

	pc.lastSuggestion = prompt.FilterHasPrefix(suggestions, in.GetWordBeforeCursor(), true)
	return pc.lastSuggestion
}
