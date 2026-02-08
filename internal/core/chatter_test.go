package core

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/creativerezz/aio/internal/chat"
	"github.com/creativerezz/aio/internal/domain"
	"github.com/creativerezz/aio/internal/plugins/db/fsdb"
)

// mockVendor implements the ai.Vendor interface for testing
type mockVendor struct {
	sendStreamError error
	streamChunks    []domain.StreamUpdate
	sendFunc        func(context.Context, []*chat.ChatCompletionMessage, *domain.ChatOptions) (string, error)
}

func (m *mockVendor) GetName() string {
	return "mock"
}

func (m *mockVendor) GetSetupDescription() string {
	return "mock vendor"
}

func (m *mockVendor) IsConfigured() bool {
	return true
}

func (m *mockVendor) Configure() error {
	return nil
}

func (m *mockVendor) Setup() error {
	return nil
}

func (m *mockVendor) SetupFillEnvFileContent(*bytes.Buffer) {
}

func (m *mockVendor) ListModels() ([]string, error) {
	return []string{"test-model"}, nil
}

func (m *mockVendor) SendStream(messages []*chat.ChatCompletionMessage, opts *domain.ChatOptions, responseChan chan domain.StreamUpdate) error {
	// Send chunks if provided (for successful streaming test)
	if m.streamChunks != nil {
		for _, chunk := range m.streamChunks {
			responseChan <- chunk
		}
	}
	// Close the channel like real vendors do
	close(responseChan)
	return m.sendStreamError
}

func (m *mockVendor) Send(ctx context.Context, messages []*chat.ChatCompletionMessage, opts *domain.ChatOptions) (string, error) {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, messages, opts)
	}
	return "test response", nil
}

func (m *mockVendor) NeedsRawMode(modelName string) bool {
	return false
}

func TestChatter_Send_SuppressThink(t *testing.T) {
	tempDir := t.TempDir()
	db := fsdb.NewDb(tempDir)

	mockVendor := &mockVendor{}

	chatter := &Chatter{
		db:     db,
		Stream: false,
		vendor: mockVendor,
		model:  "test-model",
	}

	request := &domain.ChatRequest{
		Message: &chat.ChatCompletionMessage{
			Role:    chat.ChatMessageRoleUser,
			Content: "test",
		},
	}

	opts := &domain.ChatOptions{
		Model:         "test-model",
		SuppressThink: true,
		ThinkStartTag: "<think>",
		ThinkEndTag:   "</think>",
	}

	// custom send function returning a message with think tags
	mockVendor.sendFunc = func(ctx context.Context, msgs []*chat.ChatCompletionMessage, o *domain.ChatOptions) (string, error) {
		return "<think>hidden</think> visible", nil
	}

	session, err := chatter.Send(request, opts)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if session == nil {
		t.Fatal("expected session")
	}
	last := session.GetLastMessage()
	if last.Content != "visible" {
		t.Errorf("expected filtered content 'visible', got %q", last.Content)
	}
}

func TestChatter_Send_StreamingErrorPropagation(t *testing.T) {
	// Create a temporary database for testing
	tempDir := t.TempDir()
	db := fsdb.NewDb(tempDir)

	// Create a mock vendor that will return an error from SendStream
	expectedError := errors.New("streaming error")
	mockVendor := &mockVendor{
		sendStreamError: expectedError,
	}

	// Create chatter with streaming enabled
	chatter := &Chatter{
		db:     db,
		Stream: true, // Enable streaming to trigger SendStream path
		vendor: mockVendor,
		model:  "test-model",
	}

	// Create a test request
	request := &domain.ChatRequest{
		Message: &chat.ChatCompletionMessage{
			Role:    chat.ChatMessageRoleUser,
			Content: "test message",
		},
	}

	// Create test options
	opts := &domain.ChatOptions{
		Model: "test-model",
	}

	// Call Send and expect it to return the streaming error
	session, err := chatter.Send(request, opts)

	// Verify that the error from SendStream is propagated
	if err == nil {
		t.Fatal("Expected error to be returned, but got nil")
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error %q, but got %q", expectedError, err)
	}

	// Session should still be returned (it was built successfully before the streaming error)
	if session == nil {
		t.Error("Expected session to be returned even when streaming error occurs")
	}
}

