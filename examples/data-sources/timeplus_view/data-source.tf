data "timeplus_view" "example" {
  name = "speeding_vehcles"
}

output "example_view" {
  value = data.timeplus_view.example
}
