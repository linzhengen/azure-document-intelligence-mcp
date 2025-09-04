package analysis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain/analysis"
)

// MockRoundTripper is a mock implementation of http.RoundTripper for testing.
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip executes the mock round trip function.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// MockHTTPClient is a mock implementation of the HTTPClient interface.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}, nil
}

func TestAnalyzeDocument_Success(t *testing.T) {
	ctx := context.Background()
	operationLocation := "http://test.com/operation/123"

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodPost {
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Operation-Location": []string{operationLocation}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
			if req.Method == http.MethodGet && req.URL.String() == operationLocation {
				result := &analysis.AnalyzeResult{Status: "succeeded"}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{ModelID: "test-model", DocURL: "http://test.com/doc.pdf"}

	result, err := repo.AnalyzeDocument(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "succeeded", result.Status)
}

func TestAnalyzeDocument_Failed(t *testing.T) {
	ctx := context.Background()
	operationLocation := "http://test.com/operation/123"

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodPost {
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Operation-Location": []string{operationLocation}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
			if req.Method == http.MethodGet && req.URL.String() == operationLocation {
				result := &analysis.AnalyzeResult{Status: "failed"}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{ModelID: "test-model", DocURL: "http://test.com/doc.pdf"}

	_, err := repo.AnalyzeDocument(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "analysis failed")
}

func TestAnalyzeDocument_PollingTimeout(t *testing.T) {
	ctx := context.Background()
	operationLocation := "http://test.com/operation/123"

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodPost {
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Operation-Location": []string{operationLocation}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
			if req.Method == http.MethodGet && req.URL.String() == operationLocation {
				result := &analysis.AnalyzeResult{Status: "running"}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	// Override retryDelay for faster testing
	originalRetryDelay := retryDelay
	retryDelay = 1 * time.Millisecond
	defer func() { retryDelay = originalRetryDelay }()

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{ModelID: "test-model", DocURL: "http://test.com/doc.pdf"}

	_, err := repo.AnalyzeDocument(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "polling timed out")
}

func TestAnalyzeDocument_InitiateAnalysisFails(t *testing.T) {
	ctx := context.Background()

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
			}, nil
		},
	}

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{ModelID: "test-model", DocURL: "http://test.com/doc.pdf"}

	_, err := repo.AnalyzeDocument(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initiate analysis")
}

func TestAnalyzeDocument_UnsupportedStatus(t *testing.T) {
	ctx := context.Background()
	operationLocation := "http://test.com/operation/123"

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodPost {
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Operation-Location": []string{operationLocation}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
			if req.Method == http.MethodGet && req.URL.String() == operationLocation {
				result := &analysis.AnalyzeResult{Status: "weird_status"}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{ModelID: "test-model", DocURL: "http://test.com/doc.pdf"}

	_, err := repo.AnalyzeDocument(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown status: weird_status")
}

func TestAnalyzeDocument_WithContent(t *testing.T) {
	ctx := context.Background()
	operationLocation := "http://test.com/operation/123"

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodPost {
				// Check that the content type is passed correctly
				assert.Equal(t, "application/pdf", req.Header.Get("Content-Type"))
				bodyBytes, _ := io.ReadAll(req.Body)
				assert.Equal(t, "dummy-content", string(bodyBytes))

				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Operation-Location": []string{operationLocation}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}
			if req.Method == http.MethodGet && req.URL.String() == operationLocation {
				result := &analysis.AnalyzeResult{Status: "succeeded"}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	repo := NewRepositoryWithClient("http://test.com", "dummy-key", mockClient)
	req := analysis.AnalyzeDocumentRequest{
		ModelID:     "test-model",
		Content:     []byte("dummy-content"),
		ContentType: "application/pdf",
	}

	result, err := repo.AnalyzeDocument(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "succeeded", result.Status)
}
