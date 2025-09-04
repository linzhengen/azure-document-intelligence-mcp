package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain"
)

const (
	apiVersion = "2023-07-31"
	maxRetries = 10
	retryDelay = 5 * time.Second
)

type diClient struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Document Intelligence client.
func NewClient(endpoint, apiKey string) domain.Client {
	return &diClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AnalyzeDocument analyzes the specified document URL.
func (c *diClient) AnalyzeDocument(ctx context.Context, modelID string, docURL string) (*domain.AnalyzeResult, error) {
	// 1. Send analysis request
	operationLocation, err := c.initiateAnalysis(ctx, modelID, docURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate analysis: %w", err)
	}

	// 2. Poll for the result
	result, err := c.pollForResult(ctx, operationLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for result: %w", err)
	}

	return result, nil
}

func (c *diClient) initiateAnalysis(ctx context.Context, modelID string, docURL string) (string, error) {
	requestURL := fmt.Sprintf("%s/documentintelligence/documentModels/%s:analyze?api-version=%s", c.endpoint, modelID, apiVersion)

	requestBody, err := json.Marshal(map[string]string{"urlSource": docURL})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

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

func (c *diClient) pollForResult(ctx context.Context, operationLocation string) (*domain.AnalyzeResult, error) {
	var result domain.AnalyzeResult

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, operationLocation, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create polling request: %w", err)
		}
		req.Header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send polling request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code during polling: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
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
			time.Sleep(retryDelay)
		default:
			return nil, fmt.Errorf("unknown status: %s", result.Status)
		}
	}

	return nil, fmt.Errorf("polling timed out")
}