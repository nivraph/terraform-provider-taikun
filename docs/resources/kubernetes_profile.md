---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_kubernetes_profile Resource - terraform-provider-taikun"
subcategory: ""
description: |-
  Taikun Kubernetes Profile
---

# taikun_kubernetes_profile (Resource)

Taikun Kubernetes Profile

## Example Usage

```terraform
resource "taikun_kubernetes_profile" "foo" {
  # Required
  name = "foo"

  # Optional
  organization_id         = "42"
  load_balancing_solution = "Taikun"
  bastion_proxy_enabled   = true
  is_locked               = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the Kubernetes profile.

### Optional

- **bastion_proxy_enabled** (Boolean) Exposes the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`. Defaults to `false`.
- **is_locked** (Boolean) Indicates whether the Kubernetes profile is locked or not. Defaults to `false`.
- **load_balancing_solution** (String) Load-balancing solution: `None`, `Octavia` or `Taikun`. `Octavia` and `Taikun` are only available for OpenStack cloud. Defaults to `Octavia`.
- **organization_id** (String) The id of the organization which owns the Kubernetes profile.

### Read-Only

- **cni** (String) Container Network Interface(CNI) of the Kubernetes profile.
- **created_by** (String) The creator of the Kubernetes profile.
- **id** (String) The id of the Kubernetes profile.
- **last_modified** (String) Time of last modification.
- **last_modified_by** (String) The last user who modified the Kubernetes profile.
- **organization_name** (String) The name of the organization which owns the Kubernetes profile.

## Import

Import is supported using the following syntax:

```shell
# import with Taikun ID
terraform import taikun_kubernetes_profile.myprofile 42
```