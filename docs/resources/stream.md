---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_stream Resource - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus streams are similar to tables in the traditional SQL databases. Both of them are essentially datasets. The key difference is that Timeplus stream is an append-only (by default), unbounded, constantly changing events group.
---

# timeplus_stream (Resource)

Timeplus streams are similar to tables in the traditional SQL databases. Both of them are essentially datasets. The key difference is that Timeplus stream is an append-only (by default), unbounded, constantly changing events group.

## Example Usage

```terraform
resource "timeplus_stream" "basic_example" {
  name = "basic_example"

  description = "A simple stream with three columns"

  column {
    name = "col_1"
    type = "string"
  }

  column {
    name = "col_2"
    type = "int32"
  }

  column {
    name    = "col_3"
    type    = "datetime64(3)"
    default = "now()"
  }
}

resource "timeplus_stream" "codec_example" {
  name = "codec_example"

  description = "An example shows how to use codec on columns"

  column {
    name  = "col_1"
    type  = "string"
    codec = "LZ4"
  }

  column {
    name = "col_2"
    type = "int32"
  }

  column {
    name    = "col_3"
    type    = "datetime64(3)"
    default = "now()"
  }
}

resource "timeplus_stream" "retention_example" {
  name = "retention_example"

  description = "An example shows how to customize retention plicy on a stream"

  column {
    name = "col_1"
    type = "string"
  }

  column {
    name = "col_2"
    type = "int32"
  }

  column {
    name    = "col_3"
    type    = "datetime64(3)"
    default = "now()"
  }

  retention_bytes = 10 * 1024 * 1204 * 1024 // 10Gi in bytes
  retention_ms    = 7 * 24 * 60 * 60 * 1000 // 7 days in ms 
  history_ttl     = "to_datetime(_tp_time) + INTERVAL 30 DAY"
}

resource "timeplus_stream" "mode_example" {
  name = "mode_example"

  description = "An example shows how to use different mode to create stream"

  mode = "versioned_kv"

  column {
    name        = "id"
    type        = "string"
    primary_key = true
  }

  column {
    name = "value"
    type = "int32"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The stream name

### Optional

- `column` (Block List) Define the columns of the stream (see [below for nested schema](#nestedblock--column))
- `description` (String) A detailed text describes the stream
- `history_ttl` (String) A SQL expression defines the maximum age of data that are persisted in the historical store
- `mode` (String) The stream mode. Options: append, changelog, changelog_kv, versioned_kv. Default: "append"
- `retention_bytes` (Number) The retention size threadhold in bytes indicates how many data could be kept in the streaming store
- `retention_ms` (Number) The retention period threadhold in millisecond indicates how long data could be kept in the streaming store

<a id="nestedblock--column"></a>
### Nested Schema for `column`

Required:

- `name` (String) The column name
- `type` (String) The type name of the column

Optional:

- `codec` (String) The codec for value encoding
- `default` (String) The default value for the column
- `primary_key` (Boolean) If set to `true`, this column will be used as the primary key, or part of the combined primary key if multiple columns are marked as primary keys.
- `use_as_event_time` (Boolean) If set to `true`, this column will be used as the event time column (by default ingest time will be used as event time). Only one column can be marked as the event time column in a stream.
