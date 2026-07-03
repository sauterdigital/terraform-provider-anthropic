resource "claudeadmin_workspace" "example" {
  name = "engineering"

  tags = {
    env  = "prod"
    team = "platform"
  }
}

output "workspace_id" {
  value = claudeadmin_workspace.example.id
}
