package aws

import (
	"fmt"
)

// Validate validates the IAMStatement
func (s *IAMStatement) Validate() error {
	if s.Effect == "" {
		return fmt.Errorf("effect is required")
	}
	if len(s.Action) == 0 {
		return fmt.Errorf("action is required")
	}
	return nil
}

// Validate validates the IAMPolicy
func (p *IAMPolicy) Validate() error {
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	if len(p.Statement) == 0 {
		return fmt.Errorf("statement is required")
	}
	for i, stmt := range p.Statement {
		if err := stmt.Validate(); err != nil {
			return fmt.Errorf("statement %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates the KMSStatement
func (s *KMSStatement) Validate() error {
	if s.Effect == "" {
		return fmt.Errorf("effect is required")
	}
	if len(s.Action) == 0 {
		return fmt.Errorf("action is required")
	}
	if s.Resource == "" {
		return fmt.Errorf("resource is required")
	}
	return nil
}

// Validate validates the KMSPolicy
func (p *KMSPolicy) Validate() error {
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	if len(p.Statement) == 0 {
		return fmt.Errorf("statement is required")
	}
	for i, stmt := range p.Statement {
		if err := stmt.Validate(); err != nil {
			return fmt.Errorf("statement %d: %w", i, err)
		}
	}
	return nil
}

func validateDb(db PostgresDbArgs) error {
	if db.DbSubnetGroupName == "" {
		return fmt.Errorf("dbSubnetGroupName is required")
	}

	if db.DbUsername == "" {
		return fmt.Errorf("dbUsername is required")
	}

	if db.DbPassword == "" {
		return fmt.Errorf("dbPassword is required")
	}

	if db.DbName == "" {
		return fmt.Errorf("dbName is required")
	}

	return nil
}
