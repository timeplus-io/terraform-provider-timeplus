---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "timeplus_javascript_function Resource - terraform-provider-timeplus"
subcategory: ""
description: |-
  Timeplus javascript functions are one of the supported user defined function types. Javascript functions allow users to implement functions with the javascript programming language, and be called in queries.
---

# timeplus_javascript_function (Resource)

Timeplus javascript functions are one of the supported user defined function types. Javascript functions allow users to implement functions with the javascript programming language, and be called in queries.

## Example Usage

```terraform
resource "timeplus_javascript_function" "example" {
  name        = "add"
  description = "Adds two integer and returns the sum"

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
  function add(left, right) {
    return left.map((n, i) => [n, right[i]]).map(ins => ins[0] * ins[0] + ins[1] * ins[1])
  }
  EOS
}

resource "timeplus_javascript_function" "aggregate_example" {
  name                  = "second_max"
  description           = "Returns the unique second largest number"
  is_aggregate_function = true
  return_type           = "float64"

  arg {
    name = "value"
    type = "float64"
  }

  source = <<-EOS
{
    initialize: function() {
        this.max = -Infinity;
        this.sec_max = -Infinity;
    },

    process: function(values) {
        for (const v of values) {
            this._update(v);
        }
    },

    _update: function(value) {
        if (value == this.max) {
            // skip
        } else if (value > this.max) {
            this.sec_max = this.max;
            this.max = value;
        } else if (value > this.sec_max) {
            this.sec_max = value;
        }
    },

    finalize: function() {
        return this.sec_max
    },

    serialize: function() {
        return JSON.stringify([this.max, this.sec_max]);
    },

    deserialize: function(state_str) {
        let s = JSON.parse(state_str);
        this.max = s[0];
        this.sec_max = s[1]
    },

    merge: function(state_str) {
        let s = JSON.parse(state_str);
        this._update(s[0]);
        this._update(s[1]);
    }
};
  EOS
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The javascript function name
- `return_type` (String) The type of the function's return value
- `source` (String) The javascript function source code

### Optional

- `arg` (Block List) Describe an argument of the javascript function, argument order matters (see [below for nested schema](#nestedblock--arg))
- `description` (String) A detailed text describes the javascript function
- `is_aggregate_function` (Boolean) Indecates if the javascript function an aggregate function

<a id="nestedblock--arg"></a>
### Nested Schema for `arg`

Required:

- `name` (String) The argument name
- `type` (String) The argument type