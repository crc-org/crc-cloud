// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package iam

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides a resource for adding an IAM User to IAM Groups. This
// resource can be used multiple times with the same user for non-overlapping
// groups.
//
// To exclusively manage the users in a group, see the
// `iam.GroupMembership` resource.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			user1, err := iam.NewUser(ctx, "user1", nil)
//			if err != nil {
//				return err
//			}
//			group1, err := iam.NewGroup(ctx, "group1", nil)
//			if err != nil {
//				return err
//			}
//			group2, err := iam.NewGroup(ctx, "group2", nil)
//			if err != nil {
//				return err
//			}
//			_, err = iam.NewUserGroupMembership(ctx, "example1", &iam.UserGroupMembershipArgs{
//				User: user1.Name,
//				Groups: pulumi.StringArray{
//					group1.Name,
//					group2.Name,
//				},
//			})
//			if err != nil {
//				return err
//			}
//			group3, err := iam.NewGroup(ctx, "group3", nil)
//			if err != nil {
//				return err
//			}
//			_, err = iam.NewUserGroupMembership(ctx, "example2", &iam.UserGroupMembershipArgs{
//				User: user1.Name,
//				Groups: pulumi.StringArray{
//					group3.Name,
//				},
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
// IAM user group membership can be imported using the user name and group names separated by `/`.
//
// ```sh
//
//	$ pulumi import aws:iam/userGroupMembership:UserGroupMembership example1 user1/group1/group2
//
// ```
type UserGroupMembership struct {
	pulumi.CustomResourceState

	// A list of IAM Groups to add the user to
	Groups pulumi.StringArrayOutput `pulumi:"groups"`
	// The name of the IAM User to add to groups
	User pulumi.StringOutput `pulumi:"user"`
}

// NewUserGroupMembership registers a new resource with the given unique name, arguments, and options.
func NewUserGroupMembership(ctx *pulumi.Context,
	name string, args *UserGroupMembershipArgs, opts ...pulumi.ResourceOption) (*UserGroupMembership, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Groups == nil {
		return nil, errors.New("invalid value for required argument 'Groups'")
	}
	if args.User == nil {
		return nil, errors.New("invalid value for required argument 'User'")
	}
	var resource UserGroupMembership
	err := ctx.RegisterResource("aws:iam/userGroupMembership:UserGroupMembership", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetUserGroupMembership gets an existing UserGroupMembership resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetUserGroupMembership(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *UserGroupMembershipState, opts ...pulumi.ResourceOption) (*UserGroupMembership, error) {
	var resource UserGroupMembership
	err := ctx.ReadResource("aws:iam/userGroupMembership:UserGroupMembership", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering UserGroupMembership resources.
type userGroupMembershipState struct {
	// A list of IAM Groups to add the user to
	Groups []string `pulumi:"groups"`
	// The name of the IAM User to add to groups
	User *string `pulumi:"user"`
}

type UserGroupMembershipState struct {
	// A list of IAM Groups to add the user to
	Groups pulumi.StringArrayInput
	// The name of the IAM User to add to groups
	User pulumi.StringPtrInput
}

func (UserGroupMembershipState) ElementType() reflect.Type {
	return reflect.TypeOf((*userGroupMembershipState)(nil)).Elem()
}

type userGroupMembershipArgs struct {
	// A list of IAM Groups to add the user to
	Groups []string `pulumi:"groups"`
	// The name of the IAM User to add to groups
	User string `pulumi:"user"`
}

// The set of arguments for constructing a UserGroupMembership resource.
type UserGroupMembershipArgs struct {
	// A list of IAM Groups to add the user to
	Groups pulumi.StringArrayInput
	// The name of the IAM User to add to groups
	User pulumi.StringInput
}

func (UserGroupMembershipArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*userGroupMembershipArgs)(nil)).Elem()
}

type UserGroupMembershipInput interface {
	pulumi.Input

	ToUserGroupMembershipOutput() UserGroupMembershipOutput
	ToUserGroupMembershipOutputWithContext(ctx context.Context) UserGroupMembershipOutput
}

func (*UserGroupMembership) ElementType() reflect.Type {
	return reflect.TypeOf((**UserGroupMembership)(nil)).Elem()
}

func (i *UserGroupMembership) ToUserGroupMembershipOutput() UserGroupMembershipOutput {
	return i.ToUserGroupMembershipOutputWithContext(context.Background())
}

func (i *UserGroupMembership) ToUserGroupMembershipOutputWithContext(ctx context.Context) UserGroupMembershipOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserGroupMembershipOutput)
}

// UserGroupMembershipArrayInput is an input type that accepts UserGroupMembershipArray and UserGroupMembershipArrayOutput values.
// You can construct a concrete instance of `UserGroupMembershipArrayInput` via:
//
//	UserGroupMembershipArray{ UserGroupMembershipArgs{...} }
type UserGroupMembershipArrayInput interface {
	pulumi.Input

	ToUserGroupMembershipArrayOutput() UserGroupMembershipArrayOutput
	ToUserGroupMembershipArrayOutputWithContext(context.Context) UserGroupMembershipArrayOutput
}

