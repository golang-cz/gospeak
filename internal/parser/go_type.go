package parser

import (
	"go/types"
	"path/filepath"
	"strings"
	"unicode"
)

func (p *Parser) GoTypeName(typ types.Type) string {
	name := typ.String() // []*github.com/golang-cz/gospeak/Pkg.Typ

	firstLetter := findFirstLetter(name)
	prefix := name[:firstLetter] // []*
	name = name[firstLetter:]    // github.com/golang-cz/gospeak/Pkg.Typ

	name = strings.TrimPrefix(name, p.SchemaPkgName+".")       // Typ (ignore root Pkg)
	name = strings.TrimPrefix(name, "command-line-arguments.") // Typ (ignore "command-line-arguments" Pkg autogenerated by Go tool chain)
	name = filepath.Base(name)                                 // Pkg.Typ

	if name == "invalid type" {
		name = "invalidType"
	}

	return prefix + name // []*Pkg.Typ
}

func (p *Parser) GoTypeImport(typ types.Type) string {
	name := typ.String() // []*github.com/golang-cz/gospeak/Pkg.Typ

	firstLetter := findFirstLetter(name)
	name = name[firstLetter:] // github.com/golang-cz/gospeak/Pkg.Typ

	lastDot := strings.LastIndex(name, ".")
	if lastDot <= 0 {
		return ""
	}

	name = name[:lastDot] // github.com/golang-cz/gospeak/Pkg
	switch name {
	case p.SchemaPkgName, "command-line-arguments", "time":
		return ""
	}

	return name
}

func findFirstLetter(s string) int {
	for i, char := range s {
		if unicode.IsLetter(char) {
			return i
		}
	}
	return 0
}
