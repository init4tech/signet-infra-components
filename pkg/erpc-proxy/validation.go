package erpcproxy

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Validate validates the ErpcProxyComponentArgs
func (args *ErpcProxyComponentArgs) Validate() error {
	if args.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	if args.Name == "" {
		return fmt.Errorf("name is required")
	}

	if args.Image == "" {
		return fmt.Errorf("image is required")
	}

	// Validate config
	if err := args.Config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Validate resources if provided
	if err := args.Resources.Validate(); err != nil {
		return fmt.Errorf("invalid resources: %w", err)
	}

	// Validate replicas if provided
	if args.Replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}

	return nil
}

// Validate validates the ErpcProxyConfig
func (c *ErpcProxyConfig) Validate() error {
	// Validate log level if provided
	if c.LogLevel != "" {
		validLogLevels := []string{"debug", "info", "warn", "error"}
		valid := false
		for _, level := range validLogLevels {
			if strings.ToLower(c.LogLevel) == level {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid log level: %s, must be one of: %v", c.LogLevel, validLogLevels)
		}
	}

	// Validate database config if provided
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("invalid database config: %w", err)
	}

	// Validate server config
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("invalid server config: %w", err)
	}

	// Validate projects
	if len(c.Projects) == 0 {
		return fmt.Errorf("at least one project is required")
	}

	for i, project := range c.Projects {
		if err := project.Validate(); err != nil {
			return fmt.Errorf("invalid project at index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates the ErpcProxyDatabaseConfig
func (d *ErpcProxyDatabaseConfig) Validate() error {
	if d.Type != "" {
		validTypes := []string{"postgres", "memory", "redis"}
		valid := false
		for _, t := range validTypes {
			if d.Type == t {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid database type: %s, must be one of: %v", d.Type, validTypes)
		}

		// If type is postgres or redis, connection URL is required
		if (d.Type == "postgres" || d.Type == "redis") && d.ConnectionUrl == "" {
			return fmt.Errorf("connection URL is required for database type %s", d.Type)
		}
	}

	return nil
}

// Validate validates the ErpcProxyServerConfig
func (s *ErpcProxyServerConfig) Validate() error {
	if s.HttpPortV4 < 0 || s.HttpPortV4 > 65535 {
		return fmt.Errorf("invalid HTTP port: %d, must be between 0 and 65535", s.HttpPortV4)
	}

	if s.MaxTimeout != "" {
		_, err := time.ParseDuration(s.MaxTimeout)
		if err != nil {
			return fmt.Errorf("invalid max timeout: %w", err)
		}
	}

	return nil
}

// Validate validates the ErpcProxyProjectConfig
func (p *ErpcProxyProjectConfig) Validate() error {
	if p.Id == "" {
		return fmt.Errorf("project ID is required")
	}

	if len(p.Networks) == 0 {
		return fmt.Errorf("at least one network is required")
	}

	for i, network := range p.Networks {
		if err := network.Validate(); err != nil {
			return fmt.Errorf("invalid network at index %d: %w", i, err)
		}
	}

	// Validate upstreams
	if len(p.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream is required")
	}

	for i, upstream := range p.Upstreams {
		if err := upstream.Validate(); err != nil {
			return fmt.Errorf("invalid upstream at index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates the ErpcProxyNetworkConfig
func (n *ErpcProxyNetworkConfig) Validate() error {
	if n.ChainId <= 0 {
		return fmt.Errorf("chain ID must be positive")
	}

	if n.Architecture == "" {
		return fmt.Errorf("architecture is required")
	}

	// Validate architecture
	validArchitectures := []string{"evm", "non-evm"}
	valid := false
	for _, arch := range validArchitectures {
		if n.Architecture == arch {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid architecture: %s, must be one of: %v", n.Architecture, validArchitectures)
	}

	// Validate failover config if provided
	if err := n.Failover.Validate(); err != nil {
		return fmt.Errorf("invalid failover config: %w", err)
	}

	return nil
}

// Validate validates the ErpcProxyFailoverConfig
func (f *ErpcProxyFailoverConfig) Validate() error {
	if f.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}

	if f.BackoffMs < 0 {
		return fmt.Errorf("backoff must be non-negative")
	}

	if f.BackoffMaxMs < 0 {
		return fmt.Errorf("max backoff must be non-negative")
	}

	if f.BackoffMaxMs > 0 && f.BackoffMs > 0 && f.BackoffMaxMs < f.BackoffMs {
		return fmt.Errorf("max backoff must be greater than or equal to backoff")
	}

	if f.BackoffFactor < 0 {
		return fmt.Errorf("backoff factor must be non-negative")
	}

	return nil
}

// Validate validates the ErpcProxyUpstreamConfig
func (u *ErpcProxyUpstreamConfig) Validate() error {
	if u.Id == "" {
		return fmt.Errorf("upstream ID is required")
	}

	if u.Type == "" {
		return fmt.Errorf("upstream type is required")
	}

	// Validate upstream type
	validTypes := []string{"evm", "non-evm"}
	valid := false
	for _, t := range validTypes {
		if u.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid upstream type: %s, must be one of: %v", u.Type, validTypes)
	}

	if u.Endpoint == "" {
		return fmt.Errorf("upstream endpoint is required")
	}

	if u.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}

	return nil
}

// Validate validates the ErpcProxyResources
func (r *ErpcProxyResources) Validate() error {
	// Basic validation for resource strings
	// In a production environment, you might want to parse and validate the actual values
	// For now, we just ensure they're in the correct format if provided

	resourceFields := map[string]string{
		"memoryRequest": r.MemoryRequest,
		"memoryLimit":   r.MemoryLimit,
		"cpuRequest":    r.CpuRequest,
		"cpuLimit":      r.CpuLimit,
	}

	for name, value := range resourceFields {
		if value != "" && !isValidResourceString(value) {
			return fmt.Errorf("invalid %s format: %s", name, value)
		}
	}

	return nil
}

// isValidResourceString checks if a resource string is in valid format
func isValidResourceString(s string) bool {
	// Basic validation - should end with a unit like Mi, Gi, m, etc.
	// This is a simplified check
	if s == "" {
		return true // Empty is valid (uses defaults)
	}

	// Check for valid memory suffixes (Ki, Mi, Gi, Ti, Pi, Ei)
	memorySuffixes := []string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	for _, suffix := range memorySuffixes {
		if strings.HasSuffix(s, suffix) {
			prefix := strings.TrimSuffix(s, suffix)
			if prefix == "" {
				return false
			}
			num, err := strconv.ParseFloat(prefix, 64)
			return err == nil && num >= 0
		}
	}

	// Check for CPU suffixes (m for millicores, or k, M, G, T, P, E for decimal)
	cpuSuffixes := []string{"m", "k", "M", "G", "T", "P", "E"}
	for _, suffix := range cpuSuffixes {
		if strings.HasSuffix(s, suffix) {
			prefix := strings.TrimSuffix(s, suffix)
			if prefix == "" {
				return false
			}
			num, err := strconv.ParseFloat(prefix, 64)
			return err == nil && num >= 0
		}
	}

	// Also allow plain numbers (for CPU cores)
	num, err := strconv.ParseFloat(s, 64)
	return err == nil && num >= 0
}
