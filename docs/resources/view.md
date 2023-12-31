---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_view Resource - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus views are named queries. When you create a view, you basically create a query and assign a name to the query. Therefore, a view is useful for wrapping a commonly used complex query.
---

# timeplus_view (Resource)

Timeplus views are named queries. When you create a view, you basically create a query and assign a name to the query. Therefore, a view is useful for wrapping a commonly used complex query.

## Example Usage

```terraform
resource "timeplus_stream" "traffic" {
  name        = "traffic"
  description = "A stream for view example"

  column {
    name = "license_no"
    type = "string"
  }

  column {
    name = "road_speed_limit_mph"
    type = "uint8"
  }

  column {
    name = "speed_mph"
    type = "uint8"
  }
}

resource "timeplus_view" "example" {
  name  = "speeding_vehcles"
  query = "select * from ${resource.timeplus_stream.traffic.name} where speed_mph > road_speed_limit_mph"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The view name
- `query` (String) The query SQL of the view

### Optional

- `description` (String) A detailed text describes the view
