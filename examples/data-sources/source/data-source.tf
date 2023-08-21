data "timeplus_source" "example" {
  id = "f50750ca-8535-42c5-9622-baba08df55d9"
}

output "existing" {
  value = data.timeplus_source.example
}
