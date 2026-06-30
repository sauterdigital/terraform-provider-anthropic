data "anthropic_chat_projects_usage" "last_week" {
  starting_date = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-168h"))
}
