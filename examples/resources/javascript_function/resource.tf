resource "timeplus_javascript_function" "example" {
  name        = "add"
  description = "a Timeplus javascript function example"

  return_type = "int64"

  arg {
    name = "a"
    type = "int64"
  }

  arg {
    name = "b"
    type = "int64"
  }

  source = <<-EOS
  function add(first_values, second_values) {
    results = []
    for (const [i, first] of first_values.entries()) {
      results.push(first + second_values[i])
    }
    return results
  }
  EOS
}
