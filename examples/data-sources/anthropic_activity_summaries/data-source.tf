data "anthropic_activity_summaries" "last_week" {
  starting_date = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-240h")) # 10 days ago
  ending_date   = formatdate("YYYY-MM-DD", timeadd(timestamp(), "-72h"))  # 3 days ago
}
