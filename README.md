# Terraform Provider for OpsRamp

<!--
(C) Copyright 2026 Hewlett Packard Enterprise Development LP
-->

A [Terraform](https://www.terraform.io/) provider to manage resources via the [OpsRamp](https://www.opsramp.com/) REST API, built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Features

- Full CRUD lifecycle management for OpsRamp resources
- Multi-tenant support (partner and client scopes)
- Rate-limited API transport to avoid throttling
- Import support for most resource types

## Supported Resources

| Resource | Type Name | Description |
|----------|-----------|-------------|
| Resource | `opsramp_resource` | Managed devices/resources |
| Service Map | `opsramp_servicemap` | Service group topology |
| Service Map Link | `opsramp_servicemap_link` | Links between service maps |
| Device Group | `opsramp_device_group` | Logical grouping of devices |
| Client | `opsramp_client` | Sub-tenant management |
| User | `opsramp_user` | User accounts |
| Role | `opsramp_role` | Access control roles |
| Permission Set | `opsramp_permission_set` | Permission bundles for roles |
| User Group | `opsramp_user_group` | User group management |
| Custom Integration | `opsramp_custom_integration` | OAuth2 API integrations |
| Script | `opsramp_script` | Executable scripts |
| Script Category | `opsramp_script_category` | Script category organization |
| KB Category | `opsramp_kb_category` | Knowledge base categories |
| KB Article | `opsramp_kb_article` | Knowledge base articles |
| ServiceDesk Category | `opsramp_servicedesk_category` | Ticket categories |
| ServiceDesk Urgency | `opsramp_servicedesk_urgency` | Urgency levels |
| ServiceDesk Business Impact | `opsramp_servicedesk_business_impact` | Business impact levels |

## Supported Data Sources

| Data Source | Type Name | Description |
|-------------|-----------|-------------|
| Resource Lookup | `opsramp_resource_lookup` | Query resources by OpsQL filter |
| Tenant | `opsramp_tenant` | Current tenant information |
| Role | `opsramp_role` | Look up a role by name |

## Installation

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    opsramp = {
      source  = "terraform.local/local/opsramp"
      version = "~> 0.1"
    }
  }
}
```

## Authentication

### Option 1: Environment Variables (Recommended)

```bash
export OPSRAMP_CLIENT_ID="your-client-id"
export OPSRAMP_CLIENT_SECRET="your-client-secret"
export OPSRAMP_TENANT="your-tenant-id"
export OPSRAMP_ENDPOINT="your-api-endpoint.opsramp.com"
```

### Option 2: Provider Block

```hcl
provider "opsramp" {
  client_id     = var.opsramp_client_id
  client_secret = var.opsramp_client_secret
  endpoint      = "your-api-endpoint.opsramp.com"
  tenant        = "your-tenant-id"
}
```

> **Note:** Environment variables take precedence unless explicitly overridden in the provider block.

## Usage Examples

### Managing a Resource

```hcl
resource "opsramp_resource" "server" {
  resource_name = "my-server"
  hostname      = "server01.example.com"
  resource_type = "Linux"
}
```

### Creating a Service Map

```hcl
resource "opsramp_servicemap" "app" {
  name = "My Application"
  type = "APPLICATION"

  threshold_type  = "count"
  threshold_limit = 1
}
```

### Multi-Tenant User Management

```hcl
resource "opsramp_client" "customer" {
  name      = "Acme Corp"
  country   = "United States"
  time_zone = "America/New_York"
}

resource "opsramp_user" "admin" {
  client     = opsramp_client.customer.unique_id
  login_name = "admin@acme.com"
  first_name = "Admin"
  last_name  = "User"
  email      = "admin@acme.com"
  password   = var.admin_password
  roles      = ["Client Full Permissions Admin"]
}
```

### Querying Resources

```hcl
data "opsramp_resource_lookup" "search" {
  query = "type = 'Linux' AND name LIKE 'prod%'"
}

output "resource_exists" {
  value = data.opsramp_resource_lookup.search.exists
}
```

For more complete examples, see the [`examples/`](examples/) directory.

## Development

### Prerequisites

- Go 1.24+
- Terraform 1.5+

### Build

```bash
go build -o terraform-provider-opsramp .
```

### Run Tests

```bash
go test -v ./...
```

### Run Acceptance Tests

Set required provider authentication variables (or define them in a repository-root `.env` file):

```bash
export OPSRAMP_CLIENT_ID="your-client-id"
export OPSRAMP_CLIENT_SECRET="your-client-secret"
export OPSRAMP_TENANT="your-tenant-id"
export OPSRAMP_ENDPOINT="your-api-endpoint.opsramp.com"
```

Then run acceptance tests:

```bash
make testacc
```

Optional scoped run examples:

```bash
make testacc TESTACC_PATH=./internal/resources TESTACC_RUN=TestAccDeviceGroupResource
```

### Install Locally

```bash
make install
```

## Project Structure

```
├── main.go                      # Provider entry point
├── internal/
│   ├── provider/
│   │   ├── provider.go          # Provider configuration and registration
│   │   └── version.go           # Version constant
│   ├── client/                  # API client layer
│   │   ├── common.go            # HTTP client, auth, request handling
│   │   ├── common_models.go     # Shared API models
│   │   ├── rate_limited_transport.go
│   │   ├── <resource>.go        # API methods per resource
│   │   └── <resource>_models.go # API request/response models
│   ├── resources/               # Terraform resource implementations
│   │   └── <resource>.go        # Schema, CRUD, model conversion
│   └── data/                    # Terraform data source implementations
│       └── <datasource>.go
├── docs/                        # Resource documentation from tfplugindocs
├── examples/                    # Example Terraform configurations
└── modules/                     # Reusable Terraform modules
```

### Architecture Conventions

- **Client layer** (`internal/client/`): Pure API communication. Each file contains CRUD methods and their corresponding request/response models. No Terraform types here.
- **Resource layer** (`internal/resources/`): Terraform schema definitions and CRUD handlers. Each resource uses `apiClient` as the field name for the API client and follows a consistent pattern:
  1. Interface compliance checks (required to embed standard configuration from BaseResource)
  ```
  var _ resource.Resource = &MyResource{}
  var _ resource.ResourceWithModifyPlan = &MyResource{}
  type MyResource struct {
	  BaseResource
  }
  ```
  3. Model struct with `tfsdk` tags
  4. Constructor (`NewMyResource()`)
  5. `Metadata`, `Schema`, `Configure` methods
  6. `Create`, `Read`, `Update`, `Delete` methods with consistent error handling
  7. `ImportState` where applicable
  8. Translation helpers (`buildMyResource`, `mapMyResource`)
- **Data source layer** (`internal/data/`): Read-only lookups following the same client/schema pattern.

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) for the Contributor License Agreement (CLA) requirement, coding standards, and pull request process before submitting changes.

## License

This project is property of HPE OpsRamp / Hewlett Packard Enterprise.

```
(C) Copyright 2026 Hewlett Packard Enterprise Development LP
```
