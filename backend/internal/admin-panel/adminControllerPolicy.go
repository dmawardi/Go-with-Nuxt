package adminpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/dmawardi/Go-Template/internal/helpers"
	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	"github.com/go-chi/chi/v5"
)

// Table headers to show on find all pages
var authPolicyTableHeaders = []TableHeader{
	{Label: "resource", ColumnSortLabel: "resource", Pointer: false, DataType: "string"},
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string"},
	{Label: "action", ColumnSortLabel: "action", Pointer: false, DataType: "string"},
}

var inheritanceTableHeaders = []TableHeader{
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string", Sortable: false},
	{Label: "inherits_from", ColumnSortLabel: "inherits_from", Pointer: false, DataType: "string", Sortable: false},
}

var roleTableHeaders = []TableHeader{
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string", Sortable: false},
}

// Constructor
func NewAdminAuthPolicyController(service coreservices.AuthPolicyService) AdminAuthPolicyController {
	return &adminAuthPolicyController{
		service: service,
		// Use values from above
		adminHomeUrl:            "/admin/policy",
		schemaName:              "Policy",
		pluralSchemaName:        "Policies",
		tableHeaders:            authPolicyTableHeaders,
		inheritanceTableHeaders: inheritanceTableHeaders,
		roleTableheaders:        roleTableHeaders,
	}
}

type AdminAuthPolicyController interface {
	// Policy
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CreateSuccess(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	// Roles
	FindAllRoles(w http.ResponseWriter, r *http.Request)
	CreateRole(w http.ResponseWriter, r *http.Request)
	CreateRoleSuccess(w http.ResponseWriter, r *http.Request)
	// Inheritance
	FindAllRoleInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritance(w http.ResponseWriter, r *http.Request)
	DeleteInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritanceSuccess(w http.ResponseWriter, r *http.Request)
	DeleteInheritanceSuccess(w http.ResponseWriter, r *http.Request)
	// For sidebar
	ObtainURLDetails() models.URLDetails
}
type adminAuthPolicyController struct {
	service coreservices.AuthPolicyService
	// For link generation
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders            []TableHeader
	inheritanceTableHeaders []TableHeader
	roleTableheaders        []TableHeader
}

// POLICIES
func (c adminAuthPolicyController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab search parameter
	searchQuery := r.URL.Query().Get("search")

	// Find all policies from database
	groupsSlice, err := c.service.FindAll(searchQuery)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}

	// Build the roles table data
	tableData := BuildPolicyTableData(groupsSlice, c.adminHomeUrl, c.tableHeaders)
	// Add the row span attribute to the table based on resource grouping (for UX)
	editTableDataRowSpan(tableData.TableRows)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl,
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) Edit(w http.ResponseWriter, r *http.Request) {
	// Grab slug from URL
	policySlug := chi.URLParam(r, "id")
	// Unslug
	policyUnslug := webapi.Unslugify(policySlug)
	// Detect request method
	method := r.Method

	// If form is being submitted (method = POST)
	if method == "POST" || method == "DELETE" {
		// Init new policy
		pol := &models.PolicyRule{}
		// Decode request body as JSON and store
		err := json.NewDecoder(r.Body).Decode(&pol)
		if err != nil {
			http.Error(w, "Invalid policy", http.StatusBadRequest)
			return
		}

		// Validate the incoming DTO
		pass, _ := request.GoValidateStruct(pol)

		// If passes
		if pass {
			// and method is post
			if method == "POST" {
				// Create policy
				err = c.service.Create(*pol)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
					return
				}
			} else if method == "DELETE" {
				// Else if method is delete
				// Delete policy
				err = c.service.Delete(*pol)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
					return
				}
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/%s", c.adminHomeUrl, policySlug), http.StatusSeeOther)
			return
		}
	}

	// If not POST, ie. GET
	// Find all policies
	found, err := c.service.FindByResource(policyUnslug)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Prepare policies for rendering
	policies := buildEditPolicyTable(found)

	// Init new role selector values
	rolesCurrentlyInPolicy := RoleSelection()
	// Remove roles that are already in the policy
	rolesCurrentlyInPolicy = genRolesLeftOnlySelection(policies, rolesCurrentlyInPolicy)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Edit %s: %s", c.schemaName, policyUnslug),
		SectionTitle: fmt.Sprintf("Edit %s: %s", c.schemaName, policyUnslug),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   true,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		PolicySection: PolicySection{
			FocusedPolicies: policies,
			PolicyResource:  policyUnslug,
			Selectors: PolicyEditSelectors{
				RoleSelection:   rolesCurrentlyInPolicy,
				ActionSelection: ActionSelection()},
		},
		FormData: FormData{
			FormDetails: FormDetails{},
			FormFields:  []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form submission
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Convert to policy rule struct
		toValidate := models.PolicyRule{
			Role:     formFieldMap["role"],
			Resource: formFieldMap["resource"],
			Action:   formFieldMap["action"],
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			err = c.service.Create(toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field with errors
		SetValidationErrorsInForm(createForm, *valErrors)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s", c.schemaName), fmt.Sprintf("%s Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Resource", Label: "Resource", Name: "resource", Placeholder: "eg. '/api/posts'", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "First Role", Name: "role", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
		{DbLabel: "Action", Label: "Action", Name: "action", Placeholder: "", Value: "", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: ActionSelection()},
	}
}

// Roles
func (c adminAuthPolicyController) FindAllRoles(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")

	// Find all with options from database
	rolesSlice, err := c.service.FindAllRoles()
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}

	filteredSlice := []string{}
	// Iterate through roles slice and remove items that do not match search query
	for _, role := range rolesSlice {
		if utility.ContainsString(role, searchQuery) {
			filteredSlice = append(filteredSlice, role)
		}
	}

	// // Build the roles table data
	tableData := BuildRoleTableData(filteredSlice, c.adminHomeUrl, c.roleTableheaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl + "/roles",
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateRoleForm()
	formFieldMap, err := adminpanel.ParseFormToMap(r)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Convert to role assignment struct
		toValidate := models.CasbinRoleAssignment{
			Role:   formFieldMap["role"],
			UserId: formFieldMap["user"],
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			success, err := c.service.CreateRole(toValidate.UserId, toValidate.Role)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			if !*success {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create-role/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s Role", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s Role", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create-role", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) CreateRoleSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s Role", c.schemaName), fmt.Sprintf("%s Role Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) generateCreateRoleForm() []FormField {
	return []FormField{
		{DbLabel: "Role", Label: "New Role Name", Name: "role", Placeholder: "eg. 'Moderator'", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "First Member", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: UserSelection()},
	}
}

// Inheritance
func (c adminAuthPolicyController) FindAllRoleInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")

	// Find all with options from database
	inheritanceSlice, err := c.service.FindAllRoleInheritance()
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// // Filter by search query
	inheritanceSlice = helpers.SearchGRecords(inheritanceSlice, searchQuery)

	// Sort by resource alphabetically
	sort.Slice(inheritanceSlice, func(i, j int) bool {
		// Give two items to compare to alphabetic sorter
		return inheritanceSlice[i].Role < inheritanceSlice[j].Role
	})

	// // Build the roles table data
	tableData := BuildRoleInheritanceTableData(inheritanceSlice, c.adminHomeUrl, c.inheritanceTableHeaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "inheritance",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl + "/inheritance",
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) CreateInheritance(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateInheritanceForm()
	// Extract user form submission
	formFieldMap, err := adminpanel.ParseFormToMap(r)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Convert to role assignment struct
		submittedForm := models.GRecord{
			Role:         formFieldMap["role"],
			InheritsFrom: formFieldMap["inherits_from"],
		}
		// Validate struct
		pass, valErrors := request.GoValidateStruct(submittedForm)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			err := c.service.CreateInheritance(models.GRecord{Role: submittedForm.Role, InheritsFrom: submittedForm.InheritsFrom})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}

			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create-inheritance/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s Inheritance", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s Inheritance", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create-inheritance", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) DeleteInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab params from URL
	inheritSlug := chi.URLParam(r, "inherit-slug")
	// Split by comma to separate the two params
	inheritArray := strings.Split(inheritSlug, ",")
	// Assign individually
	role := inheritArray[0]
	inherits := inheritArray[1]

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Delete user
		err := c.service.DeleteInheritance(models.GRecord{Role: role, InheritsFrom: inherits})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
			return
		}
		// Redirect to success page
		http.Redirect(w, r, fmt.Sprintf("%s/delete-inheritance/success", c.adminHomeUrl), http.StatusSeeOther)
		return
	}

	// GET request
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Delete %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Are you sure you wish to delete: %s?", fmt.Sprintf("%s inherits from %s", role, inherits)),
		SidebarList:  sidebar,
		SchemaHome:   c.adminHomeUrl,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/delete-inheritance/%s", c.adminHomeUrl, inheritSlug),
				FormMethod: "post",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "policy.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) CreateInheritanceSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s Role Inheritance", c.schemaName), fmt.Sprintf("%s Role Inheritance Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) DeleteInheritanceSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Delete %s Inheritance Inheritance", c.schemaName), fmt.Sprintf("%s Inheritance Inheritance Created Successfully!", c.schemaName))
}

