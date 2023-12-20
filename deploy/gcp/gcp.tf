provider "google" {
  region = "us-west1"
}

data "google_project" "project" {}

locals {
  services = [
    "iam.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "run.googleapis.com",
    "cloudbuild.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}

resource "google_project_service" "project" {
  count   = length(local.services)
  service = local.services[count.index]

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
  disable_on_destroy         = true
}
