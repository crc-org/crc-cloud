package util

import (
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const passwordOverrideSpecial = "!#%&*()-_=+[]{}<>.?"

func CreatePassword(ctx *pulumi.Context, name string) (*random.RandomPassword, error) {
	return CreateDependingPassword(ctx, name, []pulumi.Resource{})
}

func CreateDependingPassword(ctx *pulumi.Context, name string, dependecies []pulumi.Resource) (*random.RandomPassword, error) {
	return random.NewRandomPassword(ctx,
		name,
		&random.RandomPasswordArgs{
			Length:          pulumi.Int(16),
			Special:         pulumi.Bool(true),
			OverrideSpecial: pulumi.String(passwordOverrideSpecial),
		},
		pulumi.DependsOn(dependecies))
}
