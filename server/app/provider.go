package app

import "github.com/smtdfc/nagare/server/taskers"

type AppError struct {
	error
}

// @Injectable
// @Root
func StartApp(app *App, _ *AppRoute, _ taskers.PluginTasker) *AppError {
	return app.StartServer()
}
