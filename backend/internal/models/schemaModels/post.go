package schemamodels

import (
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

type CreatePost struct {
	Title string  `json:"title,omitempty" valid:"length(3|36),required"`
	Body  string  `json:"body,omitempty" valid:"length(10|1200),required"`
	User  db.User `json:"user,omitempty" valid:"required"`
}

type UpdatePost struct {
	Title string  `json:"title,omitempty" valid:"length(3|36)"`
	Body  string  `json:"body,omitempty" valid:"length(10|1200)"`
	User  db.User `json:"user,omitempty" valid:""`
}

type PaginatedPosts struct {
	Data *[]db.Post            `json:"data"`
	Meta models.SchemaMetaData `json:"meta"`
}
