package adminpanel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/dmawardi/Go-Template/internal/db"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/go-chi/chi/v5"
)

// Table headers to show on find all page
var actionTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int", Sortable: true},
	{Label: "Admin", ColumnSortLabel: "admin_id", Pointer: false, DataType: "int", Sortable: false},
	{Label: "Description", ColumnSortLabel: "description", Pointer: false, DataType: "string", Sortable: true},
	{Label: "EntityType", ColumnSortLabel: "entity_type", Pointer: false, DataType: "string"},
	{Label: "EntityID", ColumnSortLabel: "entity_id", Pointer: true, DataType: "string", Sortable: true},
}

func NewAdminActionController(service webapi.ActionService) AdminActionController {
	return &adminActionController{
		service: service,
		// Use values from above
		adminHomeUrl:     "/admin/actions",
		schemaName:       "Action",
		pluralSchemaName: "Actions",
		tableHeaders:     actionTableHeaders,
	}
}

type adminActionController struct {
	service webapi.ActionService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders []TableHeader
}

type AdminActionController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request)
}

func (c adminActionController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := core.ActionConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Grab all items
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Convert data to AdminPanelSchema
	schemaSlice := *found.Data
	var adminSchemaSlice []models.AdminPanelSchema
	for _, item := range schemaSlice {
		// Append to schemaSlice
		adminSchemaSlice = append(adminSchemaSlice, item)
	}

	// Build the table data
	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.adminHomeUrl, c.tableHeaders, false)

	// Generate Find All page render data
	data := GenerateFindAllRenderData(tableData, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, searchQuery)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminActionController) View(w http.ResponseWriter, r *http.Request) {
	// Init new User Edit form
	editForm := c.generateEditForm()

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find Current record
	found := &db.Action{}
	// Search for by ID and store in found
	found, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", c.schemaName), http.StatusNotFound)
		return
	}

	// Convert db struct to map for placeholder population
	currentData := getValuesUsingFieldMap(*found)
	// Populate form field placeholders with data from database
	err = populateValuessWithDBData(&editForm, currentData)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		return
	}

	data := GenerateEditRenderData(editForm, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, stringParameter, false)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Form generation
// Used to build Edit form
func (c adminActionController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "AdminID", Label: "Admin ID", Name: "admin", Placeholder: "Enter admin", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "ActionType", Label: "Action Type", Name: "action_type", Placeholder: "Enter action type", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "EntityType", Label: "Entity Type", Name: "entity_type", Placeholder: "Enter entity type", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "EntityID", Label: "Entity ID", Name: "entity_id", Placeholder: "Enter entity id", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "IPAddress", Label: "IP Address", Name: "ip_address", Placeholder: "Enter ip address", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "Changes", Label: "Changes", Name: "changes", Placeholder: "Enter changes", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "Description", Label: "Description", Name: "description", Placeholder: "Enter description", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},

		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}
