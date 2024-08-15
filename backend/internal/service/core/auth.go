package coreservices

import (
	"fmt"
	"sort"
	"strings"

	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/models"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
)

type AuthPolicyService interface {
	// Policies
	FindAll(searchQuery string) ([]models.PolicyRuleCombinedActions, error)
	FindByResource(policyResource string) ([]models.PolicyRuleCombinedActions, error)
	Create(policy models.PolicyRule) error
	Update(oldPolicy, newPolicy models.PolicyRule) error
	Delete(policy models.PolicyRule) error
	// Roles
	FindAllRoles() ([]string, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	CreateRole(userId, roleToApply string) (*bool, error)
	// Inheritance
	FindAllRoleInheritance() ([]models.GRecord, error)
	CreateInheritance(inherit models.GRecord) error
	DeleteInheritance(inherit models.GRecord) error
	// Not for controller usage (used in auth)
	FindRoleByUserId(userId int) (string, error)
}

type authPolicyService struct {
	repo corerepositories.AuthPolicyRepository
}

func NewAuthPolicyService(repo corerepositories.AuthPolicyRepository) AuthPolicyService {
	return &authPolicyService{repo}
}

// Policies
//

func (s *authPolicyService) FindAll(searchQuery string) ([]models.PolicyRuleCombinedActions, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Transform data for easier consumption
	groupsSlice := transformDataToResponse(data)

	if searchQuery != "" {
		// Filter policies []map[string]interface{} by search query
		groupsSlice = adminpanel.SearchPoliciesByResource(groupsSlice, searchQuery)
	}

	// Sort by resource alphabetically
	sort.Slice(groupsSlice, func(i, j int) bool {
		return groupsSlice[i].Resource < groupsSlice[j].Resource
	})

	return groupsSlice, nil
}
func (s *authPolicyService) FindByResource(policyResource string) ([]models.PolicyRuleCombinedActions, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Transform data for easier consumption
	organizedData := transformDataToResponse(data)
	// Filter by search query
	resourceMatchRecords := adminpanel.SearchPoliciesForExactResouceMatch(organizedData, policyResource)

	return resourceMatchRecords, nil
}
func (s *authPolicyService) Create(policy models.PolicyRule) error {
	casbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    policy.Role,
		V1:    policy.Resource,
		V2:    policy.Action,
	}

	return s.repo.Create(casbinPolicy)
}
func (s *authPolicyService) Update(oldPolicy, newPolicy models.PolicyRule) error {
	oldCasbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    oldPolicy.Role,
		V1:    oldPolicy.Resource,
		V2:    oldPolicy.Action,
	}
	newCasbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    newPolicy.Role,
		V1:    newPolicy.Resource,
		V2:    newPolicy.Action,
	}
	return s.repo.Update(oldCasbinPolicy, newCasbinPolicy)
}
func (s *authPolicyService) Delete(policy models.PolicyRule) error {
	casbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    policy.Role,
		V1:    policy.Resource,
		V2:    policy.Action,
	}
	return s.repo.Delete(casbinPolicy)
}

// Roles
//

func (s *authPolicyService) FindAllRoles() ([]string, error) {
	return s.repo.FindAllRoles()
}
func (s *authPolicyService) FindRoleByUserId(userId int) (string, error) {
	// Convert the userId to string then pass to repo
	return s.repo.FindRoleByUserId(fmt.Sprint(userId))
}
func (s *authPolicyService) CreateRole(userId, roleToApply string) (*bool, error) {
	return s.repo.CreateRole(userId, roleToApply)
}
func (s *authPolicyService) AssignUserRole(userId, roleToApply string) (*bool, error) {
	success, err := s.repo.AssignUserRole(userId, roleToApply)
	if err != nil {
		return nil, err
	}
	return success, nil
}

// Inheritance
//

func (s *authPolicyService) FindAllRoleInheritance() ([]models.GRecord, error) {
	return s.repo.FindAllRoleInheritance()
}
func (s *authPolicyService) CreateInheritance(inherit models.GRecord) error {
	return s.repo.CreateInheritance(inherit)
}

func (s *authPolicyService) DeleteInheritance(inherit models.GRecord) error {
	return s.repo.DeleteInheritance(inherit)
}

// Transform data from enforcer policies to User friendly response
// (removes prefix from role as well)
func transformDataToResponse(data [][]string) []models.PolicyRuleCombinedActions {
	// Response format
	var response []models.PolicyRuleCombinedActions
	// Init policy dictionary for sorting
	policyDict := make(map[string]*models.PolicyRuleCombinedActions)

	// Loop through data and build policy dictionary
	for _, item := range data {
		// Assign policy vars
		role, resource, action := item[0], item[1], item[2]
		key := role + resource

		// If key does not exist, create new entry
		if _, ok := policyDict[key]; !ok {
			policyDict[key] = &models.PolicyRuleCombinedActions{
				// Remove role: prefix
				Role:     strings.TrimPrefix(role, "role:"),
				Resource: resource,
				Action:   []string{action},
			}

		} else {
			// Else, if record exists with resource, append action to action slice
			policyDict[key].Action = append(policyDict[key].Action, action)
		}
	}

	// Loop through policyDict and append to response
	for _, policy := range policyDict {
		response = append(response, *policy)
	}

	return response
}
