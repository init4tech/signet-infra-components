package erpcproxy

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMarshalErpcConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   ErpcProxyConfig
		expected map[string]interface{}
	}{
		{
			name: "minimal config",
			config: ErpcProxyConfig{
				LogLevel: "info",
				Projects: []ErpcProxyProjectConfig{
					{
						Id: "main",
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
			expected: map[string]interface{}{
				"logLevel": "info",
				"projects": []interface{}{
					map[string]interface{}{
						"id": "main",
						"networks": []interface{}{
							map[string]interface{}{
								"evm": map[string]interface{}{
									"chainId": 1,
								},
								"architecture": "evm",
							},
						},
						"upstreams": []interface{}{
							map[string]interface{}{
								"id":       "upstream1",
								"type":     "evm",
								"endpoint": "https://eth.example.com",
							},
						},
					},
				},
			},
		},
		{
			name: "config with CORS",
			config: ErpcProxyConfig{
				LogLevel: "debug",
				Server: ErpcProxyServerConfig{
					HttpHostV4: "0.0.0.0",
					HttpPortV4: 4000,
					MaxTimeout: "30s",
				},
				Database: ErpcProxyDatabaseConfig{
					Type:          "postgresql",
					ConnectionUrl: "postgres://user:pass@localhost/db",
				},
				Projects: []ErpcProxyProjectConfig{
					{
						Id:              "main",
						RateLimitBudget: "frontend-budget",
						Cors: &ErpcProxyCorsConfig{
							AllowedOrigins:   []string{"https://example.com", "https://*.example.com"},
							AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
							AllowedHeaders:   []string{"Content-Type", "Authorization"},
							ExposedHeaders:   []string{"X-Request-ID"},
							AllowCredentials: true,
							MaxAge:           3600,
						},
						Networks: []ErpcProxyNetworkConfig{
							{
								ChainId:      1,
								Architecture: "evm",
								Failover: ErpcProxyFailoverConfig{
									MaxRetries:    3,
									BackoffMs:     100,
									BackoffMaxMs:  5000,
									BackoffFactor: 2,
									Duration:      "30s",
								},
							},
						},
						Upstreams: []ErpcProxyUpstreamConfig{
							{
								Id:              "upstream1",
								Type:            "evm",
								Endpoint:        "https://eth.example.com",
								RateLimitBudget: "global-budget",
								MaxRetries:      2,
								Timeout:         "15s",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"logLevel": "debug",
				"server": map[string]interface{}{
					"httpHostV4": "0.0.0.0",
					"httpPortV4": 4000,
					"maxTimeout": "30s",
				},
				"database": map[string]interface{}{
					"type":          "postgresql",
					"connectionUrl": "postgres://user:pass@localhost/db",
				},
				"projects": []interface{}{
					map[string]interface{}{
						"id":              "main",
						"rateLimitBudget": "frontend-budget",
						"cors": map[string]interface{}{
							"allowedOrigins":   []interface{}{"https://example.com", "https://*.example.com"},
							"allowedMethods":   []interface{}{"GET", "POST", "OPTIONS"},
							"allowedHeaders":   []interface{}{"Content-Type", "Authorization"},
							"exposedHeaders":   []interface{}{"X-Request-ID"},
							"allowCredentials": true,
							"maxAge":           3600,
						},
						"networks": []interface{}{
							map[string]interface{}{
								"evm": map[string]interface{}{
									"chainId": 1,
								},
								"architecture": "evm",
								"failover": map[string]interface{}{
									"maxRetries":    3,
									"backoffMs":     100,
									"backoffMaxMs":  5000,
									"backoffFactor": 2,
									"duration":      "30s",
								},
							},
						},
						"upstreams": []interface{}{
							map[string]interface{}{
								"id":              "upstream1",
								"type":            "evm",
								"endpoint":        "https://eth.example.com",
								"rateLimitBudget": "global-budget",
								"maxRetries":      2,
								"timeout":         "15s",
							},
						},
					},
				},
			},
		},
		{
			name: "config with partial CORS",
			config: ErpcProxyConfig{
				LogLevel: "warn",
				Projects: []ErpcProxyProjectConfig{
					{
						Id: "main",
						Cors: &ErpcProxyCorsConfig{
							AllowedOrigins: []string{"https://example.com"},
							AllowedMethods: []string{"GET", "POST"},
							// Other fields left empty to test partial marshalling
						},
						Networks: []ErpcProxyNetworkConfig{
							{
								ChainId:      42161,
								Architecture: "evm",
							},
						},
						Upstreams: []ErpcProxyUpstreamConfig{
							{
								Id:       "upstream1",
								Type:     "evm",
								Endpoint: "https://arb.example.com",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"logLevel": "warn",
				"projects": []interface{}{
					map[string]interface{}{
						"id": "main",
						"cors": map[string]interface{}{
							"allowedOrigins": []interface{}{"https://example.com"},
							"allowedMethods": []interface{}{"GET", "POST"},
						},
						"networks": []interface{}{
							map[string]interface{}{
								"evm": map[string]interface{}{
									"chainId": 42161,
								},
								"architecture": "evm",
							},
						},
						"upstreams": []interface{}{
							map[string]interface{}{
								"id":       "upstream1",
								"type":     "evm",
								"endpoint": "https://arb.example.com",
							},
						},
					},
				},
			},
		},
		{
			name: "config without CORS",
			config: ErpcProxyConfig{
				LogLevel: "error",
				Projects: []ErpcProxyProjectConfig{
					{
						Id: "main",
						// No CORS config
						Networks: []ErpcProxyNetworkConfig{
							{
								ChainId:      10,
								Architecture: "evm",
							},
						},
						Upstreams: []ErpcProxyUpstreamConfig{
							{
								Id:       "upstream1",
								Type:     "evm",
								Endpoint: "https://opt.example.com",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"logLevel": "error",
				"projects": []interface{}{
					map[string]interface{}{
						"id": "main",
						// No cors field should be present
						"networks": []interface{}{
							map[string]interface{}{
								"evm": map[string]interface{}{
									"chainId": 10,
								},
								"architecture": "evm",
							},
						},
						"upstreams": []interface{}{
							map[string]interface{}{
								"id":       "upstream1",
								"type":     "evm",
								"endpoint": "https://opt.example.com",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the config
			yamlStr, err := marshalErpcConfig(tt.config)
			assert.NoError(t, err)
			assert.NotEmpty(t, yamlStr)

			// Parse the YAML back to verify structure
			var result map[string]interface{}
			err = yaml.Unmarshal([]byte(yamlStr), &result)
			assert.NoError(t, err)

			// Verify the structure matches expected
			assert.Equal(t, tt.expected, result)

			// Additional validation: ensure YAML is valid and contains expected sections
			assert.Contains(t, yamlStr, "projects:")
			assert.Contains(t, yamlStr, "logLevel:")

			// For CORS tests, verify CORS section is present/absent as expected
			if tt.config.Projects[0].Cors != nil {
				assert.Contains(t, yamlStr, "cors:")
				assert.Contains(t, yamlStr, "allowedOrigins:")
				assert.Contains(t, yamlStr, "allowedMethods:")
			} else {
				assert.NotContains(t, yamlStr, "cors:")
			}

			// Verify YAML structure conforms to eRPC specification
			// Check that projects array is properly formatted
			assert.True(t, strings.Contains(yamlStr, "- id:") || strings.Contains(yamlStr, "id:"))

			// Check that networks are properly nested under projects
			assert.Contains(t, yamlStr, "networks:")
			assert.Contains(t, yamlStr, "upstreams:")
		})
	}
}

func TestMarshalErpcConfig_CORSValidation(t *testing.T) {
	t.Run("CORS with all fields", func(t *testing.T) {
		config := ErpcProxyConfig{
			LogLevel: "debug",
			Projects: []ErpcProxyProjectConfig{
				{
					Id: "main",
					Cors: &ErpcProxyCorsConfig{
						AllowedOrigins:   []string{"https://example.com", "https://*.example.com"},
						AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
						AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Custom-Header"},
						ExposedHeaders:   []string{"X-Request-ID", "X-Response-Time"},
						AllowCredentials: true,
						MaxAge:           86400,
					},
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
		}

		yamlStr, err := marshalErpcConfig(config)
		assert.NoError(t, err)

		// Verify all CORS fields are present in YAML
		assert.Contains(t, yamlStr, "allowedOrigins:")
		assert.Contains(t, yamlStr, "allowedMethods:")
		assert.Contains(t, yamlStr, "allowedHeaders:")
		assert.Contains(t, yamlStr, "exposedHeaders:")
		assert.Contains(t, yamlStr, "allowCredentials: true")
		assert.Contains(t, yamlStr, "maxAge: 86400")

		// Verify specific values
		assert.Contains(t, yamlStr, "https://example.com")
		assert.Contains(t, yamlStr, "https://*.example.com")
		assert.Contains(t, yamlStr, "Content-Type")
		assert.Contains(t, yamlStr, "X-Request-ID")
	})

	t.Run("CORS with empty arrays should be omitted", func(t *testing.T) {
		config := ErpcProxyConfig{
			LogLevel: "info",
			Projects: []ErpcProxyProjectConfig{
				{
					Id: "main",
					Cors: &ErpcProxyCorsConfig{
						AllowedOrigins:   []string{}, // Empty array
						AllowedMethods:   []string{}, // Empty array
						AllowedHeaders:   []string{}, // Empty array
						ExposedHeaders:   []string{}, // Empty array
						AllowCredentials: false,      // Default value
						MaxAge:           0,          // Default value
					},
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
		}

		yamlStr, err := marshalErpcConfig(config)
		assert.NoError(t, err)

		// CORS section should not be present since all fields are empty/default
		assert.NotContains(t, yamlStr, "cors:")
	})

	t.Run("CORS with only some fields populated", func(t *testing.T) {
		config := ErpcProxyConfig{
			LogLevel: "info",
			Projects: []ErpcProxyProjectConfig{
				{
					Id: "main",
					Cors: &ErpcProxyCorsConfig{
						AllowedOrigins:   []string{"https://example.com"},
						AllowCredentials: true,
						MaxAge:           3600,
						// Other fields left empty
					},
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
		}

		yamlStr, err := marshalErpcConfig(config)
		assert.NoError(t, err)

		// Only populated fields should be present
		assert.Contains(t, yamlStr, "cors:")
		assert.Contains(t, yamlStr, "allowedOrigins:")
		assert.Contains(t, yamlStr, "allowCredentials: true")
		assert.Contains(t, yamlStr, "maxAge: 3600")

		// Empty fields should not be present
		assert.NotContains(t, yamlStr, "allowedMethods:")
		assert.NotContains(t, yamlStr, "allowedHeaders:")
		assert.NotContains(t, yamlStr, "exposedHeaders:")
	})
}

func TestMarshalErpcConfig_ConformsToERPCSpec(t *testing.T) {
	t.Run("YAML structure matches eRPC specification", func(t *testing.T) {
		config := ErpcProxyConfig{
			LogLevel: "debug",
			Server: ErpcProxyServerConfig{
				HttpHostV4: "0.0.0.0",
				HttpPortV4: 4000,
				MaxTimeout: "30s",
			},
			Database: ErpcProxyDatabaseConfig{
				Type:          "postgresql",
				ConnectionUrl: "postgres://user:pass@localhost/db",
			},
			Projects: []ErpcProxyProjectConfig{
				{
					Id:              "main",
					RateLimitBudget: "frontend-budget",
					Cors: &ErpcProxyCorsConfig{
						AllowedOrigins:   []string{"https://example.com", "https://*.example.com"},
						AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
						AllowedHeaders:   []string{"Content-Type", "Authorization"},
						ExposedHeaders:   []string{"X-Request-ID"},
						AllowCredentials: true,
						MaxAge:           3600,
					},
					Networks: []ErpcProxyNetworkConfig{
						{
							ChainId:      1,
							Architecture: "evm",
							Failover: ErpcProxyFailoverConfig{
								MaxRetries:    3,
								BackoffMs:     100,
								BackoffMaxMs:  5000,
								BackoffFactor: 2,
								Duration:      "30s",
							},
						},
					},
					Upstreams: []ErpcProxyUpstreamConfig{
						{
							Id:              "upstream1",
							Type:            "evm",
							Endpoint:        "https://eth.example.com",
							RateLimitBudget: "global-budget",
							MaxRetries:      2,
							Timeout:         "15s",
						},
					},
				},
			},
		}

		yamlStr, err := marshalErpcConfig(config)
		assert.NoError(t, err)

		// Verify the YAML structure matches the eRPC specification from the documentation
		// Based on https://docs.erpc.cloud/config/example

		// Top-level sections
		assert.Contains(t, yamlStr, "logLevel:")
		assert.Contains(t, yamlStr, "server:")
		assert.Contains(t, yamlStr, "database:")
		assert.Contains(t, yamlStr, "projects:")

		// Server configuration structure
		assert.Contains(t, yamlStr, "httpHostV4:")
		assert.Contains(t, yamlStr, "httpPortV4:")
		assert.Contains(t, yamlStr, "maxTimeout:")

		// Database configuration structure
		assert.Contains(t, yamlStr, "type:")
		assert.Contains(t, yamlStr, "connectionUrl:")

		// Projects array structure
		assert.Contains(t, yamlStr, "id:")
		assert.Contains(t, yamlStr, "rateLimitBudget:")

		// CORS configuration structure (matches the specification)
		assert.Contains(t, yamlStr, "cors:")
		assert.Contains(t, yamlStr, "allowedOrigins:")
		assert.Contains(t, yamlStr, "allowedMethods:")
		assert.Contains(t, yamlStr, "allowedHeaders:")
		assert.Contains(t, yamlStr, "exposedHeaders:")
		assert.Contains(t, yamlStr, "allowCredentials:")
		assert.Contains(t, yamlStr, "maxAge:")

		// Networks structure
		assert.Contains(t, yamlStr, "networks:")
		assert.Contains(t, yamlStr, "architecture:")
		assert.Contains(t, yamlStr, "evm:")
		assert.Contains(t, yamlStr, "chainId:")
		assert.Contains(t, yamlStr, "failover:")

		// Upstreams structure
		assert.Contains(t, yamlStr, "upstreams:")
		assert.Contains(t, yamlStr, "id:")
		assert.Contains(t, yamlStr, "type:")
		assert.Contains(t, yamlStr, "endpoint:")

		// Verify YAML is properly formatted (no syntax errors)
		var result map[string]interface{}
		err = yaml.Unmarshal([]byte(yamlStr), &result)
		assert.NoError(t, err, "Generated YAML should be valid")

		// Verify the structure can be parsed back correctly
		assert.Equal(t, "debug", result["logLevel"])
		assert.Contains(t, result, "server")
		assert.Contains(t, result, "database")
		assert.Contains(t, result, "projects")

		// Verify projects array structure
		projects, ok := result["projects"].([]interface{})
		assert.True(t, ok, "projects should be an array")
		assert.Len(t, projects, 1, "should have one project")

		project := projects[0].(map[string]interface{})
		assert.Equal(t, "main", project["id"])
		assert.Equal(t, "frontend-budget", project["rateLimitBudget"])

		// Verify CORS structure
		cors, ok := project["cors"].(map[string]interface{})
		assert.True(t, ok, "cors should be a map")
		assert.Equal(t, []interface{}{"https://example.com", "https://*.example.com"}, cors["allowedOrigins"])
		assert.Equal(t, []interface{}{"GET", "POST", "OPTIONS"}, cors["allowedMethods"])
		assert.Equal(t, []interface{}{"Content-Type", "Authorization"}, cors["allowedHeaders"])
		assert.Equal(t, []interface{}{"X-Request-ID"}, cors["exposedHeaders"])
		assert.Equal(t, true, cors["allowCredentials"])
		assert.Equal(t, 3600, cors["maxAge"])
	})
}
