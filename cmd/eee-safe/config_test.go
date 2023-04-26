package main

import (
	"github.com/joeig/eee-safe/pkg/storage/dynamodb"
	"github.com/joeig/eee-safe/pkg/storage/filesystem"
	"reflect"
	"testing"
)

func Test_mapStorageBackendType(t *testing.T) {
	backends := &StorageBackends{
		Filesystem: filesystem.Filesystem{},
		DynamoDB:   dynamodb.DynamoDB{},
	}

	type args struct {
		config   *Config
		backends *StorageBackends
	}

	tests := []struct {
		name string
		args args
		want StorageBackend
	}{
		{
			name: "filesystem",
			args: args{
				config:   &Config{StorageBackendType: StorageBackendTypeFilesystem},
				backends: backends,
			},
			want: &backends.Filesystem,
		},
		{
			name: "dynamoDB",
			args: args{
				config:   &Config{StorageBackendType: StorageBackendTypeDynamoDB},
				backends: backends,
			},
			want: &backends.DynamoDB,
		},
		{
			name: "unknown",
			args: args{
				config:   &Config{StorageBackendType: "unknown"},
				backends: backends,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := mapStorageBackendType(tt.args.config, tt.args.backends); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapStorageBackendType() = %v, want %v", got, tt.want)
			}
		})
	}
}
