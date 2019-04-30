package main

import (
	"time"
	"github.com/caseymrm/menuet"
)

func update() {
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "Hello World " + time.Now().Format(":05"),
		})
		time.Sleep(time.Second)
	}
}

func main() {
	go update()
	menuet.App().Label = "com.github.nigelzor.teamcity-status-reporter"
	menuet.App().RunApplication()
}
