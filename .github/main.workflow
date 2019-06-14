workflow "On pull request run go fmt" {
  on = "pull_request"
  resolves = ["run gofmt"]
}

action "run gofmt" {
  uses = "docker://golang:1.12.5-alpine3.9"
  args = "gofmt -d ."
}
