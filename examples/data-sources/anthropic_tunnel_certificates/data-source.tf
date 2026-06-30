data "anthropic_tunnel_certificates" "primary_certs" {
  tunnel_id        = "tunnel_01ABC..."
  include_archived = false
}
