package adminpanel

import (
	"strconv"

	"github.com/dmawardi/Go-Template/internal/db"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	schemamodels "github.com/dmawardi/Go-Template/internal/models/schemaModels"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
)

func NewAdminPostController(service moduleservices.PostService, actionService webapi.ActionService) models.BasicAdminController{
	return &basicAdminController[db.Post, schemamodels.CreatePost, schemamodels.UpdatePost]{
		Service: service,
		ActionService: actionService,
		// Use values from above
		AdminHomeUrl:     "/admin/posts",
		SchemaName:       "Post",
		PluralSchemaName: "Posts",
		tableHeaders:     []TableHeader{
			{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int", Sortable: true},
			{Label: "Title", ColumnSortLabel: "title", Pointer: false, DataType: "string", Sortable: true},
			{Label: "Body", ColumnSortLabel: "body", Pointer: false, DataType: "string", Sortable: true},
			{Label: "User", ColumnSortLabel: "user", Pointer: false, DataType: "foreign", ForeignKeyRepKeyName: "Username"},
		},
		generateCreateForm: func () []FormField  {
			return []FormField{
				{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
				{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "rich-text-editor", Required: true, Disabled: false, Errors: []ErrorMessage{}},
				{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: UserSelection()},
			}
		},
		generateEditForm: func () []FormField  {
			return []FormField{
				{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
				{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "rich-text-editor", Required: false, Disabled: false, Errors: []ErrorMessage{}},
				{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: UserSelection()},
			}
		},
		prepareSubmittedFormForCreation: func(formFieldMap map[string]string) (*schemamodels.CreatePost, error) {
			// Convert relationsip to int
			userId, err := strconv.Atoi(formFieldMap["user"])
			if err != nil {
				return nil, err
			}
			// Convert submitted form field map to struct for validation/creation
			toValidate := schemamodels.CreatePost{
				Title: formFieldMap["title"],
				Body:  formFieldMap["body"],
				User:  db.User{ID: uint(userId)},
			}
			return &toValidate, nil
		},
		prepareSubmittedFormForUpdate: func(formFieldMap map[string]string) (*schemamodels.UpdatePost, error) {
			// Convert relationsip to int
			userId, err := strconv.Atoi(formFieldMap["user"])
			if err != nil {
				return nil, err
			}
			// Convert submitted form field map to struct for validation/creation
			toValidate := schemamodels.UpdatePost{
				Title: formFieldMap["title"],
				Body:  formFieldMap["body"],
				User:  db.User{ID: uint(userId)},
			}
			return &toValidate, nil
		},
		newEmptySchema: func(params ...uint) *db.Post {
			// If there is a parameter
			if len(params) > 0 {
				// Grab paramter and add to struct as ID for return
				id := params[0]
				return &db.Post{ID: id}
			} else {
				return &db.Post{}
			}
		},
		getIDFromSchema: func(schema *db.Post) uint {
			return schema.ID
		},
		
	}
}