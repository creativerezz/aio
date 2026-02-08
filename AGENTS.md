# AGENTS.md - Aio Codebase Context

## Project Overview

Aio is an open-source CLI tool and REST API for augmenting humans using AI. It organizes AI prompts into reusable "patterns" and provides a unified interface across 15+ LLM vendors. Written in Go 1.25.1, licensed MIT. Forked from [creativerezz/aio](https://github.com/creativerezz/aio).

**Repository:** `github.com/creativerezz/aio`

## Architecture

```
cmd/aio/                 Main CLI entry point (thin - delegates to internal/cli)
cmd/code2context/        Helper: scans code directories for AI context
cmd/generate_changelog/  Changelog generation from git history + GitHub PRs
cmd/to_pdf/              PDF conversion utility

internal/
  cli/                   CLI flags (jessevdk/go-flags), command routing, handler chain
  core/                  PluginRegistry (central hub) + Chatter (chat orchestration)
  domain/                Data models: ChatRequest, ChatOptions, StreamUpdate
  chat/                  ChatCompletionMessage, ChatMessagePart, ToolCall types
  plugins/
    ai/                  Vendor interface + 15+ implementations:
                           anthropic/, openai/, gemini/, ollama/, azure/, bedrock/,
                           perplexity/, copilot/, digitalocean/, vertexai/, exolab/,
                           lmstudio/, openai_compatible/, gemini_openai/, dryrun/
    db/fsdb/             Filesystem DB: patterns, sessions, contexts (no external DB)
    template/            Template variable substitution ({{var}} syntax)
    strategy/            Named prompt strategies (prepended to system messages)
  server/                REST API (Gin framework, Swagger docs at /swagger/index.html)
  tools/                 Integrations: YouTube, Spotify, Jina, patterns loader, git helper
  i18n/                  10+ languages with BCP 47 locale support
  log/                   Debug logging
  util/                  OAuth storage, helpers

data/
  patterns/              ~240+ AI prompt patterns (each = directory with system.md)
  strategies/            Processing strategies (JSON files)

web/                     Web UI frontend
docs/                    Documentation, Swagger files, CONTRIBUTING.md
scripts/                 Docker, installer, Python UI utilities
```

## Key Data Flow

1. User runs `aio -p pattern_name < input`
2. `internal/cli` parses flags, builds `domain.ChatRequest`
3. `core.Chatter.BuildSession()` loads pattern from fsdb, applies template variables, context, strategy, language
4. `core.PluginRegistry` routes to configured AI vendor
5. Vendor (e.g., `plugins/ai/openai/`) makes API call, streams response

## Core Types

- **`domain.ChatRequest`** - Pattern name, context, session, variables, message, language, strategy
- **`domain.ChatOptions`** - Model, temperature, top_p, max_tokens, thinking level, image/audio options
- **`chat.ChatCompletionMessage`** - Role + content (supports multi-part: text, images)
- **`ai.Vendor` interface** - `ListModels()`, `Send()`, `SendStream()`, `NeedsRawMode()`
- **`core.PluginRegistry`** - Central registry for all vendors, tools, DB, patterns
- **`core.Chatter`** - Chat orchestration: session building, vendor dispatch, streaming
- **`fsdb.Db`** - File-based storage: patterns, sessions, contexts, .env config

## Patterns System

Patterns live in `data/patterns/<name>/system.md` with sections:
- `# IDENTITY` / `# IDENTITY and PURPOSE`
- `# STEPS`
- `# OUTPUT` / `# OUTPUT INSTRUCTIONS`
- `# EXAMPLE`
- `# INPUT`

Template variables: `{{input}}` (user input), `{{variable_name}}` (custom via `-v`)
Custom patterns: user-configurable directory via `CUSTOM_PATTERNS_DIRECTORY` env var.

## Plugin/Vendor System

All AI vendors implement the `Vendor` interface (`internal/plugins/ai/vendor.go`).
New vendors extend `plugins.PluginBase` for setup questions and env var handling.
OpenAI-compatible vendors reuse `openai.NewClientCompatible()`.

**Vendor registration** happens in `core.NewPluginRegistry()` (`internal/core/plugin_registry.go`).

## REST API

Server entry: `internal/server/serve.go` using Gin framework.
Handlers: chat, patterns, contexts, sessions, models, strategies, YouTube, config.
Auth: optional API key via `X-API-Key` header.
Swagger annotations required for new endpoints (see `docs/CONTRIBUTING.md`).

## Build & Test

```bash
go build -o aio ./cmd/aio           # Build
go test ./...                        # Run all tests (57 test files)
go test -cover ./...                 # With coverage
go test -tags=integration ./...      # Integration tests (require credentials)
```

Tests use `testify` for assertions, table-driven patterns, and mock vendors.
No Makefile - builds automated via GoReleaser (`.goreleaser.yaml`).

## CI/CD

- **ci.yml** - Tests + modernize check on push/PR to main
- **release.yml** - GoReleaser builds on version tags
- **update-version-and-create-tag.yml** - Auto-bumps patch version on merge to main
- **patterns.yaml** - Zips patterns on changes to `data/patterns/`

## Configuration

User config stored in `~/.config/aio/`:
- `.env` - API keys and settings
- `patterns/` - Downloaded patterns
- `sessions/` - Chat session history (JSON)
- `contexts/` - Named contexts
- `strategies/` - Prompt strategies

Per-pattern model override: `AIO_MODEL_{PATTERN_NAME}=vendor|model`

## Conventions

- Standard Go conventions: `gofmt`, `golint`
- Commit messages: `feat:`, `fix:`, `docs:`, `chore:`
- PRs require changelog entry: `go run ./cmd/generate_changelog --incoming-pr PR_NUMBER`
- Keep PRs focused (<50 files)
- Vendor plugins follow setup question pattern with env var prefixes
- i18n keys used for all user-facing strings in setup flows
