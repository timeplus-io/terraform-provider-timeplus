---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_source Data Source - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus sources run queries in background and send query results to the target system continuously.
---

# timeplus_source (Data Source)

Timeplus sources run queries in background and send query results to the target system continuously.

## Example Usage

```terraform
data "timeplus_source" "example" {
  id = "f2fc567c-2c9c-44e0-a648-ba2c20938e86"
}

output "example_source" {
  // properties are sensitive
  value = {
    name        = data.timeplus_source.example.name
    description = data.timeplus_source.example.description
    stream      = data.timeplus_source.example.stream
    type        = data.timeplus_source.example.type
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The source immutable ID, generated by Timeplus

### Read-Only

- `description` (String) A detailed text describes the source
- `name` (String) The human-friendly name for the source
- `properties` (String, Sensitive) A JSON object defines the configurations for the specific source type. The properites could contain sensitive information like password, secret, etc.
- `stream` (String) The target stream the source ingests data to
- `type` (String) The type of the source, refer to the Timeplus document for supported source types