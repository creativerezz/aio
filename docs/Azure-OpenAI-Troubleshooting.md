# Azure OpenAI Troubleshooting

This document describes a known issue with Azure OpenAI integration and its fix.

## Issue: DeploymentNotFound Error (404)

### Symptoms

When using Fabric with Azure OpenAI, you may encounter this error:

```
POST "https://{resource}.cognitiveservices.azure.com/openai/chat/completions?api-version=...": 404 Not Found
{
  "code": "DeploymentNotFound",
  "message": "The API deployment for this resource does not exist..."
}
```

### Root Cause

Azure OpenAI requires deployment names in the URL path:

```
✅ Correct:  /openai/deployments/{deployment-name}/chat/completions
❌ Incorrect: /openai/chat/completions
```

The OpenAI Go SDK's `azure.WithEndpoint()` middleware has a bug in its URL transformation logic:

1. The SDK's `jsonRoutes` map expects paths like `/openai/chat/completions`
2. But the SDK actually sends paths like `/chat/completions` (without the `/openai/` prefix)
3. The `/openai/` prefix is included in the base URL, not the request path
4. This causes the route matching to **always fail**, so deployment names are never injected into the URL

### Technical Details

In the SDK's `azure/azure.go`:

```go
// SDK checks for these routes:
var jsonRoutes = map[string]bool{
    "/openai/chat/completions": true,  // Expects /openai/ prefix
    // ...
}

// But actual request path is:
path := "chat/completions"  // No /openai/ prefix!
```

The mismatch means `jsonRoutes[req.URL.Path]` never matches, and the deployment name transformation never happens.

## Fix

The fix in `internal/plugins/ai/azure/azure.go` adds custom middleware that:

1. Intercepts outgoing requests
2. Extracts the deployment name from the request body's `model` field
3. Transforms the URL path to include `/deployments/{name}/`

```go
// Transform: /chat/completions -> /openai/deployments/{name}/chat/completions
func azureDeploymentMiddleware(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
    // Routes that need deployment name injection
    deploymentRoutes := map[string]bool{
        "/chat/completions":     true,
        "/completions":          true,
        "/embeddings":           true,
        "/audio/speech":         true,
        "/audio/transcriptions": true,
        "/audio/translations":   true,
        "/images/generations":   true,
    }

    // Extract deployment from body and transform URL...
}
```

## Additional Fix: StreamOptions Error

### Symptom

```
400 Bad Request
{
  "message": "The 'stream_options' parameter is only allowed when 'stream' is enabled."
}
```

### Cause

The Chat Completions API was sending `stream_options` for all requests, but Azure only accepts this parameter when `stream: true` is also set.

### Fix

Moved `StreamOptions` to only be set for streaming requests in `internal/plugins/ai/openai/chat_completions.go`.

## Configuration

Ensure your Azure OpenAI configuration is correct:

```bash
# In ~/.config/aio/.env
AZURE_API_KEY=your-api-key
AZURE_API_BASE_URL=https://{your-resource}.cognitiveservices.azure.com/
AZURE_DEPLOYMENTS=your-deployment-1,your-deployment-2  # Comma-separated deployment names
AZURE_API_VERSION=2024-12-01-preview  # Optional, defaults to 2024-05-01-preview
```

**Note:** The deployment name is what you specified when deploying a model in Azure AI Foundry (formerly Azure OpenAI Studio), not the model name itself (e.g., `my-gpt4-deployment` rather than `gpt-4`).

## Verification

Test your Azure OpenAI setup:

```bash
fabric --model <your-deployment-name> --pattern summarize "Hello world"
```

Replace `<your-deployment-name>` with the actual deployment name from your Azure configuration.

You should see a successful response from your Azure OpenAI deployment.

## References

- GitHub Issue: [#1954](https://github.com/creativerezz/aio/issues/1954)
- Pull Request: [#1965](https://github.com/creativerezz/aio/pull/1965)
- [Azure OpenAI REST API Reference](https://learn.microsoft.com/en-us/azure/ai-services/openai/reference)
