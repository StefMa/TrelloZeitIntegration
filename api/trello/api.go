package trello

import(
  "net/url"
  "net/http"
  "encoding/json"
  "io/ioutil"
)

const API_KEY = "[API_KEY]"

var LinkForAuthKey = "https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&name=Treit&key=" + API_KEY

func CreateNewBoard(authKey string, boardname string) (board Board) {
  response, err := http.Post("https://api.trello.com/1/boards?name=" + boardname + "&key=" + API_KEY + "&token=" + authKey, "", nil)
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

func GetBoardsByUsername(authKey string, username string) (boards []Board) {
  response, err := http.Get("https://api.trello.com/1/members/" + username + "/boards?key=" + API_KEY + "&token=" + authKey)
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

func GetListsFromBoardId(authKey string, boardId string) (lists []List) {
  response, err := http.Get("https://api.trello.com/1/boards/" + boardId + "/lists?cards=all&card_fields=id,name,shortUrl&key=" + API_KEY + "&token=" + authKey)
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

func MoveCardToList(authKey string, cardId string, listId string) {
  client := &http.Client{}
  req, err := http.NewRequest("PUT", "https://api.trello.com/1/cards/" + cardId + "?idList=" + listId + "&key=" + API_KEY + "&token=" + authKey, nil)
  response, err := client.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func DeleteCard(authKey string, cardId string) {
  client := &http.Client{}
  req, err := http.NewRequest("DELETE", "https://api.trello.com/1/cards/" + cardId + "?key=" + API_KEY + "&token=" + authKey, nil)
  response, err := client.Do(req)
  defer response.Body.Close()
  if err != nil {
    // TODO: Handle error
  }
}

func AddCard(authKey string, cardName string, listId string) {
  response, err := http.Post("https://api.trello.com/1/cards?name=" + url.QueryEscape(cardName) + "&idList=" + listId + "&key=" + API_KEY + "&token=" + authKey, "", nil)
  if err != nil {
    // TODO: Handle error
  }
  defer response.Body.Close()
}
