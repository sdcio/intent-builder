package clients

import (
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

type DataServerClient struct {
	c sdcpb.DataServerClient
}

func NewDataServerClient(target string) (*DataServerClient, error) {
	// setup grpc connection
	conn, err := NewGrpcClient(target)
	if err != nil {
		return nil, err
	}

	// instantiate SchemaServerClient
	ssc := &DataServerClient{
		c: sdcpb.NewDataServerClient(conn),
	}
	return ssc, nil
}
