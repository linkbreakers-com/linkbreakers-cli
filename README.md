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

### Homebrew (macOS and Linux)

```bash
brew install linkbreakers-com/tap/linkbreakers
```

### npm (any OS with Node.js 18+)

```bash
npm install -g linkbreakers-cli
```

The npm package downloads the prebuilt binary for your platform from GitHub Releases and verifies it against the release `checksums.txt`. Note that the npm package `linkbreakers` is the TypeScript SDK; the CLI is `linkbreakers-cli`.

### Install script (macOS and Linux)

```bash
curl -fsSL https://cli.linkbreakers.com/install.sh | bash
```

The installer detects OS and architecture and installs the latest GitHub Release.

### Windows (PowerShell)

Use npm, or download the binary directly:

```powershell
$version = "<version>"
Invoke-WebRequest -Uri "https://github.com/linkbreakers-com/linkbreakers-cli/releases/download/v$version/linkbreakers-cli_${version}_windows_amd64.zip" -OutFile "linkbreakers.zip"
Expand-Archive -Path "linkbreakers.zip" -DestinationPath ".\\linkbreakers"
Move-Item ".\\linkbreakers\\linkbreakers.exe" "$HOME\\bin\\linkbreakers.exe"
```

Replace `<version>` with a real release like `1.140.0`. All archives are on the [Releases page](https://github.com/linkbreakers-com/linkbreakers-cli/releases).

## Updating

The CLI checks periodically for new releases and tells you how to update for the way it was installed:

- Homebrew: `brew upgrade linkbreakers`
- npm: `npm install -g linkbreakers-cli@latest`
- Install script or direct download: `linkbreakers self-update` (or rerun the installer)

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
6. GoReleaser publishes macOS, Linux, and Windows binaries to GitHub Releases and the Homebrew cask to `linkbreakers-com/homebrew-tap`.
7. The `npm/` wrapper is published to npm as `linkbreakers-cli` with the same version.

Publishing to Homebrew and npm needs the `HOMEBREW_TAP_SSH_KEY` (private half of a write deploy key on `linkbreakers-com/homebrew-tap`) and `NPM_TOKEN` repository secrets. Without them those two steps are skipped and the GitHub Release still ships.

## Local Development

```bash
make generate
make docs
make test
make build
```
