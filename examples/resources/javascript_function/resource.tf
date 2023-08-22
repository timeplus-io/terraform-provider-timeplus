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
  description           = "Returns the second largest number"
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
            console.log("this is a log from second_max UDA")
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
