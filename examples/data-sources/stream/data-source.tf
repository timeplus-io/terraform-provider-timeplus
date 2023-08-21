data "timeplus_stream" "example" {
  name = "example"
}

output "example_stream_data" {
  value = data.timeplus_stream.example
}
