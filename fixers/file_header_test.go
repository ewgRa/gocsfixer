package fixers_test

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"testing"
)

func TestFileHeaderFixerLint(t *testing.T) {
	assertLint(t, createFileHeaderFixer(), fileHeaderTestTable())
}

func TestFileHeaderFix(t *testing.T) {
	assertFix(t, createFileHeaderFixer(), fileHeaderTestTable())
}

func fileHeaderTestTable() []fixerTestCase {
	return []fixerTestCase{
		{`
package main

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
		fixers.FixerOptions{"header": "// Header\n", "lintText": "License header required"},
	)

	return mapFixer.(*fixers.FileHeaderCsFixer)
}
