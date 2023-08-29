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
