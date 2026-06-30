# Every member's effective spend limit + period-to-date spend.
data "anthropic_effective_spend_limits" "monthly" {
  period = ["monthly"]
}
