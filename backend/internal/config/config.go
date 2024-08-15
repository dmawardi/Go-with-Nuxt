package config

import (
	"context"
	"html/template"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dmawardi/Go-Template/internal/cache"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type AppConfig struct {
	// TemplateCache map[string]*template.Template
	// UseCache      bool
	InProduction   bool
	Ctx            context.Context
	DbClient       *gorm.DB
	Session        *sessions.CookieStore
	Auth           AuthEnforcer
	AdminTemplates *template.Template
	// Should be set to the base url of the app upon server start
	BaseURL string
	// Cache
	Cache *cache.CacheMap
	// Core modules
	User models.ModuleSet
	Policy models.ModuleSet
}

type AuthEnforcer struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}
