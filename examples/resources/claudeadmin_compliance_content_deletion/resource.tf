# eDiscovery / DLP hard-delete. IRREVERSIBLE. Requires Compliance Access
# Key with scope delete:compliance_user_data.
#
# Common target_types:
#   chat, chat_file, chat_generated_file, project, project_document

# LGPD art. 18 erasure — delete an entire chat
resource "claudeadmin_compliance_content_deletion" "lgpd_ticket_4471" {
  target_type = "chat"
  target_id   = "chat_01ABC..."
  reason      = "LGPD art. 18 erasure request — ticket #4471"
}

# DLP incident — delete a leaked file from a project
resource "claudeadmin_compliance_content_deletion" "dlp_incident_2026_07_10" {
  target_type = "project_document"
  target_id   = "doc_01XYZ..."
  reason      = "DLP incident 2026-07-10 — customer secrets accidentally uploaded"
}
