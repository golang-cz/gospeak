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
	schemaDir, _, err := collectCliArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	if schemaDir == "" {
		fmt.Fprintf(os.Stderr, "<schema> is required: try gospeak --help\n")
		os.Exit(1)
	}

	targets, err := gospeak.Parse(schemaDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse Go schema: %v\n", err)
		os.Exit(1)
	}

	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "no interfaces with //go:webrpc found, see https://github.com/golang-cz/gospeak\n")
		os.Exit(1)
	}

	for _, target := range targets {
		if target.Generator == "json" {
			jsonSchema, _ := target.Schema.ToJSON()
			if err := os.WriteFile(target.OutFile, []byte(jsonSchema), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write to %q file: %v\n", target.OutFile, err)
				os.Exit(1)
			}
			fmt.Printf("%20v => %v ✓\n", target.InterfaceName, target.OutFile)
			continue
		}

		config := &gen.Config{
			RefreshCache:    false,
			Format:          false,
			TemplateOptions: target.Opts,
		}

		generated, err := gen.Generate(target.Schema, target.Generator, config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(target.OutFile, []byte(generated.Code), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to %q file: %v\n", target.OutFile, err)
			os.Exit(1)
		}
		fmt.Printf("%20v => %v ✓\n", target.InterfaceName, target.OutFile)
	}
}

type Target struct {
	Name string
	Out  string
	Opts map[string]interface{}
}

// gospeak <schema.go> <target> [-targetOpts...] -out=<file> ... [<targetN> [-targetOpts] -out=<file>...]
func collectCliArgs(args []string) (schema string, targets []*Target, err error) {
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
func collectTargets(args []string) (targets []*Target, err error) {
	currentTarget := -1

	for _, arg := range args {
		name, value, _ := strings.Cut(arg, "=")

		// CLI flags or target options
		if strings.HasPrefix(name, "-") {
			name = strings.TrimLeft(name, "-")

			// target options
			if name == "out" {
				targets[currentTarget].Out = value
			} else {
				targets[currentTarget].Opts[name] = value
			}
		} else {
			currentTarget++
			targets = append(targets, &Target{
				Name: name,
				Opts: map[string]interface{}{},
			})
		}
	}

	for i, target := range targets {
		if target.Out == "" {
			return nil, fmt.Errorf("target[%v] %v must have -out=<path> flag", i, target.Name)
		}
	}

	return
}

const usage = `
Usage: gospeak <schema.go>
  -h, --help
        print this help
  -v, --version
        print gospeak version and exit

Finds all Go interfaces annotated with the special //go:webrpc target command comment.
Creates Webrpc schema from the Go interface.
Executes webrpc code generation for the given targets.

Example:

package api

//go:webrpc golang@v0.11.0 -client -out=../client.gen.go
//go:webrpc typescript@v0.11.0 -client -out=../client.gen.ts
type ExampleAPI interface {
	Ping(context.Context) (*Pong, error)
}

type Pong struct {
	Message string
}
`
