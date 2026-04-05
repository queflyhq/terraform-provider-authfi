# Cloud Access — federate identity to AWS, GCP, Azure

# GCP Workload Identity Federation
resource "authfi_cloud_account" "gcp_prod" {
  name     = "GCP Production"
  provider_type = "gcp"
  config = {
    project_id      = "my-gcp-project"
    workforce_pool   = "authfi-workforce"
    provider_id      = "authfi-provider"
    service_account  = "app-sa@my-gcp-project.iam.gserviceaccount.com"
  }
}

# AWS IAM Role Federation
resource "authfi_cloud_account" "aws_prod" {
  name     = "AWS Production"
  provider_type = "aws"
  config = {
    account_id = "123456789012"
    role_arn   = "arn:aws:iam::123456789012:role/AuthFIFederated"
    region     = "us-east-1"
  }
}

# Azure AD Workload Identity
resource "authfi_cloud_account" "azure_prod" {
  name     = "Azure Production"
  provider_type = "azure"
  config = {
    tenant_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    subscription  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}

# Map roles to cloud access
resource "authfi_cloud_mapping" "doctor_gcp" {
  cloud_account_id = authfi_cloud_account.gcp_prod.id
  role_id          = authfi_role.doctor.id
  cloud_role       = "roles/viewer"
}

resource "authfi_cloud_mapping" "admin_aws" {
  cloud_account_id = authfi_cloud_account.aws_prod.id
  role_id          = authfi_role.admin.id
  cloud_role       = "arn:aws:iam::123456789012:policy/AdminAccess"
}

# eBPF Access Policy — zero trust service mesh
resource "authfi_access_policy" "api_to_db" {
  name   = "API to Database"
  action = "allow"
  source = {
    service = "api-server"
    role    = "doctor"
  }
  destination = {
    service = "patient-db"
    port    = 5432
  }
}

resource "authfi_access_policy" "deny_external" {
  name     = "Deny External Access"
  action   = "deny"
  priority = 100
  source = {
    external = true
  }
  destination = {
    service = "internal-api"
  }
}

# Agent enrollment token
resource "authfi_agent_token" "k8s_cluster" {
  name        = "K8s Production Cluster"
  description = "eBPF agent for production Kubernetes"
}

output "agent_enrollment_token" {
  value     = authfi_agent_token.k8s_cluster.token
  sensitive = true
}
