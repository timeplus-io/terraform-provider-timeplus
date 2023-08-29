data "timeplus_remote_function" "example" {
  name = "add"
}

output "example_remote_func" {
  value = data.timeplus_remote_function.example
}
