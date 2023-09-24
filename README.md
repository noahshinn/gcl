# gcl

## Git commit lookup

This command enables lexical search for git commit messages.

## Installation

Build from source

```bash
go install github.com/noahshinn024/gcl/cmd@latest
```

## Usage

Make sure that your current working directory is in a git project, then run:

```bash
gcl --query "something change about search and ranking"
```

`gcl` passes the `since` flag to `git log` (the default query is 1 week), so you can do:

```bash
gcl --query "something change about search and ranking" --since "1 day ago"
```

or

```bash
gcl --query "something change about search and ranking" --since "2014-02-12T16:36:00-07:00"
```

`gcl` returns the top 10 results by default. Pass `--n` to get more or less results:

```bash
gcl --query "something change about search and ranking" --n 1000
```
