package models

import (
	"fmt"
	"sync"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
	Cart    map[*Item]int
	//Cart map[int]int
	mu sync.Mutex
}

//cart: make([]map[*Item]int, 0)
//Cart: make(map[int]int)

func AccountInit(id int, name string, balance float64, market *Market) *Account {
	account := Account{ID: id, Name: name, Balance: balance, Cart: make(map[*Item]int, 0)}
	market.addAccount(&account)
	return &account
}

func (account *Account) AddToCart(item *Item, quantity int) {
	account.Cart[item] += quantity
}

func (account *Account) ListTheCart() {
	fmt.Printf("items in cart for user: %s\n", account.Name)
	for item, quantity := range account.Cart {
		fmt.Printf("	Item ID: %d, Item Name: %s, Quantity: %d, Price per item: %.2f, Total Price: %.2f\n", item.ID, item.Name, quantity, item.Price, item.Price*float64(quantity))
	}
}

func (account *Account) CompletePurchase(purchaseChannel chan PurchaseResult) {
	account.mu.Lock()
	defer account.mu.Unlock()
	cost := 0.0
	for item, quantity := range account.Cart {
		item.mu.Lock()
		if quantity <= item.Stock {
			cost += (item.Price * float64(quantity))

		} else {
			text := fmt.Sprintf("only left %d from item %s, you requested %d", item.Stock, item.Name, quantity)
			purchaseChannel <- PurchaseResult{account.ID, account.Name, false, nil, 0, text}
			return
		}
	}
	if cost <= account.Balance {
		account.Balance -= cost
		ItemNameAndQuantity := make(map[string]int)
		for item, quantity := range account.Cart {
			item.Stock -= quantity
			ItemNameAndQuantity[item.Name] = quantity
		}
		text := "Your order is preparing"
		purchaseChannel <- PurchaseResult{account.ID, account.Name, true, ItemNameAndQuantity, cost, text}
	} else {
		text := fmt.Sprintf("insufficient balance; your balance is %.2f, order amount is %.2f", account.Balance, cost)
		purchaseChannel <- PurchaseResult{account.ID, account.Name, false, nil, 0, text}
	}
	for item := range account.Cart {
		item.mu.Unlock()
	}
}
