data "anthropic_user_activity" "last_month" {
  starting_date = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-720h"))
}
