package repository_test

import (
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_Create(t *testing.T) {
	// Build hashed password from user password input
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	// Create a new user
	user := &db.User{
		Email: "test@example.com",
		// Imitate bcrypt encryption from user service
		Password: string(hashedPassword),
	}

	// Test function
	createdUser, err := testModule.users.repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify that the created user has an ID
	if createdUser.ID == 0 {
		t.Error("created user should have an ID")
	}

	// Verify that the created user has a hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password")); err != nil {
		t.Errorf("created user has incorrect password hash: %v", err)
	}

	// Attempt to create duplicate user
	duplicateUser := &db.User{
		Email: "test@example.com",
		// Imitate bcrypt encryption from user service
		Password: string(hashedPassword),
	}
	_, err = testModule.users.repo.Create(duplicateUser)
	if err == nil {
		t.Fatalf("Creating duplicate email should have failed but it didn't: %v", err)
	}

	// Clean up: Delete created user
	testModule.dbClient.Delete(createdUser)
	// In case duplicate was created, delete
	testModule.dbClient.Delete(duplicateUser)
}

func TestUserRepository_FindById(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	foundUser, err := testModule.users.repo.FindById(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	helpers.CompareObjects(createdUser, foundUser, t, []string{"ID", "Email", "Username"})

	// Clean up: Delete created user
	testModule.dbClient.Delete(createdUser)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "elon@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	foundUser, err := testModule.users.repo.FindByEmail(createdUser.Email)
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	// Verify that the found user matches the original user
	helpers.CompareObjects(createdUser, foundUser, t, []string{"ID", "Email", "Username"})

	// Clean up: Delete created user
	testModule.dbClient.Delete(createdUser)
}

func TestUserRepository_FindByVerificationCode(t *testing.T) {
	// Generate verification code and set expiry
	verificationDetails, err := helpers.GenerateVerificationCodeAndSetExpiry()
	if err != nil {
		t.Fatalf("failed to generate verification code: %v", err)
	}
	userToCreate := &db.User{
		Username:               "Jabar",
		Email:                  "scribble@ymail.com",
		Password:               "password",
		Name:                   "Bamba",
		Verified:               verificationDetails.Verified,
		VerificationCode:       verificationDetails.VerificationCode,
		VerificationCodeExpiry: verificationDetails.VerificationCodeExpiry,
	}

	createdUser, err := hashPassAndGenerateUserInDb(userToCreate, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Test function
	foundUser, err := testModule.users.repo.FindByVerificationCode(createdUser.VerificationCode)
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	// Verify that the found user matches the original user
	helpers.CompareObjects(createdUser, foundUser, t, []string{"ID", "Email", "Username", "VerificationCode", "VerificationCodeExpiry"})

	// Clean up: Delete created user
	testModule.dbClient.Delete(createdUser)
}

func TestUserRepository_Delete(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "delete@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Delete the created user
	err = testModule.users.repo.Delete(int(createdUser.ID))
	if err != nil {
		t.Fatalf("failed to delete created user: %v", err)
	}

	// Check to see if user has been deleted
	_, err = testModule.users.repo.FindById(int(createdUser.ID))
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
}

func TestUserRepository_BulkDelete(t *testing.T) {
	createdUser1, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Joe",
		Email:    "swagodahlia@gmail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user1: %v", err)
	}

	createdUser2, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "resolutesoldier@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user2: %v", err)
	}

	// Test function
	err = testModule.users.repo.BulkDelete([]int{int(createdUser1.ID), int(createdUser2.ID)})
	if err != nil {
		t.Fatalf("failed to bulk delete users: %v", err)
	}

	users := &[]db.User{}

	// Check to see if users have been deleted
	result := testModule.dbClient.Find(users)
	if result.Error != nil {
		t.Fatalf("failed to find users: %v", result.Error)
	}

	// If records still found
	if len(*users) != 0 {
		t.Fatalf("failed to delete users: %+v", users)
	}
}

func TestUserRepository_Update(t *testing.T) {
	createdUser, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "update@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	createdUser.Username = "Al-Amal"

	updatedUser, err := testModule.users.repo.Update(int(createdUser.ID), createdUser)
	if err != nil {
		t.Fatalf("An error was encountered while updating: %v", err)
	}

	foundUser, err := testModule.users.repo.FindById(int(updatedUser.ID))
	if err != nil {
		t.Errorf("An error was encountered while finding updated user: %v", err)
	}

	// Verify that the found user matches the original user
	helpers.CompareObjects(createdUser, foundUser, t, []string{"ID", "Email", "Username", "Name"})

	// Clean up: Delete created user
	testModule.dbClient.Delete(updatedUser)
}

func TestUserRepository_FindAll(t *testing.T) {
	createdUser1, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Joe",
		Email:    "crazy@gmail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user1: %v", err)
	}

	createdUser2, err := hashPassAndGenerateUserInDb(&db.User{
		Username: "Jabar",
		Email:    "scuba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, t)
	if err != nil {
		t.Fatalf("failed to create test user2: %v", err)
	}

	// Test function
	users, err := testModule.users.repo.FindAll(10, 0, "", []models.QueryConditionParameters{})
	if err != nil {
		t.Fatalf("failed to find all: %v", err)
	}
	// Make sure both users are in database
	if len(*users.Data) != 2 {
		t.Fatalf("Length of []users is not as expected. Got: %v", len(*users.Data))
	}

	// Iterate through users and ensure details match
	for _, u := range *users.Data {
		// If it's the first user
		if int(u.ID) == int(createdUser1.ID) {
			helpers.CompareObjects(createdUser1, &u, t, []string{"ID", "Email", "Username", "Name"})
		} else {
			helpers.CompareObjects(createdUser2, &u, t, []string{"ID", "Email", "Username", "Name"})
		}
	}
	// Clean up created users
	usersToDelete := []db.User{{ID: createdUser1.ID}, {ID: createdUser2.ID}}
	testModule.dbClient.Delete(usersToDelete)
}

// Test helper function
func hashPassAndGenerateUserInDb(user *db.User, t *testing.T) (*db.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		t.Fatalf("Couldn't create user")
	}
	user.Password = string(hashedPass)
	return testModule.users.repo.Create(user)
}
