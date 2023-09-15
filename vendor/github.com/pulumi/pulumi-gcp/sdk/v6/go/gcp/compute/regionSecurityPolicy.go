// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ## Example Usage
// ### Region Security Policy Basic
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
//			_, err := compute.NewRegionSecurityPolicy(ctx, "region-sec-policy-basic", &compute.RegionSecurityPolicyArgs{
//				Description: pulumi.String("basic region security policy"),
//				Type:        pulumi.String("CLOUD_ARMOR"),
//			}, pulumi.Provider(google_beta))
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
// ### Region Security Policy With Ddos Protection Config
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
//			_, err := compute.NewRegionSecurityPolicy(ctx, "region-sec-policy-ddos-protection", &compute.RegionSecurityPolicyArgs{
//				Description: pulumi.String("with ddos protection config"),
//				Type:        pulumi.String("CLOUD_ARMOR_NETWORK"),
//				DdosProtectionConfig: &compute.RegionSecurityPolicyDdosProtectionConfigArgs{
//					DdosProtection: pulumi.String("ADVANCED_PREVIEW"),
//				},
//			}, pulumi.Provider(google_beta))
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
// ### Region Security Policy With User Defined Fields
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
//			_, err := compute.NewRegionSecurityPolicy(ctx, "region-sec-policy-user-defined-fields", &compute.RegionSecurityPolicyArgs{
//				Description: pulumi.String("with user defined fields"),
//				Type:        pulumi.String("CLOUD_ARMOR_NETWORK"),
//				UserDefinedFields: compute.RegionSecurityPolicyUserDefinedFieldArray{
//					&compute.RegionSecurityPolicyUserDefinedFieldArgs{
//						Name:   pulumi.String("SIG1_AT_0"),
//						Base:   pulumi.String("UDP"),
//						Offset: pulumi.Int(8),
//						Size:   pulumi.Int(2),
//						Mask:   pulumi.String("0x8F00"),
//					},
//					&compute.RegionSecurityPolicyUserDefinedFieldArgs{
//						Name:   pulumi.String("SIG2_AT_8"),
//						Base:   pulumi.String("UDP"),
//						Offset: pulumi.Int(16),
//						Size:   pulumi.Int(4),
//						Mask:   pulumi.String("0xFFFFFFFF"),
//					},
//				},
//			}, pulumi.Provider(google_beta))
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
// # RegionSecurityPolicy can be imported using any of these accepted formats
//
// ```sh
//
//	$ pulumi import gcp:compute/regionSecurityPolicy:RegionSecurityPolicy default projects/{{project}}/regions/{{region}}/securityPolicies/{{name}}
//
// ```
//
// ```sh
//
//	$ pulumi import gcp:compute/regionSecurityPolicy:RegionSecurityPolicy default {{project}}/{{region}}/{{name}}
//
// ```
//
// ```sh
//
//	$ pulumi import gcp:compute/regionSecurityPolicy:RegionSecurityPolicy default {{region}}/{{name}}
//
// ```
//
// ```sh
//
//	$ pulumi import gcp:compute/regionSecurityPolicy:RegionSecurityPolicy default {{name}}
//
// ```
type RegionSecurityPolicy struct {
	pulumi.CustomResourceState

	// Configuration for Google Cloud Armor DDOS Proctection Config.
	// Structure is documented below.
	DdosProtectionConfig RegionSecurityPolicyDdosProtectionConfigPtrOutput `pulumi:"ddosProtectionConfig"`
	// An optional description of this resource. Provide this property when you create the resource.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Fingerprint of this resource. This field is used internally during
	// updates of this resource.
	Fingerprint pulumi.StringOutput `pulumi:"fingerprint"`
	// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
	// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
	//
	// ***
	Name pulumi.StringOutput `pulumi:"name"`
	// The unique identifier for the resource. This identifier is defined by the server.
	PolicyId pulumi.StringOutput `pulumi:"policyId"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringOutput `pulumi:"project"`
	// The Region in which the created Region Security Policy should reside.
	// If it is not provided, the provider region is used.
	Region pulumi.StringOutput `pulumi:"region"`
	// Server-defined URL for the resource.
	SelfLink pulumi.StringOutput `pulumi:"selfLink"`
	// Server-defined URL for this resource with the resource id.
	SelfLinkWithPolicyId pulumi.StringOutput `pulumi:"selfLinkWithPolicyId"`
	// The type indicates the intended use of the security policy.
	// - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
	// - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
	// - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
	//   This field can be set only at resource creation time.
	//   Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
	Type pulumi.StringPtrOutput `pulumi:"type"`
	// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
	// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
	// Rules may then specify matching values for these fields.
	// Structure is documented below.
	UserDefinedFields RegionSecurityPolicyUserDefinedFieldArrayOutput `pulumi:"userDefinedFields"`
}

// NewRegionSecurityPolicy registers a new resource with the given unique name, arguments, and options.
func NewRegionSecurityPolicy(ctx *pulumi.Context,
	name string, args *RegionSecurityPolicyArgs, opts ...pulumi.ResourceOption) (*RegionSecurityPolicy, error) {
	if args == nil {
		args = &RegionSecurityPolicyArgs{}
	}

	opts = internal.PkgResourceDefaultOpts(opts)
	var resource RegionSecurityPolicy
	err := ctx.RegisterResource("gcp:compute/regionSecurityPolicy:RegionSecurityPolicy", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetRegionSecurityPolicy gets an existing RegionSecurityPolicy resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetRegionSecurityPolicy(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *RegionSecurityPolicyState, opts ...pulumi.ResourceOption) (*RegionSecurityPolicy, error) {
	var resource RegionSecurityPolicy
	err := ctx.ReadResource("gcp:compute/regionSecurityPolicy:RegionSecurityPolicy", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering RegionSecurityPolicy resources.
type regionSecurityPolicyState struct {
	// Configuration for Google Cloud Armor DDOS Proctection Config.
	// Structure is documented below.
	DdosProtectionConfig *RegionSecurityPolicyDdosProtectionConfig `pulumi:"ddosProtectionConfig"`
	// An optional description of this resource. Provide this property when you create the resource.
	Description *string `pulumi:"description"`
	// Fingerprint of this resource. This field is used internally during
	// updates of this resource.
	Fingerprint *string `pulumi:"fingerprint"`
	// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
	// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
	//
	// ***
	Name *string `pulumi:"name"`
	// The unique identifier for the resource. This identifier is defined by the server.
	PolicyId *string `pulumi:"policyId"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
	// The Region in which the created Region Security Policy should reside.
	// If it is not provided, the provider region is used.
	Region *string `pulumi:"region"`
	// Server-defined URL for the resource.
	SelfLink *string `pulumi:"selfLink"`
	// Server-defined URL for this resource with the resource id.
	SelfLinkWithPolicyId *string `pulumi:"selfLinkWithPolicyId"`
	// The type indicates the intended use of the security policy.
	// - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
	// - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
	// - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
	//   This field can be set only at resource creation time.
	//   Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
	Type *string `pulumi:"type"`
	// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
	// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
	// Rules may then specify matching values for these fields.
	// Structure is documented below.
	UserDefinedFields []RegionSecurityPolicyUserDefinedField `pulumi:"userDefinedFields"`
}

