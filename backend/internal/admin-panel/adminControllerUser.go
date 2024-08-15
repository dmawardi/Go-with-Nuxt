package adminpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/dmawardi/Go-Template/internal/helpers"
	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	data "github.com/dmawardi/Go-Template/internal/helpers/data"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/go-chi/chi/v5"
)

// Table headers to show on find all page
var userTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int", Sortable: true},
	{Label: "Username", ColumnSortLabel: "username", Pointer: false, DataType: "string", Sortable: true},
	{Label: "Email", ColumnSortLabel: "email", Pointer: false, DataType: "string", Sortable: true},
	{Label: "Role", ColumnSortLabel: "role", Pointer: false, DataType: "string"},
	{Label: "Verified", ColumnSortLabel: "verified", Pointer: true, DataType: "bool", Sortable: true},
}

func NewAdminUserController(service coreservices.UserService, action webapi.ActionService) AdminUserController {
	return &adminUserController{
		service:       service,
		actionService: action,
		// Use values from above
		adminHomeUrl:     "/admin/users",
		schemaName:       "User",
		pluralSchemaName: "Users",
		tableHeaders:     userTableHeaders,
	}
}

type adminUserController struct {
	service       coreservices.UserService
	actionService webapi.ActionService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders []TableHeader
}

type AdminUserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// Bulk delete (from table)
	BulkDelete(w http.ResponseWriter, r *http.Request)
	// Success pages
	CreateSuccess(w http.ResponseWriter, r *http.Request)
	EditSuccess(w http.ResponseWriter, r *http.Request)
	DeleteSuccess(w http.ResponseWriter, r *http.Request)
	// For sidebar
	ObtainUrlDetails() models.URLDetails
}

