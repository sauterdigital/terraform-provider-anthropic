package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Smallest possible smoke test — proves auth, headers, and the parsing
// path for the simplest endpoint all work end-to-end. If anything in the
// client core breaks, this is the cheapest test to fail first.
func TestAccDataSourceOrganization_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "anthropic_organization" "current" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.anthropic_organization.current", "id"),
					resource.TestCheckResourceAttrSet("data.anthropic_organization.current", "name"),
				),
			},
		},
	})
}

func TestAccDataSourceWorkspaces_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "anthropic_workspaces" "all" { include_archived = false }`,
				// `workspaces` may be empty in a fresh org — we just verify the
				// list attribute exists (returns 0+ items without error).
				Check: resource.TestCheckResourceAttrSet("data.anthropic_workspaces.all", "workspaces.#"),
			},
		},
	})
}

func TestAccDataSourceOrganizationRateLimits_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "anthropic_organization_rate_limits" "all" {}`,
				Check:  resource.TestCheckResourceAttrSet("data.anthropic_organization_rate_limits.all", "groups.#"),
			},
		},
	})
}
