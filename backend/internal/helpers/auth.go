package helpers

import (
	"strings"
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers/utility"
	"github.com/dmawardi/Go-Template/internal/models"
)

// Generates a verification code and sets expiry and returns a user object with the values
func GenerateVerificationCodeAndSetExpiry() (*db.User, error) {
	userUpdate := &db.User{}
	// Generate token for verification
	tokenCode, err := utility.GenerateRandomString(25)
	if err != nil {
		return nil, err
	}
	// Update user to be unverified
	verified := false
	userUpdate.Verified = &verified
	// Set token code
	userUpdate.VerificationCode = tokenCode
	// Set verification code expiry to 12 hours from now
	userUpdate.VerificationCodeExpiry = time.Now().Add(12 * time.Hour)
	return userUpdate, nil
}

// SearchG2Records searches through all fields of each G2Record for the searchTerm and adds any match found in the results
func SearchGRecords(records []models.GRecord, searchTerm string) []models.GRecord {
	var result []models.GRecord

	// Iterate through the slice of G2Records
	for _, record := range records {
		// Check if the searchTerm is in any of the record's fields
		if strings.Contains(record.Role, searchTerm) || strings.Contains(record.InheritsFrom, searchTerm) {
			result = append(result, record)
		}
	}

	return result
}

// FilterOnlyRolesToGRecords filters out the role inheritances from a slice of roles and assignments (All G records)
func FilterOnlyInheritanceToGRecords(rolesAndAssignments [][]string) ([]models.GRecord, error) {
	var roleInheritancePolicies []models.GRecord

	for _, policy := range rolesAndAssignments {
		// Assuming policy[0] contains the role/subject and policy[1] contains the inherited role
		// Adjust the indexing based on your actual policy structure
		if strings.HasPrefix(policy[0], "role:") && strings.HasPrefix(policy[1], "role:") {
			roleInheritancePolicies = append(roleInheritancePolicies, models.GRecord{Role: policy[0], InheritsFrom: policy[1]})
		}
	}

	return roleInheritancePolicies, nil
}

func FilterOnlyRolesToList(rolesAndAssignments [][]string) []string {
	var roles []string
	// Filter out the roles that are not user assigned
	for _, policy := range rolesAndAssignments {
		if strings.HasPrefix(policy[0], "role:") {
			// If not already contained within the slice, add it
			if !utility.ArrayContainsString(roles, policy[0]) {
				roles = append(roles, policy[0])
			}
		}
		// If inherits from is a role, add it to the slice
		if strings.HasPrefix(policy[1], "role:") {
			// If not already contained within the slice, add it
			if !utility.ArrayContainsString(roles, policy[1]) {
				roles = append(roles, policy[1])
			}
		}
	}
	return roles
}

// Grabs a slice of GRecords and filters out the roles into a string slice
func ConvertInheritanceGRecordsToRoleList(roles []models.GRecord) []string {
	var roleList []string
	// Iterate through inheritance policies and add to roles slice
	for _, role := range roles {
		if !utility.ArrayContainsString(roleList, role.Role) {
			roleList = append(roleList, role.Role)
		}
		if !utility.ArrayContainsString(roleList, role.InheritsFrom) {
			roleList = append(roleList, role.InheritsFrom)
		}
	}

	return roleList
}
