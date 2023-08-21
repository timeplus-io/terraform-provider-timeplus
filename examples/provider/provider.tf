terraform {
  required_providers {
    timeplus = {
      source = "timeplus.io/terraform/timeplus"
    }
  }
}

provider "timeplus" {
  workspace = "my-workspace-id"
  api_key   = "my-api-key"
}
