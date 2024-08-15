package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dmawardi/Go-Template/internal/helpers"
	data "github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestAuthController_FindAll(t *testing.T) {
	testName := "Auth Find all"
	expectedStatusResponse := http.StatusOK
	req, err := helpers.BuildApiRequest("GET", "auth", nil, true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != expectedStatusResponse {
		t.Errorf("%v: got %v want %v: %v", testName,
			status, expectedStatusResponse, rr.Body.String())
	}

	// Convert response JSON to struct
	var body []models.PolicyRuleCombinedActions
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check if type matches
	match := data.CheckSliceType(body, reflect.TypeOf(models.PolicyRuleCombinedActions{}))
	if match == false {
		t.Errorf("Expected %v, got %v", reflect.TypeOf(models.PolicyRuleCombinedActions{}), reflect.TypeOf(body))
	}
}
func TestAuthController_FindByResource(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}
	// Create policy
	err := testModule.auth.serv.Create(policy1)
	if err != nil {
		t.Error(err)
	}
	// Build slug
	slug := webapi.Slugify(policy1.Resource)
	requestUrl := fmt.Sprintf("auth/%s", slug)

	req, err := helpers.BuildApiRequest("GET", requestUrl, nil, true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v", "Auth Find by resource",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []models.PolicyRuleCombinedActions
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check details
	checkPolicyDetails(t, body[0], policy1)

	// Delete policy
	err = testModule.auth.serv.Delete(policy1)
	if err != nil {
		t.Error(err)
	}
}
func TestAuthController_Delete(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}
	// Create policy
	err := testModule.auth.serv.Create(policy1)
	if err != nil {
		t.Error(err)
	}

	// Build slug
	requestUrl := "auth"

	req, err := helpers.BuildApiRequest("DELETE", requestUrl, helpers.BuildReqBody(policy1), true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Delete",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is deleted
	found, err := testModule.auth.serv.FindByResource(policy1.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) > 0 {
		t.Errorf("Expected to not find resource, however, found: %v", found)
	}
}
func TestAuthController_Create(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}

	// Build slug
	requestUrl := "auth"

	req, err := helpers.BuildApiRequest("POST", requestUrl, helpers.BuildReqBody(policy1), true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Create",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is created
	found, err := testModule.auth.serv.FindByResource(policy1.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) == 0 {
		t.Errorf("Expected to find resource, however, not found: %v", found)
	}

	// Check details
	checkPolicyDetails(t, found[0], policy1)

	// Delete policy
	err = testModule.auth.serv.Delete(policy1)
	if err != nil {
		t.Error(err)
	}
}
func TestAuthController_Update(t *testing.T) {
	policy1 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "read",
	}
	policy2 := models.PolicyRule{
		Role:     "admin",
		Resource: "/api/gustav",
		Action:   "update",
	}

	// Create policy
	err := testModule.auth.serv.Create(policy1)
	if err != nil {
		t.Error(err)
	}

	// Build slug
	requestUrl := "auth"

	req, err := helpers.BuildApiRequest("PUT", requestUrl, helpers.BuildReqBody(models.UpdateCasbinRule{OldPolicy: policy1, NewPolicy: policy2}), true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v.\nResp:%s", "Auth Update",
			status, http.StatusOK, rr.Body.String())
	}

	// Check if the record is updated
	found, err := testModule.auth.serv.FindByResource(policy2.Resource)

	if err != nil {
		t.Errorf("Error detected when finding resource: %v", err)
	}
	if len(found) == 0 {
		t.Errorf("Expected to find resource, however, not found: %v", found)
	}

	// Check details
	checkPolicyDetails(t, found[0], policy2)

	// Delete policy
	err = testModule.auth.serv.Delete(policy2)
	if err != nil {
		t.Error(err)
	}
}

