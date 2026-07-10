data "claudeadmin_compliance_projects" "all" {}

output "top_users_by_projects" {
  value = {
    for p in data.claudeadmin_compliance_projects.all.projects :
    p.user_id => p.chat_count...
  }
}
