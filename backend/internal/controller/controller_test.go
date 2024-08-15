package controller_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	adminpanel "github.com/dmawardi/Go-Template/internal/admin-panel"
	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/cache"
	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/dmawardi/Go-Template/internal/helpers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/modules"
	"github.com/dmawardi/Go-Template/internal/queue"
	"github.com/dmawardi/Go-Template/internal/repository"
	corerepositories "github.com/dmawardi/Go-Template/internal/repository/core"
	"github.com/dmawardi/Go-Template/internal/routes"
	coreservices "github.com/dmawardi/Go-Template/internal/service/core"

	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/service"
	"gorm.io/gorm"
)

var testModule controllerTestModule

var app config.AppConfig

type controllerTestModule struct {
	dbClient *gorm.DB
	users    userModule
	admin    adminpanel.AdminPanelController
	auth     authModule
	router   http.Handler
	api      routes.Api
	// For authentication mocking
	accounts userAccounts
}

// Module structures
type userModule struct {
	repo corerepositories.UserRepository
	serv coreservices.UserService
	cont core.UserController
}
type authModule struct {
	repo corerepositories.AuthPolicyRepository
	serv coreservices.AuthPolicyService
	cont core.AuthPolicyController
}

// Account structures
type userAccounts struct {
	admin dummyAccount
	user  dummyAccount
}
type dummyAccount struct {
	details *models.UserWithRole
	token   string
}

// Initial setup before running e2e tests in controllers_test package
func TestMain(m *testing.M) {
	fmt.Printf("Setting up test connection\n")
	// Set URL in app state
	app.BaseURL = webapi.BuildBaseUrl()

	// Setup DB
	testModule.dbClient = helpers.SetupTestDatabase()
	// Set Gorm client
	app.DbClient = testModule.dbClient

	fmt.Printf("Preparing to build enforcer\n")
	// Build enforcer
	enforcer, err := auth.EnforcerSetup(testModule.dbClient, true)
	if err != nil {
		fmt.Println("Error building enforcer")
	}
	// Set enforcer in state
	app.Auth.Enforcer = enforcer.Enforcer
	app.Auth.Adapter = enforcer.Adapter

	// Setup new cache
	app.Cache = &cache.CacheMap{}

	// Sync app in authentication package for usage in authentication functions
	SetAppWideState(&app)

	// build API for serving requests
	testModule.api = testModule.TestApiSetup(testModule.dbClient)
	testModule.router = testModule.api.Routes()

	// Setup accounts for mocking authentication
	testModule.setupDummyAccounts(&models.CreateUser{
		Username: "Jabar",
		Email:    "Jabal@ymail.com",
		Password: "password",
		Name:     "Bamba",
	}, &models.CreateUser{
		Username: "Jabar",
		Email:    "Juba@ymail.com",
		Password: "password",
		Name:     "Bamba",
	})

	// Run the rest of the tests
	exitCode := m.Run()
	// exit with the same exit code as the tests
	os.Exit(exitCode)
}

// Builds new API using routes package
func (t *controllerTestModule) TestApiSetup(client *gorm.DB) routes.Api {
	mail := &helpers.EmailMock{}

	// Create job queue
	jobQueue := queue.NewQueue(client, mail)
	// Setup module stack
	// Auth
	t.auth.repo = corerepositories.NewAuthPolicyRepository(client)
	t.auth.serv = coreservices.NewAuthPolicyService(t.auth.repo)
	t.auth.cont = core.NewAuthPolicyController(t.auth.serv)
	// Users
	t.users.repo = corerepositories.NewUserRepository(client)
	t.users.serv = coreservices.NewUserService(t.users.repo, t.auth.repo, jobQueue)
	t.users.cont = core.NewUserController(t.users.serv)

	// Action
	actionRepo := corerepositories.NewActionRepository(client)
	actionService := coreservices.NewActionService(actionRepo)
	adminActionController := adminpanel.NewAdminActionController(actionService)

	// Setup basic modules with new implementation
	moduleMap := modules.SetupModules(modules.ModulesToSetup, client, actionService)

	// Admin panel
	t.admin = adminpanel.NewAdminPanelController(
		adminpanel.NewAdminCoreController(t.users.serv),
		adminpanel.NewAdminUserController(t.users.serv, actionService),
		adminpanel.NewAdminAuthPolicyController(t.auth.serv),
		adminActionController,
		// Additional modules
		moduleMap,
	)

	// Generate admin sidebar list from admin controller
	adminpanel.GenerateAndSetAdminSidebar(t.admin)

	// Setup API using controllers
	api := routes.NewApi(
		t.admin,
		t.users.cont,
		t.auth.cont,
		moduleMap,
	)

	return api
}

// Setup functions
//
// Setup dummy admin and user account and apply to test connection
func (t *controllerTestModule) setupDummyAccounts(adminUser *models.CreateUser, basicUser *models.CreateUser) {
	adminUser.Role = "admin"
	// Build admin user
	createdAdminUser, adminToken := t.generateUserWithRoleAndToken(
		adminUser)
	// Store credentials
	t.accounts.admin.details = createdAdminUser
	t.accounts.admin.token = adminToken

	basicUser.Role = "user"
	// Build normal user
	createdBasicUser, userToken := t.generateUserWithRoleAndToken(
		basicUser)
	// Store credentials
	t.accounts.user.details = createdBasicUser
	t.accounts.user.token = userToken
}

// Generates a new user, changes its role to admin and returns it with token
func (t *controllerTestModule) generateUserWithRoleAndToken(user *models.CreateUser) (*models.UserWithRole, string) {
	// Create user (will create new roles if not found)
	// The creation of user is incorrect. not applying naming convention to role
	createdUser, err := t.users.serv.Create(user)
	if err != nil {
		fmt.Println("Failed to create user: ", err)
		return nil, ""

	}

	// If successful, generate token
	fmt.Println("Generating token for: ", createdUser.Email)
	// Set login status to true
	tokenString, err := auth.GenerateJWT(int(createdUser.ID), createdUser.Email, createdUser.Role)
	if err != nil {
		fmt.Println("Failed to create JWT")
	}

	// Add unhashed password to returned object
	createdUser.Password = user.Password
	// Send to user in body
	return createdUser, tokenString
}

// Sets app config state to all packages for usage
func SetAppWideState(appConfig *config.AppConfig) {
	controller.SetStateInHandlers(appConfig)
	auth.SetStateInAuth(appConfig)
	adminpanel.SetStateInAdminPanel(appConfig)
	service.SetAppConfig(appConfig)
	repository.SetAppConfig(appConfig)
	routes.BuildRouteState(appConfig)
}
