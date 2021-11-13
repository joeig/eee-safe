package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
	"github.com/spf13/viper"
)

// StorageBackend is an interface for basic storage operations
type StorageBackend interface {
	PutBackup(ctx context.Context, backupInput *threema.BackupInput) error
	GetBackup(ctx context.Context, backupID threema.BackupID) (*threema.BackupOutput, error)
	DeleteBackup(ctx context.Context, backupID threema.BackupID) error
}

// AppContext contains the application context.
type AppContext struct {
	Config             *Config
	StorageBackend     StorageBackend
	RequestIDGenerator RequestIDGenerator
}

// InitializeStorageBackend takes the configured storage backend type and initializes the storage backend.
func (a *AppContext) InitializeStorageBackend() error {
	storageBackend := mapStorageBackendType(a.Config, &a.Config.StorageBackends)
	if storageBackend == nil {
		return errors.New("invalid storage backend")
	}

	a.StorageBackend = storageBackend

	return nil
}

// BuildVersion is set at linking time
var BuildVersion string

// BuildGitCommit is set at linking time
var BuildGitCommit string

func main() {
	initialiseSeed()

	configFile := flag.String("config", "config.yml", "Configuration file")
	debug := flag.Bool("debug", false, "Debug mode")
	version := flag.Bool("version", false, "Prints the version name")
	flag.Parse()

	if *version {
		printVersionAndExit(BuildVersion, BuildGitCommit)
	}

	appCtx := &AppContext{
		Config:             &Config{},
		RequestIDGenerator: NewRandomRequestIDGenerator(),
	}

	if err := appCtx.Config.Read(viper.New(), *configFile); err != nil {
		panic(err)
	}

	if err := appCtx.InitializeStorageBackend(); err != nil {
		panic(err)
	}

	if *debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	router := getGinEngine(appCtx)

	if appCtx.Config.Server.TLS.Enable {
		log.Fatal(router.RunTLS(
			appCtx.Config.Server.ListenAddress,
			appCtx.Config.Server.TLS.CertFile,
			appCtx.Config.Server.TLS.KeyFile,
		))
	}

	log.Fatal(router.Run(appCtx.Config.Server.ListenAddress))
}

func printVersionAndExit(version, commit string) {
	log.Printf("Build Version: %s\n", version)
	log.Printf("Build Git Commit: %s\n", commit)
	os.Exit(0)
}