type UserGroupMembershipArray []UserGroupMembershipInput

func (UserGroupMembershipArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*UserGroupMembership)(nil)).Elem()
}

func (i UserGroupMembershipArray) ToUserGroupMembershipArrayOutput() UserGroupMembershipArrayOutput {
	return i.ToUserGroupMembershipArrayOutputWithContext(context.Background())
}

func (i UserGroupMembershipArray) ToUserGroupMembershipArrayOutputWithContext(ctx context.Context) UserGroupMembershipArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserGroupMembershipArrayOutput)
}

// UserGroupMembershipMapInput is an input type that accepts UserGroupMembershipMap and UserGroupMembershipMapOutput values.
// You can construct a concrete instance of `UserGroupMembershipMapInput` via:
//
//	UserGroupMembershipMap{ "key": UserGroupMembershipArgs{...} }
type UserGroupMembershipMapInput interface {
	pulumi.Input

	ToUserGroupMembershipMapOutput() UserGroupMembershipMapOutput
	ToUserGroupMembershipMapOutputWithContext(context.Context) UserGroupMembershipMapOutput
}

type UserGroupMembershipMap map[string]UserGroupMembershipInput

func (UserGroupMembershipMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*UserGroupMembership)(nil)).Elem()
}

func (i UserGroupMembershipMap) ToUserGroupMembershipMapOutput() UserGroupMembershipMapOutput {
	return i.ToUserGroupMembershipMapOutputWithContext(context.Background())
}

func (i UserGroupMembershipMap) ToUserGroupMembershipMapOutputWithContext(ctx context.Context) UserGroupMembershipMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserGroupMembershipMapOutput)
}

type UserGroupMembershipOutput struct{ *pulumi.OutputState }

func (UserGroupMembershipOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**UserGroupMembership)(nil)).Elem()
}

func (o UserGroupMembershipOutput) ToUserGroupMembershipOutput() UserGroupMembershipOutput {
	return o
}

func (o UserGroupMembershipOutput) ToUserGroupMembershipOutputWithContext(ctx context.Context) UserGroupMembershipOutput {
	return o
}

// A list of IAM Groups to add the user to
func (o UserGroupMembershipOutput) Groups() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *UserGroupMembership) pulumi.StringArrayOutput { return v.Groups }).(pulumi.StringArrayOutput)
}

// The name of the IAM User to add to groups
func (o UserGroupMembershipOutput) User() pulumi.StringOutput {
	return o.ApplyT(func(v *UserGroupMembership) pulumi.StringOutput { return v.User }).(pulumi.StringOutput)
}

type UserGroupMembershipArrayOutput struct{ *pulumi.OutputState }

func (UserGroupMembershipArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*UserGroupMembership)(nil)).Elem()
}

func (o UserGroupMembershipArrayOutput) ToUserGroupMembershipArrayOutput() UserGroupMembershipArrayOutput {
	return o
}

func (o UserGroupMembershipArrayOutput) ToUserGroupMembershipArrayOutputWithContext(ctx context.Context) UserGroupMembershipArrayOutput {
	return o
}

func (o UserGroupMembershipArrayOutput) Index(i pulumi.IntInput) UserGroupMembershipOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *UserGroupMembership {
		return vs[0].([]*UserGroupMembership)[vs[1].(int)]
	}).(UserGroupMembershipOutput)
}

type UserGroupMembershipMapOutput struct{ *pulumi.OutputState }

func (UserGroupMembershipMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*UserGroupMembership)(nil)).Elem()
}

func (o UserGroupMembershipMapOutput) ToUserGroupMembershipMapOutput() UserGroupMembershipMapOutput {
	return o
}

func (o UserGroupMembershipMapOutput) ToUserGroupMembershipMapOutputWithContext(ctx context.Context) UserGroupMembershipMapOutput {
	return o
}

func (o UserGroupMembershipMapOutput) MapIndex(k pulumi.StringInput) UserGroupMembershipOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *UserGroupMembership {
		return vs[0].(map[string]*UserGroupMembership)[vs[1].(string)]
	}).(UserGroupMembershipOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*UserGroupMembershipInput)(nil)).Elem(), &UserGroupMembership{})
	pulumi.RegisterInputType(reflect.TypeOf((*UserGroupMembershipArrayInput)(nil)).Elem(), UserGroupMembershipArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*UserGroupMembershipMapInput)(nil)).Elem(), UserGroupMembershipMap{})
	pulumi.RegisterOutputType(UserGroupMembershipOutput{})
	pulumi.RegisterOutputType(UserGroupMembershipArrayOutput{})
	pulumi.RegisterOutputType(UserGroupMembershipMapOutput{})
}