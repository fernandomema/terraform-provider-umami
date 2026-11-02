terraform {
  required_providers {
    umami = {
      source = "fernandomema/umami"
    }
  }
}

provider "umami" {
  endpoint = "https://umami.example.com"
}
