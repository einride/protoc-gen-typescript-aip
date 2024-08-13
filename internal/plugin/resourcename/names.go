package resourcename

import (
	"strings"

	"github.com/stoewer/go-strcase"
	"go.einride.tech/aip/reflect/aipreflect"
	"go.einride.tech/aip/resourcename"
)

func interfaceName(resource aipreflect.ResourceType) string {
	return resource.Type() + "ResourceName"
}

func ancestorAccessorName(resource aipreflect.ResourceType) string {
	ifi := interfaceName(resource)
	return strings.ToLower(string(ifi[0])) + ifi[1:]
}

func constructorName(resource aipreflect.ResourceType) string {
	return interfaceName(resource) + "Constructor"
}

func variableSegName(seg resourcename.Segment) string {
	return strcase.LowerCamelCase(string(seg.Literal()))
}
