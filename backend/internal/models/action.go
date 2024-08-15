package models

type CreateAction struct {
	ActionType  string `json:"action_type,omitempty" valid:"required"`
	EntityType  string `json:"entity_type,omitempty" valid:"required"`
	EntityID    string   `json:"entity_id,omitempty" valid:"required"`
	Changes     string `json:"changes,omitempty" valid:""`
	Description string `json:"description,omitempty" valid:""`
	AdminID     uint   `json:"admin_id,omitempty" valid:"required"`
	IPAddress   string `json:"ip_address,omitempty" valid:""`
}

type UpdateAction struct {
	ActionType  string `json:"action_type,omitempty" valid:""`
	EntityType  string `json:"entity_type,omitempty" valid:""`
	EntityID    string   `json:"entity_id,omitempty" valid:""`
	Changes     string `json:"changes,omitempty" valid:""`
	Description string `json:"description,omitempty" valid:""`
}

type RecordedAction struct {
	ActionType  string `json:"action_type,omitempty" valid:""`
	EntityType  string `json:"entity_type,omitempty" valid:""`
	EntityID    string   `json:"entity_id,omitempty" valid:""`
	Changes     string `json:"changes,omitempty" valid:""`
}