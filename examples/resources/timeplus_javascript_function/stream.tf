resource "timeplus_stream" "two_numbers" {
  name = "two_numbers"
  column {
    name = "a"
    type = "int64"
  }
  column {
    name = "b"
    type = "int64"
  }
}
