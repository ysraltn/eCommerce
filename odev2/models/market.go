package models

import "fmt"

type Market struct {
	Accounts []*Account
	Items    []*Item
}

func MarketInit() *Market {
	return &Market{make([]*Account, 0), make([]*Item, 0)}
}

func (market *Market) addAccount(account *Account) {
	market.Accounts = append(market.Accounts, account)
}

func (market *Market) addItem(item *Item) {
	market.Items = append(market.Items, item)
}

func (market *Market) ListAccounts() {
	fmt.Println("All registered accounts are:")
	for _, account := range market.Accounts {
		fmt.Printf("	ID: %d, Name: %s, Balance: %.2f\n", account.ID, account.Name, account.Balance)
	}
}

func (market *Market) ListItems() {
	fmt.Println("All registered items are:")
	for _, item := range market.Items {
		fmt.Printf("	ID: %d, Name: %s, Stock: %d, Price: %.2f\n", item.ID, item.Name, item.Stock, item.Price)
	}
}

