package trello

type Board struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type List struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type Card struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ShortUrl string `json:"shortUrl"`
}
