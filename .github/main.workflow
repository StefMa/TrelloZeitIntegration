workflow "On pull request run go fmt" {
  on = "pull_request"
  resolves = ["run go fmt"]
}

action "run go fmt" {
  uses = "docker://golang:1.12.5-alpine3.9"
  args = "go fmt"
}
