resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name                        = "${random_id.default.hex}-gcf-source" # Every bucket name must be globally unique
  location                    = "US"
  uniform_bucket_level_access = true
}

locals {
  function_filename_timestamp = formatdate("ZZZhhmmDDMMMYYYY", timestamp())
  function_filename           = "/tmp/function-source-${local.function_filename_timestamp}.zip"
}

data "archive_file" "default" {
  type        = "zip"
  output_path = local.function_filename
  source_dir  = "../../functions/accounting_bot/"
}
resource "google_storage_bucket_object" "object" {
  name   = "function-source-${data.archive_file.default.output_sha512}.zip"
  bucket = google_storage_bucket.default.name
  source = data.archive_file.default.output_path # Add path to the zipped function source code
}



locals {
  google_firestore_db_id = split("/", google_firestore_database.database.id)
}
resource "google_cloudfunctions2_function" "default" {
  name        = "function-v2"
  location    = "us-west1"
  description = "accounting bot dispatchMessages function"

  build_config {
    runtime     = "go121"
    entry_point = "dispatchMessages" # Set the entry point
    source {
      storage_source {
        bucket = google_storage_bucket.default.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    service_account_email = google_service_account.accounting_bot_sa.email
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60

    environment_variables = {
      "GOOGLE_PROJECT_ID"      = split("/", data.google_project.project.id)[1]
      "GOOGLE_FIRESTORE_DB_ID" = element(local.google_firestore_db_id, length(local.google_firestore_db_id) - 1)
    }

    secret_environment_variables {
      key        = "TELEGRAM_BOT_TOKEN"
      project_id = data.google_project.project.project_id
      secret     = google_secret_manager_secret.bot_cred_tg_secret.secret_id
      version    = google_secret_manager_secret_version.bot_cred_tg_secret_version.version
    }

    secret_environment_variables {
      key        = "PP_STRIPE_TOKEN"
      project_id = data.google_project.project.project_id
      secret     = google_secret_manager_secret.bot_cred_stripe_secret.secret_id
      version    = google_secret_manager_secret_version.bot_cred_stripe_secret_version.version
    }
  }

  depends_on = [
    google_project_service.project,
    google_secret_manager_secret_version.bot_cred_tg_secret_version,
    google_secret_manager_secret_version.bot_cred_stripe_secret_version
  ]

}

resource "google_cloud_run_service_iam_member" "member" {
  location = google_cloudfunctions2_function.default.location
  service  = google_cloudfunctions2_function.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "function_uri" {
  value = google_cloudfunctions2_function.default.url
}
