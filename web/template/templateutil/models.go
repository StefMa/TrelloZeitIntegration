package templateutil

type Setup struct {
	ClientStateAuthKey        string
	LinkForAuthKey            string
	ClientStateTrelloUsername string
	ActionFinishSetup         string
}

type TrelloBoards struct {
	ClientStateUseTrelloBoardId string
	ActionUseTrelloBoard        string
	Boards                      []TrelloBoard
	ClientStateTrelloBoardName  string
}

type TrelloBoard struct {
	BoardId   string
	BoardName string
}

type TrelloLists struct {
	ClientStateAddCardName        string
	ClientStateAddCardIdInListId  string
	ActionAddCard                 string
	Lists                         []TrelloList
	ActionDeleteCard              string
	BoardId                       string
	ClientStateUpdateCardInListId string
	ActionMoveCardToList          string
}

type TrelloList struct {
	ListId    string
	ListName  string
	ListCards []TrelloCard
}

type TrelloCard struct {
	CardId       string
	CardName     string
	CardShortUrl string
}
