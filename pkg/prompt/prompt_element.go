package prompt

import "github.com/sdcio/intent-builder/pkg/prompt/prompts"

type PromptTreeElement interface {
	prompts.PromptElement
	GetName() string
	// GetPromptElement() prompts.PromptElement
	AddPromptElement(p prompts.PromptElement, path []string, level int) error
	Navigate(path []string, level_in int) (pte PromptTreeElement, level int)
}
