resource "timeplus_remote_function" "example" {
  name        = "add"
  description = "a Timeplus remote function example"
  url         = "http://localhost:9090"

  auth_header = {
    name  = "Authorization"
    value = "Token ABCDEFG1234567890"
  }

  return_type = "int64"

  arg {
    name = "a"
    type = "int64"
  }

  arg {
    name = "b"
    type = "int64"
  }
}

output "example_remote_function" {
  value = resource.timeplus_remote_function.example
}
