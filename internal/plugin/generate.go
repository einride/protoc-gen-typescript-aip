package plugin

import (
	"fmt"
	"path"
	"sort"

	"go.einride.tech/aip/reflect/aipreflect"
	"go.einride.tech/aip/reflect/aipregistry"
	"go.einride.tech/protoc-gen-typescript-aip/internal/codegen"
	"go.einride.tech/protoc-gen-typescript-aip/internal/plugin/resourcename"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	var opts Options
	if err := opts.Unmarshal(request.Parameter); err != nil {
		return nil, err
	}

	generate := make(map[string]struct{})
	for _, f := range request.FileToGenerate {
		generate[f] = struct{}{}
	}

	protoFiles, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: request.ProtoFile,
	})
	if err != nil {
		return nil, fmt.Errorf("create proto registry: %w", err)
	}

	aipRegistry, err := aipregistry.NewResources(protoFiles)
	if err != nil {
		return nil, fmt.Errorf("create aip registry: %w", err)
	}
	packageResources := make(map[string][]*aipreflect.ResourceDescriptor)
	aipRegistry.RangeResources(func(descriptor *aipreflect.ResourceDescriptor) bool {
		if _, ok := generate[descriptor.ParentFile]; !ok {
			return true
		}
		dir := path.Dir(descriptor.ParentFile)
		packageResources[dir] = append(packageResources[dir], descriptor)
		return true
	})

	var response pluginpb.CodeGeneratorResponse
	for pkg, resources := range packageResources {
		resources := resources
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Type.Type() < resources[j].Type.Type()
		})

		var file codegen.File
		if err := resourcename.GeneratePackage(&file, resources, aipRegistry); err != nil {
			return nil, fmt.Errorf("generate resource name: %w", err)
		}
		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:           proto.String(path.Join(pkg, opts.Filename)),
			Content:        proto.String(string(file.Content())),
			InsertionPoint: proto.String(opts.InsertionPoint),
		})
	}

	return &response, nil
}
