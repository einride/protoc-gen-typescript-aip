package resourcename

import (
	"strings"

	"github.com/stoewer/go-strcase"
	"go.einride.tech/aip/reflect/aipreflect"
)

func interfaceName(resource aipreflect.ResourceTypeName) string {
	return resource.Type() + "ResourceName"
}

func ancestorAccessorName(resource aipreflect.ResourceTypeName) string {
	ifi := interfaceName(resource)
	return strings.ToLower(string(ifi[0])) + ifi[1:]
}

func constructorName(resource aipreflect.ResourceTypeName) string {
	return interfaceName(resource) + "Constructor"
}

func variableSegName(value string) string {
	return strcase.LowerCamelCase(value)
}
