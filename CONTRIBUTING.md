# Contributing

## Repository Overview

This is a Pulumi component library for deploying Signet blockchain infrastructure to Kubernetes. It provides reusable infrastructure-as-code components for various blockchain services.

## Development Commands

```bash
# Build all packages
go build -v ./...

# Run all tests
go test -v ./...

# Run tests for a specific package (example)
go test -v ./pkg/builder

# Install/update dependencies
go mod download
go mod tidy

# Verify module dependencies
go mod verify
```

## Architecture

The repository follows a component-based architecture where each blockchain service is implemented as a self-contained Pulumi component under `pkg/`:

- **AWS Components** (`pkg/aws/`): IAM roles and Postgres database components
- **Blockchain Services**:
  - `pkg/builder/`: Builder service for blockchain operations
  - `pkg/erpc-proxy/`: eRPC proxy for RPC load balancing and failover
  - `pkg/ethereum/`: Composite Ethereum node management
    - `pkg/ethereum/consensus/`: Ethereum consensus layer (beacon chain)
    - `pkg/ethereum/execution/`: Ethereum execution layer
  - `pkg/pylon/`: Pylon service component
  - `pkg/quincey/`: Quincey service component
  - `pkg/signet_node/`: Core Signet node component
  - `pkg/txcache/`: Transaction cache component
- **Utilities** (`pkg/utils/`): Shared helper functions

### Component Structure Pattern

Each component follows this consistent structure:

#### Core Files (Required)
1. **`types.go`** - Type definitions using dual-struct pattern:
   - Public structs with base Go types for external API
   - Internal structs with Pulumi types for resource creation
   - Conversion functions (`toInternal()`) between public and internal types
2. **`[component].go`** - Main component implementation:
   - `New[Component]()` constructor function
   - Resource creation logic (ServiceAccount, ConfigMap, Deployment, Service)
   - Helper methods for accessing component outputs
3. **`validation.go`** - Input validation logic:
   - `Validate()` methods for all argument structs
   - Field-level validation functions
   - Resource format validators (memory, CPU, etc.)
4. **`validation_test.go`** - Comprehensive unit tests:
   - Table-driven tests for all validation functions
   - Coverage of valid and invalid cases
   - Proper error message assertions

#### Optional Files
5. **`constants.go`** - Component-specific constants:
   - Component kind identifier
   - Resource name suffixes
   - Default values for optional fields
   - Environment variable names
6. **`helpers.go`** - Utility functions specific to the component:
   - Complex data transformations
   - Custom marshaling logic
   - Component-specific business logic

### Key Design Patterns

#### 1. Dual-Struct Pattern
All components use a dual-struct pattern to separate public API from internal implementation:
```go
// Public struct with base Go types
type ComponentArgs struct {
    Namespace string
    Name      string
    Config    ConfigStruct
}

// Internal struct with Pulumi types
type componentArgsInternal struct {
    Namespace pulumi.StringInput
    Name      pulumi.StringInput
    Config    configStructInternal
}

// Conversion function
func (args ComponentArgs) toInternal() componentArgsInternal {
    return componentArgsInternal{
        Namespace: pulumi.String(args.Namespace),
        Name:      pulumi.String(args.Name),
        Config:    args.Config.toInternal(),
    }
}
```

#### 2. Component Registration Pattern
```go
func NewComponent(ctx *pulumi.Context, args ComponentArgs, opts ...pulumi.ResourceOption) (*Component, error) {
    // 1. Apply defaults to optional fields
    // 2. Validate arguments
    if err := args.Validate(); err != nil {
        return nil, fmt.Errorf("invalid component args: %w", err)
    }
    // 3. Convert to internal types
    internalArgs := args.toInternal()
    // 4. Register component
    component := &Component{}
    err := ctx.RegisterComponentResource(ComponentKind, args.Name, component)
    // 5. Create resources with pulumi.Parent(component)
}
```

#### 3. Resource Creation Pattern
- Always create ServiceAccount first
- Use consistent naming with suffixes (`-sa`, `-config`, `-deployment`, `-service`)
- Apply standard labels using `utils.CreateResourceLabels()`
- Set proper parent relationships with `pulumi.Parent(component)`

#### 4. Validation Pattern
```go
func (args *ComponentArgs) Validate() error {
    // Required field validation
    if args.Namespace == "" {
        return fmt.Errorf("namespace is required")
    }
    // Nested validation
    if err := args.Config.Validate(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    // Range validation
    if args.Port < 0 || args.Port > 65535 {
        return fmt.Errorf("invalid port: %d", args.Port)
    }
    return nil
}
```

#### 5. Complex Component Composition
Some components (like `ethereum` and `pylon`) compose other components:
- The `ethereum` package creates both execution and consensus clients
- The `pylon` package wraps an ethereum node with additional configuration
- Use clear ownership and dependency chains between sub-components

### Best Practices for New Components

1. **Start with Types**: Define your public API in `types.go` first
2. **Implement Validation Early**: Write validation logic and tests before resource creation
3. **Use Constants**: Define all magic numbers and strings in `constants.go`
4. **Follow Naming Conventions**: 
   - Public types/functions: PascalCase
   - Internal types: camelCase with "Internal" suffix
   - File names: lowercase with underscores
5. **Document Complex Logic**: Add comments for non-obvious transformations
6. **Test Thoroughly**: Aim for 100% coverage of validation logic
7. **Handle Errors Gracefully**: Wrap errors with context using `fmt.Errorf`
8. **Consider Defaults**: Provide sensible defaults for optional fields
9. **Support Configuration**: Use ConfigMaps for configuration, Secrets for sensitive data
10. **Implement Health Checks**: Add liveness, readiness, and startup probes where applicable

## Testing Approach

- Unit tests focus on validation logic and helper functions
- Test files use the `testify/assert` library for assertions
- No integration tests - components are tested when used in actual Pulumi programs

## Adding New Components

When creating a new component:
1. Create a new directory under `pkg/`
2. Follow the standard file structure (types.go, component.go, validation.go, etc.)
3. Implement the `New[Component]` function that returns a `ComponentResource`
4. Add comprehensive validation for all inputs
5. Write unit tests for validation logic
6. Update imports in consuming Pulumi programs to use the new component

## Common Pitfalls

1. **Pulumi Context**: Always use the provided `*pulumi.Context` for resource creation
2. **Resource Dependencies**: Use Pulumi's dependency system via `pulumi.DependsOn` when needed
3. **Kubernetes Namespaces**: Components don't create namespaces - they expect them to exist