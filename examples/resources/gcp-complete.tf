# Complete GCP + AuthFI setup
# Shows both sides: GCP infrastructure + AuthFI configuration

terraform {
  required_providers {
    authfi = {
      source = "queflyhq/authfi"
    }
    google = {
      source = "hashicorp/google"
    }
  }
}

variable "gcp_org_id" {
  default = "279469116393"
}

variable "gcp_project" {
  default = "my-gcp-project"
}

# ============================================================
# GCP Side — Workforce Identity Federation (for console login)
# ============================================================

# Workforce pool (one per org)
resource "google_iam_workforce_pool" "authfi" {
  workforce_pool_id = "authfi-workforce"
  parent            = "organizations/${var.gcp_org_id}"
  location          = "global"
  display_name      = "AuthFI Identity Pool"
  description       = "Allows AuthFI users to login to GCP Console"
}

# OIDC provider (connects AuthFI as identity source)
resource "google_iam_workforce_pool_provider" "authfi" {
  workforce_pool_id = google_iam_workforce_pool.authfi.workforce_pool_id
  location          = "global"
  provider_id       = "authfi-provider"
  display_name      = "AuthFI OIDC"

  oidc {
    issuer_uri = "https://ayush.authfi.app"
    client_id  = authfi_application.gcp_login.client_id
    client_secret {
      value = authfi_application.gcp_login.client_secret
    }
    web_sso_config {
      response_type             = "CODE"
      assertion_claims_behavior = "MERGE_USER_INFO_OVER_ID_TOKEN_CLAIMS"
    }
  }

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "google.display_name"  = "assertion.name"
    "google.groups"        = "assertion.groups"
    "attribute.email"      = "assertion.email"
    "attribute.department" = "assertion.department"
  }

  attribute_condition = "assertion.email_verified == true"
}

# IAM binding — map AuthFI users to GCP roles
resource "google_project_iam_member" "authfi_viewers" {
  project = var.gcp_project
  role    = "roles/viewer"
  member  = "principalSet://iam.googleapis.com/${google_iam_workforce_pool.authfi.name}/attribute.department/engineering"
}

resource "google_project_iam_member" "authfi_editors" {
  project = var.gcp_project
  role    = "roles/editor"
  member  = "principalSet://iam.googleapis.com/${google_iam_workforce_pool.authfi.name}/group/admin"
}

# ============================================================
# GCP Side — Workload Identity Federation (for SDK/apps)
# ============================================================

# Workload identity pool (for app-to-GCP access)
resource "google_iam_workload_identity_pool" "authfi_apps" {
  workload_identity_pool_id = "authfi-apps"
  project                   = var.gcp_project
  display_name              = "AuthFI App Pool"
  description               = "Apps authenticated via AuthFI can access GCP resources"
}

# OIDC provider for workload identity
resource "google_iam_workload_identity_pool_provider" "authfi_apps" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.authfi_apps.workload_identity_pool_id
  workload_identity_pool_provider_id = "authfi-app-provider"
  project                            = var.gcp_project
  display_name                       = "AuthFI App OIDC"

  oidc {
    issuer_uri = "https://ayush.authfi.app"
    allowed_audiences = [
      authfi_application.backend_api.client_id,
    ]
  }

  attribute_mapping = {
    "google.subject"  = "assertion.sub"
    "attribute.role"   = "assertion.role"
    "attribute.tenant" = "assertion.tenant"
  }
}

# Service account for app impersonation
resource "google_service_account" "app_sa" {
  account_id   = "authfi-app-sa"
  display_name = "AuthFI App Service Account"
  project      = var.gcp_project
}

# Allow AuthFI-authenticated apps to impersonate the SA
resource "google_service_account_iam_binding" "authfi_app_binding" {
  service_account_id = google_service_account.app_sa.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.authfi_apps.name}/*",
  ]
}

# Grant SA permissions on GCP resources
resource "google_project_iam_member" "app_sa_storage" {
  project = var.gcp_project
  role    = "roles/storage.objectViewer"
  member  = "serviceAccount:${google_service_account.app_sa.email}"
}

# ============================================================
# AuthFI Side — Applications + Cloud Accounts
# ============================================================

# Application for GCP Console login flow
resource "authfi_application" "gcp_login" {
  name     = "GCP Console Login"
  app_type = "regular_web"
  callback_urls = [
    "https://auth.cloud.google/signin-callback/locations/global/workforcePools/authfi-workforce/providers/authfi-provider"
  ]
}

# Application for backend API (workload identity)
resource "authfi_application" "backend_api" {
  name     = "Backend API"
  app_type = "m2m"
}

# Register GCP Console access in AuthFI
resource "authfi_cloud_account" "gcp_console" {
  name          = "GCP Console"
  provider_type = "gcp"
  access_type   = "workforce_identity"
  config = {
    organization_id = var.gcp_org_id
    workforce_pool  = google_iam_workforce_pool.authfi.workforce_pool_id
    provider_id     = google_iam_workforce_pool_provider.authfi.provider_id
    project_id      = var.gcp_project
    console_url     = "https://console.cloud.google.com/?workforce_pool=${google_iam_workforce_pool.authfi.workforce_pool_id}"
  }
}

# Register GCP Workload access in AuthFI
resource "authfi_cloud_account" "gcp_workload" {
  name          = "GCP App Access"
  provider_type = "gcp"
  access_type   = "workload_identity"
  config = {
    project_id       = var.gcp_project
    workload_pool    = google_iam_workload_identity_pool.authfi_apps.workload_identity_pool_id
    provider_id      = google_iam_workload_identity_pool_provider.authfi_apps.workload_identity_pool_provider_id
    service_account  = google_service_account.app_sa.email
  }
}

# Map roles to cloud access
resource "authfi_cloud_mapping" "admin_gcp_editor" {
  cloud_account_id = authfi_cloud_account.gcp_console.id
  role_id          = authfi_role.admin.id
  cloud_role       = "roles/editor"
}

resource "authfi_cloud_mapping" "doctor_gcp_viewer" {
  cloud_account_id = authfi_cloud_account.gcp_console.id
  role_id          = authfi_role.doctor.id
  cloud_role       = "roles/viewer"
}
