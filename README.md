teamcity-status-reporter
========================

it's a MacOS menu app that shows current running/failed build status for a TeamCity server

| Icon | Meaning                                |
|:----:| -------------------------------------- |
|  ❕  | Failed to fetch build status           |
|  ➰  | Building...                            |
|  ✔️   | Everything's ok!                       |
|  ❗️  | One build is failing                   |
|  ❓  | One failed build, rebuilding           | 
|  ‼️   | More than one failed build             |
|  ⁉️   | More than one failed build, rebuilding |

### setup:
 - `go build`
 - add credentials to your [.netrc](https://ec.haxx.se/usingcurl-netrc.html)
 - `./teamcity-status-reporter -host https://my.teamcity.server`
