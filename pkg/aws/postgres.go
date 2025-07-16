package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewPostgresDbComponent(ctx *pulumi.Context, args *PostgresDbArgs, opts ...pulumi.ResourceOption) (*PostgresDbComponent, error) {
	if err := validateDb(*args); err != nil {
		return nil, err
	}

	internalArgs := args.toInternal()
	component := &PostgresDbComponent{}

	err := ctx.RegisterComponentResource("signet:index:PostgresDb", args.DbName, component, opts...)
	if err != nil {
		return nil, err
	}

	// Create a VPC subnet group for the database
	subnetGroup, err := rds.NewSubnetGroup(ctx, "dbClusterSubnetGroup", &rds.SubnetGroupArgs{
		SubnetIds: internalArgs.DbSubnetGroupIds,
		Name:      pulumi.Sprintf("%s-subnet-group", args.DbName),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	dbCluster, err := rds.NewCluster(ctx, "dbCluster", &rds.ClusterArgs{
		Engine:         pulumi.String(rds.EngineTypeAuroraPostgresql),
		EngineVersion:  pulumi.String("16.4"),
		DatabaseName:   internalArgs.DbName,
		MasterUsername: internalArgs.DbUsername,
		MasterPassword: internalArgs.DbPassword,
		Serverlessv2ScalingConfiguration: &rds.ClusterServerlessv2ScalingConfigurationArgs{
			MaxCapacity: pulumi.Float64(2),
			MinCapacity: pulumi.Float64(1),
		},
		DbSubnetGroupName: subnetGroup.Name,
		SkipFinalSnapshot: pulumi.Bool(true),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	dbClusterInstance, err := rds.NewClusterInstance(ctx, "dbClusterInstance", &rds.ClusterInstanceArgs{
		ClusterIdentifier: dbCluster.ID(),
		InstanceClass:     pulumi.String("db.serverless"),
		Engine:            dbCluster.Engine,
		EngineVersion:     dbCluster.EngineVersion,
		DbSubnetGroupName: subnetGroup.Name,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	component.DbCluster = dbCluster
	component.DbClusterInstance = dbClusterInstance
	component.DbClusterEndpoint = dbCluster.Endpoint

	return component, nil
}
