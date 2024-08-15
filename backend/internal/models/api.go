package models

import "net/http"

type BasicController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type BasicAdminController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// Bulk delete (from table)
	BulkDelete(w http.ResponseWriter, r *http.Request)
	// Success pages
	CreateSuccess(w http.ResponseWriter, r *http.Request)
	EditSuccess(w http.ResponseWriter, r *http.Request)
	DeleteSuccess(w http.ResponseWriter, r *http.Request)
	// Obtain URL details for sidebar
	ObtainUrlDetails() URLDetails
}

type BasicService interface {
	FindAll() ([]interface{}, error)
	Find(id string) (interface{}, error)
	Create(data interface{}) (interface{}, error)
	Update(id string, data interface{}) (interface{}, error)
	Delete(id string) error
}

type BasicRepository interface {
	FindAll() ([]interface{}, error)
	Find(id string) (interface{}, error)
	Create(data interface{}) (interface{}, error)
	Update(id string, data interface{}) (interface{}, error)
	Delete(id string) error
}
