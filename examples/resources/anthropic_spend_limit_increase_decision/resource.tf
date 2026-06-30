# Approve a pending spend-limit increase request, setting a new $200/mo cap.
resource "anthropic_spend_limit_increase_decision" "approve_alice" {
  request_id = "slir_01ABC..."
  decision   = "approve"
  amount     = "20000"
  period     = "monthly"
}

# Or deny:
# resource "anthropic_spend_limit_increase_decision" "deny_bob" {
#   request_id = "slir_01XYZ..."
#   decision   = "deny"
# }
