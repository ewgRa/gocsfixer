package fixers_test

import (
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func TestUsePathJoinCsFixerLint(t *testing.T) {
	assertLint(t, &fixers.UsePathJoinCsFixer{}, usePathJoinCsFixerTestTable())
}

func TestUsePathJoinFix(t *testing.T) {
	assertFix(t, &fixers.UsePathJoinCsFixer{}, usePathJoinCsFixerTestTable())
}

func usePathJoinCsFixerTestTable() []fixerTestCase {
	cases := []fixerTestCase{
		{`os.Readlink("fine")`, `os.Readlink("fine")`, fixers.Problems{}},
		{`os.Readlink("fine" + "shine")`, `os.Readlink("fine" + "shine")`, fixers.Problems{}},
		{`os.Readlink(gosigar.Procd + "self")`, `os.Readlink(gosigar.Procd + "self")`, fixers.Problems{}},
		{`os.Readlink(path.Join("a", "b"))`, `os.Readlink(path.Join("a", "b"))`, fixers.Problems{}},
		{`os.Readlink("foo/self")`, `os.Readlink(path.Join("foo", "self"))`, fixers.Problems{}},
		{`os.Readlink(gosigar.Procd + "se/lf")`, `os.Readlink(gosigar.Procd + path.Join("se", "lf"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/ne")`, `os.Readlink("fine" + path.Join("shi", "ne"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/ne" + "wi/ne")`, `os.Readlink("fine" + path.Join("shi", "ne") + path.Join("wi", "ne"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shine" + "wi/ne")`, `os.Readlink("fine" + "shine" + path.Join("wi", "ne"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/ne" + "wine")`, `os.Readlink("fine" + path.Join("shi", "ne") + "wine")`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/n/e" + "w/i/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e") + path.Join("w", "i", "ne"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/n/e" + "/wi/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))`, fixers.Problems{}},
		{`os.Readlink("fine/" + "shi/n/e" + "/wi/ne")`, `os.Readlink(path.Join("fine", "shi", "n", "e", "wi", "ne"))`, fixers.Problems{}},
		{`os.Readlink(gosigar.Procd + "/self")`, `os.Readlink(path.Join(gosigar.Procd, "self"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "/shine")`, `os.Readlink(path.Join("fine", "shine"))`, fixers.Problems{}},
		{`os.Readlink("/self")`, `os.Readlink(path.Join("", "self"))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shi/n/e/" + "wi/ne")`, `os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))`, fixers.Problems{}},
		{`os.Readlink(gosigar.Procd + "self/")`, `os.Readlink(gosigar.Procd + path.Join("self", ""))`, fixers.Problems{}},
		{`os.Readlink("fine" + "shine/")`, `os.Readlink("fine" + path.Join("shine", ""))`, fixers.Problems{}},
		{`os.Readlink("self/")`, `os.Readlink(path.Join("self", ""))`, fixers.Problems{}},
		{`os.Readlink("/self/")`, `os.Readlink(path.Join("", "self", ""))`, fixers.Problems{}},
	}

	for k, _ := range cases {
		if cases[k].expected != cases[k].test && len(cases[k].problems) == 0 {
			cases[k].problems = append(cases[k].problems, &fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Use path.Join"})
		}
	}

	cases = append(cases, getTotalTestCase(cases))

	for k, _ := range cases {
		cases[k].expected = getExpectedContentForUsePathJoinCsFixer(cases[k])
		cases[k].test = getTestContentForUsePathJoinCsFixer(cases[k].test)
	}

	return cases
}

func getTestContentForUsePathJoinCsFixer(content string) string {
	return `package main

		func main() {
			` + content + `
		}
	`
}

func getExpectedContentForUsePathJoinCsFixer(testCase fixerTestCase) string {
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
