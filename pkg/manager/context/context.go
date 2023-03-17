package context

import (
	"github.com/crc/crc-cloud/pkg/util/maps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// store details for the current execution
type context struct {
	tags pulumi.StringMap
}

var c context

func Init(tags map[string]string) {
	c = context{
		tags: maps.Convert(tags,
			func(name string) string { return name },
			func(value string) pulumi.StringInput { return pulumi.String(value) }),
	}
}

func GetTags() pulumi.StringMap {
	return c.tags
}
