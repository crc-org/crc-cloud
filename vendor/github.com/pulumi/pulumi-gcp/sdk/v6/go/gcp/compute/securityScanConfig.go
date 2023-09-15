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

// A ScanConfig resource contains the configurations to launch a scan.
//
// To get more information about ScanConfig, see:
//
// * [API documentation](https://cloud.google.com/security-scanner/docs/reference/rest/v1beta/projects.scanConfigs)
// * How-to Guides
//   - [Using Cloud Security Scanner](https://cloud.google.com/security-scanner/docs/scanning)
//
// > **Warning:** All arguments including `authentication.google_account.password` and `authentication.custom_account.password` will be stored in the raw state as plain-text.
//
// ## Example Usage
// ### Scan Config Basic
//
// ```go
// package main
//
// import (
//
//	"fmt"
//
//	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			scannerStaticIp, err := compute.NewAddress(ctx, "scannerStaticIp", nil, pulumi.Provider(google_beta))
//			if err != nil {
//				return err
//			}
//			_, err = compute.NewSecurityScanConfig(ctx, "scan-config", &compute.SecurityScanConfigArgs{
//				DisplayName: pulumi.String("scan-config"),
//				StartingUrls: pulumi.StringArray{
//					scannerStaticIp.Address.ApplyT(func(address string) (string, error) {
//						return fmt.Sprintf("http://%v", address), nil
//					}).(pulumi.StringOutput),
//				},
//				TargetPlatforms: pulumi.StringArray{
//					pulumi.String("COMPUTE"),
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
// # ScanConfig can be imported using any of these accepted formats
//
// ```sh
//
//	$ pulumi import gcp:compute/securityScanConfig:SecurityScanConfig default projects/{{project}}/scanConfigs/{{name}}
//
// ```
//
// ```sh
//
//	$ pulumi import gcp:compute/securityScanConfig:SecurityScanConfig default {{project}}/{{name}}
//
// ```
//
// ```sh
//
//	$ pulumi import gcp:compute/securityScanConfig:SecurityScanConfig default {{name}}
//
// ```
type SecurityScanConfig struct {
	pulumi.CustomResourceState

	// The authentication configuration.
	// If specified, service will use the authentication configuration during scanning.
	// Structure is documented below.
	Authentication SecurityScanConfigAuthenticationPtrOutput `pulumi:"authentication"`
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns pulumi.StringArrayOutput `pulumi:"blacklistPatterns"`
	// The user provider display name of the ScanConfig.
	DisplayName pulumi.StringOutput `pulumi:"displayName"`
	// Controls export of scan configurations and results to Cloud Security Command Center.
	// Default value is `ENABLED`.
	// Possible values are: `ENABLED`, `DISABLED`.
	ExportToSecurityCommandCenter pulumi.StringPtrOutput `pulumi:"exportToSecurityCommandCenter"`
	// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
	// Defaults to 15.
	MaxQps pulumi.IntPtrOutput `pulumi:"maxQps"`
	// A server defined name for this index. Format:
	// `projects/{{project}}/scanConfigs/{{server_generated_id}}`
	Name pulumi.StringOutput `pulumi:"name"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringOutput `pulumi:"project"`
	// The schedule of the ScanConfig
	// Structure is documented below.
	Schedule SecurityScanConfigSchedulePtrOutput `pulumi:"schedule"`
	// The starting URLs from which the scanner finds site pages.
	//
	// ***
	StartingUrls pulumi.StringArrayOutput `pulumi:"startingUrls"`
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
	// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
	TargetPlatforms pulumi.StringArrayOutput `pulumi:"targetPlatforms"`
	// Type of the user agents used for scanning
	// Default value is `CHROME_LINUX`.
	// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
	UserAgent pulumi.StringPtrOutput `pulumi:"userAgent"`
}

// NewSecurityScanConfig registers a new resource with the given unique name, arguments, and options.
func NewSecurityScanConfig(ctx *pulumi.Context,
	name string, args *SecurityScanConfigArgs, opts ...pulumi.ResourceOption) (*SecurityScanConfig, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.DisplayName == nil {
		return nil, errors.New("invalid value for required argument 'DisplayName'")
	}
	if args.StartingUrls == nil {
		return nil, errors.New("invalid value for required argument 'StartingUrls'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource SecurityScanConfig
	err := ctx.RegisterResource("gcp:compute/securityScanConfig:SecurityScanConfig", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetSecurityScanConfig gets an existing SecurityScanConfig resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetSecurityScanConfig(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *SecurityScanConfigState, opts ...pulumi.ResourceOption) (*SecurityScanConfig, error) {
	var resource SecurityScanConfig
	err := ctx.ReadResource("gcp:compute/securityScanConfig:SecurityScanConfig", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering SecurityScanConfig resources.
type securityScanConfigState struct {
	// The authentication configuration.
	// If specified, service will use the authentication configuration during scanning.
	// Structure is documented below.
	Authentication *SecurityScanConfigAuthentication `pulumi:"authentication"`
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns []string `pulumi:"blacklistPatterns"`
	// The user provider display name of the ScanConfig.
	DisplayName *string `pulumi:"displayName"`
	// Controls export of scan configurations and results to Cloud Security Command Center.
	// Default value is `ENABLED`.
	// Possible values are: `ENABLED`, `DISABLED`.
	ExportToSecurityCommandCenter *string `pulumi:"exportToSecurityCommandCenter"`
	// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
	// Defaults to 15.
	MaxQps *int `pulumi:"maxQps"`
	// A server defined name for this index. Format:
	// `projects/{{project}}/scanConfigs/{{server_generated_id}}`
	Name *string `pulumi:"name"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
	// The schedule of the ScanConfig
	// Structure is documented below.
	Schedule *SecurityScanConfigSchedule `pulumi:"schedule"`
	// The starting URLs from which the scanner finds site pages.
	//
	// ***
	StartingUrls []string `pulumi:"startingUrls"`
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
	// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
	TargetPlatforms []string `pulumi:"targetPlatforms"`
	// Type of the user agents used for scanning
	// Default value is `CHROME_LINUX`.
	// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
	UserAgent *string `pulumi:"userAgent"`
}

