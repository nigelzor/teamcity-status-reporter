package main

import (
	"cmd/browser"
	"cmd/go/auth"
	"encoding/xml"
	"fmt"
	"github.com/caseymrm/menuet"
	"io/ioutil"
	"log"
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

func status(projects Projects) (status string) {
	status = "✔️"
	failed := 0
	fixing := 0
	running := 0
	for _, project := range projects.Projects {
		if project.Activity != "Sleeping" {
			running += 1
		}
		if project.LastBuildStatus == "Unknown" {
			continue
		}
		if project.LastBuildStatus != "Success" {
			fmt.Println(project)
			failed += 1
			if project.Activity != "Sleeping" {
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

func interesting(projects Projects) (res []menuet.MenuItem) {
	for _, project := range projects.Projects {
		if project.LastBuildStatus == "Unknown" {
			continue
		}
		if project.LastBuildStatus != "Success" {
			res = append(res, menuet.MenuItem{
				Text: project.Name + ": " + project.LastBuildStatus,
				Clicked: func() {
					browser.Open(project.WebUrl)
				},
			})
		}
	}
	return
}

func update() {
	for {
		icon := "❕"
		projects, err := fetchProjects()
		if err == nil {
			icon = status(projects)
			items = interesting(projects)
		} else {
			log.Println(err)
			items = append(items, menuet.MenuItem{
				Text: err.Error(),
			})
		}
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: icon,
		})
		time.Sleep(30 * time.Second)
	}
}

func menu() []menuet.MenuItem {
	return append(items, menuet.MenuItem{
		Type: menuet.Separator,
	}, menuet.MenuItem{
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

type Projects struct {
	Projects []Project `xml:"Project"`
}

func main() {
	go update()
	menuet.App().Label = "com.github.nigelzor.teamcity-status-reporter"
	menuet.App().Children = menu
	menuet.App().RunApplication()
}
