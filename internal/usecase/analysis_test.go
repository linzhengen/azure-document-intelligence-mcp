
package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain/analysis"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAnalysisRepository is a mock implementation of the analysis.Repository interface.
type MockAnalysisRepository struct {
	AnalyzeDocumentFunc func(ctx context.Context, req analysis.AnalyzeDocumentRequest) (*analysis.AnalyzeResult, error)
}

func (m *MockAnalysisRepository) AnalyzeDocument(ctx context.Context, req analysis.AnalyzeDocumentRequest) (*analysis.AnalyzeResult, error) {
	if m.AnalyzeDocumentFunc != nil {
		return m.AnalyzeDocumentFunc(ctx, req)
	}
	return &analysis.AnalyzeResult{Status: "succeeded"}, nil
}

func TestAnalysisHandler_SuccessWithURL(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:     "prebuilt-read",
		DocumentURL: "http://example.com/doc.pdf",
	}

	_, result, err := handler(ctx, nil, params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "succeeded", result.Status)
}

func TestAnalysisHandler_SuccessWithContent(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	content := base64.StdEncoding.EncodeToString([]byte("dummy content"))
	params := &AnalysisParams{
		ModelID:         "prebuilt-layout",
		DocumentContent: content,
		ContentType:     "application/pdf",
	}

	_, result, err := handler(ctx, nil, params)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "succeeded", result.Status)
}

func TestAnalysisHandler_UnsupportedModelID(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:     "unsupported-model",
		DocumentURL: "http://example.com/doc.pdf",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported modelId")
}

func TestAnalysisHandler_MissingDocumentSource(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID: "prebuilt-read",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Equal(t, "either documentUrl or documentContent must be provided, but not both", err.Error())
}

func TestAnalysisHandler_BothDocumentSourcesProvided(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:         "prebuilt-read",
		DocumentURL:     "http://example.com/doc.pdf",
		DocumentContent: "dummy-content",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Equal(t, "either documentUrl or documentContent must be provided, but not both", err.Error())
}

func TestAnalysisHandler_MissingContentType(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:         "prebuilt-read",
		DocumentContent: "dummy-content",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Equal(t, "contentType must be provided when using documentContent", err.Error())
}

func TestAnalysisHandler_InvalidBase64Content(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAnalysisRepository{}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:         "prebuilt-read",
		DocumentContent: "invalid-base64",
		ContentType:     "application/pdf",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode documentContent")
}

func TestAnalysisHandler_AnalyzerError(t *testing.T) {
	ctx := context.Background()
	analyzerErr := errors.New("analyzer error")
	mockRepo := &MockAnalysisRepository{
		AnalyzeDocumentFunc: func(ctx context.Context, req analysis.AnalyzeDocumentRequest) (*analysis.AnalyzeResult, error) {
			return nil, analyzerErr
		},
	}
	handler := NewAnalysisHandler(mockRepo)

	params := &AnalysisParams{
		ModelID:     "prebuilt-read",
		DocumentURL: "http://example.com/doc.pdf",
	}

	_, _, err := handler(ctx, nil, params)

	require.Error(t, err)
	assert.Equal(t, analyzerErr, err)
}
