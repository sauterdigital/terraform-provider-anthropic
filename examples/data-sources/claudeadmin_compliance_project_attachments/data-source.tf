data "claudeadmin_compliance_project_attachments" "under_review" {
  project_id = "proj_01ABC..."
}

output "big_attachments" {
  value = [
    for a in data.claudeadmin_compliance_project_attachments.under_review.attachments :
    a if a.size_bytes > 10 * 1024 * 1024
  ]
}