func (c adminUserController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := core.UserConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Grab all users from database
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
	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.adminHomeUrl, c.tableHeaders, true)

	// Generate Find All page render data
	data := GenerateFindAllRenderData(tableData, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, searchQuery)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new Create form
	createForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Preparation for validation
		// Extract form submission
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		verified := false
		if formFieldMap["verified"] != "" {
			verified = true
		}
		// Build struct for validation
		toValidate := models.CreateUser{
			Name:     formFieldMap["name"],
			Username: formFieldMap["username"],
			Email:    formFieldMap["email"],
			Password: formFieldMap["password"],
			Role:     formFieldMap["role"],
			Verified: verified,
		}
		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create user
			createdUser, err := c.service.Create(&toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Record action
			err = c.actionService.RecordAction(r, c.schemaName, createdUser.ID, &models.RecordedAction{
				ActionType: "create",
				EntityType: c.schemaName,
				EntityID:   fmt.Sprint(createdUser.ID),
			}, helpers.ChangeLogInput{OldObj: &models.UserWithRole{}, NewObj: createdUser})
			if err != nil {
				fmt.Printf("Error recording action: %s", err)
			}

			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)

		// Populate previously entered values (Avoids password)
		err = populateFormValuesWithSubmittedFormMap(&createForm, formFieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}

	data := GenerateCreateRenderData(createForm, c.schemaName, c.pluralSchemaName, c.adminHomeUrl)

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
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
	found := &models.UserWithRole{}
	// Search for by ID and store in found
	found, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", c.schemaName), http.StatusNotFound)
		return
	}

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		verified := false
		if formFieldMap["verified"] != "" {
			verified = true
		}
		toValidate := models.UpdateUser{
			Name:     formFieldMap["name"],
			Username: formFieldMap["username"],
			Email:    formFieldMap["email"],
			Password: formFieldMap["password"],
			Role:     formFieldMap["role"],
			Verified: verified,
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Update user
			updated, err := c.service.Update(idParameter, &toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error updating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			updated.Password = ""

			// Record action
			err = c.actionService.RecordAction(r, c.schemaName, uint(idParameter), &models.RecordedAction{
				ActionType: "update",
				EntityType: c.schemaName,
				EntityID:   stringParameter,
			}, helpers.ChangeLogInput{OldObj: found, NewObj: updated})
			if err != nil {
				fmt.Printf("Error recording action: %s", err)
			}

			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/edit/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(editForm, *valErrors)

		// Populate previously entered values (Avoids password)
		err = populateFormValuesWithSubmittedFormMap(&editForm, formFieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}

	// If not POST, ie. GET
	// Remove password from struct for field population (Usually, this is done during JSON marshalling)
	found.Password = ""

	// Convert db struct to map for placeholder population
	currentData := getValuesUsingFieldMap(*found)
	// Populate form field placeholders with data from database
	err = populateValuessWithDBData(&editForm, currentData)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		return
	}

	data := GenerateEditRenderData(editForm, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, stringParameter, true)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminUserController) Delete(w http.ResponseWriter, r *http.Request) {
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		serveAdminError(w, "Unable to interpret ID")
		return
	}
	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Delete
		err = c.service.Delete(idParameter)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
			return
		}

		// Record action
		err = c.actionService.RecordAction(r, c.schemaName, uint(idParameter), &models.RecordedAction{
			ActionType: "delete",
			EntityType: c.schemaName,
			EntityID:   stringParameter,
		}, helpers.ChangeLogInput{OldObj: &models.UserWithRole{ID: uint(idParameter)}, NewObj: &models.UserWithRole{}})
		if err != nil {
			fmt.Printf("Error recording action: %s", err)
		}

		// Redirect to success page
		http.Redirect(w, r, fmt.Sprintf("%s/delete/success", c.adminHomeUrl), http.StatusSeeOther)
		return
	}

	data := GenerateDeleteRenderData(c.schemaName, c.pluralSchemaName, c.adminHomeUrl, stringParameter)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminUserController) BulkDelete(w http.ResponseWriter, r *http.Request) {
	// Grab body of request
	// Init
	var listOfIds BulkDeleteRequest

	// Prepare response
	bulkResponse := models.BulkDeleteResponse{
		// Set deleted records to length of selected items
		DeletedRecords: len(listOfIds.SelectedItems),
		Errors:         []error{},
	}

	// Decode request body as JSON and store
	err := json.NewDecoder(r.Body).Decode(&listOfIds)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Convert string slice to int slice
	intIdList, err := data.ConvertStringSliceToIntSlice(listOfIds.SelectedItems)
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		request.WriteAsJSON(w, bulkResponse)
		return
	}

	// Bulk Delete
	err = c.service.BulkDelete(intIdList)
	// If error detected send error response
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		request.WriteAsJSON(w, bulkResponse)
		return
	}
	// Record Bulk delete
	err = c.actionService.RecordBulkDelete(r, c.schemaName, c.pluralSchemaName, intIdList, &models.RecordedAction{
		ActionType: "bulk-delete",
		EntityType: c.schemaName,
		EntityID:   fmt.Sprint(intIdList),
	})
	if err != nil {
		fmt.Printf("Error recording action: %s", err)
	}

	// else if successful
	bulkResponse.Success = true
	request.WriteAsJSON(w, bulkResponse)
}

// Success handlers
func (c adminUserController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s", c.schemaName), fmt.Sprintf("%s Created Successfully!", c.schemaName))
}
func (c adminUserController) EditSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Edit %s", c.schemaName), fmt.Sprintf("%s Updated Successfully!", c.schemaName))
}
func (c adminUserController) DeleteSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Delete %s", c.schemaName), fmt.Sprintf("%s Deleted Successfully!", c.schemaName))
}

// Form generation
// Used to build Create form
func (c adminUserController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "true", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
	}
}

// Used to build Edit form
func (c adminUserController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCode", Label: "Verification Code", Name: "verification_code", Placeholder: "Enter verification code", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCodeExpiry", Label: "Verification Code Expiry", Name: "verification_code_expiry", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}

// Used to build standardize controller fields for admin panel sidebar generation
func (c adminUserController) ObtainUrlDetails() models.URLDetails {
	return models.URLDetails{
		AdminHomeUrl:     c.adminHomeUrl,
		SchemaName:       c.schemaName,
		PluralSchemaName: c.pluralSchemaName,
	}
}
