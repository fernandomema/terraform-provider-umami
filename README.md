# terraform-provider-umami

Provider de OpenTofu/Terraform para gestionar una instancia de Umami (self-hosted o cloud).

Recursos: `umami_website` (name, domain, share_id; expone `id` = data-website-id,
`script_url` y `tracking_script`). Siguientes candidatos: `umami_team`, `umami_user`, data sources.

Credenciales por argumentos o entorno: `UMAMI_ENDPOINT`, `UMAMI_USERNAME`, `UMAMI_PASSWORD`
(o `UMAMI_API_KEY` en Umami Cloud).

```bash
go test ./... && go build -o terraform-provider-umami .
```

Uso local con `dev_overrides` en `~/.tofurc`:
`"registry.opentofu.org/fernandomema/umami" = "/ruta/al/directorio/del/binario"`

## Instalacion

```hcl
terraform {
  required_providers {
    umami = { source = "fernandomema/umami" }
  }
}
```
