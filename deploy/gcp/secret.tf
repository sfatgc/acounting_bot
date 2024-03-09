resource "google_secret_manager_secret" "bot_cred_tg_secret" {
  secret_id = "bot_cred_tg"

  replication {
    auto {}
  }
}


resource "google_secret_manager_secret_version" "bot_cred_tg_secret_version" {
  secret = google_secret_manager_secret.bot_cred_tg_secret.id

  secret_data = var.telegram_bot_token
}

resource "google_secret_manager_secret" "bot_cred_stripe_secret" {
  secret_id = "bot_cred_stripe"

  replication {
    auto {}
  }
}


resource "google_secret_manager_secret_version" "bot_cred_stripe_secret_version" {
  secret = google_secret_manager_secret.bot_cred_stripe_secret.id

  secret_data = var.pp_stripe_token
}