type SecurityScanConfigState struct {
	// The authentication configuration.
	// If specified, service will use the authentication configuration during scanning.
	// Structure is documented below.
	Authentication SecurityScanConfigAuthenticationPtrInput
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns pulumi.StringArrayInput
	// The user provider display name of the ScanConfig.
	DisplayName pulumi.StringPtrInput
	// Controls export of scan configurations and results to Cloud Security Command Center.
	// Default value is `ENABLED`.
	// Possible values are: `ENABLED`, `DISABLED`.
	ExportToSecurityCommandCenter pulumi.StringPtrInput
	// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
	// Defaults to 15.
	MaxQps pulumi.IntPtrInput
	// A server defined name for this index. Format:
	// `projects/{{project}}/scanConfigs/{{server_generated_id}}`
	Name pulumi.StringPtrInput
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput
	// The schedule of the ScanConfig
	// Structure is documented below.
	Schedule SecurityScanConfigSchedulePtrInput
	// The starting URLs from which the scanner finds site pages.
	//
	// ***
	StartingUrls pulumi.StringArrayInput
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
	// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
	TargetPlatforms pulumi.StringArrayInput
	// Type of the user agents used for scanning
	// Default value is `CHROME_LINUX`.
	// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
	UserAgent pulumi.StringPtrInput
}

func (SecurityScanConfigState) ElementType() reflect.Type {
	return reflect.TypeOf((*securityScanConfigState)(nil)).Elem()
}

type securityScanConfigArgs struct {
	// The authentication configuration.
	// If specified, service will use the authentication configuration during scanning.
	// Structure is documented below.
	Authentication *SecurityScanConfigAuthentication `pulumi:"authentication"`
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns []string `pulumi:"blacklistPatterns"`
	// The user provider display name of the ScanConfig.
	DisplayName string `pulumi:"displayName"`
	// Controls export of scan configurations and results to Cloud Security Command Center.
	// Default value is `ENABLED`.
	// Possible values are: `ENABLED`, `DISABLED`.
	ExportToSecurityCommandCenter *string `pulumi:"exportToSecurityCommandCenter"`
	// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
	// Defaults to 15.
	MaxQps *int `pulumi:"maxQps"`
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project *string `pulumi:"project"`
	// The schedule of the ScanConfig
	// Structure is documented below.
	Schedule *SecurityScanConfigSchedule `pulumi:"schedule"`
	// The starting URLs from which the scanner finds site pages.
	//
	// ***
	StartingUrls []string `pulumi:"startingUrls"`
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
	// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
	TargetPlatforms []string `pulumi:"targetPlatforms"`
	// Type of the user agents used for scanning
	// Default value is `CHROME_LINUX`.
	// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
	UserAgent *string `pulumi:"userAgent"`
}

