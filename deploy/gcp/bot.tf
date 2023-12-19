provider "telegram" {}

resource "telegram_bot_webhook" "example" {
  url             = google_cloudfunctions2_function.default.url
  max_connections = 5
}

resource "telegram_bot_commands" "example" {
  commands = [
    {
      command     = "start",
      description = "View welcome message"
    },
    {
      command     = "help",
      description = "Show help"
    }
  ]
}

