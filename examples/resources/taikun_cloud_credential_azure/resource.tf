resource "taikun_cloud_credential_azure" "foo" {
  # Required
  name = "foo"

  client_id         = "client_id"
  client_secret     = "client_secret"
  subscription_id   = "subscription_id"
  tenant_id         = "tenant_id"
  location          = "location"
  availability_zone = "availability_zone"

  # Optional
  organization_id = "42"
  is_locked       = false
}