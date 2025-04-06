package main

func (app *Config) setupRoutes() {
	app.router.POST("/log", app.WriteLog)
	app.router.GET("/logs", app.GetAllLogs)
}