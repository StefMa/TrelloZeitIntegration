package templateutil

import (
	"github.com/stefma/TrelloZeitIntegration/api/trello"
	"html/template"
	"os"
)

func GenerateForSetup(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (tmpl *template.Template, model Setup) {
	tmpl = templateWithName("setup.html")
	model = createSetupModel(clientStateAuthKey, linkForAuthKey, clientStateTrelloUserName, actionFinishSetup)
	return
}

func createSetupModel(clientStateAuthKey string, linkForAuthKey string, clientStateTrelloUserName string, actionFinishSetup string) (model Setup) {
	model = Setup{
		ClientStateAuthKey:        clientStateAuthKey,
		LinkForAuthKey:            linkForAuthKey,
		ClientStateTrelloUsername: clientStateTrelloUserName,
		ActionFinishSetup:         actionFinishSetup,
	}
	return
}

func GenerateForTrelloBoards(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (tmpl *template.Template, model TrelloBoards) {
	tmpl = templateWithName("trelloBoards.html")
	model = createTrelloBoardsModel(clientStateUseTrelloBoardId, actionUseTrelloBoard, boards, clientStateTrelloBoardName)
	return
}

func createTrelloBoardsModel(clientStateUseTrelloBoardId string, actionUseTrelloBoard string, boards []trello.Board, clientStateTrelloBoardName string) (model TrelloBoards) {
	templateBoards := []TrelloBoard{}
	for _, board := range boards {
		templateBoards = append(templateBoards, TrelloBoard{board.Id, board.Name})
	}
	model = TrelloBoards{
		ClientStateUseTrelloBoardId: clientStateUseTrelloBoardId,
		ActionUseTrelloBoard:        actionUseTrelloBoard,
		Boards:                      templateBoards,
		ClientStateTrelloBoardName:  clientStateTrelloBoardName,
	}
	return
}

func GenerateForTrelloLists(clientStateAddCardName string, clientStateAddCardIdInListId string, actionAddCard string, lists []trello.List, actionDeleteCard string, boardId string, clientStateUpdateCardInListId string, actionMoveCardToList string) (tmpl *template.Template, model TrelloLists) {
	tmpl = templateWithName("trelloLists.html")
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
		templateList = append(templateList, TrelloList{list.Id, list.Name, templateCard})
	}
	model = TrelloLists{
		ClientStateAddCardName:        clientStateAddCardName,
		ClientStateAddCardIdInListId:  clientStateAddCardIdInListId,
		ActionAddCard:                 actionAddCard,
		Lists:                         templateList,
		ActionDeleteCard:              actionDeleteCard,
		BoardId:                       boardId,
		ClientStateUpdateCardInListId: clientStateUpdateCardInListId,
		ActionMoveCardToList:          actionMoveCardToList,
	}
	return
}

func templateWithName(name string) *template.Template {
	templateLocationDir := os.Getenv("TEMPLATES_LOCATION_DIR")
	return template.Must(template.ParseFiles(templateLocationDir + name))
}
