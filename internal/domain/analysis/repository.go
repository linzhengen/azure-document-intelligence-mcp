package analysis

import "context"

type AnalyzeDocumentOptions struct {
	DocURL      string
	Content     []byte
	ContentType string
}

type Repository interface {
	AnalyzeDocument(ctx context.Context, modelID string, options AnalyzeDocumentOptions) (*AnalyzeOperationResult, error)
}