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

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "TICKET_LIST", "--value", ticketListStrigns).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
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
