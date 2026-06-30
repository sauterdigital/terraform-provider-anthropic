# Register a CA certificate for an MCP tunnel. Beta — requires OAuth Bearer.
resource "anthropic_tunnel_certificate" "primary" {
  tunnel_id          = "tunnel_01ABC..."
  ca_certificate_pem = file("${path.module}/certs/internal-ca.pem")
}
