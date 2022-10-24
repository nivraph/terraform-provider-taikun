---
page_title: "taikun_images_openstack Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve images for a given OpenStack cloud credential.
---

# taikun_images_openstack (Data Source)

Retrieve images for a given OpenStack cloud credential.

~> **Role Requirement** To use the `taikun_images_openstack` data source, you need a Manager or Partner account.

## Example Usage

```terraform
resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo"
}

data "taikun_images_openstack" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_credential_id` (String) OpenStack cloud credential ID.

### Read-Only

- `id` (String) The ID of this resource.
- `images` (List of Object) List of retrieved OpenStack images. (see [below for nested schema](#nestedatt--images))

<a id="nestedatt--images"></a>
### Nested Schema for `images`

Read-Only:

- `id` (String)
- `name` (String)

