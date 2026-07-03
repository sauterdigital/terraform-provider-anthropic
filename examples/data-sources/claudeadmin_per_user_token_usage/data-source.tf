data "claudeadmin_per_user_token_usage" "last_week" {
  starting_at  = formatdate("YYYY-MM-DD'T'00:00:00'Z'", timeadd(timestamp(), "-168h"))
  bucket_width = "1d"
}
