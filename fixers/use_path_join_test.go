package fixers

import (
	"testing"
	"fmt"
)

func TestUsePathJoinCsFixerLint(t *testing.T) {
	expected := []struct {
		line int
		text string
	}{
		{5, "Use path.Join"},
		{6, "Use path.Join"},
		{8, "Use path.Join"},
	}

	fixer := &UsePathJoinCsFixer{}

	problems, _ := fixer.Lint(contentForUsePathJoinCsFixer())

	if (len(problems) != len(expected)) {
		fmt.Println("Expected three problems, got", len(problems))
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

func TestUsePathJoinFix(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	contentFix, err := fixer.Fix(contentForUsePathJoinCsFixer())

	if nil != err {
		fmt.Println("Error when perform fix:", err)
		t.Fail()
		return
	}

	expectedFixedContent := fixedContentForUsePathJoinCsFixer()

	if expectedFixedContent != contentFix {
		fmt.Println("Fixed content differ from expected")
		fmt.Println(contentFix)
		t.Fail()
		return
	}
}

func contentForUsePathJoinCsFixer() string {
	return `
		package main

		func main() {
			os.Readlink(gosigar.Procd + "self")
			os.Readlink("foo/self")
			os.Readlink("fine")
			os.Readlink(gosigar.Procd + "se/lf")
			os.Readlink(path.Join("a", "b"))
		}
	`
}

func fixedContentForUsePathJoinCsFixer() string {
	// FIXME: import must be added?
	return `
		package main

		func main() {
			os.Readlink(path.Join(gosigar.Procd, "self"))
			os.Readlink(path.Join("foo", "self")
			os.Readlink("fine")
			os.Readlink(path.Join(gosigar.Procd, "se", "lf")
			os.Readlink(path.Join("a", "b"))
		}
	`
}