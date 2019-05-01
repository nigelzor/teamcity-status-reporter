package main

import (
	"cmd/browser"
	"cmd/go/auth"
	"encoding/xml"
	"fmt"
	"github.com/caseymrm/menuet"
	"io/ioutil"
	"net/http"
	"time"
)

var items []menuet.MenuItem
var host string = "https://teamcity.newhippo.com"

func fetchProjects() (projects Projects, err error) {
	req, err := http.NewRequest("GET", host+"/app/rest/cctray/projects.xml", nil)
	if err != nil {
		return
	}
	auth.AddCredentials(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed to fetch projects: %v", resp.Status)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	xml.Unmarshal(body, &projects)
	return
}

func status(projects []Project) (status string) {
	status = "✔️"
	failed := 0
	fixing := 0
	running := 0
	for _, project := range projects {
		if project.Running() {
			running += 1
		}
		if project.Failed() {
			failed += 1
			if project.Running() {
				fixing += 1
			}
		}
	}
	if running > 0 {
		status = "➰"
	}
	if failed == 1 {
		status = "❗️"
		if fixing > 0 {
			status = "❓"
		}
	}
	if failed > 1 {
		status = "‼️"
		if fixing > 0 {
			status = "⁉️"
		}
	}
	return
}

func asMenuItem(p *Project) menuet.MenuItem {
	var prefix string
	if p.Running() && p.Failed() {
		prefix = "❓"
	} else if p.Running() {
		prefix = "➰"
	} else if p.Failed() {
		prefix = "❗️"
	}
	text := prefix + p.Name
	url := p.WebUrl
	return menuet.MenuItem{
		Text: text,
		Clicked: func() {
			browser.Open(url)
		},
	}
}

func interesting(projects []Project) (res []menuet.MenuItem) {
	for _, project := range projects {
		if project.Running() || project.Failed() {
			res = append(res, asMenuItem(&project))
		}
	}
	return
}

func update() {
	delay := 30 * time.Second
	for {
		icon := "❕"
		projects, err := fetchProjects()
		if err == nil {
			delay = 30 * time.Second
			icon = status(projects.Projects)
			items = interesting(projects.Projects)
		} else {
			if delay < time.Hour {
				delay *= 2
			}
			items = []menuet.MenuItem{menuet.MenuItem{
				Text: err.Error(),
			}}
		}
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: icon,
		})
		time.Sleep(delay)
	}
}

func menu() []menuet.MenuItem {
	return append(items,
		menuet.MenuItem{
			Type: menuet.Separator,
		},
		menuet.MenuItem{
			Text: "Open TeamCity",
			Clicked: func() {
				browser.Open(host)
			},
		})
}

type Project struct {
	Activity        string `xml:"activity,attr"`
	LastBuildLabel  string `xml:"lastBuildLabel,attr"`
	LastBuildStatus string `xml:"lastBuildStatus,attr"`
	LastBuildTime   string `xml:"lastBuildTime,attr"`
	Name            string `xml:"name,attr"`
	WebUrl          string `xml:"webUrl,attr"`
}

func (p *Project) Running() bool {
	return p.Activity != "Sleeping"
}

func (p *Project) Failed() bool {
	return p.LastBuildStatus != "Unknown" && p.LastBuildStatus != "Success"
}

type Projects struct {
	Projects []Project `xml:"Project"`
}

func main() {
	go update()
	menuet.App().Label = "com.github.nigelzor.teamcity-status-reporter"
	menuet.App().Children = menu
	menuet.App().RunApplication()
}