type RegionSecurityPolicyState struct {
	// Configuration for Google Cloud Armor DDOS Proctection Config.
	// Structure is documented below.
	DdosProtectionConfig RegionSecurityPolicyDdosProtectionConfigPtrInput
	// An optional description of this resource. Provide this property when you create the resource.
	Description pulumi.StringPtrInput
	// Fingerprint of this resource. This field is used internally during
	// updates of this resource.
	Fingerprint pulumi.StringPtrInput
	// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
	// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
	//
	// ***
	Name pulumi.StringPtrInput
	// The unique identifier for the resource. This identifier is defined by the server.
	PolicyId pulumi.StringPtrInput
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput
	// The Region in which the created Region Security Policy should reside.
	// If it is not provided, the provider region is used.
	Region pulumi.StringPtrInput
	// Server-defined URL for the resource.
	SelfLink pulumi.StringPtrInput
	// Server-defined URL for this resource with the resource id.
	SelfLinkWithPolicyId pulumi.StringPtrInput
	// The type indicates the intended use of the security policy.
	// - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
	// - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
	// - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
	//   This field can be set only at resource creation time.
	//   Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
	Type pulumi.StringPtrInput
	// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
	// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
	// Rules may then specify matching values for these fields.
	// Structure is documented below.
	UserDefinedFields RegionSecurityPolicyUserDefinedFieldArrayInput
}

