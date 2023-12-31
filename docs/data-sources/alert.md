---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_alert Data Source - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus alerts run queries in background and send query results to the target system continuously.
---

# timeplus_alert (Data Source)

Timeplus alerts run queries in background and send query results to the target system continuously.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The alert immutable ID, generated by Timeplus

### Read-Only

- `action` (String) The type of action the alert should take, i.e. the name of the target system, like 'slack', 'email', etc. Please refer to the Timeplus document for supported alert action types
- `description` (String) A detailed text describes the alert
- `name` (String) The human-friendly name for the alert
- `properties` (String, Sensitive) a JSON object defines the configurations for the specific alert action. The properites could contain sensitive information like password, secret, etc.
- `resolve_sql` (String) The query the alert uses to generate events that resolve the alert
- `severity` (Number) A number indicates how serious this alert is
- `trigger_sql` (String) The query the alert uses to generate events that trigger the alert
