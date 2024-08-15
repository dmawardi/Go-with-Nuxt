package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	"github.com/go-chi/chi/v5"
)

type UserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// API/ME
	GetMyUserDetails(w http.ResponseWriter, r *http.Request)
	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	// Login
	Login(w http.ResponseWriter, r *http.Request)
	// Reset password
	ResetPassword(w http.ResponseWriter, r *http.Request)
	// Email Verification
	ResendVerificationEmail(w http.ResponseWriter, r *http.Request)
	EmailVerification(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	service coreservices.UserService
}

func NewUserController(service coreservices.UserService) UserController {
	return &userController{service}
}

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"age": "int", "name": "string", "active": "bool"}
func UserConditionQueryParams() map[string]string {
	return map[string]string{
		"email":    "string",
		"name":     "string",
		"username": "string",
		"verified": "bool",
		"role":     "string",
	}
}

// API/USERS
// @Summary      Find a list of users
// @Description  Accepts limit, offset, order, search (added as non-case sensitive LIKE) and field names (eg. email=) query parameters to find a list of users. Search is applied to all string fields.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        offset   query      int  false  "offset"
// @Param        order   query      int  false  "order by eg. (asc) "id" (desc) "id_desc" )"
// @Param        search   query      string  false  "search (added to all string conditions as LIKE SQL search)"
// @Param        email query string false "email"
// @Param        name query string false "name"
// @Param        username query string false "username"
// @Param        verified query bool false "verified"
// @Param        role query string false "role"
// @Success      200 {object} models.PaginatedUsersWithRole
// @Failure      400 {string} string "Can't find users"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Failure      400 {string} string "Error extracting query params"
// @Router       /users [get]
// @Security BearerToken
func (c userController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab basic query params set defaults as needed
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := UserConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Query database for all users using query params
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Can't find users", http.StatusBadRequest)
		return
	}
	err = request.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find users", http.StatusBadRequest)
		fmt.Println("error writing users to response: ", err)
		return
	}
}

// @Summary      Find User
// @Description  Find a user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} models.UserWithRole
// @Failure      400 {string} string "Can't find user with ID: {id}"
// @Router       /users/{id} [get]
// @Security BearerToken
func (c userController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find user with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = request.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find user with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// @Summary      Create User
// @Description  Creates a new user
// @Tags         User
// @Accept       json
// @Produce      plain
// @Param        user body models.CreateUser true "New User"
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Success      201 {string} string "User creation successful!"
// @Failure      400 {string} string "User creation failed."
// @Router       /users [post]
func (c userController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var toCreate models.CreateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&toCreate)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes
	// Check if token is present
	token, err := auth.ValidateAndParseToken(r)
	// If not found or not admin
	if err != nil || token.Role != "admin" {
		// Disallow assignments to a user's role
		toCreate.Role = ""
	}

	// Create user
	_, createErr := c.service.Create(&toCreate)
	if createErr != nil {
		http.Error(w, "User creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send user success message in body
	w.Write([]byte("User creation successful!"))
}

// @Summary      Update User
// @Description  Updates an existing user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body models.UpdateUser true "Update User"
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} models.UserWithRole
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Failure      400 {string} string "Failed user update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /users/{id} [put]
// @Security BearerToken
func (c userController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var toUpdate models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&toUpdate)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&toUpdate)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Update user
	updated, createErr := c.service.Update(idParameter, &toUpdate)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed user update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write user to output
	err = request.WriteAsJSON(w, updated)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// @Summary      Delete User
// @Description  Deletes an existing user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed user deletion"
// @Router       /users/{id} [delete]
// @Security BearerToken
func (c userController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed user deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}

// API/ME
//
// @Summary      Update my profile
// @Description  Updates the currently logged in user
// @Tags         My Profile
// @Accept       json
// @Produce      json
// @Param        user body models.UpdateUser true "Update User"
// @Success      200 {object} models.UserWithRole
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Failure      400 {string} string "Failed user update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Failure      400 {string} string "Bad request"
// @Router       /me [put]
// @Security BearerToken
func (c userController) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var toUpdate models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&toUpdate)
	if err != nil {
		fmt.Println("Decoding error: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	toUpdate.Role = "" // Prevent role from being updated

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&toUpdate)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Extract the user's id from their authentication token
	tokenData, err := auth.ValidateAndParseToken(r)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
		return
	}
	// Convert to int
	userId, err := strconv.Atoi(tokenData.UserID)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
		return
	}

	// Update user
	updated, createErr := c.service.Update(userId, &toUpdate)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed user update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write updated user to output
	err = request.WriteAsJSON(w, updated)
	if err != nil {
		fmt.Println("Error writing to JSON", err)
		return
	}
}

