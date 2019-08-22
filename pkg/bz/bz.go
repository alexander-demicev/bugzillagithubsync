package bz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alexander-demichev/bugzillagithubsync/config"
)

// Bugs bug list
type Bugs struct {
	ID      int    `json:"id"`
	Summary string `json:"summary"`
}

// BugList bug list
type BugList struct {
	Bugs []Bugs `json:"bugs"`
}

// RequestBugs request bugs
func RequestBugs() ([]byte, error) {
	apiKey := config.Config().GetString("API_KEY")
	link := "https://bugzilla.redhat.com/rest/bug?component=Migration%20Tooling&list_id=10428565&product=OpenShift%20Container%20Platform&api_key=" + apiKey

	resp, err := http.Get(link)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ParseBugs parse bugs
func ParseBugs(bugs []byte) (*BugList, error) {
	bugList := &BugList{}

	err := json.Unmarshal(bugs, bugList)
	if err != nil {
		return nil, err
	}

	return bugList, nil
}

// SelectCPMABugs select CPMA bugs
func SelectCPMABugs(bugList *BugList) *BugList {
	selectedBugs := &BugList{}

	for _, bug := range bugList.Bugs {
		if strings.Contains(bug.Summary, "CPMA") || strings.Contains(bug.Summary, "cpma") {
			selectedBugs.Bugs = append(selectedBugs.Bugs, bug)
		}
	}

	return selectedBugs
}
