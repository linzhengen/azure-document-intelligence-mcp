package usecase

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain"
)

// AnalysisParams defines the parameters for the document analysis tool.
type AnalysisParams struct {
	ModelID     string `json:"modelId"`
	DocumentURL string `json:"documentUrl"`
}

// NewAnalysisHandler creates a tool handler for document analysis.
func NewAnalysisHandler(diClient domain.Client) func(context.Context, *mcp.CallToolRequest, *AnalysisParams) (*mcp.CallToolResult, *domain.AnalyzeResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *AnalysisParams) (*mcp.CallToolResult, *domain.AnalyzeResult, error) {
		if params.ModelID != "prebuilt-read" && params.ModelID != "prebuilt-layout" {
			return nil, nil, fmt.Errorf("unsupported modelId: %s", params.ModelID)
		}
		result, err := diClient.AnalyzeDocument(ctx, params.ModelID, params.DocumentURL)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}
