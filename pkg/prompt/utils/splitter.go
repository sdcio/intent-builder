package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// Split collapses multiple spaces to one and treats text in quotes as a single element.
func Split(in string) (*SplitResult, error) {

	var result []string
	inQuote := false
	var currentToken strings.Builder

	for i, r := range in {
		if r == '"' {
			inQuote = !inQuote
			if currentToken.Len() > 0 {
				result = append(result, currentToken.String())
				currentToken.Reset()
			}
			continue
		}

		if inQuote {
			currentToken.WriteRune(r)
		} else {
			if unicode.IsSpace(r) {
				if currentToken.Len() > 0 {
					result = append(result, currentToken.String())
					currentToken.Reset()
				}
				// Skip multiple spaces
				continue
			} else {
				currentToken.WriteRune(r)
			}
		}

		// Handle the last character
		if i == len(in)-1 && currentToken.Len() > 0 {
			result = append(result, currentToken.String())
		}
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("not enough parameters to split")
	}

	return NewSplitResult(result[:len(result)-1], result[len(result)-1]), nil

}

type SplitResult struct {
	path  []string
	value string
}

func NewSplitResult(path []string, value string) *SplitResult {
	return &SplitResult{
		path:  path,
		value: value,
	}
}

func (s *SplitResult) GetPath() []string {
	return s.path
}

func (s *SplitResult) GetValue() string {
	return s.value
}
