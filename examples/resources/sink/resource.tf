resource "timeplus_sink" "example" {
  name        = "example"
  description = "hello sink example"
  query       = "select * from foo"
  type        = "http"
  properties = jsonencode({
    url           = "http://localhost:4567"
    http_method   = "POST"
    content_type  = "application/json"
    payload_field = "{{ . | toJSON }}"
  })
}

output "example" {
  value = resource.timeplus_sink.example
}
