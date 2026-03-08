package matchers

import (
	"fmt"
	"os/exec"
	"strings"
)

type MatcherRG struct {
}

func (m MatcherRG) Match(dirs []string, needle string) []string {
	echoCmd := exec.Command("echo", strings.Join(dirs, "\n"))
	rgCmd := exec.Command("rg", fmt.Sprintf("/[^/]*%[1]s[^/]*$|/[^/]*%[1]s[^/]*/[^/]*$", strings.ReplaceAll(needle, " ", ".*")))

	cmdPipe, err := echoCmd.StdoutPipe()
	if err != nil {
		return dirs
	}

	rgCmd.Stdin = cmdPipe

	echoCmd.Start()
	res, err := rgCmd.CombinedOutput()

	if err != nil {
		return dirs
	}

	return strings.Split(string(res[:len(res)-1]), "\n")
}

func IsMatcherRGAvailable() bool {
	cmd := exec.Command("rg", "--version")
	err := cmd.Run()

	return err == nil
}
