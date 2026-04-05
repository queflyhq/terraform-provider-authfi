# Create a project with data residency in India
resource "authfi_project" "production" {
  name     = "Production India"
  region   = "in"
  env_type = "production"
  tags = {
    team        = "platform"
    environment = "prod"
  }
}

# Create an organization within the project
resource "authfi_organization" "apollo" {
  name          = "Apollo Hospital"
  slug          = "apollo"
  primary_color = "#1a73e8"
}

# Create an OAuth2 application
resource "authfi_application" "web_app" {
  name          = "Patient Portal"
  app_type      = "regular_web"
  callback_urls = ["https://portal.apollohospital.com/callback"]
  logout_urls   = ["https://portal.apollohospital.com"]
  allowed_origins = ["https://portal.apollohospital.com"]
}

# Create an SSO connection
resource "authfi_connection" "google" {
  name          = "Google Workspace"
  provider_type = "google"
  type          = "social"
  client_id     = var.google_client_id
}

# Create RBAC roles
resource "authfi_role" "doctor" {
  name        = "Doctor"
  slug        = "doctor"
  description = "Medical staff with patient access"
  permissions = ["read:patients", "write:records", "read:appointments"]
}

resource "authfi_role" "admin" {
  name        = "Hospital Admin"
  slug        = "hospital-admin"
  description = "Administrative access"
  permissions = ["manage:users", "manage:billing", "read:reports"]
}

variable "google_client_id" {
  type = string
}

# Outputs
output "app_client_id" {
  value = authfi_application.web_app.client_id
}

output "app_client_secret" {
  value     = authfi_application.web_app.client_secret
  sensitive = true
}
