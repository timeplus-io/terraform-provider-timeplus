terraform {
  required_providers {
    timeplus = {
      source  = "timeplus-io/timeplus"
      version = ">= 0.1.2"
    }
  }
}

provider "timeplus" {
  username = "my-username"
  password = "my-password"
}
