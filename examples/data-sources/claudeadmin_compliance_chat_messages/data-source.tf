# CAUTION: message TEXT is materialized in Terraform state. Only use this
# for one-off eDiscovery / DLP investigation on a state file you handle
# with the same care as the messages themselves.

data "claudeadmin_compliance_chat_messages" "under_investigation" {
  chat_id = "chat_01ABC..."
}

output "message_texts" {
  value     = [for m in data.claudeadmin_compliance_chat_messages.under_investigation.messages : m.text]
  sensitive = true
}
