package analysis

import "time"

// AnalyzeDocumentRequest represents the request body for analyzing a document.
type AnalyzeDocumentRequest struct {
	URLSource      *string `json:"urlSource,omitempty"`
	Base64Source   []byte  `json:"base64Source,omitempty"`
}

// AnalyzeOperationResult represents the status and result of the analyze operation.
type AnalyzeOperationResult struct {
	Status              string         `json:"status"`
	CreatedDateTime     time.Time      `json:"createdDateTime"`
	LastUpdatedDateTime time.Time      `json:"lastUpdatedDateTime"`
	Error               *Error         `json:"error,omitempty"`
	AnalyzeResult       *AnalyzeResult `json:"analyzeResult,omitempty"`
}

// AnalyzeResult represents the document analysis result.
type AnalyzeResult struct {
	ApiVersion      string          `json:"apiVersion"`
	ModelID         string          `json:"modelId"`
	StringIndexType string          `json:"stringIndexType"`
	ContentFormat   *string         `json:"contentFormat,omitempty"`
	Content         string          `json:"content"`
	Pages           []Page          `json:"pages"`
	Paragraphs      []*Paragraph    `json:"paragraphs,omitempty"`
	Tables          []*Table        `json:"tables,omitempty"`
	Figures         []*Figure       `json:"figures,omitempty"`
	Sections        []*Section      `json:"sections,omitempty"`
	KeyValuePairs   []*KeyValuePair `json:"keyValuePairs,omitempty"`
	Styles          []*Style        `json:"styles,omitempty"`
	Languages       []*Language     `json:"languages,omitempty"`
	Documents       []*Document     `json:"documents,omitempty"`
	Warnings        []*Warning      `json:"warnings,omitempty"`
}

// Page represents content and layout elements extracted from a page from the input.
type Page struct {
	PageNumber     int32            `json:"pageNumber"`
	Angle          *float32         `json:"angle,omitempty"`
	Width          *float32         `json:"width,omitempty"`
	Height         *float32         `json:"height,omitempty"`
	Unit           *string          `json:"unit,omitempty"`
	Spans          []Span           `json:"spans"`
	Words          []*Word          `json:"words,omitempty"`
	SelectionMarks []*SelectionMark `json:"selectionMarks,omitempty"`
	Lines          []*Line          `json:"lines,omitempty"`
	Barcodes       []*Barcode       `json:"barcodes,omitempty"`
	Formulas       []*Formula       `json:"formulas,omitempty"`
}

// Span represents a contiguous region of the concatenated content property.
type Span struct {
	Offset int32 `json:"offset"`
	Length int32 `json:"length"`
}

// Word represents a word object.
type Word struct {
	Content    string     `json:"content"`
	Polygon    []float32  `json:"polygon,omitempty"`
	Span       Span       `json:"span"`
	Confidence float32    `json:"confidence"`
}

// SelectionMark represents a selection mark object.
type SelectionMark struct {
	State      string    `json:"state"`
	Polygon    []float32 `json:"polygon,omitempty"`
	Span       Span      `json:"span"`
	Confidence float32   `json:"confidence"`
}

// Line represents a content line object.
type Line struct {
	Content string    `json:"content"`
	Polygon []float32 `json:"polygon,omitempty"`
	Spans   []Span    `json:"spans"`
}

// Barcode represents a barcode object.
type Barcode struct {
	Kind       string    `json:"kind"`
	Value      string    `json:"value"`
	Polygon    []float32 `json:"polygon,omitempty"`
	Span       Span      `json:"span"`
	Confidence float32   `json:"confidence"`
}

// Formula represents a formula object.
type Formula struct {
	Kind       string    `json:"kind"`
	Value      string    `json:"value"`
	Polygon    []float32 `json:"polygon,omitempty"`
	Span       Span      `json:"span"`
	Confidence float32   `json:"confidence"`
}

// Paragraph represents a paragraph object.
type Paragraph struct {
	Role            *string          `json:"role,omitempty"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
}

// BoundingRegion represents a bounding polygon on a specific page.
type BoundingRegion struct {
	PageNumber int32   `json:"pageNumber"`
	Polygon    []float32 `json:"polygon"`
}

// Table represents a table object.
type Table struct {
	RowCount        int32            `json:"rowCount"`
	ColumnCount     int32            `json:"columnCount"`
	Cells           []Cell           `json:"cells"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
	Caption         *Caption         `json:"caption,omitempty"`
	Footnotes       []*Footnote      `json:"footnotes,omitempty"`
}

// Cell represents a cell in a table.
type Cell struct {
	Kind            *string          `json:"kind,omitempty"`
	RowIndex        int32            `json:"rowIndex"`
	ColumnIndex     int32            `json:"columnIndex"`
	RowSpan         *int32           `json:"rowSpan,omitempty"`
	ColumnSpan      *int32           `json:"columnSpan,omitempty"`
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
	Elements        []string         `json:"elements,omitempty"`
}

// Figure represents a figure in the document.
type Figure struct {
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
	Elements        []string         `json:"elements,omitempty"`
	Caption         *Caption         `json:"caption,omitempty"`
	Footnotes       []*Footnote      `json:"footnotes,omitempty"`
	ID              *string          `json:"id,omitempty"`
}

