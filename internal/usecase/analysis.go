package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain/analysis"
)

// AnalysisParams defines the parameters for the document analysis tool.
type AnalysisParams struct {
	ModelID         string `json:"modelId"`
	DocumentURL     string `json:"documentUrl,omitempty"`
	DocumentContent string `json:"documentContent,omitempty"` // Base64 encoded content
	ContentType     string `json:"contentType,omitempty"`     // Required when documentContent is provided
}

// NewAnalysisHandler creates a tool handler for document analysis.
func NewAnalysisHandler(analyzerRepo analysis.Repository) func(context.Context, *mcp.CallToolRequest, *AnalysisParams) (*mcp.CallToolResult, *analysis.AnalyzeResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *AnalysisParams) (*mcp.CallToolResult, *analysis.AnalyzeResult, error) {
		if params.ModelID != "prebuilt-read" && params.ModelID != "prebuilt-layout" {
			return nil, nil, fmt.Errorf("unsupported modelId: %s", params.ModelID)
		}

		if (params.DocumentURL == "" && params.DocumentContent == "") || (params.DocumentURL != "" && params.DocumentContent != "") {
			return nil, nil, errors.New("either documentUrl or documentContent must be provided, but not both")
		}

		var content []byte
		var err error
		if params.DocumentContent != "" {
			if params.ContentType == "" {
				return nil, nil, errors.New("contentType must be provided when using documentContent")
			}
			content, err = base64.StdEncoding.DecodeString(params.DocumentContent)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to decode documentContent: %w", err)
			}
		}

		analysisReq := analysis.AnalyzeDocumentRequest{
			ModelID:     params.ModelID,
			DocURL:      params.DocumentURL,
			Content:     content,
			ContentType: params.ContentType,
		}

		result, err := analyzerRepo.AnalyzeDocument(ctx, analysisReq)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}
