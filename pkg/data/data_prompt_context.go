package data

import (
	"strings"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type DataPromptContext struct {
	actualPath   *sdcpb.Path
	outputFormat OutputFormat
}

func NewDataPromptContext() *DataPromptContext {
	return &DataPromptContext{
		actualPath:   &sdcpb.Path{},
		outputFormat: OutputFormatJSON,
	}
}

type OutputFormat string

const (
	OutputFormatString    OutputFormat = "string"
	OutputFormatJSON      OutputFormat = "json"
	OutputFormatJSON_IETF OutputFormat = "json_ietf"
	OutputFormatXML       OutputFormat = "xml"
)

func ParseOutputFormat(s string) OutputFormat {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "string":
		return OutputFormatString
	case "json":
		return OutputFormatJSON
	case "json_ietf":
		return OutputFormatJSON_IETF
	case "xml":
		return OutputFormatXML
	}
	// fallback
	return OutputFormatString
}