// // Role
func TestAuthController_FindAllRoles(t *testing.T) {
	numberOfDetaultRoles := 3
	req, err := helpers.BuildApiRequest("GET", "auth/roles", nil, true, testModule.accounts.admin.token)
	if err != nil {
		t.Error(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%v: got %v want %v", "Auth Find all roles",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []string
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Checks if the number of roles is correct
	if len(body) != numberOfDetaultRoles {
		t.Errorf("Expected %v, got %v: %v", numberOfDetaultRoles, len(body), body)
	}
}
func TestAuthController_AssignUserRole(t *testing.T) {
	requestUrl := "auth/roles"

	var tests = []struct {
		name             string
		role             string
		expectedResponse int
	}{
		{"Assign user role", "admin", http.StatusOK},
		{"Fail: Assign inexistent role", "jester", http.StatusBadRequest},
	}

	for _, v := range tests {
		// Build API request
		req, err := helpers.BuildApiRequest("PUT", requestUrl, helpers.BuildReqBody(models.CasbinRoleAssignment{
			UserId: fmt.Sprint(testModule.accounts.user.details.ID),
			Role:   v.role}), true, testModule.accounts.admin.token)
		if err != nil {
			t.Error(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponse {
			t.Errorf("%v: got %v want %v.\nResp:%s", v.name,
				status, v.expectedResponse, rr.Body.String())
		}

		// If a successful response is expected
		if v.expectedResponse == http.StatusOK {
			// Check if the user role was reassigned
			found, err := testModule.auth.serv.FindRoleByUserId(int(testModule.accounts.user.details.ID))
			if err != nil {
				t.Error(err)
			}

			if found != v.role {
				t.Errorf("Expected %v, got %v", v.role, found)
			}

			// Return user to default role
			success, err := testModule.auth.serv.AssignUserRole(fmt.Sprint(testModule.accounts.user.details.ID), "user")
			if err != nil {
				t.Error(err)
			}
			// Convert to bool
			successValue := *success
			if !successValue {
				t.Errorf("Expected to reset role reassignment, however, failed")
			}
		}
	}
}
func TestAuthController_CreateRole(t *testing.T) {
	requestUrl := "auth/roles"

	var tests = []struct {
		name             string
		role             string
		expectedResponse int
	}{
		{"Create new role", "jester", http.StatusCreated},
		{"Fail: Creating existent role", "admin", http.StatusBadRequest},
	}

	for _, v := range tests {
		// Build API request
		req, err := helpers.BuildApiRequest("POST", requestUrl, helpers.BuildReqBody(models.CasbinRoleAssignment{
			UserId: fmt.Sprint(testModule.accounts.user.details.ID),
			Role:   v.role}), true, testModule.accounts.admin.token)
		if err != nil {
			t.Error(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponse {
			t.Errorf("%v: got %v want %v.\nResp:%s", v.name,
				status, v.expectedResponse, rr.Body.String())
		}

		// If a successful response is expected
		if v.expectedResponse == http.StatusCreated {
			// Check if the user role was reassigned
			found, err := testModule.auth.serv.FindRoleByUserId(int(testModule.accounts.user.details.ID))
			if err != nil {
				t.Error(err)
			}

			if found != v.role {
				t.Errorf("Expected %v, got %v", v.role, found)
			}

			// Return user to default role
			success, err := testModule.auth.serv.AssignUserRole(fmt.Sprint(testModule.accounts.user.details.ID), "user")
			if err != nil {
				t.Fatal(err)
			}
			// Convert to bool
			successValue := *success
			if !successValue {
				t.Fatal("Expected to reset role reassignment, however, failed")
			}

			// Check user role to ensure completed correctly
			foundRole, err := testModule.auth.serv.FindRoleByUserId(int(testModule.accounts.user.details.ID))
			if err != nil {
				t.Error(err)
			}
			if foundRole != "user" {
				t.Fatalf("Expected %v, got %v", "user", foundRole)
			}
		}
	}
}

// Role Inheritance
func TestAuthController_FindAllRoleInheritance(t *testing.T) {
	numberOfDetaultInheritances := 2
	var tests = []struct {
		name             string
		tokenToUse       string
		expectedResponse int
	}{
		{"Admin: Find all role inheritance", testModule.accounts.admin.token, http.StatusOK},
		{"User: Fail: Find all role inheritance", testModule.accounts.user.token, http.StatusForbidden},
	}

	for _, v := range tests {
		// Build request
		req, err := helpers.BuildApiRequest("GET", "auth/inheritance", nil, true, v.tokenToUse)
		if err != nil {
			t.Error(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponse {
			t.Errorf("%v: got %v want %v", v.name,
				status, v.expectedResponse)
		}

		// If a successful response is expected
		if v.expectedResponse == http.StatusOK {
			// Convert response JSON to struct
			var body []models.GRecord
			json.Unmarshal(rr.Body.Bytes(), &body)

			// Checks if the number of roles is correct
			if len(body) != numberOfDetaultInheritances {
				t.Errorf("Expected %v, got %v", numberOfDetaultInheritances, len(body))
			}
			// Checks if the type of the records are correct
			if data.CheckSliceType(body, reflect.TypeOf(models.GRecord{})) == false {
				t.Errorf("Expected %v, got %v", reflect.TypeOf(models.GRecord{}), reflect.TypeOf(body[0]))
			}
		}

	}
}
func TestAuthController_CreateInheritance(t *testing.T) {
	requestUrl := "auth/inheritance"

	// Setup
	userToCreate := &models.CreateUser{Email: "krusty@gmail.com", Password: "password"}
	// Create user
	createdUser, err := testModule.users.serv.Create(userToCreate)
	if err != nil {
		t.Error(err)
	}
	success, err := testModule.auth.serv.CreateRole(fmt.Sprint(createdUser.ID), "jester")
	if err != nil {
		t.Error(err)
	}
	if *success != true {
		t.Error("Expected to create role, however, failed")
	}

	var tests = []struct {
		name             string
		policy           models.GRecord
		tokenToUse       string
		expectedResponse int
	}{
		{"Success: Admin Create role inheritance", models.GRecord{
			Role:         "jester",
			InheritsFrom: "admin",
		}, testModule.accounts.admin.token, http.StatusCreated},
		{"Fail: Admin Create existent role inheritance", models.GRecord{
			Role:         "admin",
			InheritsFrom: "moderator",
		}, testModule.accounts.admin.token, http.StatusBadRequest},
		{"Fail: User Create role inheritance", models.GRecord{
			Role:         "superadmin",
			InheritsFrom: "admin",
		}, testModule.accounts.user.token, http.StatusForbidden},
	}

	for _, v := range tests {
		// Build request
		req, err := helpers.BuildApiRequest("POST", requestUrl, helpers.BuildReqBody(v.policy), true, v.tokenToUse)
		if err != nil {
			t.Error(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponse {
			t.Errorf("%v: got %v want %v.\nResp:%s", v.name,
				status, v.expectedResponse, rr.Body.String())
		}

		if v.expectedResponse == http.StatusCreated {
			// Check if the role inheritance was created
			foundInheritances, err := testModule.auth.serv.FindAllRoleInheritance()
			if err != nil {
				t.Error(err)
			}

			// Iterate through found inheritances
			foundCreatedPolicy := false
			for _, inheritance := range foundInheritances {
				// See if match found
				if inheritance.Role == v.policy.Role && inheritance.InheritsFrom == v.policy.InheritsFrom {
					foundCreatedPolicy = true
				}
			}
			if !foundCreatedPolicy {
				t.Errorf("Expected to find created role inheritance. Only found: %v", foundInheritances)
			}

			// Delete role inheritance
			err = testModule.auth.serv.DeleteInheritance(v.policy)
			if err != nil {
				t.Error(err)
			}
		}
	}

	// Cleanup
	// Delete user
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Error(err)
	}
}
func TestAuthController_DeleteInheritance(t *testing.T) {
	requestUrl := "auth/inheritance"

	// Setup
	userToCreate := &models.CreateUser{Email: "edible@gmail.com", Password: "password"}
	// Create user
	createdUser, err := testModule.users.serv.Create(userToCreate)
	if err != nil {
		t.Error(err)
	}
	success, err := testModule.auth.serv.CreateRole(fmt.Sprint(createdUser.ID), "fester")
	if err != nil {
		t.Error(err)
	}
	if *success != true {
		t.Error("Expected to create role, however, failed")
	}

	var tests = []struct {
		name             string
		tokenToUse       string
		expectedResponse int
	}{
		{"Success: Admin Delete role inheritance", testModule.accounts.admin.token, http.StatusOK},
		{"Fail: User Delete inexistent role inheritance", testModule.accounts.user.token, http.StatusForbidden},
	}
	for _, v := range tests {
		// Create role inheritance to delete
		policy := models.GRecord{
			Role:         "fester",
			InheritsFrom: "admin",
		}
		err = testModule.auth.serv.CreateInheritance(policy)
		if err != nil {
			t.Error(err)
		}

		// Build request
		req, err := helpers.BuildApiRequest("DELETE", requestUrl, helpers.BuildReqBody(policy), true, v.tokenToUse)
		if err != nil {
			t.Error(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponse {
			t.Errorf("%v: got %v want %v.\nResp:%s", v.name,
				status, v.expectedResponse, rr.Body.String())
		}

		// If a successful response is expected
		if v.expectedResponse == http.StatusOK {
			// Check if the role inheritance was deleted
			foundInheritances, err := testModule.auth.serv.FindAllRoleInheritance()
			if err != nil {
				t.Error(err)
			}
			// Iterate through found inheritances
			foundDeletedPolicy := false
			for _, inheritance := range foundInheritances {
				// See if match found
				if inheritance.Role == policy.Role && inheritance.InheritsFrom == policy.InheritsFrom {
					foundDeletedPolicy = true
				}
			}
			// If found fail test
			if foundDeletedPolicy {
				t.Errorf("Expected to not find deleted role inheritance, however, found: %v", policy)
			}

			// Else if failure is detected, delete manually before next test
		} else if rr.Code == http.StatusBadRequest {
			// Delete role inheritance
			err = testModule.auth.serv.DeleteInheritance(policy)
			if err != nil {
				t.Error(err)
			}
		}

	}
	// Cleanup
	// Delete user
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Error(err)
	}
}

// Checks if the policy details are a match
func checkPolicyDetails(t *testing.T, body models.PolicyRuleCombinedActions, policy models.PolicyRule) {
	// Check details
	if body.Resource != policy.Resource {
		t.Errorf("Body: %+v. Expected %v, got %v", body, policy.Resource, body.Resource)
	}
	if body.Role != policy.Role {
		t.Errorf("Expected %v, got %v", policy.Role, body.Role)
	}
	// Check if array of strings contains the created record
	if utility.ArrayContainsString(body.Action, policy.Action) == false {
		t.Errorf("Expected %v, got %v", policy.Action, body.Action)
	}
}
