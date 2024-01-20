# DOC: https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference
# gcloud auth application-default login
provider "google" {
  region = "us-west1"
}
data "google_project" "project" {}

# DOC: https://registry.terraform.io/providers/hashicorp/tfe/latest/docs
provider "tfe" {

}
