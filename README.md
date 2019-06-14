# Trello Zeit Integration

**Warning:** This project is currently just a show case.
You can use as a base if you are familiar with how to create a
Zeit Integration and a little bit with the Trello API.

## How to build

* Replace the [`API_KEY`](api/trello/api.go) with your [**Trello API-KEY**](https://trello.com/app-key)
* Run `now`

Now you should be able to use the deployment as your Zeit Integration UIHook

## Development setup

### Automatic gofmt

I have added a [pre-commit](.githooks/pre-commit) to automatically
format the Go code for each changed file with `gofmt`.

To enable it simply run the following command:
```
git config core.hooksPath .githooks
```
