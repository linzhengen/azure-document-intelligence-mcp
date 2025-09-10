package analysis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain/analysis"
)

const (
	apiVersion = "2024-11-30"
	maxRetries = 10
)

var retryDelay = 5 * time.Second

// HTTPClient is an interface for making HTTP requests.
// It's implemented by *http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type repository struct {
	endpoint   string
	apiKey     string
	httpClient HTTPClient
}

// NewRepository creates a new Document Intelligence client.
func NewRepository(endpoint, apiKey string, timeout int) analysis.Repository {
	return &repository{
		endpoint: endpoint,
		apiKey:   apiKey,
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// NewRepositoryWithClient creates a new Document Intelligence client with a custom http client.
func NewRepositoryWithClient(endpoint, apiKey string, httpClient HTTPClient) analysis.Repository {
	return &repository{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// AnalyzeDocument analyzes the specified document URL.
func (r *repository) AnalyzeDocument(ctx context.Context, modelID string, options analysis.AnalyzeDocumentOptions) (*analysis.AnalyzeOperationResult, error) {
	// 1. Send analysis request
	operationLocation, err := r.initiateAnalysis(ctx, modelID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate analysis: %w", err)
	}

	// 2. Poll for the result
	result, err := r.pollForResult(ctx, operationLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for result: %w", err)
	}

	return result, nil
}

func (r *repository) initiateAnalysis(ctx context.Context, modelID string, options analysis.AnalyzeDocumentOptions) (string, error) {
	requestURL := fmt.Sprintf("%s/documentintelligence/documentModels/%s:analyze?api-version=%s", r.endpoint, modelID, apiVersion)

	var requestBody io.Reader
	var contentType string

	if options.DocURL != "" {
		jsonBody, err := json.Marshal(map[string]string{"urlSource": options.DocURL})
		if err != nil {
			return "", fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
		contentType = "application/json"
	} else if options.Content != nil {
		requestBody = bytes.NewBuffer(options.Content)
		contentType = options.ContentType
	} else {
		return "", fmt.Errorf("no document source provided (URL or content)")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Ocp-Apim-Subscription-Key", r.apiKey)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	operationLocation := resp.Header.Get("Operation-Location")
	if operationLocation == "" {
		return "", fmt.Errorf("Operation-Location header not found")
	}

	return operationLocation, nil
}

func (r *repository) pollForResult(ctx context.Context, operationLocation string) (*analysis.AnalyzeOperationResult, error) {
	var result analysis.AnalyzeOperationResult

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, operationLocation, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create polling request: %w", err)
		}
		req.Header.Set("Ocp-Apim-Subscription-Key", r.apiKey)

		resp, err := r.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send polling request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code during polling: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read polling response body: %w", err)
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal polling response: %w", err)
		}

		switch result.Status {
		case "succeeded":
			return &result, nil
		case "failed":
			return nil, fmt.Errorf("analysis failed")
		case "running", "notStarted":
			// Continue
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
		default:
			return nil, fmt.Errorf("unknown status: %s", result.Status)
		}
	}

	return nil, fmt.Errorf("polling timed out")
}
