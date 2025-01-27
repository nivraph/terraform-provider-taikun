---
page_title: "taikun_images Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve images for a given cloud credential.
---

# taikun_images (Data Source)

Retrieve images for a given cloud credential.

~> **Role Requirement** To use the `taikun_images` data source, you need a Manager or Partner account.

!> **Deprecated** The `taikun_images` data source is deprecated in favour of
`taikun_images_aws`, `taikun_images_azure`, `taikun_images_gcp` and
`taikun_images_openstack`.

## Example Usage

```terraform
resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo"
}

data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_credential_id` (String) Cloud credential ID.

### Optional

- `aws_limit` (Number) Limit the number of listed AWS images (highly recommended as fetching the entire list of images can take a long time) (only valid with AWS cloud credential ID).
- `azure_offer` (String) Azure offer (only valid with Azure Cloud Credential ID).
- `azure_publisher` (String) Azure publisher (only valid with Azure Cloud Credential ID).
- `azure_sku` (String) Azure sku (only valid with Azure Cloud Credential ID).

### Read-Only

- `id` (String) The ID of this resource.
- `images` (List of Object) List of retrieved images. (see [below for nested schema](#nestedatt--images))

<a id="nestedatt--images"></a>
### Nested Schema for `images`

Read-Only:

- `id` (String)
- `name` (String)


