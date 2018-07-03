package fixers

import (
	"strings"
	"errors"
)

func init() {
	AddFixer("file_header", func (options FixerOptions) (CsFixer, error) {
		headerOption, ok := options["header"]

		if (!ok) {
			return nil, errors.New("Header option is required")
		}

		header, ok := headerOption.(string)

		if (!ok) {
			return nil, errors.New("Wrong header option")
		}

		lintText := "File header"

		lintTextOption, ok := options["lintText"]

		if (ok) {
			lintText, ok = lintTextOption.(string)

			if (!ok) {
				return nil, errors.New("Wrong header option")
			}
		}

		return &FileHeaderCsFixer{header: header, lintText: lintText}, nil
	})
}

type FileHeaderCsFixer struct {
	header string
	lintText string
}

func (l *FileHeaderCsFixer) Lint(content string) (Problems, error) {
	var problems []*Problem

	if (!strings.HasPrefix(content, l.header)) {
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
