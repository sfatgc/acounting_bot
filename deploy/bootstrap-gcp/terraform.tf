terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.34.0"
    }
    telegram = {
      source  = "yi-jiayu/telegram"
      version = "0.3.1"
    }
    tfe = {
      source  = "hashicorp/tfe"
      version = "0.51.1"
    }
  }
}
