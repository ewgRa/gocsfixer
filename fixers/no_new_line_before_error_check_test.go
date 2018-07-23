package fixers_test

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"testing"
)

func TestNoNewLineBeforeErrorLint(t *testing.T) {
	assertLint(t, &fixers.NoNewLineBeforeErrorCsFixer{}, noNewLineBeforeErrorCheckTestTable())
}

func TestNoNewLineBeforeErrorFix(t *testing.T) {
	assertFix(t, &fixers.NoNewLineBeforeErrorCsFixer{}, noNewLineBeforeErrorCheckTestTable())
}

func noNewLineBeforeErrorCheckTestTable() []fixerTestCase {
	cases := []fixerTestCase{
		{
			`err = test()

			if err == nil {
			}`,
			`err = test()
			if err == nil {
			}`,
			fixers.Problems{},
		}, {
			`err = test()
			if err == nil {
			}`,
			`err = test()
			if err == nil {
			}`,
			fixers.Problems{},
		}, {
			`
			if anotherErr = test(); anotherErr == nil {
			}`,
			`
			if anotherErr = test(); anotherErr == nil {
			}`,
			fixers.Problems{},
		}, {
			`err = test()
			if err == nil {
			}

			if err != nil {
			}

			if nil==err {
			}`,
			`err = test()
			if err == nil {
			}
			if err != nil {
			}
			if nil==err {
			}`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 11}, Text: "No newline before check error"},
				&fixers.Problem{Position: &fixers.Position{Line: 14}, Text: "No newline before check error"},
			},
		}}

	for k, _ := range cases {
		if cases[k].expected != cases[k].test && len(cases[k].problems) == 0 {
			cases[k].problems = append(cases[k].problems, &fixers.Problem{Position: &fixers.Position{Line: 9}, Text: "No newline before check error"})
		}
	}

	cases = append(cases, getTotalTestCase(cases))

	for k, _ := range cases {
		cases[k].expected = getTestContentForNoNewLineBeforeErrorCheckCsFixer(cases[k].expected)
		cases[k].test = getTestContentForNoNewLineBeforeErrorCheckCsFixer(cases[k].test)
	}

	return cases
}

func getTestContentForNoNewLineBeforeErrorCheckCsFixer(content string) string {
	return `package main

		func test() error {
			return nil
		}

		func b() {
			` + content + `
		}

		func main() {
			b()
		}
	`
}
