data "timeplus_alert" "example" {
  id = "3bd2e35a-54d1-4540-a0cf-70589e2e4b5a"
}

// properties are sensitive
output "example_alert" {
  value = {
    name        = data.timeplus_alert.example.name
    description = data.timeplus_alert.example.description
    severity    = data.timeplus_alert.example.severity
    action      = data.timeplus_alert.example.action
    trigger_sql = data.timeplus_alert.example.trigger_sql
    resolve_sql = data.timeplus_alert.example.resolve_sql
  }
}
