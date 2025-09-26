package interfaces

import (
	"context"

	"github.com/sdcio/intent-builder/pkg/types"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type CompleterSchemaClient interface {
	GetSchemaSdcpbPath(ctx context.Context, path *sdcpb.Path) (*sdcpb.GetSchemaResponse, error)
	ToPath(ctx context.Context, path []string) (*sdcpb.Path, error)
	GetSchemaSdcpbElemPath(ctx context.Context, path *sdcpb.Path) (*types.SchemaResponse, error)
}
