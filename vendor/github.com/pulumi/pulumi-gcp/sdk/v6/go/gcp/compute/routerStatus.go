// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get a Cloud Router's status within GCE from its name and region. This data source exposes the
// routes learned by a Cloud Router via BGP peers.
//
// For more information see [the official documentation](https://cloud.google.com/network-connectivity/docs/router/how-to/viewing-router-details)
// and
// [API](https://cloud.google.com/compute/docs/reference/rest/v1/routers/getRouterStatus).
//
// ## Example Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			_, err := compute.GetRouterStatus(ctx, &compute.GetRouterStatusArgs{
//				Name: "myrouter",
//			}, nil)
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// Deprecated: gcp.compute.RouterStatus has been deprecated in favor of gcp.compute.getRouterStatus
func RouterStatus(ctx *pulumi.Context, args *RouterStatusArgs, opts ...pulumi.InvokeOption) (*RouterStatusResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	var rv RouterStatusResult
	err := ctx.Invoke("gcp:compute/routerStatus:RouterStatus", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking RouterStatus.
type RouterStatusArgs struct {
	// The name of the router.
	Name string `pulumi:"name"`
	// The ID of the project in which the resource
	// belongs. If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
	// The region this router has been created in. If
	// unspecified, this defaults to the region configured in the provider.
	Region *string `pulumi:"region"`
}

// A collection of values returned by RouterStatus.
type RouterStatusResult struct {
	// List of best `compute#routes` configurations for this router's network. See compute.Route resource for available attributes.
	BestRoutes []RouterStatusBestRoute `pulumi:"bestRoutes"`
	// List of best `compute#routes` for this specific router. See compute.Route resource for available attributes.
	BestRoutesForRouters []RouterStatusBestRoutesForRouter `pulumi:"bestRoutesForRouters"`
	// The provider-assigned unique ID for this managed resource.
	Id   string `pulumi:"id"`
	Name string `pulumi:"name"`
	// The network name or resource link to the parent
	// network of this subnetwork.
	Network string  `pulumi:"network"`
	Project *string `pulumi:"project"`
	Region  string  `pulumi:"region"`
}

func RouterStatusOutput(ctx *pulumi.Context, args RouterStatusOutputArgs, opts ...pulumi.InvokeOption) RouterStatusResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (RouterStatusResult, error) {
			args := v.(RouterStatusArgs)
			r, err := RouterStatus(ctx, &args, opts...)
			var s RouterStatusResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(RouterStatusResultOutput)
}

// A collection of arguments for invoking RouterStatus.
type RouterStatusOutputArgs struct {
	// The name of the router.
	Name pulumi.StringInput `pulumi:"name"`
	// The ID of the project in which the resource
	// belongs. If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput `pulumi:"project"`
	// The region this router has been created in. If
	// unspecified, this defaults to the region configured in the provider.
	Region pulumi.StringPtrInput `pulumi:"region"`
}

func (RouterStatusOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*RouterStatusArgs)(nil)).Elem()
}

// A collection of values returned by RouterStatus.
type RouterStatusResultOutput struct{ *pulumi.OutputState }

func (RouterStatusResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*RouterStatusResult)(nil)).Elem()
}

func (o RouterStatusResultOutput) ToRouterStatusResultOutput() RouterStatusResultOutput {
	return o
}

func (o RouterStatusResultOutput) ToRouterStatusResultOutputWithContext(ctx context.Context) RouterStatusResultOutput {
	return o
}

// List of best `compute#routes` configurations for this router's network. See compute.Route resource for available attributes.
func (o RouterStatusResultOutput) BestRoutes() RouterStatusBestRouteArrayOutput {
	return o.ApplyT(func(v RouterStatusResult) []RouterStatusBestRoute { return v.BestRoutes }).(RouterStatusBestRouteArrayOutput)
}

// List of best `compute#routes` for this specific router. See compute.Route resource for available attributes.
func (o RouterStatusResultOutput) BestRoutesForRouters() RouterStatusBestRoutesForRouterArrayOutput {
	return o.ApplyT(func(v RouterStatusResult) []RouterStatusBestRoutesForRouter { return v.BestRoutesForRouters }).(RouterStatusBestRoutesForRouterArrayOutput)
}

// The provider-assigned unique ID for this managed resource.
func (o RouterStatusResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v RouterStatusResult) string { return v.Id }).(pulumi.StringOutput)
}

func (o RouterStatusResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v RouterStatusResult) string { return v.Name }).(pulumi.StringOutput)
}

// The network name or resource link to the parent
// network of this subnetwork.
func (o RouterStatusResultOutput) Network() pulumi.StringOutput {
	return o.ApplyT(func(v RouterStatusResult) string { return v.Network }).(pulumi.StringOutput)
}

func (o RouterStatusResultOutput) Project() pulumi.StringPtrOutput {
	return o.ApplyT(func(v RouterStatusResult) *string { return v.Project }).(pulumi.StringPtrOutput)
}

func (o RouterStatusResultOutput) Region() pulumi.StringOutput {
	return o.ApplyT(func(v RouterStatusResult) string { return v.Region }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(RouterStatusResultOutput{})
}