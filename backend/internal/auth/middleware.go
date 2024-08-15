package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dmawardi/Go-Template/internal/db"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
)

// Middleware to check whether user is authenticated
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Grab Http Method
		httpMethod := r.Method
		// Determine associated action based on HTTP method
		action := ActionFromMethod(httpMethod)
		// Extract current URL being accessed
		object := webapi.ExtractBasePath(r)

		// Validate the token
		tokenData, err := ValidateAndParseToken(r)
		// If error detected
		if err != nil {
			// Determine Redirect URL based on object
			redirectURL := determineInvalidTokenRedirectURL(object)
			// If redirect URL is  empty
			if redirectURL == "" {
				http.Error(w, "Error parsing authentication token", http.StatusForbidden)
				return
			}
			// Else, redirect to invalid token page
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Enforce RBAC policy and determine if user is authorized to perform action
		allowed := Authorize(tokenData.UserID, object, action)

		// If not allowed
		if !allowed {
			http.Error(w, "Not authorized to perform that action", http.StatusForbidden)
			return
		}

		// Else, allow through
		next.ServeHTTP(w, r)
	})
}

// Dependent on if in admin section or other section, redirect to home associated
func determineInvalidTokenRedirectURL(object string) string {
	// Split string
	stringArray := strings.Split(object, "/")
	// If user currently in admin section
	if stringArray[1] == "admin" {
		// Set redirect URL to admin login page
		return "/admin"
	} else {
		return ""
	}
}

// Middleware to check whether user is authorized
func Authorize(userId, object, action string) bool {
	permissionCheck := false
	// Load Authorization policy from Database
	err := app.Auth.Enforcer.LoadPolicy()
	if err != nil {
		fmt.Println("Error loading RBAC policy in Authorization middleware: ", err)
		return false
	}

	// Enforce policy for user's role using their ID and explicit policy (permissions assigned by user's role)
	permissionCheck, err = app.Auth.Enforcer.Enforce(userId, object, action)
	if err != nil {
		fmt.Print("Failed to enforce RBAC policy in Authorization middleware: ", err, "\nUser ID: ", userId, "\nObject: ", object, "\nAction: ", action, "\n")
		return false
	}
	// Get roles for user
	roles, err := app.Auth.Enforcer.GetRolesForUser(userId)
	if err != nil {
		fmt.Println("Error getting roles for user in Authorization middleware: ", err)
	}
	fmt.Printf("User with ID %s (Role: %s) is accessing %s to %s. Allowed? %v\n", userId, roles, object, action, permissionCheck)

	// Return result of enforcement
	return permissionCheck
}

// Find user in database by email (for authentication)
func FindByEmail(email string) (*db.User, error) {
	// Create an empty ref object of type user
	user := db.User{}
	// Check if user exists in db
	result := app.DbClient.Where("email = ?", email).First(&user)

	// If error detected
	if result.Error != nil {
		fmt.Println("error in finding user in authentication: ", result.Error)
		return nil, result.Error
	}
	// else
	return &user, nil
}
