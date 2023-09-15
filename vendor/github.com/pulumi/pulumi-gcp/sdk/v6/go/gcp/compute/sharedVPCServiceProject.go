// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Enables the Google Compute Engine
// [Shared VPC](https://cloud.google.com/compute/docs/shared-vpc)
// feature for a project, assigning it as a Shared VPC service project associated
// with a given host project.
//
// For more information, see,
// [the Project API documentation](https://cloud.google.com/compute/docs/reference/latest/projects),
// where the Shared VPC feature is referred to by its former name "XPN".
//
// > **Note:** If Shared VPC Admin role is set at the folder level, use the google-beta provider. The google provider only supports this permission at project or organizational level currently. [[0]](https://cloud.google.com/vpc/docs/provisioning-shared-vpc#enable-shared-vpc-host)
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
//			_, err := compute.NewSharedVPCServiceProject(ctx, "service1", &compute.SharedVPCServiceProjectArgs{
//				HostProject:    pulumi.String("host-project-id"),
//				ServiceProject: pulumi.String("service-project-id-1"),
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// For a complete Shared VPC example with both host and service projects, see
// [`compute.SharedVPCHostProject`](https://www.terraform.io/docs/providers/google/r/compute_shared_vpc_host_project.html).
//
// ## Import
//
// Google Compute Engine Shared VPC service project feature can be imported using the `host_project` and `service_project`, e.g.
//
// ```sh
//
//	$ pulumi import gcp:compute/sharedVPCServiceProject:SharedVPCServiceProject service1 host-project-id/service-project-id-1
//
// ```
type SharedVPCServiceProject struct {
	pulumi.CustomResourceState

	// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
	DeletionPolicy pulumi.StringPtrOutput `pulumi:"deletionPolicy"`
	// The ID of a host project to associate.
	HostProject pulumi.StringOutput `pulumi:"hostProject"`
	// The ID of the project that will serve as a Shared VPC service project.
	ServiceProject pulumi.StringOutput `pulumi:"serviceProject"`
}

