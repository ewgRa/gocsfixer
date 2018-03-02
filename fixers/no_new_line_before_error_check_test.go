package fixers

import (
	"testing"
	"fmt"
)

func TestNoNewLineBeforeErrorCsFixerLint(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	problems, _ := fixer.Lint(content())

	if (len(problems) != 2) {
		t.Fail()
	}

	if problems[0].Position.line != 23 {
		fmt.Println("First problem found on", problems[0].Position.line, "line, expected 23")
		t.Fail()
	}

	if problems[0].Text != "No newline before check error" {
		fmt.Println("First problem have text", problems[0].Text, ", 'No newline before check error' expected")
		t.Fail()
	}

	if problems[1].Position.line != 26 {
		fmt.Println("Second problem found on", problems[1].Position.line, "line, expected 26")
		t.Fail()
	}

	if problems[1].Text != "No newline before check error" {
		fmt.Println("Second problem have text", problems[1].Text, ", 'No newline before check error' expected")
		t.Fail()
	}
}

func TestNoNewLineBeforeErrorCsFixerFix(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	contentFix, err := fixer.Fix(content())

	if nil != err {
		fmt.Println("Error when perform fix:", err)
		t.Fail()
	}

	expectedFixedContent := fixedContent()

	if expectedFixedContent != contentFix {
		fmt.Println("Fixed content differ from expected")
		fmt.Println(contentFix)
		t.Fail()
	}
}

func content() string {
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

func fixedContent() string {
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