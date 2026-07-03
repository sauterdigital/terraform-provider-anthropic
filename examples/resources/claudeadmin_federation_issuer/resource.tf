# Federation REQUIRES OAuth Bearer auth.

# Common case — GitHub Actions OIDC, using discovery (default).
resource "claudeadmin_federation_issuer" "github_actions" {
  name       = "github-actions"
  issuer_url = "https://token.actions.githubusercontent.com"
  jwks_type  = "discovery"
}

# Alternative — explicit JWKS URL.
# resource "claudeadmin_federation_issuer" "self_hosted" {
#   name      = "internal-idp"
#   issuer_url = "https://idp.internal.example.com"
#   jwks_type  = "explicit_url"
#   jwks_url   = "https://idp.internal.example.com/.well-known/jwks.json"
# }

# Inline JWKS (static key set).
# resource "claudeadmin_federation_issuer" "static" {
#   name       = "static-jwks"
#   issuer_url = "https://example.com/issuer"
#   jwks_type  = "inline"
#   jwks_keys_json = jsonencode([{ kty = "RSA", n = "...", e = "AQAB", kid = "key-1" }])
# }
