package gapi

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("received a gRPC request")
	result, err := handler(ctx, req)
	return result, err
}
