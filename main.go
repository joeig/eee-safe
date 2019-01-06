package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/debug"
	"log"
	"os"
)

// BuildVersion is set at linking time
var BuildVersion string

// BuildGitCommit is set at linking time
var BuildGitCommit string

func main() {
	// Command line flags
	configFile := flag.String("config", "config.yml", "Configuration file")
	debugFlag := flag.Bool("debug", false, "Debug mode")
	version := flag.Bool("version", false, "Prints the version name")
	flag.Parse()

	// Version
	if *version {
		fmt.Printf("Build Version: %s\n", BuildVersion)
		fmt.Printf("Build Git Commit: %s\n", BuildGitCommit)
		os.Exit(0)
	}

	// Initialize configuration
	parseConfig(&config, configFile)
	setStorageBackend(&storageBackend)
	debug.Debug = *debugFlag
	if debug.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// Initialize Gin router
	router := getGinEngine()

	// Run server
	if config.Server.TLS.Enable {
		log.Fatal(router.RunTLS(config.Server.ListenAddress, config.Server.TLS.CertFile, config.Server.TLS.KeyFile))
	}
	log.Fatal(router.Run(config.Server.ListenAddress))
}
