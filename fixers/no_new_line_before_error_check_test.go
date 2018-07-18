package fixers

import (
	"testing"
	"fmt"
	"strings"
)

func TestNoNewLineBeforeErrorLint(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	for _, testCase := range noNewLineBeforeErrorCheckTestTable() {
		problems, _ := fixer.Lint(testCase.test)

		if (len(problems) != len(testCase.problems)) {
			fmt.Println("Expected", len(testCase.problems), "problem(s), got", len(problems))
			t.Fail()
			return
		}

		for k, problem := range testCase.problems {
			if problem.Position.Line != problems[k].Position.Line {
				fmt.Println("Problem found on", problems[k].Position.Line, "line, expected", problem.Position.Line)
				t.Fail()
				return
			}

			if problem.Text != problems[k].Text {
				fmt.Println("Problem have text", fmt.Sprintf("'%s'", problems[k].Text), ",", problem.Text, "expected")
				t.Fail()
				return
			}
		}
	}
}

func TestNoNewLineBeforeErrorFix(t *testing.T) {
	fixer := &NoNewLineBeforeErrorCsFixer{}

	for _, testCase := range noNewLineBeforeErrorCheckTestTable() {
		contentFix, err := fixer.Fix(testCase.test)

		if nil != err {
			fmt.Println("Error when perform fix:", err)
			t.Fail()
			return
		}

		if testCase.expected != contentFix {
			fmt.Println("Fixed content differ from expected")
			fmt.Println(contentFix, testCase.expected)
			t.Fail()
			return
		}

		problems, _ := fixer.Lint(contentFix)

		if len(problems) != 0 {
			fmt.Println("Expected no problem, got", len(problems))
			t.Fail()
			return
		}
	}
}

type noNewLineBeforeErrorCheckTestCase struct{
	test string
	expected string
	problems Problems
}

func noNewLineBeforeErrorCheckTestTable() []noNewLineBeforeErrorCheckTestCase {
	cases := []noNewLineBeforeErrorCheckTestCase {
		{
			`err = test()

			if err == nil {
			}`,
			`err = test()
			if err == nil {
			}`,
			make(Problems, 0),
		}, {
			`err = test()
			if err == nil {
			}`,
			`err = test()
			if err == nil {
			}`,
			make(Problems, 0),
		}, {
			`
			if anotherErr = test(); anotherErr == nil {
			}`,
			`
			if anotherErr = test(); anotherErr == nil {
			}`,
			make(Problems, 0),
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
			Problems{
				&Problem{Position: &Position{Line: 12}, Text: "No newline before check error"},
				&Problem{Position: &Position{Line: 15}, Text: "No newline before check error"},
			},
	}}

	var totalCase noNewLineBeforeErrorCheckTestCase

	for k, _ := range cases {
		if cases[k].expected != cases[k].test && len(cases[k].problems) == 0 {
			cases[k].problems = append(cases[k].problems, &Problem{Position: &Position{Line: 10}, Text: "No newline before check error"})
		}

		for _, problem := range cases[k].problems {
			totalCase.problems = append(
				totalCase.problems,
				&Problem{
					Position: &Position{Line: problem.Position.Line+strings.Count(totalCase.test, "\n")},
					Text: problem.Text,
				},
			)
		}

		totalCase.test += cases[k].test + "\n\t"
		totalCase.expected += cases[k].expected + "\n\t"

		cases[k].expected = getTestContentForNoNewLineBeforeErrorCheckCsFixer(cases[k].expected)
		cases[k].test = getTestContentForNoNewLineBeforeErrorCheckCsFixer(cases[k].test)
	}

	totalCase.expected = strings.TrimRight(totalCase.expected, "\t")
	totalCase.test = strings.TrimRight(totalCase.test, "\t")

	totalCase.expected = getTestContentForNoNewLineBeforeErrorCheckCsFixer(totalCase.expected)
	totalCase.test = getTestContentForNoNewLineBeforeErrorCheckCsFixer(totalCase.test)

	cases = append(cases, totalCase)

	return cases
}

func getTestContentForNoNewLineBeforeErrorCheckCsFixer(content string) string {
	return `
		package main

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
