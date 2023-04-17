package main

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCErrorMessage(t *testing.T) {
	cases := []struct {
		name  string
		input *status.Status
		want  map[string]interface{}
	}{
		{
			name:  "gRPC error",
			input: status.New(codes.Internal, "internal server error"),
			want: map[string]interface{}{
				errorJSONKey: grpcError{
					Code:    codes.Internal.String(),
					Message: "internal server error",
				},
			},
		},
		{
			name:  "status nil",
			input: nil,
			want: map[string]interface{}{
				errorJSONKey: grpcError{
					Code:    codes.Unknown.String(),
					Message: "gRPC status is nil",
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := grpcErrorMessage(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
