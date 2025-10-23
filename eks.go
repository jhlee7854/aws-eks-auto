package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createEksCluster(ctx *pulumi.Context, subnetIds pulumi.StringArrayOutput, clusterRoleArn pulumi.StringOutput, nodeRoleArn pulumi.StringOutput, opts ...pulumi.ResourceOption) (*eks.Cluster, error) {
	clusterName := fmt.Sprintf("%s-eks", ctx.Stack())
	return eks.NewCluster(ctx, clusterName, &eks.ClusterArgs{
		AccessConfig: &eks.ClusterAccessConfigArgs{
			AuthenticationMode:                      pulumi.String("API"),
			BootstrapClusterCreatorAdminPermissions: pulumi.Bool(true),
		},
		BootstrapSelfManagedAddons: pulumi.Bool(false),
		ComputeConfig: &eks.ClusterComputeConfigArgs{
			Enabled: pulumi.Bool(true),
			NodePools: pulumi.StringArray{
				// pulumi.String("general-purpose"),
				pulumi.String("system"),
			},
			NodeRoleArn: nodeRoleArn,
		},
		EnabledClusterLogTypes: pulumi.StringArray{
			pulumi.String("api"),
			pulumi.String("audit"),
			pulumi.String("authenticator"),
			pulumi.String("controllerManager"),
			pulumi.String("scheduler"),
		},
		KubernetesNetworkConfig: &eks.ClusterKubernetesNetworkConfigArgs{
			ElasticLoadBalancing: &eks.ClusterKubernetesNetworkConfigElasticLoadBalancingArgs{
				Enabled: pulumi.Bool(true),
			},
			IpFamily:        pulumi.String("ipv4"),
			ServiceIpv4Cidr: pulumi.String("172.20.0.0/16"),
		},
		Name:    pulumi.String(clusterName),
		RoleArn: clusterRoleArn,
		StorageConfig: &eks.ClusterStorageConfigArgs{
			BlockStorage: &eks.ClusterStorageConfigBlockStorageArgs{
				Enabled: pulumi.Bool(true),
			},
		},
		UpgradePolicy: &eks.ClusterUpgradePolicyArgs{
			SupportType: pulumi.String("STANDARD"),
		},
		Version: pulumi.String("1.33"),
		VpcConfig: &eks.ClusterVpcConfigArgs{
			EndpointPrivateAccess: pulumi.Bool(true),
			PublicAccessCidrs: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
			SubnetIds: subnetIds,
		},
		ZonalShiftConfig: &eks.ClusterZonalShiftConfigArgs{
			Enabled: pulumi.Bool(true),
		},
	}, opts...)
}

// Create the KubeConfig Structure as per https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
func generateKubeconfig(clusterEndpoint pulumi.StringOutput, certData pulumi.StringOutput, clusterName pulumi.StringOutput, awsProfile pulumi.StringOutput) pulumi.StringOutput {
	return pulumi.Sprintf(`{
        "apiVersion": "v1",
        "clusters": [{
            "cluster": {
                "server": "%s",
                "certificate-authority-data": "%s"
            },
            "name": "kubernetes",
        }],
        "contexts": [{
            "context": {
                "cluster": "kubernetes",
                "user": "aws",
            },
            "name": "aws",
        }],
        "current-context": "aws",
        "kind": "Config",
        "users": [{
            "name": "aws",
            "user": {
                "exec": {
                    "apiVersion": "client.authentication.k8s.io/v1beta1",
                    "command": "aws",
                    "args": [
                        "eks",
                        "get-token",
                        "--cluster-name",
                        "%s",
                    ],
                    "env": [
                        {
                            "name": "AWS_PROFILE",
                            "value": "%s",
                        }
                    ],
                },
            },
        }],
    }`, clusterEndpoint, certData, clusterName, awsProfile)
}
