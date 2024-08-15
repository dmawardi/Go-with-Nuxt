package models

// Used for service and DB
type CasbinRule struct {
	PType string `json:"ptype" gorm:"size:100;uniqueIndex:unique_index" valid:"required,in(p|g|g2)"`
	// role
	V0 string `json:"v0" gorm:"size:100;uniqueIndex:unique_index" valid:"required"`
	// resource
	V1 string `json:"v1" gorm:"size:100;uniqueIndex:unique_index" valid:"required"`
	// action
	V2 string `json:"v2" gorm:"size:100;uniqueIndex:unique_index" valid:"in(read|create|update|delete)"`
}

type PolicyRule struct {
	Role     string `json:"role" valid:"required"`
	Resource string `json:"resource" valid:"required"`
	Action   string `json:"action" valid:"required,in(read|create|update|delete)"`
}

type UpdateCasbinRule struct {
	OldPolicy PolicyRule `json:"old_policy" valid:"required"`
	NewPolicy PolicyRule `json:"new_policy" valid:"required"`
}

type CasbinRoleAssignment struct {
	UserId string `json:"user_id"  valid:"required"`
	Role   string `json:"role"  valid:"required"`
}

// Received from service: { "role": "admin", "resource": "/api/gustav", "action": ["read", "update"] }
type ReceivedPolicyRule map[string]interface{}

// Service received format in struct form. Used in admin section for easier rendering
type PolicyRuleCombinedActions struct {
	Role     string   `json:"role" valid:"required"`
	Resource string   `json:"resource" valid:"required"`
	Action   []string `json:"action" valid:"required"`
}

type GRecord struct {
	Role         string `json:"role" valid:"required"`
	InheritsFrom string `json:"inherits_from" valid:"required"`
}
