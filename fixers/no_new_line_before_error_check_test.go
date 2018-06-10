package fixers

import (
	"testing"
	"fmt"
)

func TestNoNewLineBeforeErrorCsFixerLint(t *testing.T) {
	expected := []struct {
		line int
		text string
	}{
		{22, "No newline before check error"},
		{25, "No newline before check error"},
	}

	fixer := &NoNewLineBeforeErrorCsFixer{}

	problems, _ := fixer.Lint(contentForNoNewLineBeforeErrorCsFixer())

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

func TestNoNewLineBeforeErrorCsFixerFix(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	contentFix, err := fixer.Fix(contentForNoNewLineBeforeErrorCsFixer())

	if nil != err {
		fmt.Println("Error when perform fix:", err)
		t.Fail()
		return
	}

	expectedFixedContent := fixedContentForNoNewLineBeforeErrorCsFixer()

	if expectedFixedContent != contentFix {
		fmt.Println("Fixed content differ from expected")
		fmt.Println(contentFix)
		t.Fail()
		return
	}
}

func contentForNoNewLineBeforeErrorCsFixer() string {
	return `
		package main

		func test() error {
			return nil
		}

		func b() {
			err := test()
			if err == nil {
			}

			if anotherErr := test(); anotherErr == nil {
			}

			if anotherErr := test(); anotherErr == nil {
			}

			err = test()
			if err == nil {
			}

			if err != nil {
			}

			if nil==err {
			}
		}

		func main() {
			b()
		}
	`
}

func fixedContentForNoNewLineBeforeErrorCsFixer() string {
	return `
		package main

		func test() error {
			return nil
		}

		func b() {
			err := test()
			if err == nil {
			}

			if anotherErr := test(); anotherErr == nil {
			}

			if anotherErr := test(); anotherErr == nil {
			}

			err = test()
			if err == nil {
			}
			if err != nil {
			}
			if nil==err {
			}
		}

		func main() {
			b()
		}
	`
}