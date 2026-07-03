data "claudeadmin_connectors_usage" "last_week" {
  starting_date = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-168h"))
}
