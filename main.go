package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strings"
)

const TRELLO_API_KEY = "[TRELLO_API_KEY]"
const TRELLO_USERNAME = "[TRELLO_USERNAME]"

type payload struct {
  Action string                       `json:"action"`
  ClientState map[string]interface{}  `json:"clientState"`
  ConfigurationId string              `json:"configurationId"`
  Token string                        `json:"token"`
}

const ACTION_VIEW = "view"
const ACTION_SET_AUTH_KEY = "setAuthKey"
const ACTION_USE_TRELLO_BOARD = "useTrelloBoard"

const CLIENT_STATE_AUTH_KEY = "authKey"
const CLIENT_STATE_TRELLO_BOARD_ID = "trelloBoardId"

type trelloBoard struct {
  Id string   `json:"id"`
  Name string `json:"name"`
}

type trelloList struct {
  Name string         `json:"name"`
  Cards []trelloCard  `json:"cards"`
}

type trelloCard struct {
  Name string `json:"name"`
}

type metadata struct {
  AuthKey string  `json:"authKey"`
}

func main() {
  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      // TODO: Handle error
    }

    payload := payload{}
    err = json.Unmarshal(body, &payload)
    if err != nil {
      // TODO: Handle error
    }
    action := payload.Action
    clientState := payload.ClientState
    configurationId := payload.ConfigurationId
    token := payload.Token
    metadata := getMetadata(configurationId, token)

    header := w.Header()
    header.Add("Access-Control-Allow-Origin", "*")
    header.Add("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    header.Add("Access-Control-Allow-Headers", "Authorization, Accept, Content-Type")

    if (action == ACTION_VIEW && metadata.AuthKey == "") {
      fmt.Fprint(w, buildOutputForViewWithoutKey())
      return
    }

    if (action == ACTION_VIEW && metadata.AuthKey != "") {
      // TODO: I need the username (or id :eyeroll:)
      boards := getTrelloBoardsByUsername(TRELLO_USERNAME, metadata.AuthKey)
      output := buildOutputForTrelloBoards(boards)
      fmt.Fprint(w, output)
      return
    }

    if (action == ACTION_SET_AUTH_KEY && metadata.AuthKey == "") {
      // This is just a check if this is really a string...
      // see https://stackoverflow.com/a/14289568
      str, _ := clientState[CLIENT_STATE_AUTH_KEY].(string)
      // Save the authKey
      var jsonStr = "{\"authKey\":\"" + str + "\"}"
      saveMetadata(configurationId, token, jsonStr)
      // Create the new board...
      //board := createNewTrelloBoard(authKey)
      // Return the HTML code
      fmt.Fprint(w, buildOutputForSavingAuthKey())
      return
    }

    if (strings.HasPrefix(action, ACTION_USE_TRELLO_BOARD) && metadata.AuthKey != "") {
      boardId := strings.TrimPrefix(action, ACTION_USE_TRELLO_BOARD)
      lists := getTrelloListsFromBoardId(boardId, metadata.AuthKey)
      fmt.Fprint(w, buildOutputForTrelloLists(lists))
    }
  })
  http.ListenAndServe(":80", nil)
}

func buildOutputForViewWithoutKey() (output string) {
  linkForAuth := "https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&name=Zeit%20Trello%20Integration&key=" + TRELLO_API_KEY

  output = "<Page>"
  output += "<Container><Input name=\"" + CLIENT_STATE_AUTH_KEY + "\" label=\"Auth Key\"/><Button action=\"" + ACTION_SET_AUTH_KEY + "\">Set Auth Key</Button></Container>"
  output += "<Container><Link href=\"" + linkForAuth + "\" target=\"_blank\">Get Auth key</Link></Container>"
  output += "</Page>"
  return
}

func buildOutputForSavingAuthKey() (output string) {
  output = "<Page>"
  output += "Received authKey. Redirect in 3 seconds."
  output += "<AutoRefresh timeout=\"3000\" />"
  output += "</Page>"
  return
}

func buildOutputForTrelloBoards(boards []trelloBoard) (output string) {
  output = "<Page>"
  output += "<H2>Your Trello boards:</H2>"
  for _, board := range boards {
    output += "<Link action=\"" + ACTION_USE_TRELLO_BOARD + board.Id + "\">" + board.Name + "</Link><BR/>"
  }
  output += "</Page>"
  return
}

func buildOutputForTrelloLists(lists []trelloList) (output string) {
  output = "<Page>"
  for _, list := range lists {
    output += "<H2>" + list.Name + "</H2>"
    for _, card := range list.Cards {
      output += card.Name + "<BR/>"
    }
  }
  output += "</Page>"
  return
}

func createNewTrelloBoard(authKey string) (board trelloBoard) {
  response, err := http.Post("https://api.trello.com/1/boards?name=TreDo&key=" + TRELLO_API_KEY + "&token=" + authKey, "", nil)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    // TODO: Handle error
  }
  board = trelloBoard{}
  json.Unmarshal(body, &board)
  return
}

func getMetadata(configurationId string, token string) (metadata metadata) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://api.zeit.co/v1/integrations/configuration/" + configurationId + "/metadata", nil)
  req.Header.Add("Authorization", "Bearer " + token)
  response, err := client.Do(req)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    // TODO: Handle error
  }
  json.Unmarshal(body, &metadata)
  return
}

func saveMetadata(configurationId string, token string, metadataJsonString string) {
  client := &http.Client{}
  metadataJsonByte := []byte(metadataJsonString)
  req, err := http.NewRequest("POST", "https://api.zeit.co/v1/integrations/configuration/" + configurationId + "/metadata", bytes.NewBuffer(metadataJsonByte))
  req.Header.Add("Authorization", "Bearer " + token)
  response, err := client.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func getTrelloBoardsByUsername(username string, authKey string) (boards []trelloBoard) {
  response, err := http.Get("https://api.trello.com/1/members/" + username + "/boards?key=" + TRELLO_API_KEY + "&token=" + authKey)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    // TODO: Handle error
  }
  json.Unmarshal(body, &boards)
  return
}


func getTrelloListsFromBoardId(boardId string, authKey string) (lists []trelloList) {
  response, err := http.Get("https://api.trello.com/1/boards/" + boardId + "/lists?cards=all&key=" + TRELLO_API_KEY + "&token=" + authKey)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    // TODO: Handle error
  }
  json.Unmarshal(body, &lists)
  return
}
