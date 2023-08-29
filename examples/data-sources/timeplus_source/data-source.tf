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
