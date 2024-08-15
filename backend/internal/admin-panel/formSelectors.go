package adminpanel

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
)

// Form Selectors
// For role selection in form
// Default Records Displayed on find all pages
var recordsPerPage = []int{10, 25, 50, 100}

var possibleActions = []string{"create", "read", "update", "delete"}


// Form Selectors
func RoleSelection() []FormFieldSelector {
	roles, err := app.Policy.Service.(coreservices.AuthPolicyService).FindAllRoles()
	if err != nil {
		// Return default selector
		return []FormFieldSelector{
			{Value: "role:user", Label: "User", Selected: true},
			{Value: "role:admin", Label: "Admin", Selected: false},
			{Value: "role:moderator", Label: "Moderator", Selected: false},
		}
	}
	// Init form field selector
	var roleSelector []FormFieldSelector
	// Build []FormFieldSelector from []string
	for _, r := range roles {
		roleSelector = append(roleSelector, FormFieldSelector{Value: r, Label: utility.CapitalizeFirstLetter(r)})
	}

	// Set basic default as user
	setDefaultSelected(roleSelector, "user")

	return roleSelector
}
func ActionSelection() []FormFieldSelector {
	return []FormFieldSelector{
		{Value: "create", Label: "Create", Selected: true},
		{Value: "read", Label: "Read", Selected: false},
		{Value: "update", Label: "Update", Selected: false},
		{Value: "delete", Label: "Delete", Selected: false},
	}
}
func UserSelection() []FormFieldSelector {
	var users []db.User
	// Query all users
	result := app.DbClient.Select("id, username").Find(&users)
	if result.Error != nil {
		fmt.Printf("Error finding users: %v\n", result.Error)
		return nil
	}

	// Init
	var selector []FormFieldSelector
	// Build []FormFieldSelector from []string DB output
	for _, user := range users {
		selector = append(selector, FormFieldSelector{Value: fmt.Sprint(user.ID), Label: utility.CapitalizeFirstLetter(user.Username)})
	}

	return selector
}

// Helpers
// Takes a slice of FormFieldSelector and sets the Selected field to true for the value that matches valueToSelect
func setDefaultSelected(selector []FormFieldSelector, valueToSelect string) {
	for i, s := range selector {
		if s.Value == valueToSelect {
			selector[i].Selected = true
		} else {
			selector[i].Selected = false
		}
	}
}
