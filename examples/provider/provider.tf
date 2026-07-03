terraform {
  required_providers {
    claudeadmin = {
      source  = "sauterdigital/claudeadmin"
      version = "~> 0.4"
    }
  }
}

provider "claudeadmin" {
  # admin_api_key      = "sk-ant-admin-..." # or ANTHROPIC_ADMIN_API_KEY
  # oauth_token        = "..."              # or ANTHROPIC_OAUTH_TOKEN (Service Accounts, Federation, MCP Tunnels)
  # compliance_api_key = "sk-ant-api01-..." # or ANTHROPIC_COMPLIANCE_API_KEY (Compliance data sources)
}
