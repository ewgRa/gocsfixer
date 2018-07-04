package fixers

import (
	"testing"
	"fmt"
)

func TestFileHeaderFixerLint(t *testing.T) {
	expected := []struct {
		line int
		text string
	}{
		{1, "License header required"},
	}

	fixer := &FileHeaderCsFixer{header: "// Header\n", lintText: "License header required"}

	problems, _ := fixer.Lint(contentForFileHeaderCsFixer())

	if (len(problems) != len(expected)) {
		fmt.Println("Expected one problem, got", len(problems))
		t.Fail()
		return
	}

	for k, exp := range expected {
		if problems[k].Position.Line != exp.line {
			fmt.Println("Problem", k,"found on", problems[k].Position.Line, "line, expected", exp.line)
			t.Fail()
			return
		}

		if problems[k].Text != exp.text {
			fmt.Println("Problem", k, "have text", fmt.Sprintf("'%s'", problems[k].Text), ",", fmt.Sprintf("'%s'", exp.text), "expected")
			t.Fail()
			return
		}
	}
}

func TestFileHeaderFix(t *testing.T) {
	fixer := &FileHeaderCsFixer{header: "// Header\n"}

	contentFix, err := fixer.Fix(contentForFileHeaderCsFixer())

	if nil != err {
		fmt.Println("Error when perform fix:", err)
		t.Fail()
		return
	}

	expectedFixedContent := fixedContentForFileHeaderCsFixer()

	if expectedFixedContent != contentFix {
		fmt.Println("Fixed content differ from expected")
		fmt.Println(contentFix)
		t.Fail()
		return
	}

	problems, _ := fixer.Lint(expectedFixedContent)

	if len(problems) != 0 {
		fmt.Println("Expected no problem, got", len(problems))
		t.Fail()
		return
	}
}

func contentForFileHeaderCsFixer() string {
	return `
		package main

		func main() {
		}
	`
}

func fixedContentForFileHeaderCsFixer() string {
	return `// Header

		package main

		func main() {
		}
	`
}