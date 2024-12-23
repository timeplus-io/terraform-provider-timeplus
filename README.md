<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Timeplus Provider for Terraform

The Timeplus provider for Terraform is a plugin that enables full lifecycle management of Timeplus resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20.0

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

To use the provider, simply add it to your terraform file, for example:

```terraform
terraform {
  required_providers {
    timeplus = {
      source  = "timeplus-io/timeplus"
      version = ">= 0.1.2"
    }
  }
}

provider "timeplus" {
  endpoint = "http://localhost:8000"
  workspace = "default"
  username  = "proton"
  password  = "proton@t+"
}
```

Then you can start provisioning Timeplus resources, and below is an example of stream:

```terraform
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
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory. Please follow [Prepare Terraform for local provider install](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install) to use the locally-built provider to test it.

To generate or update documentation, run `go generate`.

## Useful documentations for provider development

- Timeplus document web site: https://docs.timeplus.com/
- Terraform plugin framework doc: https://developer.hashicorp.com/terraform/plugin/framework
