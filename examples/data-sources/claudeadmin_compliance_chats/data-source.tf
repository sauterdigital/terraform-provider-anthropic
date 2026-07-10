# Enterprise + Compliance Access Key with scope read:compliance_user_data.
# CAUTION: chat metadata + counts land in Terraform state.

data "claudeadmin_compliance_chats" "user_alice_last_month" {
  user_id     = "user_01ABC..."
  starting_at = "2026-06-01T00:00:00Z"
  ending_at   = "2026-07-01T00:00:00Z"
}

output "alice_recent_chats" {
  value = [
    for c in data.claudeadmin_compliance_chats.user_alice_last_month.chats : {
      id       = c.id
      title    = c.title
      messages = c.message_count
    }
  ]
}
