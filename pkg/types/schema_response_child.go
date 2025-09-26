package types

type SchemaResponseChild struct {
	Name        string
	Description string
}

func NewSchemaResponseChild(name, description string) *SchemaResponseChild {
	return &SchemaResponseChild{
		Name:        name,
		Description: description,
	}
}
