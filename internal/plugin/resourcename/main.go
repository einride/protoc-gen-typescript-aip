package resourcename

import (
	"cmp"
	"slices"

	"go.einride.tech/aip/reflect/aipreflect"
	"go.einride.tech/protoc-gen-typescript-aip/internal/codegen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func GeneratePackage(
	f *codegen.File,
	pkg protoreflect.FullName,
	resources []*annotations.ResourceDescriptor,
	files *protoregistry.Files,
) error {
	if len(resources) == 0 {
		return nil
	}

	packageResources := make(map[string]struct{})
	for _, resource := range resources {
		packageResources[resource.GetType()] = struct{}{}
	}

	// keep track of what resource names have been generated in the package
	seen := make(map[string]struct{})
	queue := resources

	for len(queue) > 0 {
		resource := queue[0]
		queue = queue[1:]

		if _, ok := seen[resource.GetType()]; ok {
			continue
		}
		seen[resource.GetType()] = struct{}{}

		parents := make([]*annotations.ResourceDescriptor, 0, 16)
		aipreflect.RangeParentResourcesInPackage(
			files,
			pkg,
			resource.GetPattern()[0],
			func(parent *annotations.ResourceDescriptor) bool {
				parents = append(parents, parent)
				return true
			},
		)
		// iteration order of RangeParentResourcesInPackage is undefined.
		slices.SortFunc(parents, func(a, b *annotations.ResourceDescriptor) int {
			return cmp.Compare(b.GetType(), a.GetType())
		})
		ancestors := make([]resourceName, 0, 16)
		for _, parent := range parents {
			queue = append(queue, parent)
			ancestors = append(ancestors, resourceName{
				typeName: aipreflect.ResourceType(parent.GetType()),
				pattern:  parent.GetPattern()[0],
			})
		}

		// resource declarations "imported" from other packages
		// should not be re-exported.
		_, shouldExport := packageResources[resource.GetType()]
		resourceName{
			export:    shouldExport,
			typeName:  aipreflect.ResourceType(resource.GetType()),
			pattern:   resource.GetPattern()[0],
			ancestors: ancestors,
		}.Generate(f)
	}
	return nil
}
