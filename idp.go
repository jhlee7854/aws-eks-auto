package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func addIdentityProvider(ctx *pulumi.Context, issuerUrl pulumi.StringPtrOutput, opts ...pulumi.ResourceOption) (*iam.OpenIdConnectProvider, error) {
	idpName := fmt.Sprintf("%s-eks-idp", ctx.Stack())
	return iam.NewOpenIdConnectProvider(ctx, idpName, &iam.OpenIdConnectProviderArgs{
		Url: issuerUrl.Elem(),
		ClientIdLists: pulumi.StringArray{
			pulumi.String("sts.amazonaws.com"),
		},
	}, opts...)
}
