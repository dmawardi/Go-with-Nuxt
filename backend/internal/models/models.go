package models

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Conditions when using a findall route
type QueryConditionParameters struct {
	Condition string
	Value     interface{}
}

// Login
type Login struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type ValidationError struct {
	Validation_errors map[string][]string `json:"validation_errors"`
}

type BaseFindAllQueryParams struct {
	Page       int
	Limit      int
	Offset     int
	Order      string
	Conditions map[string]interface{}
}

type BulkDeleteResponse struct {
	Success        bool    `json:"success"`
	DeletedRecords int     `json:"deleted_records"`
	Errors         []error `json:"errors"`
}

// Schema meta data (attached to find all requests)
type SchemaMetaData interface {
	CalculateCurrentlyShowingRecords() int
	// CalculateTotalPages() int
	GetMetaData() schemaMetaData
}
type schemaMetaData struct {
	Total_Records    int64 `json:"total_records"`    // Total number of records in the entire dataset
	Records_Per_Page int   `json:"records_per_page"` // Number of records displayed per page
	Total_Pages      int   `json:"total_pages"`      // Total number of pages
	Current_Page     int   `json:"current_page"`     // Current page number
	Next_Page        *int  `json:"next_page"`        // Next page number (null if there is no next page)
	Prev_Page        *int  `json:"prev_page"`        // Previous page number (null if there is no previous page)
}

// Constructor
func NewSchemaMetaData(totalRecords int64, recordsPerPage int, totalPages int, currentPage int, nextPage *int, prevPage *int) SchemaMetaData {
	return &schemaMetaData{
		Total_Records:    totalRecords,
		Records_Per_Page: recordsPerPage,
		Total_Pages:      totalPages,
		Current_Page:     currentPage,
		Next_Page:        nextPage,
		Prev_Page:        prevPage,
	}
}

// Receiver functions
func (s *schemaMetaData) CalculateCurrentlyShowingRecords() int {
	// If there's a next page, then this page is fully filled
	if s.Next_Page != nil {
		return s.Records_Per_Page
	}

	// On the last page or if there's only one page
	remainder := s.Total_Records % int64(s.Records_Per_Page)
	if remainder == 0 {
		return s.Records_Per_Page
	}
	return int(remainder)
}

// Returns a copy of the schema meta data
func (s *schemaMetaData) GetMetaData() schemaMetaData {
	return schemaMetaData{
		Total_Records:    s.Total_Records,
		Records_Per_Page: s.Records_Per_Page,
		Total_Pages:      s.Total_Pages,
		Current_Page:     s.Current_Page,
		Next_Page:        s.Next_Page,
		Prev_Page:        s.Prev_Page,
	}
}

// Extended Schema Meta Data
// Not used in application Only used in admin panel for rendering pagination
type ExtendedSchemaMetaData interface {
}

// Extended for admin pagination rendering
type extendedSchemaMetaData struct {
	Total_Records    int64 `json:"total_records"`    // Total number of records in the entire dataset
	Records_Per_Page int   `json:"records_per_page"` // Number of records displayed per page
	Current_Page     int   `json:"current_page"`     // Current page number
	Next_Page        *int  `json:"next_page"`        // Next page number (null if there is no next page)
	Prev_Page        *int  `json:"prev_page"`        // Previous page number (null if there is no previous page)
	// Used for pagination
	Total_Pages       int `json:"total_pages"`
	Currently_Showing int `json:"currently_showing"`
}

// Constructor
func NewExtendedSchemaMetaData(schemaMetaData SchemaMetaData, currentlyShowing int) ExtendedSchemaMetaData {
	metaData := schemaMetaData.GetMetaData()
	return &extendedSchemaMetaData{
		Total_Records:     metaData.Total_Records,
		Records_Per_Page:  metaData.Records_Per_Page,
		Current_Page:      metaData.Current_Page,
		Next_Page:         metaData.Next_Page,
		Prev_Page:         metaData.Prev_Page,
		Total_Pages:       metaData.Total_Pages,
		Currently_Showing: currentlyShowing,
	}
}
