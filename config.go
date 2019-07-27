package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/storage"
	"github.com/joeig/eee-safe/storage/dynamodb"
	"github.com/joeig/eee-safe/storage/filesystem"
	"github.com/spf13/viper"
)

// Config contains the primary configuration structure of the application
type Config struct {
	Server             Server                     `mapstructure:"server"`
	StorageBackendType storage.StorageBackendType `mapstructure:"storageBackendType"`
	StorageBackends    StorageBackends            `mapstructure:"storageBackends"`
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

var config Config

func parseConfig(config *Config, configFile *string) {
	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("%s", err))
	}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("%s", err))
	}
}

const (
	// StorageBackendTypeFilesystem sets the storage backend type to "filesystem"
	StorageBackendTypeFilesystem storage.StorageBackendType = "filesystem"
	// StorageBackendTypeDynamoDB sets the storage backend type to "filesystem"
	StorageBackendTypeDynamoDB storage.StorageBackendType = "dynamodb"
)

var storageBackend storage.StorageBackend

func setStorageBackend(s *storage.StorageBackend) {
	switch config.StorageBackendType {
	case StorageBackendTypeFilesystem:
		*s = &config.StorageBackends.Filesystem
		break
	case StorageBackendTypeDynamoDB:
		*s = &config.StorageBackends.DynamoDB
		break
	default:
		panic(fmt.Errorf("Invalid dnsProviderType \"%s\"", config.StorageBackendType))
	}
}
