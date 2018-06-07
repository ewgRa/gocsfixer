package fixers

import (
	"testing"
	"fmt"
)

func TestNoNewLineBeforeErrorCsFixerLint(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	problems, _ := fixer.Lint(content())

	if (len(problems) != 2) {
		fmt.Println("Expected two problems, got", len(problems))
		t.Fail()
		return
	}

	if problems[0].Position.Line != 22 {
		fmt.Println("First problem found on", problems[0].Position.Line, "line, expected 22")
		t.Fail()
		return
	}

	if problems[0].Text != "No newline before check error" {
		fmt.Println("First problem have text", problems[0].Text, ", 'No newline before check error' expected")
		t.Fail()
		return
	}

	if problems[1].Position.Line != 25 {
		fmt.Println("Second problem found on", problems[1].Position.Line, "line, expected 25")
		t.Fail()
		return
	}

	if problems[1].Text != "No newline before check error" {
		fmt.Println("Second problem have text", problems[1].Text, ", 'No newline before check error' expected")
		t.Fail()
		return
	}
}

func TestNoNewLineBeforeErrorCsFixerFix(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	contentFix, err := fixer.Fix(content())

	if nil != err {
		fmt.Println("Error when perform fix:", err)
		t.Fail()
		return
	}

	expectedFixedContent := fixedContent()

	if expectedFixedContent != contentFix {
		fmt.Println("Fixed content differ from expected")
		fmt.Println(contentFix)
		t.Fail()
		return
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