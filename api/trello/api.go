package trello

import(
  "net/url"
  "net/http"
  "encoding/json"
  "io/ioutil"
)

const API_KEY = "[API_KEY]"

type trelloClient struct {
  apiKey string
  authKey string
}

func NewClient(authKey string) trelloClient {
  return trelloClient {
    apiKey: API_KEY,
    authKey: authKey,
  }
}

var LinkForAuthKey = "https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&name=Treit&key=" + API_KEY

func (client trelloClient) CreateNewBoard(boardname string) (board Board) {
  response, err := http.Post("https://api.trello.com/1/boards?name=" + boardname + "&key=" + client.apiKey + "&token=" + client.authKey, "", nil)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    // TODO: Handle error
  }
  board = Board{}
  json.Unmarshal(body, &board)
  return
}

func (client trelloClient) GetBoardsByUsername(username string) (boards []Board) {
  response, err := http.Get("https://api.trello.com/1/members/" + username + "/boards?key=" + client.apiKey + "&token=" + client.authKey)
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

func (client trelloClient) GetListsFromBoardId(boardId string) (lists []List) {
  response, err := http.Get("https://api.trello.com/1/boards/" + boardId + "/lists?cards=all&card_fields=id,name,shortUrl&key=" + client.apiKey + "&token=" + client.authKey)
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

func (client trelloClient) MoveCardToList(cardId string, listId string) {
  httpClient := &http.Client{}
  req, err := http.NewRequest("PUT", "https://api.trello.com/1/cards/" + cardId + "?idList=" + listId + "&key=" + client.apiKey + "&token=" + client.authKey, nil)
  response, err := httpClient.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func (client trelloClient) DeleteCard(cardId string) {
  httpClient := &http.Client{}
  req, err := http.NewRequest("DELETE", "https://api.trello.com/1/cards/" + cardId + "?key=" + client.apiKey + "&token=" + client.authKey, nil)
  response, err := httpClient.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func (client trelloClient) AddCard(cardName string, listId string) {
  response, err := http.Post("https://api.trello.com/1/cards?name=" + url.QueryEscape(cardName) + "&idList=" + listId + "&key=" + client.apiKey + "&token=" + client.authKey, "", nil)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
}
