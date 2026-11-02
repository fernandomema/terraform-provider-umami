package main

import (
	"context"
	"flag"
	"log"

	"github.com/fernandomema/terraform-provider-umami/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider in debug mode")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.opentofu.org/fernandomema/umami",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
