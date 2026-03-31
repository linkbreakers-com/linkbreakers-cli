# Linkbreakers CLI

Official command line interface for the Linkbreakers API.

The binary name is `linkbreakers`, with commands shaped for direct usage by humans, scripts, and LLMs:

```bash
linkbreakers auth set-token --token <api-token>
linkbreakers links list --page-size 20
linkbreakers links create --destination https://example.com --name "Launch link"
linkbreakers custom-domains check <domain-id>
linkbreakers raw GET /v1/links?pageSize=5
```

## Installation

The easiest install path for macOS and Linux is:

```bash
curl -fsSL https://cli.linkbreakers.com/install.sh | bash
```

That installer:

- detects OS and architecture automatically
- downloads the latest GitHub Release
- installs `linkbreakers` without requiring you to pick a version manually

Direct downloads remain available from GitHub Releases:

- Repository: `linkbreakers-com/linkbreakers-cli`
- Releases: `https://github.com/linkbreakers-com/linkbreakers-cli/releases`

No extra package registry is required. GitHub Releases remains the canonical binary distribution.

### Quick Install

#### macOS and Linux

```bash
curl -fsSL https://cli.linkbreakers.com/install.sh | bash
```

#### Windows (PowerShell)

```powershell
$version = "<version>"
Invoke-WebRequest -Uri "https://github.com/linkbreakers-com/linkbreakers-cli/releases/download/v$version/linkbreakers-cli_$version_windows_amd64.zip" -OutFile "linkbreakers.zip"
Expand-Archive -Path "linkbreakers.zip" -DestinationPath ".\\linkbreakers"
Move-Item ".\\linkbreakers\\linkbreakers.exe" "$HOME\\bin\\linkbreakers.exe"
```

Replace `<version>` with a real release like `1.42.8`, or download the right archive from the Releases page directly.

## Updating

The CLI checks periodically for new releases and will let users know when an update is available.

To update on supported platforms:

```bash
linkbreakers self-update
```

You can also rerun the installer:

```bash
curl -fsSL https://cli.linkbreakers.com/install.sh | bash
```

## Authentication

Use either:

- `LINKBREAKERS_TOKEN`
- `linkbreakers auth set-token --token <api-token>`

Optional overrides:

- `LINKBREAKERS_BASE_URL`
- `LINKBREAKERS_OUTPUT=json|table`

## Commands

First-class commands currently included:

- `linkbreakers links ...`
- `linkbreakers directories ...`
- `linkbreakers custom-domains ...`
- `linkbreakers raw METHOD PATH`
- `linkbreakers completion ...`
- `linkbreakers self-update`
- `linkbreakers version`

The `raw` command is the fallback for any endpoint that does not yet have a dedicated subcommand.

## Docs for LLMs

This repo includes:

- `linkbreakers help`
- per-command markdown docs in `docs/commands/`
- `llms.txt` at repo root

To regenerate docs after CLI changes:

```bash
go run ./cmd/linkbreakers gendocs
```

## Releases

Releases are automated through GitHub Actions:

1. The API repo dispatches `update-sdk`.
2. This repo fetches the latest Swagger version.
3. The internal Go client is regenerated from the OpenAPI spec.
4. Command docs are regenerated.
5. A git tag is created.
6. GoReleaser publishes macOS, Linux, and Windows binaries to GitHub Releases.

## Local Development

```bash
make generate
make docs
make test
make build
```
