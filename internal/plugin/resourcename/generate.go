package resourcename

import (
	"fmt"
	"strconv"
	"strings"

	"go.einride.tech/aip/reflect/aipreflect"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/protoc-gen-typescript-aip/internal/codegen"
)

type resourceName struct {
	export    bool
	typeName  aipreflect.ResourceType
	pattern   string
	ancestors []resourceName
}

func (r resourceName) Generate(f *codegen.File) {
	r.generateInterface(f)
	r.generateConstructorInterface(f)
	r.generateConstructor(f)
}

func (r resourceName) generateInterface(f *codegen.File) {
	f.P("// Example: '", r.pattern, "'")
	f.P(r.maybeExport(), "interface ", interfaceName(r.typeName), " {")

	// segment accessors
	var sc resourcename.Scanner
	sc.Init(r.pattern)
	for sc.Scan() {
		seg := sc.Segment()
		if !seg.IsVariable() {
			continue
		}
		f.P(t(1), variableSegName(seg), ": string;")
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

	var segmentCount int
	var sc resourcename.Scanner
	sc.Init(r.pattern)
	for sc.Scan() {
		segmentCount++
	}

	isWrongLen := "segments.length !== " + strconv.Itoa(segmentCount)
	wrongLenErr := "`${errPrefix}: " +
		"invalid segment count ${segments.length} " +
		"(expected " + strconv.Itoa(segmentCount) + ")`"
	f.P(t(indent+1), "if (", isWrongLen, ") {")
	f.P(t(indent+2), "throw new Error(", wrongLenErr, ")")
	f.P(t(indent+1), "}")

	sc.Init(r.pattern)
	var i int
	for sc.Scan() {
		segment := sc.Segment()
		if segment.IsVariable() {
			f.P(t(indent+1), "const ", variableSegName(segment), " = segments[", i, "]")
		} else {
			isWrongConstSegment := "segments[" + strconv.Itoa(i) + "] !== " + strconv.Quote(string(segment))
			wrongConstSegmentErr := "`${errPrefix}: " +
				"invalid constant segment ${segments[" + strconv.Itoa(i) + "]} " +
				"(expected " + string(segment) + ")`"
			f.P(t(indent+1), "if (", isWrongConstSegment, ") {")
			f.P(t(indent+2), "throw new Error(", wrongConstSegmentErr, ")")
			f.P(t(indent+1), "}")
		}
		i++
	}
	f.P(t(indent+1), "return this.from(", fromArgs(r.pattern), ")")
	f.P(t(indent), "},")
}

func (r resourceName) generateConstructorFrom(f *codegen.File, indent int) {
	var sc resourcename.Scanner
	sc.Init(r.pattern)
	f.P(t(indent+0), "from(", fromArgsDecl(r.pattern), "): ", interfaceName(r.typeName), " {")
	for sc.Scan() {
		seg := sc.Segment()
		if !seg.IsVariable() {
			continue
		}
		isEmpty := variableSegName(seg) + " === \"\""
		containsSlash := variableSegName(seg) + ".indexOf(\"/\") > -1"
		invalidVariableSegmentErr := "`invalid variable segment for " + string(seg.Literal()) + ": " +
			"'${" + variableSegName(seg) + "}'`"
		f.P(t(indent+1), "if (", isEmpty, " || ", containsSlash, ") {")
		f.P(t(indent+2), "throw new Error(", invalidVariableSegmentErr, ")")
		f.P(t(indent+1), "}")
	}
	f.P(t(indent+1), "return {")
	// segment accessors
	sc.Init(r.pattern)
	for sc.Scan() {
		seg := sc.Segment()
		if !seg.IsVariable() {
			continue
		}
		f.P(t(indent+2), variableSegName(seg), ",")
	}
	// ancestor accessors
	for _, ancestor := range r.ancestors {
		f.P(t(indent+2), ancestorAccessorName(ancestor.typeName), "(): ", interfaceName(ancestor.typeName), " {")
		f.P(t(indent+3), "return ", interfaceName(ancestor.typeName), ".from(", fromArgs(ancestor.pattern), ")")
		f.P(t(indent+2), "},")
	}
	// toString
	var toStringParts []string
	sc.Init(r.pattern)
	for sc.Scan() {
		seg := sc.Segment()
		if seg.IsVariable() {
			toStringParts = append(toStringParts, variableSegName(seg))
		} else {
			toStringParts = append(toStringParts, strconv.Quote(string(seg)))
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

func fromArgsDecl(pattern string) string {
	args := make([]string, 0, 8)
	var sc resourcename.Scanner
	sc.Init(pattern)
	for sc.Scan() {
		seg := sc.Segment()
		if !seg.IsVariable() {
			continue
		}
		args = append(args, fmt.Sprintf("%s: string", variableSegName(seg)))
	}
	return strings.Join(args, ", ")
}

func fromArgs(pattern string) string {
	args := make([]string, 0, 8)
	var sc resourcename.Scanner
	sc.Init(pattern)
	for sc.Scan() {
		seg := sc.Segment()
		if !seg.IsVariable() {
			continue
		}
		args = append(args, variableSegName(seg))
	}
	return strings.Join(args, ", ")
}
