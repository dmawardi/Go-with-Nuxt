package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestUserController_Find(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "greenie@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create test user for test: %v", err)
	}

	// Create a request with an "id" URL parameter
	req, err := helpers.BuildApiRequest("GET", fmt.Sprintf("users/%v", createdUser.ID), nil, true, testModule.accounts.admin.token)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body models.UserWithRole
	json.Unmarshal(rr.Body.Bytes(), &body)

	// check user details for match
	helpers.CompareObjects(body, createdUser, t, []string{"ID", "Username", "Email", "Name"})

	// Delete the created user
	delResult := testModule.users.serv.Delete(int(createdUser.ID))
	if delResult != nil {
		t.Fatalf("Error clearing created user")
	}
}

func TestUserController_FindAll(t *testing.T) {
	// Test finding two already created users for authentication mocking
	// Create a new request
	req, err := helpers.BuildApiRequest("GET", "users?limit=10&offset=0&order=id_desc", nil, true, testModule.accounts.admin.token)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testModule.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body *models.BasicPaginatedResponse[db.User]
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of user array
	if len(*body.Data) != 2 {
		t.Errorf("Users array in findAll failed: expected %d, got %d", 2, len(*body.Data))
	}

	// Iterate through users array received
	for _, item := range *body.Data {
		// If id is admin id
		if item.ID == testModule.accounts.admin.details.ID {
			helpers.CompareObjects(item, testModule.accounts.admin.details, t, []string{"ID", "Username", "Email", "Name"})

		} else {
			helpers.CompareObjects(item, testModule.accounts.user.details, t, []string{"ID", "Username", "Email", "Name"})

		}
	}

	// Test parameter input
	var failParameterTests = []struct {
		test_name              string
		limit                  string
		offset                 string
		order                  string
		expectedResponseStatus int
	}{
		// Only limit
		{test_name: "Only limit", limit: "5", offset: "", order: "", expectedResponseStatus: http.StatusOK},
		// Check normal parameters functional with order by
		{test_name: "Normal test", limit: "20", offset: "1", order: "id", expectedResponseStatus: http.StatusOK},
		// Descending order
		{test_name: "Normal test", limit: "20", offset: "1", order: "id_desc", expectedResponseStatus: http.StatusOK},
	}
	for _, v := range failParameterTests {
		req, err := helpers.BuildApiRequest("GET", fmt.Sprintf("users?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order), nil, true, testModule.accounts.admin.token)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("In test '%s': handler returned wrong status code: got %v want %v", v.test_name,
				status, v.expectedResponseStatus)
		}
	}
}

func TestUserController_Delete(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "zubayle@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create user for test: %v", err)
	}

	// Test parameter input
	var tests = []struct {
		testName               string
		tokenToUse             string
		expectedResponseStatus int
	}{
		{testName: "Normal user delete failure", tokenToUse: testModule.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Put last to also replace test user deletion
		{testName: "Admin user delete success", tokenToUse: testModule.accounts.admin.token, expectedResponseStatus: http.StatusOK},
	}

	for _, v := range tests {
		// Create a request
		req, err := helpers.BuildApiRequest("DELETE", fmt.Sprintf("users/%v", createdUser.ID), nil, true, v.tokenToUse)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testModule.router.ServeHTTP(rr, req)

		// Check response is failed for normal user
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Deletion test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

	}
}

