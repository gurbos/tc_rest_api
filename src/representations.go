package main

import tcm "github.com/gurbos/tcmodels"

// Constants representing resource identifiers.
const (
	ProductLinesURL        = "http://127.0.0.1:8000/productLine"                                 // Collection of product line info representations.
	ProductLineCardSetsURL = "http://127.0.0.1:8000/{productLineName}/sets"                      // Collection of set info representations.
	ProductLineCardsURL    = "http://127.0.0.1:8000/{productLineName}/cards?{setName,from,size}" // Collection of card representations from a specified product line.
	CardsFromSetURL        = "http://	127.0.0.1:8000/{productLineName}/"
)

type TradingCardAPI struct {
	ProductLinesUrl        string `json:"product_lines_url"`
	ProductLineCardSetsURL string `json:"product_line_card_sets_url"`
	ProductLineCardsURL    string `json:"product_line_cards_url"`
	CardsFromSetURL        string `json:"cards_from_set_url"`
}

func (tca *TradingCardAPI) Init() {
	tca.ProductLinesUrl = ProductLinesURL
	tca.ProductLineCardSetsURL = ProductLineCardSetsURL
	tca.ProductLineCardsURL = ProductLineCardsURL
	tca.CardsFromSetURL = CardsFromSetURL
}

type ProductLineRep struct {
	Title     string `json:"title"`
	QueryName string `json:"queryName"`
	SetCount  uint   `json:"setCount"`
	CardCount uint   `json:"cardCount"`
	Sets      Link   `json:"sets"`
	Cards     Link   `json:"cards"`
}

func (plr *ProductLineRep) Set(pl tcm.ProductLine) {
	plr.Title = pl.Title
	plr.QueryName = pl.Name
	plr.CardCount = pl.CardCount
	plr.Sets.Rel = "collection/sets"
	plr.Sets.Href = "/" + pl.Name + "/sets"
	plr.Cards.Rel = "collection/cards"
	plr.Cards.Href = "/" + pl.Name + "/cards"
}

type CardSetRep struct {
	Title            string `json:"title"`
	Name             string `json:"name"`
	CardCount        uint   `json:"cardCount"`
	ProductLineTitle string `json:"productLineTitle"`
	ProductLineName  string `json:"productLineName"`
	Cards            Link   `json:"cards"`
}

func (csr *CardSetRep) Set(setInfo tcm.SetInfo) {
	csr.Title = setInfo.Title
	csr.Name = setInfo.Name
	csr.CardCount = setInfo.CardCount
	csr.ProductLineTitle = setInfo.ProductLine.Title
	csr.ProductLineName = setInfo.ProductLine.Name
	csr.Cards.Rel = "collection/cards"
	csr.Cards.Href = "/" + csr.ProductLineName + "/" + csr.Name + "/cards"
}

type CardSetCollection struct {
	ProdutLineTitle string `json:"productLineTitle"`
	ProductLineName string `json:"productLineName"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
