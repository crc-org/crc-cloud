package tls

import (
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateKey(ctx *pulumi.Context) (*tls.PrivateKey, error) {
	pk, err := tls.NewPrivateKey(
		ctx,
		"OpenshiftLocal-OCP",
		&tls.PrivateKeyArgs{
			Algorithm: pulumi.String("RSA"),
			RsaBits:   pulumi.Int(4096),
		})
	if err != nil {
		return nil, err
	}
	return pk, nil
}
