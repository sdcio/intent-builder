package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/sdcio/data-server/pkg/tree"
	"github.com/sdcio/data-server/pkg/tree/types"
	"github.com/sdcio/data-server/pkg/utils"
	backend_interfaces "github.com/sdcio/intent-builder/pkg/data/interfaces"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type Backend struct {
	dc   backend_interfaces.CompleterDataClient
	sc   backend_interfaces.CompleterSchemaClient
	root *tree.RootEntry
}

func NewBackend(ctx context.Context, dc backend_interfaces.CompleterDataClient, sc backend_interfaces.CompleterSchemaClient, tc *tree.TreeContext) (*Backend, error) {
	root, err := tree.NewTreeRoot(ctx, tc)
	if err != nil {
		return nil, err
	}

	return &Backend{
		dc:   dc,
		sc:   sc,
		root: root,
	}, nil
}

func (b *Backend) Complete(ctx context.Context, input string) ([]prompt.Suggest, error) {
	skipLastElems := 1
	if len(input) > 0 && input[len(input)-1] == ' ' {
		skipLastElems = 0
	}

	path := strings.Split(strings.TrimSpace(input), " ")

	sdcpbPath, err := b.sc.ToPath(ctx, path[:len(path)-skipLastElems])
	if err != nil {
		return nil, err
	}
	rsp, err := b.sc.GetSchemaSdcpbElemPath(ctx, sdcpbPath)
	if err != nil {
		return nil, err
	}

	return rsp.ChildsToSuggestSlice(), nil
}

func (b *Backend) AddLine(ctx context.Context, path []string, value string, prio int32) error {

	sdcpbPath, err := b.sc.ToPath(ctx, path)
	if err != nil {
		return err
	}

	schemaResp, err := b.sc.GetSchemaSdcpbPath(ctx, sdcpbPath)
	if err != nil {
		return err
	}

	tv, err := utils.ConvertToTypedValue(schemaResp.GetSchema(), value, 0)
	if err != nil {
		return err
	}

	_, err = b.root.AddUpdateRecursive(ctx, sdcpbPath, types.NewUpdate(sdcpbPath, tv, prio, "NewIntent", 0), types.NewUpdateInsertFlags().SetNewFlag())
	if err != nil {
		return err
	}
	return nil
}

func (b *Backend) String() string {
	return b.root.String()
}

func (b *Backend) ToFormat(ctx context.Context, path *sdcpb.Path, of OutputFormat) (string, error) {

	e, err := b.root.NavigateSdcpbPath(ctx, path)
	if err != nil {
		return "", err
	}

	switch of {
	case OutputFormatJSON:
		j, err := e.ToJson(false)
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		b, err := json.MarshalIndent(j, "", " ")
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		return string(b), nil
	case OutputFormatJSON_IETF:
		j, err := e.ToJson(false)
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		b, err := json.MarshalIndent(j, "", " ")
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		return string(b), nil
	case OutputFormatXML:
		et, err := e.ToXML(false, true, false, false)
		if err != nil {
			return "", err
		}
		et.Indent(2)
		return et.WriteToString()
	case OutputFormatString:
		return strings.Join(e.StringIndent(nil), "\n"), nil
	}
	return "", fmt.Errorf("unknown format")
}
