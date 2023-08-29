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
