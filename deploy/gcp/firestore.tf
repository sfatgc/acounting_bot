resource "google_firestore_database" "database" {
  name                    = "bot_database"
  location_id             = provider.google.region
  type                    = "FIRESTORE_NATIVE"
  delete_protection_state = "DELETE_PROTECTION_ENABLED"
  deletion_policy         = "DELETE"
}
