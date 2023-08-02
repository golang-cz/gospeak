package test

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func coloredDiff(x, y interface{}, opts ...cmp.Option) string {
	escapeCode := func(code int) string {
		return fmt.Sprintf("\x1b[%dm", code)
	}
	diff := cmp.Diff(x, y, opts...)
	if diff == "" {
		return ""
	}
	ss := strings.Split(diff, "\n")
	for i, s := range ss {
		switch {
		case strings.HasPrefix(s, "-"):
			ss[i] = escapeCode(31) + s + escapeCode(0)
		case strings.HasPrefix(s, "+"):
			ss[i] = escapeCode(32) + s + escapeCode(0)
		}
	}
	return strings.Join(ss, "\n")
}

func prefixLinesWithLineNumber(input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("%d: %s", i+1, lines[i])
	}
	return strings.Join(lines, "\n")
}