// NewSharedVPCServiceProject registers a new resource with the given unique name, arguments, and options.
func NewSharedVPCServiceProject(ctx *pulumi.Context,
	name string, args *SharedVPCServiceProjectArgs, opts ...pulumi.ResourceOption) (*SharedVPCServiceProject, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.HostProject == nil {
		return nil, errors.New("invalid value for required argument 'HostProject'")
	}
	if args.ServiceProject == nil {
		return nil, errors.New("invalid value for required argument 'ServiceProject'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource SharedVPCServiceProject
	err := ctx.RegisterResource("gcp:compute/sharedVPCServiceProject:SharedVPCServiceProject", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetSharedVPCServiceProject gets an existing SharedVPCServiceProject resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetSharedVPCServiceProject(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *SharedVPCServiceProjectState, opts ...pulumi.ResourceOption) (*SharedVPCServiceProject, error) {
	var resource SharedVPCServiceProject
	err := ctx.ReadResource("gcp:compute/sharedVPCServiceProject:SharedVPCServiceProject", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering SharedVPCServiceProject resources.
type sharedVPCServiceProjectState struct {
	// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
	DeletionPolicy *string `pulumi:"deletionPolicy"`
	// The ID of a host project to associate.
	HostProject *string `pulumi:"hostProject"`
	// The ID of the project that will serve as a Shared VPC service project.
	ServiceProject *string `pulumi:"serviceProject"`
}

type SharedVPCServiceProjectState struct {
	// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
	DeletionPolicy pulumi.StringPtrInput
	// The ID of a host project to associate.
	HostProject pulumi.StringPtrInput
	// The ID of the project that will serve as a Shared VPC service project.
	ServiceProject pulumi.StringPtrInput
}

func (SharedVPCServiceProjectState) ElementType() reflect.Type {
	return reflect.TypeOf((*sharedVPCServiceProjectState)(nil)).Elem()
}

type sharedVPCServiceProjectArgs struct {
	// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
	DeletionPolicy *string `pulumi:"deletionPolicy"`
	// The ID of a host project to associate.
	HostProject string `pulumi:"hostProject"`
	// The ID of the project that will serve as a Shared VPC service project.
	ServiceProject string `pulumi:"serviceProject"`
}

// The set of arguments for constructing a SharedVPCServiceProject resource.
type SharedVPCServiceProjectArgs struct {
	// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
	DeletionPolicy pulumi.StringPtrInput
	// The ID of a host project to associate.
	HostProject pulumi.StringInput
	// The ID of the project that will serve as a Shared VPC service project.
	ServiceProject pulumi.StringInput
}

func (SharedVPCServiceProjectArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*sharedVPCServiceProjectArgs)(nil)).Elem()
}

type SharedVPCServiceProjectInput interface {
	pulumi.Input

	ToSharedVPCServiceProjectOutput() SharedVPCServiceProjectOutput
	ToSharedVPCServiceProjectOutputWithContext(ctx context.Context) SharedVPCServiceProjectOutput
}

func (*SharedVPCServiceProject) ElementType() reflect.Type {
	return reflect.TypeOf((**SharedVPCServiceProject)(nil)).Elem()
}

func (i *SharedVPCServiceProject) ToSharedVPCServiceProjectOutput() SharedVPCServiceProjectOutput {
	return i.ToSharedVPCServiceProjectOutputWithContext(context.Background())
}

func (i *SharedVPCServiceProject) ToSharedVPCServiceProjectOutputWithContext(ctx context.Context) SharedVPCServiceProjectOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SharedVPCServiceProjectOutput)
}

// SharedVPCServiceProjectArrayInput is an input type that accepts SharedVPCServiceProjectArray and SharedVPCServiceProjectArrayOutput values.
// You can construct a concrete instance of `SharedVPCServiceProjectArrayInput` via:
//
//	SharedVPCServiceProjectArray{ SharedVPCServiceProjectArgs{...} }
type SharedVPCServiceProjectArrayInput interface {
	pulumi.Input

	ToSharedVPCServiceProjectArrayOutput() SharedVPCServiceProjectArrayOutput
	ToSharedVPCServiceProjectArrayOutputWithContext(context.Context) SharedVPCServiceProjectArrayOutput
}

type SharedVPCServiceProjectArray []SharedVPCServiceProjectInput

func (SharedVPCServiceProjectArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SharedVPCServiceProject)(nil)).Elem()
}

func (i SharedVPCServiceProjectArray) ToSharedVPCServiceProjectArrayOutput() SharedVPCServiceProjectArrayOutput {
	return i.ToSharedVPCServiceProjectArrayOutputWithContext(context.Background())
}

func (i SharedVPCServiceProjectArray) ToSharedVPCServiceProjectArrayOutputWithContext(ctx context.Context) SharedVPCServiceProjectArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SharedVPCServiceProjectArrayOutput)
}

// SharedVPCServiceProjectMapInput is an input type that accepts SharedVPCServiceProjectMap and SharedVPCServiceProjectMapOutput values.
// You can construct a concrete instance of `SharedVPCServiceProjectMapInput` via:
//
//	SharedVPCServiceProjectMap{ "key": SharedVPCServiceProjectArgs{...} }
type SharedVPCServiceProjectMapInput interface {
	pulumi.Input

	ToSharedVPCServiceProjectMapOutput() SharedVPCServiceProjectMapOutput
	ToSharedVPCServiceProjectMapOutputWithContext(context.Context) SharedVPCServiceProjectMapOutput
}

type SharedVPCServiceProjectMap map[string]SharedVPCServiceProjectInput

func (SharedVPCServiceProjectMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SharedVPCServiceProject)(nil)).Elem()
}

func (i SharedVPCServiceProjectMap) ToSharedVPCServiceProjectMapOutput() SharedVPCServiceProjectMapOutput {
	return i.ToSharedVPCServiceProjectMapOutputWithContext(context.Background())
}

