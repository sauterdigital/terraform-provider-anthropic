data "claudeadmin_skills_usage" "last_week" {
  starting_date = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-168h"))
}
