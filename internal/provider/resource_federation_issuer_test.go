package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Uses GitHub Actions' public OIDC issuer URL as a test target — it's
// always reachable, has stable discovery metadata, and creating a
// Federation Issuer pointing at it is harmless (no rules attached, so no
// trust granted). The destroy step archives the issuer.
func TestAccFederationIssuer_discovery(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fi")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOAuth(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_federation_issuer" "gh" {
  name       = %q
  issuer_url = "https://token.actions.githubusercontent.com"
  jwks_type  = "discovery"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_federation_issuer.gh", "name", name),
					resource.TestCheckResourceAttr("anthropic_federation_issuer.gh", "issuer_url", "https://token.actions.githubusercontent.com"),
					resource.TestCheckResourceAttr("anthropic_federation_issuer.gh", "jwks_type", "discovery"),
					resource.TestCheckResourceAttrSet("anthropic_federation_issuer.gh", "id"),
				),
			},
			{
				ResourceName:      "anthropic_federation_issuer.gh",
				ImportState:       true,
				ImportStateVerify: true,
				// max_jwt_lifetime_seconds + check_jti get filled by API defaults;
				// import won't see them if the user didn't set them.
				ImportStateVerifyIgnore: []string{"max_jwt_lifetime_seconds", "check_jti"},
			},
		},
	})
}

// Exercises the inline JWKS path with a synthetic empty key list to prove
// serialization round-trips. The API will accept an inline JWKS with empty
// keys for issuer creation (it only fails token verification, not setup).
func TestAccFederationIssuer_inlineJWKS(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-fi-inline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOAuth(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "anthropic_federation_issuer" "inline" {
  name           = %q
  issuer_url     = "https://example.invalid/issuer"
  jwks_type      = "inline"
  jwks_keys_json = jsonencode([])
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("anthropic_federation_issuer.inline", "jwks_type", "inline"),
					resource.TestCheckResourceAttrSet("anthropic_federation_issuer.inline", "id"),
				),
			},
		},
	})
}
