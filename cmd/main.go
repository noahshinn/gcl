package main

import (
	"flag"
	"fmt"

	"github.com/noahshinn024/gcl/cmd/gcl"
)

func main() {
	query := flag.String("query", "", "The query to search for")
	maxResults := flag.Int("n", 10, "The maximum number of results to return")
	since := flag.String("since", "1 week ago", "A since query to pass; see https://git-scm.com/docs/git-log for more information")
	author := flag.String("author", "", "The author pattern to filter by")
	mode := flag.String("mode", "commits", "The mode to run in; either 'commits' or 'issues'")
	flag.Parse()
	if *query == "" {
		panic("Must provide a query")
	}
	lookupMode, err := gcl.Mode(*mode)
	if err != nil {
		panic(fmt.Errorf("invalid mode: %s", *mode))
	}

	lookup := gcl.Lookup{
		Query:      *query,
		MaxResults: *maxResults,
		Since:      *since,
		Author:     *author,
		Mode:       lookupMode,
	}
	lookup.Run()
}
