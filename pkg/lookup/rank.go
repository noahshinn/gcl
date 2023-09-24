package lookup

import (
	"container/heap"
	"math"
	"strings"
	"sync"

	"github.com/noahshinn024/gcl/pkg/lookup/utils"
)

const k1 = 1.2
const b = 0.75

type CommitRankResult struct {
	Score  float64
	Commit *Commit
}

const rankCommitsSemSize = 100

func RankCommits(query string, commits []*Commit, numHits int) ([]CommitRankResult, error) {
	textcommits := utils.Map(commits, func(commit *Commit, _ int) string {
		return commit.Message
	})

	commitsLength := 0
	for _, document := range textcommits {
		commitsLength += len(document)
	}
	avgDocumentLength := float64(commitsLength) / float64(len(commits))

	type result struct {
		commitRankResult CommitRankResult
		err              error
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, rankCommitsSemSize)
	results := make([]*result, len(commits))
	for i, inputItem := range commits {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, inputItem *Commit) {
			defer func() { <-sem }()
			defer wg.Done()
			res, err := computeResult(query, inputItem, commits, k1, b, avgDocumentLength)
			results[i] = &result{
				commitRankResult: res,
				err:              err,
			}
		}(i, inputItem)
	}
	wg.Wait()

	for _, result := range results {
		if result.err != nil {
			return nil, result.err
		}
	}

	rankHeap := resultHeap{}
	heap.Init(&rankHeap)
	for _, res := range results {
		heap.Push(&rankHeap, res.commitRankResult)
	}

	maxNumResults := min(numHits, len(rankHeap))
	rankHeapResults := make([]CommitRankResult, maxNumResults)
	for i := 0; i < maxNumResults; i++ {
		rankHeapResults[i] = heap.Pop(&rankHeap).(CommitRankResult)
	}
	return rankHeapResults, nil
}

func idf(queryTerm string, commits []string) float64 {
	numDocuments := len(commits)
	numContaining := 0
	for _, document := range commits {
		if numOccurrence(queryTerm, document) > 0 {
			numContaining++
		}
	}
	numerator := float64(numDocuments-numContaining) + 0.5
	denominator := float64(numContaining) + 0.5
	return math.Log(numerator/denominator + 1)
}

func numOccurrence(queryTerm string, document string) int {
	return strings.Count(document, queryTerm)
}

func computeResult(query string, commit *Commit, commits []*Commit, k1 float64, b float64, avgDocumentLength float64) (CommitRankResult, error) {
	score := 0.0
	queryTerms, err := Tokenize(query)
	if err != nil {
		return CommitRankResult{}, err
	}

	textcommits := utils.Map(commits, func(commit *Commit, _ int) string {
		return commit.Message
	})
	for _, queryTerm := range queryTerms {
		document := commit.Message
		left := idf(queryTerm, textcommits)
		numerator := float64(numOccurrence(queryTerm, document)) * (k1 + 1)
		denominator := float64(numOccurrence(queryTerm, document)) + k1*(1-b+b*float64(len(document))/float64(avgDocumentLength))
		right := numerator / denominator
		score += left * right
	}
	return CommitRankResult{
		Score:  score,
		Commit: commit,
	}, nil
}

type resultHeap []CommitRankResult

func (h resultHeap) Len() int {
	return len(h)
}

func (h resultHeap) Less(i, j int) bool {
	return h[i].Score > h[j].Score
}
func (h resultHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *resultHeap) Push(x interface{}) {
	*h = append(*h, x.(CommitRankResult))
}
func (h *resultHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h resultHeap) Peek() interface{} {
	return h[0]
}
