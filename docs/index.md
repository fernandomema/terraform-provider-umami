---
page_title: "Provider: Umami"
description: |-
  Manage Umami analytics websites with OpenTofu or Terraform.
---

# Umami Provider

Manage [Umami](https://umami.is) (self-hosted or Cloud) from OpenTofu/Terraform.

## Example Usage

```terraform
terraform {
  required_providers {
    umami = {
      source = "fernandomema/umami"
    }
  }
}

# Credentials via UMAMI_USERNAME / UMAMI_PASSWORD (or UMAMI_API_KEY for Umami Cloud).
provider "umami" {
  endpoint = "https://umami.example.com"
}
```

## Schema

### Optional

- `endpoint` (String) Base URL of the instance, e.g. `https://umami.example.com`. Environment variable: `UMAMI_ENDPOINT`.
- `username` (String) Username (self-hosted). Environment variable: `UMAMI_USERNAME`.
- `password` (String, Sensitive) Password (self-hosted). Environment variable: `UMAMI_PASSWORD`.
- `api_key` (String, Sensitive) Umami Cloud API key, an alternative to username/password. Environment variable: `UMAMI_API_KEY`.
- `insecure` (Boolean) Skip TLS certificate verification. Defaults to `false`.

Either `api_key`, or both `username` and `password`, must be provided.
