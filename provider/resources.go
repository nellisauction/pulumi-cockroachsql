package cockroachsql

import (
	"fmt"
	"path/filepath"
	"unicode"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	cockroachsql "github.com/nellisauction/terraform-provider-cockroachsql/cockroachsql" // Import the upstream provider

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	tfbridgetokens "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfgen"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"

	"github.com/nellisauction/pulumi-cockroachsql/provider/pkg/version"
)

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	// registries for nodejs and python:
	mainPkg = "cockroachsql"
	// modules:
	mainMod = "index" // the cockroachsql module
)

// makeMember manufactures a type token for the package and the given module and type.
func makeMember(mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(mainPkg + ":" + mod + ":" + mem)
}

// makeType manufactures a type token for the package and the given module and type.
func makeType(mod string, typ string) tokens.Type {
	return tokens.Type(makeMember(mod, typ))
}

// makeResource manufactures a standard resource token given a module and resource name.  It
// automatically uses the main package and names the file by simply lower casing the resource's
// first character.
func makeResource(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return makeType(mod+"/"+fn, res)
}

//go:embed cmd/pulumi-resource-cockroachsql/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the provider.
func Provider() tfbridge.ProviderInfo {
	// Instantiate the Terraform provider
	p := shimv2.NewProvider(cockroachsql.Provider())

	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		P:           p,
		Name:        "cockroachsql",
		Version:     version.Version,
		DisplayName: "CockroachSQL",
		Publisher:         "NellisAuction",
		LogoURL:           "",
		PluginDownloadURL: "github://api.github.com/nellisauction/pulumi-cockroachsql",
		Description:       "A Pulumi package for creating and managing CockroachSQL resources.",
		Keywords:    []string{"pulumi", "cockroachsql", "category/database"},
		License:     "Apache-2.0",
		Homepage:    "https://github.com/nellisauction/pulumi-cockroachsql",
		Repository:  "https://github.com/nellisauction/pulumi-cockroachsql",
		GitHubOrg:   "nellisauction",
		DocRules:    &tfbridge.DocRuleInfo{EditRules: docEditRules},
		Config: map[string]*tfbridge.SchemaInfo{
			"sslmode": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"PGSSLMODE"},
				},
			},
			"connect_timeout": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"PGCONNECT_TIMEOUT"},
					Value:   180,
				},
			},
		},
		Resources: map[string]*tfbridge.ResourceInfo{
			"cockroachsql_grant_role": {
				Tok: makeResource(mainMod, "GrantRole"),
				Fields: map[string]*tfbridge.SchemaInfo{
					"grant_role": {
						CSharpName: "GrantRoleName",
					},
				},
			},
			"cockroachsql_grant": {
				Tok:                 makeResource(mainMod, "Grant"),
				DeleteBeforeReplace: true,
			},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName: "@nellisauction/pulumi-cockroachsql",
			RespectSchemaVersion: true,
		},
		Python: (func() *tfbridge.PythonInfo {
			i := &tfbridge.PythonInfo{
				RespectSchemaVersion: true,
			}
			i.PyProject.Enabled = true
			return i
		})(),
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/nellisauction/pulumi-%[1]s/sdk/", mainPkg),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RespectSchemaVersion: true,
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
			Namespaces: map[string]string{
				mainPkg: "CockroachSql",
			},
		},
		MetadataInfo:                   tfbridge.NewProviderMetadata(metadata),
		EnableZeroDefaultSchemaVersion: true,
		EnableAccurateBridgePreview:    true,
	}

	prov.MustComputeTokens(tfbridgetokens.SingleModule("cockroachsql_", mainMod,
		tfbridgetokens.MakeStandard(mainPkg)))

	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}

func docEditRules(defaults []tfbridge.DocsEdit) []tfbridge.DocsEdit {
	return append(
		defaults,
		skipInstallationSections...,
	)
}

var skipInstallationSections = []tfbridge.DocsEdit{
	// TF Variable do not apply to Pulumi
	{
		Path: "index.html.markdown",
		Edit: func(_ string, content []byte) ([]byte, error) {
			return tfgen.SkipSectionByHeaderContent(content, func(headerText string) bool {
				return headerText == "Terraform Variables"
			})
		},
	},
	// This section had TF specific instructions as well
	{
		Path: "index.html.markdown",
		Edit: func(_ string, content []byte) ([]byte, error) {
			return tfgen.SkipSectionByHeaderContent(content, func(headerText string) bool {
				return headerText == "Data Sources and Resources"
			})
		},
	},
}
