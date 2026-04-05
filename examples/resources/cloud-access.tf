# Cloud Access — two modes:
# 1. Workforce Identity: humans login to cloud consoles via AuthFI
# 2. Workload Identity: apps/SDKs get cloud credentials via AuthFI tokens

# ============================================================
# GCP — Workforce Identity (console login)
# ============================================================
resource "authfi_cloud_account" "gcp_console" {
  name          = "GCP Console Access"
  provider_type = "gcp"
  access_type   = "workforce_identity"
  config = {
    organization_id = "279469116393"
    workforce_pool  = "authfi-workforce"
    provider_id     = "authfi-provider"
    project_id      = "my-gcp-project"
  }
}

# GCP — Workload Identity (SDK / app-to-cloud)
resource "authfi_cloud_account" "gcp_workload" {
  name          = "GCP Workload Access"
  provider_type = "gcp"
  access_type   = "workload_identity"
  config = {
    project_id          = "my-gcp-project"
    workload_pool       = "authfi-workload-pool"
    provider_id         = "authfi-app-provider"
    service_account     = "app-sa@my-gcp-project.iam.gserviceaccount.com"
    # Apps exchange AuthFI JWT → GCP token via STS
  }
}

# ============================================================
# AWS — IAM Identity Center (console login)
# ============================================================
resource "authfi_cloud_account" "aws_console" {
  name          = "AWS Console Access"
  provider_type = "aws"
  access_type   = "iam_identity_center"
  config = {
    account_id   = "123456789012"
    sso_role_arn = "arn:aws:iam::123456789012:role/AuthFIFederated"
    region       = "us-east-1"
  }
}

# AWS — IAM Roles Anywhere (SDK / app-to-cloud)
resource "authfi_cloud_account" "aws_workload" {
  name          = "AWS Workload Access"
  provider_type = "aws"
  access_type   = "roles_anywhere"
  config = {
    account_id     = "123456789012"
    trust_anchor   = "arn:aws:rolesanywhere:us-east-1:123456789012:trust-anchor/xxxx"
    profile_arn    = "arn:aws:rolesanywhere:us-east-1:123456789012:profile/xxxx"
    role_arn       = "arn:aws:iam::123456789012:role/AppRole"
    # Apps exchange AuthFI JWT → AWS STS credentials
  }
}

# ============================================================
# Azure — External Identity (console login)
# ============================================================
resource "authfi_cloud_account" "azure_console" {
  name          = "Azure Portal Access"
  provider_type = "azure"
  access_type   = "external_identity"
  config = {
    tenant_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    subscription = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}

# Azure — Workload Identity Federation (SDK / app-to-cloud)
resource "authfi_cloud_account" "azure_workload" {
  name          = "Azure Workload Access"
  provider_type = "azure"
  access_type   = "workload_identity"
  config = {
    tenant_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    audience     = "api://authfi-workload"
    # Apps exchange AuthFI JWT → Azure AD token
  }
}

# ============================================================
# Role → Cloud IAM Mappings
# ============================================================

# Console access mappings (who can login to cloud consoles)
resource "authfi_cloud_mapping" "admin_gcp_console" {
  cloud_account_id = authfi_cloud_account.gcp_console.id
  role_id          = authfi_role.admin.id
  cloud_role       = "roles/editor"
}

resource "authfi_cloud_mapping" "admin_aws_console" {
  cloud_account_id = authfi_cloud_account.aws_console.id
  role_id          = authfi_role.admin.id
  cloud_role       = "AdministratorAccess"
}

resource "authfi_cloud_mapping" "admin_azure_console" {
  cloud_account_id = authfi_cloud_account.azure_console.id
  role_id          = authfi_role.admin.id
  cloud_role       = "Contributor"
}

# Workload access mappings (what apps/SDKs can do in cloud)
resource "authfi_cloud_mapping" "app_gcp_storage" {
  cloud_account_id = authfi_cloud_account.gcp_workload.id
  role_id          = authfi_role.doctor.id
  cloud_role       = "roles/storage.objectViewer"
}

resource "authfi_cloud_mapping" "app_aws_s3" {
  cloud_account_id = authfi_cloud_account.aws_workload.id
  role_id          = authfi_role.doctor.id
  cloud_role       = "arn:aws:iam::123456789012:policy/S3ReadOnly"
}
