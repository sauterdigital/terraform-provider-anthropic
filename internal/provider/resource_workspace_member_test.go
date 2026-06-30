package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Needs an existing user_id in the org. Set via ANTHROPIC_TEST_USER_ID env
// var — skipped otherwise. Creates a fresh workspace for the test so the
// destroy step actually removes the membership (not just the user).
func TestAccWorkspaceMember_basic(t *testing.T) {
	userID := testAccRequireEnv(t, "ANTHROPIC_TEST_USER_ID")
	wsName := acctest.RandomWithPrefix("tf-acc-wm")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_workspace" "test" {
  name = %q
}

resource "anthropic_workspace_member" "test" {
  workspace_id   = anthropic_workspace.test.id
  user_id        = %q
  workspace_role = "workspace_developer"
}
`, wsName, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_workspace_member.test", "user_id", userID),
					resource.TestCheckResourceAttr("anthropic_workspace_member.test", "workspace_role", "workspace_developer"),
					resource.TestCheckResourceAttrSet("anthropic_workspace_member.test", "id"),
				),
			},
			// Role update — should NOT force replace, just POST to update endpoint.
			{
				Config: fmt.Sprintf(`
resource "anthropic_workspace" "test" {
  name = %q
}

resource "anthropic_workspace_member" "test" {
  workspace_id   = anthropic_workspace.test.id
  user_id        = %q
  workspace_role = "workspace_admin"
}
`, wsName, userID),
				Check: resource.TestCheckResourceAttr("anthropic_workspace_member.test", "workspace_role", "workspace_admin"),
			},
			{
				ResourceName:      "anthropic_workspace_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
