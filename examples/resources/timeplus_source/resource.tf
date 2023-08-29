variable "password" {
  type        = string
  description = "password to connect to the kafka broker"
  sensitive   = true
}

resource "timeplus_stream" "example" {
  name        = "from_kafka"
  description = "A stream for the kafka source example"

  column {
    name = "name"
    type = "string"
  }

  column {
    name = "value"
    type = "float64"
  }
}

resource "timeplus_source" "kafka_example" {
  name        = "kafka example"
  description = "A source example connects to a locally deployed insecure kafka cluster"
  stream      = timeplus_stream.example.name
  type        = "kafka"
  properties = jsonencode({
    brokers = "127.0.0.1:19092"
    topic   = "some-topic"
    offset  = "latest"
    tls = {
      disable = true
    }
    data_type = "json"
  })
}

resource "timeplus_source" "secure_kafka_example" {
  name        = "secure kafka example"
  description = "A source example connects to a locally deployed secure kafka cluster"
  stream      = timeplus_stream.example.name
  type        = "kafka"
  properties = jsonencode({
    brokers = "127.0.0.1:19092"
    topic   = "some-topic"
    offset  = "latest"
    tls = {
      disable = false
    }
    sasl      = "plain"
    username  = "some-username"
    password  = var.password
    data_type = "json"
  })
}
