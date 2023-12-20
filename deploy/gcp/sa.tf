resource "google_service_account" "accounting_bot_sa" {
  account_id   = "gcf-accounting-bot-sa"
  display_name = "Accounting Bot CF Service Account"
}

resource "google_project_iam_binding" "project" {
  project = data.google_project.project.id
  role    = "roles/secretmanager.secretAccessor"

  members = [
    "serviceAccount:${google_service_account.accounting_bot_sa.email}"
  ]
}
