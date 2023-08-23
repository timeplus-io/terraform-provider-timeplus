variable "password" {
  type        = string
  description = "password to connect to the kafka broker"
  sensitive   = true
}

resource "timeplus_source" "example" {
  name        = "example"
  description = "hello source example"
  stream      = "foo"
  type        = "kafka"
  properties = jsonencode({
    brokers = "localhost:9092"
    tls = {
      disable            = false
      skip_verify_server = false
    }
    sasl      = "plain"
    username  = "some-username"
    password  = var.password
    data_type = "json"
  })
}
