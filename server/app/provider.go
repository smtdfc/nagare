package app

type AppError struct {
	error
}

// @Injectable
// @Root
func StartApp(app *App) *AppError {
	return app.StartServer()
}
