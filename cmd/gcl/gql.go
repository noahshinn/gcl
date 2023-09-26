package gcl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/noahshinn024/gcl/pkg/lookup"
	"github.com/noahshinn024/gcl/pkg/lookup/utils"
)

type Lookup struct {
	Query      string
	MaxResults int
	Since      string
	Mode       LookupMode
}

type LookupMode string

const (
	Commits LookupMode = "commits"
	Issues  LookupMode = "issues"
)

func Mode(s string) (LookupMode, error) {
	switch s {
	case "commits":
		return Commits, nil
	case "issues":
		return Issues, nil
	default:
		return "", fmt.Errorf("invalid mode: %s", s)
	}
}

func (l *Lookup) Run() error {
	var displays []string
	switch l.Mode {
	case Commits:
		results, err := l.runCommits()
		if err != nil {
			return err
		}
		displays = results
	case Issues:
		results, err := l.runIssues()
		if err != nil {
			return err
		}
		displays = results
	default:
		return fmt.Errorf("invalid mode: %s", l.Mode)
	}
	l.Display(strings.Join(displays, "\n"))
	return nil
}

func (l *Lookup) runCommits() ([]string, error) {
	gc := lookup.GitClient{}
	if !gc.IsInsideGitWorkTree() {
		return nil, fmt.Errorf("not inside a git work tree")
	}
	commits, err := gc.GetCommitsFromLastDay(l.Since)
	if err != nil {
		return nil, err
	}
	rankInputs := utils.Map(commits, func(commit *lookup.Commit, _ int) lookup.BM25InputItem {
		return &lookup.RankCommitItem{
			Commit: commit,
		}
	})
	rankedResults, err := lookup.BM25RankN(l.Query, rankInputs, l.MaxResults)
	if err != nil {
		return nil, err
	}
	displays := make([]string, len(rankedResults))
	for i, result := range rankedResults {
		displays[i] = result.(*lookup.Commit).Display()
	}
	return displays, nil
}

func (l *Lookup) runIssues() ([]string, error) {
	gc := lookup.GitClient{}
	ghc := lookup.GHClient{}
	if !gc.IsInsideGitWorkTree() {
		return nil, fmt.Errorf("not inside a git work tree")
	}
	repoOwner, repoName, err := gc.GetCurrentRepoInfo()
	if err != nil {
		return nil, err
	}
	issues, err := ghc.GetOpenIssues(repoOwner, repoName)
	if err != nil {
		return nil, err
	}
	rankInputs := utils.Map(issues, func(issue *lookup.Issue, _ int) lookup.BM25InputItem {
		return &lookup.RankIssueItem{
			Issue: issue,
		}
	})
	rankedResults, err := lookup.BM25RankN(l.Query, rankInputs, l.MaxResults)
	if err != nil {
		return nil, err
	}
	displays := make([]string, len(rankedResults))
	for i, result := range rankedResults {
		displays[i] = result.(*lookup.Issue).Display(issues)
	}
	return displays, nil
}

func (l *Lookup) Display(s string) {
	lessCmd := exec.Command("less")
	in, _ := lessCmd.StdinPipe()
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr
	go func() {
		defer in.Close()
		fmt.Fprint(in, s)
	}()

	err := lessCmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run pager: %s\n", err)
	}
}
