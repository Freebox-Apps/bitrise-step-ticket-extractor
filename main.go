package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	RepoDirectoryEnv   = "repo_dir"
	TicketDelimiterEnv = "output_delimiter"
	TicketIdRegex      = "[#]([^,\\s\n]+)"
)

func main() {
	commitStrList := getCommitStringList()
	fmt.Printf("Found %d commit candidates\n", len(commitStrList))

	allTickets := []string{}
	for i := 0; i < len(commitStrList); i++ {
		allTickets = append(allTickets, extractSovledTickets(commitStrList[i])...)
	}

	ticketListStrigns := strings.Join(allTickets[:], getOutputDelimiter())
	//entries := createEntries(prefixStrList)
	//fillCommitInfo(commitStrList, entries)

	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "TICKET_LIST", "--value", ticketListStrigns).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}

func extractSovledTickets(message string) []string {
	regex := regexp.MustCompile(TicketIdRegex)
	matches := regex.FindAllStringSubmatch(message, -1)
	var result []string
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func getOutputDelimiter() string {
	return os.Getenv(TicketDelimiterEnv)
}
