package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createClusterRole(ctx *pulumi.Context, opts ...pulumi.ResourceOption) (*iam.Role, error) {
	roleName := fmt.Sprintf("%s-eks-cluster-role", ctx.Stack())
	return iam.NewRole(ctx, roleName, &iam.RoleArgs{
		AssumeRolePolicy: pulumi.Any("{\"Statement\":[{\"Action\":[\"sts:AssumeRole\",\"sts:TagSession\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"eks.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"),
		Description:      pulumi.String("Allows access to other AWS service resources that are required to operate Auto Mode clusters managed by EKS."),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy"),
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSComputePolicy"),
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy"),
			pulumi.String("arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy"),
		},
		Name: pulumi.Sprintf("%sAmazonEKSAutoClusterRole", ctx.Stack()),
	}, opts...)
}
