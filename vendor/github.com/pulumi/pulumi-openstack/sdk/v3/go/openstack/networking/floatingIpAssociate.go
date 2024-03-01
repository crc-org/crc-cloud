// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Associates a floating IP to a port. This is useful for situations
// where you have a pre-allocated floating IP or are unable to use the
// `networking.FloatingIp` resource to create a floating IP.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/networking"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			port1, err := networking.NewPort(ctx, "port1", &networking.PortArgs{
//				NetworkId: pulumi.String("a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"),
//			})
//			if err != nil {
//				return err
//			}
//			_, err = networking.NewFloatingIpAssociate(ctx, "fip1", &networking.FloatingIpAssociateArgs{
//				FloatingIp: pulumi.String("1.2.3.4"),
//				PortId:     port1.ID(),
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
// ## Import
//
// Floating IP associations can be imported using the `id` of the floating IP, e.g.
//
// ```sh
//
//	$ pulumi import openstack:networking/floatingIpAssociate:FloatingIpAssociate fip 2c7f39f3-702b-48d1-940c-b50384177ee1
//
// ```
type FloatingIpAssociate struct {
	pulumi.CustomResourceState

	FixedIp pulumi.StringOutput `pulumi:"fixedIp"`
	// IP Address of an existing floating IP.
	FloatingIp pulumi.StringOutput `pulumi:"floatingIp"`
	// ID of an existing port with at least one IP address to
	// associate with this floating IP.
	PortId pulumi.StringOutput `pulumi:"portId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a floating IP that can be used with
	// another networking resource, such as a load balancer. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// floating IP (which may or may not have a different address).
	Region pulumi.StringOutput `pulumi:"region"`
}

// NewFloatingIpAssociate registers a new resource with the given unique name, arguments, and options.
func NewFloatingIpAssociate(ctx *pulumi.Context,
	name string, args *FloatingIpAssociateArgs, opts ...pulumi.ResourceOption) (*FloatingIpAssociate, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.FloatingIp == nil {
		return nil, errors.New("invalid value for required argument 'FloatingIp'")
	}
	if args.PortId == nil {
		return nil, errors.New("invalid value for required argument 'PortId'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource FloatingIpAssociate
	err := ctx.RegisterResource("openstack:networking/floatingIpAssociate:FloatingIpAssociate", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetFloatingIpAssociate gets an existing FloatingIpAssociate resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetFloatingIpAssociate(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *FloatingIpAssociateState, opts ...pulumi.ResourceOption) (*FloatingIpAssociate, error) {
	var resource FloatingIpAssociate
	err := ctx.ReadResource("openstack:networking/floatingIpAssociate:FloatingIpAssociate", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering FloatingIpAssociate resources.
type floatingIpAssociateState struct {
	FixedIp *string `pulumi:"fixedIp"`
	// IP Address of an existing floating IP.
	FloatingIp *string `pulumi:"floatingIp"`
	// ID of an existing port with at least one IP address to
	// associate with this floating IP.
	PortId *string `pulumi:"portId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a floating IP that can be used with
	// another networking resource, such as a load balancer. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// floating IP (which may or may not have a different address).
	Region *string `pulumi:"region"`
}

type FloatingIpAssociateState struct {
	FixedIp pulumi.StringPtrInput
	// IP Address of an existing floating IP.
	FloatingIp pulumi.StringPtrInput
	// ID of an existing port with at least one IP address to
	// associate with this floating IP.
	PortId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a floating IP that can be used with
	// another networking resource, such as a load balancer. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// floating IP (which may or may not have a different address).
	Region pulumi.StringPtrInput
}

func (FloatingIpAssociateState) ElementType() reflect.Type {
	return reflect.TypeOf((*floatingIpAssociateState)(nil)).Elem()
}