func (RegionSecurityPolicyState) ElementType() reflect.Type {
	return reflect.TypeOf((*regionSecurityPolicyState)(nil)).Elem()
}

type regionSecurityPolicyArgs struct {
	// Configuration for Google Cloud Armor DDOS Proctection Config.
	// Structure is documented below.
	DdosProtectionConfig *RegionSecurityPolicyDdosProtectionConfig `pulumi:"ddosProtectionConfig"`
	// An optional description of this resource. Provide this property when you create the resource.
	Description *string `pulumi:"description"`
	// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
	// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
	//
	// ***
	Name *string `pulumi:"name"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
	// The Region in which the created Region Security Policy should reside.
	// If it is not provided, the provider region is used.
	Region *string `pulumi:"region"`
	// The type indicates the intended use of the security policy.
	// - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
	// - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
	// - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
	//   This field can be set only at resource creation time.
	//   Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
	Type *string `pulumi:"type"`
	// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
	// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
	// Rules may then specify matching values for these fields.
	// Structure is documented below.
	UserDefinedFields []RegionSecurityPolicyUserDefinedField `pulumi:"userDefinedFields"`
}

// The set of arguments for constructing a RegionSecurityPolicy resource.
type RegionSecurityPolicyArgs struct {
	// Configuration for Google Cloud Armor DDOS Proctection Config.
	// Structure is documented below.
	DdosProtectionConfig RegionSecurityPolicyDdosProtectionConfigPtrInput
	// An optional description of this resource. Provide this property when you create the resource.
	Description pulumi.StringPtrInput
	// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
	// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
	//
	// ***
	Name pulumi.StringPtrInput
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput
	// The Region in which the created Region Security Policy should reside.
	// If it is not provided, the provider region is used.
	Region pulumi.StringPtrInput
	// The type indicates the intended use of the security policy.
	// - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
	// - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
	// - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
	//   This field can be set only at resource creation time.
	//   Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
	Type pulumi.StringPtrInput
	// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
	// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
	// Rules may then specify matching values for these fields.
	// Structure is documented below.
	UserDefinedFields RegionSecurityPolicyUserDefinedFieldArrayInput
}

func (RegionSecurityPolicyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*regionSecurityPolicyArgs)(nil)).Elem()
}

type RegionSecurityPolicyInput interface {
	pulumi.Input

	ToRegionSecurityPolicyOutput() RegionSecurityPolicyOutput
	ToRegionSecurityPolicyOutputWithContext(ctx context.Context) RegionSecurityPolicyOutput
}

func (*RegionSecurityPolicy) ElementType() reflect.Type {
	return reflect.TypeOf((**RegionSecurityPolicy)(nil)).Elem()
}

func (i *RegionSecurityPolicy) ToRegionSecurityPolicyOutput() RegionSecurityPolicyOutput {
	return i.ToRegionSecurityPolicyOutputWithContext(context.Background())
}

func (i *RegionSecurityPolicy) ToRegionSecurityPolicyOutputWithContext(ctx context.Context) RegionSecurityPolicyOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RegionSecurityPolicyOutput)
}

// RegionSecurityPolicyArrayInput is an input type that accepts RegionSecurityPolicyArray and RegionSecurityPolicyArrayOutput values.
// You can construct a concrete instance of `RegionSecurityPolicyArrayInput` via:
//
//	RegionSecurityPolicyArray{ RegionSecurityPolicyArgs{...} }
type RegionSecurityPolicyArrayInput interface {
	pulumi.Input

	ToRegionSecurityPolicyArrayOutput() RegionSecurityPolicyArrayOutput
	ToRegionSecurityPolicyArrayOutputWithContext(context.Context) RegionSecurityPolicyArrayOutput
}

type RegionSecurityPolicyArray []RegionSecurityPolicyInput

func (RegionSecurityPolicyArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*RegionSecurityPolicy)(nil)).Elem()
}

func (i RegionSecurityPolicyArray) ToRegionSecurityPolicyArrayOutput() RegionSecurityPolicyArrayOutput {
	return i.ToRegionSecurityPolicyArrayOutputWithContext(context.Background())
}

func (i RegionSecurityPolicyArray) ToRegionSecurityPolicyArrayOutputWithContext(ctx context.Context) RegionSecurityPolicyArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RegionSecurityPolicyArrayOutput)
}

// RegionSecurityPolicyMapInput is an input type that accepts RegionSecurityPolicyMap and RegionSecurityPolicyMapOutput values.
// You can construct a concrete instance of `RegionSecurityPolicyMapInput` via:
//
//	RegionSecurityPolicyMap{ "key": RegionSecurityPolicyArgs{...} }
type RegionSecurityPolicyMapInput interface {
	pulumi.Input

	ToRegionSecurityPolicyMapOutput() RegionSecurityPolicyMapOutput
	ToRegionSecurityPolicyMapOutputWithContext(context.Context) RegionSecurityPolicyMapOutput
}

type RegionSecurityPolicyMap map[string]RegionSecurityPolicyInput

func (RegionSecurityPolicyMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*RegionSecurityPolicy)(nil)).Elem()
}

func (i RegionSecurityPolicyMap) ToRegionSecurityPolicyMapOutput() RegionSecurityPolicyMapOutput {
	return i.ToRegionSecurityPolicyMapOutputWithContext(context.Background())
}

func (i RegionSecurityPolicyMap) ToRegionSecurityPolicyMapOutputWithContext(ctx context.Context) RegionSecurityPolicyMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RegionSecurityPolicyMapOutput)
}

type RegionSecurityPolicyOutput struct{ *pulumi.OutputState }

func (RegionSecurityPolicyOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**RegionSecurityPolicy)(nil)).Elem()
}

func (o RegionSecurityPolicyOutput) ToRegionSecurityPolicyOutput() RegionSecurityPolicyOutput {
	return o
}

func (o RegionSecurityPolicyOutput) ToRegionSecurityPolicyOutputWithContext(ctx context.Context) RegionSecurityPolicyOutput {
	return o
}

// Configuration for Google Cloud Armor DDOS Proctection Config.
// Structure is documented below.
func (o RegionSecurityPolicyOutput) DdosProtectionConfig() RegionSecurityPolicyDdosProtectionConfigPtrOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) RegionSecurityPolicyDdosProtectionConfigPtrOutput {
		return v.DdosProtectionConfig
	}).(RegionSecurityPolicyDdosProtectionConfigPtrOutput)
}

// An optional description of this resource. Provide this property when you create the resource.
func (o RegionSecurityPolicyOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Fingerprint of this resource. This field is used internally during
// updates of this resource.
func (o RegionSecurityPolicyOutput) Fingerprint() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.Fingerprint }).(pulumi.StringOutput)
}

// Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035.
// Specifically, the name must be 1-63 characters long and match the regular expression a-z? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
//
// ***
func (o RegionSecurityPolicyOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// The unique identifier for the resource. This identifier is defined by the server.
func (o RegionSecurityPolicyOutput) PolicyId() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.PolicyId }).(pulumi.StringOutput)
}

// The ID of the project in which the resource belongs.
// If it is not provided, the provider project is used.
func (o RegionSecurityPolicyOutput) Project() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.Project }).(pulumi.StringOutput)
}

// The Region in which the created Region Security Policy should reside.
// If it is not provided, the provider region is used.
func (o RegionSecurityPolicyOutput) Region() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.Region }).(pulumi.StringOutput)
}

// Server-defined URL for the resource.
func (o RegionSecurityPolicyOutput) SelfLink() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.SelfLink }).(pulumi.StringOutput)
}

// Server-defined URL for this resource with the resource id.
func (o RegionSecurityPolicyOutput) SelfLinkWithPolicyId() pulumi.StringOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringOutput { return v.SelfLinkWithPolicyId }).(pulumi.StringOutput)
}

// The type indicates the intended use of the security policy.
//   - CLOUD_ARMOR: Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services. They filter requests before they hit the origin servers.
//   - CLOUD_ARMOR_EDGE: Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage). They filter requests before the request is served from Google's cache.
//   - CLOUD_ARMOR_NETWORK: Cloud Armor network policies can be configured to filter packets targeting network load balancing resources such as backend services, target pools, target instances, and instances with external IPs. They filter requests before the request is served from the application.
//     This field can be set only at resource creation time.
//     Possible values are: `CLOUD_ARMOR`, `CLOUD_ARMOR_EDGE`, `CLOUD_ARMOR_NETWORK`.
func (o RegionSecurityPolicyOutput) Type() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) pulumi.StringPtrOutput { return v.Type }).(pulumi.StringPtrOutput)
}

// Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies.
// A user-defined field consists of up to 4 bytes extracted from a fixed offset in the packet, relative to the IPv4, IPv6, TCP, or UDP header, with an optional mask to select certain bits.
// Rules may then specify matching values for these fields.
// Structure is documented below.
func (o RegionSecurityPolicyOutput) UserDefinedFields() RegionSecurityPolicyUserDefinedFieldArrayOutput {
	return o.ApplyT(func(v *RegionSecurityPolicy) RegionSecurityPolicyUserDefinedFieldArrayOutput {
		return v.UserDefinedFields
	}).(RegionSecurityPolicyUserDefinedFieldArrayOutput)
}

type RegionSecurityPolicyArrayOutput struct{ *pulumi.OutputState }

func (RegionSecurityPolicyArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*RegionSecurityPolicy)(nil)).Elem()
}

func (o RegionSecurityPolicyArrayOutput) ToRegionSecurityPolicyArrayOutput() RegionSecurityPolicyArrayOutput {
	return o
}

func (o RegionSecurityPolicyArrayOutput) ToRegionSecurityPolicyArrayOutputWithContext(ctx context.Context) RegionSecurityPolicyArrayOutput {
	return o
}

func (o RegionSecurityPolicyArrayOutput) Index(i pulumi.IntInput) RegionSecurityPolicyOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *RegionSecurityPolicy {
		return vs[0].([]*RegionSecurityPolicy)[vs[1].(int)]
	}).(RegionSecurityPolicyOutput)
}

type RegionSecurityPolicyMapOutput struct{ *pulumi.OutputState }

func (RegionSecurityPolicyMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*RegionSecurityPolicy)(nil)).Elem()
}

func (o RegionSecurityPolicyMapOutput) ToRegionSecurityPolicyMapOutput() RegionSecurityPolicyMapOutput {
	return o
}

func (o RegionSecurityPolicyMapOutput) ToRegionSecurityPolicyMapOutputWithContext(ctx context.Context) RegionSecurityPolicyMapOutput {
	return o
}

func (o RegionSecurityPolicyMapOutput) MapIndex(k pulumi.StringInput) RegionSecurityPolicyOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *RegionSecurityPolicy {
		return vs[0].(map[string]*RegionSecurityPolicy)[vs[1].(string)]
	}).(RegionSecurityPolicyOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*RegionSecurityPolicyInput)(nil)).Elem(), &RegionSecurityPolicy{})
	pulumi.RegisterInputType(reflect.TypeOf((*RegionSecurityPolicyArrayInput)(nil)).Elem(), RegionSecurityPolicyArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*RegionSecurityPolicyMapInput)(nil)).Elem(), RegionSecurityPolicyMap{})
	pulumi.RegisterOutputType(RegionSecurityPolicyOutput{})
	pulumi.RegisterOutputType(RegionSecurityPolicyArrayOutput{})
	pulumi.RegisterOutputType(RegionSecurityPolicyMapOutput{})
}