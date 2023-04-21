package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/storage/dynamodb"
	"github.com/joeig/eee-safe/pkg/storage/filesystem"
	"github.com/spf13/viper"
)

// Config contains the primary configuration structure of the application
type Config struct {
	Server             Server          `mapstructure:"server"`
	StorageBackendType string          `mapstructure:"storageBackendType"`
	StorageBackends    StorageBackends `mapstructure:"storageBackends"`
}

// Read reads uses a Viper to read a configuration file into a Config instance.
func (c *Config) Read(viperCtx *viper.Viper, configFile string) error {
	viperCtx.SetConfigFile(configFile)

	if err := viperCtx.ReadInConfig(); err != nil {
		return err
	}

	if err := viperCtx.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

// Server defines the structure of the server configuration
type Server struct {
	ListenAddress string       `mapstructure:"listenaddress"`
	TLS           TLS          `mapstructure:"tls"`
	Accounts      gin.Accounts `mapstructure:"accounts"`
	Backups       Backups      `mapstructure:"backups"`
}

// TLS defines the structure of the TLS configuration
type TLS struct {
	Enable   bool   `mapstructure:"enable"`
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
}

// Backups defines the structure of the backup configuration
type Backups struct {
	MaxBackupBytes uint `mapstructure:"maxBackupBytes"`
	RetentionDays  uint `mapstructure:"retentionDays"`
}

// StorageBackends defines the structure for the storage backend configuration
type StorageBackends struct {
	Filesystem filesystem.Filesystem `mapstructure:"filesystem"`
	DynamoDB   dynamodb.DynamoDB     `mapstructure:"dynamodb"`
}

const (
	// StorageBackendTypeFilesystem sets the storage backend type to "filesystem"
	StorageBackendTypeFilesystem = "filesystem"
	// StorageBackendTypeDynamoDB sets the storage backend type to "filesystem"
	StorageBackendTypeDynamoDB = "dynamodb"
)

func mapStorageBackendType(config *Config, backends *StorageBackends) StorageBackend {
	switch config.StorageBackendType {
	case StorageBackendTypeFilesystem:
		if config.Server.Backups.RetentionDays != 0 {
			log.Println("Filesystem storage backend does currently not support backup retention.")
		}

		return &backends.Filesystem

	case StorageBackendTypeDynamoDB:
		backends.DynamoDB.InitializeService(session.Must(session.NewSession()))
		return &backends.DynamoDB

	default:
		return nil
	}
}
