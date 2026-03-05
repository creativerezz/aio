# Changelog

## v0.0.5 (2026-03-05)

### Direct commits

- Fix: create empty .env file in Docker image
The app fails to start if ~/.config/aio/.env doesn't exist
because godotenv.Load() returns an error for missing files.
Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>

## v0.0.4 (2026-03-05)

### Direct commits

- Fix: add Dockerfile and fix Railway deployment config

- Add Dockerfile with multi-stage build (no CGO dependency)
- Bake patterns and strategies into the image

- Fix railway.toml to point to correct Dockerfile path
- Fix healthcheck path to /patterns/names and increase timeout
Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>

## v0.0.3 (2026-02-08)

### Direct commits

- Feat: add Railway deployment configuration
Bake patterns and strategies into Docker image, add railway.toml config,
and add GitHub Actions workflow to auto-deploy on release.
Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>

## v0.0.2 (2026-02-08)

### Direct commits

- Refactor: replace all remaining fabric references with aio

- Config paths: .config/fabric → .config/aio (runtime bug fix)
- Strategies paths: .config/fabric/strategies → .config/aio/strategies

- i18n locale dir: fabric/locales → aio/locales
- Shell completions: remove fabric/fabric-ai registrations, fix function name

- Web env vars: FABRIC_BASE_URL → AIO_BASE_URL, __FABRIC_CONFIG__ → __AIO_CONFIG__
- CLI: update flag descriptions, error messages, help text, example config

- i18n: update all 10 locale files (command examples and product name)
- Server: rename FabricResponseFormat, buildFabricChatURL, fabricDb vars

- Tests: update temp dir prefixes and test config paths
- Rename scripts/setup_fabric.bat → scripts/setup_aio.bat

- Update Nix descriptions and GitHub issue template URL
Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>

All notable changes to **aio** will be documented in this file.

_Forked from [creativerezz/aio](https://github.com/creativerezz/aio)._
