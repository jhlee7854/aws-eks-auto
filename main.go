package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		slug := fmt.Sprintf("organization/aws-vpc/%v", ctx.Stack())
		stackRef, err := pulumi.NewStackReference(ctx, slug, nil)
		if err != nil {
			return err
		}

		subnetIds := stackRef.GetOutput(pulumi.String("subnetIds")).AsStringArrayOutput()

		clusterRole, err := createClusterRole(ctx, nil)
		if err != nil {
			return err
		}

		nodeRole, err := createNodeRole(ctx, nil)
		if err != nil {
			return err
		}

		cluster, err := createEksCluster(ctx, subnetIds, clusterRole.Arn, nodeRole.Arn, pulumi.DependsOn([]pulumi.Resource{clusterRole, nodeRole}))
		if err != nil {
			return err
		}

		issuerUrl := cluster.Identities.Index(pulumi.Int(0)).Oidcs().Index(pulumi.Int(0)).Issuer()
		oidcp, err := addIdentityProvider(ctx, issuerUrl, pulumi.DependsOn([]pulumi.Resource{cluster}))
		if err != nil {
			return err
		}

		kubeconfig := generateKubeconfig(cluster.Endpoint, cluster.CertificateAuthority.Data().Elem(), cluster.Name, pulumi.String(config.New(ctx, "aws").Require("profile")).ToStringOutput())

		k8sProvider, err := kubernetes.NewProvider(ctx, "k8sporivder", &kubernetes.ProviderArgs{
			Kubeconfig: kubeconfig,
		}, pulumi.DependsOn([]pulumi.Resource{cluster}))
		if err != nil {
			return err
		}

		_, err = NewDefaultStorageClass(ctx, pulumi.DependsOn([]pulumi.Resource{k8sProvider}), pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		ctx.Export("clusterName", cluster.Name)
		ctx.Export("kubeconfig", kubeconfig)
		ctx.Export("oidcpUrl", oidcp.Url)
		ctx.Export("oidcpId", oidcp.Arn)

		return nil
	})
}