func TestUserController_Update(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.serv.Create(&models.CreateUser{
		Username: "Jabar",
		Email:    "greenthumb@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})
	if err != nil {
		t.Fatalf("failed to create test user for test: %v", err)
	}

	var updateTests = []struct {
		testName string
		// To be converted to string for URL
		urlExtension           interface{}
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		{"Fail case: User updating another user", createdUser.ID, map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testModule.accounts.user.token, http.StatusForbidden, false},
		{"Admin updating a user", createdUser.ID, map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
		}, testModule.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{"Fail case: Update with validation errors", createdUser.ID, map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testModule.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating
		{"Fail case: User invalid update should fail due to permissions", createdUser.ID, map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testModule.accounts.user.token, http.StatusForbidden, false},
		// Should fail as url extension is incorrect
		{"Fail case: Bad url parameter", "x", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testModule.accounts.admin.token, http.StatusBadRequest, false},
		// Should fail as id is above currently created
		{"Fail case: Bad url parameter", "99", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, testModule.accounts.admin.token, http.StatusBadRequest, false},
	}

	for _, v := range updateTests {
		req, err := helpers.BuildApiRequest("PUT", fmt.Sprintf("users/%v", createdUser.ID), helpers.BuildReqBody(v.data), true, v.tokenToUse)
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			helpers.UpdateModelFields(createdUser, v.data)

			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			// check user details for match
			helpers.CompareObjects(body, createdUser, t, []string{"ID", "Username", "Email", "Name"})
		}
	}

	// Delete the created user
	err = testModule.users.serv.Delete(int(createdUser.ID))
	if err != nil {
		t.Fatalf("Error clearing created user")
	}
}

func TestUserController_Create(t *testing.T) {
	var tests = []struct {
		testName               string
		data                   models.CreateUser
		expectedResponseStatus int
	}{
		{"Successful user creation", models.CreateUser{
			Username: "Jabarnam",
			Email:    "gabor@ymail.com",
			Password: "password",
			Name:     "Bambaliya",
		}, http.StatusCreated},
		{"Successful user creation", models.CreateUser{
			Username: "Swalanim",
			Email:    "salvia@ymail.com",
			Password: "seradfasdf",
			Name:     "CreditTomyaA",
		}, http.StatusCreated},
		{"Failure: Not email", models.CreateUser{
			Username: "Yukon",
			Email:    "Sylvio",
			Password: "wowogsdfg",
			Name:     "Sosawsdfgsdfg",
		}, http.StatusBadRequest},
		{"Failure: Pass/Name field length", models.CreateUser{
			Username: "Jabarnam",
			Email:    "Cakawu@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
		{"Failure: Duplicate user", models.CreateUser{
			Username: "Jabarnam",
			Email:    "Jabal@ymail.com",
			Password: "as",
			Name:     "df",
		}, http.StatusBadRequest},
	}

	for _, v := range tests {
		req, err := helpers.BuildApiRequest("POST", "users", helpers.BuildReqBody(v.data), false, "")

		// Make new request with user update in body
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {

			t.Errorf("%s: Got %v want %v.", v.data.Name,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body models.UserWithRole
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Delete the created user
		err = testModule.users.serv.Delete(int(body.ID))
		if err != nil {
			t.Fatalf("Error clearing created user")
		}
	}
}

// My Profile
func TestUserController_GetMyUserDetails(t *testing.T) {
	// Create a request url
	requestUrl := "me"
	var tests = []struct {
		testName               string
		checkDetails           bool
		useToken               bool
		tokenToUse             string
		userToCheck            models.UserWithRole
		expectedResponseStatus int
	}{
		{"User checking own profile", true, true, testModule.accounts.user.token, *testModule.accounts.user.details, http.StatusOK},
		{"Admin checking own profile", true, true, testModule.accounts.admin.token, *testModule.accounts.admin.details, http.StatusOK},
		// Deny access to user that doesn't have authentication
		{"Logged out user checking profile", false, false, "", models.UserWithRole{}, http.StatusForbidden},
	}

	for _, v := range tests {
		req, err := helpers.BuildApiRequest("GET", requestUrl, nil, v.useToken, v.tokenToUse)
		// Make new request with user update in body
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)
		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%s: Got %v want %v.(%v)", v.testName,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If need to check details
		if v.checkDetails == true {
			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			// Check user details using updated object
			helpers.CompareObjects(body, &v.userToCheck, t, []string{"ID", "Username", "Email", "Name"})
		}
	}
}

func TestUserController_UpdateMyProfile(t *testing.T) {
	// Create a request url with an "id" URL parameter
	requestUrl := "me"

	var tests = []struct {
		testName               string
		data                   map[string]string
		useToken               bool
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		loggedInDetails        models.UserWithRole
	}{
		{"Admin self update", map[string]string{
			"Username": "JabarCindi",
			"Name":     "Bambaloonie",
		}, true, testModule.accounts.admin.token, http.StatusOK, true, *testModule.accounts.admin.details},
		{"User self update", map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, true, testModule.accounts.user.token, http.StatusOK, true, *testModule.accounts.user.details},
		{"User self update with invalid email", map[string]string{
			"Email": "JabarHindi",
		}, true, testModule.accounts.user.token, http.StatusBadRequest, false, *testModule.accounts.user.details},
		{"Fail: User self update with duplicate email", map[string]string{
			"Username": "Swahili",
			"Email":    testModule.accounts.admin.details.Email,
		}, true, testModule.accounts.user.token, http.StatusBadRequest, false, *testModule.accounts.user.details},
		{"Fail: User update without token", map[string]string{
			"Username": "JabarHindi",
			"Name":     "Bambaloonie",
			"Password": "YeezusChris",
		}, false, "", http.StatusForbidden, false, *testModule.accounts.user.details},
		{"Fail: Admin update with invalid validation", map[string]string{
			"Username": "Gobod",
			"Name":     "solu",
		}, true, testModule.accounts.admin.token, http.StatusBadRequest, false, *testModule.accounts.admin.details},
		{"Fail: User update with invalid validation", map[string]string{
			"Username": "Gabor",
			"Name":     "solu",
		}, true, testModule.accounts.user.token, http.StatusBadRequest, false, *testModule.accounts.user.details},
	}

	for _, v := range tests {
		// Make new request with user update in body
		req, err := helpers.BuildApiRequest("PUT", requestUrl, helpers.BuildReqBody(v.data), v.useToken, v.tokenToUse)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()
		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)

		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%v: got %v want %v. Msg: %s\n", v.testName,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// If need to check details
		if v.checkDetails == true {
			// Update created user struct with the changes pushed through API
			helpers.UpdateModelFields(&v.loggedInDetails, v.data)

			// Convert response JSON to struct
			var body models.UserWithRole
			json.Unmarshal(rr.Body.Bytes(), &body)

			helpers.CompareObjects(body, v.loggedInDetails, t, []string{"ID", "Username", "Email", "Name"})
		}

		// Return updates to original state
		testModule.users.serv.Update(int(testModule.accounts.admin.details.ID), &models.UpdateUser{
			Username: testModule.accounts.admin.details.Username,
			Password: testModule.accounts.admin.details.Password,
			Email:    testModule.accounts.admin.details.Email,
			Name:     testModule.accounts.admin.details.Name,
		})
		testModule.users.serv.Update(int(testModule.accounts.user.details.ID), &models.UpdateUser{
			Username: testModule.accounts.user.details.Username,
			Password: testModule.accounts.user.details.Password,
			Email:    testModule.accounts.user.details.Email,
			Name:     testModule.accounts.user.details.Name,
		})
	}
}

// Login
func TestUserController_Login(t *testing.T) {
	// Create a request url with an "id" URL parameter
	requestUrl := "users/login"

	var tests = []struct {
		testName               string
		data                   models.Login
		expectedResponseStatus int
		failureExpected        bool
		expectedMessage        string
	}{
		{"Admin user login", models.Login{
			Email:    testModule.accounts.admin.details.Email,
			Password: testModule.accounts.admin.details.Password,
		}, http.StatusOK, false, ""},
		{"Fail: Admin user incorrect details", models.Login{
			Email:    testModule.accounts.admin.details.Email,
			Password: "wrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Basic user login", models.Login{
			Email:    testModule.accounts.user.details.Email,
			Password: testModule.accounts.user.details.Password,
		}, http.StatusOK, false, ""},
		{"Fail: Basic user incorrect details", models.Login{
			Email:    testModule.accounts.user.details.Email,
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Fail: Non existent user login", models.Login{
			Email:    "jester@gmail.com",
			Password: "VeryWrongPassword",
		}, http.StatusUnauthorized, true, "Invalid Credentials\n"},
		{"Fail: Invalid email user login", models.Login{
			Email:    "jester",
			Password: "VeryWrongPassword",
		}, http.StatusBadRequest, false, ""},
		{"Fail: Empty credentials", models.Login{
			Email:    "jester",
			Password: "",
		}, http.StatusBadRequest, false, ""},
	}

	for _, v := range tests {
		req, err := helpers.BuildApiRequest("POST", requestUrl, helpers.BuildReqBody(v.data), false, "")
		// Make request with update in body
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)

		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%v: Got %v want %v. \nResp: %v", v.testName,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If failure is expected
		if v.failureExpected {
			// Form req body
			reqBody := rr.Body.String()
			// Check if matches with expectation
			if reqBody != v.expectedMessage {
				t.Errorf("%v: The body is: %v. expected: %v.", v.testName, rr.Body.String(), v.expectedMessage)
			}

		}

	}
}

func TestUserController_ResetPassword(t *testing.T) {
	// Create a request url with an "id" URL parameter
	requestUrl := "users/forgot-password"

	var tests = []struct {
		testName               string
		data                   models.ResetPasswordAndEmailVerification
		expectedResponseStatus int
		checkMessage           bool
		expectedMessage        string
	}{
		{"Successful reset password", models.ResetPasswordAndEmailVerification{
			Email: testModule.accounts.admin.details.Email,
		}, http.StatusOK, false, ""},
		{"Fail: Non existent user reset password", models.ResetPasswordAndEmailVerification{
			Email: "baffoon@snailmail.com"}, http.StatusBadRequest, true, "Password reset request failed\n"},
		// The below tests should return a validation errors object
		{"Fail: Invalid email user reset password", models.ResetPasswordAndEmailVerification{
			Email: "baffoon"}, http.StatusBadRequest, false, ""},
		{"Fail: Empty email user reset password", models.ResetPasswordAndEmailVerification{
			Email: ""}, http.StatusBadRequest, false, ""},
	}

	for _, v := range tests {

		req, err := helpers.BuildApiRequest("POST", requestUrl, helpers.BuildReqBody(v.data), false, "")
		// Make request with update in body
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)

		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%v: Got %v want %v. \nResp: %v", v.testName,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If failure is expected
		if v.checkMessage {
			// Form req body
			reqBody := rr.Body.String()
			// Check if matches with expectation
			if reqBody != v.expectedMessage {
				t.Errorf("%v: The body is: %v. expected: %v.", v.testName, rr.Body.String(), v.expectedMessage)
			}

		}

	}
}

func TestUserController_ResendVerificationEmail(t *testing.T) {
	// Create a request url with an "id" URL parameter
	requestUrl := "users/send-verification-email"

	var tests = []struct {
		testName               string
		useToken               bool
		tokenToUse             string
		expectedResponseStatus int
		checkMessage           bool
		expectedMessage        string
	}{
		{"Successful (user): resend verification email", true, testModule.accounts.user.token, http.StatusOK, false, ""},
		{"Successful (admin): resend verification email", true, testModule.accounts.admin.token, http.StatusOK, false, ""},
		{"Fail: Not logged in", false, "", http.StatusForbidden, true, "Error parsing authentication token\n"},
	}

	for _, v := range tests {
		req := &http.Request{}
		err := error(nil)

		// If using token
		if v.useToken {
			req, err = helpers.BuildApiRequest("POST", requestUrl, nil, v.useToken, v.tokenToUse)
			// Make request with update in body
			if err != nil {
				t.Fatal(err)
			}
			// Else build without token
		} else {
			req, err = helpers.BuildApiRequest("POST", requestUrl, nil, false, "")
			if err != nil {
				t.Fatal(err)
			}

		}

		// Make request with update in body
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Send update request to mock server
		testModule.router.ServeHTTP(rr, req)

		// Check response is failed for normal user to update another
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("%v: Got %v want %v. \nResp: %v", v.testName,
				status, v.expectedResponseStatus, rr.Body)
		}

		// If failure is expected
		if v.checkMessage {
			// Form req body
			reqBody := rr.Body.String()
			// Check if matches with expectation
			if reqBody != v.expectedMessage {
				t.Errorf("%v: The body is: %v. expected: %v.", v.testName, rr.Body.String(), v.expectedMessage)
			}

		}

	}
}
