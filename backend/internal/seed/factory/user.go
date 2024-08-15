package factory

import (
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/seed/schemas"
	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type userFactory struct {
	db *gorm.DB
}

func NewUserFactory(db *gorm.DB) schemas.BasicFactory {
	return &userFactory{db: db}
}

// factory generates and returns a slice of randomly generated users.
// The number of users generated is determined by the `count` parameter.
func (f userFactory) GenerateRandom() *db.User {
	user := &db.User{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: faker.Word(),
	}
	return user
}

// Factory inserts uses Generate method to insert data into the database
func (f userFactory) Factory(count int) error {
	// Create an empty slice
	users := []interface{}{}
	// Loop to generate users and append them to the slice
	for i := 0; i < count; i++ {
		user := f.GenerateRandom()
		users = append(users, user)
	}
	// Insert into the database
	return data.InsertUniqueRecord(f.db, users)
}
