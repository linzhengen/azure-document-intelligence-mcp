package domain

import (
	"context"
	"time"
)

// AnalyzeDocumentRequest is an interface for abstracting interactions with the Azure Document Intelligence service.
type AnalyzeDocumentRequest struct {
	ModelID     string
	DocURL      string
	Content     []byte
	ContentType string
}

type Client interface {
	AnalyzeDocument(ctx context.Context, req AnalyzeDocumentRequest) (*AnalyzeResult, error)
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
	Tables     []Table     `json:"tables"`
	Styles     []Style     `json:"styles"`
	Languages  []Language  `json:"languages"`
}

// Page represents a single page of a document.
type Page struct {
	PageNumber     int             `json:"pageNumber"`
	Angle          float64         `json:"angle"`
	Width          float64         `json:"width"`
	Height         float64         `json:"height"`
	Unit           string          `json:"unit"`
	Words          []Word          `json:"words"`
	Lines          []Line          `json:"lines"`
	Spans          []Span          `json:"spans"`
	SelectionMarks []SelectionMark `json:"selectionMarks"`
}

// Paragraph represents an extracted paragraph.
type Paragraph struct {
	Role            string           `json:"role"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Spans           []Span           `json:"spans"`
}

// BoundingRegion represents a region in the document.
type BoundingRegion struct {
	PageNumber int       `json:"pageNumber"`
	Polygon    []float64 `json:"polygon"`
}

// Table represents a table extracted from the document.
type Table struct {
	RowCount        int              `json:"rowCount"`
	ColumnCount     int              `json:"columnCount"`
	Cells           []Cell           `json:"cells"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Spans           []Span           `json:"spans"`
}

// Cell represents a cell in a table.
type Cell struct {
	Kind            string           `json:"kind"`
	RowIndex        int              `json:"rowIndex"`
	ColumnIndex     int              `json:"columnIndex"`
	RowSpan         int              `json:"rowSpan,omitempty"`
	ColumnSpan      int              `json:"columnSpan,omitempty"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions"`
	Spans           []Span           `json:"spans"`
}

// SelectionMark represents a checkbox or radio button.
type SelectionMark struct {
	State      string    `json:"state"` // "selected", "unselected"
	Polygon    []float64 `json:"polygon"`
	Span       Span      `json:"span"`
	Confidence float64   `json:"confidence"`
}

// Language represents a detected language.
type Language struct {
	Locale     string  `json:"locale"`
	Spans      []Span  `json:"spans"`
	Confidence float64 `json:"confidence"`
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
