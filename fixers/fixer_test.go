package fixers_test

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"testing"
	"fmt"
	"strings"
)

type fixerTestCase struct{
	test string
	expected string
	problems fixers.Problems
}

func assertLint(t *testing.T, linter fixers.Linter, testTable []fixerTestCase) {
	for _, testCase := range testTable {
		problems, _ := linter.Lint(testCase.test)

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

func assertFix(t *testing.T, csFixer fixers.CsFixer, testTable []fixerTestCase) {
	fixer := csFixer.(fixers.Fixer)
	linter := csFixer.(fixers.Linter)

	for _, testCase := range testTable {
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

		problems, _ := linter.Lint(contentFix)

		if len(problems) != 0 {
			fmt.Println("Expected no problem, got", len(problems))
			t.Fail()
			return
		}
	}
}

func getTotalTestCase(cases []fixerTestCase) fixerTestCase {
	var totalCase fixerTestCase

	for k, _ := range cases {
		for _, problem := range cases[k].problems {
			totalCase.problems = append(
				totalCase.problems,
				&fixers.Problem{
					Position: &fixers.Position{Line: problem.Position.Line+strings.Count(totalCase.test, "\n")},
					Text: problem.Text,
				},
			)
		}

		totalCase.test += cases[k].test + "\n\t"
		totalCase.expected += cases[k].expected + "\n\t"
	}

	totalCase.expected = strings.TrimRight(totalCase.expected, "\t")
	totalCase.test = strings.TrimRight(totalCase.test, "\t")

	return totalCase
}