package plugin

import (
	"cmp"
	"fmt"
	"path"
	"path/filepath"
	"slices"

	"go.einride.tech/protoc-gen-typescript-aip/internal/codegen"
	"go.einride.tech/protoc-gen-typescript-aip/internal/plugin/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	var opts Options
	if err := opts.Unmarshal(request.GetParameter()); err != nil {
		return nil, err
	}

	generate := make(map[string]struct{})
	for _, f := range request.GetFileToGenerate() {
		generate[f] = struct{}{}
	}

	protoFiles, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: request.GetProtoFile(),
	})
	if err != nil {
		return nil, fmt.Errorf("create proto registry: %w", err)
	}
	packages := make([]protoreflect.FullName, 0, 8)
	packageDirs := make(map[protoreflect.FullName]string)
	seenPackages := make(map[protoreflect.FullName]struct{})
	protoFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		if _, ok := generate[file.Path()]; !ok {
			return true
		}
		pkg := file.Package()
		if _, ok := seenPackages[pkg]; ok {
			return true
		}
		packages = append(packages, pkg)
		packageDirs[pkg] = filepath.Dir(file.Path())
		seenPackages[pkg] = struct{}{}
		return true
	})

	var response pluginpb.CodeGeneratorResponse
	for _, pkg := range packages {
		descriptors := make([]*annotations.ResourceDescriptor, 0, 8)
		rangeResourceDescriptorsInPackage(protoFiles, pkg, func(resource *annotations.ResourceDescriptor) bool {
			descriptors = append(descriptors, resource)
			return true
		})
		slices.SortFunc(descriptors, func(a, b *annotations.ResourceDescriptor) int {
			return cmp.Compare(a.GetType(), b.GetType())
		})

		if len(descriptors) == 0 {
			continue
		}

		var file codegen.File
		if err := resourcename.GeneratePackage(&file, pkg, descriptors, protoFiles); err != nil {
			return nil, fmt.Errorf("generate resource name: %w", err)
		}
		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:           proto.String(path.Join(packageDirs[pkg], opts.Filename)),
			Content:        proto.String(string(file.Content())),
			InsertionPoint: proto.String(opts.InsertionPoint),
		})
	}

	return &response, nil
}

// RangeResourceDescriptorsInPackage iterates over all resource descriptors in a package while fn returns true.
// The provided registry is used for looking up files in the package.
// The iteration order is undefined.
func rangeResourceDescriptorsInPackage(
	registry *protoregistry.Files,
	packageName protoreflect.FullName,
	fn func(resource *annotations.ResourceDescriptor) bool,
) {
	registry.RangeFilesByPackage(packageName, func(file protoreflect.FileDescriptor) bool {
		for i := 0; i < file.Messages().Len(); i++ {
			resource := proto.GetExtension(
				file.Messages().Get(i).Options(), annotations.E_Resource,
			).(*annotations.ResourceDescriptor)
			if resource == nil {
				continue
			}
			if !fn(resource) {
				return false
			}
		}
		return true
	})
}
