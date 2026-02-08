// Package geminicommon provides shared utilities for Gemini API integrations.
// Used by both the standalone Gemini provider (API key auth) and VertexAI provider (ADC auth).
package geminicommon

import (
	"fmt"
	"strings"

	"github.com/creativerezz/aio/internal/chat"
	"google.golang.org/genai"
)

// Citation formatting constants
const (
	CitationHeader    = "\n\n## Sources\n\n"
	CitationSeparator = "\n"
	CitationFormat    = "- [%s](%s)"
)

// ConvertMessages converts aio chat messages to genai Content format.
// Gemini's API only accepts "user" and "model" roles, so other roles are mapped to "user".
func ConvertMessages(msgs []*chat.ChatCompletionMessage) []*genai.Content {
	var contents []*genai.Content

	for _, msg := range msgs {
		content := &genai.Content{Parts: []*genai.Part{}}

		switch msg.Role {
		case chat.ChatMessageRoleAssistant:
			content.Role = "model"
		case chat.ChatMessageRoleUser:
			content.Role = "user"
		case chat.ChatMessageRoleSystem, chat.ChatMessageRoleDeveloper, chat.ChatMessageRoleFunction, chat.ChatMessageRoleTool:
			// Gemini's API only accepts "user" and "model" roles.
			// Map all other roles to "user" to preserve instruction context.
			content.Role = "user"
		default:
			content.Role = "user"
		}

		if strings.TrimSpace(msg.Content) != "" {
			content.Parts = append(content.Parts, &genai.Part{Text: msg.Content})
		}

		// Handle multi-content messages (images, etc.)
		for _, part := range msg.MultiContent {
			switch part.Type {
			case chat.ChatMessagePartTypeText:
				content.Parts = append(content.Parts, &genai.Part{Text: part.Text})
			case chat.ChatMessagePartTypeImageURL:
				// TODO: Handle image URLs if needed
				// This would require downloading and converting to inline data
			}
		}

		contents = append(contents, content)
	}

	return contents
}

// ExtractText extracts just the text parts from a Gemini response.
func ExtractText(response *genai.GenerateContentResponse) string {
	if response == nil {
		return ""
	}

	var builder strings.Builder
	for _, candidate := range response.Candidates {
		if candidate == nil || candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part != nil && part.Text != "" {
				builder.WriteString(part.Text)
			}
		}
	}
	return builder.String()
}

// ExtractTextWithCitations extracts text content from the response and appends
// any web citations in a standardized format.
func ExtractTextWithCitations(response *genai.GenerateContentResponse) string {
	if response == nil {
		return ""
	}

	text := ExtractText(response)
	citations := ExtractCitations(response)
	if len(citations) > 0 {
		return text + CitationHeader + strings.Join(citations, CitationSeparator)
	}
	return text
}

// ExtractCitations extracts web citations from grounding metadata.
func ExtractCitations(response *genai.GenerateContentResponse) []string {
	if response == nil || len(response.Candidates) == 0 {
		return nil
	}

	citationMap := make(map[string]bool)
	var citations []string
	for _, candidate := range response.Candidates {
		if candidate == nil || candidate.GroundingMetadata == nil {
			continue
		}
		chunks := candidate.GroundingMetadata.GroundingChunks
		if len(chunks) == 0 {
			continue
		}
		for _, chunk := range chunks {
			if chunk == nil || chunk.Web == nil {
				continue
			}
			uri := chunk.Web.URI
			title := chunk.Web.Title
			if uri == "" || title == "" {
				continue
			}
			key := uri + "|" + title
			if !citationMap[key] {
				citationMap[key] = true
				citations = append(citations, fmt.Sprintf(CitationFormat, title, uri))
			}
		}
	}
	return citations
}
