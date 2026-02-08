package vertexai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	debuglog "github.com/creativerezz/aio/internal/log"
)

const (
	// API limits
	maxResponseSize    = 10 * 1024 * 1024 // 10MB
	errorResponseLimit = 1024             // 1KB for error messages

	// Default region for Model Garden API (global doesn't work for this endpoint)
	defaultModelGardenRegion = "us-central1"
)

// Supported Model Garden publishers (others can be added when SDK support is implemented)
var publishers = []string{"google", "anthropic"}

// publisherModelsResponse represents the API response from publishers.models.list
type publisherModelsResponse struct {
	PublisherModels []publisherModel `json:"publisherModels"`
	NextPageToken   string           `json:"nextPageToken"`
}

// publisherModel represents a single model in the API response
type publisherModel struct {
	Name string `json:"name"` // Format: publishers/{publisher}/models/{model}
}

// fetchModelsPage makes a single API request and returns the parsed response.
// Extracted to ensure proper cleanup of HTTP response bodies in pagination loops.
func fetchModelsPage(ctx context.Context, httpClient *http.Client, url, projectID, publisher string) (*publisherModelsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	// Set quota project header required by Vertex AI API
	req.Header.Set("x-goog-user-project", projectID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, errorResponseLimit))
		debuglog.Debug(debuglog.Basic, "API error for %s: status %d, url: %s, body: %s\n", publisher, resp.StatusCode, url, string(bodyBytes))
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if len(bodyBytes) > maxResponseSize {
		return nil, fmt.Errorf("response too large (>%d bytes)", maxResponseSize)
	}

	var response publisherModelsResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// listPublisherModels fetches models from a specific publisher via the Model Garden API
func listPublisherModels(ctx context.Context, httpClient *http.Client, region, projectID, publisher string) ([]string, error) {
	// Use default region if global or empty (Model Garden API requires a specific region)
	if region == "" || region == "global" {
		region = defaultModelGardenRegion
	}

	baseURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1beta1/publishers/%s/models", region, publisher)

	var allModels []string
	pageToken := ""

	for {
		url := baseURL
		if pageToken != "" {
			url = fmt.Sprintf("%s?pageToken=%s", baseURL, pageToken)
		}

		response, err := fetchModelsPage(ctx, httpClient, url, projectID, publisher)
		if err != nil {
			return nil, err
		}

		// Extract model names, stripping the publishers/{publisher}/models/ prefix
		for _, model := range response.PublisherModels {
			modelName := extractModelName(model.Name)
			if modelName != "" {
				allModels = append(allModels, modelName)
			}
		}

		// Check for more pages
		if response.NextPageToken == "" {
			break
		}
		pageToken = response.NextPageToken
	}

	debuglog.Debug(debuglog.Detailed, "Listed %d models from publisher %s\n", len(allModels), publisher)
	return allModels, nil
}

// extractModelName extracts the model name from the full resource path
// Input: "publishers/google/models/gemini-2.0-flash"
// Output: "gemini-2.0-flash"
func extractModelName(fullName string) string {
	parts := strings.Split(fullName, "/")
	if len(parts) >= 4 && parts[0] == "publishers" && parts[2] == "models" {
		return parts[3]
	}
	// Fallback: return the last segment
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}

// sortModels sorts models by priority: Gemini > Claude > Others
// Within each group, models are sorted alphabetically
func sortModels(models []string) []string {
	sort.Slice(models, func(i, j int) bool {
		pi := modelPriority(models[i])
		pj := modelPriority(models[j])
		if pi != pj {
			return pi < pj
		}
		// Same priority: sort alphabetically (case-insensitive)
		return strings.ToLower(models[i]) < strings.ToLower(models[j])
	})
	return models
}

// modelPriority returns the sort priority for a model (lower = higher priority)
func modelPriority(model string) int {
	lower := strings.ToLower(model)
	switch {
	case strings.HasPrefix(lower, "gemini"):
		return 1
	case strings.HasPrefix(lower, "claude"):
		return 2
	default:
		return 3
	}
}

// knownGeminiModels is a curated list of Gemini models available on Vertex AI.
// Vertex AI doesn't provide a list API for Gemini models - they must be known ahead of time.
// This list is based on Google Cloud documentation as of January 2025.
// See: https://docs.cloud.google.com/vertex-ai/generative-ai/docs/models
var knownGeminiModels = []string{
	// Gemini 3 (Preview)
	"gemini-3-pro-preview",
	"gemini-3-flash-preview",
	// Gemini 2.5 (GA)
	"gemini-2.5-pro",
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
	// Gemini 2.0 (GA)
	"gemini-2.0-flash",
	"gemini-2.0-flash-lite",
}

// getKnownGeminiModels returns the curated list of Gemini models available on Vertex AI.
// Unlike third-party models which can be listed via the Model Garden API,
// Gemini models must be known ahead of time as there's no list endpoint for them.
func getKnownGeminiModels() []string {
	return knownGeminiModels
}

// isGeminiModel returns true if the model is a Gemini model
func isGeminiModel(modelName string) bool {
	return strings.HasPrefix(strings.ToLower(modelName), "gemini")
}

// isConversationalModel returns true if the model is suitable for text generation/chat
// Filters out image generation, embeddings, and other non-conversational models
func isConversationalModel(modelName string) bool {
	lower := strings.ToLower(modelName)

	// Exclude patterns for non-conversational models
	excludePatterns := []string{
		"imagen", // Image generation models
		"imagegeneration",
		"imagetext",
		"image-segmentation",
		"embedding", // Embedding models
		"textembedding",
		"multimodalembedding",
		"text-bison", // Legacy completion models (not chat)
		"text-unicorn",
		"code-bison", // Legacy code models
		"code-gecko",
		"codechat-bison", // Deprecated chat model
		"chat-bison",     // Deprecated chat model
		"veo",            // Video generation
		"chirp",          // Audio/speech models
		"medlm",          // Medical models (restricted)
		"medical",
	}

	for _, pattern := range excludePatterns {
		if strings.Contains(lower, pattern) {
			return false
		}
	}

	return true
}

// filterConversationalModels returns only models suitable for text generation/chat
func filterConversationalModels(models []string) []string {
	var filtered []string
	for _, model := range models {
		if isConversationalModel(model) {
			filtered = append(filtered, model)
		}
	}
	return filtered
}
