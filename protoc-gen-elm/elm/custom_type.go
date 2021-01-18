package elm

import (
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/types/descriptorpb"
)

// CustomType - defines an Elm custom type (sometimes called union type)
// https://guide.elm-lang.org/types/custom_types.html
type CustomType struct {
	Name                   VariableName
	Decoder                VariableName
	Encoder                VariableName
	DefaultVariantVariable VariableName
	DefaultVariantValue    VariantName
	Variants               []CustomTypeVariant
}

// VariantName - unique camelcase identifier used for custom type variants
// https://guide.elm-lang.org/types/custom_types.html
type VariantName string

// VariantJSONName - unique JSON identifier, uppercase snake case, for a custom type variant
type VariantJSONName string

// CustomTypeVariant - a possible variant of a CustomType
// https://guide.elm-lang.org/types/custom_types.html
type CustomTypeVariant struct {
	Name     VariantName
	Number   ProtobufFieldNumber
	JSONName VariantJSONName
}

// NestedVariableName - top level Elm variable name for a possibly nested PB definition
func NestedVariableName(name string, preface []string) VariableName {
	fullName := name
	for _, p := range preface {
		fullName = fmt.Sprintf("%s_%s", p, fullName)
	}

	return VariableName(fullName)
}

// NestedVariantName - Elm variant name for a possibly nested PB definition
func NestedVariantName(name string, preface []string) VariantName {
	fullName := camelCase(strings.ToLower(name))
	for _, p := range preface {
		fullName = fmt.Sprintf("%s_%s", camelCase(strings.ToLower(p)), fullName)
	}

	return VariantName(fullName)
}

// DefaultVariantVariableName - convenient identifier for a custom types default variant
func DefaultVariantVariableName(name string, preface []string) VariableName {
	variableName := NestedVariableName(name, preface)
	return VariableName(firstLower(fmt.Sprintf("%sDefault", variableName)))
}

// EnumVariantJSONName - JSON identifier for variant decoder/encoding
func EnumVariantJSONName(pb *descriptorpb.EnumValueDescriptorProto) VariantJSONName {
	return VariantJSONName(pb.GetName())
}

// CustomTypeTemplate - defines templates for custom types
// For legacy code the definitions are split - reorganizing to reduce complexity is planned
func CustomTypeTemplate(t *template.Template) (*template.Template, error) {
	return t.Parse(`
{{- define "custom-type-definition" }}


type {{ .Name }}
{{- range $i, $v := .Variants }}
    {{ if not $i }}={{ else }}|{{ end }} {{ $v.Name }} -- {{ $v.Number }}
{{- end }}
{{- end }}
{{- define "custom-type-decoder" }}


{{ .Decoder }} : JD.Decoder {{ .Name }}
{{ .Decoder }} =
    let
        lookup s =
            case s of
{{- range .Variants }}
                "{{ .JSONName }}" ->
                    {{ .Name }}
{{ end }}
                _ ->
                    {{ .DefaultVariantValue }}
    in
        JD.map lookup JD.string


{{ .DefaultVariantVariable }} : {{ .Name }}
{{ .DefaultVariantVariable }} = {{ .DefaultVariantValue }}
{{- end }}
{{- define "custom-type-encoder" }}


{{ .Encoder }} : {{ .Name }} -> JE.Value
{{ .Encoder }} v =
    let
        lookup s =
            case s of
{{- range .Variants }}
                {{ .Name }} ->
                    "{{ .JSONName }}"
{{ end }}
    in
        JE.string <| lookup v
{{- end }}
{{- define "custom-type" }}
{{- template "custom-type-definition" . }}
{{- template "custom-type-decoder" . }}
{{- template "custom-type-encoder" . }}
{{- end }}
`)
}
