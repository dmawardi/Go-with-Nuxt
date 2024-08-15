package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connects to database and returns client
func DbConnect(useVerboseLogger bool) *gorm.DB {
	// Grab environment variables for connection
	var DB_USER string = os.Getenv("DB_USER")
	var DB_PASS string = os.Getenv("DB_PASS")
	var DB_HOST string = os.Getenv("DB_HOST")
	var DB_PORT string = os.Getenv("DB_PORT")
	var DB_NAME string = os.Getenv("DB_NAME")
	// Create connection string
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)

	// Build config based on verbose logger parameter
	var config = &gorm.Config{}
	if !useVerboseLogger {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbUrl), config)
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	for _, table := range Models {
		if err := db.AutoMigrate(table); err != nil {
			fmt.Printf("failed to migrate database schema: %v", err)
		}
	}

	return db
}

// Extract pointer value as string using data type (used in ObtainValue)
func PointerToStringWithType(ptr interface{}, dataType string) string {
	switch dataType {
	case "bool":
		if val, ok := ptr.(*bool); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%t", *val)
		}
	case "int":
		if val, ok := ptr.(*int); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%d", *val)
		}
	case "float64":
		if val, ok := ptr.(*float64); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%f", *val)
		}
	case "string":
		if val, ok := ptr.(*string); ok {
			if val == nil {
				return "nil"
			}
			return *val
		}
	}

	return ""
}
