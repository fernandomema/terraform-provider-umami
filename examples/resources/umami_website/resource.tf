resource "umami_website" "blog" {
  name   = "Blog"
  domain = "blog.example.com"
}

output "tracking_script" {
  value = umami_website.blog.tracking_script
}
