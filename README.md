# terraform-provider-anthropic

[![ci](https://github.com/sauterdigital/terraform-provider-anthropic/actions/workflows/ci.yml/badge.svg)](https://github.com/sauterdigital/terraform-provider-anthropic/actions/workflows/ci.yml)
[![license](https://img.shields.io/badge/license-MPL--2.0-blue.svg)](./LICENSE)

Terraform provider for the [Anthropic Admin API](https://platform.claude.com/docs/en/api/admin). Manages workspaces, API keys, organization and workspace members, invites, CMEK external keys, **service accounts** (OAuth Bearer), **federation issuers + rules** (workload identity federation), **spend limits**, **MCP tunnel certificates** (beta), and exposes the full analytics surface (token usage, cost, Claude Code, Skills, Connectors, Chat Projects, per-user breakdowns) as data sources for FinOps pipelines.

Covers **every documented Admin API endpoint**: 14 resources + 37 data sources spanning ~80 endpoints across 15 endpoint groups.

## Why this provider

Two existing options when this was written:

1. Use the Anthropic Console manually — not declarative, no drift detection, no audit trail.
2. The community provider [`terraform-mars/terraform-provider-anthropic`](https://github.com/terraform-mars/terraform-provider-anthropic), which covers workspaces, api_keys, workspace_members, and invites but not organization-level member management, CMEK, rate limits, or any usage / cost reports.

This provider's scope is **the full Admin API surface**, with first-class data sources for billing/observability signal so cost reports can flow into the same IaC pipeline that provisions workspaces.

## Quick start

```hcl
terraform {
  required_providers {
    anthropic = {
      source  = "sauterdigital/anthropic"
      version = "~> 0.1"
    }
  }
}

provider "anthropic" {
  # admin_api_key = "sk-ant-admin-..."   # or set ANTHROPIC_ADMIN_API_KEY
}

resource "anthropic_workspace" "platform" {
  name = "platform"
  tags = {
    env  = "prod"
    team = "platform"
  }
}

# Daily cost per workspace for the last 30 days — feed into your FinOps stack.
data "anthropic_cost_report" "monthly" {
  starting_at  = formatdate("YYYY-MM-DD'T'00:00:00'Z'", timeadd(timestamp(), "-720h"))
  bucket_width = "1d"
  group_by     = ["workspace_id"]
}
```

The admin key is distinct from regular Claude API keys — generate it in the Anthropic Console under organization settings.

## What's included

**14 resources**

Authenticated with `admin_api_key` (x-api-key):

| Resource | Notes |
|---|---|
| `anthropic_workspace` | Full CRUD. Tags mutable, `external_key_id` write-once, `data_residency` triggers replace on change. |
| `anthropic_api_key` | Update-only — the Admin API can't create keys. Supply an existing `id` and the provider manages name/status. |
| `anthropic_workspace_member` | Composite id `<workspace_id>:<user_id>`; role mutable. |
| `anthropic_invite` | Immutable after create — changes to email/role force replacement. |
| `anthropic_organization_member` | Set org role for an existing user (joined via accepted invite). |
| `anthropic_external_key` | CMEK config CRUD + validate, polymorphic across AWS / GCP / Azure. |
| `anthropic_spend_limit` | Per-user spend limit override (org/group/seat-tier limits stay in Console). |
| `anthropic_spend_limit_increase_decision` | Approve or deny a user's request to raise their cap. |

Require `oauth_token` (Bearer auth) — Admin API keys are rejected:

| Resource | Notes |
|---|---|
| `anthropic_service_account` | Named non-human identity for federation. `admin`-role creation needs interactive credential. |
| `anthropic_service_account_workspace` | Assigns an SA to a workspace with a role. |
| `anthropic_federation_issuer` | OIDC issuer registration (GitHub Actions, GitLab, etc). Polymorphic JWKS source. |
| `anthropic_federation_rule` | Workload identity federation rule binding OIDC claims to an SA. |
| `anthropic_federation_rule_workspace` | Extends a rule to an additional workspace. |
| `anthropic_tunnel_certificate` | MCP tunnel CA certificate (beta, `mcp-tunnels-2026-05-19` header added automatically). |

**37 data sources**

- Identity & membership: `anthropic_organization`, `anthropic_workspace[s]`, `anthropic_workspace_member[s]`, `anthropic_organization_member[s]`, `anthropic_invite[s]`
- Keys / CMEK: `anthropic_api_key[s]`, `anthropic_external_key[s]`
- Operational: `anthropic_organization_rate_limits`, `anthropic_workspace_rate_limits`
- FinOps reports: `anthropic_usage_report`, `anthropic_claude_code_usage_report`, `anthropic_cost_report`
- FinOps automation: `anthropic_effective_spend_limits`, `anthropic_spend_limit_increase_request[s]`
- Analytics v2 (Enterprise + `read:analytics` scope): `anthropic_activity_summaries`, `anthropic_token_usage_over_time`, `anthropic_per_user_token_usage`, `anthropic_cost_over_time`, `anthropic_per_user_cost`, `anthropic_user_activity`, `anthropic_skills_usage`, `anthropic_connectors_usage`, `anthropic_chat_projects_usage`
- Service accounts (Bearer auth): `anthropic_service_account[s]`, `anthropic_service_account_workspaces`, `anthropic_workspace_service_accounts`
- MCP Tunnels (Bearer + beta): `anthropic_tunnel[s]`, `anthropic_tunnel_certificates`

Full schema reference: [`docs/`](./docs).

## Configuration

| Argument | Env var | Description |
|---|---|---|
| `admin_api_key` | `ANTHROPIC_ADMIN_API_KEY` | Admin API key (`sk-ant-admin-...`). Used as `x-api-key` header. Required for most endpoints. |
| `oauth_token` | `ANTHROPIC_OAUTH_TOKEN` | OAuth Bearer token (user OAuth or WIF-minted SA token). **Required** for Service Accounts, Federation, and MCP Tunnels (which reject Admin API keys). When set, Bearer auth is used for ALL requests. |
| `base_url` | — | Optional. Defaults to `https://api.anthropic.com`. Override for staging or mock servers. |

At least one of `admin_api_key` or `oauth_token` must be set. When both are configured the client uses Bearer (the doc's modern preferred pattern). Every request sets `anthropic-version: 2023-06-01` and a provider-versioned `User-Agent`. HTTP 429 responses are retried with exponential backoff (capped at 30s), honoring `Retry-After` when present.

## Development

Requirements: Go 1.25, Terraform ≥ 1.0.

```bash
make build                       # compile the provider binary
make test                        # unit tests (fast, no API access)
make testacc                     # acceptance tests — requires TF_ACC=1 + ANTHROPIC_ADMIN_API_KEY, creates real workspaces
make fmt vet                     # gofmt + go vet
make docs                        # regenerate docs/ (requires tfplugindocs in PATH)

go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest   # one-time
```

To use a local build in a real config without publishing, add a `dev_overrides` block to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "sauterdigital/anthropic" = "/path/to/your/$GOPATH/bin"
  }
  direct {}
}
```

Then `make install` and run `terraform plan` against `examples/`.

## CI

- [`ci.yml`](./.github/workflows/ci.yml) runs on every push/PR: build, vet, gofmt, race-enabled unit tests, `terraform fmt` on examples, and a docs-drift check.
- [`acceptance.yml`](./.github/workflows/acceptance.yml) is `workflow_dispatch` only — acceptance tests mutate the real organization and incur API cost, so they never run automatically. Requires `ANTHROPIC_ADMIN_API_KEY` in repo secrets.
- [`release.yml`](./.github/workflows/release.yml) fires on `v*` tag push: builds signed, multi-arch artifacts via goreleaser, attaches them to a GitHub Release, and includes the `terraform-registry-manifest.json` the Terraform Registry needs to ingest the release. Requires `GPG_PRIVATE_KEY` and `PASSPHRASE` repo secrets.

## Publishing a release

1. Confirm `go test -race ./...` passes locally and `make docs` shows no diff.
2. Bump the version, commit, then tag: `git tag -a v0.X.Y -m "..."` and `git push origin v0.X.Y`.
3. The `release` workflow builds the binaries, signs the checksum file with the GPG key, and creates a draft-less GitHub Release.
4. First-time only: register the provider at https://registry.terraform.io/publish/provider — point it at this repo and upload the matching GPG public key. Subsequent releases are picked up automatically when the workflow finishes.

## Scope notes

- **MCP Tunnels** are deliberately out of scope. They are a beta surface that uses Bearer / WIF authentication instead of `x-api-key` and require an `anthropic-beta` header — supporting them would add a second auth code path to the client. Track as a future addition if there's demand.
- **Service Accounts** and **Audit Logs** are not part of the Admin API today (confirmed via sitemap, despite breadcrumbs that hint otherwise).

## License

[Mozilla Public License 2.0](./LICENSE) — the standard for HashiCorp-ecosystem Terraform providers.
