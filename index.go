package treit

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strings"
  "treit/api/trello"
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
    fmt.Fprint(w, buildOutputForSetup())
    return
  }

  if (action == ACTION_VIEW && metadata.AuthKey != "") {
    boards := trello.GetBoardsByUsername(metadata.AuthKey, metadata.Username)
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
    boards := trello.GetBoardsByUsername(authKey, username)
    fmt.Fprint(w, buildOutputForTrelloBoards(boards))
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
    fmt.Fprint(w, buildOutputForTrelloLists(lists, boardId))
    return
  }

  if (action == ACTION_MOVE_CARD_TO_LIST && metadata.AuthKey != "") {
    boardCardListId, _ := clientState[CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardCardListId, "_")
    trello.MoveCardToList(metadata.AuthKey, ids[1], ids[2])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[0])
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[0]))
    return
  }

  if (strings.HasPrefix(action, ACTION_DELETE_CARD) && metadata.AuthKey != "") {
    ids := strings.Split(action, "_")
    trello.DeleteCard(metadata.AuthKey, ids[2])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[1])
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[1]))
    return
  }

  if (action == ACTION_ADD_CARD && metadata.AuthKey != "") {
    boardListId, _ := clientState[CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID].(string)
    ids := strings.Split(boardListId, "_")
    cardName, _ := clientState[CLIENT_STATE_ADD_CARD_NAME].(string)
    trello.AddCard(metadata.AuthKey, cardName, ids[1])
    lists := trello.GetListsFromBoardId(metadata.AuthKey, ids[0])
    fmt.Fprint(w, buildOutputForTrelloLists(lists, ids[0]))
    return
  }
}

func buildOutputForSetup() (output string) {
  linkForAuth := trello.LinkForAuthKey

  output = "<Page>"
  output += "<H2>Setup</H2>"
  output += "<Container><Input name=\"" + CLIENT_STATE_AUTH_KEY + "\" label=\"Auth Key\"/><Link href=\"" + linkForAuth + "\" target=\"_blank\">Get Auth key</Link></Container>"
  output += "<Container><Input name=\"" + CLIENT_STATE_TRELLO_USERNAME + "\" label=\"Username\"/></Container>"
  output += "<Button action=\"" + ACTION_FINISH_SETUP + "\">Finish setup</Button>"
  output += "</Page>"
  return
}

func buildOutputForTrelloBoards(boards []trello.Board) (output string) {
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

func buildOutputForTrelloLists(lists []trello.List, boardId string) (output string) {
  listNames := []string{}
  listIds := []string{}
  for _, list := range lists {
    listNames = append(listNames, list.Name)
    listIds = append(listIds, list.Id)
  }

  output = "<Page>"
  output += "<Input name=\"" + CLIENT_STATE_ADD_CARD_NAME + "\" label=\"Add card\"/>"
  output += "<Select name=\"" + CLIENT_STATE_ADD_CARD_ID_IN_LIST_ID + "\" action=\"" + ACTION_ADD_CARD + "\">"
  output += "<Option selected disabled caption=\"Add card to...\"/>"
  for i, name := range listNames {
    output += "<Option value=\"" + boardId + "_" + listIds[i] + "\" caption=\"" + name + "\" />"
  }
  output += "</Select>"
  output += "<BR/><BR/>"

  for _, list := range lists {
    output += "<H2>" + list.Name + "</H2>"
    output += "<Box display=\"flex\" flex-wrap=\"wrap\" justifyContent=\"flex-start\" alignItems=\"center\">"
    for _, card := range list.Cards {
      output += "<Box display=\"flex\" flex-direction=\"column\" padding=\"12px\" margin=\"3px 3px 3px 0px\" backgroundColor=\"#FFF\" border=\"1px solid #f5f5f5\" borderRadius=\"3px\" >"
      output += "<Box display=\"flex\" alignItems=\"center\">"
      output += "<Link href=\"" + card.ShortUrl + "\">" + card.Name + "</Link>"
      output += "<Link action=\"" + ACTION_DELETE_CARD + "_" + boardId + "_" + card.Id + "\">❌</Link>"
      output += "</Box>"
      output += "<Select name=\"" + CLIENT_STATE_UPDATE_CARD_ID_IN_LIST_ID + "\" action=\"" + ACTION_MOVE_CARD_TO_LIST + "\">"
      output += "<Option selected disabled caption=\"Move card to...\"/>"
      for i, name := range listNames {
        if name != list.Name {
          output += "<Option value=\"" + boardId + "_" + card.Id + "_" + listIds[i] + "\" caption=\"" + name + "\" />"
        }
      }
      output += "</Select>"
      output += "</Box>"
    }
    output += "</Box>"
    output += "<BR/>"
  }
  output += "</Page>"
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
