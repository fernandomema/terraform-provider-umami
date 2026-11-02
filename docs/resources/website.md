---
page_title: "umami_website Resource - umami"
subcategory: ""
description: |-
  An Umami website. Its id is the data-website-id used by the tracker.
---

# umami_website (Resource)

An Umami website. Its `id` is the `data-website-id` used by the tracker.

## Example Usage

```terraform
resource "umami_website" "blog" {
  name   = "Blog"
  domain = "blog.example.com"
}

output "tracking_script" {
  value = umami_website.blog.tracking_script
}
```

## Schema

### Required

- `name` (String) Display name.
- `domain` (String) Domain of the website, without protocol.

### Optional

- `share_id` (String) ID used for the public share link (`/share/<id>`).

### Read-Only

- `id` (String) Website ID, i.e. the `data-website-id`.
- `script_url` (String) URL of the instance's `script.js`.
- `tracking_script` (String) Ready-to-paste `<script>` tag.

## Import

```shell
tofu import umami_website.blog <website-id>
```
