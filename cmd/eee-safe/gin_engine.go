package main

import "github.com/gin-gonic/gin"

// Initializes the Gin engine
func getGinEngine(appCtx *AppCtx) *gin.Engine {
	router := gin.Default()
	router.Use(requestIDMiddleware())
	router.Use(gin.BasicAuth(appCtx.Config.Server.Accounts))
	router.Use(validateUserAgentMiddleware())
	router.GET("/health", appCtx.GetHealthHandler)
	router.GET("/config", appCtx.GetConfigHandler)
	router.PUT("/backups/:threemaSafeBackupID", appCtx.PutBackupHandler)
	router.GET("/backups/:threemaSafeBackupID", appCtx.GetBackupHandler)
	router.DELETE("/backups/:threemaSafeBackupID", appCtx.DeleteBackupHandler)

	return router
}
