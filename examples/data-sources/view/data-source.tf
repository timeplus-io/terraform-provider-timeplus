data "timeplus_view" "example" {
  name = "example"
}

output "example_view_data" {
  value = data.timeplus_view.example
}
