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
	cases := []fixerTestCase {
		{
			"\n" + `package main` + "\n\n" +
			`func main() {` + "\n\n" +
			`}`,
			`// Header` + "\n\n" +
			`package main` + "\n\n" +
			`func main() {` + "\n\n" +
			`}`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 1}, Text: "License header required"},
			},
		},
	}

	return cases
}

func createFileHeaderFixer() *fixers.FileHeaderCsFixer {
	mapFixer, _ := fixers.CreateFixer(
		"file_header",
		fixers.FixerOptions{"header": "// Header\n", "lintText": "License header required"},
	)

	return mapFixer.(*fixers.FileHeaderCsFixer)
}
