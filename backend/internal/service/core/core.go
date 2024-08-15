package coreservices

import "github.com/dmawardi/Go-Template/internal/config"

var app *config.AppConfig

func SetAppConfig(appConfig *config.AppConfig) {
	app = appConfig
}
