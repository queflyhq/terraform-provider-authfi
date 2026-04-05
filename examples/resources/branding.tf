# Branding — configure login page appearance via Terraform

# Tenant-level branding (default for all orgs)
resource "authfi_branding" "tenant" {
  name             = "Ayush Healthcare Pvt. Ltd."
  logo_url         = "https://ayush.live/logo.png"
  favicon_url      = "https://ayush.live/favicon.ico"
  primary_color    = "#1a73e8"
  background_color = "#0a1628"
  text_color       = "#ffffff"
  button_radius    = "12px"
  font_family      = "Inter"
  theme_mode       = "dark"

  # Login page layout
  login_layout     = "split"        # split, centered, fullscreen
  social_position  = "bottom"       # top, bottom
  show_attribution = false          # hide "Powered by AuthFI" (paid plan)

  # Custom text
  welcome_text  = "Sign in to your account"
  subtitle_text = "Use your Ayush Healthcare Account"
  support_email = "support@ayush.live"
  privacy_url   = "https://ayush.live/privacy"
  terms_url     = "https://ayush.live/terms"

  # Custom CSS (paid plan)
  custom_css = <<-CSS
    .login-card { box-shadow: 0 8px 32px rgba(0,0,0,0.3); }
  CSS

  # Password policy
  password_min_length       = 8
  password_require_uppercase = true
  password_require_number    = true
  password_require_special   = true
  max_login_attempts        = 5
  mfa_policy                = "optional"  # optional, required, adaptive
}

# Organization branding override (inherits from tenant, overrides what's set)
resource "authfi_org_branding" "apollo" {
  organization_id  = authfi_organization.apollo.id

  # Only set what differs from tenant
  name             = "Apollo Hospital"
  logo_url         = "https://apollohospital.com/logo.png"
  primary_color    = "#e53935"
  welcome_text     = "Sign in to Apollo Hospital"
  subtitle_text    = "Healthcare staff portal"

  # Inherits from tenant: background_color, text_color, font_family,
  # button_radius, theme_mode, password_policy, etc.
}

# Custom domain for org
resource "authfi_domain" "apollo" {
  domain = "login.apollohospital.com"
  type   = "organization"
  org_id = authfi_organization.apollo.id
  # DNS: CNAME login.apollohospital.com → ayush.authfi.app
}

# Custom domain for tenant
resource "authfi_domain" "ayush" {
  domain = "login.ayush.live"
  type   = "tenant"
  # DNS: CNAME login.ayush.live → ayush.authfi.app
}
