package main

import (
	v1 "github.com/thought-machine/dracon/api/proto/v1"
	types "github.com/thought-machine/dracon/producers/npm_audit/types/npmaudit-issue"

	"fmt"
	"flag"
	"log"
	"github.com/thought-machine/dracon/producers"
	"strings"
)

var (
	PackagePath string
)

func main() {
	flag.StringVar(&PackagePath, "package-path", "", "Path to the package.json file corresponding to this audit report; will be prepended to vulnerable dependency names in issue reports if specified")

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	var results types.NpmAuditOut
	if err := producers.ParseInFileJSON(&results); err != nil {
		log.Fatal(err)
	}

	if err := producers.WriteDraconOut(
		"npm-audit",
		parseOut(&results),
	); err != nil {
		log.Fatal(err)
	}
}

func parseOut(results *types.NpmAuditOut) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, res := range results.Advisories {
		issues = append(issues, parseResult(&res))
	}
	return issues
}
func parseResult(r *types.NpmAuditAdvisories) *v1.Issue {
	targetName := r.ModuleName
	if PackagePath != "" {
		targetName = PackagePath + ":" + r.ModuleName
	}

	return &v1.Issue{
		Target:   targetName,
		Type:     "Vulnerable Dependency",
		Title:    r.Title,
		Severity: v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", strings.ToUpper(r.Severity))]),
		// Cvss:        0.0,
		Confidence: v1.Confidence_CONFIDENCE_HIGH,
		Description: fmt.Sprintf("Vulnerable Versions: %s\nRecommendation: %s\nOverview: %s\nReferences: %s\nNPM Advisory URL: %s\n",
			r.VulnerableVersions, r.Recommendation, r.Overview, r.References, r.URL),
	}
}
