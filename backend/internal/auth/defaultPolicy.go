package auth

type policySet struct {
	Subject string
	Object  string
	Action  string
}

type Policy struct {
	PType   string
	Subject string
	Object  string
	Action  string
}

type GroupingPolicy struct {
	PType string
	User  string
	Role  string
}
