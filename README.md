# OpenTofu / Terraform Provider for Umami

Manage your [Umami](https://umami.is) analytics as code. This provider works with
self-hosted Umami and Umami Cloud, and lets you create the websites your apps track
(and get their tracking snippet) from the same OpenTofu/Terraform configuration that
deploys them.

- **Registry:** `fernandomema/umami`
- **Resources:** [`umami_website`](docs/resources/website.md)
- **Compatibility:** OpenTofu >= 1.6 and Terraform >= 1.0 (plugin protocol 6)

## Quick start

```hcl
terraform {
  required_providers {
    umami = {
      source = "fernandomema/umami"
    }
  }
}

provider "umami" {
  endpoint = "https://umami.example.com"
  # username / password via UMAMI_USERNAME / UMAMI_PASSWORD
}

resource "umami_website" "blog" {
  name   = "Blog"
  domain = "blog.example.com"
}

output "tracking_script" {
  value = umami_website.blog.tracking_script
}
```

`tracking_script` is a ready-to-paste tag:

```html
<script defer src="https://umami.example.com/script.js" data-website-id="..."></script>
```

You can also feed `umami_website.blog.id` straight into your app's environment
(for example a `PUBLIC_UMAMI_WEBSITE_ID` variable) so a new service gets its analytics
site created and wired in a single `apply`.

## Authentication

| Method | Arguments | Environment variables | Use with |
|---|---|---|---|
| Username + password | `username`, `password` | `UMAMI_USERNAME`, `UMAMI_PASSWORD` | Self-hosted |
| API key | `api_key` | `UMAMI_API_KEY` | Umami Cloud |

The endpoint can be set with `endpoint` or `UMAMI_ENDPOINT`. Prefer environment
variables for credentials so they stay out of your configuration and state.

> Use a dedicated Umami user for automation rather than your personal admin account.

Set `insecure = true` only if your instance uses a self-signed TLS certificate.

## Resource: `umami_website`

| Argument | Required | Description |
|---|---|---|
| `name` | yes | Display name of the website. |
| `domain` | yes | Domain being tracked, without protocol (`example.com`). |
| `share_id` | no | Enables a public share link (`/share/<share_id>`). |

| Attribute | Description |
|---|---|
| `id` | Website ID. This is the `data-website-id` used by the tracker. |
| `script_url` | URL of the instance's `script.js`. |
| `tracking_script` | Complete `<script>` tag for the site. |

Existing websites can be imported by ID:

```sh
tofu import umami_website.blog <website-id>
```

## Development

```sh
go test ./...
go build -o terraform-provider-umami .
```

To use a local build without publishing, add a development override to `~/.tofurc`
(or `~/.terraformrc`):

```hcl
provider_installation {
  dev_overrides {
    "registry.opentofu.org/fernandomema/umami" = "/path/to/directory/with/binary"
  }
  direct {}
}
```

Releases are built and signed with GoReleaser when a `v*` tag is pushed.

## Roadmap

`umami_team`, `umami_user`, and data sources for existing websites. Issues and pull
requests are welcome.

## License

[MIT](LICENSE)
