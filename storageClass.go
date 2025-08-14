package main

import (
	meta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	storage "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/create-storage-class.html
func NewDefaultStorageClass(ctx *pulumi.Context, opts ...pulumi.ResourceOption) (*storage.StorageClass, error) {
	return storage.NewStorageClass(ctx, "ebs-sc", &storage.StorageClassArgs{
		Metadata: meta.ObjectMetaArgs{
			Annotations: pulumi.StringMap{
				"storageclass.kubernetes.io/is-default-class": pulumi.String("true"),
			},
		},
		Provisioner:       pulumi.String("ebs.csi.eks.amazonaws.com"),
		VolumeBindingMode: pulumi.String("WaitForFirstConsumer"),
		Parameters: pulumi.StringMap{
			"type": pulumi.String("gp3"),
		},
	}, opts...)
}
