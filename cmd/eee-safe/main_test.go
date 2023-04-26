package main

import "testing"

func TestAppContext_InitializeStorageBackend1(t *testing.T) {
	type fields struct {
		Config         *Config
		StorageBackend StorageBackend
	}

	tests := []struct {
		name               string
		fields             fields
		wantStorageBackend bool
		wantErr            bool
	}{
		{
			name: "filesystem",
			fields: fields{
				Config: &Config{StorageBackendType: StorageBackendTypeFilesystem},
			},
			wantStorageBackend: true,
			wantErr:            false,
		},
		{
			name: "unknown",
			fields: fields{
				Config: &Config{StorageBackendType: "unknown"},
			},
			wantStorageBackend: false,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := &AppContext{
				Config:         tt.fields.Config,
				StorageBackend: tt.fields.StorageBackend,
			}
			if err := a.InitializeStorageBackend(); (err != nil) != tt.wantErr {
				t.Errorf("InitializeStorageBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantStorageBackend && a.StorageBackend == nil {
				t.Error("storage backend is nil")
			}
			if !tt.wantStorageBackend && a.StorageBackend != nil {
				t.Error("storage backend is not nil")
			}
		})
	}
}
