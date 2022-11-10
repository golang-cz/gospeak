package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-cz/go2webrpc"
)

var flags = flag.NewFlagSet("go2webrpc", flag.ExitOnError)

func main() {
	versionFlag := flags.Bool("version", false, "print go2webrpc version and exit")
	schemaFlag := flags.String("schema", "", "Path to Go file / package folder (required)")
	interfaceFlag := flags.String("interface", "", "Name of the interface (required)")
	outFlag := flags.String("out", "", "generated output file (optional, default: stdout)")

	flags.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Printf("go2webrpc %s\n", go2webrpc.VERSION)
		os.Exit(0)
	}

	if *schemaFlag == "" {
		fmt.Fprintln(os.Stderr, "-schema flag is required")
		os.Exit(1)
	}

	if *interfaceFlag == "" {
		fmt.Fprintln(os.Stderr, "-interface flag is required")
		os.Exit(1)
	}

	schema, err := go2webrpc.Parse(*schemaFlag, *interfaceFlag)
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
