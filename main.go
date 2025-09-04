package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/linzhengen/azure-document-intelligence-mcp/config"
	"github.com/linzhengen/azure-document-intelligence-mcp/internal/domain/analysis"
	analysisinfra "github.com/linzhengen/azure-document-intelligence-mcp/internal/infrastructure/analysis"
	"github.com/linzhengen/azure-document-intelligence-mcp/internal/usecase"
)

func main() {
	ctx := context.Background()

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize infrastructure layer
	analysisRepo := analysisinfra.NewRepository(cfg.AzureEndpoint, cfg.AzureAPIKey)

	// 3. Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "azure-document-intelligence-mcp",
		Version: "1.0.0",
	}, nil)

	// 4. Create the tool handler
	analysisHandler := usecase.NewAnalysisHandler(analysisRepo)

	// 5. Register the analysis tool using ToolFor to ensure correct handler type
	analyzeToolDef := &mcp.Tool{
		Name:        "analyze_document",
		Description: "Analyzes a document using Azure Document Intelligence. Pass 'prebuilt-read' or 'prebuilt-layout' in the modelId parameter. Provide the document either via 'documentUrl' or by passing base64 encoded data in 'documentContent' with its 'contentType'.",
	}
	tool, handler := mcp.ToolFor[*usecase.AnalysisParams, *analysis.AnalyzeResult](analyzeToolDef, analysisHandler)
	server.AddTool(tool, handler)

	// 6. Run the server with StdioTransport
	log.Println("Starting MCP server over stdio")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
