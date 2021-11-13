package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/debug"
	"github.com/joeig/eee-safe/pkg/threema"
	"github.com/spf13/viper"
)

// StorageBackend is an interface for basic storage operations
type StorageBackend interface {
	PutBackup(backupInput *threema.BackupInput) error
	GetBackup(backupID threema.BackupID) (*threema.BackupOutput, error)
	DeleteBackup(backupID threema.BackupID) error
}

// AppCtx contains the application context.
type AppCtx struct {
	Config *Config
}

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
		printVersionAndExit(BuildVersion, BuildGitCommit)
	}

	// Initialize the application context
	appCtx := &AppCtx{
		Config: &Config{},
	}

	if err := appCtx.Config.Read(viper.New(), *configFile); err != nil {
		panic(err)
	}

	setStorageBackend(appCtx.Config, &storageBackend)

	debug.Debug = *debugFlag
	if debug.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// Initialize Gin router
	router := getGinEngine(appCtx)

	// Run server
	if appCtx.Config.Server.TLS.Enable {
		log.Fatal(router.RunTLS(appCtx.Config.Server.ListenAddress, appCtx.Config.Server.TLS.CertFile, appCtx.Config.Server.TLS.KeyFile))
	}

	log.Fatal(router.Run(appCtx.Config.Server.ListenAddress))
}

func printVersionAndExit(version, commit string) {
	fmt.Printf("Build Version: %s\n", version)
	fmt.Printf("Build Git Commit: %s\n", commit)
	os.Exit(0)
}
