package main

import (
	"context"

	"google.golang.org/grpc/status"
)

type grpcError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func grpcErrorMessage(ctx context.Context, err error) map[string]interface{} {
	if stat, ok := status.FromError(err); ok {
		appLogger.Warnf(ctx, "gRPC error, code=%s, message=%s", stat.Code(), stat.Message())
		// TODO: handling by error code
		return map[string]interface{}{errorJSONKey: grpcError{
			Code:    stat.Code().String(),
			Message: stat.Message(),
		}}
	}
	return map[string]interface{}{errorJSONKey: err.Error()}
}
