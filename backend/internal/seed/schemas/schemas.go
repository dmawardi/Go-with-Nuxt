package schemas

// FactoryRegistration is a struct that holds the factory and the name of the factory for factory registration
type FactoryRegistration struct {
	Factory BasicFactory
	Name    string
}

// BasicFactory is an interface that defines the methods that a factory must implement
type BasicFactory interface {
	Factory(count int) error
}
