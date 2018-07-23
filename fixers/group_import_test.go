package fixers_test

import (
	"github.com/ewgRa/gocsfixer/fixers"
	"testing"
)

func TestGroupImportFixerLint(t *testing.T) {
	assertLint(t, createGroupImportFixer(), GroupImportFixerTestTable())
}

func TestGroupImportFix(t *testing.T) {
	assertFix(t, createGroupImportFixer(), GroupImportFixerTestTable())
}

func GroupImportFixerTestTable() []fixerTestCase {
	return []fixerTestCase{
		{
			`package main

import (
	"os"
	"strconv"

	"path"
)

func main() {
}
`,
			`package main

import (
	"os"
	"path"
	"strconv"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 7}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"go/api" // foobar
	// foo doc
	// foo doc multi
	"go/token" // foo
	"testing"  // bar

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			`package main

import (
	"go/api" // foobar
	// foo doc
	// foo doc multi
	"go/token" // foo
	"testing"  // bar

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{},
		},
		{
			`package main

import (
	"testing"
	"github.com/ewgRa/gocsfixer/fixers"

	"go/token"
)

func main() {
}
`,
			`package main

import (
	"go/token"
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 7}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"testing" // bar
	// foo doc
	"go/token" // foo
	"github.com/ewgRa/gocsfixer/fixers"

	// foobar doc
	"go/api" // foobar

	// foobar doc
	// multidoc
	"go/apifoo" // foobar
)

func main() {
}
`,
			`package main

import (
	// foobar doc
	"go/api" // foobar
	// foobar doc
	// multidoc
	"go/apifoo" // foobar
	// foo doc
	"go/token" // foo
	"testing"  // bar

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 6}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 10}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 14}, Text: "Group stdlib imports"},
			},
		},
		/*{
			`package main

		// foobar
		import (
			"testing"
			"github.com/ewgRa/gocsfixer/fixers"

			"go/token"
		)

		func main() {
		}
		`,
			`package main

		// foobar
		import (
			"go/token"
			"testing"

			"github.com/ewgRa/gocsfixer/fixers"
		)

		func main() {
		}
		`,
					fixers.Problems{
						&fixers.Problem{Position: &fixers.Position{Line: 5}, Text: "Group stdlib imports"},
						&fixers.Problem{Position: &fixers.Position{Line: 8}, Text: "Group stdlib imports"},
					},
				},*/
		{
			`package main

// foo doc
import "testing" // foo comment

import (
	"github.com/ewgRa/gocsfixer/fixers"

	"go/token"
)

// bar doc
import "os" // bar comment

func main() {
}
`,
			`package main

import (
	"go/token"
	// bar doc
	"os" // bar comment
	// foo doc
	"testing" // foo comment

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 9}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 13}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"github.com/ewgRa/gocsfixer/fixers"

	"go/token"
)

import "testing"

func main() {
}
`,
			`package main

import (
	"go/token"
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 6}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 9}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"go/token"
	"github.com/ewgRa/gocsfixer/fixers"

	"testing"
)

func main() {
}
`,
			`package main

import (
	"go/token"
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 7}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			`package main

import (
	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{},
		},
		{
			`package main

import (
	"testing"
)

func main() {
}
`,
			`package main

import (
	"testing"
)

func main() {
}
`,
			fixers.Problems{},
		},
		{
			`package main

import (
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			`package main

import (
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{},
		},
		{
			`package main

import (
	"github.com/ewgRa/gocsfixer/fixers"

	"testing"
)

func main() {
}
`,
			`package main

import (
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 6}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

import (
	"testing"
	"github.com/ewgRa/gocsfixer/fixers"

	// Comment
	"go/token"
)

func main() {
}
`,
			`package main

import (
	// Comment
	"go/token"
	"testing"

	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 4}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 8}, Text: "Group stdlib imports"},
			},
		},
		{
			`package main

// foo doc
import "github.com/ewgRa/gocsfixer"

import (
	"testing"
	"github.com/ewgRa/gocsfixer/fixers"

	// Comment
	"go/token"
)

func main() {
}
`,
			`package main

import (
	// Comment
	"go/token"
	"testing"
)

// foo doc
import "github.com/ewgRa/gocsfixer"

import (
	"github.com/ewgRa/gocsfixer/fixers"
)

func main() {
}
`,
			fixers.Problems{
				&fixers.Problem{Position: &fixers.Position{Line: 7}, Text: "Group stdlib imports"},
				&fixers.Problem{Position: &fixers.Position{Line: 11}, Text: "Group stdlib imports"},
			},
		},
	}
}

func createGroupImportFixer() *fixers.GroupImportFixer {
	mapFixer, _ := fixers.CreateFixer(
		"group_import",
		fixers.FixerOptions{"stdLib": true, "lintText": "Group stdlib imports"},
	)

	return mapFixer.(*fixers.GroupImportFixer)
}
