data "timeplus_dashboard" "example" {
  id = "dcd5b165-fe72-42f0-a099-aa086bf701b9"
}

output "example_dashboard" {
  value = data.timeplus_dashboard.example
}
