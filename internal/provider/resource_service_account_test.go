package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Requires OAuth Bearer (admin keys are rejected by the API). The
// pre-check skips if ANTHROPIC_OAUTH_TOKEN is unset.
func TestAccServiceAccount_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-sa")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOAuth(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_service_account" "test" {
  name              = %q
  description       = "acceptance test SA — safe to destroy"
  organization_role = "developer"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_service_account.test", "name", name),
					resource.TestCheckResourceAttr("anthropic_service_account.test", "organization_role", "developer"),
					resource.TestCheckResourceAttrSet("anthropic_service_account.test", "id"),
					resource.TestCheckResourceAttrSet("anthropic_service_account.test", "created_at"),
				),
			},
			// Update description in place (no replace).
			{
				Config: fmt.Sprintf(`
resource "anthropic_service_account" "test" {
  name              = %q
  description       = "updated description"
  organization_role = "developer"
}
`, name),
				Check: resource.TestCheckResourceAttr("anthropic_service_account.test", "description", "updated description"),
			},
			{
				ResourceName:      "anthropic_service_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
