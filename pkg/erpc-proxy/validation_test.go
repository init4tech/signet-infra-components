package erpcproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErpcProxyComponentArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    ErpcProxyComponentArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid args",
			args: ErpcProxyComponentArgs{
				Namespace: "default",
				Name:      "erpc-proxy",
				Image:     "ghcr.io/erpc/erpc:latest",
				Config: ErpcProxyConfig{
					LogLevel: "info",
					Server: ErpcProxyServerConfig{
						HttpPortV4: 4000,
						MaxTimeout: "30s",
					},
					Projects: []ErpcProxyProjectConfig{
						{
							Id: "project1",
							Networks: []ErpcProxyNetworkConfig{
								{
									ChainId:      1,
									Architecture: "evm",
								},
							},
							Upstreams: []ErpcProxyUpstreamConfig{
								{
									Id:       "upstream1",
									Type:     "evm",
									Endpoint: "https://eth.example.com",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing namespace",
			args: ErpcProxyComponentArgs{
				Name:  "erpc-proxy",
				Image: "ghcr.io/erpc/erpc:latest",
				Config: ErpcProxyConfig{
					Projects: []ErpcProxyProjectConfig{
						{
							Id: "project1",
							Networks: []ErpcProxyNetworkConfig{
								{
									ChainId:      1,
									Architecture: "evm",
								},
							},
							Upstreams: []ErpcProxyUpstreamConfig{
								{
									Id:       "upstream1",
									Type:     "evm",
									Endpoint: "https://eth.example.com",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "namespace is required",
		},
		{
			name: "missing name",
			args: ErpcProxyComponentArgs{
				Namespace: "default",
				Image:     "ghcr.io/erpc/erpc:latest",
				Config: ErpcProxyConfig{
					Projects: []ErpcProxyProjectConfig{
						{
							Id: "project1",
							Networks: []ErpcProxyNetworkConfig{
								{
									ChainId:      1,
									Architecture: "evm",
								},
							},
							Upstreams: []ErpcProxyUpstreamConfig{
								{
									Id:       "upstream1",
									Type:     "evm",
									Endpoint: "https://eth.example.com",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing image",
			args: ErpcProxyComponentArgs{
				Namespace: "default",
				Name:      "erpc-proxy",
				Config: ErpcProxyConfig{
					Projects: []ErpcProxyProjectConfig{
						{
							Id: "project1",
							Networks: []ErpcProxyNetworkConfig{
								{
									ChainId:      1,
									Architecture: "evm",
								},
							},
							Upstreams: []ErpcProxyUpstreamConfig{
								{
									Id:       "upstream1",
									Type:     "evm",
									Endpoint: "https://eth.example.com",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "image is required",
		},
		{
			name: "invalid replicas",
			args: ErpcProxyComponentArgs{
				Namespace: "default",
				Name:      "erpc-proxy",
				Image:     "ghcr.io/erpc/erpc:latest",
				Replicas:  -1,
				Config: ErpcProxyConfig{
					Projects: []ErpcProxyProjectConfig{
						{
							Id: "project1",
							Networks: []ErpcProxyNetworkConfig{
								{
									ChainId:      1,
									Architecture: "evm",
								},
							},
							Upstreams: []ErpcProxyUpstreamConfig{
								{
									Id:       "upstream1",
									Type:     "evm",
									Endpoint: "https://eth.example.com",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "replicas must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyConfig{
				LogLevel: "info",
				Projects: []ErpcProxyProjectConfig{
					{
						Id: "project1",
						Networks: []ErpcProxyNetworkConfig{
							{
								ChainId:      1,
								Architecture: "evm",
							},
						},
						Upstreams: []ErpcProxyUpstreamConfig{
							{
								Id:       "upstream1",
								Type:     "evm",
								Endpoint: "https://eth.example.com",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: ErpcProxyConfig{
				LogLevel: "invalid",
				Projects: []ErpcProxyProjectConfig{
					{
						Id: "project1",
						Networks: []ErpcProxyNetworkConfig{
							{
								ChainId:      1,
								Architecture: "evm",
							},
						},
						Upstreams: []ErpcProxyUpstreamConfig{
							{
								Id:       "upstream1",
								Type:     "evm",
								Endpoint: "https://eth.example.com",
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name: "no projects",
			config: ErpcProxyConfig{
				LogLevel: "info",
				Projects: []ErpcProxyProjectConfig{},
			},
			wantErr: true,
			errMsg:  "at least one project is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyProjectConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyProjectConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyProjectConfig{
				Id: "project1",
				Networks: []ErpcProxyNetworkConfig{
					{
						ChainId:      1,
						Architecture: "evm",
					},
				},
				Upstreams: []ErpcProxyUpstreamConfig{
					{
						Id:       "upstream1",
						Type:     "evm",
						Endpoint: "https://eth.example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing project id",
			config: ErpcProxyProjectConfig{
				Networks: []ErpcProxyNetworkConfig{
					{
						ChainId:      1,
						Architecture: "evm",
					},
				},
				Upstreams: []ErpcProxyUpstreamConfig{
					{
						Id:       "upstream1",
						Type:     "evm",
						Endpoint: "https://eth.example.com",
					},
				},
			},
			wantErr: true,
			errMsg:  "project ID is required",
		},
		{
			name: "no networks",
			config: ErpcProxyProjectConfig{
				Id:       "project1",
				Networks: []ErpcProxyNetworkConfig{},
				Upstreams: []ErpcProxyUpstreamConfig{
					{
						Id:       "upstream1",
						Type:     "evm",
						Endpoint: "https://eth.example.com",
					},
				},
			},
			wantErr: true,
			errMsg:  "at least one network is required",
		},
		{
			name: "no upstreams",
			config: ErpcProxyProjectConfig{
				Id: "project1",
				Networks: []ErpcProxyNetworkConfig{
					{
						ChainId:      1,
						Architecture: "evm",
					},
				},
				Upstreams: []ErpcProxyUpstreamConfig{},
			},
			wantErr: true,
			errMsg:  "at least one upstream is required",
		},
		{
			name: "invalid network",
			config: ErpcProxyProjectConfig{
				Id: "project1",
				Networks: []ErpcProxyNetworkConfig{
					{
						ChainId:      0, // Invalid
						Architecture: "evm",
					},
				},
				Upstreams: []ErpcProxyUpstreamConfig{
					{
						Id:       "upstream1",
						Type:     "evm",
						Endpoint: "https://eth.example.com",
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid network at index 0",
		},
		{
			name: "invalid upstream",
			config: ErpcProxyProjectConfig{
				Id: "project1",
				Networks: []ErpcProxyNetworkConfig{
					{
						ChainId:      1,
						Architecture: "evm",
					},
				},
				Upstreams: []ErpcProxyUpstreamConfig{
					{
						Id:       "upstream1",
						Type:     "invalid", // Invalid type
						Endpoint: "https://eth.example.com",
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid upstream at index 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyDatabaseConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid postgres config",
			config: ErpcProxyDatabaseConfig{
				Type:          "postgres",
				ConnectionUrl: "postgres://user:pass@localhost/db",
			},
			wantErr: false,
		},
		{
			name: "valid memory config",
			config: ErpcProxyDatabaseConfig{
				Type: "memory",
			},
			wantErr: false,
		},
		{
			name: "invalid database type",
			config: ErpcProxyDatabaseConfig{
				Type: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid database type",
		},
		{
			name: "postgres without connection url",
			config: ErpcProxyDatabaseConfig{
				Type: "postgres",
			},
			wantErr: true,
			errMsg:  "connection URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyServerConfig{
				HttpPortV4: 4000,
				MaxTimeout: "30s",
			},
			wantErr: false,
		},
		{
			name: "invalid port - negative",
			config: ErpcProxyServerConfig{
				HttpPortV4: -1,
				MaxTimeout: "30s",
			},
			wantErr: true,
			errMsg:  "invalid HTTP port",
		},
		{
			name: "invalid port - too high",
			config: ErpcProxyServerConfig{
				HttpPortV4: 70000,
				MaxTimeout: "30s",
			},
			wantErr: true,
			errMsg:  "invalid HTTP port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyNetworkConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyNetworkConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyNetworkConfig{
				ChainId:      1,
				Architecture: "evm",
			},
			wantErr: false,
		},
		{
			name: "invalid chain id",
			config: ErpcProxyNetworkConfig{
				ChainId:      0,
				Architecture: "evm",
			},
			wantErr: true,
			errMsg:  "chain ID must be positive",
		},
		{
			name: "missing architecture",
			config: ErpcProxyNetworkConfig{
				ChainId: 1,
			},
			wantErr: true,
			errMsg:  "architecture is required",
		},
		{
			name: "invalid architecture",
			config: ErpcProxyNetworkConfig{
				ChainId:      1,
				Architecture: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid architecture",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyUpstreamConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyUpstreamConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyUpstreamConfig{
				Id:       "upstream1",
				Type:     "evm",
				Endpoint: "https://eth.example.com",
			},
			wantErr: false,
		},
		{
			name: "missing id",
			config: ErpcProxyUpstreamConfig{
				Type:     "evm",
				Endpoint: "https://eth.example.com",
			},
			wantErr: true,
			errMsg:  "upstream ID is required",
		},
		{
			name: "missing type",
			config: ErpcProxyUpstreamConfig{
				Id:       "upstream1",
				Endpoint: "https://eth.example.com",
			},
			wantErr: true,
			errMsg:  "upstream type is required",
		},
		{
			name: "invalid type",
			config: ErpcProxyUpstreamConfig{
				Id:       "upstream1",
				Type:     "invalid",
				Endpoint: "https://eth.example.com",
			},
			wantErr: true,
			errMsg:  "invalid upstream type",
		},
		{
			name: "missing endpoint",
			config: ErpcProxyUpstreamConfig{
				Id:   "upstream1",
				Type: "evm",
			},
			wantErr: true,
			errMsg:  "upstream endpoint is required",
		},
		{
			name: "negative max retries",
			config: ErpcProxyUpstreamConfig{
				Id:         "upstream1",
				Type:       "evm",
				Endpoint:   "https://eth.example.com",
				MaxRetries: -1,
			},
			wantErr: true,
			errMsg:  "max retries must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyFailoverConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyFailoverConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyFailoverConfig{
				MaxRetries:    3,
				BackoffMs:     100,
				BackoffMaxMs:  5000,
				BackoffFactor: 2,
			},
			wantErr: false,
		},
		{
			name: "negative max retries",
			config: ErpcProxyFailoverConfig{
				MaxRetries: -1,
			},
			wantErr: true,
			errMsg:  "max retries must be non-negative",
		},
		{
			name: "negative backoff",
			config: ErpcProxyFailoverConfig{
				BackoffMs: -1,
			},
			wantErr: true,
			errMsg:  "backoff must be non-negative",
		},
		{
			name: "max backoff less than backoff",
			config: ErpcProxyFailoverConfig{
				BackoffMs:    1000,
				BackoffMaxMs: 500,
			},
			wantErr: true,
			errMsg:  "max backoff must be greater than or equal to backoff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErpcProxyResources_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ErpcProxyResources
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ErpcProxyResources{
				MemoryRequest: "256Mi",
				MemoryLimit:   "2Gi",
				CpuRequest:    "100m",
				CpuLimit:      "1000m",
			},
			wantErr: false,
		},
		{
			name: "valid config with plain numbers",
			config: ErpcProxyResources{
				CpuRequest: "0.1",
				CpuLimit:   "1",
			},
			wantErr: false,
		},
		{
			name: "invalid memory format",
			config: ErpcProxyResources{
				MemoryRequest: "256invalid",
			},
			wantErr: true,
			errMsg:  "invalid memoryRequest format",
		},
		{
			name: "invalid cpu format",
			config: ErpcProxyResources{
				CpuRequest: "100invalid",
			},
			wantErr: true,
			errMsg:  "invalid cpuRequest format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
