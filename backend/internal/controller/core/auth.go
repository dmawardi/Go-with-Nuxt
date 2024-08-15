package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"
	"github.com/go-chi/chi/v5"
)

type AuthPolicyController interface {
	// Policies
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByResource(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	// Roles
	FindAllRoles(w http.ResponseWriter, r *http.Request)
	AssignUserRole(w http.ResponseWriter, r *http.Request)
	CreateRole(w http.ResponseWriter, r *http.Request)
	// Inheritance
	FindAllRoleInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritance(w http.ResponseWriter, r *http.Request)
	DeleteInheritance(w http.ResponseWriter, r *http.Request)
}

type authPolicyController struct {
	service coreservices.AuthPolicyService
}

func NewAuthPolicyController(service coreservices.AuthPolicyService) AuthPolicyController {
	return &authPolicyController{service}
}

// API/POLICY

// POLICIES
//

// @Summary      Finds a list of authorization policies
// @Description  Accepts search (added as non-case sensitive LIKE) as query parameters and returns matching policies
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        search   query      string  false  "search (added to all string conditions as LIKE SQL search)"
// @Success      200 {object} []map[string]interface{}
// @Failure      400 {string} string "Can't find policies"
// @Router       /auth [get]
// @Security BearerToken
func (c authPolicyController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab search query
	searchQuery := r.URL.Query().Get("search")
	// Find all
	policies, err := c.service.FindAll(searchQuery)
	if err != nil {
		http.Error(w, "Can't find policies", http.StatusBadRequest)
		return
	}

	// Return
	request.WriteAsJSON(w, policies)
}

// @Summary      Finds a list of authorization policies by resource
// @Description  Accepts resource as a URL parameter (in slug form: replace all '/' with '-') and returns matching policies associated with resource
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        policy-slug   path      string  true  "policy-slug"
// @Success      200 {object} []map[string]interface{}
// @Failure      400 {string} string "Can't find policies for resource"
// @Router       /auth/{policy-slug} [get]
// @Security BearerToken
func (c authPolicyController) FindByResource(w http.ResponseWriter, r *http.Request) {
	// Grab search query
	policyResource := chi.URLParam(r, "policy-slug")

	// Unslugify
	policyResource = webapi.Unslugify(policyResource)

	// Find all
	policies, err := c.service.FindByResource(policyResource)
	if err != nil || len(policies) == 0 {
		http.Error(w, "Can't find policies for resource", http.StatusBadRequest)
		return
	}
	// Return
	request.WriteAsJSON(w, policies)
}

// @Summary      Deletes an authorization policy
// @Description  Accepts a policy as a JSON body and deletes the policy
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        policy   body      models.PolicyRule  true  "policy"
// @Success      200 {string} string "Policy deletion successful!"
// @Failure      400 {string} string "Can't delete policy"
// @Router       /auth [delete]
// @Security BearerToken
func (c authPolicyController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.PolicyRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(pol)
	if err != nil {
		fmt.Printf("Error deleting policy: %v\n", err)
		http.Error(w, "Can't delete policy", http.StatusBadRequest)
		return
	}
	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Policy deletion successful!"))
}

// @Summary      Creates an authorization policy
// @Description  Accepts a policy as a JSON body and creates the policy
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        policy   body      models.PolicyRule  true  "policy"
// @Success      201 {string} string "Policy creation successful!"
// @Failure      400 {string} string "Can't create policy"
// @Router       /auth [post]
// @Security BearerToken
func (c authPolicyController) Create(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.PolicyRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create the policy
	err = c.service.Create(pol)
	if err != nil {
		http.Error(w, "Can't create policy", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Policy creation successful!"))
}

// @Summary      Updates an authorization policy
// @Description  Accepts a policy as a JSON body and updates the policy
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        policy   body      models.UpdateCasbinRule  true  "policy"
// @Success      201 {string} string "Policy update successful!"
// @Failure      400 {string} string "Can't update policy"
// @Router       /auth [put]
// @Security BearerToken
func (c authPolicyController) Update(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.UpdateCasbinRule
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	err = c.service.Update(pol.OldPolicy, pol.NewPolicy)
	if err != nil {
		fmt.Printf("Error updating policy: %v\n", err)
		http.Error(w, "Can't update policy", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Policy update successful!"))
}

// ROLES
//

// @Summary      Finds a list of roles
// @Description  Returns a list of roles
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Success      200 {object} []string
// @Failure      400 {string} string "Can't find roles"
// @Router       /auth/roles [get]
// @Security BearerToken
func (c authPolicyController) FindAllRoles(w http.ResponseWriter, r *http.Request) {
	// Find all roles
	roles, err := c.service.FindAllRoles()
	if err != nil {
		http.Error(w, "Can't find roles", http.StatusBadRequest)
		return
	}
	// Return posts
	request.WriteAsJSON(w, roles)
}

// @Summary      Assigns a role to a user
// @Description  Accepts a user_id and role as a JSON body and assigns the role to the user
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        RoleAssignment   body      models.CasbinRoleAssignment  true  "Role Assignment"
// @Success      200 {string} string "User assigned role successfully!"
// @Failure      400 {string} string "Can't assign user"
// @Router       /auth/roles [put]
// @Security BearerToken
func (c authPolicyController) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.CasbinRoleAssignment
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	success, err := c.service.AssignUserRole(pol.UserId, pol.Role)
	if err != nil {
		http.Error(w, "Can't assign user", http.StatusBadRequest)
		return
	}
	if !*success {
		http.Error(w, "Can't assign user", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User assigned role successfully!"))
}

// @Summary      Creates a new role and assigns a user to it
// @Description  Accepts a user_id and role as a JSON body and assigns the new role to the associated user
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        RoleAssignment   body      models.CasbinRoleAssignment  true  "Role Assignment"
// @Success      200 {string} string "Role creation successful!"
// @Failure      400 {string} string "Can't create role"
// @Router       /auth/roles [post]
// @Security BearerToken
func (c authPolicyController) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.CasbinRoleAssignment
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	success, err := c.service.CreateRole(pol.UserId, pol.Role)
	if err != nil {
		http.Error(w, "Can't create role", http.StatusBadRequest)
		return
	}
	if !*success {
		http.Error(w, "Can't create role", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Role creation successful!"))
}

// INHERITANCE
//

// // @Summary      Finds a list of role inheritance
// // @Description  Returns a list of role inheritance policies
// // @Tags         Authorization
// // @Accept       json
// // @Produce      json
// // @Success      200 {object} []map[string]string
// // @Failure      400 {string} string "Can't find roles"
// // @Router       /auth/inheritance [get]
// // @Security BearerToken
func (c authPolicyController) FindAllRoleInheritance(w http.ResponseWriter, r *http.Request) {
	// Find all roles
	roles, err := c.service.FindAllRoleInheritance()
	if err != nil {
		http.Error(w, "Can't find roles", http.StatusBadRequest)
		return
	}
	// Return posts
	request.WriteAsJSON(w, roles)
}

// // @Summary      Creates a role inheritance
// // @Description  Accepts a role inheritance as a JSON body and creates the role inheritance
// // @Tags         Authorization
// // @Accept       json
// // @Produce      json
// // @Param        inheritance   body      models.GRecord  true  "Inheritance Record"
// // @Success      201 {string} string "Inheritance creation successful!"
// // @Failure      400 {string} string "Can't create inheritance"
// // @Router       /auth/inheritance [post]
// // @Security BearerToken
func (c authPolicyController) CreateInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.GRecord
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	err = c.service.CreateInheritance(pol)
	if err != nil {
		http.Error(w, "Can't create inheritance", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inheritance creation successful!"))
}

// // @Summary      Deletes a role inheritance
// // @Description  Accepts a role inheritance as a JSON body and deletes the role inheritance
// // @Tags         Authorization
// // @Accept       json
// // @Produce      json
// // @Param        inheritance   body      models.GRecord  true  "Inheritance Record"
// // @Success      200 {string} string "Inheritance deletion successful!"
// // @Failure      400 {string} string "Can't delete inheritance"
// // @Router       /auth/inheritance [delete]
// // @Security BearerToken
func (c authPolicyController) DeleteInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab request body as models.CasbinRule
	var pol models.GRecord
	err := json.NewDecoder(r.Body).Decode(&pol)
	if err != nil {
		http.Error(w, "Invalid policy", http.StatusBadRequest)
		return
	}

	// Validate the incoming DTO
	pass, valErrors := request.GoValidateStruct(&pol)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		request.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	err = c.service.DeleteInheritance(pol)
	if err != nil {
		http.Error(w, "Can't delete inheritance", http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inheritance deletion successful!"))
}