// The set of arguments for constructing a SecurityScanConfig resource.
type SecurityScanConfigArgs struct {
	// The authentication configuration.
	// If specified, service will use the authentication configuration during scanning.
	// Structure is documented below.
	Authentication SecurityScanConfigAuthenticationPtrInput
	// The blacklist URL patterns as described in
	// https://cloud.google.com/security-scanner/docs/excluded-urls
	BlacklistPatterns pulumi.StringArrayInput
	// The user provider display name of the ScanConfig.
	DisplayName pulumi.StringInput
	// Controls export of scan configurations and results to Cloud Security Command Center.
	// Default value is `ENABLED`.
	// Possible values are: `ENABLED`, `DISABLED`.
	ExportToSecurityCommandCenter pulumi.StringPtrInput
	// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
	// Defaults to 15.
	MaxQps pulumi.IntPtrInput
	// The ID of the project in which the resource belongs.
	// If it is not provided, the provider project is used.
	Project pulumi.StringPtrInput
	// The schedule of the ScanConfig
	// Structure is documented below.
	Schedule SecurityScanConfigSchedulePtrInput
	// The starting URLs from which the scanner finds site pages.
	//
	// ***
	StartingUrls pulumi.StringArrayInput
	// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
	// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
	TargetPlatforms pulumi.StringArrayInput
	// Type of the user agents used for scanning
	// Default value is `CHROME_LINUX`.
	// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
	UserAgent pulumi.StringPtrInput
}

func (SecurityScanConfigArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*securityScanConfigArgs)(nil)).Elem()
}

type SecurityScanConfigInput interface {
	pulumi.Input

	ToSecurityScanConfigOutput() SecurityScanConfigOutput
	ToSecurityScanConfigOutputWithContext(ctx context.Context) SecurityScanConfigOutput
}

func (*SecurityScanConfig) ElementType() reflect.Type {
	return reflect.TypeOf((**SecurityScanConfig)(nil)).Elem()
}

func (i *SecurityScanConfig) ToSecurityScanConfigOutput() SecurityScanConfigOutput {
	return i.ToSecurityScanConfigOutputWithContext(context.Background())
}

func (i *SecurityScanConfig) ToSecurityScanConfigOutputWithContext(ctx context.Context) SecurityScanConfigOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityScanConfigOutput)
}

// SecurityScanConfigArrayInput is an input type that accepts SecurityScanConfigArray and SecurityScanConfigArrayOutput values.
// You can construct a concrete instance of `SecurityScanConfigArrayInput` via:
//
//	SecurityScanConfigArray{ SecurityScanConfigArgs{...} }
type SecurityScanConfigArrayInput interface {
	pulumi.Input

	ToSecurityScanConfigArrayOutput() SecurityScanConfigArrayOutput
	ToSecurityScanConfigArrayOutputWithContext(context.Context) SecurityScanConfigArrayOutput
}

type SecurityScanConfigArray []SecurityScanConfigInput

func (SecurityScanConfigArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SecurityScanConfig)(nil)).Elem()
}

func (i SecurityScanConfigArray) ToSecurityScanConfigArrayOutput() SecurityScanConfigArrayOutput {
	return i.ToSecurityScanConfigArrayOutputWithContext(context.Background())
}

func (i SecurityScanConfigArray) ToSecurityScanConfigArrayOutputWithContext(ctx context.Context) SecurityScanConfigArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityScanConfigArrayOutput)
}

// SecurityScanConfigMapInput is an input type that accepts SecurityScanConfigMap and SecurityScanConfigMapOutput values.
// You can construct a concrete instance of `SecurityScanConfigMapInput` via:
//
//	SecurityScanConfigMap{ "key": SecurityScanConfigArgs{...} }
type SecurityScanConfigMapInput interface {
	pulumi.Input

	ToSecurityScanConfigMapOutput() SecurityScanConfigMapOutput
	ToSecurityScanConfigMapOutputWithContext(context.Context) SecurityScanConfigMapOutput
}

type SecurityScanConfigMap map[string]SecurityScanConfigInput

func (SecurityScanConfigMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SecurityScanConfig)(nil)).Elem()
}

func (i SecurityScanConfigMap) ToSecurityScanConfigMapOutput() SecurityScanConfigMapOutput {
	return i.ToSecurityScanConfigMapOutputWithContext(context.Background())
}

func (i SecurityScanConfigMap) ToSecurityScanConfigMapOutputWithContext(ctx context.Context) SecurityScanConfigMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityScanConfigMapOutput)
}

type SecurityScanConfigOutput struct{ *pulumi.OutputState }

func (SecurityScanConfigOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**SecurityScanConfig)(nil)).Elem()
}

func (o SecurityScanConfigOutput) ToSecurityScanConfigOutput() SecurityScanConfigOutput {
	return o
}

func (o SecurityScanConfigOutput) ToSecurityScanConfigOutputWithContext(ctx context.Context) SecurityScanConfigOutput {
	return o
}

// The authentication configuration.
// If specified, service will use the authentication configuration during scanning.
// Structure is documented below.
func (o SecurityScanConfigOutput) Authentication() SecurityScanConfigAuthenticationPtrOutput {
	return o.ApplyT(func(v *SecurityScanConfig) SecurityScanConfigAuthenticationPtrOutput { return v.Authentication }).(SecurityScanConfigAuthenticationPtrOutput)
}

