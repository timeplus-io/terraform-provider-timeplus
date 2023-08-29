---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_sink Resource - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus sinks run queries in background and send query results to the target system continuously.
---

# timeplus_sink (Resource)

Timeplus sinks run queries in background and send query results to the target system continuously.

## Example Usage

```terraform
resource "timeplus_stream" "example" {
  name        = "city_temperature"
  description = "An example stream for creating a sink."

  column {
    name = "city_name"
    type = "string"
  }

  column {
    name = "temp"
    type = "float32"
  }
}

resource "timeplus_sink" "example" {
  name        = "Hot Cities"
  description = "An example sink sends data to an HTTP endpoint."
  query       = "select _tp_time, city_name, temp from ${timeplus_stream.example.name}"
  type        = "http"
  properties = jsonencode({
    url           = "http://localhost:6789"
    http_method   = "POST"
    content_type  = "application/json"
    payload_field = "City of {{ .city_name }} recorded a high temperature {{ .temp }} at {{ ._tp_time }}"
  })
}

output "example_sink_id" {
  value = timeplus_sink.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The human-friendly name for the sink
- `properties` (String, Sensitive) A JSON object defines the configurations for the specific sink type. The properites could contain sensitive information like password, secret, etc.
- `query` (String) The query the sink uses to generate data
- `type` (String) The type of the sink, refer to the Timeplus document for supported sink types

### Optional

- `description` (String) A detailed text describes the sink

### Read-Only

- `id` (String) The sink immutable ID, generated by Timeplus