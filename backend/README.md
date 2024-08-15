# Go-Template

## Written in GoLang

This repository will serve as a base web application in Go.

- Built in Go version 1.19
- Uses the [chi](https://github.com/go-chi/chi/v5/v5) router
- Uses [godotenv](https://github.com/joho/godotenv) for environmental variables
- Uses [Swaggo](https://github.com/swaggo/swag) to generate API documentation
- Uses [Go-Validator](https://github.com/asaskevich/govalidator) for validating incoming data

### Environment Setup

Create a .env file in the root of the project folder with the following content:

```
# Database Settings
DB_USER=postgres
DB_PASS=
DB_HOST=localhost
DB_PORT=5432
DB_NAME=
SESSIONS_SECRET_KEY=
HMAC_SECRET=
# SMTP Settings
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
```

### Database (Object Relational Management)

- Uses [Gorm](https://gorm.io) for ORM (Postgres)

### Security / Authentication / Authorization

- Password hashing with [bcrypt](https://golang.org/x/crypto)
- JWT authentication with [Golang-jwt](https://github.com/golang-jwt/jwt)
- Role-based access control with [casbin](https://github.com/casbin/casbin/v2)

## Running the Server

```
go run ./cmd
```

---

## Adding a Feature

1. Schema: Define schema in ./internal/db according to ORM instructions below with required receiver functions for admin panel (getID() and ObtainValue()).
2. Repository: Create a repository in ./internal/repository.
3. DTO Models: Build DTO models in ./internal/models.
4. Service: Develop the service in ./internal/service.
5. Controller: Implement the controller in ./internal/controller that accepts the request, performs data validation, then sends to the service to interact with database.
6. Validation: Add validation using govalidator in DTO definitions. This is done by adding `valid:""` key-value pairs to struct DTO definitions (/internal/models) that are being passed into the ValidateStruct function (used in controllers).
7. Routes: Update the modulesToSetup variable in the ./internal/modules/modules.go file with the created repo, service, and controller. This is used within the Routes function (./internal/routes/routes.go) automatically upon server run through the modules.SetupModules function where a modulemap is created and fed into API struct creation.
8. RBAC Policy: Add routes to the RBAC policy file (./internal/auth/rbac_policy.go).

### ADMIN PANEL

1. Admin Panel: Add admin panel file for created schema in ./internal/admin-panel folder as adminController<u>SchemaName</u>.go This file will contain the basicAdminController constructor with all the relevant details.
2. Admin Route Preparation: The db schema will need to fit the specs of the db.AdminPanelSchema, so add two receiver functions to your schema struct (in ./internal/db/<u>SchemaName</u>.go) as below.

- The first will be a function that returns the ID of the schema (GetId())
- The second will be a function that returns the value of the field given a key (ObtainValue())
  These functions will be used in the admin panel to display the data.

3. Admin Route Creation: Update the modulesToSetup in ./internal/modules/modules.go to include the new Admin panel controller.
4. Admin RBAC Policy: Add the admin routes to the RBAC authorization policy file (./internal/auth/rbac_policy.go)

5. Tests: For e2e testing, you will need to update the controllers_test.go file in ./internal/controller. Updates are required in the controllerTestModule struct, TestApiSetup, & setupTestDatabase functions. You will need to create a new module to add to the Test Module struct. This module will contain the repository, service, and controller for the new feature. You will also need to add the new module to the controllerTestModule struct in the setupTestDatabase function.

---

## Testing

Run all tests:

```
go test ./...
```

This will run all files that match the testing file naming convention (\*\_test.go).

For detailed results and coverage report:

```
go test ./... -V -cover
```

- "-V" prints more detailed results
- "-cover" will provide a test coverage report

Tests for repositories, services, and controllers should be in their respective directories. The controllers folder consists of E2E tests.  
The setup for these tests is in the controllers_test.go file.

Upon adding a new module:
-Make sure to build a DB struct that contains the modules of: repo, service, & controller. This object should then be added to the testDbRepo which serves as the test connection that will be serving the requests for the DB and API.

---

## Caching

Caching is handled by the cache package. It is stored in the app state and can be used from within services to store and retrieve details.

The caching functions should be used in the service ideally in the below functions:
Find by ID: Cache store, Cache Load
Update: Cache store
Delete, Bulk Delete: Cache Delete

Convention:
To retrieve, use prefix of table struct

```Go
// Define a key with a naming convention
	cacheKey := fmt.Sprintf("user:%d", userId)
	// Attempt to load the user from the cache first
	cachedUser, found := app.Cache.Load(cacheKey)
	if found {
		// If found and not expired, return the cached user
		return cachedUser.(*models.UserWithRole), nil
	}
```

To store, use below convention

```Go
// Cache the full user with a TTL before returning
	// Assuming you have defined a duration for the TTL
	ttl := 10 * time.Minute // For example, cache for 10 minutes
	app.Cache.Store(cacheKey, fullUser, ttl)
```

## Job Queue

The queue is handled by the Queue package.
worker.go: Contains the job worker that will complete a job every 5 seconds
email.go: Contains email associated job processing code
queue.go: Contains code to init, add, process, and mark complete jobs.

The queue struct is within the db package in the job.go file.

A new queue is init within the API creation and an async worker is initialized at this point as well to handle jobs.

---

## API documentation

API documentation is auto generated using markdown within code. This is achieved using Swag.

The below commands must be used upon making changes to the API in order to regenerate the API docs.

- "-d" directory flag allows custom directory to be used
- "-g" flag allows direct pointing to the main.go file for generation of swagger annotations from files that are imported (ie. controllers, services, repositories, etc.)
- "--pd" flag parses dependecies as well
- "--parseInternal" flag parses internal packages

Generate and update API docs using Swag:

```
<!-- Generate docs from home folder -->
<!-- Remove old API docs folder in static -->
<!-- To move the generated folder to be accessible to users -->
swag init -d ./internal/controller -g ../../cmd/main.go --pd --parseInternal
rm -rf static/docs
mv docs static/
```

This will update API documentation generated in the ./docs folder. It is served on path /swagger

## To use Database ORM

To edit schemas: Go to ./internal/db/schemas.go

The schemas are Structs based off of gorm.Model.

After creating the schema in schemas.go, make sure to add it to the models slice at the top for automigation in db.go.

For the admin panel, you will need to add two receiver functions to your schema struct in order for it to adhere to the db.AdminPanelSchema interface.

Note for creating and updating using GORM: Relationship data that does not yet exist will be created as a new entry. However, if you try to edit an existing record, it will not allow you to.

## Role based access control (RBAC) settings

The authorization model is found in the ./internal/auth/rbac_model.conf file.
This data structure is used by the setupcasbin policy to implement policy in DB upon server start.

The default policy is found in the ./internal/auth/rbac_policy.csv file.

Format of policy: Subject, Object, Action ("Who" is accessing "DB object" to commit "CRUD action")

Eg. admin, /api/v1/users, POST

In the policy implementation above:
p = Used to assign permissions to roles
eg. Assigning read permission to user role for /api/me endpoint
| p type | v0 | v1 | v2 |
| ------ | ---- | ------- | ---- |
| p | user | /api/me | read |

g = Used to assign roles to users & used to assign permissions to roles
Each g record is either a role assignment to a user or a role inheritance to another role.
All roles start with a 'role:' prefix

A role assignment is considered a role inheritance to a user: eg. user with id 2 has moderator role
An inheritance record is considered a role inheritance to another role: eg. admin role has all permissions of moderator role

eg. Assigning a moderator role to user with id 2
| p type | v0 | v1 |
| ------ | ---- | ------- |
| g | 2 | role:moderator |
| g | role:admin | role:moderator |

What about record level control?
eg. User can only edit their own profile

For sake of flexibility, reducing casbin policy model complexity, and to avoid having to create a new policy for each new record, we will use custom application logic within handlers to check if the user is the owner of the record to allow passing of the request to the service.

Another option is to use the casbin policy to check if the user is the owner of the record. However, this would require a new policy for each new record type.

### Wild cards

Policy resources allow for wild cards to be more flexible. This is useful for allowing a role to access all records of a certain type.

eg. admin, /api/v1/users/_, POST
or admin, /api/v1/users/\*\* _/, GET

\*: would allow for the admin role to access all resources under /api/v1/users. However, this doesn't proceed past the first level of the resource.(till the next '/')
**: **would allow for the admin role to access all resources under /api/v1/users. This would proceed past the first level of the resource.

Request -> Authorization (does user have permission through role?) -> Validation (does user own record?) -> Service (perform CRUD operation)

## Running with Docker

To run the application within a Docker container, you will need to build the image and run the container.

When running Docker on a Mac with an ARM processor, you will need to use the buildx command to build the image for amd64. This is where the --platform option comes in handy.

"container-name" is typically the github address of your project. (ie. dmawardi/go-template)

Build Docker container:

```
docker build -t container-name .
```

For ARM processors:

```
docker buildx build --platform linux/amd64 -t container-name .
```

Runs docker image and matches port

```
docker run --publish 8080:8080 container-name
```

The below command can be used to combine the build and to run with database in unison.
Runs Docker compose to build the server and database and run together in a docker container.

```
docker-compose up --build
```

Docker compose will expose the database for admin management using pgAdmin. The database will be accessible on port 5432.

In order to run the Docker image on a server, you will need to push the image to a Docker registry (Docker Hub). This can be done using Docker Desktop

## To deploy container on a server

First, ensure that Docker is installed on the server.

Then, pull the Docker image on the server using the container name (dmawardi/go-template:latest) and the docker pull command:

```
docker pull container-name:version
```

Then, run the Docker image on the server using the following command:

```
docker run -d -p 8080:8080 container-name:version
```
