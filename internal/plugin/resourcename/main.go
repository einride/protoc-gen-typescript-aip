package resourcename

import (
	"fmt"
	"path"

	"github.com/einride/protoc-gen-typescript-aip/internal/codegen"
	"go.einride.tech/aip/reflect/aipreflect"
	"go.einride.tech/aip/reflect/aipregistry"
)

func GeneratePackage(
	f *codegen.File,
	resources []*aipreflect.ResourceDescriptor,
	registry *aipregistry.Resources,
) error {
	if len(resources) == 0 {
		return nil
	}

	pkg := path.Dir(resources[0].ParentFile)

	// keep track of what resource names have been generated in the package
	seen := make(map[aipreflect.ResourceTypeName]struct{})
	queue := resources

	for len(queue) > 0 {
		resource := queue[0]
		queue = queue[1:]

		if _, ok := seen[resource.Type]; ok {
			continue
		}
		seen[resource.Type] = struct{}{}

		// no support for multiple resource name definitions
		name := resource.Names[0]

		ancestors := make([]resourceName, 0, len(name.Ancestors))
		for _, ancestor := range name.Ancestors {
			ancestorResource, ok := registry.FindResourceByType(ancestor)
			if !ok {
				return fmt.Errorf("unable to find resource type: %s", ancestor)
			}
			_ = pkg
			ancestors = append(ancestors, resourceName{
				typeName: ancestor,
				pattern:  ancestorResource.Names[0].Pattern,
			})
			queue = append(queue, ancestorResource)
		}

		// resource declarations "imported" from other packages
		// should not be re-exported.
		shouldExport := path.Dir(resource.ParentFile) == pkg
		resourceName{
			export:    shouldExport,
			typeName:  name.Type,
			pattern:   name.Pattern,
			ancestors: ancestors,
		}.Generate(f)
	}
	return nil
}
