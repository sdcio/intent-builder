package clients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(target string) (grpc.ClientConnInterface, error) {

	opts := []grpc.DialOption{}

	// only insecure for now
	opts = append(opts,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
