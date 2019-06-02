package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strings"
  "net/url"
)

const TRELLO_API_KEY = "[TRELLO_API_KEY]"

type payload struct {
  Action string                       `json:"action"`
  ClientState map[string]interface{}  `json:"clientState"`
  ConfigurationId string              `json:"configurationId"`
  Token string                        `json:"token"`
}

const ACTION_VIEW = "view"
const ACTION_FINISH_SETUP = "finishSetup"
const ACTION_USE_TRELLO_BOARD = "useTrelloBoard"
const ACTION_MOVE_CARD_TO_LIST = "moveCardToList"
const ACTION_DELETE_CARD = "deleteCard"
const ACTION_ADD_CARD = "addCard"

const CLIENT_STATE_AUTH_KEY = "authKey"
const CLIENT_STATE_TRELLO_USERNAME = "trelloUsername"
const CLIENT_STATE_USE_TRELLO_BOARD_ID = "trelloBoardId"
const CLIENT_STATE_TRELLO_BOARD_NAME = "trelloBoardName"
const CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID = "updateCardIdInListId"
const CLIENT_STATE_ADD_CARD_NAME = "cardName"
const CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID = "addCardTolist"

type trelloBoard struct {
  Id string   `json:"id"`
  Name string `json:"name"`
}

type trelloList struct {
  Id string           `json:"id"`
  Name string         `json:"name"`
  Cards []trelloCard  `json:"cards"`
}

type trelloCard struct {
  Id string       `json:"id"`
  Name string     `json:"name"`
  ShortUrl string `json:"shortUrl"`
}

type metadata struct {
  AuthKey string  `json:"authKey"`
  Username string `json:"username"`
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
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
    fmt.Fprint(w, buildOutputForSetup())
    return
  }

  if (action == ACTION_VIEW && metadata.AuthKey != "") {
    // TODO: I need the username (or id :eyeroll:)
    boards := getTrelloBoardsByUsername(metadata.Username, metadata.AuthKey)
    output := buildOutputForTrelloBoards(boards)
    fmt.Fprint(w, output)
    return
  }

  if (action == ACTION_FINISH_SETUP && metadata.AuthKey == "") {
    // This is just a check if this is really a string...
    // see https://stackoverflow.com/a/14289568
    authKey, _ := clientState[CLIENT_STATE_AUTH_KEY].(string)
    username, _ := clientState[CLIENT_STATE_TRELLO_USERNAME].(string)
    // Save the authKey
    var jsonStr = "{\"authKey\":\"" + authKey + "\", \"username\":\"" + username + "\"}"
    saveMetadata(configurationId, token, jsonStr)
    // Return the HTML code
    fmt.Fprint(w, buildOutputForSavingAuthKey())
    return
  }

  if (strings.HasPrefix(action, ACTION_USE_TRELLO_BOARD) && metadata.AuthKey != "") {
    boardId := strings.TrimPrefix(action, ACTION_USE_TRELLO_BOARD)
    if (boardId == "NEW") {
      // This is just a check if this is really a string...
      // see https://stackoverflow.com/a/14289568
      boardName, _ := clientState[CLIENT_STATE_TRELLO_BOARD_NAME].(string)
      board := createNewTrelloBoard(metadata.AuthKey, boardName)
      boardId = board.Id
    } else {
      // This is just a check if this is really a string...
      // see https://stackoverflow.com/a/14289568
      boardId, _ = clientState[CLIENT_STATE_USE_TRELLO_BOARD_ID].(string)
    }
    lists := getTrelloListsFromBoardId(boardId, metadata.AuthKey)
    fmt.Fprint(w, buildOutputForTrelloLists(lists, boardId))
    return
  }

  if (action == ACTION_MOVE_CARD_TO_LIST && metadata.AuthKey != "") {
    boardCardListId, _ := clientState[CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardCardListId, "_")
    moveCardToList(ids[1], ids[2], metadata.AuthKey)
    lists := getTrelloListsFromBoardId(ids[0], metadata.AuthKey)
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[0]))
    return
  }

  if (strings.HasPrefix(action, ACTION_DELETE_CARD) && metadata.AuthKey != "") {
    ids := strings.Split(action, "_")
    deleteCard(ids[2], metadata.AuthKey)
    lists := getTrelloListsFromBoardId(ids[1], metadata.AuthKey)
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[1]))
    return
  }

  if (action == ACTION_ADD_CARD && metadata.AuthKey != "") {
    boardListId, _ := clientState[CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardListId, "_")
    cardName, _ := clientState[CLIENT_STATE_ADD_CARD_NAME].(string)
    addCard(cardName, ids[1], metadata.AuthKey)
    lists := getTrelloListsFromBoardId(ids[0], metadata.AuthKey)
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[0]))
    return
  }
}

