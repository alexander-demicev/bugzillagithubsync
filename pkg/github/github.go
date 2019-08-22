package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alexander-demichev/bugzillagithubsync/config"
	"github.com/alexander-demichev/bugzillagithubsync/pkg/bz"
	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

const (
	owner = "fusor"
	repo  = "cpma"
)

// InitClient init github client
func InitClient(ctx context.Context) *github.Client {
	ghToken := config.Config().GetString("GH_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

// CreateIssues create issues
func CreateIssues(ctx context.Context, client *github.Client, bugs *bz.BugList) error {
	issueNames, err := GetIssueNames(ctx, client)
	if err != nil {
		return nil
	}

	for _, bug := range bugs.Bugs {
		if isCreated(issueNames, bug.ID) {
			continue
		}

		title := fmt.Sprintf("BZ #%d: %s", bug.ID, bug.Summary)
		body := title + "\n https://bugzilla.redhat.com/show_bug.cgi?id=" + strconv.Itoa(bug.ID)

		_, _, err := client.Issues.Create(ctx, owner, repo, &github.IssueRequest{
			Title: &title,
			Body:  &body,
		})

		if err != nil {
			return nil
		}
	}

	return nil
}

// GetIssueNames get all issues by repo
func GetIssueNames(ctx context.Context, client *github.Client) ([]string, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State: "all",
	})

	if err != nil {
		return nil, err
	}

	issueNames := []string{}

	for _, issue := range issues {
		issueNames = append(issueNames, *issue.Title)
	}

	return issueNames, nil
}

func isCreated(issueNames []string, id int) bool {
	for _, issueName := range issueNames {
		if strings.Contains(issueName, strconv.Itoa(id)) {
			return true
		}
	}

	return false
}
