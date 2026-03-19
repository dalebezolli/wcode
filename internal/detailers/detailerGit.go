package detailers

import (
	"os/exec"
	"strings"
)

type DetailerGit struct {
}

const MARKDOWN_H1_CHARS = 2

const LABEL_CURRENT_BRANCH = "Branch: "
const LABEL_CURRENT_COMMIT_HASH = "Commit hash: "
const LABEL_CURRENT_COMMIT = "Commit label: "
const LABEL_COMMITS_AHEAD = "Ahead by: "
const LABEL_COMMITS_BEHIND = "Behind by: "

func (d DetailerGit) GetDetails(path string) Details {
	var title string
	titleCmd := exec.Command("head", "-n", "1", "README.md")
	titleCmd.Dir = path

	out, err := titleCmd.Output()

	if err != nil {
		splitPath := strings.Split(path, "/")
		title = splitPath[len(splitPath)-1]
	} else {
		var startPos int = 0

		if out[0] == '#' && out[1] == ' ' {
			startPos = MARKDOWN_H1_CHARS
		}

		title = strings.TrimSpace(string(out[startPos : len(out)-1]))
	}

	rest := make(map[string]string)

	currentBranchCmd := exec.Command("git", "branch", "--show-current")
	currentBranchCmd.Dir = path
	if out, err = currentBranchCmd.Output(); err == nil {
		rest[LABEL_CURRENT_BRANCH] = string(out[:len(out)-1])
	}

	commitCmd := exec.Command("git", "log", "--oneline", "-n", "1")
	commitCmd.Dir = path
	if out, err = commitCmd.Output(); err == nil {
		line := string(out[:len(out)-1])

		if len(line) > 7 {
			rest[LABEL_CURRENT_COMMIT_HASH] = line[:7]
			// rest[LABEL_CURRENT_COMMIT] = strings.TrimSpace(line[7:])
		}
	}

	fetchCmd := exec.Command("git", "fetch")
	fetchCmd.Dir = path

	statusCmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{u}")
	statusCmd.Dir = path
	if out, err = statusCmd.Output(); err == nil {
		counts := strings.Fields(string(out[:len(out)-1]))

		if len(counts) == 2 {
			rest[LABEL_COMMITS_AHEAD] = counts[0]
			rest[LABEL_COMMITS_BEHIND] = counts[1]
		}
	}

	return Details{
		Title: title,
		Path:  path,

		Rest: rest,
	}
}

func (d DetailerGit) GetRestOrder() []string {
	return []string{
		LABEL_CURRENT_BRANCH,
		LABEL_CURRENT_COMMIT_HASH,
		LABEL_CURRENT_COMMIT,
		LABEL_COMMITS_AHEAD,
		LABEL_COMMITS_BEHIND,
	}
}

func IsDetailerGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
