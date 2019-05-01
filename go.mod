module github.com/nigelzor/teamcity-status-reporter

go 1.12

require (
	github.com/caseymrm/askm v0.0.0-20180731222743-ac4324c25a4d // indirect
	github.com/caseymrm/menuet v0.0.0-20190225170317-be05f48c376e
)

require cmd/go/auth v0.0.0

replace cmd/go/auth v0.0.0 => ./auth

require cmd/browser v0.0.0

replace cmd/browser v0.0.0 => ./browser
