---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_dashboard Data Source - terraform-provider-timeplus"
subcategory: ""
description: |-
  A dashboard is a set of one or more panels organized and arranged in one web page. A variety of panels are supported to make it easy to construct the visualization components so that you can create the dashboards for specific monitoring and analytics needs.
---

# timeplus_dashboard (Data Source)

A dashboard is a set of one or more panels organized and arranged in one web page. A variety of panels are supported to make it easy to construct the visualization components so that you can create the dashboards for specific monitoring and analytics needs.

## Example Usage

```terraform
data "timeplus_dashboard" "example" {
  id = "dcd5b165-fe72-42f0-a099-aa086bf701b9"
}

output "example_dashboard" {
  value = data.timeplus_dashboard.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The dashboard immutable ID, generated by Timeplus

### Read-Only

- `description` (String) A detailed text describes the dashboard
- `name` (String) The human-friendly name for the dashboard
- `panels` (String) A list of panels defined in a JSON array. The best way to generate such array is to copy it directly from the Timeplus console UI.
