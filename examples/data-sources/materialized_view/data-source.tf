data "timeplus_materialized_view" "example" {
  name = "example"
}

output "existing" {
  value = data.timeplus_materialized_view.example
}
