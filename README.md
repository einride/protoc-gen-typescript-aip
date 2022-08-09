protoc-gen-typescript-aip
=========================

Generates Typescript helpers for Protobuf APIs conforming to [AIP](https://aip.dev).

### Install the plugin

```bash
go get go.einride.tech/protoc-gen-typescript-aip
```

Or download a prebuilt binary from [releases](https://github.com/einride/protoc-gen-typescript-aip/releases).

### Invocation

#### protoc

```bash
protoc 
  --typescript-aip_out [OUTPUT DIR] \
  [.proto files ...]
```

#### buf

```yaml
plugins:
  - name: typescript-aip
    out: [OUTPUT DIR]
    strategy: all
```

### Configuration

```
filename                Name of the file to generate the code to.
                        Default: `index.ts`.

insertion_point         If non-empty, indicates that the named file should already exist,
                        and the content here is to be inserted into that file at a defined 
                        insertion point. 

skip_file_resource_definitions
                        If set to true, resource names will not be generated for resource definitions
                        in the file scope. Default: false.
```

---

Features
--------

### Resource names

Generates helpers for working with resource names, based on [ResourceDescriptor](https://github.com/googleapis/googleapis/blob/master/google/api/resource.proto) annotations.

#### Example

A resource annotated with

```proto
option (google.api.resource) = {
  type: "account-example.einride.tech/User"
  pattern: "tenants/{tenant}/users/{user}"
  singular: "user"
  plural: "users"
};
```

generates the following API

```ts
// Parsing a string:
const name = UserResourceName.parse("tenants/1/users/2")

// Getting variable segments:
console.log(name.tenant);     // "1"
console.log(name.user);       // "2"

// Getting the string back:
console.log(name.toString())  // "tenants/1/users/2"

// Traversing the resource hierarchy:
const tenant = name.tenantResourceName();
console.log(tenant.toString())  // "tenants/1"

// Constructing the resource name from segments:
const name = UserResourceName.from("tenant", "user")
console.log(name.tenant)        // "tenant"
console.log(name.user)          // "user"
console.log(name.toString())    // "tenants/tenant/users/user"
```
