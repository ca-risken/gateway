package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func grpcErrorMessage(stat *status.Status) map[string]interface{} {
	if stat == nil {
		return map[string]interface{}{errorJSONKey: grpcError{
			Code:    codes.Unknown.String(),
			Message: "gRPC status is nil",
		}}
	}
	return map[string]interface{}{errorJSONKey: grpcError{
		Code:    stat.Code().String(),
		Message: stat.Message(),
	}}
}

func handleGRPCError(ctx context.Context, w http.ResponseWriter, err error) error {
	if err == nil {
		return errors.New("required error")
	}
	stat, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("unknown error: %w", err)
	}
	if stat.Code() == codes.OK {
		return errors.New("no error")
	}
	appLogger.Warnf(ctx, "gRPC error, code=%s, message=%s", stat.Code(), stat.Message())

	// handling by error code
	// ref: https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md
	switch stat.Code() {
	case codes.Canceled:
		writeResponse(ctx, w, 499, grpcErrorMessage(stat)) // 499 Client Closed Request
	case codes.Unknown:
		writeResponse(ctx, w, http.StatusInternalServerError, grpcErrorMessage(stat))
	case codes.InvalidArgument:
		writeResponse(ctx, w, http.StatusBadRequest, grpcErrorMessage(stat))
	case codes.DeadlineExceeded:
		writeResponse(ctx, w, http.StatusGatewayTimeout, grpcErrorMessage(stat))
	case codes.NotFound:
		writeResponse(ctx, w, http.StatusNotFound, grpcErrorMessage(stat))
	case codes.AlreadyExists:
		writeResponse(ctx, w, http.StatusConflict, grpcErrorMessage(stat))
	case codes.PermissionDenied:
		writeResponse(ctx, w, http.StatusForbidden, grpcErrorMessage(stat))
	case codes.Unauthenticated:
		writeResponse(ctx, w, http.StatusUnauthorized, grpcErrorMessage(stat))
	case codes.ResourceExhausted:
		writeResponse(ctx, w, http.StatusTooManyRequests, grpcErrorMessage(stat))
	case codes.FailedPrecondition:
		writeResponse(ctx, w, http.StatusPreconditionFailed, grpcErrorMessage(stat))
	case codes.Aborted:
		writeResponse(ctx, w, http.StatusConflict, grpcErrorMessage(stat))
	case codes.OutOfRange:
		writeResponse(ctx, w, http.StatusBadRequest, grpcErrorMessage(stat))
	case codes.Unimplemented:
		writeResponse(ctx, w, http.StatusNotImplemented, grpcErrorMessage(stat))
	case codes.Internal:
		writeResponse(ctx, w, http.StatusInternalServerError, grpcErrorMessage(stat))
	case codes.Unavailable:
		writeResponse(ctx, w, http.StatusServiceUnavailable, grpcErrorMessage(stat))
	case codes.DataLoss:
		writeResponse(ctx, w, http.StatusInternalServerError, grpcErrorMessage(stat))
	default:
		writeResponse(ctx, w, http.StatusInternalServerError, grpcErrorMessage(stat))
	}
	return nil
}
