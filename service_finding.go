package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ca-risken/core/proto/finding"
)

func (g *gatewayService) getAISummaryStreamHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &finding.GetAISummaryRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GetAISummaryRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	stream, err := g.findingClient.GetAISummaryStream(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}

	// Server-Sent Events (SSE)
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "Streaming unsupported!"})
		return
	}
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
				appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
				writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
			}
			return
		}
		// Send data periodically
		fmt.Fprintf(w, "%s", resp.Answer)
		flusher.Flush()
		time.Sleep(1 * time.Millisecond)
	}
}
