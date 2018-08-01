package fixers_test

import (
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func TestFileHeaderFixerLint(t *testing.T) {
	assertLint(t, createFileHeaderFixer(), fileHeaderTestTable())
}

func TestFileHeaderFix(t *testing.T) {
	assertFix(t, createFileHeaderFixer(), fileHeaderTestTable())
}

func fileHeaderTestTable() []fixerTestCase {
	return []fixerTestCase{
		{
			`package main

func main() {
}`,
			`// Header

package main

func main() {
}`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 1}, Text: "License header required"},
			},
		},
	}
}

func createFileHeaderFixer() *fixers.FileHeaderCsFixer {
	mapFixer, _ := fixers.CreateFixer(
		"file_header",
		fixers.FixerOptions{"header": "// Header\n\n", "lintText": "License header required"},
	)

	return mapFixer.(*fixers.FileHeaderCsFixer)
}
