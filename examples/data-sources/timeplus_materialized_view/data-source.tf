data "timeplus_materialized_view" "example" {
  name = "speeding_vehcles_with_retention"
}

output "example_mv" {
  value = data.timeplus_materialized_view.example
}
