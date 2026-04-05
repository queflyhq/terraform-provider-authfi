terraform {
  required_providers {
    authfi = {
      source  = "queflyhq/authfi"
      version = "~> 0.1"
    }
  }
}

provider "authfi" {
  api_key = var.authfi_api_key # or set AUTHFI_API_KEY env var
  tenant  = var.authfi_tenant  # or set AUTHFI_TENANT env var
}

variable "authfi_api_key" {
  type      = string
  sensitive = true
}

variable "authfi_tenant" {
  type = string
}
