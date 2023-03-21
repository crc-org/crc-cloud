// Code generated by the Pulumi SDK Generator DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package resources

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Deployment information.
// API Version: 2021-01-01.
func LookupDeploymentAtScope(ctx *pulumi.Context, args *LookupDeploymentAtScopeArgs, opts ...pulumi.InvokeOption) (*LookupDeploymentAtScopeResult, error) {
	var rv LookupDeploymentAtScopeResult
	err := ctx.Invoke("azure-native:resources:getDeploymentAtScope", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

type LookupDeploymentAtScopeArgs struct {
	// The name of the deployment.
	DeploymentName string `pulumi:"deploymentName"`
	// The resource scope.
	Scope string `pulumi:"scope"`
}

// Deployment information.
type LookupDeploymentAtScopeResult struct {
	// The ID of the deployment.
	Id string `pulumi:"id"`
	// the location of the deployment.
	Location *string `pulumi:"location"`
	// The name of the deployment.
	Name string `pulumi:"name"`
	// Deployment properties.
	Properties DeploymentPropertiesExtendedResponse `pulumi:"properties"`
	// Deployment tags
	Tags map[string]string `pulumi:"tags"`
	// The type of the deployment.
	Type string `pulumi:"type"`
}

func LookupDeploymentAtScopeOutput(ctx *pulumi.Context, args LookupDeploymentAtScopeOutputArgs, opts ...pulumi.InvokeOption) LookupDeploymentAtScopeResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupDeploymentAtScopeResult, error) {
			args := v.(LookupDeploymentAtScopeArgs)
			r, err := LookupDeploymentAtScope(ctx, &args, opts...)
			var s LookupDeploymentAtScopeResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupDeploymentAtScopeResultOutput)
}

type LookupDeploymentAtScopeOutputArgs struct {
	// The name of the deployment.
	DeploymentName pulumi.StringInput `pulumi:"deploymentName"`
	// The resource scope.
	Scope pulumi.StringInput `pulumi:"scope"`
}

func (LookupDeploymentAtScopeOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupDeploymentAtScopeArgs)(nil)).Elem()
}

// Deployment information.
type LookupDeploymentAtScopeResultOutput struct{ *pulumi.OutputState }

func (LookupDeploymentAtScopeResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupDeploymentAtScopeResult)(nil)).Elem()
}

func (o LookupDeploymentAtScopeResultOutput) ToLookupDeploymentAtScopeResultOutput() LookupDeploymentAtScopeResultOutput {
	return o
}

func (o LookupDeploymentAtScopeResultOutput) ToLookupDeploymentAtScopeResultOutputWithContext(ctx context.Context) LookupDeploymentAtScopeResultOutput {
	return o
}

// The ID of the deployment.
func (o LookupDeploymentAtScopeResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) string { return v.Id }).(pulumi.StringOutput)
}

// the location of the deployment.
func (o LookupDeploymentAtScopeResultOutput) Location() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) *string { return v.Location }).(pulumi.StringPtrOutput)
}

// The name of the deployment.
func (o LookupDeploymentAtScopeResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) string { return v.Name }).(pulumi.StringOutput)
}

// Deployment properties.
func (o LookupDeploymentAtScopeResultOutput) Properties() DeploymentPropertiesExtendedResponseOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) DeploymentPropertiesExtendedResponse { return v.Properties }).(DeploymentPropertiesExtendedResponseOutput)
}

// Deployment tags
func (o LookupDeploymentAtScopeResultOutput) Tags() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) map[string]string { return v.Tags }).(pulumi.StringMapOutput)
}

// The type of the deployment.
func (o LookupDeploymentAtScopeResultOutput) Type() pulumi.StringOutput {
	return o.ApplyT(func(v LookupDeploymentAtScopeResult) string { return v.Type }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupDeploymentAtScopeResultOutput{})
}