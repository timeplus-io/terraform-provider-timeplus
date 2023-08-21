terraform {
  required_providers {
    timeplus = {
      source = "dev.timeplus.com/terraform/timeplus"
    }
  }
}

provider "timeplus" {
  workspace = "my-workspace-id"
  api_key   = "my-api-key"
}
