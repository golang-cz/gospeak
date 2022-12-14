package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-cz/gospeak"
	"github.com/webrpc/webrpc/gen"
)

var (
	VERSION = "v0.0.x-dev"
)

func main() {
	schemaDir, targets, err := collectCliArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	if schemaDir == "" {
		fmt.Fprintf(os.Stderr, "<schema> is required: try gospeak --help\n")
		os.Exit(1)
	}

	schema, err := gospeak.Parse(schemaDir, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse Go schema: %v\n", err)
		os.Exit(1)
	}

	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "target is required: try gospeak --help\n")
		os.Exit(1)
	}

	for _, target := range targets {
		if target.name == "json" {
			jsonSchema, _ := schema.ToJSON(true)
			if err := os.WriteFile(target.out, []byte(jsonSchema), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write to %q file: %v\n", target.out, err)
				os.Exit(1)
			}
			fmt.Printf("generated %v ✓\n", target.out)
			continue
		}

		config := &gen.Config{
			RefreshCache:    false,
			Format:          true,
			TemplateOptions: target.opts,
		}

		generated, err := gen.Generate(schema, target.name, config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(target.out, []byte(generated.Code), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to %q file: %v\n", target.out, err)
			os.Exit(1)
		}
		fmt.Printf("generated %v ✓\n", target.out)
	}
}

type target struct {
	name string
	out  string
	opts map[string]interface{}
}

// gospeak <schema.go> <target> [-targetOpts...] -out=<file> ... [<targetN> [-targetOpts] -out=<file>...]
func collectCliArgs(args []string) (schema string, targets []*target, err error) {
	for i, arg := range args {
		// CLI flags or target options
		if strings.HasPrefix(arg, "-") {
			// CLI flags
			switch strings.TrimLeft(arg, "-") {
			case "h", "help":
				fmt.Fprintf(os.Stdout, usage)
				os.Exit(0)

			case "v", "version":
				fmt.Println("gospeak", VERSION)
				os.Exit(0)

			default:
				return "", nil, fmt.Errorf("unknown option %q", arg)
			}
		} else {
			if schema == "" {
				schema = arg
				continue
			}
			targets, err = collectTargets(args[i:])
			return
		}
	}

	return
}

// <target> [-targetOpts...] ... [<targetN> [-targetOpts]...]
func collectTargets(args []string) (targets []*target, err error) {
	currentTarget := -1

	for _, arg := range args {
		name, value, _ := strings.Cut(arg, "=")

		// CLI flags or target options
		if strings.HasPrefix(name, "-") {
			name = strings.TrimLeft(name, "-")

			// target options
			if name == "out" {
				targets[currentTarget].out = value
			} else {
				targets[currentTarget].opts[name] = value
			}
		} else {
			currentTarget++
			targets = append(targets, &target{
				name: name,
				opts: map[string]interface{}{},
			})
		}
	}

	for i, target := range targets {
		if target.out == "" {
			return nil, fmt.Errorf("target[%v] %v must have -out=<path> flag", i, target.name)
		}
	}

	return
}

const usage = `
Usage: gospeak <schema.go> <target> [-targetOpts...] -out=<file> ...
  -h, --help
        print this help
  -v, --version
        print gospeak version and exit

Targets:
  json
  golang       (see https://github.com/webrpc/gen-golang)
  typescript   (see https://github.com/webrpc/gen-typescript)
  javascript   (see https://github.com/webrpc/gen-javascript)
  openapi      (see https://github.com/webrpc/gen-openapi)

Example usage:
  gospeak path/to/api.go json -out=api.json

  gospeak ./schema/api.go                                  \
    json -out ./schema.json                                \
    golang -server -pkg server -out ./server/server.gen.go \
    golang -client -pkg client -out ./client/client.gen.go \
    typescript -client -out ../frontend/src/client.gen.ts  \
`
