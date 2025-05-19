# Signet Infrastructure Components

A collection of Pulumi infrastructure components for Signet blockchain services deployed to Kubernetes.

## Overview

This repository contains reusable Pulumi components for deploying and managing Signet blockchain infrastructure. Each component is packaged as a separate Go module that can be imported into your Pulumi projects.

## Components

### Builder

The `builder` component deploys a Signet builder service to Kubernetes. It creates the necessary resources:

- Kubernetes Deployment
- Kubernetes Service
- Service Account
- IAM Role and Policies
- Prometheus monitoring

#### Usage

```go
import (
    "github.com/your-org/signet-infra-components/builder"
)

func main() {
    // Create a new builder component
    builderComponent, err := builder.NewBuilder(ctx, builder.BuilderComponentArgs{
        Namespace: "signet",
        Name:      "signet-builder",
        Image:     "your-registry/builder:latest",
        AppLabels: builder.AppLabels{
            Labels: pulumi.StringMap{
                "app": pulumi.String("builder"),
            },
        },
        BuilderEnv: builder.BuilderEnv{
            BuilderPort:    pulumi.Int(8080),
            BuilderKey:     pulumi.String("arn:aws:kms:region:account:key/keyid"),
            HostRpcUrl:     pulumi.String("https://ethereum-rpc.example.com"),
            RollupRpcUrl:   pulumi.String("https://rollup-rpc.example.com"),
            // Add other required environment variables
        },
    })
    if err != nil {
        // Handle error
    }

    // Export the service URL
    ctx.Export("builderServiceUrl", builderComponent.GetServiceURL())
}
```

## Adding New Components

To add a new component:

1. Create a new directory with the component name
2. Implement the component following the same structure as the `builder` component:
   - `types.go` - Define component types and interfaces
   - `[component].go` - Implement the main component logic
   - `validation.go` - Implement input validation
   - `helpers.go` - Add helper functions

## Development

### Prerequisites

- Go 1.20+
- Pulumi CLI
- Access to a Kubernetes cluster

### Testing Components

You can test components by creating a simple Pulumi program that uses them:

```go
package main

import (
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/your-org/signet-infra-components/builder"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Test your component here
        return nil
    })
}
```

## License

MIT License see [LICENSE](LICENSE)