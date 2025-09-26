package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/sdcio/intent-builder/pkg/prompt/utils"
)

type DataPrompt struct {
	backend *Backend
	context *DataPromptContext
}

func NewDataPrompt(backend *Backend) *DataPrompt {
	return &DataPrompt{
		backend: backend,
		context: NewDataPromptContext(),
	}
}

func (dp *DataPrompt) Complete(ctx context.Context, in []string, level int) []prompt.Suggest {
	// delegate to completer
	suggestions, err := dp.backend.Complete(ctx, strings.Join(in[level:], " "))
	if err != nil {
		fmt.Println("\n", err)
		return nil
	}
	return suggestions
}

func (dp *DataPrompt) Execute(ctx context.Context, in []string, level int) {
	switch in[level-1] {
	case "set":
		dp.set(ctx, in, level)
	case "show":
		dp.show(ctx, in, level)
	case "delete":
		dp.delete(ctx, in, level)
	}
}

func (dp *DataPrompt) show(ctx context.Context, _ []string, _ int) {
	result, err := dp.backend.ToFormat(ctx, dp.context.actualPath, dp.context.outputFormat)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func (db *DataPrompt) delete(_ context.Context, _ []string, _ int) {
	fmt.Println("! delete - not implemented yet !")
}

func (dp *DataPrompt) set(ctx context.Context, in []string, level int) {
	splitResult, err := utils.Split(strings.Join(in[level:], " "))
	if err != nil {
		// TODO Need to handle error
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = dp.backend.AddLine(ctx, splitResult.GetPath(), splitResult.GetValue(), 5)
	if err != nil {
		// TODO Need to handle error
		fmt.Printf("Error: %v\n", err)
		return
	}
}
