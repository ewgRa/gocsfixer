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
		{9, "Use path.Join"},
		{10, "Use path.Join"},
		{11, "Use path.Join"},
		{12, "Use path.Join"},
		{13, "Use path.Join"},
		{14, "Use path.Join"},
		{15, "Use path.Join"},
		{16, "Use path.Join"},
		{17, "Use path.Join"},
		{18, "Use path.Join"},
		{19, "Use path.Join"},
		{20, "Use path.Join"},
		{21, "Use path.Join"},
		{22, "Use path.Join"},
		{23, "Use path.Join"},
		{24, "Use path.Join"},
		{25, "Use path.Join"},
	}

	fixer := &UsePathJoinCsFixer{}

	problems, _ := fixer.Lint(contentForUsePathJoinCsFixer())

	if (len(problems) != len(expected)) {
		fmt.Println("Expected", len(expected), "problems, got", len(problems))
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
	fixer := &UsePathJoinCsFixer{}

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

	problems, _ := fixer.Lint(expectedFixedContent)

	if len(problems) != 0 {
		fmt.Println("Expected no problem, got", len(problems))
		t.Fail()
		return
	}
}

// FIXME XXX: use table tests
func contentForUsePathJoinCsFixer() string {
	return `
		package main

		func main() {
			os.Readlink("fine")
			os.Readlink("fine" + "shine")
			os.Readlink(gosigar.Procd + "self")
			os.Readlink(path.Join("a", "b"))
			os.Readlink("foo/self")
			os.Readlink(gosigar.Procd + "se/lf")
			os.Readlink("fine" + "shi/ne")
			os.Readlink("fine" + "shi/ne" + "wi/ne")
			os.Readlink("fine" + "shine" + "wi/ne")
			os.Readlink("fine" + "shi/ne" + "wine")
			os.Readlink("fine" + "shi/n/e" + "w/i/ne")
			os.Readlink("fine" + "shi/n/e" + "/wi/ne")
			os.Readlink("fine/" + "shi/n/e" + "/wi/ne")
			os.Readlink(gosigar.Procd + "/self")
			os.Readlink("fine" + "/shine")
			os.Readlink("/self")
			os.Readlink("fine" + "shi/n/e/" + "wi/ne")
			os.Readlink(gosigar.Procd + "self/")
			os.Readlink("fine" + "shine/")
			os.Readlink("self/")
			os.Readlink("/self/")
		}
	`
}

func fixedContentForUsePathJoinCsFixer() string {
	return `package main

import "path"

func main() {
	os.Readlink("fine")
	os.Readlink("fine" + "shine")
	os.Readlink(gosigar.Procd + "self")
	os.Readlink(path.Join("a", "b"))
	os.Readlink(path.Join("foo", "self"))
	os.Readlink(gosigar.Procd + path.Join("se", "lf"))
	os.Readlink("fine" + path.Join("shi", "ne"))
	os.Readlink("fine" + path.Join("shi", "ne") + path.Join("wi", "ne"))
	os.Readlink("fine" + "shine" + path.Join("wi", "ne"))
	os.Readlink("fine" + path.Join("shi", "ne") + "wine")
	os.Readlink("fine" + path.Join("shi", "n", "e") + path.Join("w", "i", "ne"))
	os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))
	os.Readlink(path.Join("fine", "shi", "n", "e", "wi", "ne"))
	os.Readlink(path.Join(gosigar.Procd, "self"))
	os.Readlink(path.Join("fine", "shine"))
	os.Readlink(path.Join("", "self"))
	os.Readlink("fine" + path.Join("shi", "n", "e", "wi", "ne"))
	os.Readlink(gosigar.Procd + path.Join("self", ""))
	os.Readlink("fine" + path.Join("shine", ""))
	os.Readlink(path.Join("self", ""))
	os.Readlink(path.Join("", "self", ""))
}
`
}