func (i SharedVPCServiceProjectMap) ToSharedVPCServiceProjectMapOutputWithContext(ctx context.Context) SharedVPCServiceProjectMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SharedVPCServiceProjectMapOutput)
}

type SharedVPCServiceProjectOutput struct{ *pulumi.OutputState }

func (SharedVPCServiceProjectOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**SharedVPCServiceProject)(nil)).Elem()
}

func (o SharedVPCServiceProjectOutput) ToSharedVPCServiceProjectOutput() SharedVPCServiceProjectOutput {
	return o
}

func (o SharedVPCServiceProjectOutput) ToSharedVPCServiceProjectOutputWithContext(ctx context.Context) SharedVPCServiceProjectOutput {
	return o
}

// The deletion policy for the shared VPC service. Setting ABANDON allows the resource to be abandoned rather than deleted. Possible values are: "ABANDON".
func (o SharedVPCServiceProjectOutput) DeletionPolicy() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *SharedVPCServiceProject) pulumi.StringPtrOutput { return v.DeletionPolicy }).(pulumi.StringPtrOutput)
}

// The ID of a host project to associate.
func (o SharedVPCServiceProjectOutput) HostProject() pulumi.StringOutput {
	return o.ApplyT(func(v *SharedVPCServiceProject) pulumi.StringOutput { return v.HostProject }).(pulumi.StringOutput)
}

// The ID of the project that will serve as a Shared VPC service project.
func (o SharedVPCServiceProjectOutput) ServiceProject() pulumi.StringOutput {
	return o.ApplyT(func(v *SharedVPCServiceProject) pulumi.StringOutput { return v.ServiceProject }).(pulumi.StringOutput)
}

type SharedVPCServiceProjectArrayOutput struct{ *pulumi.OutputState }

func (SharedVPCServiceProjectArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SharedVPCServiceProject)(nil)).Elem()
}

func (o SharedVPCServiceProjectArrayOutput) ToSharedVPCServiceProjectArrayOutput() SharedVPCServiceProjectArrayOutput {
	return o
}

func (o SharedVPCServiceProjectArrayOutput) ToSharedVPCServiceProjectArrayOutputWithContext(ctx context.Context) SharedVPCServiceProjectArrayOutput {
	return o
}

func (o SharedVPCServiceProjectArrayOutput) Index(i pulumi.IntInput) SharedVPCServiceProjectOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *SharedVPCServiceProject {
		return vs[0].([]*SharedVPCServiceProject)[vs[1].(int)]
	}).(SharedVPCServiceProjectOutput)
}

type SharedVPCServiceProjectMapOutput struct{ *pulumi.OutputState }

func (SharedVPCServiceProjectMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SharedVPCServiceProject)(nil)).Elem()
}

func (o SharedVPCServiceProjectMapOutput) ToSharedVPCServiceProjectMapOutput() SharedVPCServiceProjectMapOutput {
	return o
}

func (o SharedVPCServiceProjectMapOutput) ToSharedVPCServiceProjectMapOutputWithContext(ctx context.Context) SharedVPCServiceProjectMapOutput {
	return o
}

func (o SharedVPCServiceProjectMapOutput) MapIndex(k pulumi.StringInput) SharedVPCServiceProjectOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *SharedVPCServiceProject {
		return vs[0].(map[string]*SharedVPCServiceProject)[vs[1].(string)]
	}).(SharedVPCServiceProjectOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*SharedVPCServiceProjectInput)(nil)).Elem(), &SharedVPCServiceProject{})
	pulumi.RegisterInputType(reflect.TypeOf((*SharedVPCServiceProjectArrayInput)(nil)).Elem(), SharedVPCServiceProjectArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*SharedVPCServiceProjectMapInput)(nil)).Elem(), SharedVPCServiceProjectMap{})
	pulumi.RegisterOutputType(SharedVPCServiceProjectOutput{})
	pulumi.RegisterOutputType(SharedVPCServiceProjectArrayOutput{})
	pulumi.RegisterOutputType(SharedVPCServiceProjectMapOutput{})
}