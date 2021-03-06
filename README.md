Archived in favor of https://github.com/go-critic/go-critic
# gocsfixer
Golang coding style fixer, help you automize coding style checks.

[![Build Status](https://api.travis-ci.com/ewgRa/gocsfixer.svg?branch=master)](https://travis-ci.org/ewgRa/gocsfixer)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/ewgra/gocsfixer/master/LICENSE)
[![GoReportCard](http://goreportcard.com/badge/ewgra/gocsfixer)](http://goreportcard.com/report/ewgra/gocsfixer)
[![codecov.io](https://codecov.io/github/ewgRa/gocsfixer/coverage.svg?branch=master)](https://codecov.io/github/ewgRa/gocsfixer?branch=master)

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
    our_fixer_human_readable_alias:
        type: no_new_line_before_error_check
        recommend: true
        lint: true
        fix: true
```

With this configuration we want use "no_new_line_before_error_check" fixer and define it as recommendation and linter, that can be fixed by run "gocsfixer -fix" command.

Levels allow you to deal with false-positives and false-negatives results.
For example if you think, that your check never gives false-positives and false-negatives cases, you can enable it on all three levels.
If you not sure about how strict you are with check, you can define it only as recommended check.
Or for example you know that this check must be on lint level, but it can be fixed only manually - than you define your check as lint: true, fix: false.

## Available fixers
- alternative_call - use one function instead of other, for example force to use logp.Err instead of logp.Warn
- file_header - check that file have specific header at the beginning of file, for example license header
- group_import - group imports, for now on support grouping std lib imports in one import
- no_new_line_before_error_check - check that there is no new line before "if err != nil"
- use_path_join - check that instead `os.Readlink("foo/bar" + "/foobar")` path.Join used and code looks like `os.Readlink(path.Join("foo", "bar", "foobar"))`

Check cmd/gocsfixer/.gocsfixer.yml example configuration file for details and available options.
