package main

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingPackageQueryBinding(t *testing.T) {
	cases := []struct {
		name string
		org  bool
		view bool
	}{
		{"project categories", false, true},
		{"organization categories", true, true},
		{"project findings", false, false},
		{"organization findings", true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query := url.Values{"project_id": {"1001"}, "organization_id": {"1009"}, "asset_type": {"Image"}, "asset_key": {"Lxxp4l0xeIUsX2i4WZFkXUeuVWGiKO9s-7WA9FDl6U0"}, "asset_project_id": {"1001"}, "limit": {"20"}}
			key := ""
			if c.view {
				query.Set("package_view", "true")
			} else {
				key = `["NPM", "@example/library", "", "app:/app"]`
				query.Set("package_key", key)
			}
			r := httptest.NewRequest("GET", "/?"+query.Encode(), nil)
			if c.org {
				req := &finding.ListFindingForOrgRequest{}
				require.NoError(t, bindQuery(req, r))
				require.Equal(t, c.view, req.PackageView)
				require.Equal(t, key, req.PackageKey)
				require.NoError(t, req.Validate())
			} else {
				req := &finding.ListFindingRequest{}
				require.NoError(t, bindQuery(req, r))
				require.Equal(t, c.view, req.PackageView)
				require.Equal(t, key, req.PackageKey)
				require.NoError(t, req.Validate())
			}
		})
	}
}

func TestFindingPackageResponseJSON(t *testing.T) {
	packages := []*finding.FindingPackage{{PackageKey: `["NPM", "next", "", "app:/app"]`, PackageName: "next", PackageType: "NPM", WorkType: "app", WorkLocation: "/app", FindingCount: 2, Versions: []*finding.FindingPackageVersion{{CurrentVersion: "1", FixedVersion: "2", FindingCount: 2}}}}
	cases := []struct {
		name string
		body any
	}{
		{"project", &finding.ListFindingResponse{Packages: packages, Total: 1, TotalFindingCount: 2}},
		{"organization", &finding.ListFindingForOrgResponse{Packages: packages, Total: 1, TotalFindingCount: 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := json.Marshal(c.body)
			require.NoError(t, err)
			var decoded struct {
				Packages          []*finding.FindingPackage `json:"packages"`
				Total             uint32                    `json:"total"`
				TotalFindingCount uint64                    `json:"total_finding_count"`
			}
			require.NoError(t, json.Unmarshal(data, &decoded))
			require.Equal(t, uint32(1), decoded.Total)
			require.Equal(t, uint64(2), decoded.TotalFindingCount)
			require.Equal(t, packages, decoded.Packages)
		})
	}
}
