data "anthropic_token_usage_over_time" "daily_by_model" {
  starting_at  = formatdate("YYYY-MM-DD'T'00:00:00'Z'", timeadd(timestamp(), "-168h"))
  bucket_width = "1d"
  group_by     = ["model"]
}
