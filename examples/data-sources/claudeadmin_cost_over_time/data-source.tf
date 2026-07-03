data "claudeadmin_cost_over_time" "monthly_by_product" {
  starting_at  = formatdate("YYYY-MM-DD'T'00:00:00'Z'", timeadd(timestamp(), "-720h"))
  bucket_width = "1d"
  group_by     = ["product"]
}
