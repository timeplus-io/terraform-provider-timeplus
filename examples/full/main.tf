terraform {
  required_providers {
    timeplus = {
      source = "timeplus-io/timeplus"
    }
  }
}

provider "timeplus" {
  username = "proton"
  password = "proton@t+"
}

data "timeplus_stream" "bar" {
  name = "bar"
}

output "stream_bar" {
  value = data.timeplus_stream.bar
}

resource "timeplus_stream" "foo" {
  name = "foo"

  description = "a stream managed with Terraform timeplus provider"

  column {
    name    = "col_1"
    type    = "string"
    default = "foo"
  }

  column {
    name = "col_2"
    type = "int32"
  }

  column {
    name = "timestamp"
    type = "datetime64(3)"
  }
}

output "foo_stream" {
  value = resource.timeplus_stream.foo
}
