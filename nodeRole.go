package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNodeRole(ctx *pulumi.Context, opts ...pulumi.ResourceOption) (*iam.Role, error) {
	roleName := fmt.Sprintf("%s-eks-node-role", ctx.Stack())
	return iam.NewRole(ctx, roleName, &iam.RoleArgs{
		AssumeRolePolicy: pulumi.Any("{\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"),
		Description:      pulumi.String("Allows EKS nodes to connect to EKS Auto Mode clusters and to pull container images from ECR."),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly"),
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy"),
		},
		Name: pulumi.String("AmazonEKSAutoNodeRole"),
	}, opts...)
}
