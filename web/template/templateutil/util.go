package templateutil

import(
  "html/template"
  "treit/api/trello"
)

type Setup struct {
  ClientStateAuthKey string
  LinkForAuthKey string
  ClientStateTrelloUsername string
  ActionFinishSetup string
}

type TrelloBoards struct {
  ClientStateUseTrelloBoardId string
  ActionUseTrelloBoard string
  Boards []TrelloBoard
  ClientStateTrelloBoardName string
}

type TrelloBoard struct {
  BoardId string
  BoardName string
}

type TrelloLists struct {
  ClientStateAddCardName string
  ClientStateAddCardIdInListId string
  ActionAddCard string
  Lists []TrelloList
  ActionDeleteCard string
  BoardId string
  ClientStateUpdateCardInListId string
  ActionMoveCardToList string
}

type TrelloList struct {
  ListId string
  ListName string
  ListCards []TrelloCard
}

type TrelloCard struct {
  CardId string
  CardName string
  CardShortUrl string
}

// FIXME: Currently this function crash on now.
// See https://github.com/zeit/now-builders/issues/564
func GenerateForSetup(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (tmpl *template.Template, model Setup) {
  tmpl = template.Must(template.ParseFiles("../setup.html"))
  model = CreateSetupModel(clientStateAuthKey, linkForAuthKey, clientStateTrelloUserName, actionFinishSetup)
  return
}

// TODO: Inline me in GenerateForSetup
// as soon the issue is fixed.
// Currently its only available
// to "don't break the API" from the `index.go`
func CreateSetupModel(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (model Setup) {
  model = Setup {
    ClientStateAuthKey: clientStateAuthKey,
    LinkForAuthKey: linkForAuthKey,
    ClientStateTrelloUsername: clientStateTrelloUserName,
    ActionFinishSetup: actionFinishSetup,
  }
  return
}

// FIXME: Currently this function crash on now.
// See https://github.com/zeit/now-builders/issues/564
func GenerateForTrelloBoards(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (tmpl *template.Template, model TrelloBoards) {
  tmpl = template.Must(template.ParseFiles("../trelloBoards.html"))
  model = CreateTrelloBoardsModel(clientStateUseTrelloBoardId, actionUseTrelloBoard, boards, clientStateTrelloBoardName)
  return
}

// TODO: Inline me in GenerateForTrelloBoards
// as soon the issue is fixed.
// Currently its only available
// to "don't break the API" from the `index.go`
func CreateTrelloBoardsModel(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (model TrelloBoards) {
  templateBoards := []TrelloBoard{}
  for _, board := range boards {
    templateBoards = append(templateBoards , TrelloBoard{board.Id, board.Name})
  }
  model = TrelloBoards {
    ClientStateUseTrelloBoardId: clientStateUseTrelloBoardId,
    ActionUseTrelloBoard: actionUseTrelloBoard,
    Boards: templateBoards,
    ClientStateTrelloBoardName: clientStateTrelloBoardName,
  }
  return
}

// FIXME: Currently this function crash on now.
// See https://github.com/zeit/now-builders/issues/564
func GenerateForTrelloLists(clientStateAddCardName string, clientStateAddCardIdInListId string, actionAddCard string, lists []trello.List, actionDeleteCard string, boardId string, clientStateUpdateCardInListId string, actionMoveCardToList string) (tmpl *template.Template, model TrelloLists) {
  tmpl = template.Must(template.ParseFiles("../trelloLists.html"))
  model = CreateTrelloListsModel(clientStateAddCardName, clientStateAddCardIdInListId, actionAddCard, lists, actionDeleteCard, boardId, clientStateUpdateCardInListId, actionMoveCardToList)
  return
}

// TODO: Inline me in GenerateForTrelloLists
// as soon the issue is fixed.
// Currently its only available
// to "don't break the API" from the `index.go`
func CreateTrelloListsModel(clientStateAddCardName string, clientStateAddCardIdInListId string, actionAddCard string, lists []trello.List, actionDeleteCard string, boardId string, clientStateUpdateCardInListId string, actionMoveCardToList string) (model TrelloLists) {
  templateList := []TrelloList{}
  for _, list := range lists {
    templateCard := []TrelloCard{}
    for _, card := range list.Cards {
      templateCard = append(templateCard, TrelloCard{card.Id, card.Name, card.ShortUrl})
    }
    templateList = append(templateList , TrelloList{list.Id, list.Name, templateCard})
  }
  model = TrelloLists {
    ClientStateAddCardName: clientStateAddCardName,
    ClientStateAddCardIdInListId: clientStateAddCardIdInListId,
    ActionAddCard: actionAddCard,
    Lists: templateList,
    ActionDeleteCard: actionDeleteCard,
    BoardId: boardId,
    ClientStateUpdateCardInListId: clientStateUpdateCardInListId,
    ActionMoveCardToList: actionMoveCardToList,
  }
  return
}
