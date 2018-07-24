package fixers_test

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"testing"
)

func TestAlternativeCallFixerLint(t *testing.T) {
	assertLint(t, createAlternativeCallFixer(), alternativeCallTestTable())
}

func TestAlternativeCallFix(t *testing.T) {
	assertFix(t, createAlternativeCallFixer(), alternativeCallTestTable())
}

func alternativeCallTestTable() []fixerTestCase {
	cases := []fixerTestCase{
		{
			`package main

import "logp"

func main() {
	logp.Warn("foo")
}
`,
			`package main

import "logp"

func main() {
	logp.Err("foo")
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 6}, Text: "Instead logp.Warn use alternative call logp.Err"},
			},
		},
	}

	return cases
}

func createAlternativeCallFixer() *fixers.AlternativeCallCsFixer {
	mapFixer, _ := fixers.CreateFixer(
		"alternative_call",
		fixers.FixerOptions{
			"selector":    "logp.Warn",
			"alternative": "logp.Err",
		},
	)

	return mapFixer.(*fixers.AlternativeCallCsFixer)
}
