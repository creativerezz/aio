# Aio CLI Guide

Complete reference for the `aio` command-line interface.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Input Methods](#input-methods)
- [Patterns](#patterns)
- [Models & Vendors](#models--vendors)
- [Sessions & Contexts](#sessions--contexts)
- [Streaming & Output](#streaming--output)
- [Model Parameters](#model-parameters)
- [Thinking / Reasoning](#thinking--reasoning)
- [Web Tools](#web-tools)
- [Image Generation](#image-generation)
- [Audio & Transcription](#audio--transcription)
- [Web Search](#web-search)
- [Server Mode](#server-mode)
- [Configuration](#configuration)
- [Notifications](#notifications)
- [Shell Completions](#shell-completions)
- [Debugging](#debugging)
- [Flag Reference](#flag-reference)

---

## Quick Start

```bash
# Ask a question
aio "What is quantum entanglement?"

# Use a pattern
echo "long article text" | aio -p summarize

# Stream a response
aio -s -p explain "explain blockchain to a 5 year old"

# Pick a specific model
aio -m gpt-4o -s "write a haiku about code"

# First-time setup
aio --setup
```

---

## Input Methods

Aio accepts input from multiple sources, combined in order:

```bash
# Positional argument
aio "your question here"

# Piped stdin
echo "some text" | aio

# Both combined (piped content + argument)
echo "here is an article" | aio -p summarize "focus on the key takeaways"

# File via stdin
aio -p summarize < document.txt

# Attachment (image URL or file path)
aio -a https://example.com/image.png "describe this image"
aio -a ./screenshot.png "what's in this image?"
```

When both piped input and a positional argument are provided, they are combined into a single user message.

---

## Patterns

Patterns are reusable prompt templates stored in `~/.config/aio/patterns/`. Each pattern is a directory containing a `system.md` file.

```bash
# Use a pattern
aio -p summarize "text to summarize"

# List all patterns
aio -l

# List newest patterns
aio -n 10

# Update patterns from remote
aio -U
```

### Pattern Variables

Patterns can contain `{{variable}}` placeholders. Pass values with `-v`:

```bash
aio -p explain -v=#role:expert -v=#audience:beginners "explain quantum computing"
```

| Flag | Description |
|------|-------------|
| `-v=#name:value` | Set a pattern variable |
| `--input-has-vars` | Also apply variable replacement to user input |
| `--no-variable-replacement` | Disable all variable replacement |

### Per-Pattern Model Override

Override the model for specific patterns using environment variables:

```bash
# Format: AIO_MODEL_<PATTERN_NAME_UPPERCASE>=model
export AIO_MODEL_SUMMARIZE="gpt-4o"

# With vendor: AIO_MODEL_<NAME>=vendor|model
export AIO_MODEL_CODE_REVIEW="Anthropic|claude-sonnet-4-5-20250929"

# Hyphens become underscores
# Pattern "code-review" -> AIO_MODEL_CODE_REVIEW
```

---

## Models & Vendors

```bash
# Specify a model
aio -m gpt-4o "your prompt"

# Specify vendor + model
aio -V "LM Studio" -m "openai/gpt-oss-70b" "your prompt"

# List all models
aio -L

# List all vendors
aio --listvendors

# Change default model interactively
aio -d
```

### Strategies

Strategies modify the system message with special prompting techniques:

```bash
aio --strategy chain-of-thought -p analyze "complex problem"
aio --liststrategies
```

---

## Sessions & Contexts

### Sessions

Sessions maintain conversation history across invocations:

```bash
# Start or continue a session
aio --session my-chat "hello, let's discuss AI"

# Follow up in the same session
aio --session my-chat "what about the risks?"

# List sessions
aio -X

# Print session history
aio --printsession my-chat

# Save full session to file
aio --session my-chat --output-session -o conversation.md "final question"

# Delete a session
aio -W my-chat
```

### Contexts

Contexts are named system-level instructions loaded before your pattern:

```bash
# Use a context
aio -C business -p explain "AI trends"

# List contexts
aio -x

# Print context contents
aio --printcontext business

# Delete a context
aio -w business
```

When both a context and pattern are specified, the context is prepended to the pattern's system message.

---

## Streaming & Output

```bash
# Stream response in real-time
aio -s "tell me a story"

# Copy result to clipboard
aio -c -p summarize "text"

# Write to file
aio -o output.md -p summarize "text"

# Dry run - show what would be sent without calling the model
aio --dry-run -p explain "test input"
```

| Flag | Short | Description |
|------|-------|-------------|
| `--stream` | `-s` | Stream response in real-time |
| `--copy` | `-c` | Copy result to clipboard |
| `--output` | `-o` | Write result to file |
| `--output-session` | | Write full session history (not just last response) |
| `--dry-run` | | Preview what would be sent to the model |
| `--show-metadata` | | Print token usage to stderr |

---

## Model Parameters

Fine-tune model behavior:

```bash
aio -t 0.2 "precise answer please"          # Low temperature
aio -t 1.5 -s "creative story"              # High temperature
aio --seed 42 "reproducible output"          # Fixed seed
aio --raw "use model defaults"               # Skip sending parameters
```

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--temperature` | `-t` | float | 0.7 | Randomness (0.0-2.0) |
| `--topp` | `-T` | float | 0.9 | Nucleus sampling (0.0-1.0) |
| `--presencepenalty` | `-P` | float | 0.0 | Presence penalty (-2.0 to 2.0) |
| `--frequencypenalty` | `-F` | float | 0.0 | Frequency penalty (-2.0 to 2.0) |
| `--seed` | `-e` | int | 0 | Seed for reproducibility |
| `--raw` | `-r` | bool | false | Use model defaults, skip parameter sending |

---

## Thinking / Reasoning

Enable extended thinking for models that support it (Anthropic, Gemini):

```bash
# Set thinking level
aio --thinking high "solve this math problem"

# Hide thinking output
aio --thinking medium --suppress-think "just show the answer"

# Custom thinking tags
aio --thinking high --think-start-tag "<reasoning>" --think-end-tag "</reasoning>" "problem"
```

| Level | Token Budget |
|-------|-------------|
| `off` | Disabled |
| `low` | 1,024 tokens |
| `medium` | 2,048 tokens |
| `high` | 4,096 tokens |

---

## Web Tools

### YouTube

```bash
# Get transcript (default)
aio -y "https://youtube.com/watch?v=abc123"

# Transcript with timestamps
aio -y "https://youtube.com/watch?v=abc123" --transcript-with-timestamps

# Get comments
aio -y "https://youtube.com/watch?v=abc123" --comments

# Get metadata as JSON
aio -y "https://youtube.com/watch?v=abc123" --metadata

# Process a playlist
aio -y "https://youtube.com/playlist?list=abc" --playlist

# Pipe transcript into a pattern
aio -y "https://youtube.com/watch?v=abc123" -p summarize

# Custom yt-dlp args (e.g., cookies)
aio -y "https://youtube.com/watch?v=abc123" --yt-dlp-args='--cookies-from-browser brave'
```

### Web Scraping

```bash
# Scrape URL to markdown via Jina AI
aio -u "https://example.com/article" -p summarize

# Search question via Jina AI
aio -q "what is quantum computing"
```

### Spotify

```bash
# Get podcast/episode metadata
aio --spotify "https://open.spotify.com/episode/xyz"
```

### HTML Readability

```bash
# Convert HTML input to clean markdown
cat page.html | aio --readability
```

---

## Image Generation

Generate images with models like DALL-E:

```bash
aio -m dall-e-3 --image-file sunset.png "a sunset over the ocean"

# With options
aio -m dall-e-3 \
  --image-file photo.png \
  --image-size 1536x1024 \
  --image-quality high \
  "photorealistic mountain landscape"

# Transparent background (PNG/WebP only)
aio -m dall-e-3 --image-file logo.png --image-background transparent "minimalist logo"

# JPEG with compression
aio -m dall-e-3 --image-file photo.jpg --image-compression 80 "portrait photo"
```

| Flag | Default | Options |
|------|---------|---------|
| `--image-file` | | Output path (`.png`, `.jpg`, `.jpeg`, `.webp`) |
| `--image-size` | `auto` | `1024x1024`, `1536x1024`, `1024x1536`, `auto` |
| `--image-quality` | `auto` | `low`, `medium`, `high`, `auto` |
| `--image-compression` | 0 | 0-100 (JPEG/WebP only) |
| `--image-background` | `opaque` | `opaque`, `transparent` (PNG/WebP only) |

---

## Audio & Transcription

### Transcribe Audio/Video

```bash
# Basic transcription
aio --transcribe-file recording.mp3 --transcribe-model whisper-1

# Transcribe then process with a pattern
aio --transcribe-file meeting.mp3 --transcribe-model whisper-1 -p summarize

# Split large files (>25MB) automatically
aio --transcribe-file long-video.mp4 --transcribe-model whisper-1 --split-media-file

# List available transcription models
aio --list-transcription-models
```

### Text-to-Speech

```bash
# Generate speech
aio -m tts-1 -o speech.wav "Hello, this is a test."

# Choose a voice
aio -m tts-1 --voice Puck -o greeting.mp3 "Welcome to Aio"

# List available Gemini voices
aio --list-gemini-voices
```

---

## Web Search

Enable web search for supported models (Anthropic, OpenAI, Gemini):

```bash
# Enable search
aio --search "what happened in tech news today?"

# With location context
aio --search --search-location "America/Los_Angeles" "local weather forecast"
```

---

## Server Mode

Run Aio as a REST API server:

```bash
# Start server on default port
aio --serve

# Custom address and API key
aio --serve --address :9000 --api-key my-secret-key

# With Ollama-compatible endpoints
aio --serveOllama --address :8000
```

The server exposes:
- `POST /chat` - Stream chat completions (SSE)
- `GET /models/names` - List available models
- `GET /patterns/{name}` - Get a pattern
- `POST /patterns/{name}/apply` - Apply pattern with variables
- `POST /youtube/transcript` - Get YouTube transcript
- Swagger docs at `/swagger/index.html`

Authentication via `X-API-Key` header when `--api-key` is set.

---

## Configuration

### Config File

Aio supports a YAML config file for default settings:

**Location:** `~/.config/aio/config.yaml` (override with `--config`)

```yaml
model: gpt-4o-mini
vendor: OpenAI
temperature: 0.5
stream: true
pattern: summarize
notification: true
voice: Kore
debug: 0
```

CLI flags always override config file values.

### Setup

```bash
# Interactive setup (API keys, default model, patterns)
aio --setup

# Update patterns from remote
aio --updatepatterns
```

### Language

```bash
# Set output language
aio -g zh "explain AI"    # Chinese
aio -g es "explain AI"    # Spanish
```

Auto-detected from `LC_ALL`, `LC_MESSAGES`, or `LANG` environment variables when not specified.

### Storage Locations

All user data is stored in `~/.config/aio/`:

| Directory | Contents |
|-----------|----------|
| `.env` | API keys and settings |
| `patterns/` | Downloaded prompt patterns |
| `sessions/` | Chat session history |
| `contexts/` | Named contexts |
| `strategies/` | Prompt strategies |

---

## Notifications

Get notified when long-running commands complete:

```bash
# Built-in desktop notification
aio --notification -p analyze < large-document.txt

# Custom notification command
aio --notification-command='notify-send "$1" "$2"' -p analyze < doc.txt
```

---

## Shell Completions

Tab completion is available for Zsh, Bash, and Fish:

```bash
# Auto-install completions
curl -fsSL https://raw.githubusercontent.com/creativerezz/aio/main/completions/setup-completions.sh | sh

# Dry run first
curl -fsSL https://raw.githubusercontent.com/creativerezz/aio/main/completions/setup-completions.sh | sh -s -- --dry-run
```

See [Shell Completions](Shell-Completions.md) for detailed setup instructions.

---

## Debugging

```bash
# Basic debug output
aio --debug 1 "test"

# Detailed debug (includes config resolution)
aio --debug 2 "test"

# Trace-level (full request/response details)
aio --debug 3 "test"

# Show token usage
aio -s --show-metadata "test"
# Output: [Metadata] Input: 42 | Output: 128 | Total: 170
```

---

## Extensions

Register external tools as extensions:

```bash
# Add an extension
aio --addextension /path/to/extension-config.yaml

# List registered extensions
aio --listextensions

# Remove an extension
aio --rmextension my-extension
```

---

## Flag Reference

### Pattern & Input

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--pattern` | `-p` | string | Use a pattern |
| `--variable` | `-v` | key:value | Set pattern variable |
| `--context` | `-C` | string | Load a named context |
| `--session` | | string | Use/continue a session |
| `--strategy` | | string | Apply a prompt strategy |
| `--attachment` | `-a` | string[] | Attach file or URL |
| `--no-variable-replacement` | | bool | Disable variable substitution |
| `--input-has-vars` | | bool | Apply variables to user input too |
| `--language` | `-g` | string | Output language code |

### Model

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--model` | `-m` | string | | Model name |
| `--vendor` | `-V` | string | | Vendor name |
| `--temperature` | `-t` | float | 0.7 | Randomness |
| `--topp` | `-T` | float | 0.9 | Nucleus sampling |
| `--presencepenalty` | `-P` | float | 0.0 | Presence penalty |
| `--frequencypenalty` | `-F` | float | 0.0 | Frequency penalty |
| `--seed` | `-e` | int | 0 | Reproducibility seed |
| `--raw` | `-r` | bool | false | Use model defaults |
| `--modelContextLength` | | int | 0 | Context length (ollama) |

### Output

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--stream` | `-s` | bool | Stream response |
| `--copy` | `-c` | bool | Copy to clipboard |
| `--output` | `-o` | string | Write to file |
| `--output-session` | | bool | Output full session |
| `--dry-run` | | bool | Preview without sending |
| `--show-metadata` | | bool | Show token usage |

### Thinking

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--thinking` | string | | Level: `off`, `low`, `medium`, `high` |
| `--suppress-think` | bool | false | Hide thinking output |
| `--think-start-tag` | string | `<think>` | Custom start tag |
| `--think-end-tag` | string | `</think>` | Custom end tag |

### Listing

| Flag | Short | Description |
|------|-------|-------------|
| `--listpatterns` | `-l` | List patterns |
| `--listmodels` | `-L` | List models |
| `--listcontexts` | `-x` | List contexts |
| `--listsessions` | `-X` | List sessions |
| `--liststrategies` | | List strategies |
| `--listvendors` | | List vendors |
| `--listextensions` | | List extensions |
| `--latest` | `-n` | Show N newest patterns |

### YouTube

| Flag | Short | Description |
|------|-------|-------------|
| `--youtube` | `-y` | YouTube URL |
| `--transcript` | | Get transcript (default) |
| `--transcript-with-timestamps` | | Transcript with timestamps |
| `--comments` | | Get comments |
| `--metadata` | | Output metadata as JSON |
| `--playlist` | | Prefer playlist mode |
| `--yt-dlp-args` | | Extra yt-dlp arguments |

### Web & Scraping

| Flag | Short | Description |
|------|-------|-------------|
| `--scrape_url` | `-u` | Scrape URL via Jina AI |
| `--scrape_question` | `-q` | Search via Jina AI |
| `--spotify` | | Spotify podcast/episode URL |
| `--readability` | | Convert HTML to clean markdown |
| `--search` | | Enable web search |
| `--search-location` | | Search location context |

### Image

| Flag | Description |
|------|-------------|
| `--image-file` | Output path (.png, .jpg, .webp) |
| `--image-size` | Dimensions (default: `auto`) |
| `--image-quality` | Quality level (default: `auto`) |
| `--image-compression` | JPEG/WebP compression 0-100 |
| `--image-background` | `opaque` or `transparent` |

### Audio

| Flag | Description |
|------|-------------|
| `--transcribe-file` | File to transcribe |
| `--transcribe-model` | Transcription model |
| `--split-media-file` | Split files >25MB |
| `--voice` | TTS voice (default: `Kore`) |
| `--list-gemini-voices` | List Gemini voices |
| `--list-transcription-models` | List transcription models |

### Server

| Flag | Default | Description |
|------|---------|-------------|
| `--serve` | | Start REST API |
| `--serveOllama` | | Start with Ollama endpoints |
| `--address` | `:8080` | Bind address |
| `--api-key` | | API key for auth |

### Management

| Flag | Short | Description |
|------|-------|-------------|
| `--setup` | `-S` | Interactive setup |
| `--updatepatterns` | `-U` | Update patterns |
| `--changeDefaultModel` | `-d` | Change default model |
| `--wipecontext` | `-w` | Delete a context |
| `--wipesession` | `-W` | Delete a session |
| `--printcontext` | | Print context contents |
| `--printsession` | | Print session contents |
| `--config` | | Path to config YAML |

### Other

| Flag | Description |
|------|-------------|
| `--notification` | Desktop notification on completion |
| `--notification-command` | Custom notification command |
| `--addextension` | Register extension |
| `--rmextension` | Remove extension |
| `--debug` | Debug level (0-3) |
| `--version` | Print version |
| `--disable-responses-api` | Disable OpenAI Responses API |
