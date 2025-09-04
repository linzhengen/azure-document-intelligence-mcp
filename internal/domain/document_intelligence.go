package domain

import (
	"context"
	"time"
)

// Client is an interface for abstracting interactions with the Azure Document Intelligence service.
type Client interface {
	AnalyzeDocument(ctx context.Context, modelID string, docURL string) (*AnalyzeResult, error)
}

// AnalyzeResult represents the complete result of a document analysis operation.
// See: https://learn.microsoft.com/en-us/azure/ai-services/document-intelligence/how-to-guides/use-sdk-rest-api?view=doc-intel-4.0.0&pivots=programming-language-rest-api#get-result-response-body
type AnalyzeResult struct {
	Status              string             `json:"status"`
	CreatedDateTime     time.Time          `json:"createdDateTime"`
	LastUpdatedDateTime time.Time          `json:"lastUpdatedDateTime"`
	AnalyzeResult       *AnalyzeResultBody `json:"analyzeResult"`
}

// AnalyzeResultBody is the body of the analysis result.
type AnalyzeResultBody struct {
	APIVersion string      `json:"apiVersion"`
	ModelID    string      `json:"modelId"`
	Content    string      `json:"content"`
	Pages      []Page      `json:"pages"`
	Paragraphs []Paragraph `json:"paragraphs"`
	Styles     []Style     `json:"styles"`
}

// Page represents a single page of a document.
type Page struct {
	PageNumber int       `json:"pageNumber"`
	Angle      float64   `json:"angle"`
	Width      float64   `json:"width"`
	Height     float64   `json:"height"`
	Unit       string    `json:"unit"`
	Words      []Word    `json:"words"`
	Lines      []Line    `json:"lines"`
	Spans      []Span    `json:"spans"`
}

// Paragraph represents an extracted paragraph.
type Paragraph struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Spans   []Span `json:"spans"`
}

// Style represents the style of the text.
type Style struct {
	IsHandwritten bool    `json:"isHandwritten"`
	Spans         []Span  `json:"spans"`
	Confidence    float64 `json:"confidence"`
}

// Word represents a word on a page.
type Word struct {
	Content    string    `json:"content"`
	Polygon    []float64 `json:"polygon"`
	Span       Span      `json:"span"`
	Confidence float64   `json:"confidence"`
}

// Line represents a line of text on a page.
type Line struct {
	Content string    `json:"content"`
	Polygon []float64 `json:"polygon"`
	Spans   []Span    `json:"spans"`
}

// Span indicates a range of text within the content.
type Span struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}