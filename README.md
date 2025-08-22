# Signet Infrastructure Components

A comprehensive Pulumi component library for deploying and managing Signet blockchain infrastructure on Kubernetes.

## Overview

This repository provides production-ready, reusable Pulumi components for deploying blockchain infrastructure services. Built with Go and designed for Kubernetes, these components offer a consistent, validated approach to infrastructure as code for blockchain operations.

## Features

- **Type-Safe Infrastructure**: Strongly typed Go interfaces with comprehensive validation
- **AWS Integration**: Built-in support for IAM roles, S3 buckets, and RDS databases
- **Kubernetes Native**: Deploy to any Kubernetes cluster with standard resources
- **Component Composition**: Complex services built from reusable sub-components
- **Production Ready**: Health checks, resource limits, and monitoring built-in

## Components

### Core Blockchain Services

#### Builder (`pkg/builder/`)
Deploys a Signet builder service for blockchain operations.

**Resources Created:**
- Kubernetes Deployment with configurable replicas
- Service (ClusterIP) for internal access
- ServiceAccount with optional AWS IAM role
- ConfigMap for environment configuration
- Persistent Volume Claims for data storage

#### Signet Node (`pkg/signet_node/`)
Core Signet blockchain node implementation.

**Resources Created:**
- StatefulSet for node persistence
- Service for peer-to-peer and RPC access
- PersistentVolumeClaim for blockchain data
- ConfigMap for node configuration

#### Transaction Cache (`pkg/txcache/`)
High-performance transaction caching service.

**Resources Created:**
- Deployment with memory-optimized configuration
- Service for cache access
- ConfigMap for cache settings
- Optional Redis backend support

### Ethereum Infrastructure

#### Ethereum Node (`pkg/ethereum/`)
Composite component that manages both execution and consensus clients as a complete Ethereum node.

**Sub-components:**
- **Execution Client** (`pkg/ethereum/execution/`): Reth or compatible execution layer
- **Consensus Client** (`pkg/ethereum/consensus/`): Lighthouse Beacon chain consensus layer

**Features:**
- Automatic JWT secret management
- Inter-client communication setup

### Specialized Services

#### eRPC Proxy (`pkg/erpc-proxy/`)
Advanced RPC proxy with load balancing and failover capabilities.

**Resources Created:**
- Deployment with health probes
- Service for RPC access
- ConfigMap for complex routing configuration
- Secret for API keys
- ServiceAccount for pod identity

**Example:**
```go
		_, err = erpcproxy.NewErpcProxy(ctx, erpcproxy.ErpcProxyComponentArgs{
			Namespace: "default",
			Name:      "rpc-erpc-proxy",
			Config: erpcproxy.ErpcProxyConfig{
				Server: erpcproxy.ErpcProxyServerConfig{
					HttpHostV4: "0.0.0.0",
					HttpPortV4: 8545,
					MaxTimeout: "60s",
				},
				Projects: []erpcproxy.ErpcProxyProjectConfig{
					{
						Id: "rpc",
						Networks: []erpcproxy.ErpcProxyNetworkConfig{
							{
								ChainId:      1,
								Architecture: "evm",
							},
						},
						Upstreams: []erpcproxy.ErpcProxyUpstreamConfig{
							{
								Id:       "rpc",
								Type:     "evm",
								Endpoint: "http://some-rpc-service.default.svc.cluster.local:8545",
							},
						},
					},
				},
			},
			Resources: erpcproxy.ErpcProxyResources{
				MemoryRequest: "1Gi",
				MemoryLimit:   "2Gi",
				CpuRequest:    "100m",
				CpuLimit:      "200m",
			},
			Replicas: 1,
		})
		if err != nil {
			return err
		}
```

#### Pylon (`pkg/pylon/`)
Ethereum Blob cold storage client

**Features:**
- Deploys an ExEx on top of a Reth/Lighthouse pair
- S3 integration for blob storage
- PostgreSQL database support
- Custom environment configuration

#### Quincey (`pkg/quincey/`)
Quincey service component for specialized blockchain operations.

### AWS Integration (`pkg/aws/`)

#### IAM Roles
Create and manage AWS IAM roles for Kubernetes service accounts (IRSA).

#### PostgreSQL Database
Provision RDS PostgreSQL instances with:
- Automated backups
- Security group configuration
- Parameter group customization
- Connection string management

### Utilities (`pkg/utils/`)
Shared helper functions for:
- Resource labeling
- ConfigMap creation
- Port parsing with defaults
- Environment variable management

