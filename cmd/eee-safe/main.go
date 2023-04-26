package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
	"github.com/spf13/viper"
	"log"
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

func main() {
	configFile := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()

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

	gin.SetMode("release")

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