// @Summary      Get my user profile details
// @Description  Return my user details
// @Tags         My Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} models.UserWithRole
// @Failure      400 {string} string "Can't find user details"
// @Failure      403 {string} string "Error parsing authentication token"
// @Router       /me [get]
// @Security BearerToken
func (c userController) GetMyUserDetails(w http.ResponseWriter, r *http.Request) {
	// Grab ID from cookie
	// Validate the token
	tokenData, err := auth.ValidateAndParseToken(r)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		return
	}

	// Convert to int
	idParameter, err := strconv.Atoi(tokenData.UserID)
	// If error detected
	if err != nil {
		http.Error(w, "Error parsing authentication token", http.StatusForbidden)
		return
	}

	// Find user by id from cookie
	found, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		return
	}

	// Write found user data to Response
	err = request.WriteAsJSON(w, found)
	if err != nil {
		http.Error(w, "Can't find user details", http.StatusBadRequest)
		return
	}
}

// Login
// Handler to login with existing user
// @Summary      Login
// @Description  Log in to user account
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param        user body models.Login true "Login Form"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Failure      401 {string} string "Invalid Credentials"
// @Failure      405 {string} string "Method not supported"
// @Router       /users/login [post]
func (c userController) Login(w http.ResponseWriter, r *http.Request) {
	// Deny any request that is not a post
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Init models for decoding
	var login models.Login
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&login)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through
	tokenString, err := c.service.LoginUser(&login)
	if err != nil {
		fmt.Printf("Error logging in: %s", err)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// Build login response
	var loginResponse = models.LoginResponse{Token: tokenString}
	// Send to user in body
	request.WriteAsJSON(w, loginResponse)
}

// Reset password
// Handler to reset password
// @Summary      Reset password
// @Description  Reset password
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param        email body models.ResetPasswordAndEmailVerification true "Reset Password Form"
// @Success      200 {string} string "Password reset request successful!"
// @Failure      400 {string} string "Password reset request failed"
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Router       /users/forgot-password [post]
func (c userController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Grab email from request body
	var resetPassword models.ResetPasswordAndEmailVerification
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&resetPassword)
	if err != nil {
		fmt.Println("Decoding error: ", err)
		http.Error(w, "Password reset request failed", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&resetPassword)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through
	err = c.service.ResetPasswordAndSendEmail(resetPassword.Email)
	if err != nil {
		http.Error(w, "Password reset request failed", http.StatusBadRequest)
		return
	}

	// Else
	request.WriteAsJSON(w, "Password reset request successful!")
}

// Email Verification
// @Summary      Email Verification
// @Description  Email Verification
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param		token path string true "Token"
// @Success      200 {string} string "Email verified successfully"
// @Failure      400 {string} string "Token is required"
// @Failure      401 {string} string "Invalid or expired token"
// @Router       /users/verify-email/{token} [get]
func (c userController) EmailVerification(w http.ResponseWriter, r *http.Request) {
	// The token is expected to be in the query string, e.g., /verify-email?token=12345
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Call the service to verify the token
	err := c.service.VerifyEmailCode(token)
	if err != nil {
		fmt.Printf("Error verifying email: %s", err)
		// Handle the error
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Token is valid; you might want to redirect the user to a confirmation page or back to the app
	w.WriteHeader(http.StatusOK)
	request.WriteAsJSON(w, "Email verified successfully")
}

// Send Verification Email
// @Summary      Send Verification Email
// @Description  Send Verification Email
// @Tags         Login
// @Accept       json
// @Produce      json
// @Success      200 {string} string "Email sent successfully"
// @Failure      401 {string} string "Invalid email"
// @Failure      401 {string} string "Email already verified"
// @Failure      400 {string} string "Verification request failed"
// @Failure      400 {object} models.ValidationError "Validation Errors"
// @Router       /users/send-verification-email [post]
// @Security BearerToken
func (c userController) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	tokenData, err := auth.ValidateAndParseToken(r)
	if err != nil {
		http.Error(w, "Authentication Token not detected", http.StatusForbidden)
		return
	}

	// If validation passes
	found, err := c.service.FindByEmail(tokenData.Email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	// If user is already verified
	if *found.Verified {
		http.Error(w, "Email already verified", http.StatusUnauthorized)
		return
	}

	// Call the service to resend a verification email for the associated user
	err = c.service.ResendVerificationEmail(int(found.ID))
	if err != nil {
		// Handle the error
		http.Error(w, "Error sending verification email", http.StatusUnauthorized)
		return
	}

	// Write successful response
	w.WriteHeader(http.StatusOK)
	request.WriteAsJSON(w, "Verification request sent successfully")
}
