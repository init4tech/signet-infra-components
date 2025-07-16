package txcache

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type TxCacheComponent struct {
	pulumi.ResourceState
}

type TxCacheComponentArgs struct {
	Namespace    string     `pulumi:"txCacheNamespace" validate:"required"`
	Name         string     `pulumi:"txCacheName" validate:"required"`
	Image        string     `pulumi:"txCacheImage" validate:"required"`
	Port         int        `pulumi:"txCachePort" validate:"required"`
	OauthIssuer  string     `pulumi:"txCacheOauthIssuer" validate:"required"`
	OauthJwksUri string     `pulumi:"txCacheOauthJwksUri" validate:"required"`
	Env          TxCacheEnv `pulumi:"txCacheEnv" validate:"required"`
}

type TxCacheComponentArgsInternal struct {
	Namespace    pulumi.StringInput `pulumi:"txCacheNamespace" validate:"required"`
	Name         pulumi.StringInput `pulumi:"txCacheName" validate:"required"`
	Image        pulumi.StringInput `pulumi:"txCacheImage" validate:"required"`
	Port         pulumi.IntInput    `pulumi:"txCachePort" validate:"required"`
	OauthIssuer  pulumi.StringInput `pulumi:"txCacheOauthIssuer" validate:"required"`
	OauthJwksUri pulumi.StringInput `pulumi:"txCacheOauthJwksUri" validate:"required"`
	Env          TxCacheEnvInternal `pulumi:"txCacheEnv" validate:"required"`
}

type TxCacheEnv struct {
	HttpPort                  string `pulumi:"txCacheHttpPort" validate:"required"`
	AwsAccessKeyId            string `pulumi:"txCacheAwsAccessKeyId" validate:"required"`
	AwsSecretAccessKey        string `pulumi:"txCacheAwsSecretAccessKey" validate:"required"`
	AwsRegion                 string `pulumi:"txCacheAwsRegion" validate:"required"`
	RustLog                   string `pulumi:"txCacheRustLog" validate:"required"`
	BlockQueryStart           string `pulumi:"txCacheBlockQueryStart" validate:"required"`
	BlockQueryCutoff          string `pulumi:"txCacheBlockQueryCutoff" validate:"required"`
	SlotOffset                string `pulumi:"txCacheSlotOffset" validate:"required"`
	ExpirationTimestampOffset string `pulumi:"txCacheExpirationTimestampOffset" validate:"required"`
	NetworkName               string `pulumi:"txCacheNetworkName" validate:"required"`
	Builders                  string `pulumi:"txCacheBuilders" validate:"required"`
	SlotDuration              string `pulumi:"txCacheSlotDuration" validate:"required"`
	StartTimestamp            string `pulumi:"txCacheStartTimestamp" validate:"required"`
	OtelExporterOtlpProtocol  string `pulumi:"txCacheOtelExporterOtlpProtocol"` // optional
	OtelExporterOtlpEndpoint  string `pulumi:"txCacheOtelExporterOtlpEndpoint"` // optional
}

type TxCacheEnvInternal struct {
	HttpPort                  pulumi.StringInput `pulumi:"txCacheHttpPort" validate:"required"`
	AwsAccessKeyId            pulumi.StringInput `pulumi:"txCacheAwsAccessKeyId" validate:"required"`
	AwsSecretAccessKey        pulumi.StringInput `pulumi:"txCacheAwsSecretAccessKey" validate:"required"`
	AwsRegion                 pulumi.StringInput `pulumi:"txCacheAwsRegion" validate:"required"`
	RustLog                   pulumi.StringInput `pulumi:"txCacheRustLog" validate:"required"`
	BlockQueryStart           pulumi.StringInput `pulumi:"txCacheBlockQueryStart" validate:"required"`
	BlockQueryCutoff          pulumi.StringInput `pulumi:"txCacheBlockQueryCutoff" validate:"required"`
	SlotOffset                pulumi.StringInput `pulumi:"txCacheSlotOffset" validate:"required"`
	ExpirationTimestampOffset pulumi.StringInput `pulumi:"txCacheExpirationTimestampOffset" validate:"required"`
	NetworkName               pulumi.StringInput `pulumi:"txCacheNetworkName" validate:"required"`
	Builders                  pulumi.StringInput `pulumi:"txCacheBuilders" validate:"required"`
	SlotDuration              pulumi.StringInput `pulumi:"txCacheSlotDuration" validate:"required"`
	StartTimestamp            pulumi.StringInput `pulumi:"txCacheStartTimestamp" validate:"required"`
	OtelExporterOtlpProtocol  pulumi.StringInput `pulumi:"txCacheOtelExporterOtlpProtocol"`
	OtelExporterOtlpEndpoint  pulumi.StringInput `pulumi:"txCacheOtelExporterOtlpEndpoint"`
}

func (args TxCacheComponentArgs) toInternal() TxCacheComponentArgsInternal {
	return TxCacheComponentArgsInternal{
		Namespace:    pulumi.String(args.Namespace),
		Name:         pulumi.String(args.Name),
		Image:        pulumi.String(args.Image),
		Port:         pulumi.Int(args.Port),
		OauthIssuer:  pulumi.String(args.OauthIssuer),
		OauthJwksUri: pulumi.String(args.OauthJwksUri),
		Env:          args.Env.toInternal(),
	}
}

func (env TxCacheEnv) toInternal() TxCacheEnvInternal {
	return TxCacheEnvInternal{
		HttpPort:                  pulumi.String(env.HttpPort),
		AwsAccessKeyId:            pulumi.String(env.AwsAccessKeyId),
		AwsSecretAccessKey:        pulumi.String(env.AwsSecretAccessKey),
		AwsRegion:                 pulumi.String(env.AwsRegion),
		RustLog:                   pulumi.String(env.RustLog),
		BlockQueryStart:           pulumi.String(env.BlockQueryStart),
		BlockQueryCutoff:          pulumi.String(env.BlockQueryCutoff),
		SlotOffset:                pulumi.String(env.SlotOffset),
		ExpirationTimestampOffset: pulumi.String(env.ExpirationTimestampOffset),
		NetworkName:               pulumi.String(env.NetworkName),
		Builders:                  pulumi.String(env.Builders),
		SlotDuration:              pulumi.String(env.SlotDuration),
		StartTimestamp:            pulumi.String(env.StartTimestamp),
		OtelExporterOtlpProtocol:  pulumi.String(env.OtelExporterOtlpProtocol),
		OtelExporterOtlpEndpoint:  pulumi.String(env.OtelExporterOtlpEndpoint),
	}
}
