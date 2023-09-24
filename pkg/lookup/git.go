package lookup

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/noahshinn024/gcl/pkg/lookup/utils"
)

type Commit struct {
	AuthorName  string
	AuthorEmail string
	Hash        string
	Message     string
	Timestamp   time.Time
	Refs        string
	RemoteURL   string
}

type GitClient struct {
}

func (c *Commit) Display() string {
	refsDisplay := ""
	if c.Refs != "" {
		refsParts := strings.Split(c.Refs, ", ")
		coloredRefs := []string{}
		for _, ref := range refsParts {
			if strings.HasPrefix(ref, "HEAD ->") {
				ref = utils.CYAN + "HEAD ->" + utils.RESET + utils.GREEN + strings.TrimPrefix(ref, "HEAD ->") + utils.RESET
			} else if strings.HasPrefix(ref, "origin/") {
				ref = utils.RED + ref + utils.RESET
			} else {
				ref = utils.GREEN + ref + utils.RESET
			}
			coloredRefs = append(coloredRefs, ref)
		}
		refsDisplay = utils.YELLOW + " (" + utils.RESET + strings.Join(coloredRefs, utils.YELLOW+", "+utils.RESET) + utils.YELLOW + ")" + utils.RESET
	}

	return fmt.Sprintf("%scommit %s%s%s\nAuthor: %s <%s>\nDate:   %s\nURL: %s\n\n    %s\n\n",
		utils.YELLOW, c.Hash, utils.RESET, refsDisplay, c.AuthorName, c.AuthorEmail, c.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"), c.RemoteURL, c.Message)
}

func (gc *GitClient) getRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Convert git SSH URLs to HTTPS for better compatibility with terminals
	url := strings.TrimSpace(string(output))
	url = strings.Replace(url, "git@", "", 1)
	url = strings.Replace(url, "github.com:", "github.com/", 1)
	url = strings.Replace(url, ".git", "", 1)
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	return url, nil
}

func (gc *GitClient) GetCommitsFromLastDay(since string) ([]*Commit, error) {
	const separator = "<SEPARATOR>"
	cmd := exec.Command("git", "log", fmt.Sprintf("--since=\"%s\"", since), "--pretty=format:%H"+separator+"%an"+separator+"%ae"+separator+"%s"+separator+"%ct"+separator+"%D")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var commits []*Commit
	remoteBaseURL, err := gc.getRemoteURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote URL: %w", err)
	}
	for _, line := range lines {
		parts := strings.Split(line, separator)
		if len(parts) < 6 {
			continue
		}

		unixTimestamp, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			return nil, err
		}

		commits = append(commits, &Commit{
			Hash:        parts[0],
			AuthorName:  parts[1],
			AuthorEmail: parts[2],
			Message:     parts[3],
			Timestamp:   time.Unix(unixTimestamp, 0),
			Refs:        parts[5],
			RemoteURL:   fmt.Sprintf("%s/commit/%s", remoteBaseURL, parts[0]),
		})
	}

	return commits, nil
}

func (gc *GitClient) IsInsideGitWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}
