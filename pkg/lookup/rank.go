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

const rankSemSize = 100

type BM25Result struct {
	Score      float64
	QueryTerms []string
	InputItem  BM25InputItem
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

func computeResult(query string, inputItem BM25InputItem, corpus []BM25InputItem, k1 float64, b float64, avgDocumentLength float64) (BM25Result, error) {
	score := 0.0
	queryTerms, err := Tokenize(query)
	if err != nil {
		return BM25Result{}, err
	}

	textCorpus := utils.Map(corpus, func(inputItem BM25InputItem, _ int) string {
		return inputItem.GetText()
	})
	for _, queryTerm := range queryTerms {
		document := inputItem.GetText()
		left := idf(queryTerm, textCorpus)
		numerator := float64(numOccurrence(queryTerm, document)) * (k1 + 1)
		denominator := float64(numOccurrence(queryTerm, document)) + k1*(1-b+b*float64(len(document))/float64(avgDocumentLength))
		right := numerator / denominator
		score += left * right
	}
	return BM25Result{
		Score:      score,
		QueryTerms: queryTerms,
		InputItem:  inputItem,
	}, nil
}

type resultHeap []BM25Result

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
	*h = append(*h, x.(BM25Result))
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

type BM25InputItem interface {
	GetText() string
	GetItem() interface{}
}

func BM25RankN(query string, corpus []BM25InputItem, numHits int) ([]interface{}, error) {
	textCorpus := utils.Map(corpus, func(inputItem BM25InputItem, _ int) string {
		return inputItem.GetText()
	})

	corpusLength := 0
	for _, document := range textCorpus {
		corpusLength += len(document)
	}
	avgDocumentLength := float64(corpusLength) / float64(len(corpus))
	type result struct {
		rankResult BM25Result
		err        error
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, rankSemSize)
	results := make([]*result, len(corpus))
	for i, inputItem := range corpus {
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, inputItem BM25InputItem) {
			defer func() { <-sem }()
			defer wg.Done()
			res, err := computeResult(query, inputItem, corpus, k1, b, avgDocumentLength)
			results[i] = &result{
				rankResult: res,
				err:        err,
			}
		}(i, inputItem)
	}
	wg.Wait()

	scoredResults := make([]BM25Result, len(results))
	for i, result := range results {
		if result.err != nil {
			return nil, result.err
		}
		scoredResults[i] = result.rankResult
	}

	rankHeap := resultHeap{}
	heap.Init(&rankHeap)
	for _, result := range scoredResults {
		heap.Push(&rankHeap, result)
	}

	maxNumResults := min(numHits, len(rankHeap))
	rankHeapResults := make([]BM25Result, maxNumResults)
	for i := 0; i < maxNumResults; i++ {
		rankHeapResults[i] = heap.Pop(&rankHeap).(BM25Result)
	}
	return utils.Map(rankHeapResults, func(result BM25Result, _ int) interface{} { return result.InputItem.GetItem() }), nil
}
