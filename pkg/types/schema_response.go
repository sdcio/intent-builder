package types

import "github.com/c-bata/go-prompt"

type SchemaResponse struct {
	SchemaResponseChild
	Childs []*SchemaResponseChild
	Keys   []string
}

func NewSchemaResponse(name, description string) *SchemaResponse {
	return &SchemaResponse{
		SchemaResponseChild: SchemaResponseChild{
			Name:        name,
			Description: description,
		},
	}
}

func (sr *SchemaResponse) ChildsToSuggestSlice() []prompt.Suggest {
	result := make([]prompt.Suggest, 0, len(sr.Childs))

	for _, x := range sr.Childs {
		result = append(result, prompt.Suggest{Text: x.Name, Description: x.Description})
	}
	return result
}

func (sr *SchemaResponse) Merge(o *SchemaResponse) {
	sr.Childs = append(sr.Childs, o.Childs...)
	sr.Keys = append(sr.Keys, o.Keys...)
}
