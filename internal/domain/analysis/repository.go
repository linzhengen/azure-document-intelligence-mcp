package analysis

import "context"

type Repository interface {
	AnalyzeDocument(ctx context.Context, req AnalyzeDocumentRequest) (*AnalyzeResult, error)
}
