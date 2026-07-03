resource "claudeadmin_spend_limit" "alice_monthly" {
  user_id = "user_01WCz1FkmYMm4gnmykNKUu3Q"
  amount  = "10000" # $100.00 in minor units
  period  = "monthly"
}
