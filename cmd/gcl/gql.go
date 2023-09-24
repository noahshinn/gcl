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
}

func (l *Lookup) Run() error {
	gc := lookup.GitClient{}
	if !gc.IsInsideGitWorkTree() {
		return fmt.Errorf("not inside a git work tree")
	}
	commits, err := gc.GetCommitsFromLastDay(l.Since)
	if err != nil {
		return err
	}
	rankedResults, err := lookup.RankCommits(l.Query, commits, l.MaxResults)
	if err != nil {
		return err
	}
	displays := utils.Map(rankedResults, func(result lookup.CommitRankResult, _ int) string {
		return result.Commit.Display()
	})
	l.Display(strings.Join(displays, "\n"))
	return nil
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