func buildOutputForSetup() (output string) {
  linkForAuth := "https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&name=Zeit%20Trello%20Integration&key=" + TRELLO_API_KEY

  output = "<Page>"
  output += "<H2>Setup</H2>"
  output += "<Container><Input name=\"" + CLIENT_STATE_AUTH_KEY + "\" label=\"Auth Key\"/><Link href=\"" + linkForAuth + "\" target=\"_blank\">Get Auth key</Link></Container>"
  output += "<Container><Input name=\"" + CLIENT_STATE_TRELLO_USERNAME + "\" label=\"Username\"/></Container>"
  output += "<Button action=\"" + ACTION_FINISH_SETUP + "\">Finish setup</Button>"
  output += "</Page>"
  return
}

func buildOutputForSavingAuthKey() (output string) {
  output = "<Page>"
  output += "Setup everything ✅<BR/>Redirect in 3 seconds..."
  output += "<AutoRefresh timeout=\"3000\" />"
  output += "</Page>"
  return
}

func buildOutputForTrelloBoards(boards []trelloBoard) (output string) {
  output = "<Page>"
  output += "<H2>Choose board</H2>"
  output += "<Select name=\"" + CLIENT_STATE_USE_TRELLO_BOARD_ID + "\" action=\"" + ACTION_USE_TRELLO_BOARD + "\">"
  output += "<Option selected disabled caption=\"Select board...\"/>"
  for _, board := range boards {
    output += "<Option value=\"" + board.Id + "\" caption=\"" + board.Name + "\" />"
  }
  output += "</Select>"
  output += "<BR/><BR/>"
  output += "Or create a new one"
  output += "<BR/>"
  output += "<Input name=\"" + CLIENT_STATE_TRELLO_BOARD_NAME + "\" value=\"Boardname\" />"
  output += "  "
  output += "<Button action=\"" + ACTION_USE_TRELLO_BOARD + "NEW\" secondary>Create</Button>"
  output += "</Page>"
  return
}

func buildOutputForTrelloLists(lists []trelloList, boardId string) (output string) {
  listNames := []string{}
  listIds := []string{}
  for _, list := range lists {
    listNames = append(listNames, list.Name)
    listIds = append(listIds, list.Id)
  }

  output = "<Page>"
  for _, list := range lists {
    output += "<H2>" + list.Name + "</H2>"
    for _, card := range list.Cards {
      output += "<Link href=\"" + card.ShortUrl + "\">" + card.Name + "</Link>"
      output += "<Link action=\"" + ACTION_DELETE_CARD + "_" + boardId + "_" + card.Id + "\">❌</Link>"
      output += "<Select name=\"" + CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID + "\" action=\"" + ACTION_MOVE_CARD_TO_LIST + "\">"
      output += "<Option selected disabled caption=\"Move card...\"/>"
      for i, name := range listNames {
        if name != list.Name {
          output += "<Option value=\"" + boardId + "_" + card.Id + "_" + listIds[i] + "\" caption=\"" + name + "\" />"
        }
      }
      output += "</Select>"
      output += "<BR/>"
    }
  }
  output += "<Input name=\"" + CLIENT_STATE_ADD_CARD_NAME + "\" label=\"Add card\"/>"
  output += "<Select name=\"" + CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID + "\" action=\"" + ACTION_ADD_CARD + "\">"
  output += "<Option selected disabled caption=\"Add card to...\"/>"
  for i, name := range listNames {
    output += "<Option value=\"" + boardId + "_" + listIds[i] + "\" caption=\"" + name + "\" />"
  }
  output += "</Select>"
  output += "</Page>"
  return
}

func createNewTrelloBoard(authKey string, boardname string) (board trelloBoard) {
  response, err := http.Post("https://api.trello.com/1/boards?name=" + boardname + "&key=" + TRELLO_API_KEY + "&token=" + authKey, "", nil)
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
  response, err := http.Get("https://api.trello.com/1/boards/" + boardId + "/lists?cards=all&card_fields=id,name,shortUrl&key=" + TRELLO_API_KEY + "&token=" + authKey)
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

func moveCardToList(cardId string, listId string, authKey string) {
  client := &http.Client{}
  req, err := http.NewRequest("PUT", "https://api.trello.com/1/cards/" + cardId + "?idList=" + listId + "&key=" + TRELLO_API_KEY + "&token=" + authKey, nil)
  response, err := client.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func deleteCard(cardId string, authKey string) {
  client := &http.Client{}
  req, err := http.NewRequest("DELETE", "https://api.trello.com/1/cards/" + cardId + "?key=" + TRELLO_API_KEY + "&token=" + authKey, nil)
  response, err := client.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func addCard(cardName string, listId string, authKey string) {
  response, err := http.Post("https://api.trello.com/1/cards?name=" + url.QueryEscape(cardName) + "&idList=" + listId + "&key=" + TRELLO_API_KEY + "&token=" + authKey, "", nil)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
}
