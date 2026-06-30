package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Spend limits are scope.type=user only via the API. Needs a real user_id
// in the org. The destroy step is idempotent — deleting a per-user override
// reverts the user to inherited limits, which is the desired clean-up.
func TestAccSpendLimit_basic(t *testing.T) {
	userID := testAccRequireEnv(t, "ANTHROPIC_TEST_USER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_spend_limit" "test" {
  user_id = %q
  amount  = "1000" # $10.00 USD
  period  = "monthly"
}
`, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_spend_limit.test", "user_id", userID),
					resource.TestCheckResourceAttr("anthropic_spend_limit.test", "amount", "1000"),
					resource.TestCheckResourceAttr("anthropic_spend_limit.test", "period", "monthly"),
					resource.TestCheckResourceAttr("anthropic_spend_limit.test", "scope_type", "user"),
					resource.TestCheckResourceAttrSet("anthropic_spend_limit.test", "id"),
					resource.TestCheckResourceAttrSet("anthropic_spend_limit.test", "currency"),
				),
			},
			// Amount update — re-Set (upsert by scope+period).
			{
				Config: fmt.Sprintf(`
resource "anthropic_spend_limit" "test" {
  user_id = %q
  amount  = "2000"
  period  = "monthly"
}
`, userID),
				Check: resource.TestCheckResourceAttr("anthropic_spend_limit.test", "amount", "2000"),
			},
		},
	})
}

// Effective spend limits is read-only; tests parsing + auth.
func TestAccDataSourceEffectiveSpendLimits_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "anthropic_effective_spend_limits" "all" {}`,
				Check:  resource.TestCheckResourceAttrSet("data.anthropic_effective_spend_limits.all", "summaries.#"),
			},
		},
	})
}