// Section represents a section in the document.
type Section struct {
	Spans    []Span   `json:"spans"`
	Elements []string `json:"elements,omitempty"`
}

// Caption represents a caption for a table or figure.
type Caption struct {
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
	Elements        []string         `json:"elements,omitempty"`
}

// Footnote represents a footnote.
type Footnote struct {
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
	Elements        []string         `json:"elements,omitempty"`
}

// KeyValuePair represents a key-value pair.
type KeyValuePair struct {
	Key        KeyValueElement  `json:"key"`
	Value      *KeyValueElement `json:"value,omitempty"`
	Confidence float32          `json:"confidence"`
}

// KeyValueElement represents the key or value in a key-value pair.
type KeyValueElement struct {
	Content         string           `json:"content"`
	BoundingRegions []BoundingRegion `json:"boundingRegions,omitempty"`
	Spans           []Span           `json:"spans"`
}

// Style represents observed text styles.
type Style struct {
	IsHandwritten     *bool    `json:"isHandwritten,omitempty"`
	SimilarFontFamily *string  `json:"similarFontFamily,omitempty"`
	FontStyle         *string  `json:"fontStyle,omitempty"`
	FontWeight        *string  `json:"fontWeight,omitempty"`
	Color             *string  `json:"color,omitempty"`
	BackgroundColor   *string  `json:"backgroundColor,omitempty"`
	Spans             []Span   `json:"spans"`
	Confidence        float32  `json:"confidence"`
}

// Language represents a detected language.
type Language struct {
	Locale     string  `json:"locale"`
	Spans      []Span  `json:"spans"`
	Confidence float32 `json:"confidence"`
}

// Document represents an extracted document.
type Document struct {
	DocType         string                   `json:"docType"`
	BoundingRegions []BoundingRegion         `json:"boundingRegions,omitempty"`
	Spans           []Span                   `json:"spans"`
	Fields          map[string]*DocumentField `json:"fields,omitempty"`
	Confidence      float32                  `json:"confidence"`
}

// DocumentField represents the content and location of a field value.
type DocumentField struct {
	Type               string                   `json:"type"`
	ValueString        *string                  `json:"valueString,omitempty"`
	ValueDate          *time.Time               `json:"valueDate,omitempty"`
	ValueTime          *time.Time               `json:"valueTime,omitempty"`
	ValuePhoneNumber   *string                  `json:"valuePhoneNumber,omitempty"`
	ValueNumber        *float64                 `json:"valueNumber,omitempty"`
	ValueInteger       *int64                   `json:"valueInteger,omitempty"`
	ValueSelectionMark *string                  `json:"valueSelectionMark,omitempty"`
	ValueSignature     *string                  `json:"valueSignature,omitempty"`
	ValueCountryRegion *string                  `json:"valueCountryRegion,omitempty"`
	ValueArray         []*DocumentField         `json:"valueArray,omitempty"`
	ValueObject        map[string]*DocumentField `json:"valueObject,omitempty"`
	ValueCurrency      *CurrencyValue           `json:"valueCurrency,omitempty"`
	ValueAddress       *AddressValue            `json:"valueAddress,omitempty"`
	ValueBoolean       *bool                    `json:"valueBoolean,omitempty"`
	ValueSelectionGroup []string			 `json:"valueSelectionGroup,omitempty"`
	Content            *string                  `json:"content,omitempty"`
	BoundingRegions    []BoundingRegion         `json:"boundingRegions,omitempty"`
	Spans              []Span                   `json:"spans,omitempty"`
	Confidence         *float32                 `json:"confidence,omitempty"`
}

// CurrencyValue represents a currency field value.
type CurrencyValue struct {
	Amount         float64 `json:"amount"`
	CurrencySymbol *string `json:"currencySymbol,omitempty"`
	CurrencyCode   *string `json:"currencyCode,omitempty"`
}

// AddressValue represents an address field value.
type AddressValue struct {
	HouseNumber   *string `json:"houseNumber,omitempty"`
	PoBox         *string `json:"poBox,omitempty"`
	Road          *string `json:"road,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	PostalCode    *string `json:"postalCode,omitempty"`
	CountryRegion *string `json:"countryRegion,omitempty"`
	StreetAddress *string `json:"streetAddress,omitempty"`
	Unit          *string `json:"unit,omitempty"`
	CityDistrict  *string `json:"cityDistrict,omitempty"`
	StateDistrict *string `json:"stateDistrict,omitempty"`
	Suburb        *string `json:"suburb,omitempty"`
	House         *string `json:"house,omitempty"`
	Level         *string `json:"level,omitempty"`
}

// Warning represents a warning from the service.
type Warning struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Target  *string `json:"target,omitempty"`
}

// Error represents the error object from the service.
type Error struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Target     *string     `json:"target,omitempty"`
	Details    []*Error    `json:"details,omitempty"`
	InnerError *InnerError `json:"innererror,omitempty"`
}

// InnerError represents more specific information about an error.
type InnerError struct {
	Code       *string     `json:"code,omitempty"`
	Message    *string     `json:"message,omitempty"`
	InnerError *InnerError `json:"innererror,omitempty"`
}