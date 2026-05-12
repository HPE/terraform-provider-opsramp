// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"flag"
	"log"

	"github.com/HPE/terraform-provider-opsramp/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/HPE/opsramp",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(provider.Version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
