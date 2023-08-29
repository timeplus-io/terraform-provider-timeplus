data "timeplus_stream" "example" {
  name = "basic_example"
}

output "example_stream" {
  value = data.timeplus_stream.example
}
