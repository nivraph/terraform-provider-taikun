---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_cloud_credential_gcp Resource - terraform-provider-taikun"
subcategory: ""
description: |-
  Taikun Google Cloud Platform Credential
---

# taikun_cloud_credential_gcp (Resource)

Taikun Google Cloud Platform Credential



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `az_count` (String) The number of GCP availability zone expected for the region.
- `config_file` (String) The path of the GCP credential's configuration file.
- `name` (String) The name of the GCP credential.
- `region` (String) The region of the GCP credential.

### Optional

- `billing_account_id` (String) The ID of the GCP credential's billing account. Conflicts with: `import_project`.
- `folder_id` (String) The folder ID of the GCP credential. Conflicts with: `import_project`.
- `import_project` (Boolean) Whether to import a project or not Defaults to `false`. Conflicts with: `billing_account_id`, `folder_id`.
- `lock` (Boolean) Indicates whether to lock the GCP cloud credential. Defaults to `false`.
- `organization_id` (String) The ID of the organization which owns the GCP credential.

### Read-Only

- `billing_account_name` (String) The name of the GCP credential's billing account.
- `id` (String) The ID of the GCP credential.
- `is_default` (Boolean) Indicates whether the GCP cloud credential is the default one.
- `organization_name` (String) The name of the organization which owns the GCP credential.
- `zones` (Set of String) The given zones of the GCP credential.

