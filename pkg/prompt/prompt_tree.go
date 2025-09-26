package prompt

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/sdcio/intent-builder/pkg/prompt/prompts"
)

type PromptTreeElementImpl struct {
	name    string
	childs  map[string]PromptTreeElement
	element prompts.PromptElement
}

func NewPromptTreeRoot() *PromptTreeElementImpl {
	return &PromptTreeElementImpl{
		name:   "Root",
		childs: map[string]PromptTreeElement{},
	}
}

func newPromptTreeElement(name string) *PromptTreeElementImpl {
	return &PromptTreeElementImpl{
		name:   name,
		childs: map[string]PromptTreeElement{},
	}
}

func (pte *PromptTreeElementImpl) GetName() string {
	return pte.name
}
func (pte *PromptTreeElementImpl) GetChildren() map[string]PromptTreeElement {
	return pte.childs
}
func (pte *PromptTreeElementImpl) GetPromptElement() prompts.PromptElement {
	return pte.element
}

func (pte *PromptTreeElementImpl) AddPromptElement(p prompts.PromptElement, path []string, level int) error {
	if level == len(path) {
		pte.element = p
		return nil
	}

	child, exists := pte.childs[path[level]]
	if !exists {
		child = newPromptTreeElement(path[level])
		pte.childs[path[level]] = child
	}
	return child.AddPromptElement(p, path, level+1)
}

func (pte *PromptTreeElementImpl) Complete(ctx context.Context, in []string, level int) []prompt.Suggest {
	if pte.element == nil {
		result := make([]prompt.Suggest, 0, len(pte.childs))
		for n := range pte.childs {
			result = append(result, prompt.Suggest{Text: n})
		}
		return result
	}
	return pte.element.Complete(ctx, in[level:], 0)
}

func (pte *PromptTreeElementImpl) Execute(ctx context.Context, in []string, level int) {
	if pte.element == nil {
		fmt.Println("NOOP")
		return
	}
	pte.element.Execute(ctx, in, level)
}

func (pte *PromptTreeElementImpl) Navigate(path []string, level int) (PromptTreeElement, int) {
	if level == len(path) {
		return pte, level
	}
	child, exist := pte.childs[path[level]]
	if exist {
		return child.Navigate(path, level+1)
	}
	// does not exist return actual as the PromptElement
	return pte, level

}
