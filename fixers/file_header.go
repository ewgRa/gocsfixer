package fixers

import (
	"strings"
)

func init() {
	AddFixer("file_header", func(options FixerOptions) (CsFixer, error) {
		header, err := options.extractRequiredString("header")

		if err != nil {
			return nil, err
		}

		lintText, err := options.extractString("lintText", "File header")

		if err != nil {
			return nil, err
		}

		return &FileHeaderCsFixer{header: header, lintText: lintText}, nil
	})
}

type FileHeaderCsFixer struct {
	header   string
	lintText string
}

func (l *FileHeaderCsFixer) Lint(content string) (Problems, error) {
	var problems Problems

	if !strings.HasPrefix(content, l.header) {
		lines := strings.Split(content, "\n")
		problems = append(problems, &Problem{Position: NewPosition(1), Text: l.lintText, LineText: lines[0]})
	}

	return problems, nil
}

func (l *FileHeaderCsFixer) Fix(content string) (string, error) {
	return l.header + content, nil
}

func (l *FileHeaderCsFixer) String() string {
	return "File header"
}
