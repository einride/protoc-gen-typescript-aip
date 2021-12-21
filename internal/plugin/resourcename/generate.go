package resourcename

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/einride/protoc-gen-typescript-aip/internal/codegen"
	"go.einride.tech/aip/reflect/aipreflect"
)

type resourceName struct {
	export    bool
	typeName  aipreflect.ResourceTypeName
	pattern   aipreflect.ResourceNamePatternDescriptor
	ancestors []resourceName
}

func (r resourceName) Generate(f *codegen.File) {
	r.generateInterface(f)
	r.generateConstructorInterface(f)
	r.generateConstructor(f)
}

func (r resourceName) generateInterface(f *codegen.File) {
	f.P("// Example: '", r.pattern.String(), "'")
	f.P(r.maybeExport(), "interface ", interfaceName(r.typeName), " {")

	// segment accessors
	for _, seg := range r.pattern.Segments {
		if !seg.Variable {
			continue
		}
		f.P(t(1), variableSegName(seg.Value), ": string;")
	}

	// ancestor accessors
	for _, ancestor := range r.ancestors {
		f.P(t(1), ancestorAccessorName(ancestor.typeName), "(): ", interfaceName(ancestor.typeName), ";")
	}

	f.P(t(1), "toString(): string;")
	f.P("}")
	f.P()
}

func (r resourceName) generateConstructorInterface(f *codegen.File) {
	f.P("interface ", constructorName(r.typeName), " {")
	f.P(t(1), "parse(s: string): ", interfaceName(r.typeName), ";")
	f.P(t(1), "from(", fromArgsDecl(r.pattern), "): ", interfaceName(r.typeName), ";")
	f.P("}")
	f.P()
}

func (r resourceName) generateConstructor(f *codegen.File) {
	f.P(r.maybeExport(), "const ", interfaceName(r.typeName), ": ", constructorName(r.typeName), " = {")
	r.generateConstructorParse(f, 1)
	f.P()
	r.generateConstructorFrom(f, 1)
	f.P("}")
	f.P()
}

func (r resourceName) generateConstructorParse(f *codegen.File, indent int) {
	f.P(t(indent), "parse(s: string): ", interfaceName(r.typeName), " {")
	f.P(t(indent+1), "const errPrefix = `parse resource name ${s} as ", r.typeName, "`;")
	f.P(t(indent+1), "const segments = s.split(\"/\")")

	isWrongLen := "segments.length !== " + strconv.Itoa(len(r.pattern.Segments))
	wrongLenErr := "`${errPrefix}: " +
		"invalid segment count ${segments.length} " +
		"(expected " + strconv.Itoa(len(r.pattern.Segments)) + ")`"
	f.P(t(indent+1), "if (", isWrongLen, ") {")
	f.P(t(indent+2), "throw new Error(", wrongLenErr, ")")
	f.P(t(indent+1), "}")

	for i, segment := range r.pattern.Segments {
		if segment.Variable {
			f.P(t(indent+1), "const ", variableSegName(segment.Value), " = segments[", i, "]")
		} else {
			isWrongConstSegment := "segments[" + strconv.Itoa(i) + "] !== " + strconv.Quote(segment.Value)
			wrongConstSegmentErr := "`${errPrefix}: " +
				"invalid constant segment ${segments[" + strconv.Itoa(i) + "]} " +
				"(expected " + segment.Value + ")`"
			f.P(t(indent+1), "if (", isWrongConstSegment, ") {")
			f.P(t(indent+2), "throw new Error(", wrongConstSegmentErr, ")")
			f.P(t(indent+1), "}")
		}
	}
	f.P(t(indent+1), "return this.from(", fromArgs(r.pattern), ")")
	f.P(t(indent), "},")
}

func (r resourceName) generateConstructorFrom(f *codegen.File, indent int) {
	f.P(t(indent+0), "from(", fromArgsDecl(r.pattern), "): ", interfaceName(r.typeName), " {")
	for _, seg := range r.pattern.Segments {
		if !seg.Variable {
			continue
		}
		isEmpty := variableSegName(seg.Value) + " === \"\""
		containsSlash := variableSegName(seg.Value) + ".indexOf(\"/\") > -1"
		invalidVariableSegmentErr := "`invalid variable segment for " + seg.Value + ": " +
			"'${" + variableSegName(seg.Value) + "}'`"
		f.P(t(indent+1), "if (", isEmpty, " || ", containsSlash, ") {")
		f.P(t(indent+2), "throw new Error(", invalidVariableSegmentErr, ")")
		f.P(t(indent+1), "}")
	}
	f.P(t(indent+1), "return {")
	// segment accessors
	for _, seg := range r.pattern.Segments {
		if !seg.Variable {
			continue
		}
		f.P(t(indent+2), variableSegName(seg.Value), ",")
	}
	// ancestor accessors
	for _, ancestor := range r.ancestors {
		f.P(t(indent+2), ancestorAccessorName(ancestor.typeName), "(): ", interfaceName(ancestor.typeName), " {")
		f.P(t(indent+3), "return ", interfaceName(ancestor.typeName), ".from(", fromArgs(ancestor.pattern), ")")
		f.P(t(indent+2), "},")
	}
	// toString
	var toStringParts []string
	for _, seg := range r.pattern.Segments {
		if seg.Variable {
			toStringParts = append(toStringParts, variableSegName(seg.Value))
		} else {
			toStringParts = append(toStringParts, strconv.Quote(seg.Value))
		}
	}
	f.P(t(indent+2), "toString(): string {")
	f.P(t(indent+3), "// eslint-disable-next-line no-useless-concat, prefer-template")
	f.P(t(indent+3), "return ", strings.Join(toStringParts, " + \"/\" + "))
	f.P(t(indent+2), "},")
	f.P(t(indent+1), "}")
	f.P(t(indent+0), "},")
}

func (r resourceName) maybeExport() string {
	if r.export {
		return "export "
	}
	return ""
}

func fromArgsDecl(pattern aipreflect.ResourceNamePatternDescriptor) string {
	args := make([]string, 0, len(pattern.Segments))
	for _, seg := range pattern.Segments {
		if !seg.Variable {
			continue
		}
		args = append(args, fmt.Sprintf("%s: string", variableSegName(seg.Value)))
	}
	return strings.Join(args, ", ")
}

func fromArgs(pattern aipreflect.ResourceNamePatternDescriptor) string {
	args := make([]string, 0, len(pattern.Segments))
	for _, seg := range pattern.Segments {
		if !seg.Variable {
			continue
		}
		args = append(args, variableSegName(seg.Value))
	}
	return strings.Join(args, ", ")
}
