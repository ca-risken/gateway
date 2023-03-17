package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandleGrpcError(t *testing.T) {
	cases := []struct {
		name  string
		input error
		want  map[string]interface{}
	}{
		{
			name:  "gRPC error",
			input: status.Error(codes.Internal, "intternal server error"),
			want: map[string]interface{}{
				errorJSONKey: grpcError{
					Code:    codes.Internal.String(),
					Message: "intternal server error",
				},
			},
		},
		{
			name:  "other error",
			input: errors.New("something wrong"),
			want:  map[string]interface{}{errorJSONKey: "something wrong"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := grpcErrorMessage(context.TODO(), c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
