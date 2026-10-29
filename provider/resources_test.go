package cockroachsql

import (
	"context"
	"os"
	"testing"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/nellisauction/pulumi-cockroachsql/provider/pkg/version"
)

func TestReplacementDeletesBeforeCreating(t *testing.T) {
	schema, err := os.ReadFile("cmd/pulumi-resource-cockroachsql/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if version.Version == "" {
		version.Version = "0.0.0-test"
	}
	info := Provider()
	server := tfbridge.NewProvider(context.Background(), nil, "cockroachsql", "0.0.0", info.P, info, schema)

	tests := []struct {
		name string
		typ  string
		olds resource.PropertyMap
		news resource.PropertyMap
	}{
		{
			name: "GrantRole with_admin_option change",
			typ:  "cockroachsql:index/grantRole:GrantRole",
			olds: resource.NewPropertyMapFromMap(map[string]any{
				"id":              "app_user_readers",
				"role":            "app_user",
				"grantRole":       "readers",
				"withAdminOption": false,
			}),
			news: resource.NewPropertyMapFromMap(map[string]any{
				"role":            "app_user",
				"grantRole":       "readers",
				"withAdminOption": true,
			}),
		},
		{
			name: "DefaultPrivileges with_grant_option change",
			typ:  "cockroachsql:index/defaultPrivileges:DefaultPrivileges",
			olds: resource.NewPropertyMapFromMap(map[string]any{
				"id":              "app_user_appdb_public_table",
				"role":            "app_user",
				"database":        "appdb",
				"schema":          "public",
				"objectType":      "table",
				"privileges":      []any{"SELECT"},
				"forAllRoles":     false,
				"withGrantOption": false,
			}),
			news: resource.NewPropertyMapFromMap(map[string]any{
				"role":            "app_user",
				"database":        "appdb",
				"schema":          "public",
				"objectType":      "table",
				"privileges":      []any{"SELECT"},
				"forAllRoles":     false,
				"withGrantOption": true,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			olds, err := plugin.MarshalProperties(tt.olds, plugin.MarshalOptions{KeepUnknowns: true})
			if err != nil {
				t.Fatal(err)
			}
			news, err := plugin.MarshalProperties(tt.news, plugin.MarshalOptions{KeepUnknowns: true})
			if err != nil {
				t.Fatal(err)
			}
			urn := resource.NewURN("stack", "project", "", tokens.Type(tt.typ), "test")
			resp, err := server.Diff(context.Background(), &pulumirpc.DiffRequest{
				Id:   tt.olds["id"].StringValue(),
				Urn:  string(urn),
				Olds: olds,
				News: news,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.GetReplaces()) == 0 {
				t.Fatalf("expected a replacement, got diff %v", resp)
			}
			if !resp.GetDeleteBeforeReplace() {
				t.Errorf("expected DeleteBeforeReplace for replaces %v", resp.GetReplaces())
			}
		})
	}
}
