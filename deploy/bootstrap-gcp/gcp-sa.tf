resource "google_service_account" "service_account" {
  account_id   = "tfcloud"
  display_name = "TFCloud"
  description  = "Terraform Cloud SA"
}

locals {
  service_account_roles = [
    #"roles/owner",
    "roles/editor",
    "roles/resourcemanager.projectIamAdmin",
    "roles/cloudfunctions.admin",
    "roles/run.admin",
    "roles/secretmanager.admin",
    "roles/datastore.owner",
  ]
}

resource "google_project_iam_binding" "service_account_iam_binding" {
  count   = length(local.service_account_roles)
  project = data.google_project.project.id
  role    = local.service_account_roles[count.index]

  members = [
    "serviceAccount:${google_service_account.service_account.email}"
  ]
}