## Usage Example

```go
package main

import (
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/init4tech/signet-infra-components/pkg/builder"
    "github.com/init4tech/signet-infra-components/pkg/erpc-proxy"
    "github.com/init4tech/signet-infra-components/pkg/ethereum"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Deploy a builder service
        builder, err := builder.NewBuilder(ctx, builder.BuilderComponentArgs{
            Namespace: "signet",
            Name:      "signet-builder",
            Image:     "your-registry/builder:latest",
            BuilderEnv: builder.BuilderEnv{
                // Configuration
            },
        })
        if err != nil {
            return err
        }

        // Deploy an eRPC proxy
        proxy, err := erpcproxy.NewErpcProxy(ctx, erpcproxy.ErpcProxyComponentArgs{
            Namespace: "signet",
            Name:      "erpc-proxy",
            Image:     "ghcr.io/erpc/erpc:latest",
            Config: erpcproxy.ErpcProxyConfig{
                LogLevel: "info",
                Projects: []erpcproxy.ErpcProxyProjectConfig{
                    {
                        Id: "mainnet",
                        Networks: []erpcproxy.ErpcProxyNetworkConfig{
                            {
                                ChainId:      1,
                                Architecture: "evm",
                                Upstreams: []erpcproxy.ErpcProxyUpstreamConfig{
                                    {
                                        Id:       "primary",
                                        Type:     "http",
                                        Endpoint: "https://eth-mainnet.example.com",
                                    },
                                },
                            },
                        },
                    },
                },
            },
        })
        if err != nil {
            return err
        }

        // Deploy a complete Ethereum node
        ethNode, err := ethereum.NewEthereumNodeComponent(ctx, &ethereum.EthereumNodeArgs{
            Name:      "eth-node",
            Namespace: "signet",
            ExecutionClient: &execution.ExecutionClientArgs{
                // Execution configuration
            },
            ConsensusClient: &consensus.ConsensusClientArgs{
                // Consensus configuration
            },
        })
        if err != nil {
            return err
        }

        // Export service endpoints
        ctx.Export("builderUrl", builder.GetServiceURL())
        ctx.Export("erpcUrl", proxy.GetServiceURL())
        
        return nil
    })
}
```

## Development

### Prerequisites

- Go 1.22+ (CI uses 1.22, go.mod specifies 1.24.3)
- Pulumi CLI 3.x
- Access to a Kubernetes cluster
- AWS credentials (for AWS-integrated components)

### Project Structure

```
signet-infra-components/
├── pkg/
│   ├── aws/              # AWS resource components
│   ├── builder/          # Builder service
│   ├── erpc-proxy/       # eRPC proxy service
│   ├── ethereum/         # Ethereum node components
│   │   ├── consensus/    # Consensus client
│   │   └── execution/    # Execution client
│   ├── pylon/           # Pylon service
│   ├── quincey/         # Quincey service
│   ├── signet_node/     # Signet node
│   ├── txcache/         # Transaction cache
│   └── utils/           # Shared utilities
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md            # AI assistant guidelines
```

### Building and Testing

```bash
# Build all packages
go build -v ./...

# Run all tests
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/builder

# Verify dependencies
go mod verify

# Update dependencies
go mod tidy
```

### Adding New Components

1. **Create Package Directory**: `mkdir -p pkg/your-component`

2. **Implement Required Files**:
   - `types.go` - Public and internal type definitions
   - `your-component.go` - Main component implementation
   - `validation.go` - Input validation logic
   - `validation_test.go` - Unit tests
   - `constants.go` - Component constants (optional)
   - `helpers.go` - Helper functions (optional)

3. **Follow Design Patterns**:
   - Use dual-struct pattern (public/internal types)
   - Implement comprehensive validation
   - Create standard Kubernetes resources
   - Add health probes where applicable

4. **Test Thoroughly**:
   - Write table-driven tests
   - Cover validation edge cases
   - Aim for high test coverage

See [CLAUDE.md](CLAUDE.md) for detailed architectural patterns and best practices.

## Component Configuration

### Monitoring

Components expose Prometheus metrics on configurable ports, typically:
- Application metrics: `/metrics`
- Health endpoint: `/healthcheck` or `/health`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the component patterns described in CONTRIBUTING.md
4. Ensure all tests pass
5. Submit a pull request

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## License

MIT License - see [LICENSE](LICENSE) for details.