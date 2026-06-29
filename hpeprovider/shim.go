package provider

import (
	internalprovider "github.com/HPE/terraform-provider-opsramp/internal/provider"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
)

func New() fwprovider.Provider {
	return internalprovider.New(internalprovider.Version)()
}
