package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dmawardi/Go-Template/internal/auth"
	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
)

// Admin core controller (non-schema related routes)
type AdminCoreController interface {
	// Admin home handler
	Home(w http.ResponseWriter, r *http.Request)
	// Admin login handler
	Login(w http.ResponseWriter, r *http.Request)
	// Admin logout handler
	Logout(w http.ResponseWriter, r *http.Request)
	// Change Password handler
	ChangePassword(w http.ResponseWriter, r *http.Request)
	// Change Password Success handler
	ChangePasswordSuccess(w http.ResponseWriter, r *http.Request)
	// Admin redirect handler
	AdminRedirectBasedOnLoginStatus(w http.ResponseWriter, r *http.Request)
}

type adminCoreController struct {
	service coreservices.UserService
}

// Constructor
func NewAdminCoreController(userService coreservices.UserService) AdminCoreController {
	return &adminCoreController{userService}
}

// RECEIVER FUNCTIONS
// Admin home page
func (c adminCoreController) Home(w http.ResponseWriter, r *http.Request) {
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", PageRenderData{
		SidebarList: sidebar,
		PageType: PageType{
			HomePage: true,
		},
		SectionTitle:  "Welcome to the Admin Panel",
		HeaderSection: header,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Admin login page
func (c adminCoreController) Login(w http.ResponseWriter, r *http.Request) {
	// Init token string
	tokenString := ""
	loginErrorMsg := ""

	// Generate form
	loginForm := c.generateLoginForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form data
		loginMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// Build to login struct
		login := models.Login{
			Email:    loginMap["email"],
			Password: loginMap["password"],
		}
		// Validate form data
		// Validate struct
		pass, valErrors := request.GoValidateStruct(login)

		// If validation passes
		if pass {
			// Login user
			tokenString, err = c.service.LoginUser(&login)
			if err == nil {
				// Set token in cookie
				auth.CreateAndSetHeaderCookie(w, tokenString)

				// Redirect or render a success message
				http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
				return
			}
			// Else if login fails
			fmt.Printf("Error logging in for email: %s\n", login.Email)
			loginErrorMsg = "Invalid email or password"
		}

		// Else if validation fails
		// Populate form field errors
		SetValidationErrorsInForm(loginForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Printf("formFieldMap: %v\nLogin form: %+v", formFieldMap, loginForm)
		// Populate previously entered values (Avoids password)
		err = populateFormValuesWithSubmittedFormMap(&loginForm, formFieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "login.go.tmpl", PageRenderData{
		// The section title is used on this page, to display login errors
		SectionTitle: loginErrorMsg,
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/login",
				FormMethod: "POST",
			},
			FormFields: loginForm,
		},
		HeaderSection: header,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Admin logout page
func (c adminCoreController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token", // Use the name of your auth cookie
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
	})

	// Redirect to the login page, or return a success message
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// Change Password
func (c adminCoreController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Init notification string
	notification := ""
	// Generate form
	passwordForm := c.generateChangePasswordForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form data
		form, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		changePassword := models.ChangePassword{
			CurrentPassword:    form["currentPassword"],
			NewPassword:        form["newPassword"],
			ConfirmNewPassword: form["confirmNewPassword"],
		}
		// // Validate form data
		// // Validate struct
		pass, valErrors := request.GoValidateStruct(changePassword)

		// If validation passes
		if pass {
			// Perform password check
			// First, grab token from cookie
			// Validate the token
			tokenData, err := auth.ValidateAndParseToken(r)
			// If error detected
			if err != nil {
				http.Error(w, "Error parsing authentication token", http.StatusForbidden)
				return
			}
			// convert tokenData.UserID to int
			userID, err := strconv.Atoi(tokenData.UserID)
			if err != nil {
				// handle error
				fmt.Println("Error:", err)
			}

			// Find user using id found in token
			passMatch := c.service.CheckPasswordMatch(userID, []byte(changePassword.CurrentPassword))

			// If pasword match error is nil, and new password matches confirm new password
			if changePassword.NewPassword == changePassword.ConfirmNewPassword && passMatch {

				// Update the user's password
				_, err = c.service.Update(userID, &models.UpdateUser{Password: changePassword.ConfirmNewPassword})
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				// Redirect or render a success message
				http.Redirect(w, r, "/admin/change-password-success", http.StatusSeeOther)
				return
			} else {
				notification = "Old password does not match or new passwords do not match"
			}
		}

		// If validation fails
		if notification == "" {
			if changePassword.NewPassword != changePassword.ConfirmNewPassword {
				notification = "New passwords do not match"
			} else {
				// Else if validation fails, extract errors and manipulato for display
				newPasswordErrors := valErrors.Validation_errors["new_password"]
				errorSplit := strings.Split(newPasswordErrors[0], "does")
				errorString := fmt.Sprintf("Does%s", errorSplit[1])

				// Assign to notification
				notification = errorString
			}
		}
	}
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", PageRenderData{
		SectionTitle: "Change Password",
		PageTitle:    "Change Password",
		// The section detail is used on this page, to display login errors
		SectionDetail: template.HTML("<p>" + notification + "</p>"),
		PageType: PageType{
			EditPage: true,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/change-password",
				FormMethod: "POST",
			},
			FormFields: passwordForm,
		},
		HeaderSection: header,
		SidebarList:   sidebar,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminCoreController) ChangePasswordSuccess(w http.ResponseWriter, r *http.Request) {
	serveAdminSuccess(w, "Change Password Success - Admin", "Change Password Success")
}

// Redirect
func (c adminCoreController) AdminRedirectBasedOnLoginStatus(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ValidateAndParseToken(r)
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
		return
	}
}

// Form generators
func (c adminCoreController) generateLoginForm() []FormField {
	return []FormField{
		{DbLabel: "email", Label: "Email", Name: "email", Placeholder: "", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "password", Label: "Password", Name: "password", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
	}
}
func (c adminCoreController) generateChangePasswordForm() []FormField {
	return []FormField{
		{DbLabel: "Password", Label: "Current Password", Name: "currentPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "New Password", Name: "newPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Confirm new Password", Name: "confirmNewPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
	}
}
