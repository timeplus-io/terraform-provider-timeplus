terraform {
  required_providers {
    timeplus = {
      source  = "timeplus-io/timeplus"
      version = ">= 0.1.2"
    }
  }
}

provider "timeplus" {
  workspace = "my-workspace-id"
  api_key   = "my-api-key"
}
