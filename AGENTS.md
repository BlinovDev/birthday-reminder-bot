# Guidelines for AI contributions

This repository contains a small Telegram bot written in Go. To maintain
consistency and keep builds working, follow these rules when updating code
or documentation.

## Formatting and style
- Always run `gofmt -w` on any `.go` file you modify.
- Keep commit messages short and in the imperative mood ("Add feature" rather than "Added feature").

## Quality checks
- After making changes, run `go vet ./...` and `go test ./...` to ensure the
  code compiles. There are currently no unit tests, but `go test` should
  succeed without errors.

## Pull request notes
- Summarize the changes you made and mention the result of the checks above
  in the pull request description.
