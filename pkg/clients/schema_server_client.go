package clients

import (
	"context"

	schemaClient "github.com/sdcio/data-server/pkg/datastore/clients/schema"

	"github.com/sdcio/intent-builder/pkg/types"

	backend_interfaces "github.com/sdcio/intent-builder/pkg/data/interfaces"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type SchemaServerClient struct {
	scb schemaClient.SchemaClientBound
}

func NewSchemaServerClient(scb schemaClient.SchemaClientBound) (*SchemaServerClient, error) {

	// instantiate SchemaServerClient
	ssc := &SchemaServerClient{
		scb: scb,
	}
	return ssc, nil
}

// make sure SchemaServerClient implements prompt.PromptSchemaClient
var _ backend_interfaces.CompleterSchemaClient = (*SchemaServerClient)(nil)

func (ssc *SchemaServerClient) GetSchemaSdcpbPath(ctx context.Context, path *sdcpb.Path) (*sdcpb.GetSchemaResponse, error) {
	rsp, err := ssc.scb.GetSchemaSdcpbPath(ctx, path)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (ssc *SchemaServerClient) ToPath(ctx context.Context, path []string) (*sdcpb.Path, error) {

	sdcpbPath := &sdcpb.Path{Elem: []*sdcpb.PathElem{}}

	for idx := 0; idx < len(path); idx++ {

		newPathElem := sdcpb.NewPathElem(path[idx], nil)
		sdcpbPath.Elem = append(sdcpbPath.Elem, newPathElem)
		schema, err := ssc.GetSchemaSdcpbPath(ctx, sdcpbPath)
		if err != nil {
			return nil, err
		}
		switch specifcSchema := schema.GetSchema().GetSchema().(type) {
		case *sdcpb.SchemaElem_Container:

			keys := specifcSchema.Container.GetKeys()
			// if we have keys, init the keys map
			if len(keys) > 0 {
				newPathElem.Key = make(map[string]string, len(keys))
			}
			// add the keys to the map
			for keyIdx, k := range keys {
				idx++
				// break if e.g. a key is missing yet
				if idx+keyIdx > len(path)-1 {
					break
				}
				newPathElem.Key[k.Name] = path[idx]
			}
			// case *sdcpb.SchemaElem_Field, *sdcpb.SchemaElem_Leaflist:
		}
	}

	return sdcpbPath, nil

}

func (ssc *SchemaServerClient) GetSchemaSdcpbElemPath(ctx context.Context, path *sdcpb.Path) (*types.SchemaResponse, error) {

	rsp, err := ssc.scb.GetSchemaSdcpbPath(ctx, path)
	if err != nil {
		return nil, err
	}

	// check for the root path, that need special handling
	if len(path.Elem) == 0 {
		result := types.NewSchemaResponse("root", "")
		for _, module := range rsp.GetSchema().GetContainer().Children {
			modRsp, err := ssc.scb.GetSchemaSdcpbPath(ctx, &sdcpb.Path{Elem: []*sdcpb.PathElem{sdcpb.NewPathElem(module, nil)}})
			if err != nil {
				return nil, err
			}
			sr, err := SdcpbSchemaRespToSchemaResp(modRsp.GetSchema())
			if err != nil {
				return nil, err
			}
			result.Merge(sr)
		}
		return result, nil
	}

	// if it is a container that is is list (with key) then check if all the keys are defined.
	if rsp.Schema.GetContainer() != nil && len(rsp.GetSchema().GetContainer().GetKeys()) > 0 && len(path.Elem[len(path.Elem)-1].Key) < len(rsp.Schema.GetContainer().Keys) {
		return &types.SchemaResponse{}, nil
	}

	// otherwise regular processing
	return SdcpbSchemaRespToSchemaResp(rsp.GetSchema())
}

func SdcpbSchemaRespToSchemaResp(se *sdcpb.SchemaElem) (*types.SchemaResponse, error) {
	var result *types.SchemaResponse

	switch s := se.GetSchema().(type) {
	case *sdcpb.SchemaElem_Container:
		x := s.Container
		result = types.NewSchemaResponse(x.GetName(), x.GetDescription())
		for _, c := range x.GetKeys() {
			result.Childs = append(result.Childs, types.NewSchemaResponseChild(c.GetName(), "Key: "+c.GetDescription()))
		}
		for _, c := range x.GetChildren() {
			result.Childs = append(result.Childs, types.NewSchemaResponseChild(c, ""))
		}
		for _, f := range x.GetFields() {
			result.Childs = append(result.Childs, types.NewSchemaResponseChild(f.Name, f.Description))
		}
	case *sdcpb.SchemaElem_Field:
		x := s.Field
		result = types.NewSchemaResponse(x.GetName(), x.GetDescription())
		for _, elem := range x.GetType().GetEnumNames() {
			result.Childs = append(result.Childs, types.NewSchemaResponseChild(elem, ""))
		}

		if len(result.Childs) == 0 {
			result.Childs = append(result.Childs, types.NewSchemaResponseChild("", x.GetType().Type))
		}

	case *sdcpb.SchemaElem_Leaflist:
		x := s.Leaflist
		result = types.NewSchemaResponse(x.GetName(), x.GetDescription())
	}
	return result, nil
}
