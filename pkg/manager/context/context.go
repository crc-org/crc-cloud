package context

import (
	"github.com/crc/crc-cloud/pkg/util/maps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// store details for the current execution
type context struct {
	projectName string
	tags        pulumi.StringMap
}

var c context

func Init(projectName string, tags map[string]string) {
	c = context{
		projectName: projectName,
		tags: maps.Convert(tags,
			func(name string) string { return name },
			func(value string) pulumi.StringInput { return pulumi.String(value) }),
	}
}

func GetTags() pulumi.StringMap {
	return c.tags
}

func GetName() string {
	return c.projectName
}
