package main

import (
	"encoding/xml"
	"fmt"
	"time"
	"github.com/caseymrm/menuet"
	"io/ioutil"
	"log"
	"net/http"
	"cmd/go/auth"
)

func fetchProjects() (projects Projects, err error) {
	host := "https://teamcity.newhippo.com"
	req, err := http.NewRequest("GET", host + "/app/rest/cctray/projects.xml", nil)
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
	for _, project := range projects.Projects {
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

func update() {
	status := [...]string { "❕", "❗️", "❓", "‼️", "⁉️", "✔️" }
	i := 0
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: status[i] + " " + time.Now().Format(":05"),
		})
		i += 1
		if i >= len(status) {
			i = 0
		}
		time.Sleep(time.Second)
	}
}

type Projects struct {
	Projects []struct {
		Activity string `xml:"activity,attr"`
		LastBuildLabel string `xml:"lastBuildLabel,attr"`
		LastBuildStatus string `xml:"lastBuildStatus,attr"`
		LastBuildTime string `xml:"lastBuildTime,attr"`
		Name string `xml:"name,attr"`
		WebUrl string `xml:"webUrl,attr"`
	} `xml:"Project"`
}

func main() {
	projects, err := fetchProjects()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status(projects))
	return
	go update()
	menuet.App().Label = "com.github.nigelzor.teamcity-status-reporter"
	menuet.App().RunApplication()
}
