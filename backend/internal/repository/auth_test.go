package repository_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Policies
func TestAuthPolicyRepository_Create(t *testing.T) {
	policyToCreate := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	// Test function
	err := testModule.auth.repo.Create(policyToCreate)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Tear down
	err = testModule.auth.repo.Delete(policyToCreate)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_Delete(t *testing.T) {
	policyToCreate := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	err := testModule.auth.repo.Create(policyToCreate)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	err = testModule.auth.repo.Delete(policyToCreate)
	if err != nil {
		t.Fatalf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_FindAll(t *testing.T) {
	policy1 := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "create"}
	policy2 := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "update"}
	// Setup
	err := testModule.auth.repo.Create(policy1)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	err = testModule.auth.repo.Create(policy2)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	policies, err := testModule.auth.repo.FindAll()
	if err != nil {
		t.Errorf("Error finding policies: %v", err)
	}
	if len(policies) != 2 {
		t.Errorf("Expected 2 policies, found %v", len(policies))
	}

	// Iterate through policies
	for _, policy := range policies {
		// Check if policy resource matches policy1 or policy2
		if policy[1] == policy1.V0 {
			// Iterate through policy checking details against policy1
			checkArrayStringPolicyAgainstCasbinRule(policy, policy1, t)
		} else if policy[1] == policy2.V0 {
			checkArrayStringPolicyAgainstCasbinRule(policy, policy2, t)
		}

	}
	// Cleanup
	err = testModule.auth.repo.Delete(policy1)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
	err = testModule.auth.repo.Delete(policy2)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

func TestAuthPolicyRepository_Update(t *testing.T) {
	oldPolicy := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "GET"}
	newPolicy := models.CasbinRule{V0: "admin", V1: "/api/v1/users", V2: "POST"}
	// Setup
	err := testModule.auth.repo.Create(oldPolicy)
	if err != nil {
		t.Errorf("Error creating policy: %v", err)
	}
	// Test function
	err = testModule.auth.repo.Update(oldPolicy, newPolicy)
	if err != nil {
		t.Errorf("Error updating policy: %v", err)
	}

	// Cleanup
	err = testModule.auth.repo.Delete(newPolicy)
	if err != nil {
		t.Errorf("Error deleting policy: %v", err)
	}
}

// Roles
func TestAuthPolicyRepository_AssignUserRole(t *testing.T) {
	createdUser1, err := testModule.users.repo.Create(&db.User{Email: "ratbag@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Test function
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser1.ID), "admin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}
func TestAuthPolicyRepository_FindAllRoles(t *testing.T) {
	// Test function
	roles, err := testModule.auth.repo.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}

	// Check that default roles are available
	if len(roles) != 3 {
		t.Errorf("Expected 2 roles, found %v: %v", len(roles), roles)
	}
}
func TestAuthPolicyRepository_FindRoleByUserId(t *testing.T) {
	roleToCreate := "pikachu"
	// Create user
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "pikachi@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Setup
	success, err := testModule.auth.repo.CreateRole(fmt.Sprint(createdUser.ID), roleToCreate)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// Test function
	role, err := testModule.auth.repo.FindRoleByUserId(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error finding role: %v", err)
	}
	if role != roleToCreate {
		t.Errorf("Expected %v, found %v", roleToCreate, role)
	}

	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}
func TestAuthPolicyRepository_CreateRole(t *testing.T) {
	userToCreate := &db.User{Email: "marsattackso@gmail.com", Password: "password"}
	createdUser1, role := createUserAndSetRole(*userToCreate, "admin", t)

	// Test function
	roleToCreate := "alien"
	success, err := testModule.auth.repo.CreateRole(fmt.Sprint(createdUser1.ID), roleToCreate)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if success == nil {
		t.Errorf("Expected true, found %v", success)
	}

	// Check that role has been created
	roles, err := testModule.auth.repo.FindAllRoles()
	if err != nil {
		t.Errorf("Error finding roles: %v", err)
	}

	rolesContains := utility.ArrayContainsString(roles, roleToCreate)

	if !rolesContains {
		t.Errorf("Did not find %v contained in %v", role, roles)
	}

	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser1.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}
func TestAuthRoleRepository_DeleteRolesForUser(t *testing.T) {
	// Create user
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "pikachu@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	// Setup
	success, err := testModule.auth.repo.AssignUserRole(fmt.Sprint(createdUser.ID), "admin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// Test function
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// Check that role has been deleted
	gRecords, err := app.Auth.Enforcer.GetGroupingPolicy()
	if err != nil {
		t.Errorf("Error getting roles from enforcer: %v", err)
	}
	if len(gRecords) != 2 {
		t.Errorf("Expected to have only default 2, found %v", len(gRecords))
	}

	// Cleanup
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

// // Role inheritances
func TestAuthPolicyRepository_CreateInheritance(t *testing.T) {
	// Setup
	createdUser, err := testModule.users.repo.Create(&db.User{Email: "piapaika@gmail.com", Password: "password"})
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	success, err := testModule.auth.repo.CreateRole(fmt.Sprint(createdUser.ID), "superadmin")
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}

	// inheritance preparation
	inheritanceToCreate := models.GRecord{
		Role: "superadmin", InheritsFrom: "admin",
	}

	// Test function
	err = testModule.auth.repo.CreateInheritance(inheritanceToCreate)
	if err != nil {
		t.Errorf("Error adding role inheritance: %v", err)
	}

	// Check that role has been created
	rolesAndAssignments, err := app.Auth.Enforcer.GetNamedGroupingPolicy("g")
	if err != nil {
		t.Errorf("Error getting roles from enforcer: %v", err)
	}
	// Filter out only inheritance records
	inheritances, err := helpers.FilterOnlyInheritanceToGRecords(rolesAndAssignments)
	if err != nil {
		t.Errorf("Error filtering roles from enforcer: %v", err)
	}
	if len(inheritances) != 3 {
		t.Errorf("Expected 3, found %v", len(rolesAndAssignments))
	}

	// Look for created inheritance
	for _, inheritance := range inheritances {
		// Check values
		if inheritance.Role == "role:superadmin" {
			if inheritance.InheritsFrom != "role:admin" {
				t.Errorf("Expected admin, found %v", inheritance.InheritsFrom)
			}
		}
	}

	// Cleanup
	// Apply naming convention to new role record
	adminpanel.ApplyNamingConventionToRoleInheritanceRecord(&inheritanceToCreate)
	// Remove policy from enforcer
	removed, err := app.Auth.Enforcer.RemoveGroupingPolicy(inheritanceToCreate.Role, inheritanceToCreate.InheritsFrom)
	if err != nil {
		t.Errorf("Error removing role inheritance: %v", err)
	}
	if !removed {
		t.Errorf("Expected true, found %v", removed)
	}
	// Remove Delete user and role
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
	success, err = testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(createdUser.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
	}
}

func TestAuthPolicyRepository_DeleteInheritance(t *testing.T) {
	// Create non-convention role names
	inheritanceToCreate := &models.GRecord{Role: "superadmin", InheritsFrom: "admin"}
	// The created user roles contain the naming convention
	createdUser1, createdUser1Role := createUserAndSetRole(db.User{Email: "wysiwyg@gmail.com", Password: "password"}, inheritanceToCreate.Role, t)
	if createdUser1 == nil {
		t.Errorf("Error creating users")
	}
	fmt.Printf("Created user 1: %+v role: %v\n", createdUser1, createdUser1Role)

	// Apply naming convention to new role record
	conventionalInheritanceToCreate := *inheritanceToCreate
	// Apply naming convention to new role record
	adminpanel.ApplyNamingConventionToRoleInheritanceRecord(&conventionalInheritanceToCreate)
	// Add role inheritance
	success, err := app.Auth.Enforcer.AddGroupingPolicy(conventionalInheritanceToCreate.Role, conventionalInheritanceToCreate.InheritsFrom)
	if err != nil {
		t.Fatalf("Error adding role inheritance: %v", err)
	}
	if !success {
		t.Fatalf("Expected true, found %v", success)
	}

	// Test function
	err = testModule.auth.repo.DeleteInheritance(models.GRecord{Role: inheritanceToCreate.Role, InheritsFrom: inheritanceToCreate.InheritsFrom})
	if err != nil {
		t.Errorf("Error deleting role inheritance: %v", err)
	}

	// Check that role has been deleted (Should only be 2 (users with roles) left)
	inheritance, err := app.Auth.Enforcer.GetNamedGroupingPolicy("g")
	if err != nil {
		t.Errorf("Error getting roles from enforcer: %v", err)
	}

	if len(inheritance) != 3 {
		t.Errorf("Expected 3, found %v.\nInheritance: %v", len(inheritance), inheritance)
	}

	// Cleanup
	err = deleteUserAndRole(createdUser1, t)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}
}

func TestAuthPolicyRepository_FindAllInheritances(t *testing.T) {
	// Test function
	inheritances, err := testModule.auth.repo.FindAllRoleInheritance()
	if err != nil {
		t.Errorf("Error finding role inheritances: %v", err)
	}
	if len(inheritances) != 2 {
		t.Errorf("Expected 2 inheritances, found %v", len(inheritances))
	}

	expectedInheritance1 := models.GRecord{Role: "admin", InheritsFrom: "moderator"}
	expectedInheritance2 := models.GRecord{Role: "moderator", InheritsFrom: "user"}

	// Check details of each inheritance
	for _, inheritance := range inheritances {
		// If match found to inheritance 1
		if inheritance.Role == expectedInheritance1.Role {
			helpers.CompareObjects(expectedInheritance1, inheritance, t, []string{"Role", "InheritsFrom"})
			// Else if match found to inheritance 2
		} else if inheritance.Role == expectedInheritance2.Role {
			helpers.CompareObjects(expectedInheritance2, inheritance, t, []string{"Role", "InheritsFrom"})
		}
	}
}

// Creates a user and assigns a role. Returns the user and the role if successful
func createUserAndSetRole(user db.User, role string, t *testing.T) (*db.User, string) {
	createdUser, err := testModule.users.repo.Create(&user)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
		return nil, ""
	}

	// Add role naming convention
	role = "role:" + role
	success, err := app.Auth.Enforcer.AddRoleForUser(fmt.Sprint(createdUser.ID), role)
	if err != nil {
		t.Errorf("Error assigning role to user: %v", err)
		return nil, ""
	}
	if !success {
		t.Errorf("Expected true, found %v", success)
		return nil, ""
	}
	return createdUser, role
}

// Deletes a user and the role manually
func deleteUserAndRole(user *db.User, t *testing.T) error {
	err := testModule.users.repo.Delete(int(user.ID))
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
		return err
	}
	success, err := testModule.auth.repo.DeleteRolesForUser(fmt.Sprint(user.ID))
	if err != nil {
		t.Errorf("Error deleting roles for user: %v", err)
		return err
	}
	if !*success {
		t.Errorf("Expected true, found %v", *success)
		// Return new error
		return errors.New("Error deleting roles for user")
	}
	return nil
}

// Helper functions
func checkArrayStringPolicyAgainstCasbinRule(policy []string, casbinRule models.CasbinRule, t *testing.T) bool {
	if policy[0] != casbinRule.V0 {
		t.Errorf("Expected %v, found %v", casbinRule.V0, policy[0])
	}
	if policy[1] != casbinRule.V1 {
		t.Errorf("Expected %v, found %v", casbinRule.V1, policy[1])
		return false
	}
	if policy[2] != casbinRule.V2 {
		t.Errorf("Expected %v, found %v", casbinRule.V2, policy[2])
		return false
	}
	return true
}
