package treit

import(
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strings"
  "treit/api/trello"
  "treit/web/template/templateutil"
)

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
    tmpl, model := templateutil.GenerateForSetup(CLIENT_STATE_AUTH_KEY, trello.LinkForAuthKey, CLIENT_STATE_TRELLO_USERNAME, ACTION_FINISH_SETUP)
    tmpl.Execute(w, model)
    return
  }

  if (action == ACTION_VIEW && metadata.AuthKey != "") {
    boards := trello.GetBoardsByUsername(metadata.AuthKey, metadata.Username)
    tmpl, model := templateutil.GenerateForTrelloBoards(CLIENT_STATE_USE_TRELLO_BOARD_ID, ACTION_USE_TRELLO_BOARD, boards, CLIENT_STATE_TRELLO_BOARD_NAME)
    tmpl.Execute(w, model)
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
    boards := trello.GetBoardsByUsername(authKey, username)
    tmpl, model := templateutil.GenerateForTrelloBoards(CLIENT_STATE_USE_TRELLO_BOARD_ID, ACTION_USE_TRELLO_BOARD, boards, CLIENT_STATE_TRELLO_BOARD_NAME)
    tmpl.Execute(w, model)
    return
  }

  if (strings.HasPrefix(action, ACTION_USE_TRELLO_BOARD) && metadata.AuthKey != "") {
    boardId := strings.TrimPrefix(action, ACTION_USE_TRELLO_BOARD)
    if (boardId == "NEW") {
      // This is just a check if this is really a string...
      // see https://stackoverflow.com/a/14289568
      boardName, _ := clientState[CLIENT_STATE_TRELLO_BOARD_NAME].(string)
      board := trello.CreateNewBoard(metadata.AuthKey, boardName)
      boardId = board.Id
    } else {
      // This is just a check if this is really a string...
      // see https://stackoverflow.com/a/14289568
      boardId, _ = clientState[CLIENT_STATE_USE_TRELLO_BOARD_ID].(string)
    }
    lists := trello.GetListsFromBoardId(metadata.AuthKey, boardId)
    tmpl, model := templateutil.GenerateForTrelloLists(CLIENT_STATE_ADD_CARD_NAME, CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID, ACTION_ADD_CARD, lists, ACTION_DELETE_CARD, boardId, CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID, ACTION_MOVE_CARD_TO_LIST)
    tmpl.Execute(w, model)
    return
  }

  if (action == ACTION_MOVE_CARD_TO_LIST && metadata.AuthKey != "") {
    boardCardListId, _ := clientState[CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardCardListId, "_")
    trello.MoveCardToList(metadata.AuthKey, ids[1], ids[2])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[0])
    tmpl, model := templateutil.GenerateForTrelloLists(CLIENT_STATE_ADD_CARD_NAME, CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID, ACTION_ADD_CARD, lists, ACTION_DELETE_CARD, ids[0], CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID, ACTION_MOVE_CARD_TO_LIST)
    tmpl.Execute(w, model)
    return
  }

  if (strings.HasPrefix(action, ACTION_DELETE_CARD) && metadata.AuthKey != "") {
    ids := strings.Split(action, "_")
    trello.DeleteCard(metadata.AuthKey, ids[2])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[1])
    tmpl, model := templateutil.GenerateForTrelloLists(CLIENT_STATE_ADD_CARD_NAME, CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID, ACTION_ADD_CARD, lists, ACTION_DELETE_CARD, ids[1], CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID, ACTION_MOVE_CARD_TO_LIST)
    tmpl.Execute(w, model)
    return
  }

  if (action == ACTION_ADD_CARD && metadata.AuthKey != "") {
    boardListId, _ := clientState[CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardListId, "_")
    cardName, _ := clientState[CLIENT_STATE_ADD_CARD_NAME].(string)
    trello.AddCard(metadata.AuthKey, cardName, ids[1])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[0])
    tmpl, model := templateutil.GenerateForTrelloLists(CLIENT_STATE_ADD_CARD_NAME, CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID, ACTION_ADD_CARD, lists, ACTION_DELETE_CARD, ids[0], CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID, ACTION_MOVE_CARD_TO_LIST)
    tmpl.Execute(w, model)
    return
  }
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
