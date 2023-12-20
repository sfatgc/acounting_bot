resource "google_secret_manager_secret" "accounting_bot_credentials_secret" {
  secret_id = "accounting_bot_credentials"

  replication {
    auto {}
  }
}


resource "google_secret_manager_secret_version" "accounting_bot_credentials_secret_version" {
  secret = google_secret_manager_secret.accounting_bot_credentials_secret.id

  secret_data = var.telegram_bot_token
}
