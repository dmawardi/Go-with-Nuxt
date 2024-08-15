package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/models"
)

// PAGE RENDER DATA
// Contains state for the page
type PageRenderData struct {
	// In HEAD
	PageTitle string
	// In BODY
	SectionTitle  string
	SectionDetail template.HTML
	SidebarList   AdminSideBar
	// Schema home used to return to the schema home page from delete
	SchemaHome string // eg. /admin/users/
	// Page type (Used for content selection)
	PageType PageType
	// Form
	FormData  FormData
	TableData TableData
	// Search
	SearchTerm             string
	RecordsPerPageSelector []int
	// Special section data for policies
	PolicySection PolicySection
	HeaderSection HeaderSection
}

// Variables for header
type HeaderSection struct {
	HomeUrl           template.URL
	ViewSiteUrl       template.URL
	ChangePasswordUrl template.URL
	LogOutUrl         template.URL
}

// Variables for policy section
type PolicySection struct {
	FocusedPolicies []PolicyEditDataRow
	PolicyResource  string
	Selectors       PolicyEditSelectors
}

// Page type (Used for dynamic selective rendering)
type PageType struct {
	HomePage    bool
	EditPage    bool
	ViewPage 	bool
	ReadPage    bool
	CreatePage  bool
	DeletePage  bool
	SuccessPage bool
	// Used for policy section
	PolicyMode string // eg. "policy" or "inheritance"
}

// TEMPLATES
//
// Parses all the template files in the templates directory
func ParseAdminTemplates() (*template.Template, error) {
	// Parse the base template
	tmpl := template.New("/internal/admin-panel/templates/layout.go.tmpl")

	// Walk through all files in the templates directory
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If the file is not a directory and has the .html extension
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			// Parse the file
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Return error if there is filepath walk issue
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Function to render the Admin error page to the response
func serveAdminError(w http.ResponseWriter, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Error - Admin",
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
		},
		FormData: FormData{},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Function to render the Admin success page to the response
func serveAdminSuccess(w http.ResponseWriter, pageTitle string, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    pageTitle,
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData:      FormData{},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// SIDEBAR
//
// Build item list for sidebar (Add for every module)
var sidebar = AdminSideBar{
	Main: []sidebarItem{
		// This list is filled upon runtime by GenerateAndSetAdminSidebar
	},
	Auth: BuildAuthSidebarSection(),
}

// Generate and set sidebar list
// Accepts current Admin controller and generates sidebar list based on controllers
func GenerateAndSetAdminSidebar(adminCont AdminPanelController) {
	// Iterate through all controllers and add to sidebar list
	// Get the reflect.Value of the struct.
	valueOfCont := reflect.ValueOf(adminCont)

	// Iterate through the struct fields to build sidebar list for base admin controllers
	for i := 0; i < valueOfCont.NumField(); i++ {
		// Get the field name and value.
		fieldName := valueOfCont.Type().Field(i).Name
		fieldValue := valueOfCont.Field(i).Interface()

		// If not base controller, add to sidebar list
		if fieldName != "Base" && fieldName != "Auth" && fieldName != "ModuleMap" && fieldName != "Action" {
			currentController := ObtainUrlDetailsForBasicAdminController(fieldValue)
			// Create sidebar item
			item := sidebarItem{
				Name:        currentController.PluralSchemaName,
				AddLink:     fmt.Sprintf("%s/create", currentController.AdminHomeUrl),
				FindAllLink: currentController.AdminHomeUrl,
			}

			// append to sidebar list
			sidebar.Main = append(sidebar.Main, item)
		}
	}

	// Iterate through module map of created modules and add to sidebar list
	for _, moduleSet := range adminCont.ModuleMap {
		// Check if admin controller is not nil
		if moduleSet.AdminController != nil {
			currentController := ObtainUrlDetailsForBasicAdminController(moduleSet.AdminController)
			// Create sidebar item
			item := sidebarItem{
				Name:        currentController.PluralSchemaName,
				AddLink:     fmt.Sprintf("%s/create", currentController.AdminHomeUrl),
				FindAllLink: currentController.AdminHomeUrl,
			}

			// append to sidebar list
			sidebar.Main = append(sidebar.Main, item)
		}
	}
}

// Build auth section for sidebar in admin panel
func BuildAuthSidebarSection() []sidebarItem {
	return []sidebarItem{
		{
			Name:        "Permissions",
			FindAllLink: "/admin/policy",
			AddLink:     "/admin/policy/create",
		},
		{
			Name:        "Roles",
			FindAllLink: "/admin/policy/roles",
			AddLink:     "/admin/policy/create-role",
		},
		{
			Name:        "Inheritance",
			FindAllLink: "/admin/policy/inheritance",
			AddLink:     "/admin/policy/create-inheritance",
		},
	}
}

// Page Renders
// 
// Generate the render data for the FindAll page
func GenerateFindAllRenderData(tableData TableData, SchemaName, PluralSchemaName, AdminHomeUrl, searchQuery string) PageRenderData {
	return PageRenderData{
		// Input data
		TableData:              tableData,
		SchemaHome:             AdminHomeUrl,
		SearchTerm:             searchQuery,
		PageTitle:              "Admin: " + PluralSchemaName,
		SectionTitle:           fmt.Sprintf("Select a %s to edit", SchemaName),
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: AdminHomeUrl,
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		// Admin panel standard variables
		HeaderSection: header,
		SidebarList:            sidebar,
		RecordsPerPageSelector: recordsPerPage,
		PageType: PageType{
			ReadPage:   true,
		},
	}
}

// Generate the render data for the Create page
func GenerateCreateRenderData(formFields []FormField, SchemaName, PluralSchemaName, AdminHomeUrl string) PageRenderData {
	return PageRenderData{
		// Input data
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create", AdminHomeUrl),
				FormMethod: "post",
			},
			FormFields: formFields,
		},
		PageTitle:    "Create " + SchemaName,
		SectionTitle: fmt.Sprintf("Create a new %s", SchemaName),
		SchemaHome:   AdminHomeUrl,
		// Admin panel standard variables
		HeaderSection: header,
		SidebarList:   sidebar,
		PageType: PageType{
			CreatePage: true,
		},
	}
}

