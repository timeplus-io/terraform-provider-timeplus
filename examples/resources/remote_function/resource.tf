resource "timeplus_remote_function" "simple_example" {
  name        = "add"
  description = "a Timeplus remote function example that accepts two integers and returns one"
  url         = "https://some.domain/that/hosts/my/function"

  return_type = "int64"

  arg {
    name = "left"
    type = "int64"
  }

  arg {
    name = "right"
    type = "int64"
  }
}

resource "timeplus_remote_function" "http_header_example" {
  name        = "add_with_header"
  description = "a Timeplus remote function example that uses HTTP header and accepts two integers and returns one"
  url         = "https://some.domain/that/hosts/my/function"

  auth_header = {
    name  = "Authorization"
    value = "Token my-secret-token"
  }

  return_type = "int64"

  arg {
    name = "left"
    type = "int64"
  }

  arg {
    name = "right"
    type = "int64"
  }
}
