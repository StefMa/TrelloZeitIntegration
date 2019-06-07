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

func GenerateForSetup(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (tmpl *template.Template, model Setup) {
  tmpl = template.Must(template.ParseFiles("web/template/setup.html"))
  model = createSetupModel(clientStateAuthKey, linkForAuthKey, clientStateTrelloUserName, actionFinishSetup)
  return
}

func createSetupModel(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (model Setup) {
  model = Setup {
    ClientStateAuthKey: clientStateAuthKey,
    LinkForAuthKey: linkForAuthKey,
    ClientStateTrelloUsername: clientStateTrelloUserName,
    ActionFinishSetup: actionFinishSetup,
  }
  return
}

func GenerateForTrelloBoards(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (tmpl *template.Template, model TrelloBoards) {
  tmpl = template.Must(template.ParseFiles("web/template/trelloBoards.html"))
  model = createTrelloBoardsModel(clientStateUseTrelloBoardId, actionUseTrelloBoard, boards, clientStateTrelloBoardName)
  return
}

func createTrelloBoardsModel(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (model TrelloBoards) {
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

func GenerateForTrelloLists(clientStateAddCardName string, clientStateAddCardIdInListId string, actionAddCard string, lists []trello.List, actionDeleteCard string, boardId string, clientStateUpdateCardInListId string, actionMoveCardToList string) (tmpl *template.Template, model TrelloLists) {
  tmpl = template.Must(template.ParseFiles("web/template/trelloLists.html"))
  model = createTrelloListsModel(clientStateAddCardName, clientStateAddCardIdInListId, actionAddCard, lists, actionDeleteCard, boardId, clientStateUpdateCardInListId, actionMoveCardToList)
  return
}

func createTrelloListsModel(clientStateAddCardName string, clientStateAddCardIdInListId string, actionAddCard string, lists []trello.List, actionDeleteCard string, boardId string, clientStateUpdateCardInListId string, actionMoveCardToList string) (model TrelloLists) {
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