func TestChatter_Send_StreamingSuccessfulAggregation(t *testing.T) {
	// Create a temporary database for testing
	tempDir := t.TempDir()
	db := fsdb.NewDb(tempDir)

	// Create test chunks that should be aggregated
	chunks := []string{"Hello", " ", "world", "!", " This", " is", " a", " test."}
	testChunks := make([]domain.StreamUpdate, len(chunks))
	for i, c := range chunks {
		testChunks[i] = domain.StreamUpdate{Type: domain.StreamTypeContent, Content: c}
	}
	expectedMessage := "Hello world! This is a test."

	// Create a mock vendor that will send chunks successfully
	mockVendor := &mockVendor{
		sendStreamError: nil, // No error for successful streaming
		streamChunks:    testChunks,
	}

	// Create chatter with streaming enabled
	chatter := &Chatter{
		db:     db,
		Stream: true, // Enable streaming to trigger SendStream path
		vendor: mockVendor,
		model:  "test-model",
	}

	// Create a test request
	request := &domain.ChatRequest{
		Message: &chat.ChatCompletionMessage{
			Role:    chat.ChatMessageRoleUser,
			Content: "test message",
		},
	}

	// Create test options
	opts := &domain.ChatOptions{
		Model: "test-model",
	}

	// Call Send and expect successful aggregation
	session, err := chatter.Send(request, opts)

	// Verify no error occurred
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Verify session was returned
	if session == nil {
		t.Fatal("Expected session to be returned")
	}

	// Verify the message was aggregated correctly
	messages := session.GetVendorMessages()
	if len(messages) != 2 { // user message + assistant response
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Check the assistant's response (last message)
	assistantMessage := messages[len(messages)-1]
	if assistantMessage.Role != chat.ChatMessageRoleAssistant {
		t.Errorf("Expected assistant role, got %s", assistantMessage.Role)
	}

	if assistantMessage.Content != expectedMessage {
		t.Errorf("Expected aggregated message %q, got %q", expectedMessage, assistantMessage.Content)
	}
}

func TestChatter_Send_StreamingMetadataPropagation(t *testing.T) {
	// Create a temporary database for testing
	tempDir := t.TempDir()
	db := fsdb.NewDb(tempDir)

	// Create test chunks: one content, one usage metadata
	testChunks := []domain.StreamUpdate{
		{
			Type:    domain.StreamTypeContent,
			Content: "Test content",
		},
		{
			Type: domain.StreamTypeUsage,
			Usage: &domain.UsageMetadata{
				InputTokens:  10,
				OutputTokens: 5,
				TotalTokens:  15,
			},
		},
	}

	// Create a mock vendor
	mockVendor := &mockVendor{
		sendStreamError: nil,
		streamChunks:    testChunks,
	}

	// Create chatter with streaming enabled
	chatter := &Chatter{
		db:     db,
		Stream: true,
		vendor: mockVendor,
		model:  "test-model",
	}

	// Create a test request
	request := &domain.ChatRequest{
		Message: &chat.ChatCompletionMessage{
			Role:    chat.ChatMessageRoleUser,
			Content: "test message",
		},
	}

	// Create an update channel to capture stream events
	updateChan := make(chan domain.StreamUpdate, 10)

	// Create test options with UpdateChan
	opts := &domain.ChatOptions{
		Model:      "test-model",
		UpdateChan: updateChan,
		Quiet:      true, // Suppress stdout/stderr
	}

	// Call Send
	_, err := chatter.Send(request, opts)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	close(updateChan)

	// Verify we received the metadata event
	var usageReceived bool
	for update := range updateChan {
		if update.Type == domain.StreamTypeUsage {
			usageReceived = true
			if update.Usage == nil {
				t.Error("Expected usage metadata to be non-nil")
			} else {
				if update.Usage.TotalTokens != 15 {
					t.Errorf("Expected 15 total tokens, got %d", update.Usage.TotalTokens)
				}
			}
		}
	}

	if !usageReceived {
		t.Error("Expected to receive a usage metadata update, but didn't")
	}
}