// Generate the render data for the Edit page
func GenerateEditRenderData(formFields []FormField, SchemaName, PluralSchemaName, AdminHomeUrl, idParameter string, canEdit bool) PageRenderData {
	// Determine which page to show (Edit/View)
	pageTypeToShow := PageType{}
	if canEdit {
		pageTypeToShow.EditPage = true
	} else {
		pageTypeToShow.ViewPage = true
	}
	return PageRenderData{
		// Input data
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/%s", AdminHomeUrl, idParameter),
				FormMethod: "post",
			},
			FormFields: formFields,
		},
		PageTitle:    fmt.Sprintf("Edit %s: %s", SchemaName, idParameter),
		SectionTitle: fmt.Sprintf("Edit %s: %s", SchemaName, idParameter),
		SchemaHome:   AdminHomeUrl,
		// Admin panel standard variables
		HeaderSection: header,
		SidebarList:   sidebar,
		PageType: pageTypeToShow,
	}
}

// Generate the render data for the Delete page
func GenerateDeleteRenderData(SchemaName, PluralSchemaName, AdminHomeUrl, idParameter string) PageRenderData {
	return PageRenderData{
		// Input data
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/delete/%s", AdminHomeUrl, idParameter),
				FormMethod: "post",
			},
			FormFields: []FormField{},
		},
		PageTitle:    fmt.Sprintf("Delete %s", SchemaName),
		SectionTitle: fmt.Sprintf("Are you sure you wish to delete %s: %s?", SchemaName, idParameter),
		SchemaHome:   AdminHomeUrl,
		// Admin panel standard variables
		HeaderSection: header,
		SidebarList:   sidebar,
		PageType: PageType{
			DeletePage: true,
		},
	}
}

// Sidebar helpers
// Uses the ObtainUrlDetails method to get the sidebar details of any Basic Admin Controller type
func ObtainUrlDetailsForBasicAdminController(input interface{}) models.URLDetails {
	// Use reflection to call ObtainUrlDetails method if it exists.
	value := reflect.ValueOf(input)
	// ObtainUrlDetails method
	method := value.MethodByName("ObtainUrlDetails")
	if !method.IsValid() {
		return models.URLDetails{}
	}

	// Call ObtainUrlDetails method
	result := method.Call(nil)

	// Check if result is valid (if it is a BasicAdminController)
	// If it has a result
	if len(result) == 1 {
		// Assign the result as an interface to resultFields
		interfaceFields := result[0].Interface()
		// Assign the fields of the resultFields to sidebarDetails
		sidebarDetails := models.URLDetails{
			AdminHomeUrl:     interfaceFields.(models.URLDetails).AdminHomeUrl,
			SchemaName:       interfaceFields.(models.URLDetails).SchemaName,
			PluralSchemaName: interfaceFields.(models.URLDetails).PluralSchemaName,
		}
		return sidebarDetails
	}

	return models.URLDetails{}
}



// Used for rendering admin sidebar
type sidebarItem struct {
	Name        string
	FindAllLink string
	AddLink     string
}

type AdminSideBar struct {
	Main []sidebarItem
	Auth []sidebarItem
}
