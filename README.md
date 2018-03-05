# gocsfixer
Golang coding standards fixer, help you automize coding style checks.

## Run
gocsfixer have several flags, that allow you to choose, which one levels of checks you want to perform and control exit code:
- *"- recommend"* - when you want just show recommendations. In this case exit code will be always 0, and in output you will have recommendation from fixers configured as "recommend: true".
- *"- lint"* - if one of fixer, configured with "lint: true" found coding style problem, exit code will be 1 and you will have output with error message.
- *"- fix"* - if you run fixer with this flag, it just run fixers, that configured as "fix: true" and they will change your source code with apply fixes on it.

You can combine flags, like "-recommend -lint".

## Configuration
Example of configuration file, where you need configure fixers and desired levels, where they work:
```
fixers:
    no_new_line_before_error_check:
        recommend: true
        lint: true
        fix: true
```

With this configuration we want use "no_new_line_before_error_check" fixer and define it as recommendation and linter, that can be fixed by run "gocsfixer -fix" command.

Each check have three levels, that allow you to deal with false-positives and false-negatives results.
For example if you think, that your check never gives false-positives and false-negatives cases, you can enable it on all three levels.
If you not sure about how strict you are with check, you can define it only as recommended check.
Or for example you know that this check must be on lint level, but it can be fixed only manually - than you define your check as lint: true, fix: false.





Early alpha version

Inspired by:

- https://github.com/FriendsOfPHP/PHP-CS-Fixer

- https://github.com/golang/lint

- https://github.com/golang/go/tree/master/src/fmt

- https://github.com/golang/go/tree/master/src/cmd/fix

- https://github.com/elastic/beats/issues/6273

- https://github.com/golang/lint/issues/263

- https://github.com/golang/example/tree/master/gotypes (CheckNilFuncComparison)

- https://github.com/golang/tools/commit/6d70fb2e85323e81c89374331d3d2b93304faa36 (tests)
