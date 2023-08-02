package parser

import (
	"go/types"
	"regexp"
	"strings"
)

// This regex will return the following three submatches,
// given `db:"id,omitempty,pk" json:"id,string"` struct tag:
//
//	[0]: json:"id,string"
//	[1]: id
//	[2]: ,string
var jsonTagRegex, _ = regexp.Compile(`\s?json:\"([^,\"]*)(,[^\"]*)?\"`)

type JsonTag struct {
	Name      string
	Value     string
	IsString  bool
	Omitempty bool
}

func GetJsonTag(structTags string) (JsonTag, bool) {
	if !strings.Contains(structTags, `json:"`) {
		return JsonTag{}, false
	}

	submatches := jsonTagRegex.FindStringSubmatch(structTags)

	// Submatches from the jsonTagRegex:
	// [0]: json:"deleted_by,omitempty,string"
	// [1]: deleted_by
	// [2]: ,omitempty,string
	if len(submatches) != 3 {
		return JsonTag{}, false
	}

	jsonTag := JsonTag{
		Name:      submatches[1],
		Value:     submatches[1] + submatches[2],
		IsString:  strings.Contains(submatches[2], ",string"),
		Omitempty: strings.Contains(submatches[2], ",omitempty"),
	}

	return jsonTag, true
}

var textMarshalerRegex = regexp.MustCompile(`^func \((.+)\)\.MarshalText\(\) \((.+ )?\[\]byte, ([a-z]+ )?error\)$`)
var textUnmarshalerRegex = regexp.MustCompile(`^func \((.+)\)\.UnmarshalText\((.+ )?\[\]byte\) \(?(.+ )?error\)?$`)

// Returns true if the given type implements encoding.TextMarshaler/TextUnmarshaler interfaces.
func isTextMarshaler(typ types.Type, pkg *types.Package) bool {
	marshalTextMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "MarshalText")
	if marshalTextMethod == nil || !textMarshalerRegex.MatchString(marshalTextMethod.String()) {
		return false
	}

	unmarshalTextMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "UnmarshalText")
	if unmarshalTextMethod == nil || !textUnmarshalerRegex.MatchString(unmarshalTextMethod.String()) {
		return false
	}

	return true
}

var jsonMarshalerRegex = regexp.MustCompile(`^func \((.+)\)\.MarshalJSON\(\) \((.+ )?\[\]byte, ([a-z]+ )?error\)$`)
var jsonUnmarshalerRegex = regexp.MustCompile(`^func \((.+)\)\.UnmarshalJSON\((.+ )?\[\]byte\) \(?(.+ )?error\)?$`)

// Returns true if the given type implements json.Marshaler/Unmarshaler interfaces.
func isJsonMarshaller(typ types.Type, pkg *types.Package) bool {
	marshalJsonMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "MarshalJSON")
	if marshalJsonMethod == nil || !jsonMarshalerRegex.MatchString(marshalJsonMethod.String()) {
		return false
	}

	unmarshalJsonMethod, _, _ := types.LookupFieldOrMethod(typ, true, pkg, "UnmarshalJSON")
	if unmarshalJsonMethod == nil || !jsonUnmarshalerRegex.MatchString(unmarshalJsonMethod.String()) {
		return false
	}

	return true
}
