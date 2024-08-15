package auth

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dmawardi/Go-Template/internal/config"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var app *config.AppConfig

var JWTKey = []byte(os.Getenv("HMAC_SECRET"))

// Function called in main.go to connect app state to current file
func SetStateInAuth(a *config.AppConfig) {
	app = a
}

// Authorization

type AuthToken struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Setup RBAC enforcer based using gorm client. Connects to DB and builds base policy
func EnforcerSetup(db *gorm.DB, setupDefaultPolicy bool) (*config.AuthEnforcer, error) {
	// Grab environment variables for connection
	var DB_PORT string = os.Getenv("DB_PORT")

	// Build adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	// If error
	if err != nil {
		log.Fatal("Couldn't build adapter for enforcer: ", err, "\nDB PORT", DB_PORT)
		return nil, err
	}

	// Build path to policy model
	rbacModelPath := webapi.BuildPathFromWorkingDirectory("/internal/auth/rbac_model.conf")

	// Initialize RBAC Authorization
	enforcer, err := casbin.NewEnforcer(rbacModelPath, adapter)

	// If error
	if err != nil {
		log.Fatal("Couldn't build RBAC enforcer: ", err)
		return nil, err
	}

	// If setupDefaultPolicy is true
	if setupDefaultPolicy {
		// Create default policies if not already detected within system
		SetupDefaultCasbinPolicy(enforcer)
	}

	// else
	return &config.AuthEnforcer{Enforcer: enforcer, Adapter: adapter}, nil
}

// Generates a JSON web token based on user's details
func GenerateJWT(userID int, email, roleName string) (string, error) {
	// Build expiration time
	expirationTime := time.Now().Add(12 * time.Hour)

	// Build claims to be stored in token
	claims := &AuthToken{
		Email: email,
		// Convert ID to string
		UserID: fmt.Sprint(userID),
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create new token using built claims and signing method
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Decrypt token using key to generate string
	tokenString, err := authToken.SignedString(JWTKey)
	// If error
	if err != nil {
		return "", err
	}
	// else, return token string
	return tokenString, nil
}

// Validates and parses signed token and checks if expired
func ValidateAndParseToken(r *http.Request) (tokenData *AuthToken, err error) {
	// Grab request header
	header := r.Header
	// Extract token string from Authorization header by removing prefix "Bearer "
	_, tokenString, _ := strings.Cut(header.Get("Authorization"), " ")

	// If token string is empty
	if tokenString == "" {
		// Check if found in cookie (Admin Panel SSR)
		cookie, err := r.Cookie("jwt_token")
		// If error found
		if err != nil {
			// And error is no cookie found
			if err == http.ErrNoCookie {
				// No token in the Authorization header and no cookie
				return nil, errors.New("Authentication Token not detected")
			}
			// Else, return other error
			return nil, err
		}
		// Else, set token string from found cookie
		tokenString = cookie.Value
	}

	// Parse token string and claims. Filter through auth token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthToken{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTKey), nil
		},
	)
	if err != nil {
		err = errors.New("couldn't parse token")
		return &AuthToken{}, err
	}

	// Extract claims from parsed tocken
	claims, ok := token.Claims.(*AuthToken)
	// If failed
	if !ok {
		err = errors.New("couldn't parse claims")
		return &AuthToken{}, err
	}

	// Token expiry check
	// Generate the current time in numeric date format
	currentTime := jwt.NewNumericDate(time.Now())
	// Check if expired
	if claims.RegisteredClaims.ExpiresAt != nil && claims.RegisteredClaims.ExpiresAt.Before(currentTime.Time) {
		err = errors.New("token expired")
		return &AuthToken{}, err
	}
	// else return claims
	return claims, nil
}

// Takes the http method and returns a string based on it
// for authorization assessment
func ActionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

// Set up policy settings in DB for casbin rules
func SetupDefaultCasbinPolicy(enforcer *casbin.Enforcer) {
	// Build path to default policies CSV file
	pathToPolicies := webapi.BuildPathFromWorkingDirectory("/internal/auth/rbac_policy.csv")
	// Open the CSV file
	f, err := os.Open(pathToPolicies)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a new reader
	r := csv.NewReader(f)

	// Set the comment character for the reader to ignore
	r.Comment = '#'
	// Set the fields per record to -1 to allow for variable number of fields
	// This is because we have both policies and grouping policies in the same file
	r.FieldsPerRecord = -1

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err != nil {
			break
		}

		// Switch based on the value of the first column
		switch record[0] {
		case "p":
			// If the first column is "p", then it is a policy
			// Map the record to a Policy struct
			policy := Policy{
				PType:   record[0],
				Subject: record[1],
				Object:  record[2],
				Action:  record[3],
			}

			// Check if the policy already exists
			hasPolicy, err := enforcer.HasPolicy(policy.Subject, policy.Object, policy.Action)
			if err != nil {
				log.Printf("Error checking policy: %v", err)
				continue
			}

			// If the policy does not exist, add it
			if !hasPolicy {
				success, err := enforcer.AddPolicy(policy.Subject, policy.Object, policy.Action)
				if err != nil {
					log.Printf("Error adding policy: %v", err)
					continue
				}
				if !success {
					log.Printf("Policy was not added: %v", policy)
					continue
				}
			}

		case "g":
			// If the first column is "g", then it is a grouping policy
			// Map the record to a GroupingPolicy struct
			groupingPolicy := GroupingPolicy{
				PType: record[0],
				User:  record[1],
				Role:  record[2],
			}

			// Check if the grouping policy already exists
			hasGroupingPolicy, err := enforcer.HasGroupingPolicy(groupingPolicy.User, groupingPolicy.Role)
			if err != nil {
				log.Printf("Error checking grouping policy: %v", err)
				continue
			}

			// If the grouping policy does not exist, add it
			if !hasGroupingPolicy {
				success, err := enforcer.AddGroupingPolicy(groupingPolicy.User, groupingPolicy.Role)
				if err != nil {
					log.Printf("Error adding grouping policy: %v", err)
					continue
				}
				if !success {
					log.Printf("Grouping policy was not added: %v", groupingPolicy)
					continue
				}
			}

		case "g2":
			// Form a named grouping policy
			namedGroupingPolicy := GroupingPolicy{
				PType: record[0],
				User:  record[1],
				Role:  record[2],
			}

			// Check if the grouping policy already exists
			hasGroupingPolicy, err := enforcer.HasNamedGroupingPolicy(record[0], namedGroupingPolicy.User, namedGroupingPolicy.Role)
			if err != nil {
				log.Printf("Error checking named grouping policy: %v", err)
				continue
			}

			// If the grouping policy does not exist, add it
			if !hasGroupingPolicy {
				success, err := enforcer.AddNamedGroupingPolicy(record[0], namedGroupingPolicy.User, namedGroupingPolicy.Role)
				if err != nil {
					log.Printf("Error adding grouping policy: %v", err)
					continue
				}
				if !success {
					log.Printf("Grouping policy was not added: %v", namedGroupingPolicy)
					continue
				}
			}
		}
	}
}

// Used to set header in admin panel for SSR authentication
// Create and set jwt token for SSR authentication
func CreateAndSetHeaderCookie(w http.ResponseWriter, tokenString string) {
	// Create the cookie
	expire := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name: "jwt_token",
		// Token string contians user info
		Value:    tokenString,
		Expires:  expire,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
	}

	// Set the cookie in the response header
	http.SetCookie(w, &cookie)
}
