package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envlens/internal/comparator"
	"github.com/yourorg/envlens/internal/loader"
	"github.com/yourorg/envlens/internal/reporter"
)

func main() {
	var (
		refFile  = flag.String("ref", ".env.example", "Reference .env file (source of truth)")
		targets  = flag.String("targets", "", "Comma-separated list of target .env files to check")
		format   = flag.String("format", "text", "Output format: text or json")
		failFast = flag.Bool("fail", false	, "Exit with non-zero status if issues are found")
	)
	flag.Parse()

	if *targets == "" {
		fmt.Fprintln(os.Stderr, "error: --targets is required")
		flag.Usage()
		os.Exit(2)
	}

	targetFiles := splitAndTrim(*targets)

	ref, err := loader.Load(*refFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading reference file %q: %v\n", *refFile, err)
		os.Exit(2)
	}

	loadedTargets, err := loader.LoadAll(targetFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading target files: %v\n", err)
		os.Exit(2)
	}

	reports := comparator.Compare(ref, loadedTargets)

	if err := reporter.Write(os.Stdout, reports, *format); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(2)
	}

	if *failFast && hasIssues(reports) {
		os.Exit(1)
	}
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func hasIssues(reports []comparator.Report) bool {
	for _, r := range reports {
		if len(r.Missing) > 0 || len(r.Extra) > 0 {
			return true
		}
	}
	return false
}