type floatingIpAssociateArgs struct {
	FixedIp *string `pulumi:"fixedIp"`
	// IP Address of an existing floating IP.
	FloatingIp string `pulumi:"floatingIp"`
	// ID of an existing port with at least one IP address to
	// associate with this floating IP.
	PortId string `pulumi:"portId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a floating IP that can be used with
	// another networking resource, such as a load balancer. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// floating IP (which may or may not have a different address).
	Region *string `pulumi:"region"`
}

// The set of arguments for constructing a FloatingIpAssociate resource.
type FloatingIpAssociateArgs struct {
	FixedIp pulumi.StringPtrInput
	// IP Address of an existing floating IP.
	FloatingIp pulumi.StringInput
	// ID of an existing port with at least one IP address to
	// associate with this floating IP.
	PortId pulumi.StringInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a floating IP that can be used with
	// another networking resource, such as a load balancer. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// floating IP (which may or may not have a different address).
	Region pulumi.StringPtrInput
}

func (FloatingIpAssociateArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*floatingIpAssociateArgs)(nil)).Elem()
}

type FloatingIpAssociateInput interface {
	pulumi.Input

	ToFloatingIpAssociateOutput() FloatingIpAssociateOutput
	ToFloatingIpAssociateOutputWithContext(ctx context.Context) FloatingIpAssociateOutput
}

func (*FloatingIpAssociate) ElementType() reflect.Type {
	return reflect.TypeOf((**FloatingIpAssociate)(nil)).Elem()
}

func (i *FloatingIpAssociate) ToFloatingIpAssociateOutput() FloatingIpAssociateOutput {
	return i.ToFloatingIpAssociateOutputWithContext(context.Background())
}

func (i *FloatingIpAssociate) ToFloatingIpAssociateOutputWithContext(ctx context.Context) FloatingIpAssociateOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FloatingIpAssociateOutput)
}

// FloatingIpAssociateArrayInput is an input type that accepts FloatingIpAssociateArray and FloatingIpAssociateArrayOutput values.
// You can construct a concrete instance of `FloatingIpAssociateArrayInput` via:
//
//	FloatingIpAssociateArray{ FloatingIpAssociateArgs{...} }
type FloatingIpAssociateArrayInput interface {
	pulumi.Input

	ToFloatingIpAssociateArrayOutput() FloatingIpAssociateArrayOutput
	ToFloatingIpAssociateArrayOutputWithContext(context.Context) FloatingIpAssociateArrayOutput
}

type FloatingIpAssociateArray []FloatingIpAssociateInput

func (FloatingIpAssociateArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*FloatingIpAssociate)(nil)).Elem()
}

func (i FloatingIpAssociateArray) ToFloatingIpAssociateArrayOutput() FloatingIpAssociateArrayOutput {
	return i.ToFloatingIpAssociateArrayOutputWithContext(context.Background())
}

func (i FloatingIpAssociateArray) ToFloatingIpAssociateArrayOutputWithContext(ctx context.Context) FloatingIpAssociateArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FloatingIpAssociateArrayOutput)
}

// FloatingIpAssociateMapInput is an input type that accepts FloatingIpAssociateMap and FloatingIpAssociateMapOutput values.
// You can construct a concrete instance of `FloatingIpAssociateMapInput` via:
//
//	FloatingIpAssociateMap{ "key": FloatingIpAssociateArgs{...} }
type FloatingIpAssociateMapInput interface {
	pulumi.Input

	ToFloatingIpAssociateMapOutput() FloatingIpAssociateMapOutput
	ToFloatingIpAssociateMapOutputWithContext(context.Context) FloatingIpAssociateMapOutput
}

type FloatingIpAssociateMap map[string]FloatingIpAssociateInput

func (FloatingIpAssociateMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*FloatingIpAssociate)(nil)).Elem()
}

func (i FloatingIpAssociateMap) ToFloatingIpAssociateMapOutput() FloatingIpAssociateMapOutput {
	return i.ToFloatingIpAssociateMapOutputWithContext(context.Background())
}

func (i FloatingIpAssociateMap) ToFloatingIpAssociateMapOutputWithContext(ctx context.Context) FloatingIpAssociateMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FloatingIpAssociateMapOutput)
}

type FloatingIpAssociateOutput struct{ *pulumi.OutputState }

func (FloatingIpAssociateOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**FloatingIpAssociate)(nil)).Elem()
}

func (o FloatingIpAssociateOutput) ToFloatingIpAssociateOutput() FloatingIpAssociateOutput {
	return o
}

func (o FloatingIpAssociateOutput) ToFloatingIpAssociateOutputWithContext(ctx context.Context) FloatingIpAssociateOutput {
	return o
}

func (o FloatingIpAssociateOutput) FixedIp() pulumi.StringOutput {
	return o.ApplyT(func(v *FloatingIpAssociate) pulumi.StringOutput { return v.FixedIp }).(pulumi.StringOutput)
}

// IP Address of an existing floating IP.
func (o FloatingIpAssociateOutput) FloatingIp() pulumi.StringOutput {
	return o.ApplyT(func(v *FloatingIpAssociate) pulumi.StringOutput { return v.FloatingIp }).(pulumi.StringOutput)
}

// ID of an existing port with at least one IP address to
// associate with this floating IP.
func (o FloatingIpAssociateOutput) PortId() pulumi.StringOutput {
	return o.ApplyT(func(v *FloatingIpAssociate) pulumi.StringOutput { return v.PortId }).(pulumi.StringOutput)
}

// The region in which to obtain the V2 Networking client.
// A Networking client is needed to create a floating IP that can be used with
// another networking resource, such as a load balancer. If omitted, the
// `region` argument of the provider is used. Changing this creates a new
// floating IP (which may or may not have a different address).
func (o FloatingIpAssociateOutput) Region() pulumi.StringOutput {
	return o.ApplyT(func(v *FloatingIpAssociate) pulumi.StringOutput { return v.Region }).(pulumi.StringOutput)
}

type FloatingIpAssociateArrayOutput struct{ *pulumi.OutputState }

func (FloatingIpAssociateArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*FloatingIpAssociate)(nil)).Elem()
}

func (o FloatingIpAssociateArrayOutput) ToFloatingIpAssociateArrayOutput() FloatingIpAssociateArrayOutput {
	return o
}

func (o FloatingIpAssociateArrayOutput) ToFloatingIpAssociateArrayOutputWithContext(ctx context.Context) FloatingIpAssociateArrayOutput {
	return o
}

func (o FloatingIpAssociateArrayOutput) Index(i pulumi.IntInput) FloatingIpAssociateOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *FloatingIpAssociate {
		return vs[0].([]*FloatingIpAssociate)[vs[1].(int)]
	}).(FloatingIpAssociateOutput)
}

type FloatingIpAssociateMapOutput struct{ *pulumi.OutputState }

func (FloatingIpAssociateMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*FloatingIpAssociate)(nil)).Elem()
}

func (o FloatingIpAssociateMapOutput) ToFloatingIpAssociateMapOutput() FloatingIpAssociateMapOutput {
	return o
}

func (o FloatingIpAssociateMapOutput) ToFloatingIpAssociateMapOutputWithContext(ctx context.Context) FloatingIpAssociateMapOutput {
	return o
}

func (o FloatingIpAssociateMapOutput) MapIndex(k pulumi.StringInput) FloatingIpAssociateOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *FloatingIpAssociate {
		return vs[0].(map[string]*FloatingIpAssociate)[vs[1].(string)]
	}).(FloatingIpAssociateOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*FloatingIpAssociateInput)(nil)).Elem(), &FloatingIpAssociate{})
	pulumi.RegisterInputType(reflect.TypeOf((*FloatingIpAssociateArrayInput)(nil)).Elem(), FloatingIpAssociateArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*FloatingIpAssociateMapInput)(nil)).Elem(), FloatingIpAssociateMap{})
	pulumi.RegisterOutputType(FloatingIpAssociateOutput{})
	pulumi.RegisterOutputType(FloatingIpAssociateArrayOutput{})
	pulumi.RegisterOutputType(FloatingIpAssociateMapOutput{})
}