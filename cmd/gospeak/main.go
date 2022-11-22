package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-cz/gospeak"
)

var (
	VERSION = "v0.0.x-dev"
	flags   = flag.NewFlagSet("gospeak", flag.ExitOnError)
)

func main() {
	versionFlag := flags.Bool("version", false, "print gospeak version and exit")
	schemaFlag := flags.String("schema", "", "path to Go package (required)")
	interfaceFlag := flags.String("interface", "", "Go interface name (required)")
	outFlag := flags.String("out", "", "generated output file (optional, default: stdout)")

	flags.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Printf("gospeak %s\n", VERSION)
		os.Exit(0)
	}

	if *schemaFlag == "" {
		fmt.Fprintln(os.Stderr, "-schema flag is required")
		os.Exit(1)
	}

	schema, err := gospeak.Parse(*schemaFlag, *interfaceFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse Go schema: %v\n", err)
		os.Exit(1)
	}

	json, err := schema.ToJSON(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to print webrpc schema as JSON: %v\n", err)
		os.Exit(1)
	}

	if *outFlag == "" {
		fmt.Println(json)
		os.Exit(0)
	}

	if err := os.WriteFile(*outFlag, []byte(json), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to %q file: %v\n", *outFlag, err)
		os.Exit(1)
	}
}
