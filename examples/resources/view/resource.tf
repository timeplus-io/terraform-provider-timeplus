resource "timeplus_view" "example" {
  name        = "example"
  description = "this is a view example"
  query       = "select * from foo where name like 'example_%'"
}

output "example" {
  value = resource.timeplus_view.example
}
