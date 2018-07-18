package fixers

import (
	"testing"
	"fmt"
	"strings"
)

func TestUsePathJoinCsFixerLint(t *testing.T) {
	fixer := &UsePathJoinCsFixer{}

	for _, testCase := range usePathJoinCsFixerTestTable() {
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

func TestUsePathJoinFix(t *testing.T) {
	fixer := &UsePathJoinCsFixer{}

	for _, testCase := range usePathJoinCsFixerTestTable() {
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

type pathJoinTestCase struct{
	test string
	expected string
	problems Problems
}

func usePathJoinCsFixerTestTable() []pathJoinTestCase {
	cases := []pathJoinTestCase {
		{`os.Readlink("fine")`, `os.Readlink("fine")`, make(Problems, 0)},
		{`os.Readlink("fine" + "shine")`, `os.Readlink("fine" + "shine")`, make(Problems, 0)},
		{`os.Readlink(gosigar.Procd + "self")`, `os.Readlink(gosigar.Procd + "self")`, make(Problems, 0)},
		{`os.Readlink(path.Join("a", "b"))`, `os.Readlink(path.Join("a", "b"))`, make(Problems, 0)},
		{`os.Readlink("foo/self")`, `os.Readlink(path.Join("foo", "self"))`, make(Problems, 0)},
		{`os.Readlink(gosigar.Procd + "se/lf")`, `os.Readlink(gosigar.Procd + path.Join("se", "lf"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/ne")`, `os.Readlink("fine" + path.Join("shi", "ne"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/ne" + "wi/ne")`, `os.Readlink("fine" + path.Join("shi", "ne") + path.Join("wi", "ne"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shine" + "wi/ne")`, `os.Readlink("fine" + "shine" + path.Join("wi", "ne"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/ne" + "wine")`, `os.Readlink("fine" + path.Join("shi", "ne") + "wine")`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/n/e" + "w/i/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e") + path.Join("w", "i", "ne"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/n/e" + "/wi/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))`, make(Problems, 0)},
		{`os.Readlink("fine/" + "shi/n/e" + "/wi/ne")`, `os.Readlink(path.Join("fine", "shi", "n", "e", "wi", "ne"))`, make(Problems, 0)},
		{`os.Readlink(gosigar.Procd + "/self")`, `os.Readlink(path.Join(gosigar.Procd, "self"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "/shine")`, `os.Readlink(path.Join("fine", "shine"))`, make(Problems, 0)},
		{`os.Readlink("/self")`, `os.Readlink(path.Join("", "self"))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shi/n/e/" + "wi/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))`, make(Problems, 0)},
		{`os.Readlink(gosigar.Procd + "self/")`, `os.Readlink(gosigar.Procd + path.Join("self", ""))`, make(Problems, 0)},
		{`os.Readlink("fine" + "shine/")`, `os.Readlink("fine" + path.Join("shine", ""))`, make(Problems, 0)},
		{`os.Readlink("self/")`, `os.Readlink(path.Join("self", ""))`, make(Problems, 0)},
		{`os.Readlink("/self/")`, `os.Readlink(path.Join("", "self", ""))`, make(Problems, 0)},
	}

	var totalCase pathJoinTestCase

	for k, _ := range cases {
		totalCase.test += cases[k].test + "\n\t"
		totalCase.expected += cases[k].expected + "\n\t"

		if cases[k].expected != cases[k].test {
			cases[k].problems = append(cases[k].problems, &Problem{Position: &Position{Line: 5}, Text: "Use path.Join"})
			totalCase.problems = append(totalCase.problems, &Problem{Position: &Position{Line: 5+k}, Text: "Use path.Join"})
		}

		cases[k].expected = getExpectedContentForUsePathJoinCsFixer(cases[k])
		cases[k].test = getTestContentForUsePathJoinCsFixer(cases[k].test)
	}

	totalCase.expected = strings.TrimRight(totalCase.expected, "\t")
	totalCase.test = strings.TrimRight(totalCase.test, "\t")

	totalCase.expected = getExpectedContentForUsePathJoinCsFixer(totalCase)
	totalCase.test = getTestContentForUsePathJoinCsFixer(totalCase.test)

	cases = append(cases, totalCase)

	return cases
}

func getTestContentForUsePathJoinCsFixer(content string) string {
	return `
		package main

		func main() {
			` + content + `
		}
	`
}

func getExpectedContentForUsePathJoinCsFixer(testCase pathJoinTestCase) string {
	if testCase.expected == testCase.test {
		return `package main

func main() {
	` + testCase.test + `
}
`
	}

	return `package main

import "path"

func main() {
	` + testCase.expected + `
}
`
}