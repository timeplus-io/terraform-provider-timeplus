resource "timeplus_stream" "example" {
  name = "example"

  description = "a stream managed with Terraform timeplus provider"

  column {
    name    = "col_1"
    type    = "string"
    default = "foo"
    codec   = "LZ4"
  }

  column {
    name = "col_2"
    type = "int32"
  }

  column {
    name = "timestamp"
    type = "datetime64(3)"
  }
  historical_data_ttl = "to_datetime(_tp_time)     +    INTERVAL 14 DAY"
}
