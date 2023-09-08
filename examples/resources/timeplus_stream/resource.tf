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
