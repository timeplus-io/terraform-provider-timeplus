terraform {
  required_providers {
    timeplus = {
      source = "timeplus-io/timeplus"
    }
  }
}

provider "timeplus" {
  workspace = "my-workspace-id"
  api_key   = "my-api-key"
}
