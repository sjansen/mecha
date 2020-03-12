package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

// https://github.com/google/cel-spec

func main() {
	env, err := cel.NewEnv(cel.Declarations(
		decls.NewIdent("name", decls.String, nil),
		decls.NewIdent("group", decls.String, nil),
	))
	if err != nil {
		log.Fatalf("env construction error: %s", err)
	}

	parsed, issues := env.Parse(`name.startsWith("/groups/" + group)`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("parse error: %s", issues.Err())
	}

	checked, issues := env.Check(parsed)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	prg, err := env.Program(checked)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	out, details, err := prg.Eval(map[string]interface{}{
		"name":  "/groups/acme.co/documents/secret-stuff",
		"group": "acme.co",
	})
	if err != nil {
		log.Fatalf("program evaluation error: %s", err)
	}
	fmt.Println(out)
	fmt.Printf("%#v\n", details)
}
