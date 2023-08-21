<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Timeplus Provider for Terraform

The Timeplus provider for Terraform is a plugin that enables full lifecycle management of Timeplus resources.

⚠️ Attention: this plugin is still in its early phase and under heavy deveployment, thus it is subject to errors, lack of features and breaking changes. You are welcome to try it out and report issues, but production usage is not recommended at the moment.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

This plug is not published to the terraform registry yet. To use it, please either follow the [Explicit Installation Method Configuration](https://developer.hashicorp.com/terraform/cli/config/config-file#explicit-installation-method-configuration) method or the [Prepare Terraform for local provider install](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install) from the tutorial. Below is an example of using the explicit installation method configuration:

```hcl
# file: .terraformrc

provider_installation {

  dev_overrides {
      "dev.timeplus.io/terraform/timeplus" = "/Users/gimi/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Once your `.terrraformrc` is properly configured, you can add the `timeplus` provider to your terraform config file and start managing Timeplus resources with it. Below is an example of using the provider to create a simeple stream. More examples can be found in the [examples](./tree/main/examples) folder.

```terraform
terraform {
  required_providers {
    timeplus = {
      source = "dev.timeplus.io/terraform/timeplus"
    }
  }
}

provider "timeplus" {
  # the workspace ID can be found in the URL https://us.timeplus.cloud/<my-workspace-id>
  workspace = "my-workspace-id"
  # API key is required to use the provider
  api_key   = "my-api-key"
}

resource "timeplus_stream" "example" {
    name = "example"
    description = "the example stream from the provider README file"
    column {
      name = "col_1"
      type = "string"
    }
    column {
      name = "col_2"
      type = "int64"
    }
}

output "example-stream" {
  value = resource.timeplus_stream.example
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.
