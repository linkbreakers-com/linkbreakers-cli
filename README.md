# Linkbreakers CLI

Official command line interface for the Linkbreakers API.

The binary name is `linkbreakers`. Every API operation is a command, generated from the API's OpenAPI document, so the CLI always matches the API: when the API gains an operation, the next release has the command.

```bash
linkbreakers signup --email you@example.com          # a workspace and its token, in one step
linkbreakers links create --destination https://example.com --name "Launch link"
linkbreakers links list --page-size 20
linkbreakers analytics-data get-top-links
```

## Quick start for AI agents

No account and nothing to install:

```bash
npx -y linkbreakers-cli signup --email you@example.com
npx -y linkbreakers-cli links create --destination https://example.com --name "My first link"
```

`signup` creates a free workspace, saves its API token, and emails the person a link and a 6-digit code to claim it. Until it is claimed, the token can create and manage links, QR codes, workflows and page themes and read analytics; an unclaimed workspace is deleted with everything in it 24 hours after signup. Every command run in the meantime reminds you on stderr, and once the person claims it the saved token is upgraded automatically. `linkbreakers claim --wait` waits for the claim.

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

- No account yet: `linkbreakers signup --email <email>`.
- Existing account: create a token at https://app.linkbreakers.com/workspace/dashboard/api-tokens, then `linkbreakers auth set-token --token <api-token>` or set `LINKBREAKERS_TOKEN`.

Optional overrides:

- `LINKBREAKERS_BASE_URL`
- `LINKBREAKERS_AGENT_BASE_URL` (where `signup` registers, `https://mcp.linkbreakers.com` by default)

## Commands

`linkbreakers <service> <operation>`, one per API operation: `linkbreakers links create`, `linkbreakers workflow-steps update`, `linkbreakers analytics-data get-breakdown`, and so on. `--help` on any of them lists the flags, read from the request's fields.

- Scalar fields are flags: `--destination`, `--qrcode-enabled=false`, `--page-size 20`.
- List fields repeat: `--tags spring --tags print`.
- Structured fields take JSON: `--metadata '{"campaign":"spring"}'`.
- `--data` takes the whole body as JSON, `@file.json`, or `-` for stdin; flags set fields on top of it.

Output is JSON. An API error prints its code and message, which say what to do next.

Written by hand, because they are not API operations: `signup`, `claim`, `auth set-token|status|clear`, `completion`, `self-update`, `version`.

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
3. The commands are regenerated from the OpenAPI spec (`scripts/generate-commands.sh`, which runs `internal/gen/commands`).
4. Command docs are regenerated.
5. A git tag is created.
6. GoReleaser publishes macOS, Linux, and Windows binaries to GitHub Releases and the Homebrew cask to `linkbreakers-com/homebrew-tap`.
7. The [`npm/`](./npm) wrapper is published to npm as `linkbreakers-cli` with the same version, through npm trusted publishing (no token). A failure there does not block the GitHub release; it opens a `release-failure` issue.

Publishing to Homebrew needs the `HOMEBREW_TAP_SSH_KEY` repository secret (private half of a write deploy key on `linkbreakers-com/homebrew-tap`).

## Local Development

```bash
make generate
make docs
make test
make build
```
