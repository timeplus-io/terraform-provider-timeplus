resource "timeplus_stream" "traffic" {
  name        = "traffic"
  description = "A stream for materialized view example"

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

resource "timeplus_materialized_view" "basic_example" {
  name        = "speeding_vehcles"
  description = "A basic example with default retention policy"
  query       = <<-SQL
  select * from ${resource.timeplus_stream.traffic.name} where speed_mph > road_speed_limit_mph
  SQL
}

resource "timeplus_materialized_view" "retention_example" {
  name            = "speeding_vehcles_with_retention"
  description     = "A materialized view with custom retention policy"
  query           = <<-SQL
  select * from ${resource.timeplus_stream.traffic.name} where speed_mph > road_speed_limit_mph
  SQL
  retention_bytes = 10 * 1024 * 1204 * 1024 // 10Gi in bytes
  retention_ms    = 7 * 24 * 60 * 60 * 1000 // 7 days in ms 
  history_ttl     = "to_datetime(_tp_time) + INTERVAL 30 DAY"
}
