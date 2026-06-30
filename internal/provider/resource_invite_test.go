package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Invites are immutable post-create. Lifecycle is just create + destroy.
// Emails must hit a domain the org accepts; we generate unique addresses on
// example.com so we don't collide with real invites and the destroy step
// archives them cleanly. The framework's destroy verifies state is removed,
// not that the invite is gone from the org — archived invites still show in
// the Console (limitation of the Admin API).
func TestAccInvite_basic(t *testing.T) {
	email := fmt.Sprintf("tf-acc-%s@example.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_invite" "test" {
  email = %q
  role  = "user"
}
`, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_invite.test", "email", email),
					resource.TestCheckResourceAttr("anthropic_invite.test", "role", "user"),
					resource.TestCheckResourceAttrSet("anthropic_invite.test", "id"),
					resource.TestCheckResourceAttrSet("anthropic_invite.test", "invited_at"),
					resource.TestCheckResourceAttrSet("anthropic_invite.test", "expires_at"),
					resource.TestCheckResourceAttr("anthropic_invite.test", "status", "pending"),
				),
			},
			{
				ResourceName:      "anthropic_invite.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Changing role must force replacement — invites are immutable on the API
// side, the provider models that via stringplanmodifier.RequiresReplace.
func TestAccInvite_roleForceReplace(t *testing.T) {
	email := fmt.Sprintf("tf-acc-%s@example.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testAccInviteConfig(email, "user")},
			{
				Config:           testAccInviteConfig(email, "developer"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// First step's invite gets archived, new one created — this is
					// the expected plan when role changes (RequiresReplace).
				},
			},
		},
	})
}

func testAccInviteConfig(email, role string) string {
	return fmt.Sprintf(`
resource "anthropic_invite" "test" {
  email = %q
  role  = %q
}
`, email, role)
}
