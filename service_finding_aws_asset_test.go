package main

import (
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestAWSAssetExpansionValidation(t *testing.T) {
	cases := []struct {
		name, typ string
		org       bool
	}{
		{"EC2", "EC2Instance", false}, {"user", "IAMUser", false}, {"role org", "IAMRole", true}, {"bucket org", "S3Bucket", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			scope := "project_id=1001"
			if c.org {
				scope = "organization_id=1009"
			}
			r := httptest.NewRequest("GET", "/?"+scope+"&asset_type="+c.typ+"&asset_key=Lxxp4l0xeIUsX2i4WZFkXUeuVWGiKO9s-7WA9FDl6U0&asset_project_id=1001&package_view=true", nil)
			if c.org {
				req := &finding.ListFindingForOrgRequest{}
				require.NoError(t, bindQuery(req, r))
				require.NoError(t, req.Validate())
				require.Equal(t, c.typ, req.AssetType)
			} else {
				req := &finding.ListFindingRequest{}
				require.NoError(t, bindQuery(req, r))
				require.NoError(t, req.Validate())
				require.Equal(t, c.typ, req.AssetType)
			}
		})
	}
}
