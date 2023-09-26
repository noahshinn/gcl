package lookup

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/noahshinn024/gcl/pkg/lookup/utils"
)

type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (i *Issue) GetText() string {
	return fmt.Sprintf("%s\n%s", i.Title, i.Body)
}

func (i *Issue) Display(issues []*Issue) string {
	maxNumber, maxTitle, maxTimeDiff := getMaxIssueLengths(issues)
	numberFormat := fmt.Sprintf("#%%-%dd", max(maxNumber+5, 10))
	titleFormat := fmt.Sprintf("%%-%ds", max(maxTitle+10, 50))
	timeDiffFormat := fmt.Sprintf("%%-%ds", max(maxTimeDiff+5, 10))

	return utils.GREEN + fmt.Sprintf(numberFormat, i.Number) + utils.RESET +
		fmt.Sprintf(titleFormat, i.Title) + utils.RESET +
		utils.BOLD_MAGENTA + fmt.Sprintf(timeDiffFormat, utils.BuildTimeDiffDisplay(i.CreatedAt, "created")) + utils.RESET
}

type RankIssueItem struct {
	Issue *Issue
}

func (rii *RankIssueItem) GetItem() interface{} {
	return rii.Issue
}

func (rii *RankIssueItem) GetText() string {
	return rii.Issue.GetText()
}

type GHClient struct {
}

func (ghc *GHClient) IsAvailable() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func (ghc *GHClient) GetOpenIssues(owner, repo string) ([]*Issue, error) {
	if !ghc.IsAvailable() {
		return nil, fmt.Errorf("'gh' command is not available")
	}

	cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/%s/issues", owner, repo))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var issues []*Issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func getMaxIssueLengths(issues []*Issue) (int, int, int) {
	maxNumber := 0
	maxTitle := 0
	maxTimeDiff := 0
	for _, issue := range issues {
		numberLen := len(fmt.Sprintf("%d", issue.Number))
		if numberLen > maxNumber {
			maxNumber = numberLen
		}

		titleLen := len(issue.Title)
		if titleLen > maxTitle {
			maxTitle = titleLen
		}

		timeDiffLen := len(utils.BuildTimeDiffDisplay(issue.CreatedAt, "created"))
		if timeDiffLen > maxTimeDiff {
			maxTimeDiff = timeDiffLen
		}
	}
	return maxNumber, maxTitle, maxTimeDiff
}
