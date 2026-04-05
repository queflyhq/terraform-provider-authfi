# Email & SMS Templates — configure via Terraform

# Email Templates
resource "authfi_email_template" "verification" {
  type    = "verification"
  subject = "Verify your email — {{tenant_name}}"
  body    = <<-HTML
    <h2>Welcome to {{tenant_name}}</h2>
    <p>Click the link below to verify your email address:</p>
    <a href="{{verification_url}}" style="background:{{primary_color}};color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;">
      Verify Email
    </a>
    <p>This link expires in 24 hours.</p>
  HTML
  from_name  = "Ayush Healthcare"
  from_email = "noreply@ayush.live"
}

resource "authfi_email_template" "password_reset" {
  type    = "password_reset"
  subject = "Reset your password — {{tenant_name}}"
  body    = <<-HTML
    <h2>Password Reset</h2>
    <p>Click below to reset your password:</p>
    <a href="{{reset_url}}" style="background:{{primary_color}};color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;">
      Reset Password
    </a>
    <p>If you didn't request this, ignore this email.</p>
  HTML
}

resource "authfi_email_template" "invite" {
  type    = "invite"
  subject = "You've been invited to {{tenant_name}}"
  body    = <<-HTML
    <h2>You're invited!</h2>
    <p>{{inviter_name}} has invited you to join {{tenant_name}}.</p>
    <a href="{{invite_url}}" style="background:{{primary_color}};color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;">
      Accept Invitation
    </a>
  HTML
}

resource "authfi_email_template" "mfa_code" {
  type    = "mfa_code"
  subject = "Your verification code — {{tenant_name}}"
  body    = <<-HTML
    <h2>Verification Code</h2>
    <p>Your code is: <strong>{{code}}</strong></p>
    <p>This code expires in 10 minutes.</p>
  HTML
}

resource "authfi_email_template" "magic_link" {
  type    = "magic_link"
  subject = "Sign in to {{tenant_name}}"
  body    = <<-HTML
    <h2>Sign In</h2>
    <p>Click below to sign in:</p>
    <a href="{{magic_link_url}}" style="background:{{primary_color}};color:#fff;padding:12px 24px;border-radius:8px;text-decoration:none;">
      Sign In
    </a>
    <p>This link expires in 15 minutes.</p>
  HTML
}

# SMS Templates
resource "authfi_sms_template" "otp" {
  type = "otp"
  body = "Your {{tenant_name}} code is {{code}}. Expires in 10 min."
}

resource "authfi_sms_template" "magic_link" {
  type = "magic_link"
  body = "Sign in to {{tenant_name}}: {{magic_link_url}}"
}

resource "authfi_sms_template" "verification" {
  type = "verification"
  body = "Verify your {{tenant_name}} account: {{verification_url}}"
}
