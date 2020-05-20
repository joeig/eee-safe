package main

import "github.com/gin-gonic/gin"

// Initializes the Gin engine
func getGinEngine() *gin.Engine {
	router := gin.Default()
	router.Use(requestIDMiddleware())
	router.Use(gin.BasicAuth(config.Server.Accounts))
	router.Use(validateUserAgentMiddleware())
	router.GET("/health", GetHealthHandler)
	router.GET("/config", GetConfigHandler)
	router.PUT("/backups/:threemaSafeBackupID", PutBackupHandler)
	router.GET("/backups/:threemaSafeBackupID", GetBackupHandler)
	router.DELETE("/backups/:threemaSafeBackupID", DeleteBackupHandler)

	return router
}