// The blacklist URL patterns as described in
// https://cloud.google.com/security-scanner/docs/excluded-urls
func (o SecurityScanConfigOutput) BlacklistPatterns() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringArrayOutput { return v.BlacklistPatterns }).(pulumi.StringArrayOutput)
}

// The user provider display name of the ScanConfig.
func (o SecurityScanConfigOutput) DisplayName() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringOutput { return v.DisplayName }).(pulumi.StringOutput)
}

// Controls export of scan configurations and results to Cloud Security Command Center.
// Default value is `ENABLED`.
// Possible values are: `ENABLED`, `DISABLED`.
func (o SecurityScanConfigOutput) ExportToSecurityCommandCenter() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringPtrOutput { return v.ExportToSecurityCommandCenter }).(pulumi.StringPtrOutput)
}

// The maximum QPS during scanning. A valid value ranges from 5 to 20 inclusively.
// Defaults to 15.
func (o SecurityScanConfigOutput) MaxQps() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.IntPtrOutput { return v.MaxQps }).(pulumi.IntPtrOutput)
}

// A server defined name for this index. Format:
// `projects/{{project}}/scanConfigs/{{server_generated_id}}`
func (o SecurityScanConfigOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// The ID of the project in which the resource belongs.
// If it is not provided, the provider project is used.
func (o SecurityScanConfigOutput) Project() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringOutput { return v.Project }).(pulumi.StringOutput)
}

// The schedule of the ScanConfig
// Structure is documented below.
func (o SecurityScanConfigOutput) Schedule() SecurityScanConfigSchedulePtrOutput {
	return o.ApplyT(func(v *SecurityScanConfig) SecurityScanConfigSchedulePtrOutput { return v.Schedule }).(SecurityScanConfigSchedulePtrOutput)
}

// The starting URLs from which the scanner finds site pages.
//
// ***
func (o SecurityScanConfigOutput) StartingUrls() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringArrayOutput { return v.StartingUrls }).(pulumi.StringArrayOutput)
}

// Set of Cloud Platforms targeted by the scan. If empty, APP_ENGINE will be used as a default.
// Each value may be one of: `APP_ENGINE`, `COMPUTE`.
func (o SecurityScanConfigOutput) TargetPlatforms() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringArrayOutput { return v.TargetPlatforms }).(pulumi.StringArrayOutput)
}

// Type of the user agents used for scanning
// Default value is `CHROME_LINUX`.
// Possible values are: `USER_AGENT_UNSPECIFIED`, `CHROME_LINUX`, `CHROME_ANDROID`, `SAFARI_IPHONE`.
func (o SecurityScanConfigOutput) UserAgent() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *SecurityScanConfig) pulumi.StringPtrOutput { return v.UserAgent }).(pulumi.StringPtrOutput)
}

type SecurityScanConfigArrayOutput struct{ *pulumi.OutputState }

func (SecurityScanConfigArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SecurityScanConfig)(nil)).Elem()
}

func (o SecurityScanConfigArrayOutput) ToSecurityScanConfigArrayOutput() SecurityScanConfigArrayOutput {
	return o
}

func (o SecurityScanConfigArrayOutput) ToSecurityScanConfigArrayOutputWithContext(ctx context.Context) SecurityScanConfigArrayOutput {
	return o
}

func (o SecurityScanConfigArrayOutput) Index(i pulumi.IntInput) SecurityScanConfigOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *SecurityScanConfig {
		return vs[0].([]*SecurityScanConfig)[vs[1].(int)]
	}).(SecurityScanConfigOutput)
}

type SecurityScanConfigMapOutput struct{ *pulumi.OutputState }

func (SecurityScanConfigMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SecurityScanConfig)(nil)).Elem()
}

func (o SecurityScanConfigMapOutput) ToSecurityScanConfigMapOutput() SecurityScanConfigMapOutput {
	return o
}

func (o SecurityScanConfigMapOutput) ToSecurityScanConfigMapOutputWithContext(ctx context.Context) SecurityScanConfigMapOutput {
	return o
}

func (o SecurityScanConfigMapOutput) MapIndex(k pulumi.StringInput) SecurityScanConfigOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *SecurityScanConfig {
		return vs[0].(map[string]*SecurityScanConfig)[vs[1].(string)]
	}).(SecurityScanConfigOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityScanConfigInput)(nil)).Elem(), &SecurityScanConfig{})
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityScanConfigArrayInput)(nil)).Elem(), SecurityScanConfigArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityScanConfigMapInput)(nil)).Elem(), SecurityScanConfigMap{})
	pulumi.RegisterOutputType(SecurityScanConfigOutput{})
	pulumi.RegisterOutputType(SecurityScanConfigArrayOutput{})
	pulumi.RegisterOutputType(SecurityScanConfigMapOutput{})
}