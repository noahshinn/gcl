package main

import (
	"flag"

	"github.com/noahshinn024/gcl/cmd/gcl"
)

func main() {
	query := flag.String("query", "", "The query to search for")
	maxResults := flag.Int("n", 10, "The maximum number of results to return")
	since := flag.String("since", "1 week ago", "A since query to pass; see https://git-scm.com/docs/git-log for more information")
	flag.Parse()
	if *query == "" {
		panic("Must provide a query")
	}

	lookup := gcl.Lookup{
		Query:      *query,
		MaxResults: *maxResults,
		Since:      *since,
	}
	lookup.Run()
}
