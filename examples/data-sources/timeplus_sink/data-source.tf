data "timeplus_sink" "example" {
  id = "7f10df2d-6ad4-4dea-9954-3d8934b2c329"
}

output "example_sink" {
  // properties is sensitive
  value = {
    name        = data.timeplus_sink.example.name
    description = data.timeplus_sink.example.description
    type        = data.timeplus_sink.example.type
    query       = data.timeplus_sink.example.query
  }
}
