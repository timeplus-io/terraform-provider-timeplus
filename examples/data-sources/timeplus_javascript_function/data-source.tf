data "timeplus_javascript_function" "example" {
  name = "add"
}

output "example_js_func" {
  value = data.timeplus_javascript_function.example
}