// Form
func (c adminAuthPolicyController) generateCreateInheritanceForm() []FormField {
	return []FormField{
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
		{DbLabel: "InheritsFrom", Label: "Inherits from (role)", Name: "inherits_from", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
	}
}

// Basic helper functions

// Used to build standardize controller fields for admin panel sidebar generation
func (c adminAuthPolicyController) ObtainURLDetails() models.URLDetails {
	return models.URLDetails{
		AdminHomeUrl:     c.adminHomeUrl,
		SchemaName:       c.schemaName,
		PluralSchemaName: c.pluralSchemaName,
	}
}

// Search helpers
// Convert the map received from the service to a slice of models.PolicyRule
func buildEditPolicyTable(m []models.PolicyRuleCombinedActions) []PolicyEditDataRow {
	var policies []PolicyEditDataRow
	// Iterate through map of policies
	for _, policy := range m {
		var actions []PolicyActionCell

		// Iterate through actions
		for _, action := range possibleActions {
			// Create policy action cell
			actionToAdd := PolicyActionCell{
				Action: action,
				// Make false as default
				Granted: false,
			}

			// Check if array contains a string
			if utility.ArrayContainsString(policy.Action, action) {
				actionToAdd.Granted = true
			}

			// Append to actions
			actions = append(actions, actionToAdd)
		}

		// Build policy edit row
		policyToAdd := PolicyEditDataRow{
			Role:     policy.Role,
			Resource: policy.Resource,
			Actions:  actions,
		}

		// Append to policies
		policies = append(policies, policyToAdd)
	}
	return policies
}

// Takes a slice of PolicyEditDataRow and role selector and returns a slice of role selector with only missing roles
func genRolesLeftOnlySelection(policies []PolicyEditDataRow, roleSelector []FormFieldSelector) []FormFieldSelector {
	for _, p := range policies {
		// Iterate through roleSelector
		for i, role := range roleSelector {
			// If the role matches
			if role.Value == p.Role {
				// Remove from slice
				roleSelector = append(roleSelector[:i], roleSelector[i+1:]...)
				break
			}
		}
	}
	return roleSelector
}
