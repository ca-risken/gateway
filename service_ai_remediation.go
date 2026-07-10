package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protowire"
)

const generateRemediationProposalFullMethod = "/datasource.ai.AIService/GenerateRemediationProposal"

type remediationProposalGenerator interface {
	GenerateRemediationProposal(context.Context, *generateRemediationProposalRequest) (*generateRemediationProposalResponse, error)
}

type aiRemediationClient struct {
	cc grpc.ClientConnInterface
}

type generateRemediationProposalRequest struct {
	ProjectID uint32 `json:"project_id"`
	FindingID uint64 `json:"finding_id"`
}

type generateRemediationProposalResponse struct {
	RemediationProposalID uint32 `json:"remediation_proposal_id,omitempty"`
}

func newAIRemediationClient(cc grpc.ClientConnInterface) remediationProposalGenerator {
	return &aiRemediationClient{cc: cc}
}

func (c *aiRemediationClient) GenerateRemediationProposal(ctx context.Context, req *generateRemediationProposalRequest) (*generateRemediationProposalResponse, error) {
	resp := &generateRemediationProposalResponse{}
	if err := c.cc.Invoke(ctx, generateRemediationProposalFullMethod, req, resp, grpc.ForceCodec(aiRemediationProtoCodec{})); err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *gatewayService) generateRemediationProposalDatasourceAIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &generateRemediationProposalRequest{}
	if err := bind(req, r); err != nil {
		appLogger.Warnf(ctx, "Failed to bind request, req=%s, err=%+v", "GenerateRemediationProposalRequest", err)
	}
	if err := req.Validate(); err != nil {
		writeResponse(ctx, w, http.StatusBadRequest, map[string]interface{}{errorJSONKey: err.Error()})
		return
	}
	resp, err := g.aiRemediationClient.GenerateRemediationProposal(ctx, req)
	if err != nil {
		if handleErr := handleGRPCError(ctx, w, err); handleErr != nil {
			appLogger.Errorf(ctx, "HandleGRPCError: %+v", handleErr)
			writeResponse(ctx, w, http.StatusInternalServerError, map[string]interface{}{errorJSONKey: "InternalServerError"})
		}
		return
	}
	writeResponse(ctx, w, http.StatusOK, map[string]interface{}{successJSONKey: resp})
}

func (r *generateRemediationProposalRequest) Validate() error {
	if r.ProjectID == 0 {
		return errors.New("project_id is required")
	}
	if r.FindingID == 0 {
		return errors.New("finding_id is required")
	}
	return nil
}

func (r *generateRemediationProposalRequest) marshalProto() []byte {
	b := protowire.AppendTag(nil, 1, protowire.VarintType)
	b = protowire.AppendVarint(b, uint64(r.ProjectID))
	b = protowire.AppendTag(b, 2, protowire.VarintType)
	b = protowire.AppendVarint(b, r.FindingID)
	return b
}

func (r *generateRemediationProposalResponse) unmarshalProto(b []byte) error {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
		if num != 1 {
			n = protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			b = b[n:]
			continue
		}
		if typ != protowire.VarintType {
			return fmt.Errorf("unexpected remediation_proposal_id wire type: %d", typ)
		}
		v, n := protowire.ConsumeVarint(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		r.RemediationProposalID = uint32(v)
		b = b[n:]
	}
	return nil
}

type aiRemediationProtoCodec struct{}

func (aiRemediationProtoCodec) Name() string {
	return "proto"
}

func (aiRemediationProtoCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case *generateRemediationProposalRequest:
		return m.marshalProto(), nil
	default:
		return nil, fmt.Errorf("unsupported ai remediation marshal type: %T", v)
	}
}

func (aiRemediationProtoCodec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case *generateRemediationProposalResponse:
		return m.unmarshalProto(data)
	default:
		return fmt.Errorf("unsupported ai remediation unmarshal type: %T", v)
	}
}
