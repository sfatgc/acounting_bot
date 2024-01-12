resource "google_service_account" "accounting_bot_sa" {
  account_id   = "gcf-accounting-bot-sa"
  display_name = "Accounting Bot CF Service Account"
}

resource "google_project_iam_binding" "bot_secrets_access" {
  project = data.google_project.project.id
  role    = "roles/secretmanager.secretAccessor"

  members = [
    "serviceAccount:${google_service_account.accounting_bot_sa.email}"
  ]
}

resource "google_project_iam_binding" "bot_firestore_access" {
  project = data.google_project.project.id
  role    = "roles/datastore.user"

  members = [
    "serviceAccount:${google_service_account.accounting_bot_sa.email}"
  ]
}
