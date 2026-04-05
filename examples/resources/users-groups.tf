# Create users (also provisioned via SCIM inbound from Okta/Entra)
resource "authfi_user" "doctor_smith" {
  email    = "smith@apollohospital.com"
  name     = "Dr. Smith"
  password = var.default_password # or omit for SCIM-provisioned users
  org_id   = authfi_organization.apollo.id
}

resource "authfi_user" "nurse_jones" {
  email  = "jones@apollohospital.com"
  name   = "Nurse Jones"
  org_id = authfi_organization.apollo.id
}

# Create groups
resource "authfi_group" "medical_staff" {
  name = "Medical Staff"
  slug = "medical-staff"
}

resource "authfi_group" "admin_staff" {
  name = "Admin Staff"
  slug = "admin-staff"
}

# Assign users to groups
resource "authfi_group_member" "smith_medical" {
  group_id = authfi_group.medical_staff.id
  user_id  = authfi_user.doctor_smith.id
}

resource "authfi_group_member" "jones_medical" {
  group_id = authfi_group.medical_staff.id
  user_id  = authfi_user.nurse_jones.id
}

# Assign roles to users
resource "authfi_user_role" "smith_doctor" {
  user_id = authfi_user.doctor_smith.id
  role_id = authfi_role.doctor.id
}

resource "authfi_user_role" "jones_doctor" {
  user_id = authfi_user.nurse_jones.id
  role_id = authfi_role.doctor.id
}

# Assign roles to groups (all members inherit)
resource "authfi_group_role" "medical_staff_doctor" {
  group_id = authfi_group.medical_staff.id
  role_id  = authfi_role.doctor.id
}

# Create permissions
resource "authfi_permission" "read_patients" {
  name = "Read Patients"
  slug = "read:patients"
}

resource "authfi_permission" "write_records" {
  name = "Write Records"
  slug = "write:records"
}

resource "authfi_permission" "manage_users" {
  name = "Manage Users"
  slug = "manage:users"
}

# Map permissions to roles
resource "authfi_role_permission" "doctor_read" {
  role_id       = authfi_role.doctor.id
  permission_id = authfi_permission.read_patients.id
}

resource "authfi_role_permission" "doctor_write" {
  role_id       = authfi_role.doctor.id
  permission_id = authfi_permission.write_records.id
}

# SCIM outbound target (push users to external system)
resource "authfi_scim_target" "okta" {
  name     = "Okta Directory"
  endpoint = "https://dev-xxxxx.okta.com/scim/v2"
  token    = var.okta_scim_token
}

variable "default_password" {
  type      = string
  sensitive = true
  default   = ""
}

variable "okta_scim_token" {
  type      = string
  sensitive = true
}
